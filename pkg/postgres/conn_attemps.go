package postgres

import (
	"fmt"
	"time"

	"github.com/ellofae/payment-system-kafka/pkg/logger"
)

func ConnectionAttemps(conn_func func() error, attemps int, delay time.Duration) (err error) {
	log := logger.GetLogger()

	for i := 0; i < attemps; i++ {
		err = conn_func()
		if err != nil {
			log.Warn(fmt.Sprintf("Attempting to connect: current attemp - %d, attemps left - %d", i+1, attemps-i-1))
			time.Sleep(delay)
			continue
		}
	}
	return
}
