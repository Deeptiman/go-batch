package main

import (
	"flag"
	"fmt"
	//"time"
	log "github.com/sirupsen/logrus"
	batch "go-batch/batch"
)

// Resources Structure
type Resources struct {
	id   int
	name string
	flag bool
}

func main() {

	fmt.Println(" Job Queue Example !")

	var rFlag, mFlag int
	flag.IntVar(&rFlag, "r", 10, "No of resources")
	flag.IntVar(&mFlag, "m", 10, "Maximum items")
	flag.Parse()

	logs := log.New()

	logs.WithFields(log.Fields{"Batch": "Items"}).Info("Batch System")

	b := batch.NewBatch(batch.WithMaxItems(uint64(mFlag)))

	b.StartBatchProcessing()

	// go func() {

	// 	select {
	// 	case <-time.After(time.Duration(2) * time.Second):
	// 		fmt.Println("Run Batch processing again!")
	// 		b.StartBatchProcessing()
	// 	}

	// }()

	go func() {

		// Infinite loop to listen to the Consumer Client Supply Channel that releases
		// the []BatchItems for each iteration.
		for {
			for bt := range b.Consumer.Supply.ClientSupplyCh {
				logs.WithFields(log.Fields{"Batch": bt}).Warn("Client")
			}
		}
	}()

	for i := 1; i <= rFlag; i++ {
		b.Item <- &Resources{
			id:   i,
			name: fmt.Sprintf("%s%d", "R-", i),
			flag: false,
		}
	}
	b.Close()
}
