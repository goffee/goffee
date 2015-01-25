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

const ipReflector = "https://goffee.herokuapp.com/ip"
const ipRefreshInterval = 5 * time.Minute

type IPResponse struct {
	IP      string
	Country string
}

var (
	exit           = make(chan bool)
	lastIPChange   time.Time
	currentIP      string
	currentCountry string
)

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

		if time.Since(lastIPChange) > ipRefreshInterval {
			newip()
		}
	}
}

func newip() {
	tor.NewIP()
	body, err := tor.TorGet(ipReflector)
	if err != nil {
		return
	}
	var response IPResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return
	}
	currentIP = response.IP
	currentCountry = response.Country
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
		Country:   currentCountry,
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
