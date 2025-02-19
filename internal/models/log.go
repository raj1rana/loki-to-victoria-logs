package models

type LogEntry struct {
    Fields struct {
        Data           string `json:"Data"`
        EventRecordID int64  `json:"EventRecordID"`
        Message       string `json:"Message"`
        ProcessName   string `json:"ProcessName"`
        UserID        string `json:"UserID"`
        Version       int    `json:"Version"`
    } `json:"fields"`
    Name      string            `json:"name"`
    Tags      map[string]string `json:"tags"`
    Timestamp int64            `json:"timestamp"`
}

type ParsedData struct {
    Error         int    `json:"Error"`
    Severity      int    `json:"Severity"`
    State         int    `json:"State"`
    StartTime     string `json:"StartTime"`
    TraceType     string `json:"TraceType"`
    EventClassDesc string `json:"EventClassDesc"`
    LoginName     string `json:"LoginName"`
    HostName      string `json:"HostName"`
    TextData      string `json:"TextData"`
    AppName       string `json:"ApplicationName"`
    DatabaseName  string `json:"DatabaseName"`
    ObjectName    string `json:"ObjectName"`
    RoleName      string `json:"RoleName"`
}
