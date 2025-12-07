package server

import (
	"html/template"
	"net/http"
)

var adminTpl = template.Must(template.New("admin").Parse(`
<html><body>
<h1>Admin</h1>
<h2>Users</h2>
<pre>{{printf "%+v" .Users}}</pre>
<h2>Teams</h2>
<pre>{{printf "%+v" .Teams}}</pre>
<h2>Positions</h2>
<pre>{{printf "%+v" .Positions}}</pre>
<h2>Meetings</h2>
<pre>{{printf "%+v" .Meetings}}</pre>
</body></html>
`))

func (s *Server) handleAdminPage(w http.ResponseWriter, r *http.Request) {
	users, _ := s.store.ListUsers(r.Context())
	teams, _ := s.store.ListTeams(r.Context())
	positions, _ := s.store.ListPositions(r.Context())
	units, _ := s.store.ListUnits(r.Context())
	meetings, _ := s.store.ListMeetings(r.Context())

	_ = adminTpl.Execute(w, map[string]any{
		"Users":     users,
		"Teams":     teams,
		"Positions": positions,
		"Units":     units,
		"Meetings":  meetings,
	})
}
