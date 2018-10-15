// periodically fetch checks info
package check

// call upstream api, transformation to  checkers
// how to know if it's http or tcp?

// does it support setting for every check? info come from platform?

// sync all or, sync one?

// check info constantly change? checker need update?
// we usually check the old info, if things changed, resync?

// 5 minutes interval

import (
	"fmt"
	"time"

	"github.com/chinglinwen/checkup"
	"github.com/chinglinwen/log"
	"github.com/robfig/cron"

	"wen/svc-d/config"
)

var (
	CheckOneTime bool
	TestProject  string

	CheckInterval string
)

// start backend check
func Start(conf *config.Config) {
	log.Printf("start background check for every %v\n", CheckInterval)

	if CheckOneTime {
		check(conf)
		log.Println("check one time only, exit the loop")
		return
	}

	c := cron.New()
	err := c.AddFunc("@every "+CheckInterval, func() { check(conf) })
	if err != nil {
		log.Fatalf("add cron error, checkinterval may not the correct format: %v\n", CheckInterval)
	}
	c.Start()
}

func check(conf *config.Config) {
	// empty it for new fetch
	conf.Checkers = []checkup.Checker{}

	projects, err := Fetchs()
	if err != nil {
		log.Println("fetch upstream error", err)
		return
	}
	configs, err := FetchConfigs()
	if err != nil {
		log.Println("fetch config error", err)
		//continue
	}

	conf.Checkers = ConvertToCheck(projects, configs)
	conf.Save()

	log.Println("fetch ok")

	err = conf.CheckAndStore()
	if err != nil {
		log.Println("background check error", err)
		return
	}
	log.Printf("background check ok\n")
}

func SimpleCheck(ip, port string) error {
	check := checkup.TCPChecker{
		URL: ip + ":" + port,
		//Attempts: 3,
		//UpStatus: config.StatusCode, //452, //above 500 consider error
	}
	r, err := check.Check()
	if err != nil {
		return err
	}
	if !r.Healthy {
		return fmt.Errorf("%v:%v check failed", ip, port)
	}
	return nil
}

// provided for hook-api, so here, it usually use wk name, the underscore one
func CheckIPWithConfig(name, ip, port string) error {
	config, err := FetchConfigByK8sName(name)
	if err != nil {
		// just log, later return err to the platform?
		return fmt.Errorf("fetch config for %v, err: %v", name, err)
	}

	check := GetCheckWithConfig(name, ip, port, config)

	r, err := check.Check()
	if err != nil {
		log.Debug.Printf("%v:%v check failed,err: %v\n", ip, port, err)
		return err
	}
	if !r.Healthy {
		log.Debug.Printf("%v:%v check failed,result: %v\n", ip, port, r)
		return fmt.Errorf("%v:%v check failed", ip, port)
	}
	return nil
}

func CheckLonger(name, ip, port string, t time.Duration) (err error) {
	start := time.Now()
	var interval = 1
	var i = 0
	for {
		i++
		err = CheckIPWithConfig(name, ip, port)
		if err != nil {
			log.Printf("check with config err: %v, fallback to simple tcp check\n", err)
			err = SimpleCheck(ip, port)
		}
		if err == nil {
			return
		}
		if i >= 3 {
			interval = interval*2 + 1
		}
		err = fmt.Errorf("interval: %v, tried %v times, err: %v\n", interval, i, err)
		log.Printf("simple check err: %v\n", err)

		if time.Now().Sub(start) >= t {
			return fmt.Errorf("check longer timeout, interval: %v, tried: %v times", interval, i)
		}
		time.Sleep(time.Duration(interval) * time.Millisecond * 100)
	}
}
