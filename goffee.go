package main

import (
	"fmt"
	"github.com/gophergala/goffee/tor"

	"flag"

	"github.com/gophergala/goffee/web"
)

var webMode bool
var torFetch bool

func init() {
	flag.BoolVar(&webMode, "webmode", false, "Run goffee in webmode")
	flag.BoolVar(&torFetch, "torfetch", false, "Fetch something via Tor")
	flag.Parse()
}

func main() {
	if webMode {
		web.StartServer()
	}

	if torFetch {
		ip, _ := tor.NewIP()
		fmt.Println(ip) // may not print anything, info from tor is not reliable
		status, _ := tor.TorGetStatus("http://www.apple.com")
		fmt.Println(status)
	}
}
