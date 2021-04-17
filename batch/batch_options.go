package batch

import "time"

type BatchOptions func(b *BatchProducer)

func WithMaxItems(maxItems uint64) BatchOptions{
	return func(b *BatchProducer) {
		b.MaxItems = maxItems
	}
}

func WithMaxWait(maxWait time.Duration) BatchOptions{
	return func(b *BatchProducer) {
		b.MaxWait = maxWait
	}
}
