package msgraphgocore

import (
	"context"
	abstractions "github.com/microsoft/kiota-abstractions-go"
)

type BatchRequestCollection struct {
	*batchRequest
}

// NewBatchRequestCollection creates an instance of a BatchRequestCollection
func NewBatchRequestCollection(adapter abstractions.RequestAdapter) *BatchRequestCollection {
	return &BatchRequestCollection{
		&batchRequest{
			adapter: adapter,
		},
	}
}

// AddBatchRequestStep converts RequestInformation to a BatchItem and adds it to a BatchRequestCollection
func (b *BatchRequestCollection) AddBatchRequestStep(reqInfo abstractions.RequestInformation) (BatchItem, error) {
	return b.AddLimitedBatchRequestStep(reqInfo, -1)
}

// Send serializes and sends the batch request to the server
func (b *BatchRequestCollection) Send(ctx context.Context, adapter abstractions.RequestAdapter) (BatchResponse, error) {
	// spit request with a max of 19
	requestItems := chunkSlice(b.requests, 19)

	// execute requests
	response := NewBatchResponse()
	for _, requests := range requestItems {
		batch := NewBatchRequest(b.adapter)
		batch.SetRequests(requests)
		res, err := batch.Send(ctx, adapter)
		if err != nil {
			return nil, err
		}
		response.AddResponses(res.GetResponses())
	}

	return response, nil
}

func chunkSlice[T interface{}](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
