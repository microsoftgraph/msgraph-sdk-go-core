package msgraphgocore

import (
	"net/url"
	"testing"

	abstractions "github.com/microsoft/kiota/abstractions/go"
	"github.com/stretchr/testify/assert"
)

func TestJSONBody(t *testing.T) {
	reqInfo := getRequestInfo()

	batch := NewBatchRequest()
	item, _ := batch.AppendItem(*reqInfo)
	item.Id = "1"

	expected := `{"requests":[{"id":"1","method":"GET","url":"","headers":{"content-type":"application/json"},"body":{"username":"name"},"dependsOn":[]}]}`
	actual, _ := batch.toJson()
	assert.Equal(t, expected, string(actual))
}
func TestDependsOnRelationship(t *testing.T) {
	reqInfo := getRequestInfo()

	batch := NewBatchRequest()
	item1, _ := batch.AppendItem(*reqInfo)
	item2, _ := batch.AppendItem(*reqInfo)
	item2.DependsOnItem(*item1)

	assert.Equal(t, item2.DependsOn[0], item1.Id)
}

func getRequestInfo() *abstractions.RequestInformation {
	content := `{"username":"name"}`
	reqInfo := abstractions.NewRequestInformation()
	reqInfo.SetUri(url.URL{})
	reqInfo.Content = []byte(content)
	reqInfo.Headers = map[string]string{"content-type": "application/json"}

	return reqInfo
}
