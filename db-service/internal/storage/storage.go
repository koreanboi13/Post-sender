package storage

type Storage interface {
	Save(chatID string, messengerType MessengerType) error

	Delete(chatID string) error

	Exists(chatID string) (bool, error)

	GetAllByMessenger(messengerType MessengerType) ([]int, error)

	Close() error
}
