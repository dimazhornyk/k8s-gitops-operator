package connectors

import (
	"diploma/internal/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type github struct {
	client   *http.Client
	token    string
	username string
}

func NewGithub(token, username string) *github {
	return &github{
		client:   &http.Client{Timeout: 10 * time.Second},
		token:    token,
		username: username,
	}
}

func (g github) GetRepositoryEvents(repo string) ([]common.RepositoryEvent, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/events", g.username, repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.token))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	res, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	events := make([]common.RepositoryEvent, 0)
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, err
	}

	return events, nil
}
