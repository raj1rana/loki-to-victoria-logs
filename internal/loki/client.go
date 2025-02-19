package loki

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v4"
	"log-pipeline/internal/resilience"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	maxRetries int
	cb         *resilience.CircuitBreaker
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
		cb:         resilience.NewCircuitBreaker("loki-client"),
	}
}

func (c *Client) QueryLogs(query string, start, end time.Time) (*LogResponse, error) {
	operation := func() (*LogResponse, error) {
		params := url.Values{}
		params.Add("query", query)
		params.Add("start", fmt.Sprintf("%d", start.UnixNano()))
		params.Add("end", fmt.Sprintf("%d", end.UnixNano()))

		url := fmt.Sprintf("%s/loki/api/v1/query_range?%s", c.baseURL, params.Encode())

		// Execute request through circuit breaker
		resp, err := c.cb.Execute(func() (interface{}, error) {
			resp, err := c.httpClient.Get(url)
			if err != nil {
				return nil, fmt.Errorf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
			}

			var logResp LogResponse
			if err := json.NewDecoder(resp.Body).Decode(&logResp); err != nil {
				return nil, fmt.Errorf("failed to decode response: %v", err)
			}

			return &logResp, nil
		})

		if err != nil {
			return nil, err
		}

		return resp.(*LogResponse), nil
	}

	// Create exponential backoff
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute

	var result *LogResponse
	err := backoff.RetryNotify(func() error {
		var err error
		result, err = operation()
		return err
	}, b, func(err error, duration time.Duration) {
		log.Printf("Retrying Loki query after %v due to error: %v", duration, err)
	})

	if err != nil {
		return nil, fmt.Errorf("all retries failed: %v", err)
	}

	return result, nil
}