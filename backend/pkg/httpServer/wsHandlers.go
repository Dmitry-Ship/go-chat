package httpServer

import (
	ws "GitHub/go-chat/backend/pkg/websocket"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (s *HTTPServer) handleOpenWSConnection(wsHandlers ws.WSHandlers) http.HandlerFunc {
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
			fmt.Println("WS", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client := ws.NewClient(conn, s.WSConnectionsHub, wsHandlers, userID)

		err = s.app.Commands.NotificationsService.RegisterClient(client)

		if err != nil {
			log.Println("WS", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go client.WritePump()
		go client.ReadPump()
	}
}

func (s *HTTPServer) handleReceiveWSChatMessage(data json.RawMessage, userID uuid.UUID) {
	request := struct {
		Content        string    `json:"content"`
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		log.Println(err)
		return
	}

	err := s.app.Commands.MessagingService.SendTextMessage(request.Content, request.ConversationId, userID)

	if err != nil {
		log.Println(err)
		return
	}
}
