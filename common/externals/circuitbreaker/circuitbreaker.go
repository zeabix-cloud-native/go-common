package circuitbreaker

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
	"github.com/zeabix-cloud-native/go-common/common/config"
)

type CircuitResult struct {
	StatusCode int
	Result     interface{}
	PageTotal  int
}

func NewCircuitBreaker(cf config.Config) *gobreaker.CircuitBreaker {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        fmt.Sprintf("%s-CIRCUIT-BEAKER", cf.AppName),
		MaxRequests: 3,
		Interval:    5 * time.Second,
		Timeout:     1 * time.Second,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("Circuit breaker state changed from %s to %s\n", from, to)
		},
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})
	return cb
}
