package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/gl1n0m3c/AKSP-KR/services/feemous/internal/store"
)

type eventRequest struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ScheduledAt  time.Time `json:"scheduled_at"`
	TeamID       *string   `json:"team_id"`
	Participants []string  `json:"participants"`
}

type eventResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ScheduledAt  time.Time `json:"scheduled_at"`
	OrganizerID  string    `json:"organizer_id"`
	TeamID       *string   `json:"team_id"`
	Participants []string  `json:"participants"`
}

func (s *Server) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	var req eventRequest
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uid, _ := r.Context().Value(ctxUserID).(string)
	m := &store.Meeting{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		ScheduledAt: req.ScheduledAt,
		OrganizerID: uid,
		TeamID:      req.TeamID,
	}
	if err := s.store.CreateMeeting(r.Context(), m); err != nil {
		http.Error(w, "create: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.store.ReplaceMeetingParticipants(r.Context(), m.ID, req.Participants); err != nil {
		http.Error(w, "participants: "+err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": m.ID})
}

func (s *Server) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req eventRequest
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m := &store.Meeting{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		ScheduledAt: req.ScheduledAt,
		TeamID:      req.TeamID,
	}
	if err := s.store.UpdateMeeting(r.Context(), m); err != nil {
		http.Error(w, "update: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.store.ReplaceMeetingParticipants(r.Context(), id, req.Participants); err != nil {
		http.Error(w, "participants: "+err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (s *Server) handleListEvents(w http.ResponseWriter, r *http.Request) {
	meetings, err := s.store.ListMeetings(r.Context())
	if err != nil {
		http.Error(w, "list: "+err.Error(), http.StatusInternalServerError)
		return
	}
	parts, err := s.store.ListMeetingParticipants(r.Context())
	if err != nil {
		http.Error(w, "participants: "+err.Error(), http.StatusInternalServerError)
		return
	}
	partMap := make(map[string][]string)
	for _, p := range parts {
		partMap[p.MeetingID] = append(partMap[p.MeetingID], p.UserID)
	}
	resp := make([]eventResponse, 0, len(meetings))
	for _, m := range meetings {
		resp = append(resp, eventResponse{
			ID:           m.ID,
			Title:        m.Title,
			Description:  m.Description,
			ScheduledAt:  m.ScheduledAt,
			OrganizerID:  m.OrganizerID,
			TeamID:       m.TeamID,
			Participants: partMap[m.ID],
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.store.DeleteMeeting(r.Context(), id); err != nil {
		http.Error(w, "delete: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
