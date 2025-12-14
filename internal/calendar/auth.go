package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const tokenFile = "token.json"

func loadCredentials() ([]byte, error) {
	credFile := "credentials.json"
	b, err := os.ReadFile(credFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %w (make sure %s exists)", err, credFile)
	}
	return b, nil
}

func saveToken(token *oauth2.Token) error {
	f, err := os.Create(tokenFile)
	if err != nil {
		return fmt.Errorf("unable to create token file: %w", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}

func loadToken() (*oauth2.Token, error) {
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)
	fmt.Print("Enter authorization code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %w", err)
	}

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %w", err)
	}

	return token, nil
}

// GetClient retrieves an authenticated HTTP client for Google Calendar API
func GetClient(ctx context.Context) (*calendar.Service, error) {
	// Load OAuth2 credentials
	credentials, err := loadCredentials()
	if err != nil {
		return nil, err
	}

	// Parse credentials and create OAuth2 config
	config, err := google.ConfigFromJSON(credentials, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %w", err)
	}

	// Try to load existing token
	token, err := loadToken()
	if err != nil {
		// No token found, get one from web
		token, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		// Save token for future use
		if err := saveToken(token); err != nil {
			return nil, err
		}
	}

	// Create token source that auto-refreshes and auto-saves
	tokenSource := config.TokenSource(ctx, token)
	autoSaveSource := &autoSaveTokenSource{source: tokenSource}

	// Create authenticated HTTP client
	httpClient := oauth2.NewClient(ctx, autoSaveSource)

	// Create Calendar service
	service, err := calendar.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("unable to create calendar service: %w", err)
	}

	return service, nil
}

type autoSaveTokenSource struct {
	source oauth2.TokenSource
}

func (a *autoSaveTokenSource) Token() (*oauth2.Token, error) {
	token, err := a.source.Token()
	if err != nil {
		return nil, err
	}

	// Save the token (in case it was refreshed)
	if err := saveToken(token); err != nil {
		// Log but don't fail - token refresh worked, save is just optimization
		fmt.Printf("Warning: failed to save refreshed token: %v\n", err)
	}

	return token, nil
}
