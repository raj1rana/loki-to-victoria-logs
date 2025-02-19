package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Duration is a wrapper for time.Duration that implements json.Unmarshaler
type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return fmt.Errorf("invalid duration: %v", v)
	}
}

type Config struct {
	Loki struct {
		URL      string   `json:"url"`
		Query    string   `json:"query"`
		Interval Duration `json:"interval"`
	} `json:"loki"`
	Victoria struct {
		URL    string `json:"url"`
		Schema string `json:"schema"`
	} `json:"victoria"`
	BatchSize  int      `json:"batchSize"`
	TimeWindow Duration `json:"timeWindow"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	// Override with environment variables if present
	if url := os.Getenv("LOKI_URL"); url != "" {
		config.Loki.URL = url
	}
	if url := os.Getenv("VICTORIA_URL"); url != "" {
		config.Victoria.URL = url
	}

	// Validate required fields
	if config.Loki.URL == "" {
		return nil, fmt.Errorf("loki URL is required")
	}
	if config.Victoria.URL == "" {
		return nil, fmt.Errorf("victoria URL is required")
	}
	if config.Loki.Query == "" {
		return nil, fmt.Errorf("loki query is required")
	}
	if config.Victoria.Schema == "" {
		return nil, fmt.Errorf("victoria schema is required")
	}

	return &config, nil
}