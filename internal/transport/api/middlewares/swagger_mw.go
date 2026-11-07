package middlewares

import (
	"net/http"

	_ "github.com/kr0uch/studspace/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Swagger(prefix string) http.Handler {
	swaggerPrefix := prefix + "/swagger/"
	ui := httpSwagger.Handler(httpSwagger.URL(swaggerPrefix + "doc.json"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == swaggerPrefix {
			http.Redirect(w, r, swaggerPrefix+"index.html", http.StatusTemporaryRedirect)
			return
		}
		ui.ServeHTTP(w, r)
	})
}
