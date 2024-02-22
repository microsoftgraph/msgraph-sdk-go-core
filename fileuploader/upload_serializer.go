package fileuploader

import (
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"io"
	nethttp "net/http"
	"strconv"
	"strings"
)

// UploadSerializer TODO allowing delegation of parsing an object to net_adapter should deprecate this class
type UploadSerializer struct {
	response        *nethttp.Response
	errorMappings   abstractions.ErrorMappings
	parsableFactory serialization.ParsableFactory
}

// NewUploadSerializer creates a new UploadSerializer object
func NewUploadSerializer(response *nethttp.Response, errorMappings abstractions.ErrorMappings, parsableFactory serialization.ParsableFactory) *UploadSerializer {
	return &UploadSerializer{
		response:        response,
		errorMappings:   errorMappings,
		parsableFactory: parsableFactory,
	}
}

func ParseResponse(u *UploadSerializer) (interface{}, *string, error) {
	err := u.throwIfFailedResponse()
	if err != nil {
		return nil, nil, err
	}

	if u.response.StatusCode == 204 {
		return nil, nil, nil
	}

	parseNode, err := u.getRootParseNode()
	if err != nil {
		return nil, nil, err
	}
	if parseNode == nil {
		return nil, nil, nil
	}

	result, err := parseNode.GetObjectValue(u.parsableFactory)
	if err != nil {
		return nil, nil, err
	}

	var urlLocation *string = nil
	location := u.response.Header.Get("Location")
	if location != "" {
		urlLocation = &location
	}

	return result, urlLocation, err
}

func (u *UploadSerializer) throwIfFailedResponse() error {
	statusAsString := strconv.Itoa(u.response.StatusCode)

	var errorCtor serialization.ParsableFactory = nil
	if len(u.errorMappings) != 0 {
		if u.errorMappings[statusAsString] != nil {
			errorCtor = u.errorMappings[statusAsString]
		} else if u.response.StatusCode >= 400 && u.response.StatusCode < 500 && u.errorMappings["4XX"] != nil {
			errorCtor = u.errorMappings["4XX"]
		} else if u.response.StatusCode >= 500 && u.response.StatusCode < 600 && u.errorMappings["5XX"] != nil {
			errorCtor = u.errorMappings["5XX"]
		} else if u.errorMappings["XXX"] != nil && u.response.StatusCode >= 400 && u.response.StatusCode < 600 {
			errorCtor = u.errorMappings["XXX"]
		}
	}

	responseHeaders := abstractions.NewResponseHeaders()
	for key, values := range u.response.Header {
		for i := range values {
			responseHeaders.Add(key, values[i])
		}
	}

	if errorCtor == nil {
		err := &abstractions.ApiError{
			Message:            "The server returned an unexpected status code and no error factory is registered for this code: " + statusAsString,
			ResponseStatusCode: u.response.StatusCode,
			ResponseHeaders:    responseHeaders,
		}
		return err
	}

	rootNode, err := u.getRootParseNode()
	if err != nil {
		return err
	}
	if rootNode == nil {
		err := &abstractions.ApiError{
			Message:            "The server returned an unexpected status code with no response body: " + statusAsString,
			ResponseStatusCode: u.response.StatusCode,
			ResponseHeaders:    responseHeaders,
		}
		return err
	}

	errValue, err := rootNode.GetObjectValue(errorCtor)
	if err != nil {
		if apiErrorable, ok := err.(abstractions.ApiErrorable); ok {
			apiErrorable.SetResponseHeaders(responseHeaders)
			apiErrorable.SetStatusCode(u.response.StatusCode)
		}
		return err
	} else if errValue == nil {
		return &abstractions.ApiError{
			Message:            "The server returned an unexpected status code but the error could not be deserialized: " + statusAsString,
			ResponseStatusCode: u.response.StatusCode,
			ResponseHeaders:    responseHeaders,
		}
	}

	if apiErrorable, ok := errValue.(abstractions.ApiErrorable); ok {
		apiErrorable.SetResponseHeaders(responseHeaders)
		apiErrorable.SetStatusCode(u.response.StatusCode)
	}

	err = errValue.(error)

	return err
}

func (u *UploadSerializer) getRootParseNode() (serialization.ParseNode, error) {
	if u.response.ContentLength == 0 {
		return nil, nil
	}

	body, err := io.ReadAll(u.response.Body)
	if err != nil {
		return nil, err
	}
	contentType := u.getResponsePrimaryContentType(u.response)
	if contentType == "" {
		return nil, nil
	}
	rootNode, err := serialization.DefaultParseNodeFactoryInstance.GetRootParseNode(contentType, body)
	return rootNode, err
}

func (u *UploadSerializer) getResponsePrimaryContentType(response *nethttp.Response) string {
	if response.Header == nil {
		return ""
	}
	rawType := response.Header.Get("Content-Type")
	splat := strings.Split(rawType, ";")
	return strings.ToLower(splat[0])
}
