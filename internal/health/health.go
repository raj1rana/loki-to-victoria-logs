package health

import (
	"fmt"
	"net/http"
	"time"
)

type HealthChecker struct {
	httpClient *http.Client
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (h *HealthChecker) CheckLokiHealth(url string) error {
	resp, err := h.httpClient.Get(fmt.Sprintf("%s/ready", url))
	if err != nil {
		return fmt.Errorf("loki health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("loki returned non-200 status: %d", resp.StatusCode)
	}
	return nil
}

func (h *HealthChecker) CheckVictoriaHealth(url string) error {
	resp, err := h.httpClient.Get(fmt.Sprintf("%s/health", url))
	if err != nil {
		return fmt.Errorf("victoria health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("victoria returned non-200 status: %d", resp.StatusCode)
	}
	return nil
}
