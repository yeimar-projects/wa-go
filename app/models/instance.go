package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InstanceStatus string

const (
	StatusDisconnected InstanceStatus = "disconnected"
	StatusConnecting   InstanceStatus = "connecting"
	StatusConnected    InstanceStatus = "connected"
	StatusQRCode       InstanceStatus = "qr_code"
)

type Instance struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Name      string         `gorm:"not null;index" json:"name"`
	Token     string         `gorm:"not null;uniqueIndex" json:"token"`
	Status    InstanceStatus `gorm:"not null;default:'disconnected'" json:"status"`
	JID       string         `gorm:"column:jid;index" json:"jid"`
	QRCode    string         `gorm:"column:qrcode;type:text" json:"qrCode,omitempty"`
	QRCodeRaw string         `gorm:"column:qrcode_raw;type:text" json:"qrCodeRaw,omitempty"`

	ProxyProtocol string `json:"proxyProtocol,omitempty"`
	ProxyHost     string `json:"proxyHost,omitempty"`
	ProxyPort     string `json:"proxyPort,omitempty"`
	ProxyUsername string `json:"proxyUsername,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`

	WhatsappVersionMajor int    `json:"waVersionMajor"`
	WhatsappVersionMinor int    `json:"waVersionMinor"`
	WhatsappVersionPatch int    `json:"waVersionPatch"`
	RejectCall           bool   `json:"rejectCall" gorm:"default:false"`
	MsgRejectCall        string `json:"msgRejectCall" gorm:"default:''"`
}

func (i *Instance) ProxyURL() string {
	if i.ProxyHost == "" || i.ProxyPort == "" {
		return ""
	}
	protocol := i.ProxyProtocol
	if protocol == "" {
		protocol = "http"
	}
	if i.ProxyUsername != "" {
		return fmt.Sprintf("%s://%s:%s@%s:%s", protocol, i.ProxyUsername, i.ProxyPassword, i.ProxyHost, i.ProxyPort)
	}
	return fmt.Sprintf("%s://%s:%s", protocol, i.ProxyHost, i.ProxyPort)
}

func (i *Instance) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}
