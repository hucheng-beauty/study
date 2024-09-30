package crawler

import (
	"regexp"
	"strconv"
)

var (
	cityUrlRe  = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/[^"]=)"`)
	cityListRe = regexp.MustCompile(`<a href="(https://www.zhenai.com/zhenhun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)
	profileRe  = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<])+</a>`)

	ageRe      = regexp.MustCompile(`<td><span class="label">年龄: </span>([\d]+)岁</td>`)
	marriageRe = regexp.MustCompile(`<td><span class="label">婚况: </span>([^<]+)</td>`)
	guessRe    = regexp.MustCompile(`<a class="exp-user-name"[^>]*href="(http://album.zhenai.com)/u/([\d]+)"`)
	idUrlRe    = regexp.MustCompile(`http://album.zhenai.com/u/([\d]+)`)
)

// ParseCity 用户解析器
func ParseCity(contents []byte, url string) (result Result) {
	matches := profileRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		// url: m[1], cityName: m[2]
		result.Requests = append(result.Requests, Request{
			Url:    string(m[1]),
			Parser: NewProfileParser(string(m[2])),
		})
	}

	// parse city page others contracts
	matches = cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, Request{
			Url:    string(m[1]),
			Parser: NewFuncParser(ParseCity, "ParseCity"),
		})
	}

	return result
}

// ParseCityList 城市列表解析器
func ParseCityList(contents []byte, url string) (result Result) {
	matches := cityListRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, Request{
			Url:    string(m[1]),
			Parser: NewFuncParser(ParseCity, "ParseCity"),
		})
	}
	return result
}

// ProfileParser 用户画像解析
type ProfileParser struct {
	userName string
}

func (p *ProfileParser) Parse(contents []byte, url string) Result {
	return p.parse(contents, url)
}

func (p *ProfileParser) parse(contents []byte, url string) Result {
	var extractValue = func(contents []byte, re *regexp.Regexp) string {
		match := re.FindSubmatch(contents)
		if len(match) >= 2 {
			return string(match[1])
		}
		return ""
	}

	age, _ := strconv.Atoi(extractValue(contents, ageRe))
	result := Result{
		Items: []Item{
			{
				Url:  url,
				Type: "zhenai",
				Id:   extractValue([]byte(url), idUrlRe),
				PayLoad: Profile{
					Name:          p.userName,
					Gender:        "",
					Age:           age,
					Height:        0,
					Weight:        0,
					Income:        "",
					Marriage:      extractValue(contents, marriageRe),
					Education:     "",
					Occupation:    "",
					Household:     "",
					Constellation: "",
					House:         "",
					Car:           "",
				},
			},
		},
	}

	// to gain next page
	matches := guessRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		// url: m[1], userName: m[2]
		result.Requests = append(result.Requests, Request{
			Url:    string(m[1]),
			Parser: NewProfileParser(string(m[2])),
		})
	}

	return result
}

func (p *ProfileParser) Serialize() (name string, args interface{}) {
	return "ProfileParser", p.userName
}

func NewProfileParser(username string) *ProfileParser {
	return &ProfileParser{userName: username}
}
