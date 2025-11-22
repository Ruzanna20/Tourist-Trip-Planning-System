package services

import (
	"encoding/json"
	"fmt"
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

	if as.apiKey == "" || as.apiSecret == "" {
		return fmt.Errorf("amadues api key or secret are not in .env")
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
		return fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	as.accessToken = tokenResp.AccessToken
	as.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return nil
}

func (as *AmadeusService) GetToken() (string, error) {
	as.tokenMutex.RLock()
	if as.accessToken != "" && time.Now().Before(as.tokenExpiry.Add(-5*time.Minute)) {
		defer as.tokenMutex.RUnlock()
		return as.accessToken, nil
	}
	as.tokenMutex.RUnlock()
	if err := as.FetchToken(); err != nil {
		return "", fmt.Errorf("failed to renew amadues token: %w", err)
	}

	return as.accessToken, nil
}

func (as *AmadeusService) ExecuteGetRequest(endpoint string, Params url.Values) (*http.Response, error) {
	token, err := as.GetToken()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	apiURl := as.baseURL + endpoint
	req, err := http.NewRequest("GET", apiURl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.URL.RawQuery = Params.Encode()
	req.Header.Add("Authorization", "Bearer "+token)

	return as.client.Do(req)
}
