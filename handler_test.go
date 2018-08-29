package main

import "testing"

func TestGetNamespace(t *testing.T) {
	s := "ops_fs#qb-qa-10"
	name, ns := getNamespace(s)
	if name != "ops_fs" || ns != "qb-qa-10" {
		t.Error("transform failed")
	}
}
