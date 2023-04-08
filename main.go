package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"os"
	"strings"
)

// VisitUrl 行政区划访问地址
var VisitUrl = "http://www.stats.gov.cn/sj/tjbz/tjyqhdmhcxhfdm/2022/"

//var VisitUrl = "http://www.stats.gov.cn/sj/tjbz/tjyqhdmhcxhfdm/2022/12/01/120101.html"

// Region 行政区划
type Region struct {
	bianMa         string // 编码
	mingCheng      string // 名称
	chengXiangFlDM string // 城乡分类代码
	level          int    // 级别
	fuJiBM         string // 父级编码
}

const (
	Sheng = iota
	ZhouShi
	QuXianShi
	XiangZheng
	CunSheQu
)

func main() {
	file, err := os.OpenFile("regions.txt", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		_ = file.Close()
	}()
	c := colly.NewCollector(
		colly.AllowedDomains("www.stats.gov.cn"),
	)
	c.OnHTML("tr.provincetr > td", func(e *colly.HTMLElement) {
		query(c, e, Sheng, file)
	})
	c.OnHTML("tr.citytr", func(e *colly.HTMLElement) {
		query(c, e, ZhouShi, file)
	})
	c.OnHTML("tr.countytr", func(e *colly.HTMLElement) {
		query(c, e, QuXianShi, file)
	})
	c.OnHTML("tr.towntr", func(e *colly.HTMLElement) {
		query(c, e, XiangZheng, file)
	})
	c.OnHTML("tr.villagetr", func(e *colly.HTMLElement) {
		query(c, e, CunSheQu, file)
	})
	// 修改传入的地址可以只获取某层的行政区划
	// eg: _ = c.Visit("http://www.stats.gov.cn/sj/tjbz/tjyqhdmhcxhfdm/2022/53.html") // 抽取云南
	// eg: _ = c.Visit("http://www.stats.gov.cn/sj/tjbz/tjyqhdmhcxhfdm/2022/53/5328.html") // 抽取西双版纳州
	_ = c.Visit(VisitUrl)
}

func query(c *colly.Collector, e *colly.HTMLElement, level int, file *os.File) {
	r := &Region{
		level: level,
	}
	href := ""
	switch level {
	case Sheng:
		e.ForEach("a", func(i int, element *colly.HTMLElement) {
			href = element.Attr("href")
			r.bianMa = strings.Split(href, ".")[0] + "0000000000"
			r.mingCheng = element.Text
		})
		break
	case QuXianShi:
		e.ForEach("td", func(i int, element *colly.HTMLElement) {
			switch i {
			case 0:
				// 编码
				r.bianMa = element.Text
				break
			case 1:
				// 名称
				r.mingCheng = element.Text
				href = r.bianMa[0:2] + "/" + element.ChildAttr("a", "href")
				break
			default:
			}
		})
		break
	case XiangZheng:
		e.ForEach("td > a", func(i int, element *colly.HTMLElement) {
			switch i {
			case 0:
				// 编码
				r.bianMa = element.Text
				break
			case 1:
				// 名称
				r.mingCheng = element.Text
				href = r.bianMa[0:2] + "/" + r.bianMa[2:4] + "/" + element.Attr("href")
				break
			default:
			}
		})
		break
	case CunSheQu:
		e.ForEach("td", func(i int, element *colly.HTMLElement) {
			switch i {
			case 0:
				// 编码
				r.bianMa = element.Text
				break
			case 1:
				// 城乡分类代码
				r.chengXiangFlDM = element.Text
				break
			case 2:
				// 名称
				r.mingCheng = element.Text
				break
			default:
			}
			href = ""
		})
		break
	default:
		e.ForEach("td > a", func(i int, element *colly.HTMLElement) {
			switch i {
			case 0:
				// 编码
				r.bianMa = element.Text
				break
			case 1:
				// 名称
				r.mingCheng = element.Text
				href = element.Attr("href")
				break
			default:
			}
		})
		break
	}
	r.fuJiBM = fuJiBM(r.bianMa, level)
	fmt.Printf("%+v\n", r)
	_, _ = file.WriteString(fmt.Sprintf("%s,%s,%d,%s,%s\n", r.bianMa, r.mingCheng, r.level, r.chengXiangFlDM, r.fuJiBM))
	if len(href) > 0 {
		_ = c.Visit(VisitUrl + href)
	}
}

func fuJiBM(bianMa string, level int) string {
	if len(bianMa) <= 0 {
		return ""
	}
	switch level {
	case ZhouShi:
		return bianMa[:2] + "0000000000"
	case QuXianShi:
		return bianMa[:4] + "00000000"
	case XiangZheng:
		return bianMa[:7] + "000000"
	case CunSheQu:
		return bianMa[:10] + "000"
	default:
		return ""
	}
}
