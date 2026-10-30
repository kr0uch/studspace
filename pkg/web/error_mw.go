package web

import (
	"context"
	"net/http"

	"github.com/kr0uch/studspace/pkg/errors"
)

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpError := &errors.HttpError{}
		ctx := context.WithValue(r.Context(), CtxHttpErrorKey, httpError)

		next.ServeHTTP(w, r.WithContext(ctx))

		if httpError.HasError() {
			w.WriteHeader(httpError.Status)
			payload := errors.MarshalError(httpError)

			if _, err := w.Write(payload); err != nil {
				http.Error(w, errors.ErrInternal.Error(), http.StatusInternalServerError)
			}
		}
	})
}
