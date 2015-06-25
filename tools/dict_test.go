package tools

import "testing"

func Test_ContainsKeyInMap(t *testing.T) {
	var m1 map[interface{}]interface{} = make(map[interface{}]interface{})

	m1["k1"] = "v1"
	m1[1] = 22

	if !ContainsKeyInMap(m1, "k1") {
		t.Error("m1 有 k1")
	}

	if ContainsKeyInMap(m1, "k2") {
		t.Error("m1 没有 k2")
	}

	if !ContainsKeyInMap(m1, 1) {
		t.Error("m1 有 1")
	}
}