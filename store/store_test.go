package store

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/chinglinwen/checkup"
	"github.com/jmoiron/sqlx"
)

func init() {
	var err error
	DB, err = sqlx.Connect("postgres", "postgresql://postgres:123456@localhost/?sslmode=disable")
	if err != nil {
		fmt.Println("init err", err)
		os.Exit(1)
	}

	//db.MustExec(schema)
}
func TestSave(t *testing.T) {
	fmt.Println("start testing...")

	a := Project{
		Name:   "test",
		Region: "m7",
		Checkup: checkup.Checkup{
			Checkers: []checkup.Checker{
				checkup.HTTPChecker{
					Name:     "Website",
					URL:      "http://www.baidu.com",
					Attempts: 5,
				},
			},
			Storage: checkup.FS{
				Dir:         "./data",
				CheckExpiry: 7 * 24 * time.Hour,
			},
			Notifier: &checkup.Qianbao{Channel: "test url"},
		},
	}

	err := a.Save(DB)
	if err != nil {
		t.Error("err", err)
	}
}
