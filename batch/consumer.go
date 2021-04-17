package batch

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultWorkerPool = 10
)

type BatchConsumer struct {
	ConsumerCh 		chan []BatchItems
	BatchWorkerCh 	chan []BatchItems
	Supply 			*BatchSupply
	Workerline		*sync.WaitGroup
	TerminateCh		chan os.Signal
	Quit 			chan bool
	Log				*log.Logger
}

type BatchSupply struct {
	BatchSupplyCh 	chan chan []BatchItems
	ClientSupplyCh 	chan []BatchItems
}

func NewBatchConsumer() *BatchConsumer{

	return &BatchConsumer{
		ConsumerCh: make(chan []BatchItems, 1),
		BatchWorkerCh: make(chan []BatchItems, DefaultWorkerPool),
		Supply: NewBatchSupply(),
		Workerline:	&sync.WaitGroup{},
		TerminateCh: make(chan os.Signal, 1),
		Quit: make(chan bool, 1),
		Log: log.New(),
	}
}

func NewBatchSupply() *BatchSupply{
	return &BatchSupply {
		BatchSupplyCh: make(chan chan []BatchItems, 100),
		ClientSupplyCh: make(chan []BatchItems, 1),
	}
}

func(c *BatchConsumer) StartConsumer() {

	ctx, cancel := context.WithCancel(context.Background())
	
	go c.ConsumerBatch(ctx)

	c.Workerline.Add(DefaultWorkerPool)
	for i:=0; i < DefaultWorkerPool; i++ {
		go c.WorkerFunc(i)
	}

	signal.Notify(c.TerminateCh, syscall.SIGINT, syscall.SIGTERM)
	<-c.TerminateCh

	cancel()
	c.Workerline.Wait()
}

func (c *BatchConsumer) ConsumerFunc(items []BatchItems) {
	c.ConsumerCh <- items
}

func(c *BatchConsumer) ConsumerBatch(ctx context.Context) {

	for {
		select {
		case batchItems := <-c.ConsumerCh:
			c.BatchWorkerCh <- batchItems
		case <-ctx.Done():
			c.Log.Warn("Request cancel signal received!")
			close(c.BatchWorkerCh)
			return
		case <-c.Quit:
			c.Log.Warn("Quit BatchConsumer")			
			close(c.BatchWorkerCh)
			return 
		}
	}
}

func (c *BatchConsumer) WorkerFunc(index int) {
	defer c.Workerline.Done()

	for batch := range c.BatchWorkerCh {
	
		c.Log.WithFields(log.Fields{"Worker": index, "Batch": batch}).Info("BatchConsumer")
			
		go c.GetBatchSupply()

		select {
		case supplyCh := <-c.Supply.BatchSupplyCh:
			supplyCh <- batch
		}
	}
}

func (c *BatchConsumer) Shutdown() {

	c.Log.Warn("Shutdown signal received!")
	signal.Notify(c.TerminateCh, syscall.SIGINT, syscall.SIGTERM)
	<-c.TerminateCh
}