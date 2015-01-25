package writer

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/queue"
)

var exit = make(chan bool)

func Run() {
	go run()
}

func run() {
	for {
		results := queue.FetchResults()
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
				fmt.Printf("Added result: %v\n", result)

				if previousSuccess && !result.Success {
					queue.AddNotification(strconv.FormatInt(check.Id, 10))
				} else if !previousSuccess && result.Success {
					queue.AddNotification(strconv.FormatInt(check.Id, 10))
				}
			}
		}
	}
}

func Wait() {
	<-exit
}
