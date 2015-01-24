package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/queue"
	"github.com/gophergala/goffee/scheduler"
	"github.com/gophergala/goffee/tor"
	"github.com/gophergala/goffee/web"
	"github.com/gophergala/goffee/writer"
)

var webMode bool
var torFetch bool
var schedulerMode bool
var writerMode bool
var redisAddress string

const ipReflector = "http://stephensykes.com/ip_reflection.html"

func init() {
	flag.BoolVar(&webMode, "webmode", false, "Run goffee in webmode")
	flag.BoolVar(&torFetch, "torfetch", false, "Fetch something via Tor")
	flag.BoolVar(&schedulerMode, "scheduler", false, "Run goffee scheduler")
	flag.BoolVar(&writerMode, "writer", false, "Run goffee writer")
	flag.StringVar(&redisAddress, "redisaddress", "", "Address of redis including port")
	flag.Parse()
}

func main() {
	if webMode {
		web.StartServer()
		return
	}

	queue.RedisAddressWithPort = redisAddress

	if torFetch {
		var wg sync.WaitGroup

		for { // ever
			fmt.Println("Redis fetch")
			batch := queue.FetchBatch()
			fmt.Println(batch)
			for _, item := range batch {
				if item == "newip" {
					newip()
				} else {
					wg.Add(1)
					go check(item, &wg)
				}
			}

			wg.Wait()
		}
	}

	if schedulerMode {
		scheduler.Run()
	}

	if writerMode {
		writer.Run()
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
