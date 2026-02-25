package utils

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type AuditLog struct {
    ID             string                 `json:"id"`
    Timestamp      string                 `json:"timestamp"`
    Action         string                 `json:"action"`
    Actor          string                 `json:"actor"`
    ActorRole      string                 `json:"actor_role"`
    TargetUser     string                 `json:"target_user"`
    TargetDate     string                 `json:"target_date"`
    ChangeDetails  map[string]interface{} `json:"change_details"`
    Status         string                 `json:"status"`
    IPAddress      string                 `json:"ip_address"`
    ErrorMessage   string                 `json:"error_message,omitempty"`
}

const auditLogFile = "audit_log.json"

// LoadAuditLogs reads all audit logs
func LoadAuditLogs() ([]AuditLog, error) {
    data, err := os.ReadFile(auditLogFile)
    if err != nil {
        if os.IsNotExist(err) {
            return []AuditLog{}, nil
        }
        return nil, err
    }

    var logs []AuditLog
    if err := json.Unmarshal(data, &logs); err != nil {
        return nil, err
    }
    return logs, nil
}

// SaveAuditLogs writes audit logs to file
func SaveAuditLogs(logs []AuditLog) error {
    data, err := json.MarshalIndent(logs, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(auditLogFile, data, 0644)
}

// LogAction creates and saves an audit log entry
func LogAction(action, actor, actorRole, targetUser, targetDate, ipAddress string, changeDetails map[string]interface{}, status string) error {
    logs, err := LoadAuditLogs()
    if err != nil {
        return err
    }

    auditID := fmt.Sprintf("audit_%d", len(logs)+1)

    log := AuditLog{
        ID:            auditID,
        Timestamp:     time.Now().UTC().Format(time.RFC3339),
        Action:        action,
        Actor:         actor,
        ActorRole:     actorRole,
        TargetUser:    targetUser,
        TargetDate:    targetDate,
        ChangeDetails: changeDetails,
        Status:        status,
        IPAddress:     ipAddress,
    }

    logs = append(logs, log)
    return SaveAuditLogs(logs)
}

// LogActionWithError logs failed operations
func LogActionWithError(action, actor, actorRole, targetUser, targetDate, ipAddress string, errorMsg string) error {
    logs, err := LoadAuditLogs()
    if err != nil {
        return err
    }

    auditID := fmt.Sprintf("audit_%d", len(logs)+1)

    log := AuditLog{
        ID:           auditID,
        Timestamp:    time.Now().UTC().Format(time.RFC3339),
        Action:       action,
        Actor:        actor,
        ActorRole:    actorRole,
        TargetUser:   targetUser,
        TargetDate:   targetDate,
        Status:       "failed",
        ErrorMessage: errorMsg,
        IPAddress:    ipAddress,
    }

    logs = append(logs, log)
    return SaveAuditLogs(logs)
}