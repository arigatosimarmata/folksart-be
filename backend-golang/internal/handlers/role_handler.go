package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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

func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) error {
	roles, err := h.roleUsecase.ListRoles(r.Context())
	if err != nil {
		return err
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

	httputil.WriteSuccessResponse(w, "Success", roleDtos, nil)
	return nil
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	newRole, err := h.roleUsecase.CreateRole(r.Context(), domain.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
	})
	if err != nil {
		return err
	}

	res := dto.RoleResponse{
		ID:          newRole.ID,
		Name:        newRole.Name,
		Description: newRole.Description,
		Permissions: newRole.Permissions,
		UserCount:   newRole.UserCount,
	}

	httputil.WriteSuccessResponse(w, "Role created successfully", res, nil)
	return nil
}

func (h *RoleHandler) ListPermissions(w http.ResponseWriter, r *http.Request) error {
	perms, err := h.roleUsecase.ListPermissions(r.Context())
	if err != nil {
		return err
	}

	permDtos := make([]dto.PermissionResponse, 0)
	for _, perm := range perms {
		permDtos = append(permDtos, dto.PermissionResponse{
			ID:          perm.ID,
			Key:         perm.Key,
			Description: perm.Description,
		})
	}

	httputil.WriteSuccessResponse(w, "Success", permDtos, nil)
	return nil
}

func (h *RoleHandler) AssignUserRole(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "01", "Missing user ID")
		return nil
	}
	userID := pathParts[4]

	var req dto.AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	err := h.roleUsecase.AssignRole(r.Context(), userID, req.RoleID, req.Operator)
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "Role assigned successfully", nil, nil)
	return nil
}
