package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"react-example/backend-golang/errs"
	"react-example/backend-golang/httputil"
	"react-example/backend-golang/internal/domain"
	"react-example/backend-golang/internal/dto"
	"react-example/backend-golang/internal/validation"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(uu domain.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: uu}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	search := r.URL.Query().Get("search")
	role := r.URL.Query().Get("role")
	status := r.URL.Query().Get("status")
	department := r.URL.Query().Get("department")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr == "" {
		pageStr = "1"
	}
	if limitStr == "" {
		limitStr = "10"
	}

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	filter := domain.UserFilter{
		Search:     search,
		Role:       role,
		Status:     status,
		Department: department,
		Limit:      limit,
		Offset:     offset,
	}

	users, total, err := h.userUsecase.FetchDirectories(ctx, filter)
	if err != nil {
		return err
	}

	userDtos := make([]dto.UserResponse, 0)
	for _, u := range users {
		userDtos = append(userDtos, dto.UserResponse{
			ID:         u.ID,
			Name:       u.Name,
			Username:   u.Username,
			Email:      u.Email,
			Role:       u.Role,
			Status:     u.Status,
			KYCStatus:  u.KYCStatus,
			Department: u.Department,
			RiskScore:  u.RiskScore,
			CreatedAt:  u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPages := (total + limit - 1) / limit

	httputil.WriteSuccessResponse(w, "Success", userDtos, map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
		"pages": totalPages,
	})
	return nil
}

func (h *UserHandler) EnrollUser(w http.ResponseWriter, r *http.Request) error {
	var req dto.EnrollUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.ErrBadRequest
	}

	// 1. Centralized Validation
	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httputil.WriteValidationErrorResponse(w, validationErrors)
		return nil
	}

	ctx := r.Context()
	newUser, err := h.userUsecase.EnrollPrincipal(ctx, req.Name, req.Email, req.Role, req.Department, req.Operator)
	if err != nil {
		// 2. Automated Error Translation via httputil
		httputil.WriteErrorResponse(w, err)
		return nil
	}

	res := dto.UserResponse{
		ID:         newUser.ID,
		Name:       newUser.Name,
		Username:   newUser.Username,
		Email:      newUser.Email,
		Role:       newUser.Role,
		Status:     newUser.Status,
		KYCStatus:  newUser.KYCStatus,
		Department: newUser.Department,
		RiskScore:  newUser.RiskScore,
		CreatedAt:  newUser.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	httputil.WriteSuccessResponse(w, "User enrolled successfully", res, nil)
	return nil
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 || pathParts[4] == "" {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "01", "Missing user ID")
		return nil
	}
	userID := pathParts[4]

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	ctx := r.Context()
	patchedUser, err := h.userUsecase.PatchPrincipal(ctx, userID, req.Status, req.KYCStatus, req.RiskScore, req.MFAEnabled, req.Operator)
	if err != nil {
		return err
	}

	res := dto.UserResponse{
		ID:         patchedUser.ID,
		Name:       patchedUser.Name,
		Username:   patchedUser.Username,
		Email:      patchedUser.Email,
		Role:       patchedUser.Role,
		Status:     patchedUser.Status,
		KYCStatus:  patchedUser.KYCStatus,
		Department: patchedUser.Department,
		RiskScore:  patchedUser.RiskScore,
		CreatedAt:  patchedUser.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	httputil.WriteSuccessResponse(w, "User updated successfully", res, nil)
	return nil
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 || pathParts[4] == "" {
		httputil.WriteErrorResponse(w, http.StatusBadRequest, "01", "Missing user ID")
		return nil
	}
	userID := pathParts[4]
	operator := r.URL.Query().Get("operator")

	ctx := r.Context()
	err := h.userUsecase.DecommissionPrincipal(ctx, userID, operator)
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "User deleted successfully", nil, nil)
	return nil
}
