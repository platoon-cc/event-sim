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
	url          string
	privateToken string
	publicToken  string
}

func New() (*Client, error) {
	c := &Client{}
	var err error

	c.url, err = settings.GetAuth("server")
	if err != nil {
		return nil, err
	}

	c.privateToken, err = settings.GetAuth("token")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) GetTeamList() ([]model.Team, error) {
	teams, err := settings.GetCache[[]model.Team]("team")
	if errors.Is(err, settings.ErrExpired) {
		resp, _, err := c.serverGet("team", "")
		if err != nil {
			return nil, err
		}
		json.Unmarshal(resp, &teams)
		settings.SetCache("team", teams)
	} else {
		return teams, err
	}
	return teams, nil
}

func (c *Client) GetProjectList() ([]model.Project, error) {
	projects, err := settings.GetCache[[]model.Project]("project")
	if errors.Is(err, settings.ErrExpired) {
		teamId, err := settings.GetActive("team")
		if err != nil {
			return nil, err
		}
		resp, _, err := c.serverGet("team/"+teamId+"/project", "")
		if err != nil {
			return nil, err
		}
		json.Unmarshal(resp, &projects)
		settings.SetCache("project", projects)
	} else {
		return projects, err
	}
	return projects, nil
}

func (c *Client) GetEvents(projectId string, lastId int64) ([]model.Event, error) {
	resp, _, err := c.serverGet(fmt.Sprintf("project/%s/events", projectId), fmt.Sprintf("from=%d", lastId))
	// resp, _, err := c.serverGet(fmt.Sprintf("project/%s/events", projectId))
	if err != nil {
		return nil, err
	}

	events := []model.Event{}
	if err := json.Unmarshal(resp, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (c *Client) GetAccessToken() (string, error) {
	projectId, err := settings.GetActive("project")
	if err != nil {
		return "", err
	}
	cacheKey := "accessToken." + projectId
	accessToken, err := settings.GetCache[string](cacheKey)

	if errors.Is(err, settings.ErrExpired) {
		resp, _, err := c.serverGet(fmt.Sprintf("project/%s/accessToken", projectId), "")
		if err != nil {
			return accessToken, err
		}
		json.Unmarshal(resp, &accessToken)
		settings.SetCache(cacheKey, accessToken)
	} else {
		return accessToken, err
	}
	return accessToken, nil
}

func (c *Client) PostSimEvents(events []model.Event) error {
	token, err := c.GetAccessToken()
	if err != nil {
		return err
	}
	c.publicToken = token
	_, _, err = c.serverPost("/api/ingest", events, "")
	return err
}

func (c *Client) makeUrl(endpoint string) string {
	u, _ := url.Parse(c.url)
	if c.publicToken == "" {
		u.Path = path.Join(u.Path, "cli")
	}
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
}

func (c *Client) serverGet(endpoint string, query string) ([]byte, int, error) {
	return c.call("GET", endpoint, nil, query)
}

func (c *Client) serverPost(endpoint string, data any, query string) ([]byte, int, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}
	return c.call("POST", endpoint, bytes.NewBuffer(payload), query)
}

func (c *Client) call(verb string, endpoint string, body io.Reader, query string) ([]byte, int, error) {
	req, err := http.NewRequest(verb, c.makeUrl(endpoint), body)
	if err != nil {
		return nil, 0, err
	}

	req.URL.RawQuery = query

	if c.publicToken != "" {
		req.Header.Add("X-API-KEY", c.publicToken)
	} else {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.privateToken))
	}

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

	if res.StatusCode == 401 {
		return nil, res.StatusCode, fmt.Errorf("auth error. Please log in: `platoon-cli auth login`")
	}
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return resp, res.StatusCode, nil
	}
	return nil, res.StatusCode, fmt.Errorf("error: %d", res.StatusCode)
}
