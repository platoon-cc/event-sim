package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/platoon-cc/platoon-cli/internal/model"
	"github.com/platoon-cc/platoon-cli/internal/settings"
)

var CacheDuration int64 = 30 * 60

type Client struct {
	url   string
	token string
}

func New(token string) *Client {
	c := &Client{
		token: token,
		url:   "http://pl.localhost:9999/cli",
	}

	return c
}

func (c *Client) GetTeamList() ([]model.Team, error) {
	teams, err := settings.GetCache[[]model.Team]("team")
	if errors.Is(err, settings.ErrExpired) {
		resp, _, err := c.serverGet("team")
		if err != nil {
			return nil, err
		}
		json.Unmarshal(resp, &teams)
		settings.SetCache("team", teams)
	} else if err != nil {
		return teams, err
	}
	return teams, nil
}

func (c *Client) GetProjectList() ([]model.Project, error) {
	projects, err := settings.GetCache[[]model.Project]("project")
	if errors.Is(err, settings.ErrExpired) {
		teamId := settings.GetActive("team")
		resp, _, err := c.serverGet("team/" + teamId + "/project")
		if err != nil {
			return nil, err
		}
		json.Unmarshal(resp, &projects)
		settings.SetCache("project", projects)
	} else if err != nil {
		return projects, err
	}
	return projects, nil
}

func (c *Client) makeUrl(endpoint string) string {
	u, _ := url.Parse(c.url)
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
}

func (c *Client) serverGet(endpoint string) ([]byte, int, error) {
	return c.call("GET", endpoint, nil)
}

func (c *Client) serverPost(endpoint string, data any) ([]byte, int, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}
	return c.call("POST", endpoint, bytes.NewBuffer(payload))
}

func (c *Client) call(verb string, endpoint string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest(verb, c.makeUrl(endpoint), body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	// var response map[string]any
	// if err := json.Unmarshal(body, &response); err != nil {
	// 	return nil, 0, err
	// }

	return resp, res.StatusCode, nil
}
