package youtube

import (
	"context"
	"net/http"
	"net/url"

	"github.com/lonelyshoeh/api/pkg/cache"
)

type YouTubeVideoShortURLResolver struct {
	videoCache cache.Cache
}

func (r *YouTubeVideoShortURLResolver) Check(ctx context.Context, url *url.URL) (context.Context, bool) {
	if url.Host != "youtu.be" {
		return ctx, false
	}

	videoID := getYoutubeVideoIDFromURL2(url)

	return ctx, videoID != "" && videoID != "."
}

func (r *YouTubeVideoShortURLResolver) Run(ctx context.Context, url *url.URL, req *http.Request) (*cache.Response, error) {
	videoID := getYoutubeVideoIDFromURL2(url)

	if videoID == "" || videoID == "." {
		return nil, errInvalidVideoLink
	}

	return r.videoCache.Get(ctx, videoID, req)
}

func (r *YouTubeVideoShortURLResolver) Name() string {
	return "youtube:video:shorturl"
}

func NewYouTubeVideoShortURLResolver(videoCache cache.Cache) *YouTubeVideoShortURLResolver {
	r := &YouTubeVideoShortURLResolver{
		videoCache: videoCache,
	}

	return r
}
