package store

import "time"

type User struct {
	ID           string     `json:"id"`
	Login        string     `json:"login"`
	PasswordHash string     `json:"-"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Patronymic   string     `json:"patronymic"`
	Email        string     `json:"email"`
	Phone        *string    `json:"phone,omitempty"`
	Telegram     *string    `json:"telegram,omitempty"`
	Status       string     `json:"status"`
	PositionID   *string    `json:"position_id"`
	TeamID       *string    `json:"team_id"`
	UnitID       *string    `json:"unit_id"`
	HeadID       *string    `json:"head_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type Session struct {
	Token     string     `json:"token"`
	UserID    string     `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type Team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Position struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Meeting struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ScheduledAt time.Time  `json:"scheduled_at"`
	OrganizerID string     `json:"organizer_id"`
	TeamID      *string    `json:"team_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type MeetingParticipant struct {
	MeetingID string `json:"meeting_id"`
	UserID    string `json:"user_id"`
}

type Unit struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
