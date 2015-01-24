package main

import (
	"fmt"
	"github.com/gophergala/goffee/tor"
)

func main() {
	fmt.Println("Hello, Gopher Gala!")
	body, _ := tor.TorGet("http://www.kiskolabs.com")
        fmt.Printf("Result was '%s'\n", string(body))
}
