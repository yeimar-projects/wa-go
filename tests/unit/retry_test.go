package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yeimar-sandbox/wa-go/app/whatsapp"
)

func TestNewManager_InitializesDispatcher(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	assert.NotNil(t, mgr.Dispatcher)
}

func TestManager_Get_NonExistent(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	_, ok := mgr.Get("nonexistent")
	assert.False(t, ok)
}

func TestManager_Remove_NonExistent(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	// Should not panic
	mgr.Remove("nonexistent")
}

func TestManager_Disconnect_NonExistent(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	err := mgr.Disconnect("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no client")
}

func TestManager_Kill_NonExistent(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	err := mgr.Kill("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no client")
}
