package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"react-example/backend-golang/internal/domain"
)

type mysqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) GetByID(ctx context.Context, id string) (*domain.IAMUser, error) {
	query := `SELECT id, name, username, email, phone, role, status, kyc_status, department, risk_score, mfa_enabled, created_at FROM iam_users WHERE id = ?`
	var u domain.IAMUser
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Username, &u.Email, &u.Phone, &u.Role, &u.Status, &u.KYCStatus, &u.Department, &u.RiskScore, &u.MFAEnabled, &u.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
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

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM iam_users %s", whereClause)
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	limitOffsetClause := ""
	if filter.Limit > 0 {
		limitOffsetClause = " LIMIT ? OFFSET ?"
		args = append(args, filter.Limit, filter.Offset)
	}

	fetchQuery := fmt.Sprintf("SELECT id, name, username, email, phone, role, status, kyc_status, department, risk_score, mfa_enabled, created_at FROM iam_users %s ORDER BY created_at DESC%s", whereClause, limitOffsetClause)
	rows, err := r.db.QueryContext(ctx, fetchQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []domain.IAMUser{}
	for rows.Next() {
		var u domain.IAMUser
		err := rows.Scan(
			&u.ID, &u.Name, &u.Username, &u.Email, &u.Phone, &u.Role, &u.Status, &u.KYCStatus, &u.Department, &u.RiskScore, &u.MFAEnabled, &u.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (r *mysqlUserRepository) Store(ctx context.Context, u *domain.IAMUser) error {
	query := `INSERT INTO iam_users (id, name, username, email, phone, role, status, kyc_status, department, risk_score, mfa_enabled, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Name, u.Username, u.Email, u.Phone, u.Role, u.Status, u.KYCStatus, u.Department, u.RiskScore, u.MFAEnabled, u.CreatedAt)
	return err
}

func (r *mysqlUserRepository) Update(ctx context.Context, u *domain.IAMUser) error {
	query := `UPDATE iam_users SET name = ?, email = ?, role = ?, status = ?, kyc_status = ?, department = ?, risk_score = ?, mfa_enabled = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.Role, u.Status, u.KYCStatus, u.Department, u.RiskScore, u.MFAEnabled, u.ID)
	return err
}

func (r *mysqlUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM iam_users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
