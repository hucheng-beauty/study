package parser

import (
	"log"
	"regexp"
	"strconv"
	"study/internal/crawler/engine"
	"study/internal/crawler/model"
)

var (
	ageRe      = regexp.MustCompile(`<td><span class="label">年龄: </span>([\d]+)岁</td>`)
	marriageRe = regexp.MustCompile(`<td><span class="label">婚况: </span>([^<]+)</td>`)
	guessRe    = regexp.MustCompile(`<a class="exp-user-name"[^>]*href="(http://album.zhenai.com)/u/([\d]+)"`)
	idUrlRe    = regexp.MustCompile(`http://album.zhenai.com/u/([\d]+)`)
)

func parseProfile(contents []byte, url string, name string) engine.ParseResult {
	profile := model.Profile{}

	age, err := strconv.Atoi(extractString(contents, ageRe))
	if err != nil {
		log.Println("[ParseProfile][strconv] Error:", err)
	}
	profile.Age = age
	profile.Marriage = extractString(contents, marriageRe)
	profile.Name = name

	result := engine.ParseResult{
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "zhenai",
				Id:      extractString([]byte(url), idUrlRe),
				PayLoad: profile,
			},
		},
	}

	// to get next page
	matches := guessRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		// url: m[1], name: m[2]
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]),
			Parser: NewProfileParser(string(m[2])),
		})
	}

	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	}
	return ""
}

type ProfileParser struct {
	userName string
}

func (p *ProfileParser) Parse(contents []byte, url string) engine.ParseResult {
	return parseProfile(contents, url, p.userName)
}

func (p *ProfileParser) Serialize() (name string, args interface{}) {
	return "ProfileParser", p.userName
}

func NewProfileParser(name string) *ProfileParser {
	return &ProfileParser{userName: name}
}
