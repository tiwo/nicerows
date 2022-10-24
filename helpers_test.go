package nicerows

import "testing"

func TestAnypointers(t *testing.T) {
	length := 100
	vals, ptrs := anypointers(100)
	for i := 0; i < length; i++ {
		if ptrs[i] != &vals[i] {
			t.Fatalf("anypointers at position %v", i)
		}
	}
}

func TestBytearray2string(t *testing.T) {

	var testmatrix = []struct{ input, expect any }{
		{"strings are not changed 固定点", "strings are not changed 固定点"},
		{[]byte("byte arrays are converted. Канвертавац байт."), "byte arrays are converted. Канвертавац байт."},
		{nil, nil},
		{-12389, -12389},
		{-42.389, -42.389},
	}

	for _, testcase := range testmatrix {
		actual := bytearray2string(testcase.input)
		if actual != testcase.expect {
			t.Errorf("bytearray2string(%#v) is %#v, not %#v as expected.", testcase.input, actual, testcase.expect)
		}
	}

}
