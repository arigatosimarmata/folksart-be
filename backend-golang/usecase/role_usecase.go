package usecase

import (
	"context"
	"react-example/backend-golang/domain"
)

type roleUsecase struct{}

func NewRoleUsecase() domain.RoleUsecase {
	return &roleUsecase{}
}

func (r *roleUsecase) ListRoles(ctx context.Context) ([]domain.Role, error) {
	return []domain.Role{
		{
			ID:          "role_admin",
			Name:        "admin",
			Description: "Full system access",
			Permissions: []string{"users:read", "users:write", "audit:read", "policy:write"},
			UserCount:   5,
		},
		{
			ID:          "role_analyst",
			Name:        "analyst",
			Description: "Read-only access to reports and audit logs",
			Permissions: []string{"audit:read", "reports:read"},
			UserCount:   23,
		},
	}, nil
}

func (r *roleUsecase) CreateRole(ctx context.Context, role domain.Role) (*domain.Role, error) {
	return &role, nil
}

func (r *roleUsecase) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	return []domain.Permission{
		{ID: "perm_01", Key: "users:read", Description: "Lihat data user"},
		{ID: "perm_02", Key: "users:write", Description: "Buat/edit/hapus user"},
		{ID: "perm_03", Key: "audit:read", Description: "Lihat audit log"},
	}, nil
}

func (r *roleUsecase) AssignRole(ctx context.Context, userID, roleID, operator string) error {
	return nil
}

func (r *roleUsecase) RemoveRole(ctx context.Context, userID, roleID, operator string) error {
	return nil
}
