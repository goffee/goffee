package main

import (
	"flag"

	"github.com/gophergala/goffee/probe"
	"github.com/gophergala/goffee/queue"
	"github.com/gophergala/goffee/scheduler"
	"github.com/gophergala/goffee/web"
	"github.com/gophergala/goffee/writer"
)

var webMode bool
var probeMode bool
var schedulerMode bool
var writerMode bool
var redisAddress string

func init() {
	flag.BoolVar(&webMode, "webmode", false, "Run goffee in webmode")
	flag.BoolVar(&probeMode, "torfetch", false, "Fetch something via Tor")
	flag.BoolVar(&schedulerMode, "scheduler", false, "Run goffee scheduler")
	flag.BoolVar(&writerMode, "writer", false, "Run goffee writer")
	flag.StringVar(&redisAddress, "redisaddress", "", "Address of redis including port")
	flag.Parse()
}

func main() {
	if webMode {
		web.StartServer()
	}

	if probeMode || schedulerMode || writerMode {
		queue.RedisAddressWithPort = redisAddress
	}

	if probeMode {
		probe.Run()
	}
	if schedulerMode {
		scheduler.Run()
	}
	if writerMode {
		writer.Run()
	}

	if webMode {
		web.Wait()
	}
	if writerMode {
		writer.Wait()
	}
	if schedulerMode {
		scheduler.Wait()
	}
	if probeMode {
		probe.Wait()
	}
}
