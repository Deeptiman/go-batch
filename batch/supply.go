package batch

import (
	log "github.com/sirupsen/logrus"
)

func (c *BatchConsumer) GetBatchSupply() {

	supplyCh := make(chan []BatchItems)

	defer close(supplyCh)

	c.Supply.BatchSupplyCh <- supplyCh

		select {
		case supply := <-supplyCh:
			c.Log.WithFields(log.Fields{"Supply": supply}).Info("BatchSupply")
		
			c.Supply.ClientSupplyCh<- supply
		}
}