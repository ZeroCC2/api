package supinic

import (
	"net/http"

	"github.com/lonelyshoeh/api/pkg/resolver"
)

var (
	trackNotFoundResponse = &resolver.Response{
		Status:  http.StatusNotFound,
		Message: "No track with this ID found",
	}
)
