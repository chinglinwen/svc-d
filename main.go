package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"wen/svc-d/config"

	"github.com/chinglinwen/checkup"
	"github.com/chinglinwen/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	conf               *config.Config
	env                = flag.String("env", "test", "env includes (test,pre,pro)")
	port               = flag.String("p", "8080", "port")
	testproject        = flag.String("t", "ops_test", "test project")
	checkonetime       = flag.Bool("once", false, "check only once")
	concurrentChecks   = flag.Int("cc", 100, "number of concurrent checks")
	upstreamAPI        = flag.String("upstream", "http://upstream-pre.sched.qianbao-inc.com/get_upstream_all_instance/", "upstream fetch api url")
	upstreamnChangeAPI = flag.String("upstreamc", "http://upstream-pre.sched.qianbao-inc.com/up_nginx_state/", "upstream change api url")
)

// try have two config
// one for fetch and one for manual editing
func init() {
	flag.Parse()
	if *env == "test" {
		*upstreamAPI = ""
	}
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
}

//define a global variable
//add new check, update it, and store the config as file(update config)

func main() {
	log.Println("starting...")
	log.Debug.Println("debug is on")
	c, _ := json.MarshalIndent(conf, "", " ")
	log.Println("config:", string(c))
	go start()

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
