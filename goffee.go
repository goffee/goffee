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
var bind string

func init() {
	flag.BoolVar(&webMode, "webmode", false, "Run goffee in webmode")
	flag.BoolVar(&probeMode, "torfetch", false, "Fetch something via Tor")
	flag.BoolVar(&schedulerMode, "scheduler", false, "Run goffee scheduler")
	flag.BoolVar(&writerMode, "writer", false, "Run goffee writer")
	flag.StringVar(&redisAddress, "redisaddress", "127.0.0.1:6379", "Address of redis including port")
	flag.StringVar(&bind, "bind", "127.0.0.1:8000", "Address to bind to")
	flag.Parse()

	// If no mode has been defined, just launch them all!
	if !webMode && !probeMode && !schedulerMode && !writerMode {
		webMode = true
		probeMode = true
		schedulerMode = true
		writerMode = true
	}
}

func main() {
	if webMode {
		web.StartServer(bind)
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
