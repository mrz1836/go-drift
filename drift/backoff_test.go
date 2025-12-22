package drift

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewExponentialBackoff(t *testing.T) {
	t.Parallel()

	b := NewExponentialBackoff(
		2*time.Millisecond,
		100*time.Millisecond,
		2.0,
		5*time.Millisecond,
	)

	assert.NotNil(t, b)
	assert.Equal(t, 2*time.Millisecond, b.initialTimeout)
	assert.Equal(t, 100*time.Millisecond, b.maxTimeout)
	assert.InDelta(t, 2.0, b.exponentFactor, 0.0001)
	assert.Equal(t, 5*time.Millisecond, b.maxJitter)
}

func TestExponentialBackoffNext(t *testing.T) {
	t.Parallel()

	t.Run("basic exponential growth without jitter", func(t *testing.T) {
		b := NewExponentialBackoff(
			2*time.Millisecond,   // initial
			100*time.Millisecond, // max
			2.0,                  // exponent
			0,                    // no jitter for deterministic testing
		)

		// attempt 0: 2ms * 2^0 = 2ms
		assert.Equal(t, 2*time.Millisecond, b.Next(0))

		// attempt 1: 2ms * 2^1 = 4ms
		assert.Equal(t, 4*time.Millisecond, b.Next(1))

		// attempt 2: 2ms * 2^2 = 8ms
		assert.Equal(t, 8*time.Millisecond, b.Next(2))

		// attempt 3: 2ms * 2^3 = 16ms
		assert.Equal(t, 16*time.Millisecond, b.Next(3))

		// attempt 4: 2ms * 2^4 = 32ms
		assert.Equal(t, 32*time.Millisecond, b.Next(4))
	})

	t.Run("respects max timeout cap", func(t *testing.T) {
		b := NewExponentialBackoff(
			10*time.Millisecond,
			20*time.Millisecond, // low max
			2.0,
			0, // no jitter
		)

		// attempt 0: 10ms * 2^0 = 10ms (under max)
		assert.Equal(t, 10*time.Millisecond, b.Next(0))

		// attempt 1: 10ms * 2^1 = 20ms (at max)
		assert.Equal(t, 20*time.Millisecond, b.Next(1))

		// attempt 2: 10ms * 2^2 = 40ms, capped to 20ms
		assert.Equal(t, 20*time.Millisecond, b.Next(2))

		// attempt 10: would be huge, capped to 20ms
		assert.Equal(t, 20*time.Millisecond, b.Next(10))
	})

	t.Run("adds jitter within expected range", func(t *testing.T) {
		b := NewExponentialBackoff(
			10*time.Millisecond,
			100*time.Millisecond,
			2.0,
			5*time.Millisecond, // jitter
		)

		// Run multiple times to verify jitter range
		for i := 0; i < 100; i++ {
			delay := b.Next(0)
			// Base is 10ms, jitter is 0-5ms
			assert.GreaterOrEqual(t, delay, 10*time.Millisecond)
			assert.LessOrEqual(t, delay, 15*time.Millisecond)
		}
	})

	t.Run("handles negative attempt as zero", func(t *testing.T) {
		b := NewExponentialBackoff(
			2*time.Millisecond,
			10*time.Millisecond,
			2.0,
			0,
		)

		// Negative attempt treated as 0
		assert.Equal(t, 2*time.Millisecond, b.Next(-1))
		assert.Equal(t, 2*time.Millisecond, b.Next(-100))
	})

	t.Run("handles zero initial timeout", func(t *testing.T) {
		b := NewExponentialBackoff(
			0,
			10*time.Millisecond,
			2.0,
			0,
		)

		// 0 * anything = 0
		assert.Equal(t, time.Duration(0), b.Next(0))
		assert.Equal(t, time.Duration(0), b.Next(5))
	})

	t.Run("handles exponent factor of 1", func(t *testing.T) {
		b := NewExponentialBackoff(
			5*time.Millisecond,
			100*time.Millisecond,
			1.0, // no exponential growth
			0,
		)

		// 5ms * 1^n = 5ms always
		assert.Equal(t, 5*time.Millisecond, b.Next(0))
		assert.Equal(t, 5*time.Millisecond, b.Next(1))
		assert.Equal(t, 5*time.Millisecond, b.Next(10))
	})

	t.Run("handles fractional exponent factor", func(t *testing.T) {
		b := NewExponentialBackoff(
			10*time.Millisecond,
			100*time.Millisecond,
			1.5,
			0,
		)

		// attempt 0: 10ms * 1.5^0 = 10ms
		assert.Equal(t, 10*time.Millisecond, b.Next(0))

		// attempt 1: 10ms * 1.5^1 = 15ms
		assert.Equal(t, 15*time.Millisecond, b.Next(1))

		// attempt 2: 10ms * 1.5^2 = 22.5ms
		delay := b.Next(2)
		assert.GreaterOrEqual(t, delay, 22*time.Millisecond)
		assert.LessOrEqual(t, delay, 23*time.Millisecond)
	})

	t.Run("handles very high attempt numbers causing infinity", func(t *testing.T) {
		b := NewExponentialBackoff(
			1*time.Millisecond,
			100*time.Millisecond,
			2.0,
			0,
		)

		// Very high attempt number will cause math.Pow to return +Inf
		// which should be capped to maxTimeout
		delay := b.Next(10000)
		assert.Equal(t, 100*time.Millisecond, delay)
	})

	t.Run("handles large exponent factor causing overflow", func(t *testing.T) {
		b := NewExponentialBackoff(
			1*time.Hour,
			2*time.Hour,
			10.0, // Large factor
			0,
		)

		// Large values will cause overflow, should cap to maxTimeout
		delay := b.Next(100)
		assert.Equal(t, 2*time.Hour, delay)
	})

	t.Run("handles duration conversion overflow", func(t *testing.T) {
		b := NewExponentialBackoff(
			time.Duration(1<<62), // Very large value
			time.Duration(1<<62),
			2.0,
			0,
		)

		// When float64 to int64 conversion overflows, falls back to maxTimeout
		delay := b.Next(10)
		assert.GreaterOrEqual(t, delay, time.Duration(0))
	})
}

func TestExponentialBackoffImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ Backoff = (*ExponentialBackoff)(nil)
}

func BenchmarkExponentialBackoffNext(b *testing.B) {
	backoff := NewExponentialBackoff(
		2*time.Millisecond,
		100*time.Millisecond,
		2.0,
		5*time.Millisecond,
	)

	for i := 0; i < b.N; i++ {
		_ = backoff.Next(i % 10)
	}
}
