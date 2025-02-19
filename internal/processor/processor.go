package processor

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"log-pipeline/internal/models"
	"log-pipeline/internal/loki"
	"log-pipeline/internal/victoria"
)

type Processor struct {
	lokiClient     *loki.Client
	victoriaClient *victoria.Client
	seen           map[int64]bool
	mutex          sync.RWMutex
	stats          struct {
		processed int64
		errors    int64
		skipped   int64
	}
}

func NewProcessor(lokiClient *loki.Client, victoriaClient *victoria.Client) *Processor {
	return &Processor{
		lokiClient:     lokiClient,
		victoriaClient: victoriaClient,
		seen:           make(map[int64]bool),
	}
}

func (p *Processor) ProcessLogs(query string, startTime, endTime time.Time) error {
	logs, err := p.lokiClient.QueryLogs(query, startTime, endTime)
	if err != nil {
		p.stats.errors++
		return fmt.Errorf("failed to query Loki: %v", err)
	}

	var processingErrors []error
	for _, result := range logs.Data.Result {
		for _, value := range result.Values {
			if err := p.processLogEntry(value, result.Stream); err != nil {
				p.stats.errors++
				processingErrors = append(processingErrors, err)
			} else {
				p.stats.processed++
			}
		}
	}

	if len(processingErrors) > 0 {
		return fmt.Errorf("encountered %d errors while processing logs: %v", len(processingErrors), processingErrors)
	}

	return nil
}

func (p *Processor) processLogEntry(value []string, stream map[string]string) error {
	var logEntry models.LogEntry
	if err := json.Unmarshal([]byte(value[1]), &logEntry); err != nil {
		return fmt.Errorf("failed to unmarshal log entry: %v", err)
	}

	if p.isDuplicate(logEntry.Fields.EventRecordID) {
		p.stats.skipped++
		return nil
	}

	parsedData, err := p.parseLogData(logEntry.Fields.Data)
	if err != nil {
		return fmt.Errorf("failed to parse log data: %v", err)
	}

	victoriaData := map[string]interface{}{
		"event_record_id":   logEntry.Fields.EventRecordID,
		"timestamp":         logEntry.Timestamp,
		"computer":          logEntry.Tags["Computer"],
		"error_code":        parsedData.Error,
		"severity":          parsedData.Severity,
		"state":            parsedData.State,
		"start_time":        parsedData.StartTime,
		"trace_type":        parsedData.TraceType,
		"event_class_desc":  parsedData.EventClassDesc,
		"login_name":        parsedData.LoginName,
		"host_name":         parsedData.HostName,
		"text_data":         parsedData.TextData,
		"application_name":  parsedData.AppName,
		"database_name":     parsedData.DatabaseName,
		"object_name":       parsedData.ObjectName,
		"role_name":         parsedData.RoleName,
	}

	if err := p.victoriaClient.SendLog(victoriaData); err != nil {
		return fmt.Errorf("failed to send log to Victoria: %v", err)
	}

	return nil
}

func (p *Processor) isDuplicate(eventRecordID int64) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.seen[eventRecordID] {
		return true
	}
	p.seen[eventRecordID] = true
	return false
}

func (p *Processor) parseLogData(data string) (*models.ParsedData, error) {
	lines := strings.Split(data, "\n")
	parsed := &models.ParsedData{}

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Error":
			if errCode, err := strconv.Atoi(value); err == nil {
				parsed.Error = errCode
			}
		case "Severity":
			if severity, err := strconv.Atoi(value); err == nil {
				parsed.Severity = severity
			}
		case "State":
			if state, err := strconv.Atoi(value); err == nil {
				parsed.State = state
			}
		case "StartTime":
			parsed.StartTime = value
		case "TraceType":
			parsed.TraceType = value
		case "EventClassDesc":
			parsed.EventClassDesc = value
		case "LoginName":
			parsed.LoginName = value
		case "HostName":
			parsed.HostName = value
		case "TextData":
			parsed.TextData = value
		case "ApplicationName":
			parsed.AppName = value
		case "DatabaseName":
			parsed.DatabaseName = value
		case "ObjectName":
			parsed.ObjectName = value
		case "RoleName":
			parsed.RoleName = value
		}
	}

	return parsed, nil
}

// GetStats returns the current processing statistics
func (p *Processor) GetStats() (processed, errors, skipped int64) {
	return p.stats.processed, p.stats.errors, p.stats.skipped
}