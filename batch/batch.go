package batch

import (
	"time"
	log "github.com/sirupsen/logrus"
)

type Batch struct {
	Item 			chan interface{}
	Semaphore 		*Semaphore
	Producer		*BatchProducer
	Consumer 		*BatchConsumer
	Log 			*log.Logger
}


func NewBatch(opts ...BatchOptions) *Batch{
	
	b := &Batch{
		Item: make(chan interface{}),
		Log: log.New(),
	}

	c := NewBatchConsumer()

	p := NewBatchProducer(c.ConsumerFunc)
	
	for _, opt := range opts {
		opt(p)
	}

	b.Producer = p
	b.Consumer = c
	b.Semaphore = NewSemaphore(int(p.MaxItems))
 
	items = make([]BatchItems, 0, p.MaxItems)

	return b
}


func (b *Batch) StartBachProcessing() {

	go b.Producer.WatchProducer()
	go b.Consumer.StartConsumer()
	go b.ReadItems()
}

func (b *Batch) ReadItems() {

	for {

		select {
		case item := <-b.Item:

			go func(item interface{}){
				b.Producer.Watcher <- item
			}(item)		
			time.Sleep(time.Duration(100) * time.Millisecond)				
		}
	}
}

func (b *Batch) Stop() {
	b.Producer.Quit <- true
}

func (b *Batch) Close() {
	b.Log.WithFields(log.Fields{"Remaining Items": len(items)}).Warn("Close")

	done := make(chan bool)

	go b.Producer.CheckRemainingItems(done)

	for {
		select {
		case <-done:
			b.Log.WithFields(log.Fields{"Remaining Items": len(items)}).Warn("Done")
			b.Semaphore.Lock()
			b.Stop()
			b.Consumer.Shutdown()
			close(b.Item)	
			b.Semaphore.Unlock()
		}
	}

}