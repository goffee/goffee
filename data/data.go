package data

import (
	"time"

	"github.com/jinzhu/gorm"

	// DB adapters
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var db gorm.DB

type Check struct {
	Id  int64
	URL string
}

type Result struct {
	Timestamp time.Time
	Status    int
	Success   bool
	IP        string
}

func InitDatabase() (err error) {
	db, err = gorm.Open("sqlite3", "/tmp/goffee.db")
	if err != nil {
		return err
	}

	db.AutoMigrate(&Check{}, &Result{})

	return nil
}

func (c Check) Create() error {
	r := db.Create(c)
	return r.Error
}

func Checks() ([]Check, error) {
	checks := []Check{}
	r := db.Find(&checks)
	return checks, r.Error
}
