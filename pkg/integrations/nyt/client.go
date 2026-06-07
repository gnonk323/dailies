package nyt

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// GameFetcher defines how an individual NYT game resolves its metadata
type GameFetcher interface {
	GameKey() string
	FetchGame(ctx context.Context, client *NYTClient, dateStr string, cookies string) (interface{}, error)
}

// NYTClient encapsulates common authentication and networking patterns
type NYTClient struct {
	HTTPClient *http.Client
	Cookies    string
}

func NewNYTClient(cookies string) *NYTClient {
	return &NYTClient{
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
		Cookies:    cookies,
	}
}

// ExtractRegiID pulls the user's ID out of the NYT raw cookie block
func (c *NYTClient) ExtractRegiID() (string, error) {
	re := regexp.MustCompile(`regi_id=(\d+)`)
	matches := re.FindStringSubmatch(c.Cookies)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not locate 'regi_id' inside the configuration cookies")
	}
	return matches[1], nil
}

// DoRequest handles standard header decoration for NYT endpoints
func (c *NYTClient) DoRequest(ctx context.Context, method, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "dailies-cli")
	if c.Cookies != "" {
		req.Header.Set("Cookie", c.Cookies)
	}

	// Attach proper cookie types if raw string contains semi-colons
	if c.Cookies != "" && !strings.Contains(c.Cookies, "=") {
		return nil, fmt.Errorf("invalid cookie syntax string provided")
	}

	return c.HTTPClient.Do(req)
}
