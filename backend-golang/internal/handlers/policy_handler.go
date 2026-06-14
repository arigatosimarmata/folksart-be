package handlers

import (
	"github.com/gofiber/fiber/v2"
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

func (h *PolicyHandler) List(c *fiber.Ctx) error {
	resp, err := h.usecase.ListPolicies(c.Context())
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	dtos := make([]dto.PolicyResponse, 0)
	for _, p := range resp {
		dtos = append(dtos, mapPolicyToDTO(p))
	}

	return httputil.WriteSuccessResponse(c, "Success", dtos, nil)
}

func (h *PolicyHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePolicyRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	resp, err := h.usecase.CreatePolicy(c.Context(), domain.Policy{
		Name:        req.Name,
		Description: req.Description,
		Condition:   req.Condition,
		Action:      req.Action,
		Priority:    req.Priority,
		Active:      req.Active,
	})
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Policy created successfully", mapPolicyToDTO(*resp), nil)
}

func (h *PolicyHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return httputil.WriteErrorResponse(c, fiber.ErrBadRequest)
	}

	var req dto.CreatePolicyRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	resp, err := h.usecase.UpdatePolicy(c.Context(), id, domain.Policy{
		Name:        req.Name,
		Description: req.Description,
		Condition:   req.Condition,
		Action:      req.Action,
		Priority:    req.Priority,
		Active:      req.Active,
	})
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Policy updated successfully", mapPolicyToDTO(*resp), nil)
}

func (h *PolicyHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return httputil.WriteErrorResponse(c, fiber.ErrBadRequest)
	}

	err := h.usecase.DeletePolicy(c.Context(), id)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Policy deleted successfully", nil, nil)
}

func (h *PolicyHandler) Evaluate(c *fiber.Ctx) error {
	var req dto.EvaluatePolicyRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	resp, err := h.usecase.Evaluate(c.Context(), req.UserID, req.Resource, req.Action)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
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

	return httputil.WriteSuccessResponse(c, "Evaluation complete", res, nil)
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
