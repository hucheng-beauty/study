package main

import (
    "study/pkg/crawler"
)

/*
	功能:
		First(获取网页内容):
			使用 http.Get 获取内容
			使用 Encoding 转换编码(gbk -> utf-8)
			使用 charset.DetermineEncoding 来判断编码
		Second(获取城市名称和链接)
			使用正则表达式
	解析器Parser:
		in: utf-8编码的文本
		out: Request{ url, 对应的 Parser} 列表,Item 列表

	URL去重:
		哈希表
		计算URL的MD5等哈希,再存哈希表
		使用 bloom filter 多重哈希结构
		使用 Redis 等 K-V 存储系统实现分布式去重
*/

const (
    url   = "https://www.zhenai.com/zhenghun"
    index = "dating_profile"
)

func main() {
    itemChan, err := crawler.ItemSaver(index)
    if err != nil {
        panic(err)
    }

    crawler.NewEngine(
        crawler.WithScheduler(crawler.NewQueuedScheduler()),
        crawler.WithProcessor(crawler.FetchProcessor),
        crawler.WithWorkerCount(10),
        crawler.WithItemChan(itemChan),
    ).Run(crawler.Request{
        Url:    url,
        Parser: crawler.NewFuncParser(crawler.ParseCityList, "ParseCityList"),
    })
}
