package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Repository interface {
	FetchMeetings(ctx context.Context, window time.Duration, period time.Duration) ([]MeetingNotification, error)
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) FetchMeetings(ctx context.Context, window time.Duration, period time.Duration) ([]MeetingNotification, error) {
	now := time.Now().UTC()
	startAt := now.Add(window)
	endAt := startAt.Add(period)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"m.id",
			"m.title",
			"m.description",
			"m.scheduled_at",
			"u.id AS user_id",
			"u.email",
			"u.telegram",
		).
		From("meetings m").
		LeftJoin("users u ON u.team_id = m.team_id").
		Where("m.scheduled_at >= ?", startAt).
		Where("m.scheduled_at < ?", endAt)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	events := map[string]*MeetingNotification{}

	for rows.Next() {
		var row MeetingRow
		if err = rows.Scan(
			&row.ID,
			&row.Title,
			&row.Description,
			&row.ScheduledAt,
			&row.UserID,
			&row.Email,
			&row.Telegram,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		ev, ok := events[row.ID]
		if !ok {
			ev = &MeetingNotification{
				ID:          row.ID,
				Title:       row.Title,
				Description: row.Description,
				ScheduledAt: row.ScheduledAt,
			}
			events[row.ID] = ev
		}

		ev.Users = append(ev.Users, NotificationUser{
			ID:       nullString(row.UserID),
			Email:    nullString(row.Email),
			Telegram: nullString(row.Telegram),
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	out := make([]MeetingNotification, 0, len(events))
	for _, ev := range events {
		out = append(out, *ev)
	}

	return out, nil
}

func nullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
