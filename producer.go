package batch

import (
	"sync/atomic"
	"time"
	log "github.com/Deeptiman/go-batch/logger"
)

var (	
	DefaultMaxItems = uint64(100)                     // maximum no of items packed inside a Batch
	DefaultMaxWait  = time.Duration(30) * time.Second //seconds
	DefaultBatchNo  = int32(1)

	items []BatchItems
)

// ConsumerFunc is the callback function that invoke from Consumer
type ConsumerFunc func(items []BatchItems)

// BatchProducer struct defines the Producers fields that requires to create a []BatchItems object.
//
// Watcher: The receiver channel that gets the BatchItems marshalled object from Batch reader.
// MaxItems: Maximum no of BatchItems can be packed for a released Batch.
// BatchNo: Every []BatchItems that gets released marked with BatchNo [integer].
// MaxWait: If a batch processing takes too long, then MaxWait has the timeout that expires after an interval.
// ConsumerFunc: It's the callback function that gets invoke by the Consumer
// Quit: It's the exit channel for the Producer to end the processing
// Log: Batch processing library uses "github.com/sirupsen/logrus" as logging tool.
type BatchProducer struct {
	Watcher      chan *BatchItems
	MaxItems     uint64
	BatchNo      int32
	MaxWait      time.Duration
	ConsumerFunc ConsumerFunc
	Quit         chan bool
	Log          *log.Logger
}

// NewBatchProducer defines the producer line for creating a Batch. There will be a Watcher
// channel that receives the incoming BatchItem from the source. The ConsumerFunc works as a
// callback function to the Consumer line to release the newly created set of BatchItems.
//
//
// Each Batch is registered with a BatchNo that gets created when the Batch itemCounter++ increases
// to the MaxItems value.
func NewBatchProducer(callBackFn ConsumerFunc, opts ...BatchOptions) *BatchProducer {

	return &BatchProducer{
		Watcher:      make(chan *BatchItems),
		ConsumerFunc: callBackFn,
		MaxItems:     DefaultMaxItems,
		MaxWait:      DefaultMaxWait,
		BatchNo:      DefaultBatchNo,
		Quit:         make(chan bool),
		Log:          log.NewLogger(),
	}
}

// WatchProducer has the Watcher channel that receives the BatchItem object from the Batch read
// item channel. Watcher marks each BatchItem with a BatchNo and adds it to the []BatchItems array.
// After the batch itemCounter++ increases to the MaxItems [DefaultMaxItems: 100], the Batch gets
// releases to the Consumer callback function.
//
// If the Batch processing get to halt in the Watcher
// channel then the MaxWait [DefaultMaxWait: 30 sec] timer channel gets called to check the state
// to releases the Batch to the Consumer callback function.
func (p *BatchProducer) WatchProducer() {

	for {

		select {
		case item := <-p.Watcher:

			item.BatchNo = int(p.getBatchNo())
			p.Log.Debugln("BatchProducer", "Id=", item.Id, "Batch Break=", item.Id / int(p.MaxItems), "BatchNo=",item.BatchNo, "Item=", item.Item)

			items = append(items, *item)			
			
			if (item.Id / int(p.MaxItems)) == item.BatchNo {
				p.Log.Infoln("BatchReady", "BatchNo=", item.BatchNo)
				items = p.releaseBatch(items)
				p.createBatchNo()
			}
			
		case <-time.After(p.MaxWait):
			p.Log.Infoln("MaxWait", "Items=", len(items))
			if len(items) == 0 {
				continue
			}
			
			items = p.releaseBatch(items)
		case <-p.Quit:
			p.Log.Warn("Quit BatchProducer")

			return
		}
	}
}

// releaseBatch will call the Consumer callback function to send the prepared []BatchItems to
// the Consumer line. Also it reset the []BatchItems array (items = items[:0]) to begin the
// next set of batch processing.
func (p *BatchProducer) releaseBatch(items []BatchItems) []BatchItems {

	p.ConsumerFunc(items)
	return p.resetItem(items)
}

// resetItem to slice the []BatchItems to empty.
func (p *BatchProducer) resetItem(items []BatchItems) []BatchItems {
	items = items[:0]
	return items
}

// CheckRemainingItems is a force re-check function on remaining batch items that are available
// for processing.
func (p *BatchProducer) CheckRemainingItems(done chan bool) {

	if len(items) >= 1 {
		p.releaseBatch(items)
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	done <- true
}

// addBatchNo will increases the current BatchNo to 1 atomically.
func (p *BatchProducer) createBatchNo() {
	atomic.AddInt32(&p.BatchNo, 1)
}

// getBatchNo will get the current BatchNo from the atomic variable.
func (p *BatchProducer) getBatchNo() int32 {
	return atomic.LoadInt32(&p.BatchNo)
}
