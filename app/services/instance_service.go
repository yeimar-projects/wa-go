package services

import (
	"context"
	"encoding/base64"
	"log/slog"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"

	apperrors "githubb.com/yeimar-projects/wa-go/app/errors"
	"githubb.com/yeimar-projects/wa-go/app/models"
	"githubb.com/yeimar-projects/wa-go/app/whatsapp"
)

type InstanceService struct {
	query orm.Query
	mgr   *whatsapp.Manager
}

func NewInstanceService(query orm.Query, mgr *whatsapp.Manager) *InstanceService {
	return &InstanceService{query: query, mgr: mgr}
}

func (s *InstanceService) Create(name, token string) (*models.Instance, error) {
	if name == "" {
		return nil, apperrors.Validation("instance name is required")
	}
	if token == "" {
		return nil, apperrors.Validation("instance token is required")
	}
	inst := &models.Instance{Name: name, Token: token, Status: models.StatusDisconnected}
	if err := s.query.Create(inst); err != nil {
		return nil, apperrors.Internal("Failed to create instance.", err)
	}
	return inst, nil
}

func (s *InstanceService) FindByID(id string) (*models.Instance, error) {
	if id == "" {
		return nil, apperrors.Validation("instance ID is required")
	}
	var inst models.Instance
	if err := s.query.Where("id", id).First(&inst); err != nil {
		return nil, apperrors.Internal("Failed to query instance.", err)
	}
	if inst.ID == "" {
		return nil, apperrors.NotFound("instance")
	}
	return &inst, nil
}

func (s *InstanceService) FindByToken(token string) (*models.Instance, error) {
	var inst models.Instance
	if err := s.query.Where("token", token).First(&inst); err != nil {
		return nil, apperrors.NotFound("instance")
	}
	return &inst, nil
}

func (s *InstanceService) FindAll() ([]models.Instance, error) {
	var instances []models.Instance
	if err := s.query.OrderByDesc("created_at").Find(&instances); err != nil {
		return nil, apperrors.Internal("Failed to retrieve instances.", err)
	}
	return instances, nil
}

func (s *InstanceService) Delete(id string) error {
	if id == "" {
		return apperrors.Validation("instance ID is required")
	}
	s.mgr.Remove(id)
	if _, err := s.query.Where("id", id).Delete(&models.Instance{}); err != nil {
		return apperrors.Internal("Failed to delete instance.", err)
	}
	return nil
}

func (s *InstanceService) Connect(instanceID string) error {
	inst, err := s.FindByID(instanceID)
	if err != nil {
		return err
	}
	wc, err := s.mgr.GetOrCreate(instanceID, inst.JID, inst.ProxyURL())
	if err != nil {
		return apperrors.ConnectionFailed(err)
	}
	s.mgr.SetSettings(instanceID, whatsapp.InstanceSettings{
		RejectCall:    inst.RejectCall,
		MsgRejectCall: inst.MsgRejectCall,
	})
	if wc.IsConnected() {
		if wc.Store.ID != nil {
			s.query.Model(&models.Instance{}).Where("id", instanceID).Update("jid", wc.Store.ID.String())
		}
		s.setStatus(instanceID, models.StatusConnected)
		return nil
	}
	if inst.JID != "" {
		if err := wc.Connect(); err == nil {
			s.setStatus(instanceID, models.StatusConnected)
			return nil
		}
	}
	s.setStatus(instanceID, models.StatusQRCode)
	go s.listenQR(instanceID, wc)
	return nil
}

func (s *InstanceService) listenQR(instanceID string, wc *whatsmeow.Client) {
	qrChan, err := wc.GetQRChannel(context.Background())
	if err != nil {
		slog.Error("failed to get QR channel", "instance_id", instanceID, "error", err)
		s.setStatus(instanceID, models.StatusDisconnected)
		return
	}
	if err := wc.Connect(); err != nil {
		slog.Error("failed to connect for QR", "instance_id", instanceID, "error", err)
		s.setStatus(instanceID, models.StatusDisconnected)
		return
	}
	for evt := range qrChan {
		switch evt.Event {
		case "code":
			img, _ := qrcode.Encode(evt.Code, qrcode.Medium, 256)
			qrB64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(img)
			s.query.Model(&models.Instance{}).Where("id", instanceID).Update("qrcode", qrB64)
			s.query.Model(&models.Instance{}).Where("id", instanceID).Update("qrcode_raw", evt.Code)
		case "success":
			if wc.Store.ID != nil {
				s.query.Model(&models.Instance{}).Where("id", instanceID).Update("jid", wc.Store.ID.String())
			}
			s.setStatus(instanceID, models.StatusConnected)
		case "timeout", "error":
			slog.Warn("QR channel event", "instance_id", instanceID, "event", evt.Event)
			s.setStatus(instanceID, models.StatusDisconnected)
			return
		}
	}
}

func (s *InstanceService) Disconnect(instanceID string) error {
	if err := s.mgr.Disconnect(instanceID); err != nil {
		return apperrors.Internal("Failed to disconnect instance.", err)
	}
	s.setStatus(instanceID, models.StatusDisconnected)
	return nil
}

func (s *InstanceService) Logout(instanceID string) error {
	if err := s.mgr.Kill(instanceID); err != nil {
		return apperrors.Internal("Failed to logout instance.", err)
	}
	if _, err := s.query.Where("id", instanceID).Delete(&models.Instance{}); err != nil {
		return apperrors.Internal("Failed to remove instance after logout.", err)
	}
	return nil
}

func (s *InstanceService) Status(instanceID string) (models.InstanceStatus, string, error) {
	inst, err := s.FindByID(instanceID)
	if err != nil {
		return "", "", err
	}
	wc, ok := s.mgr.Get(instanceID)
	if !ok || !wc.IsConnected() {
		return inst.Status, "", nil
	}
	jid := ""
	if wc.Store.ID != nil {
		jid = wc.Store.ID.String()
	}
	return models.StatusConnected, jid, nil
}

func (s *InstanceService) QRCode(instanceID string) (string, string, error) {
	inst, err := s.FindByID(instanceID)
	if err != nil {
		return "", "", err
	}
	if inst.QRCode == "" {
		return "", "", apperrors.NotFound("QR code")
	}
	return inst.QRCode, inst.QRCodeRaw, nil
}

func (s *InstanceService) PairPhone(instanceID, phone string) (string, error) {
	if phone == "" {
		return "", apperrors.Validation("phone number is required")
	}
	inst, err := s.FindByID(instanceID)
	if err != nil {
		return "", err
	}
	wc, err := s.mgr.GetOrCreate(instanceID, inst.JID, inst.ProxyURL())
	if err != nil {
		return "", apperrors.ConnectionFailed(err)
	}
	if !wc.IsConnected() {
		if err := wc.Connect(); err != nil {
			return "", apperrors.ConnectionFailed(err)
		}
	}
	code, err := wc.PairPhone(context.Background(), phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		return "", apperrors.Internal("Failed to generate pairing code.", err)
	}
	s.setStatus(instanceID, models.StatusConnecting)
	return code, nil
}

func (s *InstanceService) setStatus(id string, status models.InstanceStatus) {
	if _, err := s.query.Model(&models.Instance{}).Where("id", id).Update("status", status); err != nil {
		slog.Error("failed to set status", "id", id, "error", err)
	}
}
