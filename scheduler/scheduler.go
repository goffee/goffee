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
		queue.AddJob(check.URL)
	}
}

func Wait() {
	<-exit
}
