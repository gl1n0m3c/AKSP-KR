package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gl1n0m3c/AKSP-KR/services/feemous/internal/store"
)

const tokenCookie = "X-USER-TOKEN"

type ctxKey string

const (
	ctxUserID  ctxKey = "user_id"
	ctxIsAdmin ctxKey = "is_admin"
)

type registerRequest struct {
	Login      string  `json:"login"`
	Password   string  `json:"password"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Patronymic string  `json:"patronymic"`
	Email      string  `json:"email"`
	Phone      *string `json:"phone"`
	Telegram   *string `json:"telegram"`
	Status     string  `json:"status"`
	PositionID *string `json:"position_id"`
	TeamID     *string `json:"team_id"`
	UnitID     *string `json:"unit_id"`
	HeadID     *string `json:"head_id"`
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" || req.Email == "" || req.FirstName == "" || req.LastName == "" || req.Patronymic == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "hash error", http.StatusInternalServerError)
		return
	}

	uid := uuid.NewString()
	user := &UserDTO{
		ID:           uid,
		Login:        req.Login,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Patronymic:   req.Patronymic,
		Email:        req.Email,
		Phone:        req.Phone,
		Telegram:     req.Telegram,
		Status:       req.Status,
		PositionID:   req.PositionID,
		TeamID:       req.TeamID,
		UnitID:       req.UnitID,
		HeadID:       req.HeadID,
	}

	if user.Status == "" {
		user.Status = "active"
	}

	if err := s.store.CreateUser(r.Context(), user.ToModel()); err != nil {
		http.Error(w, "create user: "+err.Error(), http.StatusBadRequest)
		return
	}

	// auto-login
	s.issueToken(w, uid)

	writeJSON(w, http.StatusCreated, map[string]string{"id": uid})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := s.store.FindUserByLogin(r.Context(), req.Login)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	s.issueToken(w, u.ID)
	writeJSON(w, http.StatusOK, map[string]string{"id": u.ID})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ctxUserID)
	if uid == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	u, err := s.store.FindUserByID(r.Context(), uid.(string))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(tokenCookie)
		if err != nil || token.Value == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		uid := token.Value
		ctx := context.WithValue(r.Context(), ctxUserID, uid)
		ctx = context.WithValue(ctx, ctxIsAdmin, uid == adminID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		admin, _ := r.Context().Value(ctxIsAdmin).(bool)
		if !admin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) issueToken(w http.ResponseWriter, userID string) {
	e := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     tokenCookie,
		Value:    userID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  e,
	})
}

// DTO helper
type UserDTO struct {
	ID           string
	Login        string
	PasswordHash string
	FirstName    string
	LastName     string
	Patronymic   string
	Email        string
	Phone        *string
	Telegram     *string
	Status       string
	PositionID   *string
	TeamID       *string
	UnitID       *string
	HeadID       *string
}

func (u *UserDTO) ToModel() *store.User {
	return &store.User{
		ID:           u.ID,
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Patronymic:   u.Patronymic,
		Email:        u.Email,
		Phone:        u.Phone,
		Telegram:     u.Telegram,
		Status:       u.Status,
		PositionID:   u.PositionID,
		TeamID:       u.TeamID,
		UnitID:       u.UnitID,
		HeadID:       u.HeadID,
	}
}

// Utils
func decodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func init() {
	_ = middleware.RequestID
}
