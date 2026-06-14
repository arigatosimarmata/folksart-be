package domain

import "context"

type ReportUsecase interface {
	GetAccessSummary(ctx context.Context) (map[string]interface{}, error)
	GetRiskTrend(ctx context.Context) (map[string]interface{}, error)
}
