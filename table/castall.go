package table

import (
	"fmt"
	"reflect"
)

// UnmarshalAll loads and unmarshals all rows returned by the iterator.
//
// The result argument must necessarily be the address for a slice. The slice
// may be nil or previously allocated.
func UnmarshalAll(iter Iterator, out interface{}) error {
	outv := reflect.ValueOf(out)
	if outv.Kind() != reflect.Ptr || outv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("out argument must be a slice address")
	}
	slicev := outv.Elem()
	slicev = slicev.Slice(0, 0) // Trucantes the passed-in slice.
	elemt := slicev.Type().Elem()
	i := 0
	for iter.Next() {
		elemp := reflect.New(elemt)
		if err := iter.UnmarshalRow(elemp.Interface()); err != nil {
			return err
		}
		slicev = reflect.Append(slicev, elemp.Elem())
		slicev = slicev.Slice(0, slicev.Len())
		i++
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	outv.Elem().Set(slicev.Slice(0, i))
	return nil
}
