package usecase

import (
	"context"
	"react-example/backend-golang/domain"
)

type reportUsecase struct{}

func NewReportUsecase() domain.ReportUsecase {
	return &reportUsecase{}
}

func (u *reportUsecase) GetAccessSummary(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{
		"period":          "2025-06",
		"total_users":     142,
		"active_users":    128,
		"suspended_users": 14,
		"by_department": []interface{}{
			map[string]interface{}{"department": "Engineering", "user_count": 45, "avg_risk_score": 18},
			map[string]interface{}{"department": "Finance", "user_count": 30, "avg_risk_score": 42},
			map[string]interface{}{"department": "Operations", "user_count": 53, "avg_risk_score": 25},
		},
		"by_role": []interface{}{
			map[string]interface{}{"role": "admin", "count": 5},
			map[string]interface{}{"role": "analyst", "count": 23},
			map[string]interface{}{"role": "viewer", "count": 114},
		},
	}, nil
}

func (u *reportUsecase) GetRiskTrend(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{
		"period": "last_30_days",
		"trend": []interface{}{
			map[string]interface{}{"date": "2025-05-11", "avg_risk_score": 21, "high_risk_users": 3},
			map[string]interface{}{"date": "2025-05-18", "avg_risk_score": 24, "high_risk_users": 5},
			map[string]interface{}{"date": "2025-05-25", "avg_risk_score": 19, "high_risk_users": 2},
			map[string]interface{}{"date": "2025-06-01", "avg_risk_score": 28, "high_risk_users": 7},
			map[string]interface{}{"date": "2025-06-08", "avg_risk_score": 22, "high_risk_users": 4},
		},
		"threshold": 70,
	}, nil
}
