package twitter

import (
	"context"
	"testing"

	"github.com/lonelyshoeh/api/internal/logger"
	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/config"
	"github.com/lonelyshoeh/api/pkg/resolver"
	"github.com/pashagolub/pgxmock"

	qt "github.com/frankban/quicktest"
)

func TestInitialize(t *testing.T) {
	ctx := logger.OnContext(context.Background(), logger.NewTest())
	c := qt.New(t)

	pool, err := pgxmock.NewPool()
	c.Assert(err, qt.IsNil)

	c.Run("No credentials", func(c *qt.C) {
		cfg := config.APIConfig{}
		customResolvers := []resolver.Resolver{}
		collageCache := cache.NewPostgreSQLDependentCache(ctx, cfg, pool, cache.NewPrefixKeyProvider("test"))
		c.Assert(customResolvers, qt.HasLen, 0)
		Initialize(ctx, cfg, pool, &customResolvers, collageCache)
		c.Assert(customResolvers, qt.HasLen, 0)
	})

	c.Run("Credentials", func(c *qt.C) {
		cfg := config.APIConfig{
			TwitterBearerToken: "test",
		}
		customResolvers := []resolver.Resolver{}
		collageCache := cache.NewPostgreSQLDependentCache(ctx, cfg, pool, cache.NewPrefixKeyProvider("test"))
		c.Assert(customResolvers, qt.HasLen, 0)
		Initialize(ctx, cfg, pool, &customResolvers, collageCache)
		c.Assert(customResolvers, qt.HasLen, 1)
	})
}
