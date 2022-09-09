package msgraphgocore

import "github.com/microsoft/kiota-abstractions-go/serialization"

// SetValue takes a generic source , reads value and writs it to a setter
//
// `source func() (*T, error)` is a generic getter with possible error response.
// `setter func(t *T)` generic function that can write a value from the source
func SetValue[T interface{}](source func() (*T, error), setter func(t *T)) error {
	val, err := source()
	if err != nil {
		return err
	}
	if val != nil {
		setter(val)
	}
	return nil
}

// SetObjectValue takes a generic source with a discriminator receiver, reads value and writes it to a setter
//
// `source func() (*T, error)` is a generic getter with possible error response.
// `setter func(t *T)` generic function that can write a value from the source
func SetObjectValue[T interface{}](source func(ctor serialization.ParsableFactory) (serialization.Parsable, error), ctor serialization.ParsableFactory, setter func(t T)) error {
	val, err := source(ctor)
	if err != nil {
		return err
	}
	if val != nil {
		res := (val).(T)
		setter(res)
	}
	return nil
}
