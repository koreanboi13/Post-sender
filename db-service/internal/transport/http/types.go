package http

type SaveChatRequest struct {
	ID        string `json:"id"`
	Messenger string `json:"messenger"`
}

type response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
