package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	DB       http.Client
	Host     string
	BasePath string
}

const (
	messangerType    = "Telegram"
	methodSaveChat   = "saveChat"
	methodDeleteChat = "deleteChat"
	methodChatExist  = "chatExist"
	methodAllChats   = "allChats"
)

func New(h string, b string) *Client {
	return &Client{
		DB:       http.Client{},
		Host:     h,
		BasePath: b,
	}
}

func (c *Client) SaveUser(ctx context.Context, chatId int) error {
	u := url.URL{
		Scheme: "http",
		Host:   c.Host,
		Path:   path.Join(c.BasePath, methodSaveChat),
	}

	data := map[string]interface{}{
		"id":        strconv.Itoa(chatId),
		"messenger": messangerType,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("can't marshal request data: %v", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u.String(),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("can't make req: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.DB.Do(req)
	if err != nil {
		return fmt.Errorf("can`t do request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	var success ErrorResponse
	if err := json.Unmarshal(body, &success); err != nil {
		return fmt.Errorf("can`t unmurshall json: %v", err)
	}

	if !success.Success {
		return fmt.Errorf("error while saving chatId: unknown error")
	}
	return nil
}

func (c *Client) DeleteUser(ctx context.Context, chatId int) error {
	u := url.URL{
		Scheme: "http",
		Host:   c.Host,
		Path:   path.Join(c.BasePath, methodDeleteChat, strconv.Itoa(chatId)),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return fmt.Errorf("can't make req: %v", err)
	}

	resp, err := c.DB.Do(req)
	if err != nil {
		return fmt.Errorf("can`t do request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	var success ErrorResponse
	if err := json.Unmarshal(body, &success); err != nil {
		return fmt.Errorf("can`t unmurshall json: %v", err)
	}

	if !success.Success {
		return fmt.Errorf("error while deleting chatId: %v", success.Success)
	}
	return nil
}

func (c *Client) ChatExists(ctx context.Context, chatId int) (bool, error) {
	u := url.URL{
		Scheme: "http",
		Host:   c.Host,
		Path:   path.Join(c.BasePath, methodChatExist, strconv.Itoa(chatId)),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return false, fmt.Errorf("can`t make req: %v", err)
	}

	resp, err := c.DB.Do(req)
	if err != nil {
		return false, fmt.Errorf("can`t do request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %v", err)
	}

	var success ErrorResponse
	if err := json.Unmarshal(body, &success); err != nil {
		return false, fmt.Errorf("can`t unmurshall json: %v", err)
	}
	log.Println(success)
	if !success.Success {
		return false, fmt.Errorf("API returned non-success response")
	}
	return success.Data, nil
}

func (c *Client) AllUsers(ctx context.Context, messangerType string) ([]int, error) {
	u := url.URL{
		Scheme: "http",
		Host:   c.Host,
		Path:   path.Join(c.BasePath, methodAllChats, messangerType),
	}
	log.Println("making request to: ", u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can`t make req: %v", err)
	}

	resp, err := c.DB.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can`t do request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("can`t unmurshall json: %v", err)
	}

	if !res.Success {
		return nil, fmt.Errorf("error while trying to get new posts: %v", err)
	}
	return res.Data, nil
}
