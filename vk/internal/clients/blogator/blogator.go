package blogator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: http.Client{},
	}
}

func (c *Client) GetNewItems(ctx context.Context, since time.Time) ([]DataItem, error) {
	url := fmt.Sprintf("http://%s?since=%s", c.baseURL, since.Format(time.RFC3339))
	log.Printf("Making request to Blogator API: %s", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	log.Printf("Blogator API response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res Response
	
	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		log.Printf("Raw response body: %s", string(body))
		return nil, fmt.Errorf("can't unmarshall: %w", err)
	}

	log.Printf("Successfully unmarshaled %d items", len(res.Data))

	return res.Data, nil
}
