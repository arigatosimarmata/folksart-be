package handlers

import (
	"github.com/gofiber/fiber/v2"
	"react-example/backend-golang/httputil"
	"react-example/backend-golang/internal/domain"
	"react-example/backend-golang/internal/dto"
)

type NotificationHandler struct {
	usecase domain.NotificationUsecase
}

func NewNotificationHandler(u domain.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{usecase: u}
}

func (h *NotificationHandler) ListRules(c *fiber.Ctx) error {
	resp, err := h.usecase.ListRules(c.Context())
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	dtos := make([]dto.NotificationRuleResponse, 0)
	for _, rule := range resp {
		dtos = append(dtos, dto.NotificationRuleResponse{
			ID:        rule.ID,
			Name:      rule.Name,
			Trigger:   rule.Trigger,
			Severity:  rule.Severity,
			Channels:  rule.Channels,
			Active:    rule.Active,
			CreatedAt: rule.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return httputil.WriteSuccessResponse(c, "Success", dtos, nil)
}

func (h *NotificationHandler) CreateRule(c *fiber.Ctx) error {
	var req dto.CreateNotificationRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	resp, err := h.usecase.CreateRule(c.Context(), domain.NotificationRule{
		Name:     req.Name,
		Trigger:  req.Trigger,
		Severity: req.Severity,
		Channels: req.Channels,
		Active:   req.Active,
	})
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	res := dto.NotificationRuleResponse{
		ID:        resp.ID,
		Name:      resp.Name,
		Trigger:   resp.Trigger,
		Severity:  resp.Severity,
		Channels:  resp.Channels,
		Active:    resp.Active,
		CreatedAt: resp.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return httputil.WriteSuccessResponse(c, "Notification rule created successfully", res, nil)
}

func (h *NotificationHandler) ListNotifications(c *fiber.Ctx) error {
	resp, err := h.usecase.ListNotifications(c.Context(), "usr_current")
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	dtos := make([]dto.NotificationResponse, 0)
	for _, n := range resp {
		dtos = append(dtos, dto.NotificationResponse{
			ID:        n.ID,
			Type:      n.Type,
			Title:     n.Title,
			Body:      n.Body,
			Severity:  n.Severity,
			Read:      n.Read,
			CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return httputil.WriteSuccessResponse(c, "Success", dtos, map[string]interface{}{"total": len(dtos)})
}
