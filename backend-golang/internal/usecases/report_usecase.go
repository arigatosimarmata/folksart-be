package usecases

import (
	"context"
	"react-example/backend-golang/internal/domain"
)

type reportUsecase struct{}

func NewReportUsecase() domain.ReportUsecase {
	return &reportUsecase{}
}

func (u *reportUsecase) GetAccessSummary(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{"total_access": 100}, nil
}

func (u *reportUsecase) GetRiskTrend(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{"trend": "downward"}, nil
}
