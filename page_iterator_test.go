package msgraphgocore

//
//import (
//	abstractions "github.com/microsoft/kiota/abstractions/go"
//	absauth "github.com/microsoft/kiota/abstractions/go/authentication"
//	"github.com/microsoft/kiota/abstractions/go/serialization"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	u "net/url"
//	testing "testing"
//)
//
//type PageItem struct {
//	DisplayName string
//}
//
//type UserPage struct {
//	Value    []Item
//	NextLink *string
//}
//
//var page2Link = "page2"
//var page3Link = "page3"
//
//var page1 = &UserPage{
//	Value: []Item{
//		PageItem{DisplayName: "User A"},
//		PageItem{DisplayName: "User B"},
//		PageItem{DisplayName: "User C"},
//		PageItem{DisplayName: "User D"},
//	},
//	NextLink: &page2Link,
//}
//
//var page2 = &UserPage{
//	Value: []Item{
//		PageItem{DisplayName: "User E"},
//		PageItem{DisplayName: "User F"},
//		PageItem{DisplayName: "User G"},
//		PageItem{DisplayName: "User H"},
//	},
//	NextLink: &page3Link,
//}
//
//var page3 = &UserPage{
//	Value: []Item{
//		PageItem{DisplayName: "User I"},
//		PageItem{DisplayName: "User J"},
//		PageItem{DisplayName: "User K"},
//		PageItem{DisplayName: "User L"},
//	},
//	NextLink: nil,
//}
//
//func (p *UserPage) GetValue() []Item {
//	return p.Value
//}
//
//func (p *UserPage) GetNextLink() *string {
//	return p.NextLink
//}
//
//func (p *UserPage) GetNextPage() Page {
//	switch nextLink := p.GetNextLink(); *nextLink {
//	case "page_2":
//		return page2
//	case "page_3":
//		return page3
//	}
//	return nil
//}
//
//func (i PageItem) GetDisplayName() string {
//	return i.DisplayName
//}
//
//func TestHasNextReturnsTrueIfNextPageIsAvailable(t *testing.T) {
//	client := http.Client{}
//	pageIterator := NewPageIterator(page1, client)
//	hasNext := pageIterator.HasNext()
//
//	assert.True(t, hasNext)
//}
//
//func TestNextReturnsNextPage(t *testing.T) {
//	client := http.Client{}
//	pageIterator := NewPageIterator(page1, client)
//	nextPage := pageIterator.Next()
//
//	assert.Equal(t, nextPage, page2)
//}
//
//func TestIterateStopsWhenCallbackReturnsFalse(t *testing.T) {
//	res := make([]string, 0)
//
//	requestInformation := abstractions.NewRequestInformation()
//	requestInformation.SetUri(u.URL{})
//
//	client, _ := NewGraphRequestAdapterBase(absauth.AuthenticationProvider(abstractions.RequestInformation{
//		Method: abstractions.HttpMethod(0),
//	}), GraphClientOptions{
//		GraphServiceVersion:        "",
//		GraphServiceLibraryVersion: "",
//	})
//
//	pageIterator := NewPageIterator(page1, client)
//	pageIterator.Iterate(func(pageItem Item) bool {
//		res = append(res, pageItem.GetDisplayName())
//		if pageItem.GetDisplayName() == "User D" {
//			return false
//		}
//		return true
//	})
//
//	assert.Equal(t, 4, len(res))
//}
//
//func TestNextAndHasNext(t *testing.T) {
//	pageCount := 1
//	client := http.Client{}
//	pageIterator := NewPageIterator(page1, client)
//
//	for pageIterator.HasNext() {
//		pageCount += 1
//		_ = pageIterator.Next()
//	}
//
//	assert.Equal(t, pageCount, 3)
//}
//func TestIterateEnumeratesAllPages(t *testing.T) {
//	res := make([]string, 0)
//	client := http.Client{}
//	pageIterator := NewPageIterator(page1, client)
//
//	pageIterator.Iterate(func(pageItem Item) bool {
//		res = append(res, pageItem.GetDisplayName())
//		return true
//	})
//
//	assert.Equal(t, 12, len(res))
//}
//
//func TestIterateCanBePausedAndResumed(t *testing.T) {
//	res := make([]string, 0)
//	res2 := make([]string, 0)
//	client := http.Client{}
//	pageIterator := NewPageIterator(page1, client)
//
//	pageIterator.Iterate(func(pageItem Item) bool {
//		res = append(res, pageItem.GetDisplayName())
//
//		if pageItem.GetDisplayName() == "User D" {
//			return false
//		}
//
//		return true
//	})
//	assert.Equal(t, res, []string{"User A", "User B", "User C", "User D"})
//
//	pageIterator.Iterate(func(pageItem Item) bool {
//		pageItem = pageItem
//		res2 = append(res2, pageItem.GetDisplayName())
//
//		if pageItem.GetDisplayName() == "User G" {
//			return false
//		}
//
//		return true
//	})
//	assert.Equal(t, res2, []string{"User D", "User E", "User F", "User G"})
//}
