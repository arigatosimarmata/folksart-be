package domain

import "context"

type Role struct {
	ID          string
	Name        string
	Description string
	Permissions []string
	UserCount   int
}

type Permission struct {
	ID          string
	Key         string
	Description string
}

type RoleUsecase interface {
	ListRoles(ctx context.Context) ([]Role, error)
	CreateRole(ctx context.Context, role Role) (*Role, error)
	ListPermissions(ctx context.Context) ([]Permission, error)
	AssignRole(ctx context.Context, userID, roleID, operator string) error
	RemoveRole(ctx context.Context, userID, roleID, operator string) error
}
