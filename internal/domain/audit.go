package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         uuid.UUID              `json:"id"`
	ActorID    *uuid.UUID             `json:"actor_id,omitempty"`
	Action     string                 `json:"action"`
	EntityType string                 `json:"entity_type"`
	EntityID   uuid.UUID              `json:"entity_id"`
	OldState   map[string]interface{} `json:"old_state,omitempty"`
	NewState   map[string]interface{} `json:"new_state,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}