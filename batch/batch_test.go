package batch

import (
	"fmt"
	//"time"
	"testing"
)

// Resources Structure
type Resources struct {
	id 		int
	name 	string
	flag 	bool
}

func TestBatch(t *testing.T){

	t.Run("Batch-1", func(t *testing.T){
		t.Parallel()
		
		mFlag := 10
		rFlag := 15

		b := NewBatch(WithMaxItems(uint64(mFlag)))

		b.StartBatchProcessing()

		for i := 1; i <= rFlag; i++ {
			b.Item <- &Resources{
				id:   i,
				name: fmt.Sprintf("%s%d","R-",  i),
				flag: false,
			}
		}
		b.Close()
	})

	t.Run("Batch-2", func(t *testing.T){
		t.Parallel()

		mFlag := 10
		rFlag := 100

		b := NewBatch(WithMaxItems(uint64(mFlag)))
		 
		b.StartBatchProcessing()
		for i := 1; i <= rFlag; i++ {
			b.Item <- &Resources{
				id:   i,
				name: fmt.Sprintf("%s%d","R-",  i),
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
					name: fmt.Sprintf("%s%d","R-",  i),
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