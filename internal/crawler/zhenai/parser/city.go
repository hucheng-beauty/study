package parser

import (
	"regexp"

	"study/pkg/crawler/engine"
)

var (
	profileRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<])+</a>`)
	cityUrlRe = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/[^"]=)"`)
)

// ParseCity :用户解析器
func ParseCity(contents []byte, url string) engine.ParseResult {
	result := engine.ParseResult{}

	matches := profileRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		// url: m[1], cityName: m[2]
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]),
			Parser: NewProfileParser(string(m[2])),
		})
	}

	// parse city page others contracts
	matches = cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]),
			Parser: engine.NewFuncParser(ParseCity, "ParseCity"),
		})
	}

	return result
}
