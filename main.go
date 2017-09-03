package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func shouldHandle(req *http.Request) bool {
	var event = req.Header.Get("X-GitHub-Event")
	return event == "issues" || event == "pull_request"
}

type Repository struct {
	FullName string `json:"full_name,omitempty"`
}

type Issue struct {
	Number int64 `json:"number,omitempty"`
}

type Payload struct {
	Action      string     `json:"action,omitempty"`
	Issue       Issue      `json:"issue,omitempty"`
	PullRequest Issue      `json:"pull_request,omitempty"`
	Repository  Repository `json:"repository,omitempty"`
}

func (p Payload) GetNumber() int64 {
	if p.Issue.Number != 0 {
		return p.Issue.Number
	}
	return p.PullRequest.Number
}

func (p Payload) IsCreate() bool {
	return p.Action == "created" || p.Action == "openend"
}

func (p Payload) CommentURL(api string) string {
	var number = p.PullRequest.Number
	if p.Issue.Number != 0 {
		number = p.Issue.Number
	}
	return fmt.Sprintf(
		"%s/repos/%s/issues/%d/comments",
		api,
		p.Repository.FullName,
		number,
	)
}

var body = []byte(`
	{
		"body": "Hi! Thanks for bringing that up, but this repo was [graveyarded](https://github.com/caarlos0/gravekeeper) and is not being actively maintained anymore."
	}
`)

func comment(p Payload) error {
	var url = p.CommentURL("https://api.github.com")
	log.Println("URL:", url)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN")))
	res, err := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	log.Println(string(body))
	return err
}

func main() {
	var addr = ":" + os.Getenv("PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !shouldHandle(r) {
			fmt.Fprintln(w, "Will not handle this event")
			return
		}
		var evt Payload
		if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if evt.IsCreate() {
			fmt.Fprintln(w, "Will not handle this event")
			return
		}
		log.Println("Event:", evt)
		if err := comment(evt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Done.")
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
