package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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

func (h *AccessRequestHandler) Submit(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	resp, err := h.usecase.SubmitRequest(r.Context(), domain.AccessRequest{
		RequesterID:   req.RequesterID,
		RequesterName: req.RequesterName,
		Resource:      req.Resource,
		AccessLevel:   req.AccessLevel,
		Justification: req.Justification,
	})
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "Request submitted successfully", mapARToDTO(*resp), nil)
	return nil
}

func (h *AccessRequestHandler) List(w http.ResponseWriter, r *http.Request) error {
	status := r.URL.Query().Get("status")
	resp, err := h.usecase.ListRequests(r.Context(), status)
	if err != nil {
		return err
	}

	dtos := make([]dto.AccessRequestResponse, 0)
	for _, req := range resp {
		dtos = append(dtos, mapARToDTO(req))
	}

	httputil.WriteSuccessResponse(w, "Success", dtos, map[string]interface{}{"total": len(dtos)})
	return nil
}

func (h *AccessRequestHandler) Approve(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "01", "Missing request ID")
		return nil
	}
	id := pathParts[4]

	var req dto.ApproveAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	resp, err := h.usecase.ApproveRequest(r.Context(), id, req.Operator, req.Note)
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "Request approved successfully", mapARToDTO(*resp), nil)
	return nil
}

func (h *AccessRequestHandler) Reject(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "01", "Missing request ID")
		return nil
	}
	id := pathParts[4]

	var req dto.RejectAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	resp, err := h.usecase.RejectRequest(r.Context(), id, req.Operator, req.Reason)
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "Request rejected successfully", mapARToDTO(*resp), nil)
	return nil
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
