package nicerows

import (
	"fmt"
)

// return a slice of values, and a slice of pointers to each of them,
// suitable to retrieve a row from database/sql via `_=sql.Rows.Scan(pointers...)`
func anypointers(length int) ([]any, []any) {
	values := make([]any, length)
	pointers := make([]any, length)
	for i := 0; i < length; i++ {
		pointers[i] = &values[i]
	}
	return values, pointers
}

// Convert arguments of []byte type to string; return anything else unchanged.
func bytearray2string(thing any) any {
	blob, ok := thing.([]byte)
	if ok {
		return string(blob)
	}
	return thing
}

func copySlice[T string | any](slice []T) []T {
	same := make([]T, len(slice))
	copy(same, slice)
	return same
}

func copySliceWithStrings(slice []any) []any {
	kopy := make([]any, len(slice))
	for i, val := range slice {
		kopy[i] = bytearray2string(val)
	}
	return kopy
}

func constructMap(names []string, slice []any) map[string]any {
	if len(names) != len(slice) {
		panic(fmt.Sprintf("constructMap: len(names)==%v != len(slice)==%v", len(names), len(slice)))
	}

	m := make(map[string]any)

	for i, name := range names {
		m[name] = bytearray2string(slice[i])
	}

	return m
}

func constructMapWithStrings(names []string, slice []any) map[string]any {
	if len(names) != len(slice) {
		panic(fmt.Sprintf("constructMap: len(names)==%v != len(slice)==%v", len(names), len(slice)))
	}

	m := make(map[string]any)

	for i, name := range names {
		m[name] = bytearray2string(slice[i])
	}

	return m
}
