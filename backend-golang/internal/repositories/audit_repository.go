package repositories

import (
	"context"
	"database/sql"

	"react-example/backend-golang/internal/domain"
)

type mysqlAuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) domain.AuditRepository {
	return &mysqlAuditRepository{db: db}
}

func (r *mysqlAuditRepository) Store(ctx context.Context, log *domain.AuditLog) error {
	query := `INSERT INTO audit_logs (id, timestamp, actor, action, target, severity) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.Timestamp, log.Actor, log.Action, log.Target, log.Severity)
	return err
}

func (r *mysqlAuditRepository) List(ctx context.Context, severity string, limit int) ([]domain.AuditLog, error) {
	var rows *sql.Rows
	var err error

	if severity != "" {
		query := `SELECT id, timestamp, actor, action, target, severity FROM audit_logs WHERE severity = ? ORDER BY timestamp DESC LIMIT ?`
		rows, err = r.db.QueryContext(ctx, query, severity, limit)
	} else {
		query := `SELECT id, timestamp, actor, action, target, severity FROM audit_logs ORDER BY timestamp DESC LIMIT ?`
		rows, err = r.db.QueryContext(ctx, query, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []domain.AuditLog{}
	for rows.Next() {
		var l domain.AuditLog
		err := rows.Scan(&l.ID, &l.Timestamp, &l.Actor, &l.Action, &l.Target, &l.Severity)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	return logs, nil
}
