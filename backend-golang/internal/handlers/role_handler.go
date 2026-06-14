package handlers

import (
	"github.com/gofiber/fiber/v2"
	"react-example/backend-golang/httputil"
	"react-example/backend-golang/internal/domain"
	"react-example/backend-golang/internal/dto"
)

type RoleHandler struct {
	roleUsecase domain.RoleUsecase
}

func NewRoleHandler(ru domain.RoleUsecase) *RoleHandler {
	return &RoleHandler{roleUsecase: ru}
}

func (h *RoleHandler) ListRoles(c *fiber.Ctx) error {
	roles, err := h.roleUsecase.ListRoles(c.Context())
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	roleDtos := make([]dto.RoleResponse, 0)
	for _, role := range roles {
		roleDtos = append(roleDtos, dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: role.Permissions,
			UserCount:   role.UserCount,
		})
	}

	return httputil.WriteSuccessResponse(c, "Success", roleDtos, nil)
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req dto.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	newRole, err := h.roleUsecase.CreateRole(c.Context(), domain.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
	})
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	res := dto.RoleResponse{
		ID:          newRole.ID,
		Name:        newRole.Name,
		Description: newRole.Description,
		Permissions: newRole.Permissions,
		UserCount:   newRole.UserCount,
	}

	return httputil.WriteSuccessResponse(c, "Role created successfully", res, nil)
}

func (h *RoleHandler) ListPermissions(c *fiber.Ctx) error {
	perms, err := h.roleUsecase.ListPermissions(c.Context())
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	permDtos := make([]dto.PermissionResponse, 0)
	for _, perm := range perms {
		permDtos = append(permDtos, dto.PermissionResponse{
			ID:          perm.ID,
			Key:         perm.Key,
			Description: perm.Description,
		})
	}

	return httputil.WriteSuccessResponse(c, "Success", permDtos, nil)
}

func (h *RoleHandler) AssignUserRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return httputil.WriteErrorResponse(c, fiber.ErrBadRequest)
	}

	var req dto.AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	err := h.roleUsecase.AssignRole(c.Context(), userID, req.RoleID, req.Operator)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Role assigned successfully", nil, nil)
}
