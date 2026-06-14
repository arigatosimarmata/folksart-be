package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"react-example/backend-golang/httputil"
	"react-example/backend-golang/internal/domain"
	"react-example/backend-golang/internal/dto"
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

	dtos := make([]dto.PolicyResponse, 0)
	for _, p := range resp {
		dtos = append(dtos, mapPolicyToDTO(p))
	}

	httputil.WriteSuccessResponse(w, "Success", dtos, nil)
	return nil
}

func (h *PolicyHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	resp, err := h.usecase.CreatePolicy(r.Context(), domain.Policy{
		Name:        req.Name,
		Description: req.Description,
		Condition:   req.Condition,
		Action:      req.Action,
		Priority:    req.Priority,
		Active:      req.Active,
	})
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "Policy created successfully", mapPolicyToDTO(*resp), nil)
	return nil
}

func (h *PolicyHandler) Update(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "01", "Missing policy ID")
		return nil
	}
	id := pathParts[4]

	var req dto.CreatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	resp, err := h.usecase.UpdatePolicy(r.Context(), id, domain.Policy{
		Name:        req.Name,
		Description: req.Description,
		Condition:   req.Condition,
		Action:      req.Action,
		Priority:    req.Priority,
		Active:      req.Active,
	})
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "Policy updated successfully", mapPolicyToDTO(*resp), nil)
	return nil
}

func (h *PolicyHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "01", "Missing policy ID")
		return nil
	}
	id := pathParts[4]

	err := h.usecase.DeletePolicy(r.Context(), id)
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "Policy deleted successfully", nil, nil)
	return nil
}

func (h *PolicyHandler) Evaluate(w http.ResponseWriter, r *http.Request) error {
	var req dto.EvaluatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	resp, err := h.usecase.Evaluate(r.Context(), req.UserID, req.Resource, req.Action)
	if err != nil {
		return err
	}

	res := dto.PolicyEvaluationResponse{
		UserID:   resp.UserID,
		Resource: resp.Resource,
		Action:   resp.Action,
		Decision: resp.Decision,
	}

	if resp.MatchedPolicy != nil {
		p := mapPolicyToDTO(*resp.MatchedPolicy)
		res.MatchedPolicy = &p
	}

	httputil.WriteSuccessResponse(w, "Evaluation complete", res, nil)
	return nil
}

func mapPolicyToDTO(p domain.Policy) dto.PolicyResponse {
	return dto.PolicyResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Condition:   p.Condition,
		Action:      p.Action,
		Priority:    p.Priority,
		Active:      p.Active,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
