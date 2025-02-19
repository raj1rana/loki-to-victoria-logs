package loki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	maxRetries int
}

type LogResponse struct {
	Data struct {
		Result []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string       `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxRetries: 3,
	}
}

func (c *Client) QueryLogs(query string, start, end time.Time) (*LogResponse, error) {
	var lastErr error
	for retry := 0; retry <= c.maxRetries; retry++ {
		if retry > 0 {
			// Exponential backoff: 1s, 2s, 4s
			backoffDuration := time.Duration(1<<uint(retry-1)) * time.Second
			time.Sleep(backoffDuration)
		}

		params := url.Values{}
		params.Add("query", query)
		params.Add("start", fmt.Sprintf("%d", start.UnixNano()))
		params.Add("end", fmt.Sprintf("%d", end.UnixNano()))

		url := fmt.Sprintf("%s/loki/api/v1/query_range?%s", c.baseURL, params.Encode())

		resp, err := c.httpClient.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: %v", retry+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("attempt %d: server returned status %d", retry+1, resp.StatusCode)
			continue
		}

		var logResp LogResponse
		if err := json.NewDecoder(resp.Body).Decode(&logResp); err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to decode response: %v", retry+1, err)
			continue
		}

		return &logResp, nil
	}

	return nil, fmt.Errorf("failed after %d retries, last error: %v", c.maxRetries, lastErr)
}