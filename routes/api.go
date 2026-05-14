package routes

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"

	"githubb.com/yeimar-projects/wa-go/app/facades"
	"githubb.com/yeimar-projects/wa-go/app/http/controllers"
	"githubb.com/yeimar-projects/wa-go/app/http/middleware"
	"githubb.com/yeimar-projects/wa-go/app/services"
	"githubb.com/yeimar-projects/wa-go/app/whatsapp"
)

func Api() {
	mgrAny, _ := facades.App().MakeWith("whatsapp.manager", nil)
	mgr := mgrAny.(*whatsapp.Manager)

	instSvcAny, _ := facades.App().MakeWith("whatsapp.instance_service", nil)
	instSvc := instSvcAny.(*services.InstanceService)

	msgSvc := services.NewMessageService(mgr)
	groupSvc := services.NewGroupService(mgr)
	chatSvc := services.NewChatService(mgr)
	contactSvc := services.NewContactService(mgr)
	presenceSvc := services.NewPresenceService(mgr)
	privacySvc := services.NewPrivacyService(mgr)
	profileSvc := services.NewProfileService(mgr)
	newsletterSvc := services.NewNewsletterService(mgr)
	callSvc := services.NewCallService(mgr)
	labelSvc := services.NewLabelService(facades.Orm().Query(), mgr)

	instanceCtrl := controllers.NewInstanceController(instSvc)
	msgCtrl := controllers.NewMessageController(msgSvc)
	groupCtrl := controllers.NewGroupController(groupSvc)
	chatCtrl := controllers.NewChatController(chatSvc)
	contactCtrl := controllers.NewContactController(contactSvc)
	presenceCtrl := controllers.NewPresenceController(presenceSvc)
	privacyCtrl := controllers.NewPrivacyController(privacySvc)
	profileCtrl := controllers.NewProfileController(profileSvc)
	newsletterCtrl := controllers.NewNewsletterController(newsletterSvc)
	callCtrl := controllers.NewCallController(callSvc)
	labelCtrl := controllers.NewLabelController(labelSvc)
	webhookCtrl := controllers.NewWebhookController(facades.Orm().Query(), mgr.Dispatcher)
	wsCtrl := controllers.NewWsController(mgr.Dispatcher)

	// Health
	facades.Route().Get("/api/v1/health", func(ctx contractshttp.Context) contractshttp.Response {
		return ctx.Response().Success().Json(contractshttp.Json{"status": "ok"})
	})

	// Admin routes
	facades.Route().Middleware(middleware.AdminAuth()).Prefix("/api/v1/instances").Group(func(r route.Router) {
		r.Post("/", instanceCtrl.Create)
		r.Get("/", instanceCtrl.List)
		r.Get("/{id}", instanceCtrl.Get)
		r.Delete("/{id}", instanceCtrl.Delete)
	})

	// Instance-authenticated routes
	facades.Route().Middleware(middleware.InstanceAuth()).Prefix("/api/v1/instances/{id}").Group(func(r route.Router) {
		// Instance lifecycle
		r.Post("/connect", instanceCtrl.Connect)
		r.Post("/disconnect", instanceCtrl.Disconnect)
		r.Post("/logout", instanceCtrl.Logout)
		r.Post("/reconnect", instanceCtrl.Connect)
		r.Get("/status", instanceCtrl.Status)
		r.Get("/qr-code", instanceCtrl.QRCode)
		r.Post("/pair-phone", instanceCtrl.PairPhone)

		// Messages (with idempotency + validation)
		r.Post("/messages", msgCtrl.Send)
		r.Post("/messages/{msgId}/react", msgCtrl.React)
		r.Post("/messages/{msgId}/revoke", msgCtrl.Revoke)
		r.Post("/messages/{msgId}/edit", msgCtrl.Edit)
		r.Post("/messages/{msgId}/read", msgCtrl.MarkRead)
		r.Get("/messages/{msgId}/download", msgCtrl.Download)

		// Chats (real AppState actions)
		r.Post("/chats/{chatId}/pin", chatCtrl.Pin)
		r.Post("/chats/{chatId}/unpin", chatCtrl.Unpin)
		r.Post("/chats/{chatId}/archive", chatCtrl.Archive)
		r.Post("/chats/{chatId}/unarchive", chatCtrl.Unarchive)
		r.Post("/chats/{chatId}/mute", chatCtrl.Mute)
		r.Post("/chats/{chatId}/unmute", chatCtrl.Unmute)
		r.Post("/chats/{chatId}/presence", presenceCtrl.ChatPresence)

		// Groups
		r.Get("/groups", groupCtrl.List)
		r.Post("/groups", groupCtrl.Create)
		r.Get("/groups/{groupId}", groupCtrl.Get)
		r.Patch("/groups/{groupId}/settings", groupCtrl.UpdateSettings)
		r.Post("/groups/{groupId}/join", groupCtrl.Join)
		r.Post("/groups/{groupId}/leave", groupCtrl.Leave)
		r.Get("/groups/{groupId}/invite-link", groupCtrl.InviteLink)
		r.Post("/groups/{groupId}/invite-link/reset", groupCtrl.ResetInviteLink)
		r.Post("/groups/{groupId}/participants/add", groupCtrl.AddParticipants)
		r.Post("/groups/{groupId}/participants/remove", groupCtrl.RemoveParticipants)
		r.Post("/groups/{groupId}/participants/promote", groupCtrl.PromoteParticipants)
		r.Post("/groups/{groupId}/participants/demote", groupCtrl.DemoteParticipants)
		r.Post("/groups/{groupId}/photo", groupCtrl.SetPhoto)
		r.Get("/groups/{groupId}/join-requests", groupCtrl.GetJoinRequests)
		r.Post("/groups/{groupId}/join-requests/handle", groupCtrl.HandleJoinRequest)

		// Contacts
		r.Post("/contacts/check", contactCtrl.Check)
		r.Get("/contacts/{jid}", contactCtrl.GetInfo)
		r.Get("/contacts/{jid}/profile-picture", contactCtrl.ProfilePicture)
		r.Get("/contacts/{jid}/business-profile", contactCtrl.BusinessProfile)
		r.Post("/contacts/{jid}/block", contactCtrl.Block)
		r.Post("/contacts/{jid}/unblock", contactCtrl.Unblock)
		r.Get("/contacts/blocklist", contactCtrl.Blocklist)

		// Presence
		r.Put("/presence", presenceCtrl.Set)
		r.Post("/presence/{jid}/subscribe", presenceCtrl.Subscribe)

		// Newsletters
		r.Get("/newsletters", newsletterCtrl.List)
		r.Get("/newsletters/{newsletterId}", newsletterCtrl.Get)
		r.Post("/newsletters/{newsletterId}/follow", newsletterCtrl.Follow)
		r.Post("/newsletters/{newsletterId}/unfollow", newsletterCtrl.Unfollow)
		r.Post("/newsletters/{newsletterId}/mute", newsletterCtrl.Mute)
		r.Post("/newsletters/{newsletterId}/unmute", newsletterCtrl.Unmute)

		// Calls
		r.Post("/calls/{callId}/reject", callCtrl.Reject)

		// Privacy
		r.Get("/privacy", privacyCtrl.Get)
		r.Patch("/privacy", privacyCtrl.Update)

		// Profile
		r.Put("/profile/status-message", profileCtrl.SetStatus)
		r.Get("/profile/qr-link", profileCtrl.QRLink)
		r.Post("/profile/qr-link/revoke", profileCtrl.RevokeQRLink)
		r.Post("/profile/avatar", instanceCtrl.SetAvatar)
		r.Post("/profile/pushname", instanceCtrl.SetPushName)

		// Webhooks
		r.Post("/webhooks", webhookCtrl.Create)
		r.Get("/webhooks", webhookCtrl.List)
		r.Get("/webhooks/{webhookId}", webhookCtrl.Get)
		r.Delete("/webhooks/{webhookId}", webhookCtrl.Delete)
		r.Post("/webhooks/{webhookId}/test", webhookCtrl.Test)

		// WebSockets (Real-time Stream)
		r.Get("/ws", wsCtrl.Connect)

		// Labels
		r.Get("/labels", labelCtrl.List)
		r.Post("/labels", labelCtrl.Create)
		r.Delete("/labels/{labelId}", labelCtrl.Delete)
		r.Post("/labels/{labelId}/chat", labelCtrl.LabelChat)
	})
}
