package nicerows

import (
	"encoding/json"
	"fmt"
)

// The core type of nicerows
type NiceRows struct {
	sqlresult  SqlResult // will be assigned in New()
	colnames   []string  // will be assigned in New()
	currentrow []any     // will be nil before a call to Next()
	Err        error
}

// An interface implemented by sql.Rows (from database/sql)
type SqlResult interface {
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...any) error
}

// Create a new NiceRows struct.
// It keeps only minimal state, but should work in the sense that eg. losing the
// network connection to the database will stop the iteration. Errors are
// signalled in Nicerows.Err!
func New(sqlresult SqlResult, err error) *NiceRows {

	nr := &NiceRows{
		sqlresult:  sqlresult,
		colnames:   nil,
		currentrow: nil,
		Err:        err,
	}

	if err != nil {
		return nr
	}

	nr.colnames, nr.Err = sqlresult.Columns()
	return nr
}

func (nr *NiceRows) Colnames() []string {
	if nr.colnames == nil {
		panic(fmt.Sprintf("Colnames() on unititalized *NiceRows %#v", nr))
	}
	return copySlice(nr.colnames)
}

// retrieve the next row. Suitable for a `for nr.Next() {...}` loop
func (nr *NiceRows) Next() bool {
	if nr.Err != nil {
		nr.currentrow = nil
		return false
	}

	ok := nr.sqlresult.Next()
	if !ok {
		nr.currentrow = nil
		return false
	}

	values, pointers := anypointers(len(nr.colnames))

	nr.Err = nr.sqlresult.Scan(pointers...)
	if nr.Err != nil {
		nr.currentrow = nil
		return false
	}
	nr.currentrow = values

	return true
}

func (nr *NiceRows) Slice() []any {
	if nr.currentrow == nil {
		panic("(*NiceRows).Slice() before Next()")
	}
	return copySlice(nr.currentrow)
}

func (nr *NiceRows) Map() map[string]any {
	if nr.currentrow == nil {
		panic("(*NiceRows).Slice() before Next()")
	}
	return constructMap(nr.colnames, nr.currentrow)
}

func (nr *NiceRows) Json() string {
	if nr.currentrow == nil {
		panic("(*NiceRows).Json() before Next()")
	}
	if nr.Err != nil {
		return ""
	}

	slice := nr.Slice()
	for i := range slice {
		slice[i] = bytearray2string(slice[i])
	}

	js, err := json.Marshal(slice)
	if err != nil {
		panic(fmt.Sprintf("(*NiceRows).Json(): Error with JSON encode: %#v", err))
	}

	return string(js)
}
