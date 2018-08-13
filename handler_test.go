package main

import (
	"fmt"
	"testing"
)

func TestParseState(t *testing.T) {
	test := `[
  true, 
  "bb 104.16.25.88:80 \u6ca1\u6709\u6ce8\u518c\u5230\u8c03\u5ea6\u4e2d\u5fc3,"
]`
	fmt.Println(parseState([]byte(test)))

}
func TestChangeState(t *testing.T) {
	tests := []struct {
		endpoint, title string
		state           bool
	}{
		{"http://104.16.25.88:80", "aa", true},
		{"tcp://104.16.25.88:80", "bb", true},
		//{"104.16.25.88:80", "cc", "80"},
	}
	for _, v := range tests {
		ok, err := ChangeState(v.endpoint, v.title)
		if err != nil || ok != v.state {
			t.Error("err", v, "got", ok, "want", v.state)
			continue
		}
	}
}

func TestEndpoint2ip(t *testing.T) {
	tests := []struct {
		endpoint, ip, port string
	}{
		{"http://104.16.25.88:80", "104.16.25.88", "80"},
		{"tcp://104.16.25.88:80", "104.16.25.88", "80"},
		{"104.16.25.88:80", "104.16.25.88", "80"},
	}
	for _, v := range tests {
		i, p := endpoint2ip(v.endpoint)
		if i != v.ip || p != v.port {
			t.Error(v.endpoint, "got", i, p, "want", v.ip, v.port)
		}
	}
}
