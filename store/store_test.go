package store

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sourcegraph/checkup"
)

func init() {
	var err error
	DB, err = sqlx.Connect("postgres", "postgresql://postgres:123456@localhost/?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
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
