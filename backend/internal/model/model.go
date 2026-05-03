package model

import (
	"time"

	"github.com/google/uuid"
)

// ─── Domain structs ────────────────────────────────────────────────────────

type User struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	BirthDate    *time.Time `json:"birth_date"`
	City         string     `json:"city"`
	AvatarURL    *string    `json:"avatar_url"`
	CreatedAt    time.Time  `json:"created_at"`
}

type UpdateProfileRequest struct {
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"` // "YYYY-MM-DD", empty = clear
	City      string `json:"city"`
}

type ProfileStats struct {
	Organized    int `json:"organized"`
	Participated int `json:"participated"`
	Completed    int `json:"completed"`
	Upcoming     int `json:"upcoming"`
}

type ReminderEvent struct {
	ID             uuid.UUID
	Title          string
	DateEnd        time.Time
	OrganizerID    uuid.UUID
	ParticipantIDs []uuid.UUID
	Needs3d        bool
	Needs1d        bool
}

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateStart   time.Time `json:"date_start"`
	DateEnd     time.Time `json:"date_end"`
	Location    string    `json:"location"`
	Category    string    `json:"category"` // ordinary | business
	DressCode   *string   `json:"dress_code,omitempty"`
	Status      string    `json:"status"` // active | completed
	IsPublic    bool      `json:"is_public"`
	MaxGuests   *int      `json:"max_guests,omitempty"`
	OrganizerID uuid.UUID `json:"organizer_id"`
	InviteToken uuid.UUID `json:"invite_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventLink struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type EventParticipant struct {
	ID      uuid.UUID `json:"id"`
	EventID uuid.UUID `json:"event_id"`
	UserID  uuid.UUID `json:"user_id"`
	Status  string    `json:"status"` // invited | confirmed | declined
}

type Collection struct {
	ID              uuid.UUID `json:"id"`
	EventID         uuid.UUID `json:"event_id"`
	Title           string    `json:"title"`
	PerPersonAmount float64   `json:"per_person_amount"`
	CreatedBy       uuid.UUID `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}

type CollectionContribution struct {
	ID           uuid.UUID  `json:"id"`
	CollectionID uuid.UUID  `json:"collection_id"`
	UserID       uuid.UUID  `json:"user_id"`
	Paid         bool       `json:"paid"`
	PaidAt       *time.Time `json:"paid_at,omitempty"`
	Status       string     `json:"status"`                  // not_paid | pending | paid
	ReceiptURL   *string    `json:"receipt_url,omitempty"`
}

type Poll struct {
	ID        uuid.UUID    `json:"id"`
	EventID   uuid.UUID    `json:"event_id"`
	Question  string       `json:"question"`
	CreatedBy uuid.UUID    `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	Options   []PollOption `json:"options,omitempty"`
}

type PollOption struct {
	ID     uuid.UUID `json:"id"`
	PollID uuid.UUID `json:"poll_id"`
	Label  string    `json:"label"`
	Votes  int       `json:"votes"`
}

type PollVote struct {
	PollOptionID uuid.UUID `json:"poll_option_id"`
	UserID       uuid.UUID `json:"user_id"`
}

type Item struct {
	ID         uuid.UUID  `json:"id"`
	EventID    uuid.UUID  `json:"event_id"`
	Name       string     `json:"name"`
	AssignedTo *uuid.UUID `json:"assigned_to,omitempty"`
}

type Carpool struct {
	ID              uuid.UUID `json:"id"`
	EventID         uuid.UUID `json:"event_id"`
	DriverID        uuid.UUID `json:"driver_id"`
	SeatsAvailable  int       `json:"seats_available"`
	DeparturePoint  string    `json:"departure_point"`
}

type CarpoolPassenger struct {
	CarpoolID uuid.UUID `json:"carpool_id"`
	UserID    uuid.UUID `json:"user_id"`
}

type EventChangeLog struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	ChangedBy uuid.UUID `json:"changed_by"`
	FieldName string    `json:"field_name"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedAt time.Time `json:"changed_at"`
}

type Notification struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	EventID   *uuid.UUID `json:"event_id,omitempty"`
	Type      string     `json:"type"`
	Message   string     `json:"message"`
	IsRead    bool       `json:"is_read"`
	CreatedAt time.Time  `json:"created_at"`
}

type PasswordResetToken struct {
	Token     uuid.UUID `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// ─── DTOs ──────────────────────────────────────────────────────────────────

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateStart   time.Time `json:"date_start"`
	DateEnd     time.Time `json:"date_end"`
	Location    string    `json:"location"`
	Category    string    `json:"category"` // ordinary | business; defaults to ordinary
	DressCode   *string   `json:"dress_code,omitempty"`
	IsPublic    bool      `json:"is_public"`
	MaxGuests   *int      `json:"max_guests,omitempty"`
}

type UpdateEventRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	DateStart   *time.Time `json:"date_start,omitempty"`
	DateEnd     *time.Time `json:"date_end,omitempty"`
	Location    *string    `json:"location,omitempty"`
	Category    *string    `json:"category,omitempty"`
	DressCode   *string    `json:"dress_code,omitempty"`
	IsPublic    *bool      `json:"is_public,omitempty"`
	MaxGuests   *int       `json:"max_guests,omitempty"`
}

type CreateEventLinkRequest struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type InviteParticipantRequest struct {
	Email string `json:"email"`
}

type JoinByInviteTokenRequest struct {
	Action string `json:"action"` // "join" | "decline"; defaults to "join"
}

type UpdateParticipantStatusRequest struct {
	Status string `json:"status"` // confirmed | declined
}

type CreateCollectionRequest struct {
	Title           string  `json:"title"`
	PerPersonAmount float64 `json:"per_person_amount"`
}

type CollectionSummary struct {
	CollectionID    uuid.UUID `json:"collection_id"`
	Title           string    `json:"title"`
	PerPersonAmount float64   `json:"per_person_amount"`
	ExpectedTotal   float64   `json:"expected_total"`
	PaidCount       int       `json:"paid_count"`
	TotalCount      int       `json:"total_count"`
	TotalPaid       float64   `json:"total_paid"`
	Remaining       float64   `json:"remaining"`
}

type EventCollectionsSummaryResponse struct {
	Collections []CollectionSummary `json:"collections"`
	GrandTotal  float64             `json:"grand_total"`
	TotalPaid   float64             `json:"total_paid"`
	Remaining   float64             `json:"remaining"`
}

type CreatePollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type VoteRequest struct {
	OptionID uuid.UUID `json:"option_id"`
}

type CreateItemRequest struct {
	Name string `json:"name"`
}

type UpdateItemRequest struct {
	AssignedTo *uuid.UUID `json:"assigned_to"` // null to unassign
}

type CreateCarpoolRequest struct {
	SeatsAvailable int    `json:"seats_available"`
	DeparturePoint string `json:"departure_point"`
}

// ─── Enriched API response types ───────────────────────────────────────────

type InviteEventInfo struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	DateStart      time.Time `json:"date_start"`
	DateEnd        time.Time `json:"date_end"`
	Location       string    `json:"location"`
	Category       string    `json:"category"`
	DressCode      *string   `json:"dress_code,omitempty"`
	IsPublic       bool      `json:"is_public"`
	Organizer      UserInfo  `json:"organizer"`
	ConfirmedCount int       `json:"confirmed_count"`
}

type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type ParticipantInfo struct {
	ID     uuid.UUID `json:"id"`
	User   UserInfo  `json:"user"`
	Status string    `json:"status"`
}

type ContributionInfo struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	Paid       bool       `json:"paid"`
	PaidAt     *time.Time `json:"paid_at,omitempty"`
	Status     string     `json:"status"`
	ReceiptURL *string    `json:"receipt_url,omitempty"`
}

type CollectionInfo struct {
	ID              uuid.UUID          `json:"id"`
	Title           string             `json:"title"`
	PerPersonAmount float64            `json:"per_person_amount"`
	CreatedBy       uuid.UUID          `json:"created_by"`
	CreatedAt       time.Time          `json:"created_at"`
	Contributions   []ContributionInfo `json:"contributions"`
}

type PollOptionInfo struct {
	ID    uuid.UUID   `json:"id"`
	Label string      `json:"label"`
	Votes []uuid.UUID `json:"votes"` // voter IDs
}

type PollInfo struct {
	ID        uuid.UUID        `json:"id"`
	Question  string           `json:"question"`
	CreatedBy uuid.UUID        `json:"created_by"`
	CreatedAt time.Time        `json:"created_at"`
	Options   []PollOptionInfo `json:"options"`
}

type ItemInfo struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	AssignedTo *UserInfo `json:"assigned_to,omitempty"`
}

type CarpoolInfo struct {
	ID             uuid.UUID  `json:"id"`
	Driver         UserInfo   `json:"driver"`
	SeatsAvailable int        `json:"seats_available"`
	DeparturePoint string     `json:"departure_point"`
	Passengers     []UserInfo `json:"passengers"`
}

type EventLinkInfo struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type ChangeLogInfo struct {
	ID        uuid.UUID `json:"id"`
	ChangedBy UserInfo  `json:"changed_by"`
	FieldName string    `json:"field_name"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedAt time.Time `json:"changed_at"`
}

type EventDetail struct {
	ID             uuid.UUID         `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	DateStart      time.Time         `json:"date_start"`
	DateEnd        time.Time         `json:"date_end"`
	Location       string            `json:"location"`
	Category       string            `json:"category"`
	DressCode      *string           `json:"dress_code,omitempty"`
	Status         string            `json:"status"`
	IsPublic       bool              `json:"is_public"`
	MaxGuests      *int              `json:"max_guests,omitempty"`
	ConfirmedCount int               `json:"confirmed_count"`
	Organizer      UserInfo          `json:"organizer"`
	InviteToken    uuid.UUID         `json:"invite_token"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Participants   []ParticipantInfo `json:"participants"`
	Collections    []CollectionInfo  `json:"collections"`
	Polls          []PollInfo        `json:"polls"`
	Items          []ItemInfo        `json:"items"`
	Carpools       []CarpoolInfo     `json:"carpools"`
	Links          []EventLinkInfo   `json:"links"`
	ChangeLog      []ChangeLogInfo   `json:"change_log"`
}

type EventListItem struct {
	ID             uuid.UUID         `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	DateStart      time.Time         `json:"date_start"`
	DateEnd        time.Time         `json:"date_end"`
	Location       string            `json:"location"`
	Category       string            `json:"category"`
	DressCode      *string           `json:"dress_code,omitempty"`
	Status         string            `json:"status"`
	IsPublic       bool              `json:"is_public"`
	MaxGuests      *int              `json:"max_guests,omitempty"`
	ConfirmedCount int               `json:"confirmed_count"`
	Organizer      UserInfo          `json:"organizer"`
	InviteToken    uuid.UUID         `json:"invite_token"`
	CreatedAt      time.Time         `json:"created_at"`
	MyStatus       string            `json:"my_status"`
	Participants   []ParticipantInfo `json:"participants"`
}

type PublicEventListItem struct {
	EventListItem
	IsParticipant bool `json:"is_participant"`
}
