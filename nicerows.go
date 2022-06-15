package nicerows

import (
	"encoding/json"
)

// The core type of nicerows
type NiceRows struct {
	sqlresult SqlResult
	colnames  []string
	err       error
}

// An interface compatible with sql.Rows (from database/sql)
type SqlResult interface {
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...any) error
}

// Create a new NiceRows struct.
// It keeps minimal state, but should work in the sense that eg. losing the
// network connection to the database will stop the iteration and signal
// the error in the err field
func New(sqlresult SqlResult, err error) *NiceRows {

	nr := &NiceRows{
		sqlresult: sqlresult,
		colnames:  nil,
		err:       err,
	}

	if err != nil {
		return nr
	}

	nr.colnames, nr.err = sqlresult.Columns()
	return nr
}

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

// Iterate over all rows, as `[]any` slices
func (nr *NiceRows) IterateSlices() chan []any {

	out := make(chan []any)

	if nr.err != nil {
		close(out)
		return out
	}

	go func() {
		defer close(out)

		length := len(nr.colnames)

		// send the column names as first slice:
		header := make([]any, length)
		for i, name := range nr.colnames {
			header[i] = name
		}
		out <- header

		for nr.sqlresult.Next() {
			values, pointers := anypointers(length)
			nr.err = nr.sqlresult.Scan(pointers...)
			if nr.err != nil {
				return
			}
			out <- values
		}
	}()

	return out
}

// Iterate over all rows, as `map[string]any` maps.
// Note that the names come from the SQL driver and ultimately from the
// SQL query or database scheme; if those column names are not unique,
// later columns will eclipse former ones with the same name.
func (nr *NiceRows) IterateMaps() chan map[string]any {

	out := make(chan map[string]any)

	if nr.err != nil {
		close(out)
		return out
	}

	go func() {
		defer close(out)
		length := len(nr.colnames)

		for nr.sqlresult.Next() {
			values, pointers := anypointers(length)
			nr.err = nr.sqlresult.Scan(pointers...)
			if nr.err != nil {
				return
			}

			m := make(map[string]any)

			for i, name := range nr.colnames {
				m[name] = values[i]
			}

			out <- m
		}
	}()

	return out

}

// Convert arguments of []byte type to string; return anything else unchanged.
func bytearray2string(thing any) any {
	blob, ok := thing.([]byte)
	if ok {
		return string(blob)
	}
	return thing
}

// Iterate over all rows, as JSON arrays.
func (nr *NiceRows) IterateJson() chan string {

	out := make(chan string)

	it := nr.IterateSlices()

	go func() {
		defer close(out)

		for slice := range it {

			for i, val := range slice {
				// SQL strings come as []byte, I want them to convert to JSON strings.
				// Maybe it would be better to create custom json.Encoder or json.Marshaler?
				slice[i] = bytearray2string(val)
			}

			js, err := json.Marshal(slice)

			if err != nil {
				nr.err = err
				return
			}
			out <- string(js)
		}
	}()

	return out
}

func (nr *NiceRows) IterateJsonObjects() chan string {
	out := make(chan string)
	it := nr.IterateMaps()

	go func() {
		defer close(out)
		for in := range it {

			m := make(map[string]any)
			for key, val := range in {
				m[key] = bytearray2string(val)
			}

			js, err := json.Marshal(m)

			if err != nil {
				nr.err = err
				return
			}

			out <- string(js)

		}

	}()

	return out
}
