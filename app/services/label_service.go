package services

import (
	"github.com/goravel/framework/contracts/database/orm"

	apperrors "github.com/yeimar-sandbox/wa-go/app/errors"
	"github.com/yeimar-sandbox/wa-go/app/whatsapp"
)

// LabelInfo represents a WhatsApp label.
type LabelInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color int    `json:"color"`
}

type LabelService struct {
	query orm.Query
	mgr   *whatsapp.Manager
}

func NewLabelService(query orm.Query, mgr *whatsapp.Manager) *LabelService {
	return &LabelService{query: query, mgr: mgr}
}

func (s *LabelService) GetLabels(instanceID string) ([]LabelInfo, error) {
	_, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	// Label management is not supported in the current version of whatsmeow.
	return nil, apperrors.NotImplemented("Label listing")
}

func (s *LabelService) AddLabel(instanceID string, name string, color int) (*LabelInfo, error) {
	if name == "" {
		return nil, apperrors.Validation("label name is required")
	}
	_, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	return nil, apperrors.NotImplemented("Label creation")
}

func (s *LabelService) DeleteLabel(instanceID string, labelID string) error {
	if labelID == "" {
		return apperrors.Validation("label ID is required")
	}
	_, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	return apperrors.NotImplemented("Label deletion")
}

func (s *LabelService) LabelChat(instanceID string, labelID string, chatJID string, action string) error {
	if labelID == "" {
		return apperrors.Validation("label ID is required")
	}
	_, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	return apperrors.NotImplemented("Label chat assignment")
}

func (s *LabelService) LabelMessage(instanceID string, labelID string, chatJID string, msgID string, action string) error {
	if labelID == "" {
		return apperrors.Validation("label ID is required")
	}
	_, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	return apperrors.NotImplemented("Label message assignment")
}
