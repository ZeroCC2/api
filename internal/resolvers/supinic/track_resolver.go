package supinic

import (
	"context"
	"net/http"
	"net/url"

	"github.com/lonelyshoeh/api/internal/db"
	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/config"
	"github.com/lonelyshoeh/api/pkg/resolver"
	"github.com/lonelyshoeh/api/pkg/utils"
)

type TrackResolver struct {
	trackCache cache.Cache
}

func (r *TrackResolver) Check(ctx context.Context, url *url.URL) (context.Context, bool) {
	if !utils.IsDomains(url, trackListDomains) {
		return ctx, false
	}

	if !trackPathRegex.MatchString(url.Path) {
		return ctx, false
	}

	return ctx, true
}

func (r *TrackResolver) Run(ctx context.Context, url *url.URL, req *http.Request) (*cache.Response, error) {
	matches := trackPathRegex.FindStringSubmatch(url.Path)
	if len(matches) != 2 {
		return nil, errInvalidTrackPath
	}

	trackID := matches[1]

	return r.trackCache.Get(ctx, trackID, req)
}

func (r *TrackResolver) Name() string {
	return "supinic:track"
}

func NewTrackResolver(ctx context.Context, cfg config.APIConfig, pool db.Pool) *TrackResolver {
	trackLoader := &TrackLoader{}

	r := &TrackResolver{
		trackCache: cache.NewPostgreSQLCache(
			ctx, cfg, pool, cache.NewPrefixKeyProvider("supinic:track"),
			resolver.NewResponseMarshaller(trackLoader), cfg.SupinicTrackCacheDuration),
	}

	return r
}
