package twitch

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/lonelyshoeh/api/internal/db"
	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/config"
	"github.com/lonelyshoeh/api/pkg/resolver"
	"github.com/lonelyshoeh/api/pkg/utils"
)

var userRegex = regexp.MustCompile(`^\/([a-zA-Z0-9_]+)$`)
var ignoredUsers = []string{
	"inventory",
	"popout",
	"subscriptions",
	"videos",
	"following",
	"directory",
	"moderator",
}

type UserResolver struct {
	userCache cache.Cache
}

func (r *UserResolver) Check(ctx context.Context, url *url.URL) (context.Context, bool) {
	if !utils.IsDomains(url, userDomains) {
		return ctx, false
	}

	userMatch := userRegex.FindStringSubmatch(url.Path)
	if len(userMatch) != 2 {
		return ctx, false
	}

	for _, ignoredUser := range ignoredUsers {
		if ignoredUser == strings.ToLower(userMatch[1]) {
			return ctx, false
		}
	}

	return ctx, true
}

func (r *UserResolver) Run(ctx context.Context, url *url.URL, req *http.Request) (*cache.Response, error) {
	return r.userCache.Get(ctx, strings.ToLower(strings.TrimLeft(url.Path, "/")), req)
}

func (r *UserResolver) Name() string {
	return "twitch:user"
}

func NewUserResolver(ctx context.Context, cfg config.APIConfig, pool db.Pool, helixAPI TwitchAPIClient) *UserResolver {
	userLoader := &UserLoader{helixAPI: helixAPI}

	r := &UserResolver{
		userCache: cache.NewPostgreSQLCache(ctx, cfg, pool, cache.NewPrefixKeyProvider("twitch:user"),
			resolver.NewResponseMarshaller(userLoader), cfg.TwitchUsernameCacheDuration),
	}

	return r
}
