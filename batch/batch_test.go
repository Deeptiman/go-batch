package batch

import (
	"fmt"
	//"time"
	"testing"
)

// Resources Structure
type Resources struct {
	id   int
	name string
	flag bool
}

func TestBatch(t *testing.T) {

	t.Run("Batch-1", func(t *testing.T) {
		t.Parallel()

		mFlag := 10
		rFlag := 15

		b := NewBatch(WithMaxItems(uint64(mFlag)))

		b.StartBatchProcessing()

		for i := 1; i <= rFlag; i++ {
			b.Item <- &Resources{
				id:   i,
				name: fmt.Sprintf("%s%d", "R-", i),
				flag: false,
			}
		}
		b.Close()
	})

	t.Run("Batch-2", func(t *testing.T) {
		t.Parallel()

		mFlag := 10
		rFlag := 100

		b := NewBatch(WithMaxItems(uint64(mFlag)))

		b.StartBatchProcessing()
		for i := 1; i <= rFlag; i++ {
			b.Item <- &Resources{
				id:   i,
				name: fmt.Sprintf("%s%d", "R-", i),
				flag: false,
			}
		}

		var panics bool
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Log("Recover--", r)
					panics = true
				}
			}()

			b.StartBatchProcessing()
			for i := 1; i <= rFlag; i++ {
				b.Item <- &Resources{
					id:   i,
					name: fmt.Sprintf("%s%d", "R-", i),
					flag: false,
				}
			}
		}()

		if !panics {
			t.Error("Testing Concurrent batch processing should fail")
		}

		b.Close()
	})
}

func benchmarkBatch(numResources, maxItems int, b *testing.B) {
	for n := 0; n < b.N; n++ {

		b := NewBatch(WithMaxItems(uint64(maxItems)))

		b.StartBatchProcessing()
		for i := 1; i <= numResources; i++ {
			b.Item <- &Resources{
				id:   i,
				name: fmt.Sprintf("%s%d", "R-", i),
				flag: false,
			}
		}
		b.Close()
	}
}

func BenchmarkBatchR100M5(b *testing.B) {
	benchmarkBatch(100, 5, b)
}

func BenchmarkBatchR10M5(b *testing.B) {
	benchmarkBatch(10, 5, b)
}

func BenchmarkBatchR100M10(b *testing.B) {
	benchmarkBatch(100, 10, b)
}

func BenchmarkBatch(b *testing.B) {

	b.ResetTimer()
	b.Run("Bench-1", func(b *testing.B) {
		benchmarkBatch(100, 5, b)
	})

	b.Run("Bench-2", func(b *testing.B) {
		benchmarkBatch(500, 50, b)
	})

	b.Run("Bench-3", func(b *testing.B) {
		benchmarkBatch(3000, 300, b)
	})

	b.Run("Bench-4", func(b *testing.B) {
		benchmarkBatch(10000, 100, b)
	})
}
