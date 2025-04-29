package rabbitmq

import "time"

type Response struct {
	Success bool       `json:"success"`
	Data    []DataItem `json:"data"`
}

type DataItem struct {
	ID            int64     `json:"id"`
	Comment       string    `json:"comment"`
	Title         string    `json:"title"`
	Link          string    `json:"link"`
	UpdatedDate   time.Time `json:"updatedDate"`
	CollectedDate time.Time `json:"collectedDate"`
}
