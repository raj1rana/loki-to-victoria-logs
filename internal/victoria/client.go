package victoria

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"log-pipeline/internal/resilience"
)

type Client struct {
	baseURL    string
	schema     string
	httpClient *http.Client
	maxRetries int
	cb         *resilience.CircuitBreaker
}

func NewClient(baseURL, schema string) *Client {
	return &Client{
		baseURL: baseURL,
		schema:  schema,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxRetries: 3,
		cb:         resilience.NewCircuitBreaker("victoria-client"),
	}
}

func (c *Client) SendLog(data map[string]interface{}) error {
	operation := func() error {
		// Add schema information
		data["_schema"] = c.schema

		payload, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %v", err)
		}

		url := fmt.Sprintf("%s/write", c.baseURL)

		// Execute request through circuit breaker
		_, err = c.cb.Execute(func() (interface{}, error) {
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := c.httpClient.Do(req)
			if err != nil {
				return nil, fmt.Errorf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
			}

			return nil, nil
		})

		return err
	}

	// Create exponential backoff
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute

	err := backoff.RetryNotify(func() error {
		return operation()
	}, b, func(err error, duration time.Duration) {
		log.Printf("Retrying Victoria log send after %v due to error: %v", duration, err)
	})

	if err != nil {
		return fmt.Errorf("all retries failed: %v", err)
	}

	return nil
}