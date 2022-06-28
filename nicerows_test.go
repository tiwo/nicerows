package nicerows

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func panicif(err error) {
	if err != nil {
		panic(err)
	}
}

func exampledb() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	panicif(err)

	err = db.Ping()
	panicif(err)

	_, err = db.Exec(`
		Create table "t1"(a integer primary key, b, c, d, e);
	`)
	panicif(err)

	_, err = db.Exec(`
		Insert into "t1"
		values (1, "foo", 0, -2.5, NULL),
		       (2, "bar", 10, -7.5, x'414243')
	`)
	panicif(err)

	return db
}

func TestAnypointers(t *testing.T) {
	length := 100
	vals, ptrs := anypointers(100)
	for i := 0; i < length; i++ {
		if ptrs[i] != &vals[i] {
			t.Fatalf("anypointers at position %v", i)
		}
	}
}

func TestSliceResults(t *testing.T) {
	db := exampledb()
	rows, err := db.Query(` Select * from "t1"; `)
	nr := New(rows, err)

	ok := nr.Next()
	if !ok {
		t.Fatalf("Query should return 2 results, returns 0.")
	}
	actual_row := fmt.Sprintf("%#v", nr.Slice())
	expected_row := `[]interface {}{1, "foo", 0, -2.5, interface {}(nil)}`
	if actual_row != expected_row {
		t.Fatalf("result should be %v, but is %#v\n", expected_row, nr.currentrow)
	}

	ok = nr.Next()
	if !ok {
		t.Fatalf("Query should return 2 results, returns 1.")
	}
	actual_row = fmt.Sprintf("%#v", nr.Slice())
	expected_row = `[]interface {}{2, "bar", 10, -7.5, []uint8{0x41, 0x42, 0x43}}`
	if actual_row != expected_row {
		t.Fatalf("result should be %v, but is %#v\n", expected_row, nr.currentrow)
	}

	ok = nr.Next()
	if ok {
		t.Fatalf("Query should return 2 results, returns more.")
	}
}

func TestJsonResult(t *testing.T) {
	db := exampledb()
	rows, err := db.Query(` Select * from "t1"; `)
	nr := New(rows, err)

	ok := nr.Next()
	if !ok {
		t.Fatalf("Query should return 2 results, returns 0.")
	}

	actual_json := nr.Json()
	expected_json := `[1,"foo",0,-2.5,null]`

	if actual_json != expected_json {
		t.Fatalf("result should be %v, but is %v\n", expected_json, actual_json)
	}

	ok = nr.Next()
	if !ok {
		t.Fatalf("Query should return 2 results, returns only 1.")
	}

	actual_json = nr.Json()
	expected_json = `[2,"bar",10,-7.5,"ABC"]`

	if actual_json != expected_json {
		t.Fatalf("result should be %v, but is %v\n", expected_json, actual_json)
	}

	ok = nr.Next()
	if ok {
		t.Fatalf("Query should return 2 results, returns more.")
	}

}

func ExampleNiceRows_Json() {
	db := exampledb()
	rows, err := db.Query(` Select * from "t1"; `)
	nr := New(rows, err)
	if nr.Err != nil {
		//return nr.Err
	}
	for nr.Next() {
		fmt.Println(nr.Json())
	}
	//return nr.Err

	// Output:
	// [1,"foo",0,-2.5,null]
	// [2,"bar",10,-7.5,"ABC"]
}

func ExampleNiceRows_Colnames() {
	db := exampledb()
	rows, err := db.Query(` Select * from "t1"; `)
	nr := New(rows, err)
	if nr.Err != nil {
		//return nr.Err
	}
	fmt.Println(nr.Colnames())

	//return nr.Err

	// Output:
	// [a b c d e]
}
