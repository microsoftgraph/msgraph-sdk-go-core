package msgraphgocore

import (
	"net/url"
	"strings"
	"testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/assert"
)

func TestJSONBody(t *testing.T) {
	reqInfo := getRequestInfo()

	batch := NewBatchRequest()
	item, _ := batch.AppendItem(*reqInfo)
	item.Id = "1"

	expected := `
    {
        "requests": [
            {
                "method": "GET",
                "url": "",
                "body": {
                    "username": "name"
                },
                "headers": {
                    "content-type": "application/json"
                }
            }
        ]
    }
    `
	actual, _ := batch.toJson()

	expectedWithoutNewLine := strings.Replace(expected, "\n", "", -1)
	expectedWithoutSpace := strings.Replace(expectedWithoutNewLine, " ", "", -1)

	assert.Equal(t, expectedWithoutSpace, string(actual))
}
func TestDependsOnRelationship(t *testing.T) {}

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
