package crawler

type Parser interface {
	Parse(contents []byte, url string) Result
	Serialize() (name string, args interface{})
}

type Request struct {
	Url    string
	Parser Parser
}

type Result struct {
	Requests []Request
	Items    []Item
}

type Item struct {
	Url     string
	Type    string
	Id      string
	PayLoad interface{}
}

type ParserFunc func(contents []byte, url string) Result

type NilParser struct{}

func (NilParser) Parse(_ []byte, _ string) Result {
	return Result{}
}

func (NilParser) Serialize() (name string, args interface{}) {
	return "NilParser", nil
}

type FuncParser struct {
	parser ParserFunc
	name   string
}

func (f *FuncParser) Parse(contents []byte, url string) Result {
	return f.parser(contents, url)
}

func (f *FuncParser) Serialize() (name string, args interface{}) {
	return f.name, nil
}

func NewFuncParser(p ParserFunc, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		name:   name,
	}
}
