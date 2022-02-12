package worker

import (
	"study/pkg/crawler/engine"
)

type CrawService struct{}

func (CrawService) Process(req Request, result *ParseResult) error {
	// Request ==> engine.Request
	engineReq, err := DeserializedRequest(req)
	if err != nil {
		return err
	}

	// engine.ParseResult ==> ParseResult
	engineResult, err := engine.Worker(engineReq)
	if err != nil {
		return err
	}

	*result = SerializedResult(engineResult)
	return nil
}
