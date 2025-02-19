package victoria

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	schema     string
	httpClient *http.Client
	maxRetries int
}

func NewClient(baseURL, schema string) *Client {
	return &Client{
		baseURL:    baseURL,
		schema:     schema,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxRetries: 3,
	}
}

func (c *Client) SendLog(data map[string]interface{}) error {
	var lastErr error
	for retry := 0; retry <= c.maxRetries; retry++ {
		if retry > 0 {
			// Exponential backoff: 1s, 2s, 4s
			backoffDuration := time.Duration(1<<uint(retry-1)) * time.Second
			time.Sleep(backoffDuration)
		}

		// Add schema information
		data["_schema"] = c.schema

		payload, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %v", err)
		}

		url := fmt.Sprintf("%s/write", c.baseURL)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: %v", retry+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("attempt %d: server returned status %d", retry+1, resp.StatusCode)
			continue
		}

		return nil
	}

	return fmt.Errorf("failed after %d retries, last error: %v", c.maxRetries, lastErr)
}