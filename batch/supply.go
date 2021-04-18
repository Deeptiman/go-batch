package batch

import (
	log "github.com/sirupsen/logrus"
)

// GetBatchSupply request the WorkerChannel for the released []BatchItems. The BatchSupplyChannel
// works as a bidirectional channel to request/response for the final []BatchItems product.
// The ClientSupplyChannel will send the []BatchItems to the client.
func (c *BatchConsumer) GetBatchSupply() {

	supplyCh := make(chan []BatchItems)

	defer close(supplyCh)

	c.Supply.BatchSupplyCh <- supplyCh

		select {
		case supply := <-supplyCh:
			c.Log.WithFields(log.Fields{"Supply": len(supply)}).Warn("BatchSupply")
		
			c.Supply.ClientSupplyCh<- supply
		}
}