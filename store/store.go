// not using for now
package store

import (
	"encoding/json"
	"fmt"

	"github.com/chinglinwen/checkup"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// project based, need to write concurrency checks
type Project struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Region    string `db:"region"`
	RawChecks string `db:"checks"`
	checkup.Checkup
}

var DB *sqlx.DB

func (p *Project) Save(db *sqlx.DB) error {
	project := `INSERT INTO projects (name, region,checks ) VALUES (?, ?, ?)`

	b, err := json.Marshal(p.Checkup)
	if err != nil {
		return err
	}
	p.RawChecks = string(b)
	_, err = db.Exec(project, p.Name, p.Region, p.RawChecks)
	return err
}

func (p *Project) Read(db *sqlx.DB) error {
	if p.Name == "" {
		return fmt.Errorf("name not provided for query")
	}

	row := db.QueryRow("SELECT id,name, region,checks FROM projects WHERE name=?", p.Name)
	var a Project
	return row.Scan(&a)
}

/*
func (p *Project) Save(db *bolt.DB) error {
	// Store the user model in the user bucket using the username as the key.
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(usersBucket)
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return b.Put([]byte(user.Name), encoded)
	})
	return err
}
*/
