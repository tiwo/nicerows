package nicerows

import (
	"database/sql"
	"encoding/json"
)

type NiceRows struct {
	sqlresult *sql.Rows
	colnames  []string
	err       error
}

func New(sqlresult *sql.Rows, err error) *NiceRows {

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

func anypointers(length int) ([]any, []any) {
	values := make([]any, length)
	pointers := make([]any, length)
	for i := 0; i < length; i++ {
		pointers[i] = &values[i]
	}
	return values, pointers
}

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

func bytearray2string(thing any) any {
	blob, ok := thing.([]byte)
	if ok {
		return string(blob)
	}
	return thing
}

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
