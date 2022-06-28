package nicerows

import (
	"testing"

	"github.com/tiwo/tigo/assert"
)

func Test_bytearray2string(t *testing.T) {

	assert.DeepEqual(t, "string unchanged", bytearray2string("string unchanged"))
	assert.DeepEqual(t, 1234, bytearray2string(1234))
	assert.DeepEqual(t, interface{}(nil), bytearray2string(interface{}(nil)))
	assert.DeepEqual(t, []uint16{4, 5, 6}, bytearray2string([]uint16{4, 5, 6}))
	assert.DeepEqual(t, bytearray2string([]byte{65, 66, 67}), "ABC")
	assert.DeepEqual(t, bytearray2string([]uint8{65, 66, 67}), "ABC")
}

func TestConstructMap(t *testing.T) {
	names := []string{"foo", "bar", "baz"}
	values := []any{"abc", -1.5, nil}
	m := constructMap(names, values)
	expect := `map[string]interface {}{"bar":-1.5, "baz":interface {}(nil), "foo":"abc"}`
	assert.Equalf(t, expect, m)

}
