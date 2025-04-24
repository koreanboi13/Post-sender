package storage

type MessengerType string

const (
	Telegram MessengerType = "Telegram"
	VK       MessengerType = "Vk"
)

type ChatEntry struct {
	ID        string        `json:"id"`
	Messenger MessengerType `json:"messenger"`
}
