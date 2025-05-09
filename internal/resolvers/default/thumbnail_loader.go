package defaultresolver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lonelyshoeh/api/internal/logger"
	"github.com/lonelyshoeh/api/internal/staticresponse"
	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/resolver"
	"github.com/lonelyshoeh/api/pkg/thumbnail"
)

type ThumbnailLoader struct {
	baseURL                  string
	maxContentLength         uint64
	enableAnimatedThumbnails bool
}

func (l *ThumbnailLoader) Load(ctx context.Context, urlString string, r *http.Request) ([]byte, *int, *string, time.Duration, error) {
	log := logger.FromContext(ctx)

	url, err := url.Parse(urlString)
	if err != nil {
		return resolver.ReturnInvalidURL()
	}

	resp, err := resolver.RequestGET(ctx, url.String())
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such host") {
			return resolver.InternalServerErrorf("Error loading thumbnail, could not resolve host %s", err.Error())
		}

		return resolver.InternalServerErrorf("Error loading thumbnail: %s", err.Error())
	}

	defer resp.Body.Close()

	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		contentLengthBytes, err := strconv.Atoi(contentLength)
		if err != nil {
			r := &resolver.Response{
				Status:  http.StatusInternalServerError,
				Message: resolver.CleanResponse(fmt.Sprintf("Invalid content length: %s - %s", contentLength, err.Error())),
			}
			marshalledPayload, err := json.Marshal(r)
			if err != nil {
				panic(err)
			}

			return marshalledPayload, nil, nil, resolver.NoSpecialDur, nil
		}

		if uint64(contentLengthBytes) > l.maxContentLength {
			return resolver.FResponseTooLarge()
		}
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusMultipleChoices {
		log.Infow("Skipping url because of status code", "url", resp.Request.URL, "status", resp.StatusCode)
		return staticresponse.SNoThumbnailFound.Return()
	}

	contentType := resp.Header.Get("Content-Type")

	if !thumbnail.IsSupportedThumbnailType(contentType) {
		return resolver.UnsupportedThumbnailType, nil, nil, cache.NoSpecialDur, nil
	}

	inputBuf, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorw("Error reading body from request", "error", err)
		return resolver.ErrorBuildingThumbnail, nil, nil, cache.NoSpecialDur, nil
	}

	var image []byte
	tryAnimatedThumb := l.enableAnimatedThumbnails && thumbnail.IsAnimatedThumbnailType(contentType)

	// attempt building an animated image
	if tryAnimatedThumb {
		image, err = thumbnail.BuildAnimatedThumbnail(ctx, inputBuf, resp)
	}

	// fallback to static image if animated image building failed or is disabled
	if !tryAnimatedThumb || err != nil {
		if err != nil {
			log.Errorw("Error trying to build animated thumbnail, falling back to static thumbnail building",
				"error", err)
		}
		image, err = thumbnail.BuildStaticThumbnail(inputBuf, resp)
		if err != nil {
			log.Errorw("Error trying to build static thumbnail", "error", err)
			return resolver.InternalServerErrorf("Error building static thumbnail: %s", err.Error())
		}
	}

	return image, nil, &contentType, 10 * time.Minute, nil
}
