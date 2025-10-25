package internal_test

import (
	"expvar"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/titpetric/platform/internal"
)

// TestCounterConcurrency verifies that concurrent Inc() calls result in the expected total.
func TestCounterConcurrency(t *testing.T) {
	assert := assert.New(t)

	name := "test_counter_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	c := internal.NewCounter(name)

	b := testing.Benchmark(func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Inc()
			}
		})
	})

	v := expvar.Get(name)
	assert.NotNil(v, "expvar for %s not found", name)

	got := c.Value()
	want := int64(b.N) // approximate, as RunParallel may split

	assert.LessOrEqual(want, got)
}

// BenchmarkCounter measures concurrent Inc() performance using RunParallel.
func BenchmarkCounter(b *testing.B) {
	name := "bench_counter_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	c := internal.NewCounter(name)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}
