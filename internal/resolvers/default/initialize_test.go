package defaultresolver

import (
	"context"
	"testing"

	"github.com/lonelyshoeh/api/internal/logger"
	"github.com/lonelyshoeh/api/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/pashagolub/pgxmock"

	qt "github.com/frankban/quicktest"
)

func TestInitialize(t *testing.T) {
	ctx := logger.OnContext(context.Background(), logger.NewTest())
	c := qt.New(t)

	pool, err := pgxmock.NewPool()
	c.Assert(err, qt.IsNil)

	r := chi.NewRouter()

	c.Run("No credentials", func(c *qt.C) {
		cfg := config.APIConfig{}
		Initialize(ctx, cfg, pool, r, nil)
	})
}
