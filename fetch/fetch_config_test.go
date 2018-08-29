package fetch

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFetchConfig(t *testing.T) {
	p, err := FetchConfig("ops_test")
	if err != nil {
		t.Error(err)
	}
	b, err := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))
}

func TestGetProjectName(t *testing.T) {
	p, err := GetProjectName("ops-fs")
	if err != nil {
		t.Error("GetProjectName err", err)
	}
	want := "ops_fs"
	if p != want {
		t.Errorf("get project name test fail, got %v, want %v\n", p, want)
	}
}

func TestFetchConfigs(t *testing.T) {
	p, err := FetchConfigs()
	if err != nil {
		t.Error(err)
	}
	b, err := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))
}

// for a simple verify api return data is ok
func TestDecodeConfig(t *testing.T) {
	p, err := decodeConfig([]byte(ademoBody))
	if err != nil {
		t.Error(err)
	}
	b, err := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))
}

func TestDecodeConfigs(t *testing.T) {
	p, err := decodeConfigs([]byte(demoBody))
	if err != nil {
		t.Error(err)
	}
	b, err := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))
}

var ademoBody = `
{
    "code": 0,
    "data": {
        "attempt_spacing": 23,
        "attempts": 7,
        "enabled": "on",
        "must_contain": "fadsfa",
        "project_name": "ops_test",
        "status_code": 504,
        "timeout": 600,
        "type": "http",
        "uri": "/"
    },
    "message": "success"
}
`

var demoBody = `
{
    "code": 1,
    "data": [
        {
            "attempt_spacing": 0,
            "attempts": 6,
            "enabled": "on",
            "must_contain": "",
            "project_name": "ops_test",
            "status_code": 200,
            "timeout": 10,
            "type": "tcp",
            "uri": "/"
        },
        {
            "attempt_spacing": 0,
            "attempts": 5,
            "enabled": "on",
            "must_contain": "afafaa",
            "project_name": "ops_test2",
            "status_code": 500,
            "timeout": 10,
            "type": "http",
            "uri": "/"
        },
        {
            "attempt_spacing": 0,
            "attempts": 5,
            "enabled": "on",
            "must_contain": "erqer",
            "project_name": "ops_test_one",
            "status_code": 200,
            "timeout": 10,
            "type": "http",
            "uri": "/"
        },
        {
            "attempt_spacing": 1,
            "attempts": 50,
            "enabled": "off",
            "must_contain": "dafadf",
            "project_name": "test",
            "status_code": 500,
            "timeout": 100,
            "type": "http",
            "uri": "/asdfasdf"
        }
    ],
    "message": "success"
}`

var demoSingleBody = `
我的电脑:
成功的情况：
{
    "code": 0,
    "data": {
        "attempt_spacing": 23,
        "attempts": 7,
        "enabled": "on",
        "must_contain": "fadsfa",
        "project_name": "ops_test",
        "status_code": 504,
        "timeout": 600,
        "type": "http",
        "uri": "/"
    },
    "message": "success"
}

我的电脑:
失败的情况：
{
    "code": 1,
    "message": "获取ops_test_two 的检测信息失败"
}

`
