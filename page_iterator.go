package msgraphgocore

import "net/http"

type Item interface {
	GetDisplayName() string
}
type Page interface {
	GetValue() []Item
	GetNextLink() string
	GetNextPage() Page
}

type PageIterator struct {
	page          Page
	client        http.Client
	keepIterating bool
	pauseIndex    int
}

func NewPageIterator(page Page, client http.Client) *PageIterator {
	return &PageIterator{
		page,
		client,
		true,
		0,
	}
}

func (pI *PageIterator) HasNext() bool {
	if pI.page.GetNextLink() == "" {
		return false
	}
	return true
}

func (pI *PageIterator) Next() Page {
	nextPage := pI.page.GetNextPage()
	pI.page = nextPage

	return nextPage
}

func (pI *PageIterator) Iterate(callback func(pageItem Item) bool) {
	for pI.page != nil {
		keepIterating := pI.enumerate(callback)

		if !keepIterating {
			return
		}

		pI.Next()
		pI.pauseIndex = 0
	}
}

func (pI *PageIterator) enumerate(callback func(item Item) bool) bool {
	keepIterating := true
	pageItems := pI.page.GetValue()

	for i := pI.pauseIndex; i < len(pageItems); i++ {
		keepIterating = callback(pageItems[i])

		if !keepIterating {
			pI.pauseIndex = i
			break
		}
	}

	return keepIterating
}
