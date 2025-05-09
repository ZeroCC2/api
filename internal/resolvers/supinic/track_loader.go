package supinic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/humanize"
	"github.com/lonelyshoeh/api/pkg/resolver"
)

type TooltipData struct {
	ID         int
	Name       string
	AuthorName string
	Tags       string
	Duration   string
}

type TrackData struct {
	ID          int       `json:"id"`
	Link        string    `json:"code"` // Youtube ID/link
	Name        string    `json:"name"`
	VideoType   int       `json:"videoType"`
	TrackType   string    `json:"trackType"`
	Duration    float32   `json:"duration"`
	Available   bool      `json:"available"`
	PublishedAt time.Time `json:"published"`
	Notes       string    `json:"notes"`
	AddedBy     string    `json:"addedBy"`
	ParsedLink  string    `json:"parsedLink"`
	Tags        []string  `json:"tags"`
	Authors     []struct {
		ID   int    `json:"ID"`
		Name string `json:"name"`
		Role string `json:"role"`
	} `json:"authors"`
}

type TrackListAPIResponse struct {
	Data TrackData `json:"data"`
}

type TrackLoader struct {
}

func (l *TrackLoader) Load(ctx context.Context, rawTrackID string, r *http.Request) (*resolver.Response, time.Duration, error) {
	trackID, _ := strconv.ParseInt(rawTrackID, 10, 32)
	apiURL := fmt.Sprintf(trackListAPIURL, trackID)

	// Execute Track list API request
	resp, err := resolver.RequestGET(ctx, apiURL)
	if err != nil {
		return &resolver.Response{
			Status:  http.StatusInternalServerError,
			Message: "Track list http request error " + resolver.CleanResponse(err.Error()),
		}, cache.NoSpecialDur, nil
	}
	defer resp.Body.Close()

	// Error out if the track isn't found or something else went wrong with the request
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusMultipleChoices {
		return trackNotFoundResponse, cache.NoSpecialDur, nil
	}

	// Read response into a string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &resolver.Response{
			Status:  http.StatusInternalServerError,
			Message: "Track list http body read error " + resolver.CleanResponse(err.Error()),
		}, cache.NoSpecialDur, nil
	}

	// Parse response into a predefined JSON blob (see TrackListAPIResponse struct above)
	var jsonResponse TrackListAPIResponse
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return &resolver.Response{
			Status:  http.StatusInternalServerError,
			Message: "Track list api unmarshal error " + resolver.CleanResponse(err.Error()),
		}, cache.NoSpecialDur, nil
	}
	if jsonResponse.Data.ID == 0 { // API responds with {..., "data": null} if nothing was found
		return trackNotFoundResponse, cache.NoSpecialDur, nil
	}

	trackData := jsonResponse.Data

	prettyAuthors := ""
	for i, elem := range trackData.Authors {
		if i != 0 {
			prettyAuthors += ", "
		}
		prettyAuthors += fmt.Sprintf("%s (ID %d - %s)", elem.Name, elem.ID, elem.Role)

	}

	// API returned no authors.
	if prettyAuthors == "" {
		prettyAuthors = "unknown"
	}

	// Build tooltip data from the API response
	data := TooltipData{
		ID:         trackData.ID,
		Name:       trackData.Name,
		AuthorName: prettyAuthors,
		Tags:       strings.Join(trackData.Tags, ", "),
		Duration:   humanize.Duration(time.Duration(trackData.Duration) * time.Second),
	}

	// Build a tooltip using the tooltip template (see tooltipTemplate) with the data we massaged above
	var tooltip bytes.Buffer
	if err := trackListTemplate.Execute(&tooltip, data); err != nil {
		return &resolver.Response{
			Status:  http.StatusInternalServerError,
			Message: "Track list template error " + resolver.CleanResponse(err.Error()),
		}, cache.NoSpecialDur, nil
	}

	return &resolver.Response{
		Status:  200,
		Tooltip: url.PathEscape(tooltip.String()),
		//Thumbnail: thumbnailURL,
	}, cache.NoSpecialDur, nil

}
