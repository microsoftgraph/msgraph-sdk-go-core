package msgraphgocore

import (
	"context"
	"fmt"
	nethttp "net/http"
	httptest "net/http/httptest"
	"strconv"
	testing "testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
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

func ParsableDeltaCons(pn serialization.ParseNode) (serialization.Parsable, error) {
	return internal.NewUsersDeltaResponse(), nil
}

func TestConstructorWithInvalidRequestAdapter(t *testing.T) {
	graphResponse := internal.NewUsersResponse()

	_, err := NewPageIterator[internal.User](graphResponse, nil, ParsableCons)

	assert.NotNil(t, err)
}

func TestConstructorWithInvalidGraphResponse(t *testing.T) {
	graphResponse := internal.NewInvalidUsersResponse()

	_, err := NewPageIterator[internal.User](graphResponse, reqAdapter, ParsableCons)

	assert.NotNil(t, err)
}

func TestConstructorWithInvalidUserGraphResponse(t *testing.T) {
	graphResponse := internal.NewInvalidUsersResponse()

	nextLink := "next-page"
	users := make([]internal.User, 0)

	graphResponse.SetNextLink(&nextLink)
	graphResponse.SetValue(users)

	_, err := NewPageIterator[internal.User](graphResponse, reqAdapter, ParsableCons)

	assert.NotNil(t, err)
}

func TestIterateStopsWhenCallbackReturnsFalse(t *testing.T) {
	res := make([]string, 0)
	graphResponse := buildGraphResponse()
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"@odata.nextLink": "",
				"value": [
	        		{
	            		"id": "10"
	        		}
	        	]
        	}
        `)
		assert.NotNil(t, req.Header["ConsistencyLevel"])
	}))
	defer testServer.Close()
	pageIterator, _ := NewPageIterator[internal.User](graphResponse, reqAdapter, ParsableCons)
	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")
	pageIterator.SetHeaders(headers)

	pageIterator.Iterate(context.Background(), func(item internal.User) bool {
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
				"@odata.nextLink": "",
				"value": [
	        		{
	            		"id": "5"
	        		}
	        	]
        	}
        `)

	}))
	defer testServer.Close()

	for _, tc := range []struct {
		description string
		useNext     bool
	}{
		{
			"using Next",
			true,
		},
		{
			"using Iterate",
			false,
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			graphResponse := buildGraphResponse()
			mockPath := testServer.URL + "/next-page"
			graphResponse.SetOdataNextLink(&mockPath)

			pageIterator, _ := NewPageIterator[internal.User](graphResponse, reqAdapter, ParsableCons)
			res := make([]string, 0)

			if tc.useNext {
				for pageIterator.HasNext() {
					item, err := pageIterator.Next(context.Background())
					assert.NoError(t, err)
					res = append(res, *item.GetId())
				}
			} else {
				err := pageIterator.Iterate(context.Background(), func(item internal.User) bool {
					res = append(res, *item.GetId())
					return true
				})
				assert.Nil(t, err)
			}

			assert.Equal(t, res, []string{"0", "1", "2", "3", "4", "5"})
		})
	}
}

func TestIterateCanBePausedAndResumed(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"@odata.nextLink": "",
				"value": [
					{
						"id": "10"
					}
				]
			}
		`)

	}))
	defer testServer.Close()

	for _, tc := range []struct {
		description string
		useNext     bool
	}{
		{
			"using Next",
			true,
		},
		{
			"using Iterate",
			false,
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			res := make([]string, 0)
			res2 := make([]string, 0)

			response := buildGraphResponse()
			mockPath := testServer.URL + "/next-page"
			response.SetOdataNextLink(&mockPath)

			pageIterator, _ := NewPageIterator[internal.User](response, reqAdapter, ParsableCons)
			if tc.useNext {
				for pageIterator.HasNext() {
					item, err := pageIterator.Next(context.Background())
					assert.NoError(t, err)

					res = append(res, *item.GetId())

					if *item.GetId() == "4" {
						break
					}
				}
			} else {
				pageIterator.Iterate(context.Background(), func(item internal.User) bool {
					res = append(res, *item.GetId())

					return *item.GetId() != "4"
				})
			}

			assert.Equal(t, res, []string{"0", "1", "2", "3", "4"})
			assert.Equal(t, pageIterator.GetOdataNextLink(), response.GetOdataNextLink())

			if tc.useNext {
				for pageIterator.HasNext() {
					item, err := pageIterator.Next(context.Background())
					assert.NoError(t, err)

					res2 = append(res2, *item.GetId())
				}
			} else {
				pageIterator.Iterate(context.Background(), func(item internal.User) bool {
					res2 = append(res2, *item.GetId())

					return true
				})
			}

			assert.Equal(t, res2, []string{"10"})
			assert.Empty(t, pageIterator.GetOdataNextLink())

			if tc.useNext {
				assert.False(t, pageIterator.HasNext())
				_, err := pageIterator.Next(context.Background())
				assert.Error(t, err)
			} else {
				pageIterator.Iterate(context.Background(), func(item internal.User) bool {
					assert.Fail(t, "Should not re-iterate over items")
					return true
				})
			}
		})
	}
}

func TestAllEnumeratesAllItems(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"@odata.nextLink": "",
				"value": [
	        		{
	            		"id": "5"
	        		}
	        	]
        	}
        `)

	}))
	defer testServer.Close()

	graphResponse := buildGraphResponse()
	mockPath := testServer.URL + "/next-page"
	graphResponse.SetOdataNextLink(&mockPath)

	pageIterator, err := NewPageIterator[internal.User](graphResponse, reqAdapter, ParsableCons)
	assert.NoError(t, err)

	res, err := pageIterator.All(context.Background())
	assert.NoError(t, err)

	assert.Len(t, res, 6)
	for i, r := range res {
		assert.Equal(t, strconv.Itoa(i), *r.GetId())
	}
}

func TestGetOdataNextLink(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"@odata.nextLink": "",
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
	graphResponse.SetOdataNextLink(&mockPath)

	pageIterator, _ := NewPageIterator[internal.User](graphResponse, reqAdapter, ParsableDeltaCons)
	pageIterator.Iterate(context.Background(), func(item internal.User) bool {
		return true
	})

	assert.Empty(t, pageIterator.GetOdataNextLink())
}

func TestGetOdataDeltaLink(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"@odata.nextLink": "",
				"@odata.deltaLink": "delta-page-2",
				"value": [
	        		{
	            		"id": "10"
	        		}
	        	]
        	}
        `)
	}))
	defer testServer.Close()

	dl := "delta-page-1"
	mockPath := testServer.URL + "/next-page"

	graphResponse := &internal.UsersDeltaResponse{
		UsersResponse: *buildGraphResponse(),
	}
	graphResponse.SetOdataDeltaLink(&dl)
	graphResponse.SetOdataNextLink(&mockPath)

	pageIterator, _ := NewPageIterator[internal.User](graphResponse, reqAdapter, ParsableDeltaCons)
	pageIterator.Iterate(context.Background(), func(item internal.User) bool {
		return true
	})

	assert.Equal(t, *pageIterator.GetOdataDeltaLink(), "delta-page-2")
}

func TestHasNext(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"@odata.nextLink": "",
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
	graphResponse.SetOdataNextLink(&mockPath)

	pageIterator, _ := NewPageIterator[internal.User](graphResponse, reqAdapter, ParsableDeltaCons)
	assert.True(t, pageIterator.HasNext())
	pageIterator.Iterate(context.Background(), func(item internal.User) bool {
		return false
	})

	assert.True(t, pageIterator.HasNext())
	pageIterator.Iterate(context.Background(), func(item internal.User) bool {
		return true
	})
	assert.False(t, pageIterator.HasNext())
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

	res.SetOdataNextLink(&nextLink)
	res.SetValue(users)

	return res
}
