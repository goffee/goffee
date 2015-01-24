package main

import (
	"flag"
	"fmt"

	"github.com/goffee/web"
)

var webMode bool

func init() {
	flag.BoolVar(&webMode, "webmode", false, "Run goffee in webmode")
	flag.Parse()
}

func main() {
	if webMode {
		web.StartServer()
	}
	fmt.Println("Hello, Gopher Gala!")
}
