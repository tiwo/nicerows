package nicerows

// Return a slice of values, and a slice of pointers to each of them,
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
