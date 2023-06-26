package main

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
	batch "github.com/wlach/go-batch"
	batchLog "github.com/wlach/go-batch/logger"
)

// Resources Structure
type Resources struct {
	id   int
	name string
	flag bool
}

func main() {

	var rFlag, mFlag int
	flag.IntVar(&rFlag, "r", 10, "No of resources")
	flag.IntVar(&mFlag, "m", 10, "Maximum items")
	debugFlag := flag.Bool("d", false, "Debug mode")
	flag.Parse()

	logLevel := batchLog.Info
	if *debugFlag {
		logLevel = batchLog.Debug
	}

	logs := log.New()

	logs.Infoln("Batch Processing Example !")

	b := batch.NewBatch(batch.WithMaxItems(uint64(mFlag)), batch.WithLogLevel(logLevel))

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
