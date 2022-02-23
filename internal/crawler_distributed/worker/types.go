package worker

import (
	"errors"
	"fmt"
	"log"

	"study/internal/crawler/engine"
	"study/internal/crawler/zhenai/parser"
)

type SerializedParser struct {
	Name string
	Args interface{}
}

type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items    []engine.Item
	Requests []Request
}

func SerializedRequest(r engine.Request) Request {
	serialize, args := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: serialize,
			Args: args},
	}
}

func SerializedResult(r engine.ParseResult) ParseResult {
	result := ParseResult{Items: r.Items}
	for _, req := range r.Requests {
		result.Requests = append(result.Requests, SerializedRequest(req))
	}
	return result
}

func DeserializedRequest(r Request) (engine.Request, error) {
	p, err := deserializedRequest(r.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{Url: r.Url, Parser: p}, nil
}

func deserializedRequest(p SerializedParser) (engine.Parser, error) {
	switch p.Name {
	case "ParseCityList":
		return engine.NewFuncParser(parser.ParseCityList, "ParseCityList"), nil
	case "ParseCity":
		return engine.NewFuncParser(parser.ParseCity, "ParseCity"), nil
	case "ParseProfile":
		if userName, ok := p.Args.(string); ok {
			return parser.NewProfileParser(userName), nil
		} else {
			return nil, fmt.Errorf("invalid arg:%v", p.Args)
		}
	case "NilParser":
		return engine.NilParser{}, nil
	default:
		return nil, errors.New("unknown parser name")

	}
}

func DeserializedResult(r ParseResult) engine.ParseResult {
	result := engine.ParseResult{Items: r.Items}
	for _, req := range r.Requests {
		engineReq, err := DeserializedRequest(req)
		if err != nil {
			log.Printf("error deserializing request:%v", err)
			continue
		}
		result.Requests = append(result.Requests, engineReq)
	}
	return result
}
