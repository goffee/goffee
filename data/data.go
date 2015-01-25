package data

import (
	"time"

	"github.com/gophergala/goffee/Godeps/_workspace/src/github.com/jinzhu/gorm"
	_ "github.com/gophergala/goffee/Godeps/_workspace/src/github.com/lib/pq"
	_ "github.com/gophergala/goffee/Godeps/_workspace/src/github.com/mattn/go-sqlite3" // DB adapters
)

var db gorm.DB

type Check struct {
	Id        int64
	URL       string `gorm:"column:url"`
	Status    int    // status code of last result
	Success   bool   // success status of last result
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Result struct {
	Id        int64
	CreatedAt time.Time
	Status    int
	Success   bool
	IP        string `gorm:"column:ip"`
	Country   string
	CheckId   int64

	URL string `sql:"-"`
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

	db.Model(&User{}).AddIndex("idx_user_github_id", "github_id")
	db.Model(&Check{}).AddIndex("idx_check_id_user_id", "id", "user_id")

	return nil
}

func Checks() (checks []Check, err error) {
	res := db.Find(&checks)
	return checks, res.Error
}

func ChecksByURL(url string) (checks []Check, err error) {
	res := db.Where("url = ?", url).Find(&checks)
	return checks, res.Error
}

func (u *User) Checks() (checks []Check, err error) {
	res := db.Model(u).Related(&checks)
	return checks, res.Error
}

func (u *User) Check(id int64) (check Check, err error) {
	res := db.Where("user_id = ? and id = ?", u.Id, id).First(&check)
	return check, res.Error
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

func (c *Check) Results() (results []Result, err error) {
	res := db.Model(c).Order("id desc").Limit(120).Related(&results)
	return results, res.Error
}

func (c *Check) User() (user User, err error) {
	res := db.Model(c).Related(&user)
	return user, res.Error
}

func (c *Check) Delete() (err error) {
	tx := db.Begin()

	res := db.Where(&Result{CheckId: c.Id}).Delete(Result{})
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	res = db.Delete(c)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	tx.Commit()
	return nil
}

func (u *User) UpdateOrCreate() error {
	res := db.Where(User{GitHubId: u.GitHubId}).Assign(*u).FirstOrInit(u)
	if res.Error != nil {
		return res.Error
	}

	res = db.Save(u)
	return res.Error
}

func FindUser(id int64) (user User, err error) {
	res := db.First(&user, id)
	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}

func FindCheck(id int64) (check Check, err error) {
	res := db.First(&check, id)
	if res.Error != nil {
		return check, res.Error
	}
	return check, nil
}
