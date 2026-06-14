package http

import (
	"encoding/json"
	"net/http"

	"react-example/backend-golang/delivery/http/middleware"
	"react-example/backend-golang/domain"
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
	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}

func (h *NotificationHandler) CreateRule(w http.ResponseWriter, r *http.Request) error {
	var rule domain.NotificationRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}
	resp, err := h.usecase.CreateRule(r.Context(), rule)
	if err != nil {
		return err
	}
	middleware.SendJSON(w, http.StatusCreated, resp, nil)
	return nil
}

func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) error {
	resp, err := h.usecase.ListNotifications(r.Context(), "usr_current")
	if err != nil {
		return err
	}
	middleware.SendJSON(w, http.StatusOK, resp, map[string]interface{}{"total": len(resp)})
	return nil
}
