package batch

import "fmt"

func (c *BatchConsumer) GetBatchSupply() {

	supplyCh := make(chan []BatchItems)

	defer close(supplyCh)

	c.Supply.BatchSupplyCh <- supplyCh

		select {
		case supply := <-supplyCh:
			fmt.Println("Supply - ", supply)
			c.Supply.ClientSupplyCh<- supply
		}
}