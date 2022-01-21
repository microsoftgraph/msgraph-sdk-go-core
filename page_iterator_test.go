package msgraphgocore

import (
	testing "testing"

	"github.com/microsoft/kiota/abstractions/go/serialization"
	"github.com/microsoftgraph/msgraph-sdk-go-core/msgraphgocore_test"
	"github.com/stretchr/testify/assert"
)

type PageItem struct {
	DisplayName string
}

type UserPage struct {
	Value    []interface{}
	NextLink *string
}

var page2Link = "page2"
var page3Link = "page3"

var page1 = &UserPage{
	Value: []interface{}{
		PageItem{DisplayName: "User A"},
		PageItem{DisplayName: "User B"},
		PageItem{DisplayName: "User C"},
		PageItem{DisplayName: "User D"},
	},
	NextLink: &page2Link,
}

var page2 = &UserPage{
	Value: []interface{}{
		PageItem{DisplayName: "User E"},
		PageItem{DisplayName: "User F"},
		PageItem{DisplayName: "User G"},
		PageItem{DisplayName: "User H"},
	},
	NextLink: &page3Link,
}

var page3 = &UserPage{
	Value: []interface{}{
		PageItem{DisplayName: "User I"},
		PageItem{DisplayName: "User J"},
		PageItem{DisplayName: "User K"},
		PageItem{DisplayName: "User L"},
	},
	NextLink: nil,
}

var reqAdapter = &GraphRequestAdapterBase{}

func ParsableCons() serialization.Parsable {
	return msgraphgocore_test.NewUsersResponse()
}

var res = msgraphgocore_test.NewUsersResponse()

func TestHasNextReturnsTrueIfNextPageIsAvailable(t *testing.T) {
	pageIterator := NewPageIterator(res, *reqAdapter, ParsableCons)
	hasNext := pageIterator.HasNext()
	assert.True(t, hasNext)
}

func TestNextReturnsNextPage(t *testing.T) {
	pageIterator := NewPageIterator(res, *reqAdapter, ParsableCons)
	nextPage := pageIterator.Next()

	assert.Equal(t, nextPage, page2)
}

func TestIterateStopsWhenCallbackReturnsFalse(t *testing.T) {
	res := make([]string, 0)

	pageIterator := NewPageIterator(res, *reqAdapter, ParsableCons)
	pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(msgraphgocore_test.User)

		res = append(res, *item.GetDisplayName())
		return *item.GetDisplayName() == "User D"
	})

	assert.Equal(t, 4, len(res))
}

func TestNextAndHasNext(t *testing.T) {
	pageCount := 1
	pageIterator := NewPageIterator(res, *reqAdapter, ParsableCons)

	for pageIterator.HasNext() {
		pageCount += 1
		_ = pageIterator.Next()
	}

	assert.Equal(t, pageCount, 3)
}
func TestIterateEnumeratesAllPages(t *testing.T) {
	pageIterator := NewPageIterator(res, *reqAdapter, ParsableCons)
	res := make([]string, 0)

	pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(msgraphgocore_test.User)
		res = append(res, *item.GetDisplayName())
		return true
	})

	assert.Equal(t, 12, len(res))
}

func TestIterateCanBePausedAndResumed(t *testing.T) {
	res := make([]string, 0)
	res2 := make([]string, 0)
	pageIterator := NewPageIterator(res, *reqAdapter, ParsableCons)

	pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(msgraphgocore_test.User)

		res = append(res, *item.GetDisplayName())

		return *item.GetDisplayName() == "User D"
	})
	assert.Equal(t, res, []string{"User A", "User B", "User C", "User D"})

	pageIterator.Iterate(func(pageItem interface{}) bool {
		item := pageItem.(msgraphgocore_test.User)
		res2 = append(res2, *item.GetDisplayName())

		return *item.GetDisplayName() == "User G"
	})
	assert.Equal(t, res2, []string{"User D", "User E", "User F", "User G"})
}
