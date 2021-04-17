package batch

import (
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)


var (
	itemRead 	uint64
	itemCounter = 1
	DefaultMaxItems = uint64(100)
	DefaultMaxWait = time.Duration(30) * time.Second //seconds
	DefaultBatchNo = int32(1)

	items []BatchItems
)

type ConsumerFunc func(items []BatchItems)

type BatchProducer struct {
	Watcher 		chan *BatchItems 
	MaxItems 		uint64
	BatchNo			int32
	MaxWait 		time.Duration
	ConsumerFunc 	ConsumerFunc
	Quit 			chan bool
	Log 			*log.Logger
}

func NewBatchProducer(callBackFn ConsumerFunc, opts ...BatchOptions) *BatchProducer{

	return &BatchProducer {
		Watcher: make(chan *BatchItems),
		ConsumerFunc: callBackFn,
		MaxItems: DefaultMaxItems,
		MaxWait: DefaultMaxWait,
		BatchNo: DefaultBatchNo,
		Quit: make(chan bool),
		Log: log.New(),
	}
}

func (p *BatchProducer) prepareBatch(items []BatchItems) []BatchItems {

	p.ConsumerFunc(items)
	return p.resetItem(items)
}

func (p *BatchProducer) WatchProducer() {

	for {

		select {
		case item := <-p.Watcher:
			
			item.BatchNo = 	int(p.getBatchNo())
			p.Log.WithFields(log.Fields{"Item": item, "GetItemCounter": itemCounter}).Info("Produce Items")
			
			items = append(items, *item)
			itemCounter++
			if  p.isBatchReady() {
				p.Log.WithFields(log.Fields{"Item Size": len(items), "MaxItems": p.MaxItems}).Warn("BatchReady")
			
				itemCounter = 0
				items = p.prepareBatch(items)
			}
			
		case <-time.After(p.MaxWait):
			p.Log.WithFields(log.Fields{"Items": len(items)}).Warn("MaxWait")
			
			if len(items) == 0 {
				return
			}
			itemCounter = 0
			items = p.prepareBatch(items)
		case <-p.Quit:
			p.Log.Warn("Quit BatchProducer")

			return
		}
	}
}

func (p *BatchProducer) CheckRemainingItems(done chan bool) {

	if len(items) >= 1 {
		p.prepareBatch(items)
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	done <- true
}

func (p *BatchProducer) resetItem(items []BatchItems) []BatchItems {
	defer p.resetItemRead()
	items = items[:0]
	return items
}

func (p *BatchProducer) isBatchReady() bool {
	return uint64(itemCounter) >= p.MaxItems
}

func (p *BatchProducer) resetItemRead() {
	atomic.StoreUint64(&itemRead, 1)
}

func (p *BatchProducer) addBatchNo() {
	atomic.AddInt32(&p.BatchNo, 1)
}

func (p *BatchProducer) getBatchNo() int32 {
	
	if itemCounter == 0 {
		p.addBatchNo()
	}

	return atomic.LoadInt32(&p.BatchNo)
}