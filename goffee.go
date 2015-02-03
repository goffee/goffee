package main

import (
	"flag"
	"log"

	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/gorilla/sessions"
	"github.com/goffee/goffee/Godeps/_workspace/src/golang.org/x/oauth2"
	"github.com/goffee/goffee/Godeps/_workspace/src/golang.org/x/oauth2/github"
	"github.com/goffee/goffee/data"
	"github.com/goffee/goffee/notifier"
	"github.com/goffee/goffee/probe"
	"github.com/goffee/goffee/queue"
	"github.com/goffee/goffee/scheduler"
	"github.com/goffee/goffee/web"
	"github.com/goffee/goffee/web/controllers"
	"github.com/goffee/goffee/writer"
)

var webMode bool
var probeMode bool
var schedulerMode bool
var writerMode bool
var notifierMode bool
var redisAddress string
var bind string
var mysql string

func init() {
	var gitHubClientID string
	var gitHubClientSecret string
	var mandrillKey string
	var sessionSecret string

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
	flag.StringVar(&mysql, "mysql", "", "MySQL connection string")
	flag.StringVar(&sessionSecret, "sessionsecret", "", "The session secret for the web UI")

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
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		}

		if sessionSecret == "" {
			log.Fatal("No session secret set!")
		}

		web.SessionStore = sessions.NewCookieStore([]byte(sessionSecret))
	}
}

func main() {
	if webMode {
		web.StartServer(bind)
	}

	if probeMode || schedulerMode || writerMode || notifierMode {
		queue.InitQueue(redisAddress)
	}

	if webMode || schedulerMode || writerMode || notifierMode {
		if mysql != "" {
			data.InitDatabase("mysql", mysql)
		} else {
			data.InitDatabase("sqlite3", "/tmp/goffee.db")
		}
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
