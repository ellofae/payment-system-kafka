package logger

import (
	"os"
	"sync"

	"github.com/hashicorp/go-hclog"
)

var logger hclog.Logger
var once sync.Once

func GetLogger() hclog.Logger {
	once.Do(func() {
		loggerOptions := &hclog.LoggerOptions{
			Name:            "payment-system",
			Level:           hclog.Info,
			Output:          os.Stderr,
			IncludeLocation: true,
			TimeFormat:      "2006-01-02 15:04:05.000",
		}

		logger = hclog.New(loggerOptions)
	})

	return logger
}
