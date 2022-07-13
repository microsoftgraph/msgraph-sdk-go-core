package msgraphgocore

import (
	"net/url"
	"testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/assert"
)

func TestGeneratesJSONFromRequestBody(t *testing.T) {
	reqInfo := getRequestInfo()

	batch := NewBatchRequest()
	item, _ := batch.AppendItem(*reqInfo)
	item.Id = "1"

	expected := "{\"requests\":[{\"id\":\"1\",\"method\":\"GET\",\"url\":\"\",\"headers\":{\"content-type\":\"application/json\"},\"body\":{\"username\":\"name\"},\"dependsOn\":[]}]}"
	actual, _ := batch.toJson()

	assert.Equal(t, expected, string(actual))
}

func TestDependsOnRelationshipInBatchRequestItems(t *testing.T) {

	reqInfo1 := getRequestInfo()
	reqInfo2 := getRequestInfo()

	batch := NewBatchRequest()
	batchItem1, _ := batch.AppendItem(*reqInfo1)
	batchItem2, _ := batch.AppendItem(*reqInfo2)
	batchItem1.Id = "1"
	batchItem2.Id = "2"

	batchItem2.DependsOnItem(*batchItem1)

	expected := "{\"requests\":[{\"id\":\"1\",\"method\":\"GET\",\"url\":\"\",\"headers\":{\"content-type\":\"application/json\"},\"body\":{\"username\":\"name\"},\"dependsOn\":[]},{\"id\":\"2\",\"method\":\"GET\",\"url\":\"\",\"headers\":{\"content-type\":\"application/json\"},\"body\":{\"username\":\"name\"},\"dependsOn\":[\"1\"]}]}"
	actual, _ := batch.toJson()

	assert.Equal(t, expected, string(actual))
}

func getRequestInfo() *abstractions.RequestInformation {
	content := `
    {
        "username": "name"
    }
    `
	reqInfo := abstractions.NewRequestInformation()
	reqInfo.SetUri(url.URL{})
	reqInfo.Content = []byte(content)
	reqInfo.Headers = map[string]string{"content-type": "application/json"}

	return reqInfo
}
