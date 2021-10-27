package msgraphgocore

import (
	nethttp "net/http"
	httptest "net/http/httptest"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
)

func TestItReplacesQueryParameters(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(res nethttp.ResponseWriter, req *nethttp.Request) {
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer func() { testServer.Close() }()
	handler := NewGraphODataQueryHandler()
	handler.SetNext(&NoopMiddleware{})
	req, err := nethttp.NewRequest(nethttp.MethodGet, testServer.URL+"/?Select=something&exPand=somethingElse(select=nested)&$top=10", nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := handler.Do(req, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, resp)
	query := req.URL.Query()
	assert.Equal(t, "something", query.Get("$Select"))
	assert.Equal(t, "somethingElse($select=nested)", query.Get("$exPand"))
	assert.Equal(t, "10", query.Get("$top"))
}
