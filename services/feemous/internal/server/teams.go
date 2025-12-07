package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/gl1n0m3c/AKSP-KR/services/feemous/internal/store"
)

type simpleReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) handleCreateTeam(w http.ResponseWriter, r *http.Request) {
	var req simpleReq
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t := &store.Team{ID: uuid.NewString(), Name: req.Name, Description: req.Description}
	if err := s.store.CreateTeam(r.Context(), t); err != nil {
		http.Error(w, "create: "+err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": t.ID})
}

func (s *Server) handleDeleteTeam(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.store.DeleteTeam(r.Context(), id); err != nil {
		http.Error(w, "delete: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleListTeams(w http.ResponseWriter, r *http.Request) {
	items, err := s.store.ListTeams(r.Context())
	if err != nil {
		http.Error(w, "list: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
