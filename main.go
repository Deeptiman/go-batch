package main

import (
	"flag"
	"fmt"
	batch "go-batch/batch"
	log "github.com/sirupsen/logrus"
)

// Resources Structure
type Resources struct {
	id 		int
	name 	string
	flag 	bool
}

func main() {

	fmt.Println(" Job Queue Example !")

	var rFlag, mFlag int
	flag.IntVar(&rFlag, "r", 10, "No of resources")
	flag.IntVar(&mFlag, "m", 10, "Maximum items")
	flag.Parse()

	logs := log.New()

	logs.WithFields(log.Fields{"Batch":"Items"}).Info("Batch System")

	b := batch.NewBatch(batch.WithMaxItems(uint64(mFlag)))

	go b.StartBachProcessing()

	go func() {

		for {
			for bt := range b.Consumer.Supply.ClientSupplyCh {
				fmt.Println("Resources", "Batch Supply - "," : ",bt)
			}
		}
	}()

	for i := 1; i <= rFlag; i++ {
		b.Item <- &Resources{
			id:   i,
			name: fmt.Sprintf("%s%d","R-",  i),
			flag: false,
		}
	}
	b.Close()
}