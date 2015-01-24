package main

import (
	"fmt"
	"github.com/gophergala/goffee/tor"

	"flag"

	"github.com/gophergala/goffee/web"
	
	"github.com/gophergala/goffee/queue"
	
	"sync"
)

var webMode bool
var torFetch bool
var redisAddress string

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
					ip, _ := tor.NewIP()
					fmt.Println(ip) // may not print anything, info from tor is not reliable
				} else {
					wg.Add(1)
					go check(item, &wg)					
				}
			}
			
			wg.Wait()
		}
	}
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
	
	queue.WriteResult(address + " " + result)
}
