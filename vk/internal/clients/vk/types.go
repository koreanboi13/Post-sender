package vk

import (
	"encoding/json"
)

type LongPollResponse struct {
	TS      json.Number     `json:"ts"`
	Updates json.RawMessage `json:"updates"`
	Failed  int             `json:"failed,omitempty"`
}

type LongPollUpdate struct {
	Type    string `json:"type"`
	Object  Object `json:"object"`
	GroupId int    `json:"group_id"`
}

type Object struct {
	Text   string `json:"text"`
	FromId int    `json:"from_id"`
}

type MessageResponse struct {
	Response Connection     `json:"response"`
	Error    *ErrorResponse `json:"error,omitempty"`
}

type Connection struct {
	TS     int    `json:"ts"`
	Key    string `json:"key"`
	Server string `json:"server"`
}

type ErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}
