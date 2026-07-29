package retry

import (
	"sync"
	"time"
)

// CircuitBreaker prevents cascading failures
type CircuitBreaker struct {
	mu              sync.Mutex
	failures        int
	successes       int
	maxFailures     int
	resetTimeout    time.Duration
	lastFailureTime time.Time
	state           State
	maxSuccesses    int
}

// State represents circuit breaker state
type State int

const (
	// StateClosed - normal operation
	StateClosed State = iota
	// StateOpen - failing, reject requests
	StateOpen
	// StateHalfOpen - testing if service recovered
	StateHalfOpen
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, maxSuccesses int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		maxSuccesses: maxSuccesses,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}
}

// Allow returns true if requests are allowed
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateOpen:
		// Check if reset timeout passed
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.successes = 0
			return true
		}
		return false
	case StateHalfOpen:
		// Allow limited requests to test recovery
		return cb.successes < cb.maxSuccesses
	default:
		return true
	}
}

// RecordFailure records a failed request
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailureTime = time.Now()

	if cb.failures >= cb.maxFailures {
		cb.state = StateOpen
	}
}

// RecordSuccess records a successful request
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.successes++

	if cb.state == StateHalfOpen && cb.successes >= cb.maxSuccesses {
		// Recovery successful
		cb.state = StateClosed
		cb.failures = 0
	}
}

// State returns current circuit breaker state
func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}
