package msgraphgocore_test

import (
	"fmt"
	nethttp "net/http"
	httptest "net/http/httptest"
	testing "testing"

	msgraphgocore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go-core/internal"

	"github.com/microsoft/kiota/abstractions/go/authentication"
	"github.com/microsoft/kiota/abstractions/go/serialization"
	"github.com/stretchr/testify/assert"
)

type PageItem struct {
	DisplayName string
}

type UserPage struct {
	Value    []interface{}
	NextLink *string
}

var reqAdapter, _ = msgraphgocore.NewGraphRequestAdapterBase(&authentication.AnonymousAuthenticationProvider{}, msgraphgocore.GraphClientOptions{
	GraphServiceVersion:        "",
	GraphServiceLibraryVersion: "",
})

func ParsableCons() serialization.Parsable {
	return internal.NewUsersResponse()
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
	pageIterator := msgraphgocore.NewPageIterator(graphResponse, *reqAdapter, ParsableCons, nil)

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

	pageIterator := msgraphgocore.NewPageIterator(graphResponse, *reqAdapter, ParsableCons, nil)
	res := make([]string, 0)

	pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(internal.User)
		res = append(res, *item.GetId())
		return true
	})

	// Initial page has 5 items and the next page has 1 item.
	assert.Equal(t, len(res), 6)
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

	pageIterator := msgraphgocore.NewPageIterator(response, *reqAdapter, ParsableCons, nil)
	pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(internal.User)
		res = append(res, *item.GetId())

		if *item.GetId() == "2" {
			return false
		}
		return true
	})
	assert.Equal(t, res, []string{"0", "1", "2"})

	pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(internal.User)
		res2 = append(res2, *item.GetId())

		return true
	})
	assert.Equal(t, res2, []string{"2", "3", "4", "10"})
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
