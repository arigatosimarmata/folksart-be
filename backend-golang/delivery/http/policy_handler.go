package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"react-example/backend-golang/delivery/http/middleware"
	"react-example/backend-golang/domain"
)

type PolicyHandler struct {
	usecase domain.PolicyUsecase
}

func NewPolicyHandler(u domain.PolicyUsecase) *PolicyHandler {
	return &PolicyHandler{usecase: u}
}

func (h *PolicyHandler) List(w http.ResponseWriter, r *http.Request) error {
	resp, err := h.usecase.ListPolicies(r.Context())
	if err != nil {
		return err
	}
	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}

func (h *PolicyHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var p domain.Policy
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}
	resp, err := h.usecase.CreatePolicy(r.Context(), p)
	if err != nil {
		return err
	}
	middleware.SendJSON(w, http.StatusCreated, resp, nil)
	return nil
}

func (h *PolicyHandler) Update(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		return middleware.NewCustomError(http.StatusBadRequest, "Missing ID", nil)
	}
	id := pathParts[4]

	var p domain.Policy
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}
	resp, err := h.usecase.UpdatePolicy(r.Context(), id, p)
	if err != nil {
		return err
	}
	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}

func (h *PolicyHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		return middleware.NewCustomError(http.StatusBadRequest, "Missing ID", nil)
	}
	id := pathParts[4]
	err := h.usecase.DeletePolicy(r.Context(), id)
	if err != nil {
		return err
	}
	middleware.SendJSON(w, http.StatusOK, map[string]string{"message": "Policy deleted"}, nil)
	return nil
}

func (h *PolicyHandler) Evaluate(w http.ResponseWriter, r *http.Request) error {
	var payload struct {
		UserID   string `json:"user_id"`
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}
	resp, err := h.usecase.Evaluate(r.Context(), payload.UserID, payload.Resource, payload.Action)
	if err != nil {
		return err
	}
	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}
