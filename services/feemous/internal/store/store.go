package store

import (
	"context"
	"database/sql"
)

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}

// Users
func (s *Store) CreateUser(ctx context.Context, u *User) error {
	q := `INSERT INTO users
	(id, login, password_hash, first_name, last_name, patronymic, email, phone, telegram,
     status, position_id, team_id, unit_id, head_id, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,NOW())`
	_, err := s.DB.ExecContext(ctx, q,
		u.ID, u.Login, u.PasswordHash, u.FirstName, u.LastName, u.Patronymic,
		u.Email, u.Phone, u.Telegram, u.Status, u.PositionID, u.TeamID, u.UnitID, u.HeadID)
	return err
}

func (s *Store) FindUserByLogin(ctx context.Context, login string) (*User, error) {
	q := `SELECT id, login, password_hash, first_name, last_name, patronymic, email, phone, telegram,
       status, position_id, team_id, unit_id, head_id, created_at, updated_at
       FROM users WHERE login=$1`
	var u User
	err := s.DB.QueryRowContext(ctx, q, login).Scan(
		&u.ID, &u.Login, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Patronymic,
		&u.Email, &u.Phone, &u.Telegram, &u.Status, &u.PositionID, &u.TeamID, &u.UnitID, &u.HeadID,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) FindUserByID(ctx context.Context, id string) (*User, error) {
	q := `SELECT id, login, password_hash, first_name, last_name, patronymic, email, phone, telegram,
       status, position_id, team_id, unit_id, head_id, created_at, updated_at
       FROM users WHERE id=$1`
	var u User
	err := s.DB.QueryRowContext(ctx, q, id).Scan(
		&u.ID, &u.Login, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Patronymic,
		&u.Email, &u.Phone, &u.Telegram, &u.Status, &u.PositionID, &u.TeamID, &u.UnitID, &u.HeadID,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, login, first_name, last_name, patronymic, email, phone, telegram, status, position_id, team_id, unit_id, head_id, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID, &u.Login, &u.FirstName, &u.LastName, &u.Patronymic,
			&u.Email, &u.Phone, &u.Telegram, &u.Status,
			&u.PositionID, &u.TeamID, &u.UnitID, &u.HeadID,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// Teams
func (s *Store) CreateTeam(ctx context.Context, t *Team) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO teams (id,name,description) VALUES ($1,$2,$3)`, t.ID, t.Name, t.Description)
	return err
}

func (s *Store) DeleteTeam(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM teams WHERE id=$1`, id)
	return err
}

func (s *Store) ListTeams(ctx context.Context) ([]Team, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,name,description FROM teams`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Team
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Description); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// Positions
func (s *Store) CreatePosition(ctx context.Context, p *Position) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO positions (id,name,description) VALUES ($1,$2,$3)`, p.ID, p.Name, p.Description)
	return err
}

func (s *Store) DeletePosition(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM positions WHERE id=$1`, id)
	return err
}

func (s *Store) ListPositions(ctx context.Context) ([]Position, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,name,description FROM positions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Position
	for rows.Next() {
		var p Position
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Meetings
func (s *Store) CreateMeeting(ctx context.Context, m *Meeting) error {
	_, err := s.DB.ExecContext(ctx, `
INSERT INTO meetings (id,title,description,scheduled_at,created_at,organizer_id,team_id)
VALUES ($1,$2,$3,$4,NOW(),$5,$6)`,
		m.ID, m.Title, m.Description, m.ScheduledAt, m.OrganizerID, m.TeamID)
	return err
}

func (s *Store) UpdateMeeting(ctx context.Context, m *Meeting) error {
	_, err := s.DB.ExecContext(ctx, `
UPDATE meetings
SET title=$2, description=$3, scheduled_at=$4, updated_at=NOW(), team_id=$5
WHERE id=$1`,
		m.ID, m.Title, m.Description, m.ScheduledAt, m.TeamID)
	return err
}

func (s *Store) DeleteMeeting(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM meetings WHERE id=$1`, id)
	return err
}

func (s *Store) ListMeetings(ctx context.Context) ([]Meeting, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,title,description,scheduled_at,organizer_id,team_id,created_at,updated_at FROM meetings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Meeting
	for rows.Next() {
		var m Meeting
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.ScheduledAt, &m.OrganizerID, &m.TeamID, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Store) ListMeetingParticipants(ctx context.Context) ([]MeetingParticipant, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT meeting_id, user_id FROM meeting_participants`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []MeetingParticipant
	for rows.Next() {
		var mp MeetingParticipant
		if err := rows.Scan(&mp.MeetingID, &mp.UserID); err != nil {
			return nil, err
		}
		out = append(out, mp)
	}
	return out, rows.Err()
}

// Meeting participants
func (s *Store) ReplaceMeetingParticipants(ctx context.Context, meetingID string, userIDs []string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM meeting_participants WHERE meeting_id=$1`, meetingID); err != nil {
		return err
	}
	for _, uid := range userIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO meeting_participants (meeting_id, user_id) VALUES ($1,$2)`, meetingID, uid); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Units
func (s *Store) CreateUnit(ctx context.Context, u *Unit) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO units (id,name,description) VALUES ($1,$2,$3)`, u.ID, u.Name, u.Description)
	return err
}

func (s *Store) DeleteUnit(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM units WHERE id=$1`, id)
	return err
}

func (s *Store) ListUnits(ctx context.Context) ([]Unit, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,name,description FROM units`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Unit
	for rows.Next() {
		var u Unit
		if err := rows.Scan(&u.ID, &u.Name, &u.Description); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
