package scheduler

import (
	"time"

	"github.com/goffee/goffee/data"
	"github.com/goffee/goffee/queue"
)

var (
	stop = make(chan bool)
	exit = make(chan bool)
)

func Run() {
	go run()
}

func run() {
	scheduleChecks()

	for {
		select {
		case <-time.Tick(15 * time.Second):
			scheduleChecks()
		case <-stop:
			exit <- true
		}
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

func Stop() {
	stop <- true
}

func Wait() {
	<-exit
}
