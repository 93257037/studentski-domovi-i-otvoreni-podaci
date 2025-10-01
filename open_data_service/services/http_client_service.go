package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClientService handles HTTP communication with other services
type HTTPClientService struct {
	client *http.Client
	baseURL string
}

// NewHTTPClientService creates a new HTTP client service
func NewHTTPClientService(baseURL string) *HTTPClientService {
	return &HTTPClientService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// PrihvacenaAplikacijaResponse represents the response from st_dom_service
type PrihvacenaAplikacijaResponse struct {
	PrihvaceneAplikacije []PrihvacenaAplikacija `json:"prihvacene_aplikacije"`
	Count                int                    `json:"count"`
}

// PrihvacenaAplikacija represents an accepted application from st_dom_service
type PrihvacenaAplikacija struct {
	ID             string `json:"id"`
	AplikacijaID   string `json:"aplikacija_id"`
	UserID         string `json:"user_id"`
	BrojIndexa     string `json:"broj_indexa"`
	Prosek         int    `json:"prosek"`
	SobaID         string `json:"soba_id"`
	AcademicYear   string `json:"academic_year"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// GetPrihvaceneAplikacije calls the st_dom_service to get all accepted applications
func (h *HTTPClientService) GetPrihvaceneAplikacije(authHeader string) (*PrihvacenaAplikacijaResponse, error) {
	url := fmt.Sprintf("%s/api/v1/prihvacene_aplikacije/", h.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response PrihvacenaAplikacijaResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetPrihvaceneAplikacijeForUser calls the st_dom_service to get accepted applications for a specific user
func (h *HTTPClientService) GetPrihvaceneAplikacijeForUser(userID string, authHeader string) (*PrihvacenaAplikacijaResponse, error) {
	url := fmt.Sprintf("%s/api/v1/prihvacene_aplikacije/user/%s", h.baseURL, userID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response PrihvacenaAplikacijaResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetPrihvaceneAplikacijeForRoom calls the st_dom_service to get accepted applications for a specific room
func (h *HTTPClientService) GetPrihvaceneAplikacijeForRoom(roomID string, authHeader string) (*PrihvacenaAplikacijaResponse, error) {
	url := fmt.Sprintf("%s/api/v1/prihvacene_aplikacije/room/%s", h.baseURL, roomID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response PrihvacenaAplikacijaResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetPrihvaceneAplikacijeForAcademicYear calls the st_dom_service to get accepted applications for a specific academic year
func (h *HTTPClientService) GetPrihvaceneAplikacijeForAcademicYear(academicYear string, authHeader string) (*PrihvacenaAplikacijaResponse, error) {
	url := fmt.Sprintf("%s/api/v1/prihvacene_aplikacije/academic_year/%s", h.baseURL, academicYear)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response PrihvacenaAplikacijaResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// HealthCheck checks if the st_dom_service is available
func (h *HTTPClientService) HealthCheck() error {
	url := fmt.Sprintf("%s/health", h.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}
