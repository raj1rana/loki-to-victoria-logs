# Log Pipeline Utility

A robust Go-based log pipeline utility designed for collecting, processing, and forwarding logs from Loki to Victoria Logs. This utility provides efficient log processing with features like deduplication, schema management, and resilient error handling.

## Features

- **Log Collection**: Collect logs from Loki using customizable queries
- **Schema Management**: Dynamic schema support for Victoria Logs integration
- **Deduplication**: Prevent duplicate log entries based on EventRecordID
- **Resilient Processing**:
  - Circuit breaker pattern to prevent cascading failures
  - Exponential backoff retry mechanism
  - Configurable retry attempts and timeouts
- **Health Monitoring**: Built-in health checks for both Loki and Victoria services
- **Structured Logging**: Comprehensive parsing and transformation of log data

## Prerequisites

- Go 1.19 or later
- Access to Loki and Victoria Logs instances

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

## Configuration

The utility can be configured through both a configuration file (`config.json`) and environment variables.

### Configuration File (config.json)
```json
{
    "loki": {
        "url": "http://localhost:3100",
        "query": "{topic=\"iaas-database-auditlogs\"}",
        "interval": "1m"
    },
    "victoria": {
        "url": "http://localhost:8428",
        "schema": "database_audit_logs"
    },
    "batchSize": 1000,
    "timeWindow": "5m"
}
```

### Environment Variables

- `LOKI_URL`: Loki server URL (overrides config file)
- `VICTORIA_URL`: Victoria Logs server URL (overrides config file)

## Usage

Run the utility with a custom configuration file:
```bash
go run main.go -config /path/to/config.json
```

## Schema

The utility uses a predefined schema for Victoria Logs:

```json
{
    "name": "database_audit_logs",
    "version": 1,
    "fields": [
        {"name": "event_record_id", "type": "int64", "required": true},
        {"name": "timestamp", "type": "timestamp", "required": true},
        {"name": "computer", "type": "string", "required": true},
        // ... other fields
    ]
}
```

## Architecture

### Components

1. **Processor**: Core component handling log processing and deduplication
2. **Loki Client**: Handles log collection with resilient error handling
3. **Victoria Client**: Manages log forwarding with circuit breaker pattern
4. **Health Checker**: Monitors service health and availability

### Resilience Features

#### Circuit Breaker
- Prevents cascading failures
- Configurable thresholds and timeouts
- State change monitoring

#### Retry Mechanism
- Exponential backoff
- Configurable max retries
- Proper error logging

## Monitoring

The utility provides real-time statistics:
- Number of processed logs
- Error counts
- Skipped (duplicate) entries
- Circuit breaker state changes

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]
