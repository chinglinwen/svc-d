// periodically fetch checks info
package main

// call upstream api, transformation to  checkers
// how to know if it's http or tcp?

// does it support setting for every check? info come from platform?

// sync all or, sync one?

// check info constantly change? checker need update?
// we usually check the old info, if things changed, resync?

// 5 minutes interval

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty"
	"github.com/sourcegraph/checkup"
)

const upstreamapi = "http://upstream-pre.sched.qianbao-inc.com/get_upstream_all_instance/"

func start() {
	c := time.Tick(1 * time.Minute)
	var stop *time.Ticker
	for now := range c {
		fmt.Printf("starting fetch at %v\n", now)
		projects, err := fetch(upstreamapi)
		if err != nil {
			fmt.Println("fetch error", err)
			continue
		}
		stop.Stop()
		convertToCheck(projects)
		stop = conf.CheckAndStoreEvery(1 * time.Minute)
	}
}

func convertToCheck(projects Projects) {
	for _, p := range projects {
		if p.Enabled == "0" {
			continue
		}

		check := checkup.HTTPChecker{
			Name:     p.Project,
			URL:      "http://" + p.IP + ":" + p.Port,
			Attempts: 5,
		}
		conf.Checkers = append(conf.Checkers, check)
	}
}

type Projects []Prjoect

type Prjoect struct {
	Project  string
	Enabled  string
	ENV      string
	IDC      string
	IP       string
	IsDocker string
	Port     string
}

func fetch(url string) (p Projects, err error) {
	resp, e := resty.R().Get(url)
	if e != nil {
		err = e
		return
	}
	err = json.Unmarshal(resp.Body(), &p)
	return
}
