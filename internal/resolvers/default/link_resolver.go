package defaultresolver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/lonelyshoeh/api/internal/db"
	"github.com/lonelyshoeh/api/internal/logger"
	"github.com/lonelyshoeh/api/internal/resolvers/betterttv"
	"github.com/lonelyshoeh/api/internal/resolvers/discord"
	"github.com/lonelyshoeh/api/internal/resolvers/frankerfacez"
	"github.com/lonelyshoeh/api/internal/resolvers/imgur"
	"github.com/lonelyshoeh/api/internal/resolvers/livestreamfails"
	"github.com/lonelyshoeh/api/internal/resolvers/oembed"
	"github.com/lonelyshoeh/api/internal/resolvers/seventv"
	"github.com/lonelyshoeh/api/internal/resolvers/supinic"
	"github.com/lonelyshoeh/api/internal/resolvers/twitch"
	"github.com/lonelyshoeh/api/internal/resolvers/twitter"
	"github.com/lonelyshoeh/api/internal/resolvers/wikipedia"
	"github.com/lonelyshoeh/api/internal/resolvers/youtube"
	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/config"
	"github.com/lonelyshoeh/api/pkg/resolver"
	"github.com/lonelyshoeh/api/pkg/utils"
	"github.com/nicklaw5/helix"
)

type LinkResolver struct {
	customResolvers []resolver.Resolver

	ignoredHosts map[string]struct{}

	linkCache      cache.Cache
	thumbnailCache cache.Cache
	generatedCache cache.DependentCache
}

func (r *LinkResolver) shouldIgnore(u *url.URL) bool {
	if _, ok := r.ignoredHosts[strings.ToLower(u.Hostname())]; ok {
		// Ignoring url because host is ignored
		return true
	}

	ip := net.ParseIP(u.Host)
	if ip != nil && ip.IsPrivate() {
		return true
	}

	return false
}

func (r *LinkResolver) HandleRequest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	log.Debugw("Handle request",
		"path", req.URL.Path,
	)
	// w.Header().Set("Content-Type", "application/json")
	urlString, err := utils.UnescapeURLArgument(req, "url")
	if err != nil {
		_, err = resolver.WriteInvalidURL(w)
		if err != nil {
			log.Errorw("Error writing response",
				"error", err,
			)
		}
		return
	}

	requestUrl, err := url.Parse(urlString)
	if err != nil {
		log.Errorw("Error parsing url",
			"url", urlString,
			"error", err,
		)
		_, err = resolver.WriteInvalidURL(w)
		if err != nil {
			log.Errorw("Error writing response",
				"error", err,
			)
		}
		return
	}

	if r.shouldIgnore(requestUrl) {
		_, err = resolver.WriteForbiddenURL(w)
		if err != nil {
			log.Errorw("Error writing response",
				"error", err,
			)
		}
		return
	}

	for _, m := range r.customResolvers {
		if ctx, result := m.Check(ctx, requestUrl); result {
			log.Debugw("Run url on custom resolver",
				"name", m.Name(),
				"url", requestUrl,
			)
			data, err := m.Run(ctx, requestUrl, req)

			if errors.Is(err, resolver.ErrDontHandle) {
				break
			}

			resolverHits.WithLabelValues(m.Name()).Inc()

			if err != nil {
				log.Errorw("Error in custom resolver, falling back to default",
					"name", m.Name(),
					"url", requestUrl,
					"error", err,
				)
				break
			}

			w.Header().Add("Content-Type", data.ContentType)
			w.WriteHeader(data.StatusCode)
			_, err = w.Write(data.Payload)
			if err != nil {
				log.Errorw("Error writing response",
					"name", m.Name(),
					"error", err,
				)
			}
			return
		}
	}

	resolverHits.WithLabelValues("default").Inc()

	response, err := r.linkCache.Get(ctx, urlString, req)
	if err != nil {
		log.Errorw("Error in default resolver",
			"url", requestUrl,
			"error", err,
		)
		_, err = resolver.WriteInternalServerErrorf(w, "Error resolving link")
		if err != nil {
			log.Errorw("Error in default resolver",
				"url", requestUrl,
				"error", err,
			)
		}
	} else {
		w.Header().Add("Content-Type", response.ContentType)
		w.WriteHeader(response.StatusCode)
		_, err = w.Write(response.Payload)
		if err != nil {
			log.Errorw("Error writing response",
				"error", err,
			)
		}
	}
}

func (r *LinkResolver) HandleThumbnailRequest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	url, err := utils.UnescapeURLArgument(req, "url")
	if err != nil {
		_, err = resolver.WriteInvalidURL(w)
		if err != nil {
			log.Errorw("Error writing response",
				"error", err,
			)
		}
		return
	}

	response, err := r.thumbnailCache.Get(ctx, url, req)
	if err != nil {
		log.Errorw("Error in thumbnail request",
			"url", url,
			"error", err,
		)
		return
	}

	w.Header().Add("Content-Type", response.ContentType)
	w.WriteHeader(response.StatusCode)
	_, err = w.Write(response.Payload)
	if err != nil {
		log.Errorw("Error writing response",
			"error", err,
		)
	}
}

func (r *LinkResolver) HandleGeneratedValueRequest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	url, err := utils.UnescapeURLArgument(req, "url")
	if err != nil {
		_, err = resolver.WriteInvalidURL(w)
		if err != nil {
			log.Errorw("Error writing response",
				"error", err,
			)
		}
		return
	}

	payload, contentType, err := r.generatedCache.Get(ctx, url)
	if err != nil {
		log.Errorw("Error in request for generated value",
			"url", url,
			"error", err,
		)
		return
	}

	if payload == nil {
		log.Warnw("Requested generated value does not exist",
			"url", url,
		)
		return
	}

	w.Header().Add("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(payload)
	if err != nil {
		log.Errorw("Error writing response",
			"error", err,
		)
	}
}

func New(ctx context.Context, cfg config.APIConfig, pool db.Pool, helixClient *helix.Client, ignoredHosts map[string]struct{}) *LinkResolver {
	generatedCache := cache.NewPostgreSQLDependentCache(ctx, cfg, pool, cache.NewPrefixKeyProvider("default:dependent"))

	customResolvers := []resolver.Resolver{}

	// Register Link Resolvers from internal/resolvers/
	betterttv.Initialize(ctx, cfg, pool, &customResolvers)
	discord.Initialize(ctx, cfg, pool, &customResolvers)
	frankerfacez.Initialize(ctx, cfg, pool, &customResolvers)
	imgur.Initialize(ctx, cfg, pool, &customResolvers)
	livestreamfails.Initialize(ctx, cfg, pool, &customResolvers)
	oembed.Initialize(ctx, cfg, pool, &customResolvers)
	supinic.Initialize(ctx, cfg, pool, &customResolvers)
	twitch.Initialize(ctx, cfg, pool, helixClient, &customResolvers)
	twitter.Initialize(ctx, cfg, pool, &customResolvers, generatedCache)
	wikipedia.Initialize(ctx, cfg, pool, &customResolvers)
	youtube.Initialize(ctx, cfg, pool, &customResolvers)
	seventv.Initialize(ctx, cfg, pool, &customResolvers)

	// The content type resolvers should match from most to least specific
	contentTypeResolvers := []ContentTypeResolver{
		NewPDFResolver(cfg.BaseURL, cfg.MaxContentLength),
		NewMediaResolver(cfg.BaseURL),
	}

	linkLoader := &LinkLoader{
		baseURL:              cfg.BaseURL,
		maxContentLength:     cfg.MaxContentLength,
		customResolvers:      customResolvers,
		contentTypeResolvers: contentTypeResolvers,
	}
	thumbnailLoader := &ThumbnailLoader{
		baseURL:                  cfg.BaseURL,
		maxContentLength:         cfg.MaxContentLength,
		enableAnimatedThumbnails: cfg.EnableAnimatedThumbnails,
	}

	thumbnailCache := cache.NewPostgreSQLCache(
		ctx, cfg, pool, cache.NewPrefixKeyProvider("default:thumbnail"), thumbnailLoader,
		cfg.ThumbnailCacheDuration,
	)
	linkCache := cache.NewPostgreSQLCache(
		ctx, cfg, pool, cache.NewPrefixKeyProvider("default:link"), linkLoader, cfg.DefaultLinkCacheDuration,
	)

	r := &LinkResolver{
		customResolvers: customResolvers,

		ignoredHosts: ignoredHosts,

		linkCache:      linkCache,
		thumbnailCache: thumbnailCache,
		generatedCache: generatedCache,
	}

	return r
}
