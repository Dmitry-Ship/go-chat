package httpHandlers

import (
	"GitHub/go-chat/backend/internal/app"
	"net/http"
	"os"

	ws "GitHub/go-chat/backend/internal/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type commandController struct {
	commands       *app.Commands
	clientRegister ws.ClientRegister
}

func NewCommandController(commands *app.Commands, clientRegister ws.ClientRegister) *commandController {
	return &commandController{
		commands:       commands,
		clientRegister: clientRegister,
	}
}

func (s *commandController) handleOpenWSConnection() http.HandlerFunc {
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
			http.Error(w, "userId not found in context", http.StatusInternalServerError)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.clientRegister.RegisterClient(conn, userID)
	}
}
