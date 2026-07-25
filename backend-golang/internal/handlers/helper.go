package handlers

import (
	"folksart-be/backend-golang/internal/domain"
	"folksart-be/backend-golang/internal/dto"
)

func toUserDTO(u *domain.IAMUser) *dto.UserResponse {
	if u == nil {
		return nil
	}
	return &dto.UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Username:   u.Username,
		Email:      u.Email,
		Role:       u.Role,
		Status:     u.Status,
		KYCStatus:  u.KYCStatus,
		Department: u.Department,
		RiskScore:  u.RiskScore,
		UserType:   u.UserType,
		TenantID:   u.TenantID,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
