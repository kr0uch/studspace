package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwagger(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/swagger/", Swagger("/api/v1"))

	t.Run("directory redirects to index", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/swagger/", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusTemporaryRedirect, rec.Code)
		assert.Equal(t, "/api/v1/swagger/index.html", rec.Header().Get("Location"))
	})

	t.Run("query still redirects to index", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/swagger/?spec=1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusTemporaryRedirect, rec.Code)
		assert.Equal(t, "/api/v1/swagger/index.html", rec.Header().Get("Location"))
	})

	t.Run("doc stays on prefixed path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/swagger/doc.json", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"basePath": "/api/v1"`)
	})

	t.Run("index stays on prefixed path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/swagger/index.html", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `url: "\/api\/v1\/swagger\/doc.json"`)
	})
}
