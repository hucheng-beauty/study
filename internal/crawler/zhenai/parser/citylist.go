package parser

import (
	"regexp"
	"study/internal/crawler/engine"
)

var cityListRe = regexp.MustCompile(`<a href="(https://www.zhenai.com/zhenhun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)

// ParseCityList 城市列表解析器
func ParseCityList(contents []byte, url string) engine.ParseResult {
	result := engine.ParseResult{}

	matches := cityListRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]),
			Parser: engine.NewFuncParser(ParseCity, "ParseCity"),
		})
	}
	return result
}
