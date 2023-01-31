package connectors

import (
	"diploma/internal/common"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type github struct {
	token  string
	client *http.Client
}

func NewGithub(conf *common.Config) Github {
	return &github{
		token:  conf.GithubToken,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (g github) GetFile(repo, path string) ([]byte, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo, path)

	resp, err := g.doRequest(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get file: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	content := data["content"].(string)
	r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(content))

	return io.ReadAll(r)
}

func (g github) doRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.token))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	return g.client.Do(req)
}
