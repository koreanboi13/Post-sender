package vk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	longPollUrl string
	apiUrl      string
	token       string
	client      http.Client
	apiVersion  string
	key         string
	server      string
}

const (
	getLongPollServerUrl = "https://api.vk.com/method/messages.getLongPollServer"
)

func New(token string) *Client {
	server, key, ts, err := GetLongPollServer(context.Background(), token, "5.199")
	if err != nil {
		log.Fatal("can`t get params: ", err)
	}
	longPoll := longPoll(server, key, ts)
	return &Client{
		longPollUrl: longPoll,
		apiUrl:      "https://api.vk.com/method",
		token:       token,
		client:      http.Client{},
		apiVersion:  "5.199",
		key:         key,
		server:      server,
	}
}

func (c *Client) Updates(ctx context.Context) ([]LongPollUpdate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.longPollUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response: %w", err)
	}

	var longPollResp LongPollResponse
	if err := json.Unmarshal(body, &longPollResp); err != nil {
		return nil, fmt.Errorf("can't unmarshal response: %w", err)
	}

	if longPollResp.Failed != 0 {
		if longPollResp.Failed == 4 {
			return nil, fmt.Errorf("long poll error code: %d", longPollResp.Failed)
		}

		ts, err := longPollResp.TS.Int64()
		if err != nil {
			return nil, fmt.Errorf("invalid TS format: %w", err)
		}
		c.setlongPoll(longPoll(c.server, c.key, int(ts)))
		return nil, fmt.Errorf("long poll error code: %d", longPollResp.Failed)
	}

	ts, err := longPollResp.TS.Int64()
	if err != nil {
		return nil, fmt.Errorf("invalid TS format: %w", err)
	}
	c.setlongPoll(longPoll(c.server, c.key, int(ts)))

	return c.parseUpdates(&longPollResp)
}

func (c *Client) parseUpdates(resp *LongPollResponse) ([]LongPollUpdate, error) {
	if resp == nil || len(resp.Updates) == 0 {
		return nil, nil
	}

	var arrayUpdates [][]interface{}
	if err := json.Unmarshal(resp.Updates, &arrayUpdates); err != nil {
		return nil, fmt.Errorf("can't parse updates: %w", err)
	}

	updates := make([]LongPollUpdate, 0, len(arrayUpdates))

	for _, arr := range arrayUpdates {
		if len(arr) < 6 {
			continue
		}

		eventID, ok := arr[0].(float64)
		if !ok || int(eventID) != 4 {
			continue
		}

		flags, ok := arr[2].(float64)
		if !ok || (int(flags)&2) != 0 {
			continue
		}

		peerID, ok := arr[3].(float64)
		if !ok {
			continue
		}

		text := ""
		if len(arr) > 5 {
			if textValue, ok := arr[5].(string); ok {
				text = textValue
			}
		}

		update := LongPollUpdate{
			Type: "message_new",
			Object: Object{
				Text:   text,
				FromId: int(peerID),
			},
		}

		updates = append(updates, update)
	}

	log.Printf("Processed %d message updates", len(updates))
	return updates, nil
}

func (c *Client) SendMessage(ctx context.Context, peerID int, message string) error {
	u, err := url.Parse(fmt.Sprintf("%s/messages.send", c.apiUrl))
	if err != nil {
		return fmt.Errorf("can't parse URL: %w", err)
	}

	q := u.Query()
	q.Add("peer_id", strconv.Itoa(peerID))
	q.Add("message", message)
	q.Add("random_id", strconv.FormatInt(int64(randomInt()), 10))
	q.Add("access_token", c.token)
	q.Add("v", c.apiVersion)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("can't create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("can't execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't read response: %w", err)
	}

	var msgResp MessageResponse
	if err := json.Unmarshal(body, &msgResp); err != nil {
		return fmt.Errorf("can't unmarshal response: %w", err)
	}

	if msgResp.Error != nil {
		return fmt.Errorf("VK API error: %d - %s", msgResp.Error.ErrorCode, msgResp.Error.ErrorMsg)
	}

	return nil
}

func GetLongPollServer(ctx context.Context, token string, version string) (string, string, int, error) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getLongPollServerUrl, nil)
	if err != nil {
		return "", "", 0, fmt.Errorf("can't create request: %w", err)
	}

	q := url.Values{}
	q.Add("access_token", token)
	q.Add("v", version)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", "", 0, fmt.Errorf("can't execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", 0, fmt.Errorf("can't read response: %w", err)
	}

	var response MessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", "", 0, fmt.Errorf("can't unmarshal response: %w", err)
	}

	if response.Error != nil {
		return "", "", 0, fmt.Errorf("VK API error: %d - %s", response.Error.ErrorCode, response.Error.ErrorMsg)
	}

	return response.Response.Server, response.Response.Key, response.Response.TS, nil
}

func (c *Client) setlongPoll(Url string) {
	c.longPollUrl = Url
}

func randomInt() int {
	return int(float64(100000) * float64(randomFloat()))
}

func randomFloat() float64 {
	return float64(time.Now().UnixNano()) / float64(1e9)
}
func longPoll(server, key string, ts int) string {
	return fmt.Sprintf("https://%s?act=a_check&key=%s&ts=%d&wait=20&mode=2&version=2", server, key, ts)
}
