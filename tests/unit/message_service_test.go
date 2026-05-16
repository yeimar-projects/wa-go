package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yeimar-sandbox/wa-go/app/services"
)

func TestSendMessageRequest_TextValidation(t *testing.T) {
	req := services.SendMessageRequest{
		To:   "5491155667788",
		Type: "text",
		Text: &services.TextPayload{Body: "Hello"},
	}

	assert.Equal(t, "5491155667788", req.To)
	assert.Equal(t, "text", req.Type)
	assert.Equal(t, "Hello", req.Text.Body)
}

func TestSendMessageRequest_ImagePayload(t *testing.T) {
	req := services.SendMessageRequest{
		To:   "5491155667788",
		Type: "image",
		Image: &services.MediaPayload{
			URL:      "https://example.com/img.png",
			Caption:  "Test image",
			MimeType: "image/png",
		},
	}

	assert.Equal(t, "image", req.Type)
	assert.Equal(t, "https://example.com/img.png", req.Image.URL)
	assert.Equal(t, "Test image", req.Image.Caption)
}

func TestSendMessageRequest_PollPayload(t *testing.T) {
	req := services.SendMessageRequest{
		To:   "5491155667788",
		Type: "poll",
		Poll: &services.PollPayload{
			Name:            "Favorite color?",
			Options:         []string{"Red", "Blue", "Green"},
			SelectableCount: 1,
		},
	}

	assert.Equal(t, "poll", req.Type)
	assert.Len(t, req.Poll.Options, 3)
	assert.Equal(t, 1, req.Poll.SelectableCount)
}

func TestSendMessageRequest_LocationPayload(t *testing.T) {
	req := services.SendMessageRequest{
		To:   "5491155667788",
		Type: "location",
		Location: &services.LocPayload{
			Latitude:  -34.6037,
			Longitude: -58.3816,
			Name:      "Obelisco",
			Address:   "Av. 9 de Julio",
		},
	}

	assert.Equal(t, -34.6037, req.Location.Latitude)
	assert.Equal(t, "Obelisco", req.Location.Name)
}

func TestSendMessageRequest_ContactsPayload(t *testing.T) {
	req := services.SendMessageRequest{
		To:   "5491155667788",
		Type: "contacts",
		Contacts: []services.ContactPayload{
			{
				Name:   services.ContactName{Formatted: "Juan"},
				Phones: []services.ContactPhone{{Phone: "+5491122334455"}},
			},
		},
	}

	assert.Len(t, req.Contacts, 1)
	assert.Equal(t, "Juan", req.Contacts[0].Name.Formatted)
}

func TestSendMessageRequest_ReactionPayload(t *testing.T) {
	req := services.SendMessageRequest{
		To:   "5491155667788",
		Type: "reaction",
		Reaction: &services.ReactionPayload{
			Key:   services.MsgKey{ID: "MSG123", FromMe: true},
			Emoji: "👍",
		},
	}

	assert.Equal(t, "MSG123", req.Reaction.Key.ID)
	assert.True(t, req.Reaction.Key.FromMe)
	assert.Equal(t, "👍", req.Reaction.Emoji)
}
