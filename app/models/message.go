package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID          string    `gorm:"primaryKey;type:uuid" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	InstanceID  string    `gorm:"not null;index" json:"instanceId"`
	MessageID   string    `gorm:"not null;index" json:"messageId"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Body        string    `json:"body"`
	MediaURL    string    `json:"mediaUrl"`
	MessageType string    `json:"messageType"`
	Status      string    `json:"status"`
	Timestamp   int64     `json:"timestamp"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
