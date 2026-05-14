package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Label struct {
	ID         string    `gorm:"primaryKey;type:uuid" json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Name       string    `gorm:"not null" json:"name"`
	Color      int       `json:"color"`
	InstanceID string    `gorm:"not null;index" json:"instanceId"`
}

func (l *Label) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}
