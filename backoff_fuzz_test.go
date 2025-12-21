package drift

import (
	"math"
	"testing"
	"time"
)

func FuzzExponentialBackoffNext(f *testing.F) {
	// Seed corpus with edge cases
	f.Add(0)
	f.Add(-1)
	f.Add(1)
	f.Add(10)
	f.Add(100)
	f.Add(1000)
	f.Add(math.MaxInt32)
	f.Add(math.MinInt32)

	// Create backoff with known parameters
	initialTimeout := 2 * time.Millisecond
	maxTimeout := 100 * time.Millisecond
	maxJitter := 5 * time.Millisecond

	b := NewExponentialBackoff(
		initialTimeout,
		maxTimeout,
		2.0,
		maxJitter,
	)

	f.Fuzz(func(t *testing.T, attempt int) {
		// Should never panic
		result := b.Next(attempt)

		// Result must be non-negative
		if result < 0 {
			t.Errorf("Next(%d) returned negative duration: %v", attempt, result)
		}

		// Result must not exceed maxTimeout + maxJitter
		maxExpected := maxTimeout + maxJitter
		if result > maxExpected {
			t.Errorf("Next(%d) = %v, exceeds max expected %v", attempt, result, maxExpected)
		}
	})
}

func FuzzExponentialBackoffNextWithParams(f *testing.F) {
	// Seed corpus with various parameter combinations
	f.Add(int64(time.Millisecond), int64(time.Second), 2.0, int64(time.Millisecond), 0)
	f.Add(int64(0), int64(time.Second), 2.0, int64(0), 0)
	f.Add(int64(time.Millisecond), int64(time.Millisecond), 1.0, int64(0), 10)
	f.Add(int64(time.Millisecond), int64(time.Second), 0.5, int64(time.Millisecond), 100)

	f.Fuzz(func(t *testing.T, initial, maxDuration int64, factor float64, jitter int64, attempt int) {
		// Skip invalid inputs that would cause issues in normal usage
		if initial < 0 || maxDuration < 0 || jitter < 0 {
			t.Skip("Skipping negative duration parameters")
		}
		if math.IsNaN(factor) || math.IsInf(factor, 0) {
			t.Skip("Skipping NaN/Inf factor")
		}
		if factor < 0 {
			t.Skip("Skipping negative factor")
		}

		b := NewExponentialBackoff(
			time.Duration(initial),
			time.Duration(maxDuration),
			factor,
			time.Duration(jitter),
		)

		// Should never panic
		result := b.Next(attempt)

		// Result must be non-negative
		if result < 0 {
			t.Errorf("Next(%d) returned negative duration: %v", attempt, result)
		}
	})
}
