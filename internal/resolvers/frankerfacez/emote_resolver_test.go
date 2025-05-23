package frankerfacez

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/lonelyshoeh/api/internal/logger"
	"github.com/lonelyshoeh/api/pkg/cache"
	"github.com/lonelyshoeh/api/pkg/config"
	"github.com/lonelyshoeh/api/pkg/utils"
	qt "github.com/frankban/quicktest"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
)

func TestEmoteResolver(t *testing.T) {
	ctx := logger.OnContext(context.Background(), logger.NewTest())
	c := qt.New(t)

	pool, _ := pgxmock.NewPool()

	cfg := config.APIConfig{}
	ts := testServer()
	defer ts.Close()
	emoteAPIURL := utils.MustParseURL(ts.URL + "/v1/emote/")
	resolver := NewEmoteResolver(ctx, cfg, pool, emoteAPIURL)

	c.Assert(resolver, qt.IsNotNil)

	c.Run("Name", func(c *qt.C) {
		c.Assert(resolver.Name(), qt.Equals, "frankerfacez:emote")
	})

	c.Run("Check", func(c *qt.C) {
		type checkTest struct {
			label    string
			input    *url.URL
			expected bool
		}

		tests := []checkTest{
			{
				label:    "Matching domain, no WWW",
				input:    utils.MustParseURL("https://frankerfacez.com/emoticon/297734-pajaSx"),
				expected: true,
			},
			{
				label:    "Matching domain, WWW",
				input:    utils.MustParseURL("https://www.frankerfacez.com/emoticon/297734-pajaSx"),
				expected: true,
			},
			{
				label:    "Matching domain, non-matching path",
				input:    utils.MustParseURL("https://frankerfacez.com/user/566ca04265dbbdab32ec054a"),
				expected: false,
			},
			{
				label:    "Non-matching domain",
				input:    utils.MustParseURL("https://example.com/emoticon/297734-pajaSx"),
				expected: false,
			},
		}

		for _, test := range tests {
			c.Run(test.label, func(c *qt.C) {
				_, output := resolver.Check(ctx, test.input)
				c.Assert(output, qt.Equals, test.expected)
			})
		}
	})

	c.Run("Run", func(c *qt.C) {
		c.Run("Error", func(c *qt.C) {
			type runTest struct {
				label          string
				inputURL       *url.URL
				inputEmoteHash string
				inputReq       *http.Request
				expectedError  error
			}

			tests := []runTest{
				{
					label:         "Non-matching link",
					inputURL:      utils.MustParseURL("https://frankerfacez.com/user/566ca04265dbbdab32ec054a"),
					expectedError: errInvalidFrankerFaceZEmotePath,
				},
			}

			const q = `SELECT value FROM cache WHERE key=$1`

			for _, test := range tests {
				c.Run(test.label, func(c *qt.C) {
					outputBytes, outputError := resolver.Run(ctx, test.inputURL, test.inputReq)
					c.Assert(outputError, qt.Equals, test.expectedError)
					c.Assert(outputBytes, qt.IsNil)
				})
			}
		})
		c.Run("Cached", func(c *qt.C) {
			type runTest struct {
				label          string
				inputURL       *url.URL
				inputEmoteHash string
				inputReq       *http.Request
				// expectedResponse will be returned from the cache, and expected to be returned in the same form
				expectedResponse *cache.Response
				expectedError    error
			}

			tests := []runTest{
				{
					label:          "Matching link - cached",
					inputURL:       utils.MustParseURL("https://frankerfacez.com/emoticon/297734-pajaSx"),
					inputEmoteHash: "297734",
					inputReq:       nil,
					expectedResponse: &cache.Response{
						Payload:     []byte(`{"status":200,"thumbnail":"https://cdn.frankerfacez.com/emoticon/297734/4","tooltip":"%3Cdiv%20style=%22text-align:%20left%3B%22%3E%0A%3Cb%3EpajaSx%3C%2Fb%3E%3Cbr%3E%0A%3Cb%3EFrankerFaceZ%20Emote%3C%2Fb%3E%3Cbr%3E%0A%3Cb%3EBy:%3C%2Fb%3E%20pajlada%3C%2Fdiv%3E"}`),
						StatusCode:  http.StatusOK,
						ContentType: "application/json",
					},
					expectedError: nil,
				},
				{
					label:          "Matching link - cached 2",
					inputURL:       utils.MustParseURL("https://frankerfacez.com/emoticon/367887-paaaajaW"),
					inputEmoteHash: "367887",
					inputReq:       nil,
					expectedResponse: &cache.Response{
						Payload:     []byte(`{"status":200,"thumbnail":"https://cdn.frankerfacez.com/emoticon/297734/4","tooltip":"%3Cdiv%20style=%22text-align:%20left%3B%22%3E%0A%3Cb%3EpaaaajaW%3C%2Fb%3E%3Cbr%3E%0A%3Cb%3EFrankerFaceZ%20Emote%3C%2Fb%3E%3Cbr%3E%0A%3Cb%3EBy:%3C%2Fb%3E%20Goran42069%3C%2Fdiv%3E"}`),
						StatusCode:  http.StatusOK,
						ContentType: "application/json",
					},
					expectedError: nil,
				},
				{
					label:          "Matching link - 404",
					inputURL:       utils.MustParseURL("https://frankerfacez.com/emoticon/404-404"),
					inputEmoteHash: "404",
					inputReq:       nil,
					expectedResponse: &cache.Response{
						Payload:     []byte(`{"status":404,"message":"No FrankerFaceZ emote with this id found"}`),
						StatusCode:  http.StatusOK,
						ContentType: "application/json",
					},
					expectedError: nil,
				},
			}

			const q = `SELECT value FROM cache WHERE key=$1`

			for _, test := range tests {
				c.Run(test.label, func(c *qt.C) {
					rows := pgxmock.NewRows([]string{"value", "http_status_code", "http_content_type"}).AddRow(test.expectedResponse.Payload, test.expectedResponse.StatusCode, test.expectedResponse.ContentType)
					pool.ExpectQuery("SELECT").
						WithArgs("frankerfacez:emote:" + test.inputEmoteHash).
						WillReturnRows(rows)
					outputBytes, outputError := resolver.Run(ctx, test.inputURL, test.inputReq)
					c.Assert(outputError, qt.Equals, test.expectedError)
					c.Assert(outputBytes, qt.DeepEquals, test.expectedResponse)
				})
			}
		})

		c.Run("Not cached", func(c *qt.C) {
			type runTest struct {
				label            string
				inputURL         *url.URL
				inputEmoteHash   string
				inputReq         *http.Request
				expectedResponse *cache.Response
				expectedError    error
				rowsReturned     int
			}

			tests := []runTest{
				{
					label:          "Matching link - not cached",
					inputURL:       utils.MustParseURL("https://frankerfacez.com/emoticon/297734-pajaSx"),
					inputEmoteHash: "297734",
					inputReq:       nil,
					expectedResponse: &cache.Response{
						Payload:     []byte(`{"status":200,"thumbnail":"https://cdn.frankerfacez.com/emoticon/297734/4","tooltip":"%3Cdiv%20style=%22text-align:%20left%3B%22%3E%0A%3Cb%3EpajaSx%3C%2Fb%3E%3Cbr%3E%0A%3Cb%3EFrankerFaceZ%20Emote%3C%2Fb%3E%3Cbr%3E%0A%3Cb%3EBy:%3C%2Fb%3E%20pajlada%3C%2Fdiv%3E"}`),
						StatusCode:  http.StatusOK,
						ContentType: "application/json",
					},
					expectedError: nil,
				},
				{
					label:          "Matching link - 404",
					inputURL:       utils.MustParseURL("https://frankerfacez.com/emoticon/404-404"),
					inputEmoteHash: "404",
					inputReq:       nil,
					expectedResponse: &cache.Response{
						Payload:     []byte(`{"status":404,"message":"No FrankerFaceZ emote with this id found"}`),
						StatusCode:  http.StatusOK,
						ContentType: "application/json",
					},
					expectedError: nil,
				},
				{
					label:          "Matching link - bad json",
					inputURL:       utils.MustParseURL("https://frankerfacez.com/emoticon/696969-badjson"),
					inputEmoteHash: "696969",
					inputReq:       nil,
					expectedResponse: &cache.Response{
						Payload:     []byte(`{"status":500,"message":"FrankerFaceZ API response decode error: invalid character \u0026#39;x\u0026#39; looking for beginning of value"}`),
						StatusCode:  http.StatusOK,
						ContentType: "application/json",
					},
					expectedError: nil,
				},
			}

			const q = `SELECT value FROM cache WHERE key=$1`

			for _, test := range tests {
				c.Run(test.label, func(c *qt.C) {
					pool.ExpectQuery("SELECT").WillReturnError(pgx.ErrNoRows)
					pool.ExpectExec("INSERT INTO cache").
						WithArgs("frankerfacez:emote:"+test.inputEmoteHash, test.expectedResponse.Payload, test.expectedResponse.StatusCode, test.expectedResponse.ContentType, pgxmock.AnyArg()).
						WillReturnResult(pgxmock.NewResult("INSERT", 1))
					outputBytes, outputError := resolver.Run(ctx, test.inputURL, test.inputReq)
					c.Assert(outputError, qt.Equals, test.expectedError)
					c.Assert(outputBytes, qt.DeepEquals, test.expectedResponse)
				})
			}
		})
	})
}
