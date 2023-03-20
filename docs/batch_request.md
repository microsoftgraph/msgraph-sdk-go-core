## BatchRequest

BatchRequest is useful when you want to make multiple requests efficiently. It batches all requests (upto 20 requests) into a json object and makes one api call. You can learn more about it on [Microsoft Docs](https://docs.microsoft.com/en-us/graph/json-batching). 

## Code Sample

```go
import "github.com/microsoftgraph/msgraph-sdk-go-core"
import abstractions "github.com/microsoft/kiota-abstractions-go"

reqInfo := client.Me().CreateGetRequestInformation()
batch := msgraphsdkcore.NewBatchRequest()
batchItem := batch.AppendBatchItem(*reqInfo)

resp, err := batch.Send(reqAdapter)

// print the first response
user := GetBatchResponseById[User](resp, "1", CreateUserFromDiscriminatorValue) // returns a serialized response
fmt.Println(user.GetDisplayName()) // Print display name
```

## Depends On Relationship

BatchItem supports constructing a dependency chain for scenarios where you want one request to be sent out before another request is made. In the example below batchItem2 will be sent before batchItem1.

```go
batchItem1 := batch.AppendBatchItem(*reqInfo)
batchItem2 := batch.AppendBatchItem(*reqInfo)

batchItem1.DependsOnItem(batchItem2)
```

## Adds BatchCollectionResponse

`BatchRequestCollection` allows users to add more than 19 requests and send them as multiple `BatchRequest`'s. The send functionality of BatchRequestCollection splits the requests and sends them in serial.

```go
batchCollection := msgraphgocore.NewBatchRequestCollection(client.GetAdapter())

meRequestItem, _ := batchCollection.AddBatchRequestStep(*meRequest)
eventsRequestItem, _ := batchCollection.AddBatchRequestStep(*eventsRequest)

batchResponse, _ := batchCollection.Send(context.Background(), client.GetAdapter())

// print the first response
user := GetBatchResponseById[User](resp, "1", CreateUserFromDiscriminatorValue) // returns a serialized response
fmt.Println(user.GetDisplayName()) // Print display name
```