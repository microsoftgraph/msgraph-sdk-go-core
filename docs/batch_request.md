## BatchRequest

BatchRequest is useful when you want to make multiple requests efficiently. It batches all requests (upto 20 requests) into a json object and makes one api call. You can learn more about it on [Microsoft Docs](https://docs.microsoft.com/en-us/graph/json-batching). 

## Code Sample

```go
reqInfo := abstractions.NewRequestInformation()
batch := msgraphsdkcore.NewBatchRequest()
batchItem := batch.AppendBatchItem(*reqInfo)

resp, err := msgraphsdkcore.SendBatch(*batch, reqAdapter)
```

## Depends On Relationship

BatchItem supports constructing a dependency chain for scenarios where you want one request to be sent out before another request is made. In the example below batchItem2 will be sent before batchItem1.

```go
batchItem1 := batch.AppendBatchItem(*reqInfo)
batchItem2 := batch.AppendBatchItem(*reqInfo)

batchItem1.DependsOnItem(batchItem2)
```
