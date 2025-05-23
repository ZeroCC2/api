package youtube

import (
	"context"
	"net/http"
	"net/url"
	"regexp"

	"github.com/lonelyshoeh/api/internal/db"
	"github.com/lonelyshoeh/api/internal/logger"
	"github.com/lonelyshoeh/api/internal/staticresponse"
	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/config"
	"github.com/lonelyshoeh/api/pkg/utils"
	youtubeAPI "google.golang.org/api/youtube/v3"
)

var youtubeChannelRegex = regexp.MustCompile(`^/(c\/|channel\/|user\/)?([a-zA-Z0-9\-]{1,})$`)

type YouTubeChannelResolver struct {
	channelCache cache.Cache
}

func (r *YouTubeChannelResolver) Check(ctx context.Context, url *url.URL) (context.Context, bool) {
	if !utils.IsSubdomainOf(url, "youtube.com") {
		return ctx, false
	}

	if url.Path == "/results" {
		return ctx, false
	}

	q := url.Query()
	// TODO(go1.18): Replace with q.Has("v") once we've transitioned to at least go 1.17 as least supported version
	if q.Has("v") {
		return ctx, false
	}

	matches := youtubeChannelRegex.MatchString(url.Path)
	return ctx, matches
}

func (r *YouTubeChannelResolver) Run(ctx context.Context, url *url.URL, req *http.Request) (*cache.Response, error) {
	log := logger.FromContext(ctx)
	channel := getChannelFromPath(url.Path)

	if channel.Type == InvalidChannel {
		log.Warnw("[YouTube] URL was incorrectly treated as a channel",
			"url", url,
		)

		return &staticresponse.RNoLinkInfoFound, nil
	}

	return r.channelCache.Get(ctx, channel.ToCacheKey(), req)
}

func (r *YouTubeChannelResolver) Name() string {
	return "youtube:channel"
}

func NewYouTubeChannelResolver(ctx context.Context, cfg config.APIConfig, pool db.Pool, youtubeClient *youtubeAPI.Service) *YouTubeChannelResolver {
	loader := NewYouTubeChannelLoader(youtubeClient)

	r := &YouTubeChannelResolver{
		channelCache: cache.NewPostgreSQLCache(
			ctx, cfg, pool, cache.NewPrefixKeyProvider("youtube:channel"), loader, cfg.YoutubeChannelCacheDuration,
		),
	}

	return r
}
