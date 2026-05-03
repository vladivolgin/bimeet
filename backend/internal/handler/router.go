package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	authhandler         "bimeet/internal/handler/auth"
	carpoolhandler      "bimeet/internal/handler/carpool"
	collectionhandler   "bimeet/internal/handler/collection"
	eventhandler        "bimeet/internal/handler/event"
	eventlinkhandler    "bimeet/internal/handler/eventlink"
	itemhandler         "bimeet/internal/handler/item"
	notificationhandler "bimeet/internal/handler/notification"
	pollhandler         "bimeet/internal/handler/poll"
	profilehandler      "bimeet/internal/handler/profile"
	"bimeet/internal/middleware"
)

func NewRouter(
	auth         *authhandler.Handler,
	event        *eventhandler.Handler,
	collection   *collectionhandler.Handler,
	poll         *pollhandler.Handler,
	item         *itemhandler.Handler,
	carpool      *carpoolhandler.Handler,
	link         *eventlinkhandler.Handler,
	notification *notificationhandler.Handler,
	profile      *profilehandler.Handler,
	jwtSecret    string,
) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logging)
	r.Use(middleware.CORS())

	// Public routes
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", auth.Register)
		r.Post("/login", auth.Login)
		r.Post("/forgot-password", auth.ForgotPassword)
		r.Post("/reset-password", auth.ResetPassword)
	})

	// Public invite link (GET only)
	r.Get("/api/events/invite/{token}", event.GetByInviteToken)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))

		// Join by invite token (POST requires auth)
		r.Post("/api/events/invite/{token}", event.JoinByInviteToken)

		// Public events list (must be before /{id} to avoid conflict)
		r.Get("/api/events/public", event.ListPublic)

		// Events
		r.Route("/api/events", func(r chi.Router) {
			r.Get("/", event.List)
			r.Post("/", event.Create)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", event.GetByID)
				r.Put("/", event.Update)
				r.Delete("/", event.Delete)
				r.Post("/complete", event.Complete)
				r.Post("/join", event.JoinPublic)

				// Participants
				r.Post("/participants", event.InviteParticipant)
				r.Patch("/participants/{userId}", event.UpdateParticipantStatus)

				// Collections
				r.Get("/collections", collection.List)
				r.Post("/collections", collection.Create)
				r.Get("/collections/summary", collection.Summary)
				r.Delete("/collections/{collectionId}", collection.Delete)
				r.Post("/collections/{collectionId}/contribute", collection.SubmitContribution)
				r.Post("/collections/{collectionId}/contributions/{contribId}/confirm", collection.ConfirmContribution)
				r.Post("/collections/{collectionId}/contributions/{contribId}/reject", collection.RejectContribution)
				r.Post("/collections/{collectionId}/contributions/mark-paid", collection.MarkPaid)

				// Polls
				r.Get("/polls", poll.List)
				r.Post("/polls", poll.Create)
				r.Post("/polls/{pollId}/vote", poll.Vote)

				// Items
				r.Get("/items", item.List)
				r.Post("/items", item.Create)
				r.Patch("/items/{itemId}", item.UpdateAssignment)

				// Carpools
				r.Get("/carpools", carpool.List)
				r.Post("/carpools", carpool.Create)
				r.Post("/carpools/{carpoolId}/join", carpool.Join)

				// Links
				r.Get("/links", link.List)
				r.Post("/links", link.Create)
				r.Delete("/links/{linkId}", link.Delete)
			})
		})

		// Notifications
		r.Get("/api/notifications", notification.List)
		r.Patch("/api/notifications/{id}/read", notification.MarkRead)
		r.Post("/api/notifications/read-all", notification.MarkAllRead)
		r.Delete("/api/notifications/{id}", notification.Delete)
		r.Delete("/api/notifications", notification.DeleteAll)

		// Profile
		r.Get("/api/auth/me", profile.GetMe)
		r.Put("/api/auth/me", profile.UpdateMe)
		r.Get("/api/auth/me/stats", profile.GetStats)
		r.Post("/api/auth/me/avatar", profile.UploadAvatar)
		r.Delete("/api/auth/me/avatar", profile.DeleteAvatar)
	})

	return r
}
