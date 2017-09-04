package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert"
)

func TestIsCreate(t *testing.T) {
	for _, state := range []string{"created", "opened"} {
		t.Run(state, func(t *testing.T) {
			assert.True(t, Payload{Action: state}.IsCreate())
		})
	}
	for _, state := range []string{"new_comment", "whatever-else"} {
		t.Run(state, func(t *testing.T) {
			assert.False(t, Payload{Action: state}.IsCreate())
		})
	}
}

func TestCommentURL(t *testing.T) {
	var repo = Repository{
		FullName: "caarlos0/gravekeeper",
	}
	var issue = Issue{
		Number: 1,
	}
	for _, payload := range []Payload{
		Payload{
			Action:      "pull",
			Repository:  repo,
			PullRequest: issue,
		},
		Payload{
			Action:     "issue",
			Repository: repo,
			Issue:      issue,
		},
	} {
		t.Run(payload.Action, func(t *testing.T) {
			var assert = assert.New(t)
			assert.Equal(
				"https://api.fakegithub.com/repos/caarlos0/gravekeeper/issues/1/comments",
				payload.commentURL("https://api.fakegithub.com"),
			)
		})
	}
}

func TestComment(t *testing.T) {
	var assert = assert.New(t)
	var token = "super-secret-token"
	var handler = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("token "+token, r.Header.Get("Authorization"))
		bts, err := ioutil.ReadAll(r.Body)
		assert.NoError(err)
		assert.Equal(`{"body": "`+msg+`"}`, string(bts))
		fmt.Fprintf(w, `{"ok":true}`)
	}
	var server = httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	var payload = Payload{
		Action: "opened",
		Repository: Repository{
			FullName: "caarlos0/gravekeeper",
		},
		PullRequest: Issue{
			Number: 10,
		},
	}
	assert.NoError(payload.Notify(server.URL, token))
}
