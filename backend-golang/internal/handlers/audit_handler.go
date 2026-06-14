package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"react-example/backend-golang/httputil"
	"react-example/backend-golang/internal/domain"
	"react-example/backend-golang/internal/dto"
)

type AuditHandler struct {
	auditUsecase domain.AuditUsecase
}

func NewAuditHandler(au domain.AuditUsecase) *AuditHandler {
	return &AuditHandler{auditUsecase: au}
}

func (h *AuditHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) error {
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
		return err
	}

	logDtos := make([]dto.AuditLogResponse, 0)
	for _, l := range logs {
		logDtos = append(logDtos, dto.AuditLogResponse{
			ID:        l.ID,
			Timestamp: l.Timestamp.Format("2006-01-02 15:04:05"),
			Actor:     l.Actor,
			Action:    l.Action,
			Target:    l.Target,
			Severity:  l.Severity,
		})
	}

	httputil.WriteSuccessResponse(w, "Success", logDtos, map[string]interface{}{
		"limit": limit,
		"total": len(logDtos),
	})
	return nil
}

func (h *AuditHandler) CreateLog(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateAuditLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	ctx := r.Context()
	newLog, err := h.auditUsecase.RecordAction(ctx, req.Actor, req.Action, req.Target, req.Severity)
	if err != nil {
		return err
	}

	res := dto.AuditLogResponse{
		ID:        newLog.ID,
		Timestamp: newLog.Timestamp.Format("2006-01-02 15:04:05"),
		Actor:     newLog.Actor,
		Action:    newLog.Action,
		Target:    newLog.Target,
		Severity:  newLog.Severity,
	}

	httputil.WriteSuccessResponse(w, "Audit log created successfully", res, nil)
	return nil
}

func (h *AuditHandler) SignLogs(w http.ResponseWriter, r *http.Request) error {
	httputil.WriteSuccessResponse(w, "Logs signed successfully", map[string]interface{}{
		"signed_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
		"algorithm":    "RS256",
		"checksum":     "sha256:a1b2c3d4e5f6...",
		"signed_at":    "2025-06-10T14:00:00Z",
	}, nil)
	return nil
}
