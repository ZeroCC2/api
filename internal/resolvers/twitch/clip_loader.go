package twitch

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/lonelyshoeh/api/internal/logger"
	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/humanize"
	"github.com/lonelyshoeh/api/pkg/resolver"
	"github.com/nicklaw5/helix"
)

type twitchClipsTooltipData struct {
	Title        string
	AuthorName   string
	ChannelName  string
	Duration     string
	CreationDate string
	Views        string
}

type ClipLoader struct {
	helixAPI TwitchAPIClient
}

func (l *ClipLoader) Load(ctx context.Context, clipSlug string, r *http.Request) (*resolver.Response, time.Duration, error) {
	log := logger.FromContext(ctx)

	log.Debugw("[Twitch] Get clip",
		"clipSlug", clipSlug,
	)

	response, err := l.helixAPI.GetClips(&helix.ClipsParams{IDs: []string{clipSlug}})
	if err != nil {
		log.Errorw("[Twitch] Error getting clip",
			"clipSlug", clipSlug,
			"error", err,
		)

		return resolver.Errorf("Twitch clip load error: %s", err)
	}

	if len(response.Data.Clips) != 1 {
		return noTwitchClipWithThisIDFound, cache.NoSpecialDur, nil
	}

	var clip = response.Data.Clips[0]

	data := twitchClipsTooltipData{
		Title:        clip.Title,
		AuthorName:   clip.CreatorName,
		ChannelName:  clip.BroadcasterName,
		Duration:     humanize.DurationSeconds(time.Duration(clip.Duration) * time.Second),
		CreationDate: humanize.CreationDateRFC3339(clip.CreatedAt),
		Views:        humanize.Number(uint64(clip.ViewCount)),
	}

	var tooltip bytes.Buffer
	if err := twitchClipsTooltip.Execute(&tooltip, data); err != nil {
		return resolver.Errorf("Twitch clip template error: %s", err)
	}

	return &resolver.Response{
		Status:    200,
		Tooltip:   url.PathEscape(tooltip.String()),
		Thumbnail: clip.ThumbnailURL,
	}, cache.NoSpecialDur, nil
}
