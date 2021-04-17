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

// NewBatchConsumer defines several types of production channels those are works at a different
// stages to release a Batch to the client. The ConsumerCh received the Batch and send it to the
// Workers channel. Then, the Workerline arranges the worker under a waitGroup to release the Batch
// to the Supply channel. 
//
//
// The BatchSupply has a bidirectional channel that requests a Batch from
// the Worker channel and receives a Batch via response channel. Also, BatchSupply has a Client 
// channel that sends the released Batch to the Client. The client needs to listen to the ClientSupplyCh
// to receive batch instantly.
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

// NewBatchSupply will create the BatchSupply object that has two sets of supply channels. The
// BatchSupplyCh will work as a bidirectional channel to request for a []BatchItems from the
// Workerline and gets the batch items from the response channel. The ClientSupplyCh will send
// received the []BatchItems from the BatchSupplyCh to the client. 
func NewBatchSupply() *BatchSupply{
	return &BatchSupply {
		BatchSupplyCh: make(chan chan []BatchItems, 100),
		ClientSupplyCh: make(chan []BatchItems, 1),
	}
}

// StartConsumer will create the Wokerpool [DefaultWorkerPool: 10] to handle the large set of 
// []BatchItems that gets created fequently in highly concurrent scenarios. Also, starts the 
// ConsumerCh channel listener to the incoming []BatchItems from the Producer line. 
//
//  signal.Notify(c.TerminateCh, syscall.SIGINT, syscall.SIGTERM)
//  <-c.TerminateCh
//
// To handle the graceful shutdown, the BatchConsumer supports os.Signal. So, the TerminateCh
// works as a terminate channel in case of certain os.Signal received [syscall.SIGINT, syscall.SIGTERM].
// This logic will help the Workerline to complete the remaining work before going for a shutdown.
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

// ConsumerFunc works as a callback function for the Producer line to send the released []BatchItems
// to the Consumer and then the batch items send to the ConsumerCh channel for further processing.
func (c *BatchConsumer) ConsumerFunc(items []BatchItems) {
	c.ConsumerCh <- items
}

// ConsumerBatch has the <-c.ConsumerCh receive channel to receives the newly created []BatchItems.
// After that, the []BatchItems gets send to the WorkerCh to send the batch item to the supply line.
//
// This also supports the termination of the Consumer line in case of graceful shutdown or to exit
// the batch processing forcefully. 
//      
//  <-ctx.Done(): get called during a gracefull shutdown scenarios and closes the worker channel
//  <-c.Quit: Exit the batch processing during a forceful request from the client.
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

// WorkerFunc is the final production of []BatchItems. Each WorkerChannel sends their released
// []BatchItems to the SupplyChannel. 
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