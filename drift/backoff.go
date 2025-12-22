package drift

import (
	"math"
	"math/rand/v2"
	"time"
)

// Backoff defines the interface for calculating retry delays.
type Backoff interface {
	// Next returns the duration to wait before the next retry attempt.
	// The attempt parameter is zero-indexed (0 = first retry).
	Next(attempt int) time.Duration
}

// ExponentialBackoff implements exponential backoff with jitter.
// It calculates delays using the formula:
// min(initialTimeout * (exponentFactor ^ attempt), maxTimeout) + random(0, maxJitter)
type ExponentialBackoff struct {
	initialTimeout time.Duration
	maxTimeout     time.Duration
	exponentFactor float64
	maxJitter      time.Duration
}

// NewExponentialBackoff creates a new exponential backoff calculator.
//
// Parameters:
//   - initialTimeout: starting delay before first retry
//   - maxTimeout: maximum delay cap
//   - exponentFactor: multiplier for each retry (typically 2.0)
//   - maxJitter: maximum random jitter to add (prevents thundering herd)
func NewExponentialBackoff(
	initialTimeout, maxTimeout time.Duration,
	exponentFactor float64,
	maxJitter time.Duration,
) *ExponentialBackoff {
	return &ExponentialBackoff{
		initialTimeout: initialTimeout,
		maxTimeout:     maxTimeout,
		exponentFactor: exponentFactor,
		maxJitter:      maxJitter,
	}
}

// Next calculates the delay for the given attempt number (zero-indexed).
func (e *ExponentialBackoff) Next(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	// Calculate base delay: initialTimeout * exponentFactor^attempt
	baseDelay := float64(e.initialTimeout) * math.Pow(e.exponentFactor, float64(attempt))

	// Cap at maxTimeout (handles normal case and +Inf)
	maxTimeoutFloat := float64(e.maxTimeout)
	if math.IsNaN(baseDelay) || math.IsInf(baseDelay, 0) || baseDelay > maxTimeoutFloat {
		baseDelay = maxTimeoutFloat
	}

	// Ensure non-negative before conversion
	if baseDelay < 0 {
		baseDelay = 0
	}

	delay := time.Duration(baseDelay)

	// Handle overflow: large float64 values can become negative when cast to int64
	if delay < 0 {
		delay = e.maxTimeout
	}

	// Add jitter (0 to maxJitter)
	if e.maxJitter > 0 {
		jitter := time.Duration(rand.Int64N(int64(e.maxJitter) + 1)) //nolint:gosec // Jitter for backoff doesn't require crypto-grade randomness
		delay += jitter
	}

	return delay
}
