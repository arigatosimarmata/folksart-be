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

// ListAuditLogs handles GET /api/v1/audit
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

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(logs)
}

type CreateLogPayload struct {
	Actor    string `json:"actor"`
	Action   string `json:"action"`
	Target   string `json:"target"`
	Severity string `json:"severity"`
}

// CreateLog handles POST /api/v1/audit
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(newLog)
}
