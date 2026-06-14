package handlers

import (
	"encoding/json"
	"net/http"

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

func (h *NotificationHandler) ListRules(w http.ResponseWriter, r *http.Request) error {
	resp, err := h.usecase.ListRules(r.Context())
	if err != nil {
		return err
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

	httputil.WriteSuccessResponse(w, "Success", dtos, nil)
	return nil
}

func (h *NotificationHandler) CreateRule(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateNotificationRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	resp, err := h.usecase.CreateRule(r.Context(), domain.NotificationRule{
		Name:     req.Name,
		Trigger:  req.Trigger,
		Severity: req.Severity,
		Channels: req.Channels,
		Active:   req.Active,
	})
	if err != nil {
		return err
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

	httputil.WriteSuccessResponse(w, "Notification rule created successfully", res, nil)
	return nil
}

func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) error {
	resp, err := h.usecase.ListNotifications(r.Context(), "usr_current")
	if err != nil {
		return err
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

	httputil.WriteSuccessResponse(w, "Success", dtos, map[string]interface{}{"total": len(dtos)})
	return nil
}
