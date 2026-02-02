package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type AmadeusService struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	apiSecret string

	tokenMutex  sync.RWMutex
	accessToken string
	tokenExpiry time.Time
}

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewAmadeusService() *AmadeusService {
	return &AmadeusService{
		client:    &http.Client{Timeout: 45 * time.Second},
		baseURL:   os.Getenv("AMADEUS_BASE_URL"),
		apiKey:    os.Getenv("AMADEUS_API_KEY"),
		apiSecret: os.Getenv("AMADEUS_API_SECRET"),
	}
}

func (as *AmadeusService) FetchToken() error {
	as.tokenMutex.Lock()
	defer as.tokenMutex.Unlock()

	slog.Info("Requesting new Amadeus OAuth2 token")

	if as.apiKey == "" || as.apiSecret == "" {
		slog.Error("Amadeus credentials missing in environment variables")
		return fmt.Errorf("amadeus api key or secret are not in .env")
	}

	tokenURL := as.baseURL + "/v1/security/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", as.apiKey)
	data.Set("client_secret", as.apiSecret)

	resp, err := as.client.Post(
		tokenURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		slog.Error("Amadeus token request HTTP error", "error", err)
		return fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Amadeus token request failed", "status", resp.StatusCode)
		return fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		slog.Error("Failed to decode Amadeus token response", "error", err)
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	as.accessToken = tokenResp.AccessToken
	as.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	slog.Info("Amadeus token renewed successfully", "expires_in", tokenResp.ExpiresIn)
	return nil
}

func (as *AmadeusService) GetToken() (string, error) {
	as.tokenMutex.RLock()
	if as.accessToken != "" && time.Now().Before(as.tokenExpiry.Add(-5*time.Minute)) {
		defer as.tokenMutex.RUnlock()
		return as.accessToken, nil
	}
	as.tokenMutex.RUnlock()

	slog.Debug("Amadeus token expired or missing, fetching new one")
	if err := as.FetchToken(); err != nil {
		return "", fmt.Errorf("failed to renew amadeus token: %w", err)
	}

	return as.accessToken, nil
}

func (as *AmadeusService) ExecuteGetRequest(endpoint string, params url.Values) (*http.Response, error) {
	l := slog.With("endpoint", endpoint)

	token, err := as.GetToken()
	if err != nil {
		l.Error("Amadeus authentication failed", "error", err)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	apiURL := as.baseURL + endpoint
	l.Debug("Executing Amadeus GET request", "params", params.Encode())

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		l.Error("Failed to create Amadeus request object", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := as.client.Do(req)
	if err != nil {
		l.Error("Amadeus API execution error", "error", err)
		return nil, err
	}

	if resp.StatusCode >= 400 {
		l.Warn("Amadeus API returned error status", "status", resp.StatusCode)
	}

	return resp, nil
}
