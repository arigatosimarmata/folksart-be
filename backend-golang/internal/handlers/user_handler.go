package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
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

// ListUsers godoc
// @Summary List all users with filtering
// @Description Get a paginated list of users with search and filter options
// @Tags Identity
// @Accept json
// @Produce json
// @Param search query string false "Search by name or email"
// @Param role query string false "Filter by role"
// @Param status query string false "Filter by status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {object} httputil.Response{data=[]dto.UserResponse}
// @Router /users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	search := c.Query("search")
	role := c.Query("role")
	status := c.Query("status")
	department := c.Query("department")
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

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

	users, total, err := h.userUsecase.FetchDirectories(c.Context(), filter)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
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

	return httputil.WriteSuccessResponse(c, "Success", userDtos, map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
		"pages": totalPages,
	})
}

// EnrollUser godoc
// @Summary Enroll a new principal
// @Description Create a new user identity in the system
// @Tags Identity
// @Accept json
// @Produce json
// @Param request body dto.EnrollUserRequest true "User Enrollment Data"
// @Success 201 {object} httputil.Response{data=dto.UserResponse}
// @Failure 400 {object} httputil.Response
// @Security ApiKeyAuth
// @Router /users [post]
func (h *UserHandler) EnrollUser(c *fiber.Ctx) error {
	var req dto.EnrollUserRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, errs.ErrBadRequest)
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		return httputil.WriteValidationErrorResponse(c, validationErrors)
	}

	newUser, err := h.userUsecase.EnrollPrincipal(c.Context(), req.Name, req.Email, req.Role, req.Department, req.Operator)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
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

	return httputil.WriteSuccessResponse(c, "User enrolled successfully", res, nil)
}

// UpdateUser godoc
// @Summary Update principal attributes
// @Description Patch specific fields of a user identity
// @Tags Identity
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "Update Data"
// @Success 200 {object} httputil.Response{data=dto.UserResponse}
// @Security ApiKeyAuth
// @Router /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return httputil.WriteErrorResponse(c, errs.ErrBadRequest)
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, errs.ErrBadRequest)
	}

	patchedUser, err := h.userUsecase.PatchPrincipal(c.Context(), userID, req.Status, req.KYCStatus, req.RiskScore, req.MFAEnabled, req.Operator)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
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

	return httputil.WriteSuccessResponse(c, "User updated successfully", res, nil)
}

// DeleteUser godoc
// @Summary Decommission principal
// @Description Remove a user from the active directory
// @Tags Identity
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param operator query string true "Operator ID"
// @Success 200 {object} httputil.Response
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return httputil.WriteErrorResponse(c, errs.ErrBadRequest)
	}
	operator := c.Query("operator")

	err := h.userUsecase.DecommissionPrincipal(c.Context(), userID, operator)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "User deleted successfully", nil, nil)
}
