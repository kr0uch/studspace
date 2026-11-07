package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	const origin = "https://app.example"

	newHandler := func(called *bool) http.Handler {
		return CORS([]string{" " + origin + " ", ""})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*called = true
			w.WriteHeader(http.StatusUnauthorized)
		}))
	}

	t.Run("options allowed origin skips next", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodOptions, "/bookings", nil)
		req.Header.Set("Origin", origin)
		req.Header.Set("Access-Control-Request-Method", http.MethodGet)
		req.Header.Set("Access-Control-Request-Headers", "authorization,content-type")
		rr := httptest.NewRecorder()

		newHandler(&called).ServeHTTP(rr, req)

		assert.False(t, called)
		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Equal(t, origin, rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "Origin", rr.Header().Get("Vary"))
		assert.Equal(t, corsAllowMethods, rr.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, corsAllowHeaders, rr.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, corsMaxAge, rr.Header().Get("Access-Control-Max-Age"))
	})

	t.Run("get keeps allow origin on error", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodGet, "/bookings", nil)
		req.Header.Set("Origin", origin)
		rr := httptest.NewRecorder()

		newHandler(&called).ServeHTTP(rr, req)

		assert.True(t, called)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, origin, rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("options foreign origin has no allow origin", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodOptions, "/bookings", nil)
		req.Header.Set("Origin", "https://evil.example")
		rr := httptest.NewRecorder()

		newHandler(&called).ServeHTTP(rr, req)

		assert.False(t, called)
		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("get foreign origin still calls next", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodGet, "/bookings", nil)
		req.Header.Set("Origin", "https://evil.example")
		rr := httptest.NewRecorder()

		newHandler(&called).ServeHTTP(rr, req)

		assert.True(t, called)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("no origin", func(t *testing.T) {
		called := false
		req := httptest.NewRequest(http.MethodGet, "/bookings", nil)
		rr := httptest.NewRecorder()

		CORS(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(rr, req)

		assert.True(t, called)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
	})
}
