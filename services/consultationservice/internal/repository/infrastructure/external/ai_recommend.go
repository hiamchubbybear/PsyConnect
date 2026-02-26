package external

import (
	"bytes"
	"consultationservice/internal/dto"
	"consultationservice/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Deprecated:Replace RecommendationApi instead
func RecommendationApi(data dto.FilterRawData) ([]model.ClientSwipe, error) {
	url := os.Getenv("RECOMMENDATION_SERVICE_URL")
	if url == "" {
		url = "http://127.0.0.1:8086/recommend"
	}
	jsonData, err := json.Marshal(data)
	log.Printf(string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recommendation API error: %s", res.Status)
	}
	var response struct {
		ClientId string              `json:"client_id"`
		Swipes   []model.ClientSwipe `json:"swipes"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return response.Swipes, nil
}

func RecommendationApiV1(data dto.FilterRawDataV1) ([]model.ClientSwipeV1, error) {
	url := os.Getenv("RECOMMENDATION_SERVICE_URL")
	if url == "" {
		url = "http://127.0.0.1:8086/recommend"
	}
	jsonData, err := json.Marshal(data)
	log.Printf(string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recommendation API error: %s", res.Status)
	}
	var response struct {
		ClientId string                `json:"client_id"`
		Swipes   []model.ClientSwipeV1 `json:"swipes"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return response.Swipes, nil
}
