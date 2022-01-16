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
	page       Page
	client     http.Client
	pauseIndex int
}

func NewPageIterator(page Page, client http.Client) *PageIterator {
	//TODO: Pass client to page

	return &PageIterator{
		page,
		client,
		0, // pauseIndex helps us remember where we paused enumeration in the page.
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
			// Callback returned false, stop iterating through pages.
			return
		}

		pI.Next()
		pI.pauseIndex = 0 // when moving to the next page reset pauseIndex
	}
}

func (pI *PageIterator) enumerate(callback func(item Item) bool) bool {
	keepIterating := true
	pageItems := pI.page.GetValue()

	// start/continue enumerating page items from  pauseIndex.
	// this makes it possible to resume iteration from where we paused iteration.
	for i := pI.pauseIndex; i < len(pageItems); i++ {
		keepIterating = callback(pageItems[i])

		if !keepIterating {
			// Callback returned false, pause! stop enumerating page items. Set pauseIndex so that we know
			// where to resume from.
			pI.pauseIndex = i
			break
		}
	}

	return keepIterating
}
