package policy_mode

import (
	"testing"
)

func TestMainer(t *testing.T) {
	// File
	fileL := NewFileLog()
	loggerManager := NewLoggerManager(fileL)
	loggerManager.Info()
	loggerManager.Error()

	// DB
	dbL := NewDBLog()
	loggerManager.Logger = dbL
	loggerManager.Info()
	loggerManager.Error()
}
