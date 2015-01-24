package writer

import (
	"strconv"
	"strings"
	"time"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/queue"
)

// type Result struct {
//   Id        int64
//   CreatedAt time.Time
//   Status    int
//   Success   bool
//   IP        string `gorm:"column:ip"`
//   CheckId   int64
// }

func Run() {
	data.InitDatabase()

	for {
		results := queue.FetchResults()
		for _, result := range results {
			parts := strings.Split(result, " ")
			if len(parts) >= 3 {
				time, err := time.Parse(time.RFC3339, parts[0])
				if err != nil {
					continue
				}
				url := parts[1]
				status, err := strconv.Atoi(parts[2])
				if err != nil {
					continue
				}

				success := status >= 200 && status < 300

				checks, err := data.ChecksByURL(url)
				if err != nil {
					continue
				}

				for _, check := range checks {
					result := &data.Result{CreatedAt: time, Status: status, Success: success}
					check.AddResult(result)
				}
			}
		}
	}
}
