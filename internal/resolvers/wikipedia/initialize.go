package wikipedia

import (
	"context"
	"errors"
	"html/template"
	"regexp"

	"github.com/lonelyshoeh/api/internal/db"
	"github.com/lonelyshoeh/api/pkg/config"
	"github.com/lonelyshoeh/api/pkg/resolver"
)

var (
	localeRegexp = regexp.MustCompile(`(?i)([a-z]+)\.wikipedia\.org`)
	titleRegexp  = regexp.MustCompile(`\/wiki\/(.+)`)

	wikipediaTooltipTemplate = template.Must(template.New("wikipediaTooltipTemplate").Parse(wikipediaTooltip))

	errLocaleMatch = errors.New("could not find locale from URL")
	errTitleMatch  = errors.New("could not find title from URL")
)

func Initialize(ctx context.Context, cfg config.APIConfig, pool db.Pool, resolvers *[]resolver.Resolver) {
	const apiURL = "https://%s.wikipedia.org/api/rest_v1/page/summary/%s?redirect=false"

	*resolvers = append(*resolvers, NewArticleResolver(ctx, cfg, pool, apiURL))
}
