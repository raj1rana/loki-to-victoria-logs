go mod download
```
3. Build the binary:
```bash
go build -o log-pipeline
```

### Using Docker

1. Build the container:
```bash
docker build -t log-pipeline:latest .
```

2. Run the container:
```bash
docker run -d \
  --name log-pipeline \
  -e LOKI_URL="http://loki:3100" \
  -e VICTORIA_URL="http://victoria:8428" \
  -v $(pwd)/config.json:/app/config.json \
  log-pipeline:latest
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
./log-pipeline -config /path/to/config.json
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