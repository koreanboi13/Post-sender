package db

import "time"

type ErrorResponse struct {
	Success bool `json:"success"`
	Data    bool `json:"data"`
}

type Response struct {
	Success bool  `json:"success"`
	Data    []int `json:"data"`
}

type DataItem struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Link          string    `json:"link"`
	UpdatedDate   time.Time `json:"updatedDate"`
	CollectedDate time.Time `json:"collectedDate"`
}
