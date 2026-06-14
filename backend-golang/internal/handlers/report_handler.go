package handlers

import (
	"github.com/gofiber/fiber/v2"
	"react-example/backend-golang/httputil"
	"react-example/backend-golang/internal/domain"
)

type ReportHandler struct {
	usecase domain.ReportUsecase
}

func NewReportHandler(u domain.ReportUsecase) *ReportHandler {
	return &ReportHandler{usecase: u}
}

func (h *ReportHandler) AccessSummary(c *fiber.Ctx) error {
	resp, err := h.usecase.GetAccessSummary(c.Context())
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}
	return httputil.WriteSuccessResponse(c, "Success", resp, nil)
}

func (h *ReportHandler) RiskTrend(c *fiber.Ctx) error {
	resp, err := h.usecase.GetRiskTrend(c.Context())
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}
	return httputil.WriteSuccessResponse(c, "Success", resp, nil)
}
