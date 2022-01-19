package msgraphgocore

import (
	abstractions "github.com/microsoft/kiota/abstractions/go"
	serialization "github.com/microsoft/kiota/abstractions/go/serialization"
	jsonserialization "github.com/microsoft/kiota/serialization/go/json"
	"log"
	url "net/url"
)

type Item interface{}

type Page interface {
	GetValue() []Item
	GetNextLink() *string
}

type PageIterator struct {
	page            Page
	reqAdapter      GraphRequestAdapterBase
	pauseIndex      int
	constructorFunc ParsableConstructor
}

type ParsableConstructor func() serialization.Parsable

func NewPageIterator(page Page, reqAdapter GraphRequestAdapterBase, constructorFunc ParsableConstructor) *PageIterator {
	abstractions.RegisterDefaultSerializer(func() serialization.SerializationWriterFactory {
		return jsonserialization.NewJsonSerializationWriterFactory()
	})
	abstractions.RegisterDefaultDeserializer(func() serialization.ParseNodeFactory {
		return jsonserialization.NewJsonParseNodeFactory()
	})

	return &PageIterator{
		page,
		reqAdapter,
		0, // pauseIndex helps us remember where we paused enumeration in the page.
		constructorFunc,
	}
}

func (pI *PageIterator) HasNext() bool {
	if pI.page.GetNextLink() == nil {
		return false
	}
	return true
}

func (pI *PageIterator) Next() Page {
	nextPage := pI.getNextPage().(Page)

	pI.page = nextPage
	return nextPage
}

func (pI *PageIterator) getNextPage() interface{} {
	nextLink, err := url.Parse(*pI.page.GetNextLink())
	if err != nil {
		log.Fatal(err)
	}

	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.GET
	requestInfo.SetUri(*nextLink)

	res, err := pI.reqAdapter.SendAsync(*requestInfo, pI.constructorFunc, nil)
	if err != nil {
		log.Fatal(err)
	}

	return res
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
