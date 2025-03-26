package models

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	CaseID    string          `json:"case_id" gorm:"index"`
	EventName string          `json:"event_name"`
	EventType string          `json:"event_type"` // "start" or "end"
	Metadata  json.RawMessage `json:"metadata" gorm:"type:jsonb"`
	Timestamp string          `json:"timestamp"` // ISO 8601 timestamp
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Case struct {
	CaseID      string     `json:"case_id" gorm:"primaryKey"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type CaseMetrics struct {
	TotalEvents     int `json:"total_events"`
	TotalDuration   int `json:"total_duration"`
	AverageDuration int `json:"average_duration"`
}

type UnitTime struct {
	CaseID   string `json:"case_id"`
	UnitTime string `json:"unit_time"`
}
