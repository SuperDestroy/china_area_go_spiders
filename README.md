# 中国全国5级行政区划（省、市、县、镇、村）爬虫
### 采用Golang语言
#### 插件使用 [colly](https://github.com/gocolly/colly)
#### 数据来源 [国家统计局网站](http://www.stats.gov.cn/sj/tjbz/tjyqhdmhcxhfdm/2022/)

程序采用单线程模式，避免IP被拦。

#### 解析代码示例：
```
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
```

#### 数据字段含义
```
type Region struct {
    bianMa  string // 编码
    mingCheng      string // 名称
    chengXiangFlDM string // 城乡分类代码
    level          int    // 级别
    fuJiBM         string // 父级编码
}
```
### regions.txt.zip
此文件是通过程序获取的数据\
每行为一个行政区划\
格式: 编码,名称,城乡分类代码,级别,父级编码

#### 级别定义
```
const (
	Sheng = iota // 0
	ZhouShi
	QuXianShi
	XiangZheng
	CunSheQu
)
```
