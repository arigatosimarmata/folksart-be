package dto

type AuditLogResponse struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Actor     string `json:"actor"`
	Action    string `json:"action"`
	Target    string `json:"target"`
	Severity  string `json:"severity"`
}

type CreateAuditLogRequest struct {
	Actor    string `json:"actor" validate:"required"`
	Action   string `json:"action" validate:"required"`
	Target   string `json:"target"`
	Severity string `json:"severity" validate:"required"`
}
