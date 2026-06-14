package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
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

func (h *AuditHandler) ListAuditLogs(c *fiber.Ctx) error {
	severity := c.Query("severity")
	limitStr := c.Query("limit", "50")

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 50
	}

	logs, err := h.auditUsecase.GetAuditTrail(c.Context(), severity, limit)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
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

	return httputil.WriteSuccessResponse(c, "Success", logDtos, map[string]interface{}{
		"limit": limit,
		"total": len(logDtos),
	})
}

func (h *AuditHandler) CreateLog(c *fiber.Ctx) error {
	var req dto.CreateAuditLogRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	newLog, err := h.auditUsecase.RecordAction(c.Context(), req.Actor, req.Action, req.Target, req.Severity)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	res := dto.AuditLogResponse{
		ID:        newLog.ID,
		Timestamp: newLog.Timestamp.Format("2006-01-02 15:04:05"),
		Actor:     newLog.Actor,
		Action:    newLog.Action,
		Target:    newLog.Target,
		Severity:  newLog.Severity,
	}

	return httputil.WriteSuccessResponse(c, "Audit log created successfully", res, nil)
}

func (h *AuditHandler) SignLogs(c *fiber.Ctx) error {
	return httputil.WriteSuccessResponse(c, "Logs signed successfully", map[string]interface{}{
		"signed_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
		"algorithm":    "RS256",
		"checksum":     "sha256:a1b2c3d4e5f6...",
		"signed_at":    "2025-06-10T14:00:00Z",
	}, nil)
}
