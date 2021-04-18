# go-batch
go-batch is a batch processing library written in Go. The process execution has multiple stages to release a Batch to the client.

## Features

1. Client can use this library as an asynchronous batch processing for their application use case.
2. There are no restrictions on applying batch processing matrices to the library. The client can define the maximum no of items for a batch using the <b>BatchOptions</b>.
3. The library has a Workerpool that will faster the batch processing in concurrent scenarios.

## Stages

 1. Batch Reader receives the resource payload from the client and marshals the payload item into the BatchItem object.

    ``````````````````````````
    type BatchItems struct {
        Id      int
        BatchNo int
        Item    interface{}
    } 

    ``````````````````````````
2. BatchProducer has a <b>Watcher</b> channel that receives the marshal payload from the Batch reader. Watcher marks each BatchItem with a BatchNo and adds it to the []BatchItems array. After the batch itemCounter++ increases to the MaxItems [DefaultMaxItems: 100], the Batch gets
releases to the Consumer callback function.

3. BatchConsumer has a <b>ConsumerFunc</b> that gets invoke by BatchProducer as a callback function to send the prepared []BatchItems arrays. Then, the Consumer channel sends the 
[]BatchItems to the Worker channel.

4. Workerline is the sync.WaitGroup synchronizes the workers to send the []BatchItems to the supply chain.

5. BatchSupplyChannel works as a bidirectional channel that requests for the []BatchItems to the Workerline and gets in the response.

6. ClientSupplyChannel is the delivery channel that works as a Supply line to sends the []BatchItems and the client receives by listening to the channel.


## Go Docs

   Documentation at <a href="https://pkg.go.dev/github.com/Deeptiman/go-batch">pkg.go.dev</a>

## Installation

    go get github.com/Deeptiman/go-batch

## Example

   ```````````````````````````````````````````````````` 
   b := batch.NewBatch(batch.WithMaxItems(100))
   go b.StartBatchProcessing()

   for i := 1; i <= 1000; i++ {
        b.Item <- &Resources{
            id:   i,
            name: fmt.Sprintf("%s%d", "R-", i),
            flag: false,
        }
    }
   b.Close() 
  ````````````````````````````````````````````````````
  [![asciicast](https://asciinema.org/a/2vi5gAHjsuTrB3tCBTGeSW6hq.svg)](https://asciinema.org/a/2vi5gAHjsuTrB3tCBTGeSW6hq)
  
## Note
- In this version release, the library doesn't support starting concurrent BatchProcessing sessions. 

## License

This project is licensed under the <a href="https://github.com/Deeptiman/go-batch/blob/main/LICENSE">MIT License</a>
