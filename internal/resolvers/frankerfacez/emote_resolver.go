package frankerfacez

import (
	"context"
	"net/http"
	"net/url"

	"github.com/lonelyshoeh/api/internal/db"
	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/config"
	"github.com/lonelyshoeh/api/pkg/resolver"
)

type EmoteResolver struct {
	emoteCache cache.Cache
}

func (r *EmoteResolver) Check(ctx context.Context, url *url.URL) (context.Context, bool) {
	if match, _ := resolver.MatchesHosts(url, domains); !match {
		return ctx, false
	}

	if !emotePathRegex.MatchString(url.Path) {
		return ctx, false
	}

	return ctx, true
}

func (r *EmoteResolver) Run(ctx context.Context, url *url.URL, req *http.Request) (*cache.Response, error) {
	matches := emotePathRegex.FindStringSubmatch(url.Path)
	if len(matches) != 4 {
		return nil, errInvalidFrankerFaceZEmotePath
	}

	emoteHash := matches[1]

	return r.emoteCache.Get(ctx, emoteHash, req)
}

func (r *EmoteResolver) Name() string {
	return "frankerfacez:emote"
}

func NewEmoteResolver(ctx context.Context, cfg config.APIConfig, pool db.Pool, emoteAPIURL *url.URL) *EmoteResolver {
	emoteLoader := NewEmoteLoader(emoteAPIURL)

	r := &EmoteResolver{
		emoteCache: cache.NewPostgreSQLCache(
			ctx, cfg, pool, cache.NewPrefixKeyProvider("frankerfacez:emote"),
			resolver.NewResponseMarshaller(emoteLoader), cfg.FfzEmoteCacheDuration),
	}

	return r
}
