package main

import (
	"flag"
	"log"

	"github.com/gophergala/goffee/Godeps/_workspace/src/golang.org/x/oauth2"
	"github.com/gophergala/goffee/Godeps/_workspace/src/golang.org/x/oauth2/github"
	"github.com/gophergala/goffee/notifier"
	"github.com/gophergala/goffee/probe"
	"github.com/gophergala/goffee/queue"
	"github.com/gophergala/goffee/scheduler"
	"github.com/gophergala/goffee/web"
	"github.com/gophergala/goffee/web/controllers"
	"github.com/gophergala/goffee/writer"
)

var webMode bool
var probeMode bool
var schedulerMode bool
var writerMode bool
var notifierMode bool
var redisAddress string
var bind string

func init() {
	var gitHubClientID string
	var gitHubClientSecret string
	var mandrillKey string

	flag.BoolVar(&webMode, "webmode", false, "Run goffee in webmode")
	flag.BoolVar(&probeMode, "torfetch", false, "Fetch something via Tor")
	flag.BoolVar(&schedulerMode, "scheduler", false, "Run goffee scheduler")
	flag.BoolVar(&writerMode, "writer", false, "Run goffee writer")
	flag.BoolVar(&notifierMode, "notifier", false, "Run goffee notifier")
	flag.StringVar(&redisAddress, "redisaddress", "127.0.0.1:6379", "Address of redis including port")
	flag.StringVar(&bind, "bind", "127.0.0.1:8000", "Address to bind to")
	flag.StringVar(&gitHubClientID, "clientid", "", "Github client ID")
	flag.StringVar(&gitHubClientSecret, "secret", "", "GitHub client Secret")
	flag.StringVar(&mandrillKey, "mandrill", "", "Mandrill API key")

	flag.Parse()

	// If no mode has been defined, just launch them all!
	if !webMode && !probeMode && !schedulerMode && !writerMode && !notifierMode {
		webMode = true
		probeMode = true
		schedulerMode = true
		writerMode = true
		notifierMode = true
	}

	if notifierMode {
		if mandrillKey == "" {
			log.Fatal("No Mandrill API key set!")
		}
		notifier.MandrillKey = mandrillKey
	}

	if webMode {
		if gitHubClientID == "" || gitHubClientSecret == "" {
			log.Fatal("No GitHub clientid or secret set!")
		}

		controllers.OAuthConf = &oauth2.Config{
			ClientID:     gitHubClientID,
			ClientSecret: gitHubClientSecret,
			Scopes:       []string{},
			Endpoint:     github.Endpoint,
		}
	}
}

func main() {
	if webMode {
		web.StartServer(bind)
	}

	if probeMode || schedulerMode || writerMode || notifierMode {
		queue.InitQueue(redisAddress)
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
	if notifierMode {
		notifier.Run()
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
	if notifierMode {
		notifier.Wait()
	}
}
