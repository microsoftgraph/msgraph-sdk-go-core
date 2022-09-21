package msgraphgocore

import "github.com/microsoft/kiota-abstractions-go/serialization"

// SetValue receives a source function and applies the results to the setter
//
// source is any function that produces (*T, error)
// setter recipient function of the result of the source if no error is produces
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

// SetCollectionValue is a utility function that receives a collection that can be cast to Parsable and a function that epects the results
//
// source is any function that receives a `ParsableFactory` and returns a slice of Parsable or error
// ctor is a ParsableFactory
// setter is a recipient of the function results
func SetCollectionValue[T interface{}](source func(ctor serialization.ParsableFactory) ([]serialization.Parsable, error), ctor serialization.ParsableFactory, setter func(t []T)) error {
	val, err := source(ctor)
	if err != nil {
		return err
	}
	if val != nil {
		res := make([]T, len(val))
		for i, v := range val {
			res[i] = v.(T)
		}
		setter(res)
	}
	return nil
}

// CollectionApply applies an operation to every element of the slice and returns a result of the modified collection
//
//  is a slice of all the elementents to be mutated
// mutator applies an operation to the collection and returns a response of type `R`
func CollectionApply[T any, R interface{}](collection []T, mutator func(t T) R) []R {
	cast := make([]R, len(collection))
	for i, v := range collection {
		cast[i] = mutator(v)
	}
	return cast
}

// SetEnumValue is a utility function that receives an enum getter , EnumFactory and applies the results to a setter
//
// source is any function that receives a `EnumFactory` and returns an interface or error
// parser is an EnumFactory
// setter is a recipient of the function results
func SetEnumValue[T interface{}](source func(parser serialization.EnumFactory) (interface{}, error), parser serialization.EnumFactory, setter func(t T)) error {
	val, err := source(parser)
	if err != nil {
		return err
	}
	if val != nil {
		res := (val).(T)
		setter(res)
	}
	return nil
}

// SetReferencedEnumValue is a utility function that receives an enum getter , EnumFactory and applies a de-referenced result of the factory to a setter
//
// source is any function that receives a `EnumFactory` and returns an interface or error
// parser is an EnumFactory
// setter is a recipient of the function results
func SetReferencedEnumValue[T interface{}](source func(parser serialization.EnumFactory) (interface{}, error), parser serialization.EnumFactory, setter func(t *T)) error {
	val, err := source(parser)
	if err != nil {
		return err
	}
	if val != nil {
		res := (val).(*T)
		setter(res)
	}
	return nil
}

// SetCollectionOfReferencedEnumValue is a utility function that receives an enum collection source , EnumFactory and applies a de-referenced result of the factory to a setter
//
// source is any function that receives a `EnumFactory` and returns an interface or error
// parser is an EnumFactory
// setter is a recipient of the function results
func SetCollectionOfReferencedEnumValue[T interface{}](source func(parser serialization.EnumFactory) ([]interface{}, error), parser serialization.EnumFactory, setter func(t []*T)) error {
	val, err := source(parser)
	if err != nil {
		return err
	}
	if val != nil {
		res := make([]*T, len(val))
		for i, v := range val {
			res[i] = (v).(*T)
		}
		setter(res)
	}
	return nil
}

// SetCollectionOfPrimitiveValue is a utility function that receives a collection of primitives , targetType and applies the result of the factory to a setter
//
// source is any function that receives a `EnumFactory` and returns an interface or error
// targetType is a string representing the type of result
// setter is a recipient of the function results
func SetCollectionOfPrimitiveValue[T interface{}](source func(targetType string) ([]interface{}, error), targetType string, setter func(t []T)) error {
	val, err := source(targetType)
	if err != nil {
		return err
	}
	if val != nil {
		res := make([]T, len(val))
		for i, v := range val {
			res[i] = v.(T)
		}
		setter(res)
	}
	return nil
}

// SetCollectionOfReferencedPrimitiveValue is a utility function that receives a collection of primitives , targetType and applies the re-referenced result of the factory to a setter
//
// source is any function that receives a `EnumFactory` and returns an interface or error
// parser is an EnumFactory
// setter is a recipient of the function results
func SetCollectionOfReferencedPrimitiveValue[T interface{}](source func(targetType string) ([]interface{}, error), targetType string, setter func(t []T)) error {
	val, err := source(targetType)
	if err != nil {
		return err
	}
	if val != nil {
		res := make([]T, len(val))
		for i, v := range val {
			res[i] = *(v.(*T))
		}
		setter(res)
	}
	return nil
}

func Point[T interface{}](t T) *T {
	return &t
}
