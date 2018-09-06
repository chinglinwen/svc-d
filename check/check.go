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

	"wen/svc-d/config"
	"wen/svc-d/fetch"
)

var (
	CheckOneTime  bool
	TestProject   string
	CheckInterval int
)

// start backend check
func Start(conf *config.Config) {
	log.Println("started fetch in the background")

	var (
		ticker = time.NewTicker(time.Duration(CheckInterval) * time.Second)
		//stop   *time.Ticker
		now  time.Time
		prev time.Time
	)

	for ; ; now = <-ticker.C {
		log.Printf("starting fetch at %v, last time elapsed %v\n", now.Format(time.UnixDate), now.Sub(prev))
		// empty it for new fetch
		conf.Checkers = []checkup.Checker{}

		projects, err := fetch.Fetchs()
		if err != nil {
			log.Println("fetch upstream error", err)
			continue
		}
		configs, err := fetch.FetchConfigs()
		if err != nil {
			log.Println("fetch config error", err)
			//continue
		}

		conf.Checkers = fetch.ConvertToCheck(projects, configs)
		conf.Save()

		log.Println("fetch ok")

		err = conf.CheckAndStore()
		if err != nil {
			log.Println("background check error", err)
			continue
		}
		log.Printf("background check ok\n\n")

		prev = now

		if CheckOneTime {
			log.Println("check one time only, exit the loop")
			ticker.Stop()
			goto EXIT
		}
	}
EXIT:
	log.Println("background check stopped")
}

/* type setting struct {
	checkonetime bool
	testproject  string
}
type Option func(*setting)

func SetTestProject(t string) Option {
	return func(c *setting) {
		c.testproject = t
	}
}

func SetCheckOneTime(a bool) Option {
	return func(c *setting) {
		c.checkonetime = a
	}
} */

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

// provided for hook-api
func CheckIPWithConfig(name, ip, port string) error {
	config, err := fetch.FetchConfig(name)
	if err != nil {
		// just log, later return err to the platform?
		return fmt.Errorf("fetch config for %v, err: %v", name, err)
	}

	var check checkup.Checker
	if config.Type == "http" {
		c := checkup.HTTPChecker{
			Name: name,
			URL:  "http://" + ip + ":" + port + config.URI,
		}
		if config.StatusCode == 0 {
			c.UpStatus = 452 //above 500 consider error
		} else {
			c.UpStatus = config.StatusCode
		}
		if config.Attempts != 0 {
			c.Attempts = config.Attempts
		}
		if config.MustContain != "" {
			c.MustContain = config.MustContain
		}
		check = c
	} else {
		c := checkup.TCPChecker{
			Name: name, // notify api will use it
			URL:  ip + ":" + port,
		}
		if config.Timeout != 0 {
			c.Timeout = time.Duration(config.Timeout) * time.Second
		}
		check = c
	}

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
			log.Debug.Printf("check with config err: %v, fallback to simple tcp check\n", err)
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
			return fmt.Errorf("simple check timeout, interval: %v, tried: %v times", interval, i)
		}
		time.Sleep(time.Duration(interval) * time.Millisecond * 100)
	}
}
