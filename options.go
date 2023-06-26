package batch

import (
	"time"

	log "github.com/Deeptiman/go-batch/logger"
)

type BatchOptions func(b *BatchProducer)

func WithMaxItems(maxItems uint64) BatchOptions {
	return func(b *BatchProducer) {
		b.MaxItems = maxItems
	}
}

func WithMaxWait(maxWait time.Duration) BatchOptions {
	return func(b *BatchProducer) {
		b.MaxWait = maxWait
	}
}

func WithLogLevel(level log.LogLevel) BatchOptions {
	return func(b *BatchProducer) {
		b.Log.SetLogLevel(level)
	}
}
