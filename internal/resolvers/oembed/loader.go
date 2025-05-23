package oembed

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/humanize"
	"github.com/lonelyshoeh/api/pkg/resolver"
	"github.com/dyatlov/go-oembed/oembed"
)

type Loader struct {
	oEmbed                 *oembed.Oembed
	facebookAppAccessToken string
}

func (l *Loader) Load(ctx context.Context, requestedURL string, r *http.Request) (*resolver.Response, time.Duration, error) {
	extraOpts := url.Values{}

	item := l.oEmbed.FindItem(requestedURL)

	if item.ProviderName == "Facebook" || item.ProviderName == "Instagram" {
		// Add facebook token if it exists
		if l.facebookAppAccessToken != "" {
			extraOpts.Set("access_token", l.facebookAppAccessToken)
			extraOpts.Set("omitscript", "true")
		}
	}

	data, err := item.FetchOembed(oembed.Options{
		URL:       requestedURL,
		ExtraOpts: extraOpts,
	})

	if err != nil {
		return &resolver.Response{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong loading this oEmbed.\noEmbed error: " + resolver.CleanResponse(err.Error()),
		}, cache.NoSpecialDur, nil
	}

	if data.Status > http.StatusOK {
		log.Printf("[oEmbed] Skipping url %s because status code is %d\n", requestedURL, data.Status)
		return &resolver.Response{
			Status:  data.Status,
			Message: fmt.Sprintf("This oEmbed couldn't be loaded in.\noEmbed status code: %d", data.Status),
		}, cache.NoSpecialDur, nil
	}

	infoTooltipData := oEmbedData{data, requestedURL}

	infoTooltipData.Title = humanize.Title(infoTooltipData.Title)
	infoTooltipData.Description = humanize.Description(infoTooltipData.Description)
	infoTooltipData.RequestedURL = requestedURL

	// Build a tooltip using the tooltip template (see tooltipTemplate) with the data we massaged above
	var tooltip bytes.Buffer
	if err := oEmbedTemplate.Execute(&tooltip, infoTooltipData); err != nil {
		return &resolver.Response{
			Status:  http.StatusInternalServerError,
			Message: "oEmbed template error: " + resolver.CleanResponse(err.Error()),
		}, cache.NoSpecialDur, nil
	}

	resolverResponse := resolver.Response{
		Status:  http.StatusOK,
		Tooltip: url.PathEscape(tooltip.String()),
	}

	if infoTooltipData.Type == "photo" {
		resolverResponse.Thumbnail = infoTooltipData.URL
	}

	if infoTooltipData.ThumbnailURL != "" {

		// Some thumbnail URLs, like Streamable's returns // with no schema.
		if strings.HasPrefix(infoTooltipData.ThumbnailURL, "//") {
			infoTooltipData.ThumbnailURL = "https:" + infoTooltipData.ThumbnailURL
		}

		resolverResponse.Thumbnail = infoTooltipData.ThumbnailURL
	}

	return &resolverResponse, cache.NoSpecialDur, nil
}
