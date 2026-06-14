package handlers

import (
	"github.com/gofiber/fiber/v2"
	"react-example/backend-golang/httputil"
	"react-example/backend-golang/internal/domain"
	"react-example/backend-golang/internal/dto"
)

type KYCHandler struct {
	usecase domain.KYCUsecase
}

func NewKYCHandler(u domain.KYCUsecase) *KYCHandler {
	return &KYCHandler{usecase: u}
}

func (h *KYCHandler) Submit(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return httputil.WriteErrorResponse(c, fiber.ErrBadRequest)
	}

	var reqDocs []dto.KYCDocumentDTO
	if err := c.BodyParser(&reqDocs); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	domainDocs := make([]domain.KYCDocument, 0)
	for _, d := range reqDocs {
		domainDocs = append(domainDocs, domain.KYCDocument{
			Type: d.Type,
		})
	}

	resp, err := h.usecase.SubmitKYC(c.Context(), userID, domainDocs)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "KYC submitted successfully", mapKYCStatusToDTO(resp), nil)
}

func (h *KYCHandler) Status(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return httputil.WriteErrorResponse(c, fiber.ErrBadRequest)
	}

	resp, err := h.usecase.GetKYCStatus(c.Context(), userID)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Success", mapKYCStatusToDTO(resp), nil)
}

func (h *KYCHandler) Review(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return httputil.WriteErrorResponse(c, fiber.ErrBadRequest)
	}

	var req dto.KYCReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	resp, err := h.usecase.ReviewKYC(c.Context(), userID, req.Operator, req.Status, req.Note)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "KYC reviewed successfully", mapKYCStatusToDTO(resp), nil)
}

func (h *KYCHandler) UploadToken(c *fiber.Ctx) error {
	resp, err := h.usecase.IssueUploadToken(c.Context(), "usr_current")
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Success", resp, nil)
}

func mapKYCStatusToDTO(s *domain.KYCStatus) dto.KYCStatusResponse {
	res := dto.KYCStatusResponse{
		UserID:          s.UserID,
		Status:          s.Status,
		ReviewedBy:      s.ReviewedBy,
		RejectionReason: s.RejectionReason,
		Documents:       make([]dto.KYCDocumentDTO, 0),
	}

	if s.SubmittedAt != nil {
		t := s.SubmittedAt.Format("2006-01-02 15:04:05")
		res.SubmittedAt = &t
	}
	if s.ReviewedAt != nil {
		t := s.ReviewedAt.Format("2006-01-02 15:04:05")
		res.ReviewedAt = &t
	}

	for _, d := range s.Documents {
		res.Documents = append(res.Documents, dto.KYCDocumentDTO{
			Type:       d.Type,
			Status:     d.Status,
			UploadedAt: d.UploadedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return res
}
