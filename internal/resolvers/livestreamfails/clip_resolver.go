package livestreamfails

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

type ClipResolver struct {
	clipCache cache.Cache
}

func (r *ClipResolver) Check(ctx context.Context, url *url.URL) (context.Context, bool) {
	if !utils.IsSubdomainOf(url, "livestreamfails.com") {
		return ctx, false
	}

	match := pathRegex.FindStringSubmatch(url.Path)
	if len(match) != 2 {
		return ctx, false
	}

	ctx = contextWithClipID(ctx, match[1])

	return ctx, true
}

func (r *ClipResolver) Run(ctx context.Context, url *url.URL, req *http.Request) (*cache.Response, error) {
	clipID, err := clipIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return r.clipCache.Get(ctx, clipID, req)
}

func (r *ClipResolver) Name() string {
	return "livestreamfails:clip"
}

func NewClipResolver(ctx context.Context, cfg config.APIConfig, pool db.Pool, apiURLFormat string) *ClipResolver {
	clipLoader := &ClipLoader{
		apiURLFormat: apiURLFormat,
	}

	r := &ClipResolver{
		clipCache: cache.NewPostgreSQLCache(
			ctx, cfg, pool, cache.NewPrefixKeyProvider("livestreamfails:clip"),
			resolver.NewResponseMarshaller(clipLoader), cfg.LivestreamfailsClipCacheDuration),
	}

	return r
}
