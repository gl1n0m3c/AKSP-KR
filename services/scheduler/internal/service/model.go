package service

import (
	"database/sql"
	"encoding/json"
	"time"
)

type MeetingRow struct {
	ID          string         `db:"id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	ScheduledAt time.Time      `db:"scheduled_at"`
	UserID      sql.NullString `db:"user_id"`
	Email       sql.NullString `db:"email"`
	Telegram    sql.NullString `db:"telegram"`
}

type MeetingNotification struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	ScheduledAt time.Time          `json:"scheduled_at"`
	Users       []NotificationUser `json:"users"`
}

type NotificationUser struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Telegram string `json:"telegram,omitempty"`
}

func marshalEvent(ev MeetingNotification) ([]byte, error) {
	return json.Marshal(ev)
}
