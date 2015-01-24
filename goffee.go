package main

import (
	"fmt"

	"github.com/gophergala/goffee/tor"

	"flag"

	"github.com/gophergala/goffee/web"

	"github.com/gophergala/goffee/queue"

	"sync"
	"time"
)

var webMode bool
var torFetch bool
var redisAddress string
const ipReflector = "http://stephensykes.com/ip_reflection.html"

func init() {
	flag.BoolVar(&webMode, "webmode", false, "Run goffee in webmode")
	flag.BoolVar(&torFetch, "torfetch", false, "Fetch something via Tor")
	flag.StringVar(&redisAddress, "redisaddress", "", "Address of redis including port")
	flag.Parse()
}

func main() {
	if webMode {
		web.StartServer()
	}

	if torFetch {
		queue.RedisAddressWithPort = redisAddress

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
	var result string
	if err != nil {
		result = err.Error()
	} else {
		result = status
	}

	queue.WriteResult(time.Now().Format(time.RFC3339) + " " + address + " " + result)
}
