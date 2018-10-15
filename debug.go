package main

import (
	"wen/svc-d/check"
	"wen/svc-d/notice"
)

func setdebug() {
	check.DEBUG = true
	notice.DEBUG = true
}

func stopdebug() {
	check.DEBUG = false
	notice.DEBUG = false
}
