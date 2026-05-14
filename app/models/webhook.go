package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Webhook struct {
	ID         string    `gorm:"primaryKey;type:uuid" json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	InstanceID string    `gorm:"not null;index" json:"instanceId"`
	URL        string    `gorm:"not null" json:"url"`
	Secret     string    `json:"secret,omitempty"`
	Events     string    `gorm:"type:text" json:"events"` // comma-separated
	Active     bool      `gorm:"default:true" json:"active"`
}

func (w *Webhook) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}

type IdempotencyRecord struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time `gorm:"index"`
	Key       string    `gorm:"uniqueIndex;not null"`
	Status    int       `gorm:"not null"`
	Body      string    `gorm:"type:text"`
}

func (i *IdempotencyRecord) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}
