package client

import (
	"errors"
	"reflect"
	"time"
)

// SliceContains returns true if a slice contains an element otherwise false
func SliceContains(slice, elem interface{}) (bool, error) {

	sv := reflect.ValueOf(slice)

	// Check that slice is actually a slice/array.
	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Array {
		return false, errors.New("not an array or slice")
	}

	// iterate the slice
	for i := 0; i < sv.Len(); i++ {

		// compare elem to the current slice element
		if elem == sv.Index(i).Interface() {
			return true, nil
		}
	}

	// nothing found
	return false, nil

}

func Epoch() uint64 {
	return uint64(time.Now().Unix())
}
