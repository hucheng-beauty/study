package initialize

import (
	"testing"

	"go.uber.org/zap"
)

func TestMainer(t *testing.T) {
	// simple use
	logger, _ := zap.NewProduction()
	url := "test.com"
	// 性能高
	logger.Info("failed to fetch url",
		zap.String("url", url),
		zap.Int("id", 3),
	)

	// log file
	var NewLogger = func() (*zap.Logger, error) {
		logCfg := zap.NewProductionConfig()
		logCfg.OutputPaths = []string{
			"./admin_api.log",
		}
		return logCfg.Build()
	}

	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}
}
