package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"react-example/backend-golang/delivery/http/middleware"
	"react-example/backend-golang/domain"
)

type AccessRequestHandler struct {
	usecase domain.AccessRequestUsecase
}

func NewAccessRequestHandler(u domain.AccessRequestUsecase) *AccessRequestHandler {
	return &AccessRequestHandler{usecase: u}
}

func (h *AccessRequestHandler) Submit(w http.ResponseWriter, r *http.Request) error {
	var req domain.AccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}

	resp, err := h.usecase.SubmitRequest(r.Context(), req)
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusCreated, resp, nil)
	return nil
}

func (h *AccessRequestHandler) List(w http.ResponseWriter, r *http.Request) error {
	status := r.URL.Query().Get("status")
	resp, err := h.usecase.ListRequests(r.Context(), status)
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, resp, map[string]interface{}{"total": len(resp)})
	return nil
}

func (h *AccessRequestHandler) Approve(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		return middleware.NewCustomError(http.StatusBadRequest, "Missing request ID", nil)
	}
	id := pathParts[4]

	var payload struct {
		Operator string `json:"operator"`
		Note     string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}

	resp, err := h.usecase.ApproveRequest(r.Context(), id, payload.Operator, payload.Note)
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}

func (h *AccessRequestHandler) Reject(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		return middleware.NewCustomError(http.StatusBadRequest, "Missing request ID", nil)
	}
	id := pathParts[4]

	var payload struct {
		Operator string `json:"operator"`
		Reason   string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}

	resp, err := h.usecase.RejectRequest(r.Context(), id, payload.Operator, payload.Reason)
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}
