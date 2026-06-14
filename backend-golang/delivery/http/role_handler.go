package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"react-example/backend-golang/delivery/http/middleware"
	"react-example/backend-golang/domain"
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
	middleware.SendJSON(w, http.StatusOK, roles, nil)
	return nil
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) error {
	var role domain.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}
	newRole, err := h.roleUsecase.CreateRole(r.Context(), role)
	if err != nil {
		return err
	}
	middleware.SendJSON(w, http.StatusCreated, newRole, nil)
	return nil
}

func (h *RoleHandler) ListPermissions(w http.ResponseWriter, r *http.Request) error {
	perms, err := h.roleUsecase.ListPermissions(r.Context())
	if err != nil {
		return err
	}
	middleware.SendJSON(w, http.StatusOK, perms, nil)
	return nil
}

func (h *RoleHandler) AssignUserRole(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		return middleware.NewCustomError(http.StatusBadRequest, "Missing user ID", nil)
	}
	userID := pathParts[4]

	var payload struct {
		RoleID   string `json:"roleId"`
		Operator string `json:"operator"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}

	err := h.roleUsecase.AssignRole(r.Context(), userID, payload.RoleID, payload.Operator)
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, map[string]string{"message": "Role assigned"}, nil)
	return nil
}
