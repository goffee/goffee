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
var currentIP string
var lastIPChange time.Time

func Run() {
	go run()
}

func run() {
	var wg sync.WaitGroup
	newip()

	for { // ever
		fmt.Println("Redis fetch")
		batch := queue.FetchBatch()
		fmt.Println(batch)
		for _, item := range batch {
			wg.Add(1)
			go check(item, &wg)
		}

		wg.Wait()
		
		if time.Since(lastIPChange) > time.Minute {
			newip()
		}
		
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
	currentIP = result
	fmt.Println("New IP obtained " + result)
	lastIPChange = time.Now()
}

func check(address string, wg *sync.WaitGroup) {
	defer wg.Done()

	status, err := tor.TorGetStatus(address)

	var statusCode int
	if err != nil {
		statusCode = -1
	} else {
		statusCode, err = strconv.Atoi(strings.Split(status, " ")[0])
		if err != nil {
			statusCode = -2
		}
	}

	result := &data.Result{
		CreatedAt: time.Now(),
		Status:    statusCode,
		Success:   statusCode >= 200 && statusCode < 300,
		URL:       address,
		IP:        currentIP,
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
