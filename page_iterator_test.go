package msgraphgocore

import (
	"fmt"
	nethttp "net/http"
	httptest "net/http/httptest"
	testing "testing"

	abstractions "github.com/microsoft/kiota/abstractions/go"
	"github.com/microsoft/kiota/abstractions/go/authentication"
	"github.com/microsoft/kiota/abstractions/go/serialization"
	jsonserialization "github.com/microsoft/kiota/serialization/go/json"
	"github.com/microsoftgraph/msgraph-sdk-go-core/internal"
	"github.com/stretchr/testify/assert"
)

func init() {
	abstractions.RegisterDefaultSerializer(func() serialization.SerializationWriterFactory {
		return jsonserialization.NewJsonSerializationWriterFactory()
	})
	abstractions.RegisterDefaultDeserializer(func() serialization.ParseNodeFactory {
		return jsonserialization.NewJsonParseNodeFactory()
	})
}

var reqAdapter, _ = NewGraphRequestAdapterBase(&authentication.AnonymousAuthenticationProvider{}, GraphClientOptions{
	GraphServiceVersion:        "",
	GraphServiceLibraryVersion: "",
})

func ParsableCons(pn serialization.ParseNode) (serialization.Parsable, error) {
	return internal.NewUsersResponse(), nil
}

func TestConstructorWithInvalidGraphResponse(t *testing.T) {
	graphResponse := internal.NewUsersResponse()

	_, err := NewPageIterator(graphResponse, reqAdapter, ParsableCons)

	assert.NotNil(t, err)
}

func TestIterateStopsWhenCallbackReturnsFalse(t *testing.T) {
	res := make([]string, 0)
	graphResponse := buildGraphResponse()
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"nextLink": "",
				"value": [
	        		{
	            		"id": "10"
	        		}
	        	]
        	}
        `)

	}))
	defer testServer.Close()
	pageIterator, _ := NewPageIterator(graphResponse, reqAdapter, ParsableCons)

	pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(internal.User)

		res = append(res, *item.GetDisplayName())
		return !(*item.GetId() == "2")
	})

	assert.Equal(t, len(res), 3)
}

func TestIterateEnumeratesAllPages(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"nextLink": "",
				"value": [
	        		{
	            		"id": "10"
	        		}
	        	]
        	}
        `)

	}))
	defer testServer.Close()

	graphResponse := buildGraphResponse()
	mockPath := testServer.URL + "/next-page"
	graphResponse.SetNextLink(&mockPath)

	pageIterator, _ := NewPageIterator(graphResponse, reqAdapter, ParsableCons)
	res := make([]string, 0)

	err := pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(internal.User)
		res = append(res, *item.GetId())
		return true
	})

	// Initial page has 5 items and the next page has 1 item.
	assert.Equal(t, len(res), 6)
	assert.Nil(t, err)
}

func TestIterateCanBePausedAndResumed(t *testing.T) {
	res := make([]string, 0)
	res2 := make([]string, 0)

	testServer := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"nextLink": "",
				"value": [
	        		{
	            		"id": "10"
	        		}
	        	]
        	}
        `)

	}))
	defer testServer.Close()

	response := buildGraphResponse()
	mockPath := testServer.URL + "/next-page"
	response.SetNextLink(&mockPath)

	pageIterator, _ := NewPageIterator(response, reqAdapter, ParsableCons)
	pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(internal.User)
		res = append(res, *item.GetId())

		return *item.GetId() != "4"
	})

	assert.Equal(t, res, []string{"0", "1", "2", "3", "4"})

	pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(internal.User)
		res2 = append(res2, *item.GetId())

		return true
	})
	assert.Equal(t, res2, []string{"10"})
}

func buildGraphResponse() *internal.UsersResponse {
	var res = internal.NewUsersResponse()

	nextLink := "next-page"
	users := make([]internal.User, 0)

	for i := 0; i < 5; i++ {
		u := internal.NewUser()
		id := fmt.Sprint(i)
		u.SetId(&id)

		users = append(users, *u)
	}

	res.SetNextLink(&nextLink)
	res.SetValue(users)

	return res
}
