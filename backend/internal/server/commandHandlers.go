package server

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) handleOpenWSConnection() http.HandlerFunc {
	var upgrader = WebSocketUpgrader

	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)

		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			returnError(w, http.StatusInternalServerError, err)
			return
		}

		s.notificationCommands.RegisterClient(r.Context(), conn, userID)
	}
}
