package scheduler

import (
	"time"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/queue"
)

var exit = make(chan bool)

func Run() {
	go run()
}

func run() {
	scheduleChecks()

	for range time.Tick(15 * time.Second) {
		scheduleChecks()
	}
}

func scheduleChecks() {
	if !queue.AcquireSchedulerLock(60, 300) {
		return
	}

	if checks, err := data.Checks(); err == nil {
		for _, check := range checks {
			queue.AddJob(check.URL)
		}
	}
	queue.ReleaseSchedulerLock()
}

func Wait() {
	<-exit
}
