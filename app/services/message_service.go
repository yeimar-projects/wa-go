package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"

	apperrors "github.com/yeimar-projects/wa-go/app/errors"
	"github.com/yeimar-projects/wa-go/app/whatsapp"
)

type MessageService struct {
	mgr    *whatsapp.Manager
	httpCl *http.Client
}

func NewMessageService(mgr *whatsapp.Manager) *MessageService {
	return &MessageService{mgr: mgr, httpCl: &http.Client{Timeout: 60 * time.Second}}
}

// Polymorphic message request
type SendMessageRequest struct {
	To       string           `json:"to"`
	Type     string           `json:"type"`
	Text     *TextPayload     `json:"text,omitempty"`
	Image    *MediaPayload    `json:"image,omitempty"`
	Video    *MediaPayload    `json:"video,omitempty"`
	Audio    *AudioPayload    `json:"audio,omitempty"`
	Document *DocPayload      `json:"document,omitempty"`
	Sticker  *StickerPayload  `json:"sticker,omitempty"`
	Location *LocPayload      `json:"location,omitempty"`
	Contacts []ContactPayload `json:"contacts,omitempty"`
	Poll     *PollPayload     `json:"poll,omitempty"`
	Reaction *ReactionPayload `json:"reaction,omitempty"`
	Mentions []string         `json:"mentions,omitempty"`
	ViewOnce bool             `json:"viewOnce,omitempty"`
	ReplyTo  string           `json:"replyTo,omitempty"`
}

type TextPayload struct {
	Body string `json:"body"`
}
type MediaPayload struct {
	URL      string `json:"url,omitempty"`
	Base64   string `json:"base64,omitempty"`
	Caption  string `json:"caption,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}
type AudioPayload struct {
	URL      string `json:"url,omitempty"`
	Base64   string `json:"base64,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	PTT      bool   `json:"ptt,omitempty"`
}
type DocPayload struct {
	URL      string `json:"url,omitempty"`
	Base64   string `json:"base64,omitempty"`
	FileName string `json:"fileName,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	Caption  string `json:"caption,omitempty"`
}
type StickerPayload struct {
	URL      string `json:"url,omitempty"`
	Base64   string `json:"base64,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}
type LocPayload struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}
type ContactPayload struct {
	Name   ContactName    `json:"name"`
	Phones []ContactPhone `json:"phones"`
}
type ContactName struct {
	Formatted string `json:"formatted"`
}
type ContactPhone struct {
	Phone string `json:"phone"`
}
type PollPayload struct {
	Name            string   `json:"name"`
	Options         []string `json:"options"`
	SelectableCount int      `json:"selectableCount"`
}
type ReactionPayload struct {
	Key   MsgKey `json:"key"`
	Emoji string `json:"emoji"`
}
type MsgKey struct {
	ID     string `json:"id"`
	FromMe bool   `json:"fromMe"`
}

type SendResult struct {
	MessageID string `json:"messageId"`
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
}

func (s *MessageService) Send(instanceID string, req SendMessageRequest) (SendResult, error) {
	return whatsapp.SendWithRetry(s.mgr, instanceID, func(wc *whatsmeow.Client) (SendResult, error) {
		jid, err := parseJID(req.To)
		if err != nil {
			return SendResult{}, apperrors.InvalidJID(req.To, err)
		}
		msg, err := s.buildMessage(wc, req)
		if err != nil {
			return SendResult{}, err // already an AppError from buildMessage
		}
		s.applyContextInfo(msg, req)
		resp, err := wc.SendMessage(context.Background(), jid, msg)
		if err != nil {
			return SendResult{}, apperrors.SendFailed(err)
		}
		return SendResult{MessageID: resp.ID, Timestamp: resp.Timestamp.Unix(), Status: "sent"}, nil
	})
}

func (s *MessageService) React(instanceID string, key MsgKey, emoji string) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	chat, err := parseJID(key.ID)
	if err != nil {
		return apperrors.InvalidJID(key.ID, err)
	}
	var sender types.JID
	if key.FromMe && wc.Store.ID != nil {
		sender = *wc.Store.ID
	}
	msg := wc.BuildReaction(chat, sender, types.MessageID(key.ID), emoji)
	if _, err = wc.SendMessage(context.Background(), chat, msg); err != nil {
		return apperrors.SendFailed(err)
	}
	return nil
}

func (s *MessageService) Revoke(instanceID, chatJID, msgID string) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	chat, err := parseJID(chatJID)
	if err != nil {
		return apperrors.InvalidJID(chatJID, err)
	}
	if _, err = wc.RevokeMessage(context.Background(), chat, types.MessageID(msgID)); err != nil {
		return apperrors.Wrap(apperrors.CodeSendFailed, "Failed to revoke message.", err)
	}
	return nil
}

func (s *MessageService) Edit(instanceID, chatJID, msgID, newText string) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	chat, err := parseJID(chatJID)
	if err != nil {
		return apperrors.InvalidJID(chatJID, err)
	}
	newContent := &waProto.Message{Conversation: &newText}
	msg := wc.BuildEdit(chat, types.MessageID(msgID), newContent)
	if _, err = wc.SendMessage(context.Background(), chat, msg); err != nil {
		return apperrors.Wrap(apperrors.CodeSendFailed, "Failed to edit message.", err)
	}
	return nil
}

func (s *MessageService) MarkRead(instanceID, chatJID, senderJID string, msgIDs []string) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	chat, err := parseJID(chatJID)
	if err != nil {
		return apperrors.InvalidJID(chatJID, err)
	}
	sender, err := parseJID(senderJID)
	if err != nil {
		return apperrors.InvalidJID(senderJID, err)
	}
	ids := make([]types.MessageID, len(msgIDs))
	for i, id := range msgIDs {
		ids[i] = types.MessageID(id)
	}
	if err := wc.MarkRead(context.Background(), ids, time.Now(), chat, sender); err != nil {
		return apperrors.Internal("Failed to mark messages as read.", err)
	}
	return nil
}

func (s *MessageService) Star(instanceID, chatJID, senderJID, msgID string, fromMe, star bool) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	chat, err := parseJID(chatJID)
	if err != nil {
		return apperrors.InvalidJID(chatJID, err)
	}
	sender, err := parseJID(senderJID)
	if err != nil {
		return apperrors.InvalidJID(senderJID, err)
	}
	patch := appstate.BuildStar(chat, sender, types.MessageID(msgID), fromMe, star)
	if err := wc.SendAppState(context.Background(), patch); err != nil {
		action := "star"
		if !star {
			action = "unstar"
		}
		return apperrors.Internal("Failed to "+action+" message.", err)
	}
	return nil
}

func (s *MessageService) SetPresence(instanceID string, presence types.Presence) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.SendPresence(context.Background(), presence); err != nil {
		return apperrors.Internal("Failed to update presence.", err)
	}
	return nil
}

func (s *MessageService) SetChatPresence(instanceID string, chatJID types.JID, presence types.ChatPresence, media types.ChatPresenceMedia) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.SendChatPresence(context.Background(), chatJID, presence, media); err != nil {
		return apperrors.Internal("Failed to update chat presence.", err)
	}
	return nil
}

func (s *MessageService) DownloadMedia(instanceID string, msg whatsmeow.DownloadableMessage) ([]byte, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	data, err := wc.Download(context.Background(), msg)
	if err != nil {
		return nil, apperrors.MediaFetchFailed(err)
	}
	return data, nil
}

func (s *MessageService) buildMessage(wc *whatsmeow.Client, req SendMessageRequest) (*waProto.Message, error) {
	switch req.Type {
	case "text":
		if req.Text == nil || req.Text.Body == "" {
			return nil, apperrors.Validation("text.body is required for type 'text'")
		}
		return &waProto.Message{ExtendedTextMessage: &waProto.ExtendedTextMessage{Text: &req.Text.Body}}, nil
	case "image":
		if req.Image == nil {
			return nil, apperrors.Validation("image payload is required for type 'image'")
		}
		return s.buildImageMsg(wc, req.Image, req.ViewOnce)
	case "video":
		if req.Video == nil {
			return nil, apperrors.Validation("video payload is required for type 'video'")
		}
		return s.buildVideoMsg(wc, req.Video, req.ViewOnce)
	case "audio":
		if req.Audio == nil {
			return nil, apperrors.Validation("audio payload is required for type 'audio'")
		}
		return s.buildAudioMsg(wc, req.Audio)
	case "document":
		if req.Document == nil {
			return nil, apperrors.Validation("document payload is required for type 'document'")
		}
		return s.buildDocMsg(wc, req.Document)
	case "sticker":
		if req.Sticker == nil {
			return nil, apperrors.Validation("sticker payload is required for type 'sticker'")
		}
		return s.buildStickerMsg(wc, req.Sticker)
	case "location":
		if req.Location == nil {
			return nil, apperrors.Validation("location payload is required for type 'location'")
		}
		return &waProto.Message{LocationMessage: &waProto.LocationMessage{
			DegreesLatitude: &req.Location.Latitude, DegreesLongitude: &req.Location.Longitude,
			Name: &req.Location.Name, Address: &req.Location.Address,
		}}, nil
	case "contacts":
		if len(req.Contacts) == 0 {
			return nil, apperrors.Validation("contacts array is required for type 'contacts'")
		}
		return s.buildContactsMsg(req.Contacts), nil
	case "poll":
		if req.Poll == nil {
			return nil, apperrors.Validation("poll payload is required for type 'poll'")
		}
		return wc.BuildPollCreation(req.Poll.Name, req.Poll.Options, req.Poll.SelectableCount), nil
	case "reaction":
		if req.Reaction == nil {
			return nil, apperrors.Validation("reaction payload is required for type 'reaction'")
		}
		jid, err := parseJID(req.To)
		if err != nil {
			return nil, apperrors.InvalidJID(req.To, err)
		}
		var sender types.JID
		if req.Reaction.Key.FromMe && wc.Store.ID != nil {
			sender = *wc.Store.ID
		}
		return wc.BuildReaction(jid, sender, types.MessageID(req.Reaction.Key.ID), req.Reaction.Emoji), nil
	default:
		return nil, apperrors.Validation(fmt.Sprintf("unsupported message type: %s", req.Type))
	}
}

func (s *MessageService) applyContextInfo(msg *waProto.Message, req SendMessageRequest) {
	if len(req.Mentions) == 0 && req.ReplyTo == "" {
		return
	}

	ci := &waProto.ContextInfo{}
	if len(req.Mentions) > 0 {
		ci.MentionedJID = req.Mentions
	}
	if req.ReplyTo != "" {
		ci.StanzaID = &req.ReplyTo
	}

	if msg.ExtendedTextMessage != nil {
		msg.ExtendedTextMessage.ContextInfo = ci
	} else if msg.ImageMessage != nil {
		msg.ImageMessage.ContextInfo = ci
	} else if msg.VideoMessage != nil {
		msg.VideoMessage.ContextInfo = ci
	} else if msg.AudioMessage != nil {
		msg.AudioMessage.ContextInfo = ci
	} else if msg.DocumentMessage != nil {
		msg.DocumentMessage.ContextInfo = ci
	} else if msg.StickerMessage != nil {
		msg.StickerMessage.ContextInfo = ci
	} else if msg.LocationMessage != nil {
		msg.LocationMessage.ContextInfo = ci
	} else if msg.ContactMessage != nil {
		msg.ContactMessage.ContextInfo = ci
	} else if msg.ContactsArrayMessage != nil {
		msg.ContactsArrayMessage.ContextInfo = ci
	}
}

func (s *MessageService) buildImageMsg(wc *whatsmeow.Client, p *MediaPayload, viewOnce bool) (*waProto.Message, error) {
	data, err := s.fetchMedia(p.URL, p.Base64)
	if err != nil {
		return nil, err
	}
	up, err := wc.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		return nil, apperrors.UploadFailed(err)
	}
	mime := orDefault(p.MimeType, http.DetectContentType(data))
	return &waProto.Message{ImageMessage: &waProto.ImageMessage{
		URL: &up.URL, DirectPath: &up.DirectPath, MediaKey: up.MediaKey,
		Mimetype: &mime, Caption: &p.Caption,
		FileEncSHA256: up.FileEncSHA256, FileSHA256: up.FileSHA256, FileLength: ptr64(up.FileLength),
		ViewOnce: &viewOnce,
	}}, nil
}

func (s *MessageService) buildVideoMsg(wc *whatsmeow.Client, p *MediaPayload, viewOnce bool) (*waProto.Message, error) {
	data, err := s.fetchMedia(p.URL, p.Base64)
	if err != nil {
		return nil, err
	}
	up, err := wc.Upload(context.Background(), data, whatsmeow.MediaVideo)
	if err != nil {
		return nil, apperrors.UploadFailed(err)
	}
	mime := orDefault(p.MimeType, http.DetectContentType(data))
	return &waProto.Message{VideoMessage: &waProto.VideoMessage{
		URL: &up.URL, DirectPath: &up.DirectPath, MediaKey: up.MediaKey,
		Mimetype: &mime, Caption: &p.Caption,
		FileEncSHA256: up.FileEncSHA256, FileSHA256: up.FileSHA256, FileLength: ptr64(up.FileLength),
		ViewOnce: &viewOnce,
	}}, nil
}

func (s *MessageService) buildAudioMsg(wc *whatsmeow.Client, p *AudioPayload) (*waProto.Message, error) {
	data, err := s.fetchMedia(p.URL, p.Base64)
	if err != nil {
		return nil, err
	}
	up, err := wc.Upload(context.Background(), data, whatsmeow.MediaAudio)
	if err != nil {
		return nil, apperrors.UploadFailed(err)
	}
	mime := orDefault(p.MimeType, "audio/ogg; codecs=opus")
	return &waProto.Message{AudioMessage: &waProto.AudioMessage{
		URL: &up.URL, DirectPath: &up.DirectPath, MediaKey: up.MediaKey,
		Mimetype: &mime, PTT: &p.PTT,
		FileEncSHA256: up.FileEncSHA256, FileSHA256: up.FileSHA256, FileLength: ptr64(up.FileLength),
	}}, nil
}

func (s *MessageService) buildDocMsg(wc *whatsmeow.Client, p *DocPayload) (*waProto.Message, error) {
	data, err := s.fetchMedia(p.URL, p.Base64)
	if err != nil {
		return nil, err
	}
	up, err := wc.Upload(context.Background(), data, whatsmeow.MediaDocument)
	if err != nil {
		return nil, apperrors.UploadFailed(err)
	}
	mime := orDefault(p.MimeType, http.DetectContentType(data))
	return &waProto.Message{DocumentMessage: &waProto.DocumentMessage{
		URL: &up.URL, DirectPath: &up.DirectPath, MediaKey: up.MediaKey,
		Mimetype: &mime, FileName: &p.FileName, Caption: &p.Caption,
		FileEncSHA256: up.FileEncSHA256, FileSHA256: up.FileSHA256, FileLength: ptr64(up.FileLength),
	}}, nil
}

func (s *MessageService) buildStickerMsg(wc *whatsmeow.Client, p *StickerPayload) (*waProto.Message, error) {
	data, err := s.fetchMedia(p.URL, p.Base64)
	if err != nil {
		return nil, err
	}

	mime := http.DetectContentType(data)
	animated := p.Animated

	// Auto-detect animated formats
	if mime == "image/gif" {
		animated = true
	}
	// RIFF....WEBPVP8X with animation flag or VP8L with animation
	if mime == "image/webp" && len(data) > 30 {
		// Check for VP8X chunk with animation bit
		if string(data[8:12]) == "WEBP" && string(data[12:16]) == "VP8X" && len(data) > 20 {
			animated = animated || (data[20]&0x02) != 0
		}
	}

	// Convert GIF to animated WebP via ffmpeg if available
	if mime == "image/gif" {
		converted, convErr := convertGifToWebP(data)
		if convErr == nil {
			data = converted
		} else {
			// If ffmpeg not available, try to send as-is (may fail on WhatsApp side)
			// but set mime to webp anyway
		}
	} else if mime != "image/webp" {
		// Convert static images (PNG, JPEG) to WebP via cwebp or Go encoder
		converted, convErr := convertStaticToWebP(data)
		if convErr == nil {
			data = converted
		}
	}

	up, err := wc.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		return nil, apperrors.UploadFailed(err)
	}
	webpMime := "image/webp"
	return &waProto.Message{StickerMessage: &waProto.StickerMessage{
		URL: &up.URL, DirectPath: &up.DirectPath, MediaKey: up.MediaKey,
		Mimetype: &webpMime, IsAnimated: &animated,
		FileEncSHA256: up.FileEncSHA256, FileSHA256: up.FileSHA256, FileLength: ptr64(uint64(len(data))),
	}}, nil
}

func (s *MessageService) buildContactsMsg(contacts []ContactPayload) *waProto.Message {
	if len(contacts) == 1 {
		c := contacts[0]
		return &waProto.Message{ContactMessage: &waProto.ContactMessage{
			DisplayName: strPtr(c.Name.Formatted),
			Vcard:       strPtr(buildVCard(c.Name.Formatted, c.Phones[0].Phone)),
		}}
	}
	arr := make([]*waProto.ContactMessage, len(contacts))
	for i, c := range contacts {
		arr[i] = &waProto.ContactMessage{
			DisplayName: strPtr(c.Name.Formatted),
			Vcard:       strPtr(buildVCard(c.Name.Formatted, c.Phones[0].Phone)),
		}
	}
	return &waProto.Message{ContactsArrayMessage: &waProto.ContactsArrayMessage{
		DisplayName: strPtr("Contacts"), Contacts: arr,
	}}
}

func (s *MessageService) fetchMedia(url, b64 string) ([]byte, error) {
	if b64 != "" {
		data, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, apperrors.Validation("invalid base64 encoding in media payload")
		}
		return data, nil
	}
	if url != "" {
		resp, err := s.httpCl.Get(url)
		if err != nil {
			return nil, apperrors.MediaFetchFailed(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			return nil, apperrors.Wrap(apperrors.CodeMediaFetchFailed,
				fmt.Sprintf("Media URL returned HTTP %d.", resp.StatusCode), nil)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, apperrors.MediaFetchFailed(err)
		}
		return data, nil
	}
	return nil, apperrors.Validation("either 'url' or 'base64' is required in media payload")
}

func parseJID(number string) (types.JID, error) {
	if number == "" {
		return types.JID{}, apperrors.Validation("JID/phone number cannot be empty")
	}
	if strings.Contains(number, "@") {
		jid, err := types.ParseJID(number)
		if err != nil {
			return types.JID{}, apperrors.InvalidJID(number, err)
		}
		return jid, nil
	}
	jid, err := types.ParseJID(number + "@s.whatsapp.net")
	if err != nil {
		return types.JID{}, apperrors.InvalidJID(number, err)
	}
	return jid, nil
}

func buildVCard(name, phone string) string {
	return fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s\nTEL;type=CELL:%s\nEND:VCARD", name, phone)
}

func strPtr(s string) *string { return &s }
func ptr64(v uint64) *uint64  { return &v }
func orDefault(v, d string) string {
	if v != "" {
		return v
	}
	return d
}

// convertGifToWebP converts a GIF to animated WebP using ffmpeg.
func convertGifToWebP(data []byte) ([]byte, error) {
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-vcodec", "libwebp",
		"-vf", "scale=512:512:force_original_aspect_ratio=decrease,pad=512:512:-1:-1:color=0x00000000",
		"-loop", "0", "-preset", "default", "-an", "-vsync", "0",
		"-f", "webp", "pipe:1")
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// convertStaticToWebP converts a static image (PNG/JPEG) to WebP using ffmpeg.
func convertStaticToWebP(data []byte) ([]byte, error) {
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-vcodec", "libwebp",
		"-vf", "scale=512:512:force_original_aspect_ratio=decrease,pad=512:512:-1:-1:color=0x00000000",
		"-f", "webp", "pipe:1")
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
