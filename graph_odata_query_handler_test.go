package msgraphgocore

import (
	nethttp "net/http"
	httptest "net/http/httptest"
	"net/url"
	"strings"
	testing "testing"

	abs "github.com/microsoft/kiota/abstractions/go"
	absauth "github.com/microsoft/kiota/abstractions/go/authentication"
	assert "github.com/stretchr/testify/assert"
)

func TestItReplacesQueryParameters(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(res nethttp.ResponseWriter, req *nethttp.Request) {
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer func() { testServer.Close() }()
	handler := NewGraphODataQueryHandler()
	req, err := nethttp.NewRequest(nethttp.MethodGet, testServer.URL+"/?Select=something&exPand=somethingElse(select=nested)&$top=10", nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := handler.Intercept(newNoopPipeline(), req)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, resp)
	query := req.URL.Query()
	assert.Equal(t, "something", query.Get("$Select"))
	assert.Equal(t, "somethingElse($select=nested)", query.Get("$exPand"))
	assert.Equal(t, "10", query.Get("$top"))
}

func TestItDoesNotReplaceWithLocalQueryOptions(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(res nethttp.ResponseWriter, req *nethttp.Request) {
		if strings.Contains(req.URL.RawQuery, "$") {
			t.Error("Query parameter $ was added")
		}
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer func() { testServer.Close() }()
	auth := &absauth.AnonymousAuthenticationProvider{}
	requestAdapter, err := NewGraphRequestAdapterBase(auth, GraphClientOptions{
		GraphServiceVersion:        "na",
		GraphServiceLibraryVersion: "na",
	})
	if err != nil {
		t.Error(err)
	}
	absRequest := abs.NewRequestInformation()
	targetUrl, err := url.Parse(testServer.URL + "/?Select=something")
	if err != nil {
		t.Error(err)
	}
	absRequest.SetUri(*targetUrl)
	absRequest.Method = abs.GET
	absRequest.AddRequestOptions(&GraphODataQueryHandlerOptions{
		ShouldReplace: func(*nethttp.Request) bool { return false },
	})
	requestAdapter.SendNoContentAsync(*absRequest, nil)
}
