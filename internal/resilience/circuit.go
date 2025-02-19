package resilience

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreaker is an alias for gobreaker.CircuitBreaker to make it part of our package
type CircuitBreaker = gobreaker.CircuitBreaker

// CircuitBreakerConfig holds the configuration for circuit breaker
type CircuitBreakerConfig struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	ReadyToTrip   func(counts gobreaker.Counts) bool
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

// NewCircuitBreaker creates a new circuit breaker with default configuration
func NewCircuitBreaker(name string) *CircuitBreaker {
	return NewCircuitBreakerWithConfig(CircuitBreakerConfig{
		Name:        name,
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("Circuit breaker %s state changed from %s to %s\n", name, from, to)
		},
	})
}

// NewCircuitBreakerWithConfig creates a new circuit breaker with custom configuration
func NewCircuitBreakerWithConfig(config CircuitBreakerConfig) *CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:          config.Name,
		MaxRequests:   config.MaxRequests,
		Interval:      config.Interval,
		Timeout:       config.Timeout,
		ReadyToTrip:   config.ReadyToTrip,
		OnStateChange: config.OnStateChange,
	})
}