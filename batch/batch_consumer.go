package batch

import (
	"fmt"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	DefaultWorkerPool = 10
)

type BatchConsumer struct {
	ConsumerCh 		chan []BatchItems
	BatchWorkerCh 	chan []BatchItems
	Supply 			*BatchSupply
	wg 				*sync.WaitGroup
	TerminateCh		chan os.Signal
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
		wg: &sync.WaitGroup{},
		TerminateCh: make(chan os.Signal, 1),
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
	//wg := &sync.WaitGroup{}

	go c.ConsumerBatch(ctx)

	c.wg.Add(DefaultWorkerPool)
	for i:=0; i < DefaultWorkerPool; i++ {
		go c.WorkerFunc(i)
	}

	signal.Notify(c.TerminateCh, syscall.SIGINT, syscall.SIGTERM)
	<-c.TerminateCh

	cancel()
	c.wg.Wait()
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
			fmt.Println("Request cancel signal received!")
			//close(c.BatchWorkerCh)
			return
		}
	}
}

func (c *BatchConsumer) WorkerFunc(index int) {
	//defer c.wg.Done()

	for batch := range c.BatchWorkerCh {
		fmt.Println("Worker-", index, " : Batch -- ", batch)

		go c.GetBatchSupply()

		select {
		case supplyCh := <-c.Supply.BatchSupplyCh:
			supplyCh <- batch
		}
	}
}

func (c *BatchConsumer) Shutdown() {

	close(c.BatchWorkerCh)

	c.wg.Done()

	fmt.Println("Shutdown signal received!")
	signal.Notify(c.TerminateCh, syscall.SIGINT, syscall.SIGTERM)
	<-c.TerminateCh
}