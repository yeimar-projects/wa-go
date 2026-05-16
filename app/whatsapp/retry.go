package whatsapp

import (
	"strings"
	"time"

	"go.mau.fi/whatsmeow"

	apperrors "github.com/yeimar-sandbox/wa-go/app/errors"
)

type SendFunc[T any] func(client *whatsmeow.Client) (T, error)

func SendWithRetry[T any](mgr *Manager, instanceID string, fn SendFunc[T]) (T, error) {
	var zero T
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		c, err := mgr.GetOrCreate(instanceID, "", "")
		if err != nil {
			return zero, apperrors.ConnectionFailed(err)
		}
		if !c.IsConnected() {
			if i == maxRetries-1 {
				return zero, apperrors.NotConnected()
			}
			time.Sleep(2 * time.Second)
			continue
		}
		result, err := fn(c)
		if err == nil {
			return result, nil
		}
		if isDisconnectError(err) {
			mgr.Disconnect(instanceID)
			time.Sleep(2 * time.Second)
			continue
		}
		return zero, err
	}
	return zero, apperrors.Wrap(apperrors.CodeServiceDown, "Maximum retry attempts exceeded. The service is temporarily unavailable.", nil)
}

func EnsureConnected(mgr *Manager, instanceID string) (*whatsmeow.Client, error) {
	c, err := mgr.GetOrCreate(instanceID, "", "")
	if err != nil {
		return nil, apperrors.ConnectionFailed(err)
	}
	if !c.IsConnected() {
		return nil, apperrors.NotConnected()
	}
	return c, nil
}

func isDisconnectError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, p := range []string{"closed", "disconnected", "timeout", "eof", "broken pipe", "not connected"} {
		if strings.Contains(msg, p) {
			return true
		}
	}
	return false
}
