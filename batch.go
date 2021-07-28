package batch

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// Batch struct defines the structure payload for a Batch.
//
//	Item: channel that contains the Resources object from the client.
//	Id: Each item that a client send for the processing marked with Id.
//	Semaphore: The ReadWrite locks handle by the Semaphore object, it helps to synchronize the batch processing session.
//	Islocked: Whenever the batch processing session starts, Islocked changes to [true], so it will restrict the concurrent batch processing.
//  Producer: The BatchItem object send to the Producer for further processing.
//  Consumer: The Consumer arranges the prepared []BatchItems for the Workerline.
//  Log: Batch processing library uses "github.com/sirupsen/logrus" as logging tool.
type Batch struct {
	Item      chan interface{}
	Id        int
	Semaphore *Semaphore
	Islocked  bool
	Producer  *BatchProducer
	Consumer  *BatchConsumer
	Log       *log.Logger
}

// NewBatch creates a new Batch object with BatchProducer & BatchConsumer. The BatchOptions
// sets the MaxItems for a batch and maximum wait time for a batch to complete set by MaxWait.
func NewBatch(opts ...BatchOptions) *Batch {

	b := &Batch{
		Item: make(chan interface{}),
		Log:  log.New(),
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

// StartBatchProcessing function to begin the BatchProcessing library and to start the Producer/
// Consumer listeners. The ReadItems goroutine will receive the item from a source that keeps
// listening infinitely.
func (b *Batch) StartBatchProcessing() {

	b.Semaphore.Lock()
	defer b.Semaphore.Unlock()

	if b.Islocked {
		panic("Concurrent batch processing is not allowed!")
	}

	go b.Producer.WatchProducer()
	go b.Consumer.StartConsumer()
	go b.ReadItems()
}

// Unlock function will allow the batch processing to start with the multiple iteration
func (b *Batch) Unlock() {
	b.Islocked = false
}

// ReadItems function will run infinitely to listen to the Resource channel and the received
// object marshaled with BatchItem and then send to the Producer Watcher channel for further
// processing.
func (b *Batch) ReadItems() {

	b.Islocked = true

	for {

		select {
		case item := <-b.Item:
			b.Id++
			go func(item interface{}) {
				b.Producer.Watcher <- &BatchItems{
					Id:   b.Id,
					Item: item,
				}
			}(item)
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}
}

// StopProducer to exit the Producer line.
func (b *Batch) StopProducer() {
	b.Producer.Quit <- true
}

// Stop to run StopProducer/StopConsumer goroutines to quit the execution.
func (b *Batch) Stop() {
	go b.StopProducer()
}

// Close is the exit function to terminate the batch processing.
func (b *Batch) Close() {
	b.Log.WithFields(log.Fields{"Remaining Items": len(items)}).Warn("CheckRemainingItems")

	done := make(chan bool)

	go b.Producer.CheckRemainingItems(done)

	select {
	case <-done:
		b.Log.Warn("Done")
		b.Semaphore.Lock()
		b.Stop()
		close(b.Item)
		b.Islocked = false
		b.Semaphore.Unlock()
	}
}
