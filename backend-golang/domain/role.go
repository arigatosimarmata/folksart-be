package domain

import "context"

type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	UserCount   int      `json:"user_count"`
}

type Permission struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Description string `json:"description"`
}

type RoleUsecase interface {
	ListRoles(ctx context.Context) ([]Role, error)
	CreateRole(ctx context.Context, role Role) (*Role, error)
	ListPermissions(ctx context.Context) ([]Permission, error)
	AssignRole(ctx context.Context, userID, roleID, operator string) error
	RemoveRole(ctx context.Context, userID, roleID, operator string) error
}
