package batch

import (
	"sync/atomic"
	"time"

	plog "github.com/sirupsen/logrus"
)


var (
	itemRead 	uint64
	itemId = 1
	DefaultMaxItems = uint64(100)
	DefaultMaxWait = time.Duration(30) * time.Second //seconds

	items []BatchItems
)

type ConsumerFunc func(items []BatchItems)

type BatchProducer struct {
	Watcher 		chan interface{} 
	MaxItems 		uint64
	BatchNo			int32
	MaxWait 		time.Duration
	CallBackFunc 	ConsumerFunc
	Quit 			chan bool
	Log 			*plog.Logger
}

func NewBatchProducer(callBackFn ConsumerFunc, opts ...BatchOptions) *BatchProducer{

	return &BatchProducer {
		Watcher: make(chan interface{}),
		CallBackFunc: callBackFn,
		MaxItems: DefaultMaxItems,
		MaxWait: DefaultMaxWait,
		Quit: make(chan bool),
		Log: plog.New(),
	}
}

func (p *BatchProducer) prepareBatch(items []BatchItems) []BatchItems {

	p.addBatchNo()
	p.CallBackFunc(items)
	return p.resetItem(items)
}

func (p *BatchProducer) WatchProducer() {

	for {

		select {
		case item := <-p.Watcher:
			
			p.Log.WithFields(plog.Fields{"Item": item, "GetItemRead": p.getItemRead()}).Info("Produce Items")

			p.itemRead()
			items = append(items, *p.getBatchItem(item))
			itemId++
			if  p.isBatchReady() {
				items = p.prepareBatch(items)
			}
			
		case <-time.After(p.MaxWait):
			p.Log.WithFields(plog.Fields{"Items": len(items)}).Warn("MaxWait")
			
			if len(items) == 0 {
				return
			}

			items = p.prepareBatch(items)
		case <-p.Quit:
			p.Log.Warn("Quit")
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

func (p *BatchProducer) itemRead() {
	atomic.AddUint64(&itemRead, 1)
}

func (p *BatchProducer) getItemRead() uint64{
	return atomic.LoadUint64(&itemRead)
}

func (p *BatchProducer) isBatchReady() bool {
	itemRead := p.getItemRead()
	return p.MaxItems > uint64(0) && itemRead >= p.MaxItems
}

func (p *BatchProducer) getBatchItem(item interface{}) *BatchItems {
	return &BatchItems {
		Id: itemId,
		Item: item,
	}
}

func (p *BatchProducer) resetItemRead() {
	atomic.StoreUint64(&itemRead, 1)
}

func (p *BatchProducer) addBatchNo() {
	atomic.AddInt32(&p.BatchNo, 1)
}

func (p *BatchProducer) getBatchNo() int32 {
	return atomic.LoadInt32(&p.BatchNo)
}