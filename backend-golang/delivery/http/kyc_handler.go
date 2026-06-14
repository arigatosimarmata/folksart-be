package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"react-example/backend-golang/delivery/http/middleware"
	"react-example/backend-golang/domain"
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
		return middleware.NewCustomError(http.StatusBadRequest, "Missing user ID", nil)
	}
	userID := pathParts[4]

	var docs []domain.KYCDocument
	if err := json.NewDecoder(r.Body).Decode(&docs); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}

	resp, err := h.usecase.SubmitKYC(r.Context(), userID, docs)
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}

func (h *KYCHandler) Status(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		return middleware.NewCustomError(http.StatusBadRequest, "Missing user ID", nil)
	}
	userID := pathParts[4]

	resp, err := h.usecase.GetKYCStatus(r.Context(), userID)
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}

func (h *KYCHandler) Review(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		return middleware.NewCustomError(http.StatusBadRequest, "Missing user ID", nil)
	}
	userID := pathParts[4]

	var payload struct {
		Operator string `json:"operator"`
		Status   string `json:"status"`
		Note     string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}

	resp, err := h.usecase.ReviewKYC(r.Context(), userID, payload.Operator, payload.Status, payload.Note)
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}

func (h *KYCHandler) UploadToken(w http.ResponseWriter, r *http.Request) error {
	resp, err := h.usecase.IssueUploadToken(r.Context(), "usr_current")
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}
