package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFetch(t *testing.T) {
	p, err := fetch(*upstreamAPI)
	if err != nil {
		t.Error(err)
	}
	b, err := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))
}
