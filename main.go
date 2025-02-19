package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"log-pipeline/config"
	"log-pipeline/internal/health"
	"log-pipeline/internal/loki"
	"log-pipeline/internal/processor"
	"log-pipeline/internal/victoria"
	"log-pipeline/pkg/utils"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Check service health
	healthChecker := health.NewHealthChecker()
	if err := healthChecker.CheckLokiHealth(cfg.Loki.URL); err != nil {
		log.Fatalf("Loki health check failed: %v", err)
	}
	if err := healthChecker.CheckVictoriaHealth(cfg.Victoria.URL); err != nil {
		log.Fatalf("Victoria health check failed: %v", err)
	}

	log.Printf("Services health check passed")

	// Initialize clients
	lokiClient := loki.NewClient(cfg.Loki.URL)
	victoriaClient := victoria.NewClient(cfg.Victoria.URL, cfg.Victoria.Schema)

	// Initialize processor
	proc := processor.NewProcessor(lokiClient, victoriaClient)

	// Start health check server
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			if err := healthChecker.CheckLokiHealth(cfg.Loki.URL); err != nil {
				http.Error(w, "Loki health check failed", http.StatusServiceUnavailable)
				return
			}
			if err := healthChecker.CheckVictoriaHealth(cfg.Victoria.URL); err != nil {
				http.Error(w, "Victoria health check failed", http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Health check server error: %v", err)
		}
	}()

	log.Printf("Starting log pipeline with query: %s", cfg.Loki.Query)
	log.Printf("Time window: %v, Interval: %v", time.Duration(cfg.TimeWindow), time.Duration(cfg.Loki.Interval))

	// Stats reporting ticker
	statsTicker := time.NewTicker(1 * time.Minute)
	defer statsTicker.Stop()

	// Main processing loop
	for {
		select {
		case <-statsTicker.C:
			processed, errors, skipped := proc.GetStats()
			log.Printf("Stats - Processed: %d, Errors: %d, Skipped: %d", processed, errors, skipped)
		default:
			start, end := utils.GetTimeRange(time.Duration(cfg.TimeWindow))
			log.Printf("Processing logs from %v to %v", start, end)

			if err := proc.ProcessLogs(cfg.Loki.Query, start, end); err != nil {
				log.Printf("Error processing logs: %v", err)
			}

			time.Sleep(time.Duration(cfg.Loki.Interval))
		}
	}
}