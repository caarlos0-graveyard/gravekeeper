package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	var github = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"ok":true}`)
	}))
	defer github.Close()
	var server = httptest.NewServer(handler(github.URL, "whatever"))
	defer server.Close()
	const (
		doneResult         = "done\n"
		invalidEventResult = "not a valid event, ignoring...\n"
	)
	for payload, event := range map[string]struct{ event, result string }{
		"ping.json":         {"ping", invalidEventResult},
		"opened_pr.json":    {"pull_request", doneResult},
		"reopened_pr.json":  {"pull_request", invalidEventResult},
		"opened_issue.json": {"issues", doneResult},
		"invalid.json":      {"pull_request", "invalid character '{' looking for beginning of object key string\n"},
	} {
		t.Run(payload, func(t *testing.T) {
			var assert = assert.New(t)
			resp, err := postJSON(t, server, payload, event.event)
			assert.NoError(err)
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(err)
			assert.Equal(event.result, string(body))
		})
	}
}

func postJSON(t *testing.T, server *httptest.Server, payload, event string) (*http.Response, error) {
	var assert = assert.New(t)
	f, err := os.Open("testdata/" + payload)
	assert.NoError(err)
	req, err := http.NewRequest(http.MethodPost, server.URL, f)
	req.Header.Add("X-GitHub-Event", event)
	assert.NoError(err)
	return http.DefaultClient.Do(req)
}
