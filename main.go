package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"wen/hook-api/upstream"
	"wen/svc-d/check"
	"wen/svc-d/config"
	"wen/svc-d/fetch"

	"github.com/chinglinwen/checkup"
	"github.com/chinglinwen/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	conf             *config.Config
	env              = flag.String("env", "qa", "env includes (qa,pre,pro)")
	port             = flag.String("p", "8089", "port")
	checkInterval    = flag.Int("i", 10, "check interval in seconds")
	concurrentChecks = flag.Int("cc", 100, "number of concurrent checks")
	testproject      = flag.String("test", "", "test project name")
	checkonetime     = flag.Bool("once", false, "check only once")
	dockerOnly       = flag.Bool("docker", true, "check docker only")

	upstreamBase = flag.String("upstream", "http://upstream-test.sched.qianbao-inc.com:8010", "upstream base api url")
)

// try have two config
// one for fetch and one for manual editing
func init() {
	flag.Parse()
	conf = config.New("config.json", *env) //try the item, project based ?  why not just name?

	_ = os.Mkdir("data", os.ModeDir)
	conf.Notifier = checkup.Qianbao{
		Username: "wen",
		Channel:  "http://localhost:" + *port + "/notify",
	}
	conf.Storage = checkup.FS{
		Dir:         "data",
		CheckExpiry: 7 * 24 * time.Hour,
	}
	conf.ConcurrentChecks = *concurrentChecks
	conf.Save()

	check.CheckInterval = *checkInterval
	if *testproject != "" {
		log.Printf("test for %v project only\n", *testproject)
		check.TestProject = *testproject
	}
	if *checkonetime {
		log.Println("check one time only")
		check.CheckOneTime = *checkonetime
	}
	fetch.Env = *env
	fetch.DockerOnly = *dockerOnly

	fetch.Init(*upstreamBase)
	upstream.Init(*upstreamBase)
	log.Println("using upstream", *upstreamBase)
}

// define a global variable
// add new check, update it, and store the config as file(update config)

func main() {
	log.Println("starting...")
	log.Debug.Println("debug is on")
	c, _ := json.MarshalIndent(conf, "", " ")
	log.Println("config:", string(c))

	go check.Start(conf)

	e := echo.New()
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Static("/static"))

	e.Static("/", "statuspage")
	e.Static("/data", "data")
	e.GET("/check", checkHandler)
	e.POST("/notify", notifyHandler)

	err := e.Start(":" + *port)
	log.Println("fatal", err)
	//e.Logger.Fatal()

	log.Println("exit")
}
