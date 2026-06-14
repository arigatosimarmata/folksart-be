package handlers

import (
	"net/http"

	"react-example/backend-golang/httputil"
	"react-example/backend-golang/internal/domain"
)

type ReportHandler struct {
	usecase domain.ReportUsecase
}

func NewReportHandler(u domain.ReportUsecase) *ReportHandler {
	return &ReportHandler{usecase: u}
}

func (h *ReportHandler) AccessSummary(w http.ResponseWriter, r *http.Request) error {
	resp, err := h.usecase.GetAccessSummary(r.Context())
	if err != nil {
		return err
	}
	httputil.WriteSuccessResponse(w, "Success", resp, nil)
	return nil
}

func (h *ReportHandler) RiskTrend(w http.ResponseWriter, r *http.Request) error {
	resp, err := h.usecase.GetRiskTrend(r.Context())
	if err != nil {
		return err
	}
	httputil.WriteSuccessResponse(w, "Success", resp, nil)
	return nil
}
