package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/gravekeeper/github"
)

func shouldHandle(req *http.Request) bool {
	var event = req.Header.Get("X-GitHub-Event")
	return event == "issues" || event == "pull_request"
}

func handler(api, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !shouldHandle(r) {
			fmt.Fprintln(w, "not a valid event, ignoring...")
			return
		}
		var evt github.Payload
		if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("received event:", evt)
		if !evt.IsCreate() {
			fmt.Fprintln(w, "not a valid event, ignoring...")
			return
		}
		if err := evt.Notify(api, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "done")
	}
}

func main() {
	var addr = ":" + os.Getenv("PORT")
	var token = os.Getenv("GITHUB_TOKEN")
	var api = os.Getenv("GITHUB_API")
	http.Handle("/", handler(api, token))
	log.Println("listening at", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
