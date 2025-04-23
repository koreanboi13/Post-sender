package http

import (
	"encoding/json"
	"net/http"

	"db/internal/service"
	"db/internal/storage"

	"github.com/gorilla/mux"
)

type Handler struct {
	chatService *service.ChatService
}

func NewHandler(chatService *service.ChatService) *Handler {
	return &Handler{
		chatService: chatService,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/saveChat", h.SaveChat).Methods("POST")
	router.HandleFunc("/api/deleteChat/{id}", h.DeleteChat).Methods("DELETE")
	router.HandleFunc("/api/chatExist/{id}", h.ChatExists).Methods("GET")
	router.HandleFunc("/api/allChats/{messenger}", h.GetChatsByMessenger).Methods("GET")
}

func (h *Handler) SaveChat(w http.ResponseWriter, r *http.Request) {
	var req SaveChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var messengerType storage.MessengerType
	switch req.Messenger {
	case string(storage.Telegram):
		messengerType = storage.Telegram
	case string(storage.VK):
		messengerType = storage.VK
	default:
		h.respondWithError(w, http.StatusBadRequest, "Invalid messenger type")
		return
	}

	if err := h.chatService.SaveChat(req.ID, messengerType); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, response{
		Success: true,
	})
}

func (h *Handler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["id"]

	if err := h.chatService.DeleteChat(chatID); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, response{
		Success: true,
	})
}

func (h *Handler) ChatExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["id"]

	exists, err := h.chatService.ChatExists(chatID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, response{
		Success: true,
		Data:    exists,
	})
}

func (h *Handler) GetChatsByMessenger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messengerStr := vars["messenger"]

	var messengerType storage.MessengerType
	switch messengerStr {
	case string(storage.Telegram):
		messengerType = storage.Telegram
	case string(storage.VK):
		messengerType = storage.VK
	default:
		h.respondWithError(w, http.StatusBadRequest, "Invalid messenger type")
		return
	}

	chatIDs, err := h.chatService.GetChatsByMessenger(messengerType)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, response{
		Success: true,
		Data:    chatIDs,
	})
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, response{
		Success: false,
		Error:   message,
	})
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error encoding response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
