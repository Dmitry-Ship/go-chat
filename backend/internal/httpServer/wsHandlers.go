package httpServer

import (
	ws "GitHub/go-chat/backend/internal/infra/websocket"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (s *CommandController) handleOpenWSConnection(wsHandlers ws.WSHandlers) http.HandlerFunc {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := os.Getenv("API_URL")

			return r.Header.Get("Origin") == origin
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.commands.NotificationsService.RegisterClient(conn, wsHandlers, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (s *CommandController) handleReceiveWSChatMessage(data json.RawMessage, userID uuid.UUID) {
	request := struct {
		Content        string    `json:"content"`
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		log.Println(err)
		return
	}

	err := s.commands.MessagingService.SendTextMessage(request.Content, request.ConversationId, userID)

	if err != nil {
		log.Println(err)
		return
	}
}
