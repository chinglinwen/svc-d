package main

import "testing"

func TestGetNamespace(t *testing.T) {
	s := "ops_fs#qb-qa-10#tcp"
	name, ns := getNamespace(s)
	if name != "ops_fs" || ns != "qb-qa-10" {
		t.Error("transform failed")
	}
}

//curl "localhost:8089/check?appname=ops_fs&env=qa"
//curl "172.28.46.201:8089/check?appname=ops_fs&env=qa"
