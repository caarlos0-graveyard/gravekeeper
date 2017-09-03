package github

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const msg = "Hi! Thanks for bringing that up, but this repository [is not being maintained anymore](https://github.com/caarlos0/gravekeeper)."

// Repository is a github repo
type Repository struct {
	FullName string `json:"full_name,omitempty"`
}

// Issue is a github issue
type Issue struct {
	Number int64 `json:"number,omitempty"`
}

// Payload is a payload from a github webhook
type Payload struct {
	Action      string     `json:"action,omitempty"`
	Issue       Issue      `json:"issue,omitempty"`
	PullRequest Issue      `json:"pull_request,omitempty"`
	Repository  Repository `json:"repository,omitempty"`
}

// IsCreate returns true if the event action is created or opened
func (p Payload) IsCreate() bool {
	return p.Action == "created" || p.Action == "openend"
}

// Notify notify the issue/pr related to the payload
func (p Payload) Notify(api, token string) error {
	var url = p.commentURL(api)
	var body = []byte(fmt.Sprintf(`{"body": "%s"}`, msg))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	resp, err := http.DefaultClient.Do(req)
	bts, rerr := ioutil.ReadAll(resp.Body)
	if rerr != nil {
		log.Println(rerr.Error())
	}
	log.Println(string(bts))
	return err
}

func (p Payload) commentURL(api string) string {
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
