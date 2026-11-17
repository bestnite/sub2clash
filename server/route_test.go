package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/bestnite/sub2clash/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Exercise the embedded production build through a prefix-stripping proxy.
func TestSubpathAssetsAndRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousLogger := logger.Logger
	logger.Logger = zap.NewNop()
	t.Cleanup(func() { logger.Logger = previousLogger })
	router := gin.New()
	SetRoute(router)
	upstream := httptest.NewServer(router)
	defer upstream.Close()
	target, err := url.Parse(upstream.URL)
	if err != nil {
		t.Fatal(err)
	}

	for _, prefix := range []string{"/", "/sub2clash/", "/tools/sub2clash/"} {
		t.Run(prefix, func(t *testing.T) {
			proxy := httputil.NewSingleHostReverseProxy(target)
			edge := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if prefix != "/" && r.URL.Path == strings.TrimSuffix(prefix, "/") {
					http.Redirect(w, r, prefix, http.StatusPermanentRedirect)
					return
				}
				if !strings.HasPrefix(r.URL.Path, prefix) {
					http.NotFound(w, r)
					return
				}
				r.URL.Path = "/" + strings.TrimPrefix(r.URL.Path, prefix)
				proxy.ServeHTTP(w, r)
			}))
			defer edge.Close()

			pageURL := edge.URL + prefix
			resp, err := http.Get(pageURL)
			if err != nil {
				t.Fatal(err)
			}
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("homepage: status=%d err=%v", resp.StatusCode, err)
			}
			assets := regexp.MustCompile(`(?:src|href)="([^"]*assets/[^"]+)"`).FindAllStringSubmatch(string(body), -1)
			if len(assets) < 2 {
				t.Fatalf("expected JS and CSS references in production HTML: %s", body)
			}
			for _, asset := range assets {
				reference, err := url.Parse(asset[1])
				if err != nil {
					t.Fatal(err)
				}
				assetURL := resp.Request.URL.ResolveReference(reference)
				if !strings.HasPrefix(assetURL.Path, prefix+"assets/") {
					t.Fatalf("asset escaped deployment prefix: %s", assetURL)
				}
				result, err := http.Get(assetURL.String())
				if err != nil {
					t.Fatal(err)
				}
				result.Body.Close()
				if result.StatusCode != http.StatusOK {
					t.Fatalf("asset %s: status=%d", assetURL, result.StatusCode)
				}
			}

			if prefix != "/" {
				result, err := http.Get(strings.TrimSuffix(pageURL, "/"))
				if err != nil {
					t.Fatal(err)
				}
				result.Body.Close()
				if result.Request.URL.Path != prefix || result.StatusCode != http.StatusOK {
					t.Fatalf("missing trailing-slash redirect: %s (%d)", result.Request.URL, result.StatusCode)
				}
			}

			// Invalid input must reach the application handler, not a proxy 404.
			for _, route := range []struct {
				method, path string
				status       int
			}{
				{http.MethodGet, "convert/invalid", http.StatusBadRequest},
				{http.MethodPost, "short", http.StatusBadRequest},
				{http.MethodPut, "short", http.StatusBadRequest},
				{http.MethodGet, "short/issue-76-missing", http.StatusUnauthorized},
				{http.MethodGet, "s/issue-76-missing", http.StatusUnauthorized},
				{http.MethodDelete, "short/issue-76-missing", http.StatusBadRequest},
			} {
				req, err := http.NewRequest(route.method, pageURL+route.path, nil)
				if err != nil {
					t.Fatal(err)
				}
				result, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatal(err)
				}
				result.Body.Close()
				if result.StatusCode != route.status {
					t.Fatalf("%s %s: status=%d", route.method, route.path, result.StatusCode)
				}
			}
		})
	}
}
