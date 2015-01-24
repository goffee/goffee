package probe

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/queue"
	"github.com/gophergala/goffee/tor"
)

const ipReflector = "http://stephensykes.com/ip_reflection.html"

var exit = make(chan bool)

func Run() {
	go run()
}

func run() {
	var wg sync.WaitGroup

	for { // ever
		fmt.Println("Redis fetch")
		batch := queue.FetchBatch()
		fmt.Println(batch)
		for _, item := range batch {
			if item == "newip" {
				wg.Wait()
				newip()
			} else {
				wg.Add(1)
				go check(item, &wg)
			}
		}

		wg.Wait()
	}
}

func newip() {
	tor.NewIP()
	body, err := tor.TorGet(ipReflector)
	var result string
	if err != nil {
		result = err.Error()
	} else {
		result = body
	}
	queue.WriteResult(time.Now().Format(time.RFC3339) + " newip " + result)
}

func check(address string, wg *sync.WaitGroup) {
	defer wg.Done()

	status, err := tor.TorGetStatus(address)
	if err != nil {
		return
	}

	statusCode, err := strconv.Atoi(strings.Split(status, " ")[0])
	if err != nil {
		return
	}

	result := &data.Result{
		CreatedAt: time.Now(),
		Status:    statusCode,
		Success:   statusCode >= 200 && statusCode < 300,
		URL:       address,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return
	}

	queue.WriteResult(string(data))
}

func Wait() {
	<-exit
}
