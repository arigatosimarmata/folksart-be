package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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

func (h *KYCHandler) Submit(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "01", "Missing user ID")
		return nil
	}
	userID := pathParts[4]

	var reqDocs []dto.KYCDocumentDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDocs); err != nil {
		return err
	}

	domainDocs := make([]domain.KYCDocument, 0)
	for _, d := range reqDocs {
		domainDocs = append(domainDocs, domain.KYCDocument{
			Type: d.Type,
		})
	}

	resp, err := h.usecase.SubmitKYC(r.Context(), userID, domainDocs)
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "KYC submitted successfully", mapKYCStatusToDTO(resp), nil)
	return nil
}

func (h *KYCHandler) Status(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "01", "Missing user ID")
		return nil
	}
	userID := pathParts[4]

	resp, err := h.usecase.GetKYCStatus(r.Context(), userID)
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "Success", mapKYCStatusToDTO(resp), nil)
	return nil
}

func (h *KYCHandler) Review(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "01", "Missing user ID")
		return nil
	}
	userID := pathParts[4]

	var req dto.KYCReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	resp, err := h.usecase.ReviewKYC(r.Context(), userID, req.Operator, req.Status, req.Note)
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "KYC reviewed successfully", mapKYCStatusToDTO(resp), nil)
	return nil
}

func (h *KYCHandler) UploadToken(w http.ResponseWriter, r *http.Request) error {
	resp, err := h.usecase.IssueUploadToken(r.Context(), "usr_current")
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "Success", resp, nil)
	return nil
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
