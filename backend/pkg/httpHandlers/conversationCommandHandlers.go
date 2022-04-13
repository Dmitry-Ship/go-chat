package httpHandlers

import (
	"GitHub/go-chat/backend/pkg/services"
	ws "GitHub/go-chat/backend/pkg/websocket"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := os.Getenv("API_URL")

		return r.Header.Get("Origin") == origin
	},
}

func HandleWSMessage(conversationService services.ConversationService) ws.WSHandler {
	return func(incomingNotification ws.IncomingNotification, data json.RawMessage) {
		request := struct {
			Content        string    `json:"content"`
			ConversationId uuid.UUID `json:"conversation_id"`
		}{}

		if err := json.Unmarshal([]byte(data), &request); err != nil {
			log.Println(err)
			return
		}

		err := conversationService.SendTextMessage(request.Content, request.ConversationId, incomingNotification.UserID)

		if err != nil {
			log.Println(err)
			return
		}
	}
}

func HandleWS(hub ws.Hub, wsHandlers ws.WSHandlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("WS", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client := ws.NewClient(conn, hub, wsHandlers, userID)

		hub.RegisterClient(client)

		go client.WritePump()
		go client.ReadPump()
	}
}

func HandleCreateConversation(conversationService services.ConversationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = conversationService.CreatePublicConversation(request.ConversationId, request.ConversationName, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleDeleteConversation(conversationService services.ConversationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ConversationId uuid.UUID `json:"conversation_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = conversationService.DeleteConversation(request.ConversationId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleJoinPublicConversation(conversationService services.ConversationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ConversationId uuid.UUID `json:"conversation_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID, _ := r.Context().Value("userId").(uuid.UUID)

		err = conversationService.JoinPublicConversation(request.ConversationId, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleLeavePublicConversation(conversationService services.ConversationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("userId").(uuid.UUID)
		request := struct {
			ConversationId uuid.UUID `json:"conversation_id"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = conversationService.LeavePublicConversation(request.ConversationId, userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}

func HandleRenamePublicConversation(conversationService services.ConversationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = conversationService.RenamePublicConversation(request.ConversationId, userID, request.ConversationName)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("OK")
	}
}
