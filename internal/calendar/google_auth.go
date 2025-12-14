package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	appconfig "github.com/Alvkoen/barely-incharge/internal/config"
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
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close token file: %v\n", closeErr)
		}
	}()

	return json.NewEncoder(f).Encode(token)
}

func loadToken() (*oauth2.Token, error) {
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close token file: %v\n", closeErr)
		}
	}()

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

func GetClient(ctx context.Context) (*calendar.Service, error) {
	credentials, err := loadCredentials()
	if err != nil {
		return nil, err
	}

	oauthConfig, err := google.ConfigFromJSON(credentials, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %w", err)
	}

	token, err := loadToken()
	if err != nil {
		token, err = getTokenFromWeb(oauthConfig)
		if err != nil {
			return nil, err
		}
		if err := saveToken(token); err != nil {
			return nil, err
		}
	}

	baseHTTPClient := &http.Client{
		Timeout: appconfig.HTTPTimeout,
	}
	ctxWithClient := context.WithValue(ctx, oauth2.HTTPClient, baseHTTPClient)
	tokenSource := oauthConfig.TokenSource(ctxWithClient, token)
	autoSaveSource := &autoSaveTokenSource{source: tokenSource}
	httpClient := oauth2.NewClient(ctxWithClient, autoSaveSource)

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

	if err := saveToken(token); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save refreshed token: %v\n", err)
	}

	return token, nil
}
