package notifier

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gophergala/goffee/Godeps/_workspace/src/github.com/keighl/mandrill"
	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/queue"
)

var (
	exit           = make(chan bool)
	MandrillKey    string
	mandrillClient *mandrill.Client
)

func Run() {
	go run()
}

func run() {
	mandrillClient = mandrill.ClientWithKey(MandrillKey)

	for {
		notifications := queue.FetchNotifications()
		for _, n := range notifications {
			checkId, err := strconv.ParseInt(n, 10, 64)

			check, err := data.FindCheck(checkId)
			if err != nil {
				continue
			}

			user, err := check.User()
			if err != nil {
				continue
			}

			sendMessage(check, user)
		}
	}
}

func sendMessage(c data.Check, u data.User) {
	println("Notifying via email: " + u.Email)

	var subject string

	if c.Success {
		subject = fmt.Sprintf("Up: %s (%d)", c.URL, c.Status)
	} else {
		subject = fmt.Sprintf("Down: %s (%d)", c.URL, c.Status)
	}

	html := `<strong>%s</strong>
  <br>
  <br>
  <p>Checked at %s by <a href='http://goffee.io/'>Goffee.io</a></p>`
	html = fmt.Sprintf(html, subject, c.UpdatedAt.Format(time.Kitchen))

	text := `%s\n\nChecked at %s by Goffee.io`
	text = fmt.Sprintf(text, subject, c.UpdatedAt.Format(time.Kitchen))

	message := &mandrill.Message{}
	message.AddRecipient(u.Email, u.Name, "to")
	message.FromEmail = "no-reply@goffee.io"
	message.FromName = "Goffee Notifier"
	message.Subject = subject
	message.HTML = html
	message.Text = text
	message.Subaccount = "goffee"

	mandrillClient.MessagesSend(message)
}

func Wait() {
	<-exit
}
