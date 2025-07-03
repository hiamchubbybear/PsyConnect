package external

import (
	"bytes"
	"consultationservice/internal/dto"
	"consultationservice/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RecommendationApi(data dto.FilterRawData) (*model.ClientSwipes, error) {
	url := "http://127.0.0.1:5000/recommend"
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
	var user model.ClientSwipes
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &user, nil
}
