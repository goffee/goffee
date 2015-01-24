package main

import (
	"fmt"
	"github.com/gophergala/goffee/tor"

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
	body, _ := tor.TorGet("http://www.kiskolabs.com")
        fmt.Printf("Result was '%s'\n", string(body))
}
