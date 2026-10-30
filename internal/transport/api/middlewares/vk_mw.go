package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/web"
)

const vkSignMaxAge = 24 * time.Hour

type UserLookup interface {
	GetUserByVKUserID(ctx context.Context, vkUserID int64) (*entities.User, error)
}

type VKAuth struct {
	secret []byte
	users  UserLookup
}

func NewVKAuth(secret string, users UserLookup) *VKAuth {
	return &VKAuth{secret: []byte(secret), users: users}
}

func (a *VKAuth) VerifySign(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpError := web.GetHttpErrorFromCtx(r.Context())

		vkUserID, err := parseVKLaunchParams(r.Header.Get("Authorization"), a.secret, time.Now())
		if err != nil {
			httpError.SetError(err)
			return
		}

		next.ServeHTTP(w, r.WithContext(web.ContextWithVKUserID(r.Context(), vkUserID)))
	})
}

func (a *VKAuth) LoadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpError := web.GetHttpErrorFromCtx(r.Context())

		vkUserID, err := web.GetVKUserIDFromCtx(r.Context())
		if err != nil {
			httpError.SetError(err)
			return
		}

		user, err := a.users.GetUserByVKUserID(r.Context(), vkUserID)
		if err != nil {
			httpError.SetError(err)
			return
		}

		next.ServeHTTP(w, r.WithContext(web.ContextWithUser(r.Context(), user)))
	})
}

func parseVKLaunchParams(header string, secret []byte, now time.Time) (int64, error) {
	scheme, rawQuery, ok := strings.Cut(header, " ")
	if !ok || scheme != "VK" || rawQuery == "" {
		return 0, errors.ErrInvalidSign
	}

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return 0, errors.ErrInvalidSign
	}

	sign := values.Get("sign")
	if sign == "" {
		return 0, errors.ErrInvalidSign
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		if strings.HasPrefix(key, "vk_") {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var payload strings.Builder
	for i, key := range keys {
		if i > 0 {
			payload.WriteByte('&')
		}
		payload.WriteString(key)
		payload.WriteByte('=')
		payload.WriteString(values.Get(key))
	}

	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(payload.String()))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(sign)) {
		return 0, errors.ErrInvalidSign
	}

	ts, err := strconv.ParseInt(values.Get("vk_ts"), 10, 64)
	if err != nil {
		return 0, errors.ErrInvalidSign
	}
	if now.Unix()-ts > int64(vkSignMaxAge.Seconds()) {
		return 0, errors.ErrInvalidSign
	}

	vkUserID, err := strconv.ParseInt(values.Get("vk_user_id"), 10, 64)
	if err != nil || vkUserID < 1 {
		return 0, errors.ErrInvalidSign
	}
	return vkUserID, nil
}
