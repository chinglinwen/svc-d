package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

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
	port             = flag.String("p", "8080", "port")
	concurrentChecks = flag.Int("cc", 100, "number of concurrent checks")
	testproject      = flag.String("test", "ops_fs", "test project name")
	checkonetime     = flag.Bool("once", false, "check only once")
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

	if *testproject != "" {
		log.Printf("test for %v project only\n", *testproject)
		check.TestProject = *testproject
	}
	if *checkonetime {
		log.Println("check one time only")
		check.CheckOneTime = *checkonetime
	}
	fetch.Env = *env
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

	e.GET("/", homeHandler)
	e.POST("/check", checkHandler)
	e.POST("/notify", notifyHandler)

	err := e.Start(":" + *port)
	log.Println("fatal", err)
	//e.Logger.Fatal()

	log.Println("exit")
}
