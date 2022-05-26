package httpHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type commandHandlers struct {
	commands   *app.Commands
	wsHandlers WSHandlers
}

func NewCommandHandlers(commands *app.Commands) *commandHandlers {
	return &commandHandlers{
		commands:   commands,
		wsHandlers: NewWSHandlers(commands),
	}
}

func (s *commandHandlers) handleOpenWSConnection() http.HandlerFunc {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := os.Getenv("CLIENT_ORIGIN")

			return r.Header.Get("Origin") == origin
		},
	}

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

		s.commands.ClientsService.RegisterClient(conn, userID, s.wsHandlers.HandleNotification)
	}
}
