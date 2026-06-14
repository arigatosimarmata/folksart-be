package http

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"react-example/backend-golang/delivery/http/middleware"
	"react-example/backend-golang/domain"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(uu domain.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: uu}
}

// ListUsers handles GET /api/v1/users
// @Summary List all user identity entries
// @Description Fetch a paginated list of users with optional filters for search, role, status, and department.
// @Tags users
// @Accept json
// @Produce json
// @Param search query string false "Search by name or email"
// @Param role query string false "Filter by role (e.g., Admin, User)"
// @Param status query string false "Filter by status (e.g., Active, Inactive)"
// @Param department query string false "Filter by department"
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of records per page" default(10)
// @Success 200 {object} map[string]interface{} "List of users with pagination metadata"
// @Failure 405 {object} middleware.APIError "Method Not Allowed"
// @Failure 500 {object} middleware.APIError "Internal Server Error"
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return middleware.NewCustomError(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
	}

	ctx := r.Context()

	// Parse filtering conditions
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
		return fmtDBError("Failed to query identity database pool", err)
	}

	totalPages := (total + limit - 1) / limit

	middleware.SendJSON(w, http.StatusOK, users, map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
		"pages": totalPages,
	})
	return nil
}

type EnrollRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Operator   string `json:"operator"`
}

// EnrollUser handles POST /api/v1/users
// @Summary Enroll a new user identity
// @Description Create a new user record in the governance directory.
// @Tags users
// @Accept json
// @Produce json
// @Param user body EnrollRequest true "User enrollment details"
// @Success 201 {object} domain.User "Newly created user object"
// @Failure 400 {object} middleware.APIError "Bad Request / Validation Error"
// @Failure 405 {object} middleware.APIError "Method Not Allowed"
// @Failure 500 {object} middleware.APIError "Internal Server Error"
// @Router /api/v1/users [post]
func (h *UserHandler) EnrollUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return middleware.NewCustomError(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
	}

	var req EnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid JSON payload structure", err)
	}

	ctx := r.Context()
	newUser, err := h.userUsecase.EnrollPrincipal(ctx, req.Name, req.Email, req.Role, req.Department, req.Operator)
	if err != nil {
		// Differentiate database validation/unique collisions vs general inputs
		if strings.Contains(err.Error(), "insufficient") {
			return middleware.NewCustomError(http.StatusBadRequest, err.Error(), nil)
		}
		return fmtDBError("Failed to enroll subject in directory", err)
	}

	middleware.SendJSON(w, http.StatusCreated, newUser, nil)
	return nil
}

type UpdateUserPayload struct {
	Status     *string `json:"status"`
	KYCStatus  *string `json:"kycStatus"`
	RiskScore  *int    `json:"riskScore"`
	MFAEnabled *bool   `json:"mfaEnabled"`
	Operator   string  `json:"operator"`
}

// UpdateUser handles PATCH /api/v1/users/:id
// @Summary Update user attributes
// @Description Patch specific fields of an existing user identity.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Reference ID"
// @Param update body UpdateUserPayload true "Fields to update"
// @Success 200 {object} domain.User "Updated user object"
// @Failure 400 {object} middleware.APIError "Bad Request"
// @Failure 405 {object} middleware.APIError "Method Not Allowed"
// @Failure 500 {object} middleware.APIError "Internal Server Error"
// @Router /api/v1/users/{id} [patch]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPatch {
		return middleware.NewCustomError(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 || pathParts[4] == "" {
		return middleware.NewCustomError(http.StatusBadRequest, "Missing identifier parameter in resource path", nil)
	}
	userID := pathParts[4]

	var req UpdateUserPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Payload decode failure: check formatting specs", err)
	}

	ctx := r.Context()
	patchedUser, err := h.userUsecase.PatchPrincipal(ctx, userID, req.Status, req.KYCStatus, req.RiskScore, req.MFAEnabled, req.Operator)
	if err != nil {
		return fmtDBError("Failed to update identity record attributes", err)
	}

	middleware.SendJSON(w, http.StatusOK, patchedUser, nil)
	return nil
}

// DeleteUser handles DELETE /api/v1/users/:id
// @Summary Decommission a user identity
// @Description Safely remove a user from the directory and log the decommissioning event.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User Reference ID"
// @Param operator query string true "Operator performing the deletion"
// @Success 200 {object} map[string]interface{} "Success confirmation message"
// @Failure 400 {object} middleware.APIError "Bad Request"
// @Failure 405 {object} middleware.APIError "Method Not Allowed"
// @Failure 500 {object} middleware.APIError "Internal Server Error"
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodDelete {
		return middleware.NewCustomError(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 || pathParts[4] == "" {
		return middleware.NewCustomError(http.StatusBadRequest, "Missing identifier parameter in resource path", nil)
	}
	userID := pathParts[4]
	operator := r.URL.Query().Get("operator")

	ctx := r.Context()
	err := h.userUsecase.DecommissionPrincipal(ctx, userID, operator)
	if err != nil {
		return fmtDBError("Failed to decommission identity from RDBMS pool", err)
	}

	middleware.SendJSON(w, http.StatusOK, map[string]interface{}{
		"message":           "Identity has been safely decommissioned.",
		"decommissioned_at": time.Now(),
	}, nil)
	return nil
}

// ExportCSV handles GET /api/v1/export/csv
// @Summary Export identity directory to CSV
// @Description Generates a CSV report of the user directory with active filters.
// @Tags export
// @Produce text/csv
// @Param department query string false "Filter by department"
// @Param role query string false "Filter by role"
// @Param status query string false "Filter by status"
// @Success 200 {file} file "CSV File Download"
// @Failure 405 {object} middleware.APIError "Method Not Allowed"
// @Failure 500 {object} middleware.APIError "Internal Server Error"
// @Router /api/v1/export/csv [get]
func (h *UserHandler) ExportCSV(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return middleware.NewCustomError(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
	}

	ctx := r.Context()

	filter := domain.UserFilter{
		Department: r.URL.Query().Get("department"),
		Role:       r.URL.Query().Get("role"),
		Status:     r.URL.Query().Get("status"),
	}

	users, err := h.userUsecase.ExportCSVStream(ctx, filter)
	if err != nil {
		return fmtDBError("Failed to capture datasets for CSV output alignment", err)
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=iam_governance_report.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write CSV headers
	if err := writer.Write([]string{"Reference ID", "Full Name", "Username", "Email Address", "Phone Number", "Assigned Role", "Account Status", "KYC Compliance Status", "Department", "Threat Risk Score", "MFA Armed Status", "Enrollment Date"}); err != nil {
		return middleware.NewCustomError(http.StatusInternalServerError, "Failed compiling CSV document headers", err)
	}

	for _, u := range users {
		mfaLabel := "Disarmed"
		if u.MFAEnabled {
			mfaLabel = "Armed"
		}
		row := []string{
			u.ID,
			u.Name,
			u.Username,
			u.Email,
			u.Phone,
			u.Role,
			u.Status,
			u.KYCStatus,
			u.Department,
			strconv.Itoa(u.RiskScore) + "%",
			mfaLabel,
			u.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return middleware.NewCustomError(http.StatusInternalServerError, "Failed compiling CSV document records", err)
		}
	}
	return nil
}

// fmtDBError evaluates deep database issues and translates them to standard CustomAppError,
// while letting the central middleware filter underlying MySQL code errors
func fmtDBError(contextMsg string, err error) error {
	if err == nil {
		return nil
	}
	// Return the raw error wrapped with custom message, letting the middleware parse constraints and SQL state
	return middleware.NewCustomError(http.StatusInternalServerError, contextMsg, err)
}
