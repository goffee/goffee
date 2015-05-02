package writer

import (
	"encoding/json"
	"strconv"

	"github.com/goffee/goffee/data"
	"github.com/goffee/goffee/queue"
)

var stop = make(chan bool)
var exit = make(chan bool)

func Run() {
	go run()
}

func run() {
	for {
		select {
		case <-stop:
			exit <- true
		case results := <-queue.FetchResults():
			for _, r := range results {
				var result data.Result
				err := json.Unmarshal([]byte(r), &result)

				checks, err := data.ChecksByURL(result.URL)
				if err != nil {
					continue
				}

				for _, check := range checks {
					previousSuccess := check.Success

					check.AddResult(&result)

					if previousSuccess && !result.Success {
						queue.AddNotification(strconv.FormatInt(check.Id, 10))
					} else if !previousSuccess && result.Success {
						queue.AddNotification(strconv.FormatInt(check.Id, 10))
					}
				}
			}
		}
	}
}

func Stop() {
	stop <- true
}

func Wait() {
	<-exit
}
