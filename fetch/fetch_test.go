package fetch

import (
	"encoding/json"
	"fmt"
	"testing"
)

func init() {
	Init("http://upstream-pre.sched.qianbao-inc.com")
}

func TestFetchs(t *testing.T) {
	p, err := Fetchs()
	if err != nil {
		t.Error(err)
	}
	b, err := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))
}

func TestFetch(t *testing.T) {
	p, err := Fetch("pre", "ops_fs")
	if err != nil {
		t.Error(err)
	}
	b, err := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))
}

func TestCheck(t *testing.T) {
	a := AProjectIps{
		Name: "ops_fs",
		Env:  "qa",
		IPs: []AProject{
			{
				IP:   "172.28.40.251",
				Port: "3000",
			},
			{
				IP:   "172.28.40.251",
				Port: "3001",
			},
		},
	}
	/*
		p, err := Fetch("qa", "ops_fs")
		if err != nil {
			t.Error(err)
		}
	*/
	r, err := a.Check()
	if err != nil {
		t.Log("check test err", err)
		return
	}
	fmt.Println(r)
}

var aprojectdemo = `
[
  {
    "fail_timeout": 30, 
    "instance_type": 1, 
    "ip_hash": 0, 
    "max_fails": 3, 
    "namespace": "qb-qa-10", 
    "ops_code": 200, 
    "resthub": 0, 
    "resthub_state": 2, 
    "sync_state": 1, 
    "upstream_createtime": "Mon, 01 Jan 2018 01:01:01 GMT", 
    "upstream_env": "qa", 
    "upstream_filename": "ops_fs_upstream.conf", 
    "upstream_group": "BJ-SH", 
    "upstream_id": 5125, 
    "upstream_ip": "172.28.137.22", 
    "upstream_name": "ops_fs", 
    "upstream_pool": "ops_fs_pool", 
    "upstream_port": "8000", 
    "upstream_state": 1, 
    "upstream_updatetime": "Tue, 21 Aug 2018 20:24:48 GMT", 
    "weight": 1000
  }
]`

var testitems = `
[
  {
    "enabled": "1", 
    "env": "test", 
    "idc": "BJ-SH", 
    "ip": "fs.qianbao-inc.com", 
    "is_docker": "1", 
    "port": "80", 
    "project": "ops_test"
  }, 
  {
    "enabled": "1", 
    "env": "test", 
    "idc": "BJ-SH", 
    "ip": "104.16.25.88", 
    "is_docker": "1", 
    "port": "80", 
    "project": "ip.cn"
  }, 
  {
    "enabled": "1", 
    "env": "test", 
    "idc": "BJ-SH", 
    "ip": "example.com", 
    "is_docker": "1", 
    "port": "80", 
    "project": "example"
  }
]
`
