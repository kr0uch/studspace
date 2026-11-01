package web

import (
	"net/http"

	"github.com/kr0uch/studspace/pkg/errors"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func HandleWithError(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpError := GetHttpErrorFromCtx(r.Context())

		err := handler(w, r)
		if err == nil {
			return
		}
		if errors.IsNotify(err) {
			GetNotifyErrorFromCtx(r.Context()).SetError(err)
			return
		}
		httpError.SetError(err)
	}
}
