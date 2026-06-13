package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"react-example/backend-golang/domain"
)

type mysqlUserRepository struct {
	db *sql.DB
}

// NewMySQLUserRepository yields a repository wrapper implementing domain.UserRepository interface (Liskov Substitution)
func NewMySQLUserRepository(db *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) GetByID(ctx context.Context, id string) (*domain.IAMUser, error) {
	query := `SELECT id, name, username, email, phone, role, status, kyc_status, department, risk_score, mfa_enabled, created_at FROM iam_users WHERE id = ?`
	var u domain.IAMUser
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Username, &u.Email, &u.Phone, &u.Role, &u.Status, &u.KYCStatus, &u.Department, &u.RiskScore, &u.MFAEnabled, &u.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Return nil with no error if record doesn't exist
	}
	if err != nil {
		return nil, fmt.Errorf("failed scanning IAMUser record: %w", err)
	}
	return &u, nil
}

func (r *mysqlUserRepository) List(ctx context.Context, filter domain.UserFilter) ([]domain.IAMUser, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Search != "" {
		conditions = append(conditions, "(name LIKE ? OR username LIKE ? OR email LIKE ? OR id = ?)")
		wildcard := "%" + filter.Search + "%"
		args = append(args, wildcard, wildcard, wildcard, filter.Search)
	}
	if filter.Role != "" {
		conditions = append(conditions, "role = ?")
		args = append(args, filter.Role)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.Department != "" {
		conditions = append(conditions, "department = ?")
		args = append(args, filter.Department)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Calculate count first
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM iam_users %s", whereClause)
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting filtered users: %w", err)
	}

	// Append limit offsets
	limitOffsetClause := ""
	if filter.Limit > 0 {
		limitOffsetClause = " LIMIT ? OFFSET ?"
		args = append(args, filter.Limit, filter.Offset)
	}

	fetchQuery := fmt.Sprintf("SELECT id, name, username, email, phone, role, status, kyc_status, department, risk_score, mfa_enabled, created_at FROM iam_users %s ORDER BY created_at DESC%s", whereClause, limitOffsetClause)
	rows, err := r.db.QueryContext(ctx, fetchQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed fetching user directories slice: %w", err)
	}
	defer rows.Close()

	users := []domain.IAMUser{}
	for rows.Next() {
		var u domain.IAMUser
		err := rows.Scan(
			&u.ID, &u.Name, &u.Username, &u.Email, &u.Phone, &u.Role, &u.Status, &u.KYCStatus, &u.Department, &u.RiskScore, &u.MFAEnabled, &u.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed mapping User entity: %w", err)
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (r *mysqlUserRepository) Store(ctx context.Context, u *domain.IAMUser) error {
	query := `INSERT INTO iam_users (id, name, username, email, phone, role, status, kyc_status, department, risk_score, mfa_enabled, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Name, u.Username, u.Email, u.Phone, u.Role, u.Status, u.KYCStatus, u.Department, u.RiskScore, u.MFAEnabled, u.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed writing new User record to mysql: %w", err)
	}
	return nil
}

func (r *mysqlUserRepository) Update(ctx context.Context, u *domain.IAMUser) error {
	query := `UPDATE iam_users SET name = ?, email = ?, role = ?, status = ?, kyc_status = ?, department = ?, risk_score = ?, mfa_enabled = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.Role, u.Status, u.KYCStatus, u.Department, u.RiskScore, u.MFAEnabled, u.ID)
	if err != nil {
		return fmt.Errorf("failed modifying user record in mysql: %w", err)
	}
	return nil
}

func (r *mysqlUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM iam_users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed executing user decommission in mysql: %w", err)
	}
	return nil
}
