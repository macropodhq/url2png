// Package url2png uses the url2png service to take screenshots of websites.
package url2png

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// Client holds the API key and secret for a given url2png account.
type Client struct {
	Key    string
	Secret string
	Base   *http.Client
}

func (c Client) client() *http.Client {
	b := c.Base
	if b == nil {
		return http.DefaultClient
	}
	return b
}

// Options represents the options that can be provided to the url2png screenshot
// API.
type Options struct {
	ThumbnailMaxWidth uint
	Viewport          [2]int
	Fullpage          bool
	Unique            string
	CustomCSSURL      string
	SayCheese         bool
	TTL               time.Duration
	AcceptLanguages   string
	UserAgent         string
}

// Screenshot submits a request to the url2png service, returning a stream of
// the PNG data, or an error.
func (c Client) Screenshot(website string, options *Options) (io.ReadCloser, error) {
	ctx := context.Background()
	return c.ScreenshotWithContext(ctx, website, options)
}

// ScreenshotWithContext submits a request with context to the url2png service, returning a stream of
// the PNG data, or an error.
func (c Client) ScreenshotWithContext(ctx context.Context, website string, options *Options) (io.ReadCloser, error) {
	q := make(url.Values)

	q.Set("url", website)

	if options != nil {
		if options.ThumbnailMaxWidth != 0 {
			q.Set("thumbnail_max_width", fmt.Sprintf("%d", options.ThumbnailMaxWidth))
		}

		if options.Viewport != [2]int{0, 0} {
			q.Set("viewport", fmt.Sprintf("%dx%d", options.Viewport[0], options.Viewport[1]))
		}

		if options.Fullpage {
			q.Set("fullpage", "true")
		}

		if options.Unique != "" {
			q.Set("unique", options.Unique)
		}

		if options.CustomCSSURL != "" {
			q.Set("custom_css_url", options.CustomCSSURL)
		}

		if options.SayCheese {
			q.Set("say_cheese", "true")
		}

		if options.TTL != time.Duration(0) {
			q.Set("ttl", fmt.Sprintf("%d", int(math.Floor(options.TTL.Seconds()))))
		}

		if options.AcceptLanguages != "" {
			q.Set("accept_languages", options.AcceptLanguages)
		}

		if options.UserAgent != "" {
			q.Set("user_agent", options.UserAgent)
		}
	}

	eq := q.Encode()

	h := crypto.MD5.New()
	_, _ = h.Write([]byte(eq + c.Secret))

	u := url.URL{
		Scheme:   "https",
		Host:     "api.url2png.com",
		Path:     fmt.Sprintf("/v6/%s/%s/png/", c.Key, hex.EncodeToString(h.Sum(nil))),
		RawQuery: eq,
	}

	res, err := ctxhttp.Get(ctx, c.client(), u.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code; expected 200 but got %d", res.StatusCode)
	}

	return res.Body, nil
}
