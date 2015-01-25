package mandrill

import (
  "fmt"
  "testing"
  "reflect"
  "net/http/httptest"
  "net/http"
  "net/url"
)

func expect(t *testing.T, a interface{}, b interface{}) {
  if a != b {
    t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
  }
}

func refute(t *testing.T, a interface{}, b interface{}) {
  if a == b {
    t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
  }
}

func testTools(code int, body string) (*httptest.Server, *Client)  {

  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(code)
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintln(w, body)
  }))

  tr := &http.Transport{
    Proxy: func(req *http.Request) (*url.URL, error) {
      return url.Parse(server.URL)
    },
  }
  httpClient := &http.Client{Transport: tr}

  client := &Client{"APIKEY", server.URL, httpClient}
  return server, client
}

// MessagesSendTemplate //////////

func Test_MessagesSendTemplate_Success(t *testing.T) {
  server, m := testTools(200, `[{"email":"bob@example.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
  defer server.Close()
  responses, apiError, _ := m.MessagesSendTemplate(&Message{}, "cheese", map[string]string{"name": "bob"})

  expect(t, len(responses), 1)
  expect(t, apiError, (*Error)(nil))

  correctResponse := &Response{
    Email: "bob@example.com",
    Status: "sent",
    RejectionReason: "hard-bounce",
    Id: "1",
  }
  expect(t, reflect.DeepEqual(correctResponse, responses[0]), true)
}

func Test_MessagesSendTemplate_Fail(t *testing.T) {
  server, m := testTools(400, `{"status":"error","code":12,"name":"Unknown_Subaccount","message":"No subaccount exists with the id 'customer-123'"}`)
  defer server.Close()
  responses, apiError, _ := m.MessagesSendTemplate(&Message{}, "cheese", map[string]string{"name": "bob"})

  expect(t, len(responses), 0)

  correctResponse := &Error{
    Status: "error",
    Code: 12,
    Name: "Unknown_Subaccount",
    Message: "No subaccount exists with the id 'customer-123'",
  }
  expect(t, reflect.DeepEqual(correctResponse, apiError), true)
}

// MessagesSend //////////

func Test_MessageSend_Success(t *testing.T) {
  server, m := testTools(200, `[{"email":"bob@example.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
  defer server.Close()
  responses, apiError, _ := m.MessagesSend(&Message{})

  expect(t, len(responses), 1)
  expect(t, apiError, (*Error)(nil))

  correctResponse := &Response{
    Email: "bob@example.com",
    Status: "sent",
    RejectionReason: "hard-bounce",
    Id: "1",
  }
  expect(t, reflect.DeepEqual(correctResponse, responses[0]), true)
}

func Test_MessageSend_Fail(t *testing.T) {
  server, m := testTools(400, `{"status":"error","code":12,"name":"Unknown_Subaccount","message":"No subaccount exists with the id 'customer-123'"}`)
  defer server.Close()
  responses, apiError, _ := m.MessagesSend(&Message{})

  expect(t, len(responses), 0)

  correctResponse := &Error{
    Status: "error",
    Code: 12,
    Name: "Unknown_Subaccount",
    Message: "No subaccount exists with the id 'customer-123'",
  }
  expect(t, reflect.DeepEqual(correctResponse, apiError), true)
}

// AddRecipient //////////

func Test_AddRecipient(t *testing.T) {
  m := &Message{}
  m.AddRecipient("bob@example.com", "Bob Johnson", "to")
  tos := []*To{&To{"bob@example.com", "Bob Johnson", "to"}}
  expect(t, reflect.DeepEqual(m.To, tos), true)
}

// ConvertMapToVariables /////

func Test_ConvertMapToVariables(t *testing.T) {
  m := map[string]string{"name": "bob"}
  target := ConvertMapToVariables(m)
  hand := []*Variable{&Variable{"name", "bob"}}
  expect(t, reflect.DeepEqual(target, hand), true)
}

// ConvertMapToVariablesForRecipient ////

func Test_ConvertMapToVariablesForRecipient(t *testing.T) {
  m := map[string]string{"name": "bob"}
  target := ConvertMapToVariablesForRecipient("bob@example.com", m)
  hand := &RcptMergeVars{"bob@example.com", ConvertMapToVariables(m)}
  expect(t, reflect.DeepEqual(target, hand), true)
}

