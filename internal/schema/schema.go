package schema

const Schema = `
{
    "name": "database_audit_logs",
    "version": 1,
    "fields": [
        {"name": "event_record_id", "type": "int64", "required": true},
        {"name": "timestamp", "type": "timestamp", "required": true},
        {"name": "computer", "type": "string", "required": true},
        {"name": "error_code", "type": "int"},
        {"name": "severity", "type": "int"},
        {"name": "state", "type": "int"},
        {"name": "start_time", "type": "timestamp"},
        {"name": "trace_type", "type": "string"},
        {"name": "event_class_desc", "type": "string"},
        {"name": "login_name", "type": "string"},
        {"name": "host_name", "type": "string"},
        {"name": "text_data", "type": "string"},
        {"name": "application_name", "type": "string"},
        {"name": "database_name", "type": "string"},
        {"name": "object_name", "type": "string"},
        {"name": "role_name", "type": "string"}
    ]
}
`
