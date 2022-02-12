package main

import (
	"study/pkg/crawler/engine"
	"study/pkg/crawler/persist"
	"study/pkg/crawler/scheduler"
	"study/pkg/crawler/zhenai/parser"
)

/*
	功能:
		First: 获取网页内容:
			使用http.Get获取内容
			使用Encoding转换编码:gbk -> utf-8
			使用 charset.DetermineEncoding来判断编码
		Second: 获取城市名称和链接
			使用正则表达式

	爬虫总体算法:
				==>	城市1 ==> 用户
		城市列表
				==>	城市2 ==> 用户


	解析器Parser:
		in: utf-8编码的文本
		out:Request{URL,对应的 Parser} 列表，Item 列表

	单任务版爬虫架构:
							Seed(Request)
			<-text			    ||		->URL
	Parser<==================Engine==================>  Fetcher
			->(requests,items)	||		<-text
				 			(任务队列)


*/

/*
	URL去重
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
	// config ES index
	itemChan, err := persist.ItemSaver(index)
	if err != nil {
		panic(err)
	}

	// config Engine
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      10,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}

	// start crawler
	e.Run(engine.Request{
		Url:    url,
		Parser: engine.NewFuncParser(parser.ParseCityList, "ParseCityList"),
	})
}
