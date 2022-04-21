package httpServer

import (
	ws "GitHub/go-chat/backend/pkg/websocket"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (s *HTTPServer) handleWSMessage(data json.RawMessage, userID uuid.UUID) {
	request := struct {
		Content        string    `json:"content"`
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	if err := json.Unmarshal([]byte(data), &request); err != nil {
		log.Println(err)
		return
	}

	err := s.app.Commands.ConversationService.SendTextMessage(request.Content, request.ConversationId, userID)

	if err != nil {
		log.Println(err)
		return
	}
}

func (s *HTTPServer) handleWS(wsHandlers ws.WSHandlers) http.HandlerFunc {
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

		s.WSConnectionsHub.RegisterClient(client)

		go client.WritePump()
		go client.ReadPump()
	}
}

func (s *HTTPServer) handleCreateConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationName string    `json:"conversation_name"`
		ConversationId   uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value("userId").(uuid.UUID)

	err = s.app.Commands.ConversationService.CreatePublicConversation(request.ConversationId, request.ConversationName, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("OK")
}

func (s *HTTPServer) handleDeleteConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Commands.ConversationService.DeleteConversation(request.ConversationId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("OK")
}

func (s *HTTPServer) handleJoinPublicConversation(w http.ResponseWriter, r *http.Request) {
	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value("userId").(uuid.UUID)

	err = s.app.Commands.ConversationService.JoinPublicConversation(request.ConversationId, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("OK")
}

func (s *HTTPServer) handleLeavePublicConversation(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userId").(uuid.UUID)
	request := struct {
		ConversationId uuid.UUID `json:"conversation_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Commands.ConversationService.LeavePublicConversation(request.ConversationId, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("OK")
}

func (s *HTTPServer) handleRenamePublicConversation(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userId").(uuid.UUID)
	request := struct {
		ConversationId   uuid.UUID `json:"conversation_id"`
		ConversationName string    `json:"new_name"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Commands.ConversationService.RenamePublicConversation(request.ConversationId, userID, request.ConversationName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("OK")
}
