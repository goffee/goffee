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
	Id        int64
	URL       string `gorm:"column:url"`
	Status    int    // status code of last result
	Success   bool   // success status of last result
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Result struct {
	Id        int64
	CreatedAt time.Time
	Status    int
	Success   bool
	IP        string `gorm:"column:ip"`
	CheckId   int64
}

type User struct {
	Id          int64
	CreatedAt   time.Time
	Name        string
	Email       string
	GitHubId    int64  `gorm:"column:github_id"`
	GitHubLogin string `gorm:"column:github_login"`
	OAuthToken  string `gorm:"column:oauth_token"`
}

func InitDatabase() (err error) {
	db, err = gorm.Open("sqlite3", "/tmp/goffee.db")
	if err != nil {
		return err
	}

	db.AutoMigrate(&Check{}, &Result{}, &User{})

	return nil
}

func Checks() ([]Check, error) {
	var checks []Check
	res := db.Find(&checks)
	return checks, res.Error
}

func (c *Check) Create() error {
	res := db.Create(c)
	return res.Error
}

func (c *Check) AddResult(r *Result) error {
	tx := db.Begin()

	r.CheckId = c.Id
	res := tx.Create(r)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	c.Status = r.Status
	c.Success = r.Success
	res = tx.Save(c)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	tx.Commit()
	return nil
}

func (c *Check) Results() ([]Result, error) {
	var results []Result
	res := db.Model(c).Related(&results)
	return results, res.Error
}

func (u *User) UpdateOrCreate() error {
	res := db.Where(User{GitHubId: u.GitHubId}).Assign(*u).FirstOrInit(u)
	if res.Error != nil {
		return res.Error
	}

	res = db.Save(u)
	return res.Error
}
