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
