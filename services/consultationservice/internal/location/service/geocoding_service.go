package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SearchResult struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ReverseResult struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type cacheItem[T any] struct {
	value     T
	expiresAt time.Time
}

type GeocodingService struct {
	client       *http.Client
	userAgent    string
	searchCache  map[string]cacheItem[[]SearchResult]
	reverseCache map[string]cacheItem[ReverseResult]
	mu           sync.RWMutex
}

type nominatimSearchItem struct {
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
}

type nominatimReverseItem struct {
	DisplayName string `json:"display_name"`
}

func NewGeocodingService() *GeocodingService {
	return &GeocodingService{
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
		userAgent:    "PsyConnect/1.0 (consultationservice geocoding proxy)",
		searchCache:  make(map[string]cacheItem[[]SearchResult]),
		reverseCache: make(map[string]cacheItem[ReverseResult]),
	}
}

func (s *GeocodingService) Search(ctx context.Context, query string) ([]SearchResult, error) {
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	if normalizedQuery == "" {
		return nil, fmt.Errorf("query is required")
	}

	if cached, ok := s.getSearchCache(normalizedQuery); ok {
		return cached, nil
	}

	endpoint := "https://nominatim.openstreetmap.org/search?format=jsonv2&limit=6&addressdetails=1&accept-language=vi,en&q=" +
		url.QueryEscape(query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding provider returned status %d", resp.StatusCode)
	}

	var payload []nominatimSearchItem
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(payload))
	for _, item := range payload {
		lat, err := parseCoordinate(item.Lat)
		if err != nil {
			continue
		}
		lon, err := parseCoordinate(item.Lon)
		if err != nil {
			continue
		}

		results = append(results, SearchResult{
			Address:   item.DisplayName,
			Latitude:  lat,
			Longitude: lon,
		})
	}

	s.setSearchCache(normalizedQuery, results)
	return results, nil
}

func (s *GeocodingService) Reverse(ctx context.Context, latitude, longitude float64) (ReverseResult, error) {
	cacheKey := fmt.Sprintf("%.6f,%.6f", latitude, longitude)
	if cached, ok := s.getReverseCache(cacheKey); ok {
		return cached, nil
	}

	endpoint := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/reverse?format=jsonv2&accept-language=vi,en&lat=%f&lon=%f",
		latitude,
		longitude,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return ReverseResult{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return ReverseResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ReverseResult{}, fmt.Errorf("reverse geocoding provider returned status %d", resp.StatusCode)
	}

	var payload nominatimReverseItem
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return ReverseResult{}, err
	}

	result := ReverseResult{
		Address:   payload.DisplayName,
		Latitude:  latitude,
		Longitude: longitude,
	}
	if result.Address == "" {
		result.Address = cacheKey
	}

	s.setReverseCache(cacheKey, result)
	return result, nil
}

func (s *GeocodingService) getSearchCache(key string) ([]SearchResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.searchCache[key]
	if !ok || time.Now().After(item.expiresAt) {
		return nil, false
	}
	return item.value, true
}

func (s *GeocodingService) setSearchCache(key string, value []SearchResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.searchCache[key] = cacheItem[[]SearchResult]{
		value:     value,
		expiresAt: time.Now().Add(10 * time.Minute),
	}
}

func (s *GeocodingService) getReverseCache(key string) (ReverseResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.reverseCache[key]
	if !ok || time.Now().After(item.expiresAt) {
		return ReverseResult{}, false
	}
	return item.value, true
}

func (s *GeocodingService) setReverseCache(key string, value ReverseResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reverseCache[key] = cacheItem[ReverseResult]{
		value:     value,
		expiresAt: time.Now().Add(10 * time.Minute),
	}
}

func parseCoordinate(raw string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(raw), 64)
}
