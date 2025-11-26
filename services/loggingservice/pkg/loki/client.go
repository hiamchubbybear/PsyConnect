package loki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/loggingservice/pkg/models"
	"github.com/loggingservice/pkg/settings"
)

// Client represents a Loki client
type Client struct {
	config     settings.LokiConfig
	httpClient *http.Client
	batch      []*models.LogEvent
	mu         sync.Mutex
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// LokiPushRequest represents the Loki push API request format
type LokiPushRequest struct {
	Streams []LokiStream `json:"streams"`
}

// LokiStream represents a single log stream
type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// NewClient creates a new Loki client
func NewClient(config settings.LokiConfig) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		batch:  make([]*models.LogEvent, 0, config.BatchSize),
		stopCh: make(chan struct{}),
	}
}

// Start starts the background batch flusher
func (c *Client) Start() {
	c.wg.Add(1)
	go c.batchFlusher()
}

// Stop stops the client and flushes remaining logs
func (c *Client) Stop() {
	close(c.stopCh)
	c.wg.Wait()
	c.flush() // Final flush
}

// Push adds a log event to the batch
func (c *Client) Push(event *models.LogEvent) error {
	if !c.config.Enabled {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.batch = append(c.batch, event)

	// Flush if batch is full
	if len(c.batch) >= c.config.BatchSize {
		return c.flush()
	}

	return nil
}

// batchFlusher periodically flushes the batch
func (c *Client) batchFlusher() {
	defer c.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			if len(c.batch) > 0 {
				_ = c.flush()
			}
			c.mu.Unlock()
		case <-c.stopCh:
			return
		}
	}
}

// flush sends the current batch to Loki
func (c *Client) flush() error {
	if len(c.batch) == 0 {
		return nil
	}

	// Group logs by labels (stream)
	streams := make(map[string]*LokiStream)

	for _, event := range c.batch {
		labels := event.GetLabels()
		streamKey := labelsToKey(labels)

		if stream, exists := streams[streamKey]; exists {
			// Add to existing stream
			stream.Values = append(stream.Values, []string{
				strconv.FormatInt(time.Now().UnixNano(), 10),
				event.GetLogLine(),
			})
		} else {
			// Create new stream
			streams[streamKey] = &LokiStream{
				Stream: labels,
				Values: [][]string{
					{
						strconv.FormatInt(time.Now().UnixNano(), 10),
						event.GetLogLine(),
					},
				},
			}
		}
	}

	// Convert map to slice
	streamSlice := make([]LokiStream, 0, len(streams))
	for _, stream := range streams {
		streamSlice = append(streamSlice, *stream)
	}

	// Create push request
	pushReq := LokiPushRequest{
		Streams: streamSlice,
	}

	// Send to Loki
	if err := c.sendToLoki(pushReq); err != nil {
		return err
	}

	// Clear batch
	c.batch = c.batch[:0]
	return nil
}

// sendToLoki sends the push request to Loki
func (c *Client) sendToLoki(pushReq LokiPushRequest) error {
	jsonData, err := json.Marshal(pushReq)
	if err != nil {
		return fmt.Errorf("failed to marshal push request: %w", err)
	}

	url := fmt.Sprintf("%s/loki/api/v1/push", c.config.URL)
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add basic auth if configured
	if c.config.Username != "" && c.config.Password != "" {
		req.SetBasicAuth(c.config.Username, c.config.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("loki returned non-2xx status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// labelsToKey converts labels map to a unique string key
func labelsToKey(labels map[string]string) string {
	// Simple concatenation for demo, in production use proper hashing
	key := ""
	for k, v := range labels {
		key += k + "=" + v + ","
	}
	return key
}
