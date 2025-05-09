package resolver

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lonelyshoeh/api/internal/logger"
	"github.com/lonelyshoeh/api/internal/version"
)

var (
	httpClient = &http.Client{
		Timeout: 15 * time.Second,
	}
)

func RequestGET(ctx context.Context, url string) (response *http.Response, err error) {
	log := logger.FromContext(ctx)

	log.Debugw("[resolver] GET",
		"url", url,
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// ensures websites return pages in english (e.g. twitter would return french preview
	// when the request came from a french IP.)
	req.Header.Add("Accept-Language", "en-US, en;q=0.9, *;q=0.5")
	req.Header.Set("User-Agent", fmt.Sprintf("chatterino-api-cache/%s link-resolver", version.Version))

	return httpClient.Do(req)
}

func RequestGETWithHeaders(url string, extraHeaders map[string]string) (response *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// ensures websites return pages in english (e.g. twitter would return french preview
	// when the request came from a french IP.)
	req.Header.Add("Accept-Language", "en-US, en;q=0.9, *;q=0.5")
	req.Header.Set("User-Agent", fmt.Sprintf("chatterino-api-cache/%s link-resolver", version.Version))

	for headerKey, headerValue := range extraHeaders {
		req.Header.Set(headerKey, headerValue)
	}

	return httpClient.Do(req)
}

func RequestPOST(url, body string) (response *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("chatterino-api-cache/%s link-resolver", version.Version))

	return httpClient.Do(req)
}

func HTTPClient() *http.Client {
	return httpClient
}
