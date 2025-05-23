package seventv

import (
	"net/http"

	"github.com/lonelyshoeh/api/pkg/resolver"
)

var (
	emoteNotFoundResponse = &resolver.Response{
		Status:  http.StatusNotFound,
		Message: "No 7TV emote with this id found",
	}
)
