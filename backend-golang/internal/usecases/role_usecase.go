package usecases

import (
	"context"
	"react-example/backend-golang/internal/domain"
)

type roleUsecase struct{}

func NewRoleUsecase() domain.RoleUsecase {
	return &roleUsecase{}
}

func (u *roleUsecase) ListRoles(ctx context.Context) ([]domain.Role, error) {
	return []domain.Role{
		{ID: "role-1", Name: "Administrator", Description: "Full system access", UserCount: 5},
	}, nil
}

func (u *roleUsecase) CreateRole(ctx context.Context, role domain.Role) (*domain.Role, error) {
	role.ID = "role-new"
	return &role, nil
}

func (u *roleUsecase) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	return []domain.Permission{
		{ID: "perm-1", Key: "user:write", Description: "Create and update users"},
	}, nil
}

func (u *roleUsecase) AssignRole(ctx context.Context, userID, roleID, operator string) error {
	return nil
}

func (u *roleUsecase) RemoveRole(ctx context.Context, userID, roleID, operator string) error {
	return nil
}
