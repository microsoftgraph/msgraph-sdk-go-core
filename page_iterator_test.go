package msgraphgocore

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	testing "testing"
)

type PageItem struct {
	DisplayName string
}

type UserPage struct {
	Value    []Item
	NextLink string
}

var page1 = &UserPage{
	Value: []Item{
		PageItem{DisplayName: "User A"},
		PageItem{DisplayName: "User B"},
		PageItem{DisplayName: "User C"},
		PageItem{DisplayName: "User D"},
	},
	NextLink: "page_2",
}

var page2 = &UserPage{
	Value: []Item{
		PageItem{DisplayName: "User E"},
		PageItem{DisplayName: "User F"},
		PageItem{DisplayName: "User G"},
		PageItem{DisplayName: "User H"},
	},
	NextLink: "page_3",
}

var page3 = &UserPage{
	Value: []Item{
		PageItem{DisplayName: "User I"},
		PageItem{DisplayName: "User J"},
		PageItem{DisplayName: "User K"},
		PageItem{DisplayName: "User L"},
	},
	NextLink: "",
}

func (p *UserPage) GetValue() []Item {
	return p.Value
}

func (p *UserPage) GetNextLink() string {
	return p.NextLink
}

func (p *UserPage) GetNextPage() Page {
	switch nextLink := p.GetNextLink(); nextLink {
	case "page_2":
		return page2
	case "page_3":
		return page3
	}
	return nil
}

func (i PageItem) GetDisplayName() string {
	return i.DisplayName
}

func TestHasNextReturnsTrueIfNextPageIsAvailable(t *testing.T) {
	client := http.Client{}
	pageIterator := NewPageIterator(page1, client)
	hasNext := pageIterator.HasNext()

	assert.True(t, hasNext)
}

func TestNextReturnsNextPage(t *testing.T) {
	client := http.Client{}
	pageIterator := NewPageIterator(page1, client)
	nextPage := pageIterator.Next()

	assert.Equal(t, nextPage, page2)
}

func TestIterateStopsWhenCallbackReturnsFalse(t *testing.T) {
	res := make([]string, 0)

	client := http.Client{}
	pageIterator := NewPageIterator(page1, client)
	pageIterator.Iterate(func(pageItem Item) bool {
		res = append(res, pageItem.GetDisplayName())
		if pageItem.GetDisplayName() == "User D" {
			return false
		}
		return true
	})

	assert.Equal(t, 4, len(res))
}

func TestNextAndHasNext(t *testing.T) {
	pageCount := 1
	client := http.Client{}
	pageIterator := NewPageIterator(page1, client)

	for pageIterator.HasNext() {
		pageCount += 1
		_ = pageIterator.Next()
	}

	assert.Equal(t, pageCount, 3)
}
func TestIterateEnumeratesAllPages(t *testing.T) {
	res := make([]string, 0)
	client := http.Client{}
	pageIterator := NewPageIterator(page1, client)

	pageIterator.Iterate(func(pageItem Item) bool {
		res = append(res, pageItem.GetDisplayName())
		return true
	})

	assert.Equal(t, 12, len(res))
}

func TestIterateCanBePausedAndResumed(t *testing.T) {
	res := make([]string, 0)
	res2 := make([]string, 0)
	client := http.Client{}
	pageIterator := NewPageIterator(page1, client)

	pageIterator.Iterate(func(pageItem Item) bool {
		res = append(res, pageItem.GetDisplayName())

		if pageItem.GetDisplayName() == "User D" {
			return false
		}

		return true
	})
	assert.Equal(t, res, []string{"User A", "User B", "User C", "User D"})

	pageIterator.Iterate(func(pageItem Item) bool {
		res2 = append(res2, pageItem.GetDisplayName())

		if pageItem.GetDisplayName() == "User G" {
			return false
		}

		return true
	})
	assert.Equal(t, res2, []string{"User D", "User E", "User F", "User G"})
}
