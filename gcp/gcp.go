package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/maito1201/githubtrend"
	"github.com/maito1201/hatebu"
)

func PostNewsCommand(w http.ResponseWriter, _ *http.Request) {
	hatebuResults, err := hatebu.ScrapeHotEntry()
	if err != nil {
		http.Error(w, fmt.Sprintf("InternalError: %v", err), http.StatusBadRequest)
		return
	}

	contents := []string{"今日のはてブ1位", hatebuResults[0].Href, "", "GitHub Trending"}

	gitHubResults, err := githubtrend.ScrapeGitHubTrend()
	if err != nil {
		http.Error(w, fmt.Sprintf("InternalError: %v", err), http.StatusBadRequest)
		return
	}

	for i, v := range gitHubResults {
		if i > 2 {
			break
		}
		contents = append(contents, fmt.Sprintf("%dStar %s", v.Star, v.Href))
	}

	endpoint := os.Getenv("DISCORD_WEBHOOK")
	param := fmt.Sprintf(`{
		"username": "tech-news",
		"avatar_url": "https://b.st-hatena.com/54e3e2fdaf3b549836dedda8cb7409f9336287d9/images/v4/public/gh-logo@2x.png",
		"content": "%s"
	}`, strings.Join(contents, "\\r"))

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(param)))
	if err != nil {
		http.Error(w, fmt.Sprintf("InternalError: %v", err), http.StatusBadRequest)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	_, err = client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("InternalError: %v", err), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "Success")
}
