package scheduler

import (
	"fmt"
	"time"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/queue"
)

var RedisAddressWithPort string

func Run() {
	data.InitDatabase()

	if checks, err := data.Checks(); err == nil {
		scheduleChecks(checks)
	}

	for range time.Tick(10 * time.Second) {
		if checks, err := data.Checks(); err == nil {
			scheduleChecks(checks)
		}
	}
}

func scheduleChecks(checks []data.Check) {
	for _, check := range checks {
		fmt.Println("Job added:", check.URL)
		queue.AddJob(check.URL)
	}
}
