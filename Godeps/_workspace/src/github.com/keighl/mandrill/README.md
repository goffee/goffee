# Mandrill Emails via Golang

[![Build Status](https://travis-ci.org/keighl/mandrill.png?branch=master)](https://travis-ci.org/keighl/mandrill)

Stripped down package for sending emails through the Mandrill API. Inspired by [@mostafah's implementation](https://github.com/mostafah/mandrill).

### Installation

    go get -u github.com/keighl/mandrill

### Documentation

http://godoc.org/github.com/keighl/mandrill

### Regular Message

https://mandrillapp.com/api/docs/messages.JSON.html#method=send

    import (
      m "github.com/keighl/mandrill"
    )

    client := m.ClientWithKey("y2cQvBBfdFoZNByVaKsJsA")

    message := &m.Message{}
    message.AddRecipient("bob@example.com", "Bob Johnson", "to")
    message.FromEmail = "kyle@example.com"
    message.FromName = "Kyle Truscott"
    message.Subject = "You won the prize!"
    message.HTML = "<h1>You won!!</h1>"
    message.Text = "You won!!"

    responses, apiError, err := client.MessagesSend(message)

### Send Template

https://mandrillapp.com/api/docs/messages.JSON.html#method=send-template

http://help.mandrill.com/entries/21694286-How-do-I-add-dynamic-content-using-editable-regions-in-my-template-

    templateContent := map[string]string{"header": "Bob! You won the prize!"}
    responses, apiError, err := client.MessagesSendTemplate(message, "you-won", templateContent)

### Including Merge Tags

http://help.mandrill.com/entries/21678522-How-do-I-use-merge-tags-to-add-dynamic-content-

    message.GlobalMergeVars := m.ConvertMapToVariables(map[string]string{"name": "Bob"})
    message.MergeVars := m.ConvertMapToVariablesForRecipient("bob@example.com", map[string]string{"name": "Bob"})
