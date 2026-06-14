package handlers

import (
	"github.com/gofiber/fiber/v2"
	"react-example/backend-golang/httputil"
	"react-example/backend-golang/internal/domain"
	"react-example/backend-golang/internal/dto"
)

type AccessRequestHandler struct {
	usecase domain.AccessRequestUsecase
}

func NewAccessRequestHandler(u domain.AccessRequestUsecase) *AccessRequestHandler {
	return &AccessRequestHandler{usecase: u}
}

func (h *AccessRequestHandler) Submit(c *fiber.Ctx) error {
	var req dto.CreateAccessRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	resp, err := h.usecase.SubmitRequest(c.Context(), domain.AccessRequest{
		RequesterID:   req.RequesterID,
		RequesterName: req.RequesterName,
		Resource:      req.Resource,
		AccessLevel:   req.AccessLevel,
		Justification: req.Justification,
	})
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Request submitted successfully", mapARToDTO(*resp), nil)
}

func (h *AccessRequestHandler) List(c *fiber.Ctx) error {
	status := c.Query("status")
	resp, err := h.usecase.ListRequests(c.Context(), status)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	dtos := make([]dto.AccessRequestResponse, 0)
	for _, req := range resp {
		dtos = append(dtos, mapARToDTO(req))
	}

	return httputil.WriteSuccessResponse(c, "Success", dtos, map[string]interface{}{"total": len(dtos)})
}

func (h *AccessRequestHandler) Approve(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return httputil.WriteErrorResponse(c, fiber.ErrBadRequest)
	}

	var req dto.ApproveAccessRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	resp, err := h.usecase.ApproveRequest(c.Context(), id, req.Operator, req.Note)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Request approved successfully", mapARToDTO(*resp), nil)
}

func (h *AccessRequestHandler) Reject(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return httputil.WriteErrorResponse(c, fiber.ErrBadRequest)
	}

	var req dto.RejectAccessRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	resp, err := h.usecase.RejectRequest(c.Context(), id, req.Operator, req.Reason)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Request rejected successfully", mapARToDTO(*resp), nil)
}

func mapARToDTO(r domain.AccessRequest) dto.AccessRequestResponse {
	res := dto.AccessRequestResponse{
		ID:            r.ID,
		RequesterID:   r.RequesterID,
		RequesterName: r.RequesterName,
		Resource:      r.Resource,
		AccessLevel:   r.AccessLevel,
		Justification: r.Justification,
		Status:        r.Status,
		RequestedAt:   r.RequestedAt.Format("2006-01-02 15:04:05"),
		ApprovedBy:    r.ApprovedBy,
	}
	if r.ApprovedAt != nil {
		t := r.ApprovedAt.Format("2006-01-02 15:04:05")
		res.ApprovedAt = &t
	}
	if r.ExpiresAt != nil {
		t := r.ExpiresAt.Format("2006-01-02 15:04:05")
		res.ExpiresAt = &t
	}
	return res
}
