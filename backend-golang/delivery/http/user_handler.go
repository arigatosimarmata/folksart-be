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

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": totalPages,
		},
	})
}

type EnrollRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Operator   string `json:"operator"`
}

// EnrollUser handles POST /api/v1/users
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(newUser)
}

type UpdateUserPayload struct {
	Status     *string `json:"status"`
	KYCStatus  *string `json:"kycStatus"`
	RiskScore  *int    `json:"riskScore"`
	MFAEnabled *bool   `json:"mfaEnabled"`
	Operator   string  `json:"operator"`
}

// UpdateUser handles PATCH /api/v1/users/:id
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

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(patchedUser)
}

// DeleteUser handles DELETE /api/v1/users/:id
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

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Identity has been safely decommissioned.",
	})
}

// ExportCSV handles GET /api/v1/export/csv
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
