package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"react-example/backend-golang/delivery/http/middleware"
	"react-example/backend-golang/domain"
)

type AuditHandler struct {
	auditUsecase domain.AuditUsecase
}

func NewAuditHandler(au domain.AuditUsecase) *AuditHandler {
	return &AuditHandler{auditUsecase: au}
}

// ListAuditLogs handles GET /api/v1/audit-logs
// @Summary Fetch administrative audit trails
// @Description Retrieve a ledger of all security and governance events within the IAM system.
// @Tags audit
// @Accept json
// @Produce json
// @Param severity query string false "Filter by severity (e.g., Critical, Warning, Info)"
// @Param limit query int false "Maximum number of logs to return" default(50)
// @Success 200 {array} domain.AuditLog "Array of audit log entries"
// @Failure 405 {object} middleware.APIError "Method Not Allowed"
// @Failure 500 {object} middleware.APIError "Internal Server Error"
// @Router /api/v1/audit-logs [get]
func (h *AuditHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return middleware.NewCustomError(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
	}

	ctx := r.Context()
	severity := r.URL.Query().Get("severity")
	limitStr := r.URL.Query().Get("limit")

	limit := 50
	if limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	logs, err := h.auditUsecase.GetAuditTrail(ctx, severity, limit)
	if err != nil {
		return middleware.NewCustomError(http.StatusInternalServerError, "Failed to capture logs ledger index", err)
	}

	middleware.SendJSON(w, http.StatusOK, logs, map[string]interface{}{
		"limit": limit,
		"total": len(logs),
	})
	return nil
}

type CreateLogPayload struct {
	Actor    string `json:"actor"`
	Action   string `json:"action"`
	Target   string `json:"target"`
	Severity string `json:"severity"`
}

// CreateLog handles POST /api/v1/audit-logs
// @Summary Record a manual audit event
// @Description Manually append a governance event to the audit trail.
// @Tags audit
// @Accept json
// @Produce json
// @Param log body CreateLogPayload true "Audit log payload"
// @Success 201 {object} domain.AuditLog "Newly created audit log record"
// @Failure 400 {object} middleware.APIError "Bad Request / Validation Error"
// @Failure 405 {object} middleware.APIError "Method Not Allowed"
// @Failure 500 {object} middleware.APIError "Internal Server Error"
// @Router /api/v1/audit-logs [post]
func (h *AuditHandler) CreateLog(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return middleware.NewCustomError(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
	}

	var req CreateLogPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid JSON payload structure", err)
	}

	ctx := r.Context()
	newLog, err := h.auditUsecase.RecordAction(ctx, req.Actor, req.Action, req.Target, req.Severity)
	if err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Failed to write administrative event record", err)
	}

	middleware.SendJSON(w, http.StatusCreated, newLog, nil)
	return nil
}

func (h *AuditHandler) SignLogs(w http.ResponseWriter, r *http.Request) error {
	middleware.SendJSON(w, http.StatusOK, map[string]interface{}{
		"signed_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
		"algorithm":    "RS256",
		"checksum":     "sha256:a1b2c3d4e5f6...",
		"signed_at":    "2025-06-10T14:00:00Z",
	}, nil)
	return nil
}
