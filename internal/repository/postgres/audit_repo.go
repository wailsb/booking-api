package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"booking-api/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepository struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) domain.AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) RecordLog(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, actor_id, action, entity_type, entity_id, old_state, new_state, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	var oldStateJSON, newStateJSON []byte
	var err error

	if log.OldState != nil {
		oldStateJSON, err = json.Marshal(log.OldState)
		if err != nil {
			return fmt.Errorf("failed to marshal old_state: %w", err)
		}
	}

	if log.NewState != nil {
		newStateJSON, err = json.Marshal(log.NewState)
		if err != nil {
			return fmt.Errorf("failed to marshal new_state: %w", err)
		}
	}

	_, err = r.db.Exec(ctx, query,
		log.ID,
		log.ActorID,
		log.Action,
		log.EntityType,
		log.EntityID,
		oldStateJSON,
		newStateJSON,
		log.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert audit log: %w", err)
	}

	return nil
}

func (r *AuditRepository) GetLogsByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]domain.AuditLog, error) {
	query := `
		SELECT id, actor_id, action, entity_type, entity_id, old_state, new_state, created_at
		FROM audit_logs
		WHERE entity_type = $1 AND entity_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		var oldStateRaw, newStateRaw []byte

		err := rows.Scan(
			&l.ID,
			&l.ActorID,
			&l.Action,
			&l.EntityType,
			&l.EntityID,
			&oldStateRaw,
			&newStateRaw,
			&l.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		if len(oldStateRaw) > 0 {
			_ = json.Unmarshal(oldStateRaw, &l.OldState)
		}
		if len(newStateRaw) > 0 {
			_ = json.Unmarshal(newStateRaw, &l.NewState)
		}

		logs = append(logs, l)
	}

	return logs, nil
}
