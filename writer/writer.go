package writer

import (
	"encoding/json"
	"fmt"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/queue"
)

func Run() {
	data.InitDatabase()

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
				check.AddResult(&result)
				fmt.Println("Added result: %v", result)
			}
		}
	}
}
