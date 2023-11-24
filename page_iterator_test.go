package msgraphgocore

import (
	"context"
	"fmt"
	nethttp "net/http"
	httptest "net/http/httptest"
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

	pageIterator, _ := NewPageIterator[internal.User](graphResponse, reqAdapter, ParsableCons)
	res := make([]string, 0)

	err := pageIterator.Iterate(context.Background(), func(item internal.User) bool {
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

	response := buildGraphResponse()
	mockPath := testServer.URL + "/next-page"
	response.SetOdataNextLink(&mockPath)

	pageIterator, _ := NewPageIterator[internal.User](response, reqAdapter, ParsableCons)
	pageIterator.Iterate(context.Background(), func(item internal.User) bool {
		res = append(res, *item.GetId())

		return *item.GetId() != "4"
	})

	assert.Equal(t, res, []string{"0", "1", "2", "3", "4"})
	assert.Equal(t, pageIterator.GetOdataNextLink(), response.GetOdataNextLink())

	pageIterator.Iterate(context.Background(), func(item internal.User) bool {
		res2 = append(res2, *item.GetId())

		return true
	})
	assert.Equal(t, res2, []string{"10"})
	assert.Empty(t, pageIterator.GetOdataNextLink())

	pageIterator.Iterate(context.Background(), func(item internal.User) bool {
		assert.Fail(t, "Should not re-iterate over items")
		return true
	})
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

func Test_convertToPage(t *testing.T) {
	type args struct {
		response interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    PageResult[validTestStruct]
		wantErr bool
	}{
		{
			name: "should pass",
			args: args{
				response: &validTestStruct{
					obj: []validTestStruct{
						{
							obj: []validTestStruct{
								{
									obj: []validTestStruct{},
								},
							},
						},
					}},
			},
			want: PageResult[validTestStruct]{
				oDataNextLink:  nil,
				oDataDeltaLink: nil,
				value: []validTestStruct{
					{
						obj: []validTestStruct{
							{
								obj: []validTestStruct{},
							},
						},
					},
				},
			},
		},
		{
			name: "should return error 'saying response cannot be nil' for nil response",
			args: args{
				response: nil,
			},
			want: PageResult[validTestStruct]{
				oDataNextLink:  nil,
				oDataDeltaLink: nil,
				value:          nil,
			},
			wantErr: true,
		},
		{
			name: "should return error 'value property missing in response object' for missing 'GetValue' method",
			args: args{
				response: &invalidTestStruct{
					obj: []invalidTestStruct{
						{
							obj: []invalidTestStruct{
								{
									obj: []invalidTestStruct{},
								},
							},
						},
					}},
			},
			want: PageResult[validTestStruct]{
				oDataNextLink:  nil,
				oDataDeltaLink: nil,
			},
			wantErr: true,
		},
		{
			name: "should return error 'response does not have next link accessor' for missing 'GetOdataNextLink() *string' method",
			args: args{
				response: &invalidTestStruct{
					obj: []invalidTestStruct{
						{
							obj: []invalidTestStruct{
								{
									obj: []invalidTestStruct{},
								},
							},
						},
					}},
			},
			want: PageResult[validTestStruct]{
				oDataNextLink:  nil,
				oDataDeltaLink: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToPage[validTestStruct](tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.oDataDeltaLink, got.oDataDeltaLink, "got %v, want %v", got.oDataNextLink, tt.want.oDataNextLink)
			assert.Equal(t, tt.want.oDataNextLink, got.oDataNextLink, "got %v, want %v", got.oDataDeltaLink, tt.want.oDataDeltaLink)
			assert.Equal(t, tt.want.value, got.value, "got %v, want %v", got.value, tt.want.value)
		})
	}
}

type validTestStruct struct {
	obj []validTestStruct
}

func (t *validTestStruct) GetValue() []validTestStruct {
	return t.obj
}

func (t *validTestStruct) GetOdataNextLink() *string {
	return nil
}

func (t *validTestStruct) GetOdataDeltaLink() *string {
	return nil
}

type invalidTestStruct struct {
	obj []invalidTestStruct
}

func (t *invalidTestStruct) GetOdataDeltaLink() *string {
	return nil
}
