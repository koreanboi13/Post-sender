package service

import (
	"fmt"

	"db/internal/storage"
)

type ChatService struct {
	storage storage.Storage
}

func NewChatService(storage storage.Storage) *ChatService {
	return &ChatService{
		storage: storage,
	}
}

func (s *ChatService) SaveChat(chatID string, messengerType storage.MessengerType) error {
	if chatID == "" {
		return fmt.Errorf("chat ID cannot be empty")
	}

	return s.storage.Save(chatID, messengerType)
}

func (s *ChatService) DeleteChat(chatID string) error {
	if chatID == "" {
		return fmt.Errorf("chat ID cannot be empty")
	}

	return s.storage.Delete(chatID)
}

func (s *ChatService) ChatExists(chatID string) (bool, error) {
	if chatID == "" {
		return false, fmt.Errorf("chat ID cannot be empty")
	}

	return s.storage.Exists(chatID)
}

func (s *ChatService) GetChatsByMessenger(messengerType storage.MessengerType) ([]string, error) {
	return s.storage.GetAllByMessenger(messengerType)
}

func (s *ChatService) Close() error {
	return s.storage.Close()
}
