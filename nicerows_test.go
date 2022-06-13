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

func TestIterateSlices(t *testing.T) {
	db := exampledb()
	rows, err := db.Query(` Select * from "t1"; `)
	nr := New(rows, err)
	it := nr.IterateSlices(false)

	actual := fmt.Sprintf("%#v", <-it)
	if actual != `[]interface {}{1, "foo", 0, -2.5, interface {}(nil)}` {
		t.Fatal("Iteration 1!")
	}

	actual = fmt.Sprintf("%#v", <-it)
	if actual != `[]interface {}{2, "bar", 10, -7.5, []uint8{0x41, 0x42, 0x43}}` {
		t.Fatal("Iteration 2!")
	}

	_, ok := <-it
	if ok {
		t.Fatal("Iteration continued after last row!")
	}

}

func TestIterateinfinity(t *testing.T) {
	db := exampledb()
	rows, err := db.Query(`
		With recursive "numbers"("n")
		as (
			Select 0
			union all
			select "n"+1 from "numbers"
		)
		select "n" from "numbers"
	`)
	it := New(rows, err).IterateSlices(false)
	for i := 0; i < 1000; i++ {
		row, ok := <-it

		if !ok {
			t.Fatalf("Infinite iterator ends too soon (at %v)", i)
		}

		target := fmt.Sprintf("[%v]", i)
		actual := fmt.Sprintf("%v", row)
		if actual != target {
			t.Fatalf("At index %v: %v != %v", i, actual, target)
		}
	}

}

func TestIncludeHeader(t *testing.T) {
	db := exampledb()
	rows, err := db.Query(`
		Select * from "t1";
	`)
	nr := New(rows, err)
	it := nr.IterateSlices(true)
	hdr := <-it

	actual := fmt.Sprintf("%v", hdr)
	if actual != "[a b c d e]" {
		t.Fatalf("Header seems not right: %v", actual)
	}

}

func TestSqlSyntaxError(t *testing.T) {
	db := exampledb()
	rows, err := db.Query(`
		Window exists for outer view. Ignore limit, then raise cross into outer vacuum.
	`)
	nr := New(rows, err)

	if nr.err == nil {
		t.Fatalf("Nonsensical SQL should produce an error, but returns %#v", nr.err)
	}

	it := nr.IterateSlices(false)
	for row := range it {
		t.Fatalf("Nonsensical SQL should not yield any rows, but does: %#v", row)
	}

	if nr.err == nil {
		t.Fatalf("Nonsensical SQL should produce an error, but returns %#v", nr.err)
	}

}

func TestJsonlines(t *testing.T) {
	db := exampledb()
	rows, err := db.Query(` Select * from "t1"; `)
	nr := New(rows, err)
	it := nr.IterateJsonlines(true)

	actual := <-it
	if actual != `["a","b","c","d","e"]` {
		t.Fatalf("header: %v", actual)
	}

	actual = <-it
	if actual != `[1,"foo",0,-2.5,null]` {
		t.Fatalf("first row: %v", actual)
	}

	actual = <-it
	if actual != `[2,"bar",10,-7.5,"ABC"]` {
		t.Fatalf("second row: %v", actual)
	}

	_, ok := <-it
	if ok {
		t.Fatalf("Iterator continued after last row")
	}
}

/*
values (1, "foo", 0, -2.5, NULL),
(2, "bar", 10, -7.5, x'414243')
*/
