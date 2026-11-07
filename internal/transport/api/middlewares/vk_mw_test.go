package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"testing"
	"time"

	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const vkTestSecret = "wvl68m4dR1UpLrVRli"

func signVK(payload string) string {
	mac := hmac.New(sha256.New, []byte(vkTestSecret))
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func TestParseVKLaunchParams(t *testing.T) {
	now := time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	ts := strconv.FormatInt(now.Unix(), 10)

	t.Run("success (single scope)", func(t *testing.T) {
		payload := "vk_access_token_settings=notify&vk_ts=" + ts + "&vk_user_id=42"
		header := "VK " + payload + "&sign=" + signVK(payload)

		id, err := parseVKLaunchParams(header, []byte(vkTestSecret), now)
		require.NoError(t, err)
		assert.Equal(t, int64(42), id)
	})

	t.Run("success (comma is percent-encoded)", func(t *testing.T) {
		payload := "vk_access_token_settings=friends%2Cnotify&vk_ts=" + ts + "&vk_user_id=42"
		header := "VK vk_user_id=42&vk_access_token_settings=friends%2Cnotify&vk_ts=" + ts + "&sign=" + signVK(payload)

		id, err := parseVKLaunchParams(header, []byte(vkTestSecret), now)
		require.NoError(t, err)
		assert.Equal(t, int64(42), id)
	})

	t.Run("success (space is percent-encoded)", func(t *testing.T) {
		payload := "vk_ref=other%20page&vk_ts=" + ts + "&vk_user_id=7"
		header := "VK vk_ref=other%20page&vk_ts=" + ts + "&vk_user_id=7&sign=" + signVK(payload)

		id, err := parseVKLaunchParams(header, []byte(vkTestSecret), now)
		require.NoError(t, err)
		assert.Equal(t, int64(7), id)
	})

	t.Run("ignores params outside vk_", func(t *testing.T) {
		payload := "vk_ts=" + ts + "&vk_user_id=42"
		header := "VK vk_user_id=42&vk_ts=" + ts + "&foo=1&sign=" + signVK(payload)

		id, err := parseVKLaunchParams(header, []byte(vkTestSecret), now)
		require.NoError(t, err)
		assert.Equal(t, int64(42), id)
	})

	t.Run("failure (sign of decoded comma)", func(t *testing.T) {
		payload := "vk_access_token_settings=friends,notify&vk_ts=" + ts + "&vk_user_id=42"
		header := "VK vk_access_token_settings=friends%2Cnotify&vk_ts=" + ts + "&vk_user_id=42&sign=" + signVK(payload)

		id, err := parseVKLaunchParams(header, []byte(vkTestSecret), now)
		assert.ErrorIs(t, err, errors.ErrInvalidSign)
		assert.Zero(t, id)
	})

	t.Run("success (exactly 24h)", func(t *testing.T) {
		old := strconv.FormatInt(now.Add(-vkSignMaxAge).Unix(), 10)
		payload := "vk_ts=" + old + "&vk_user_id=42"
		header := "VK " + payload + "&sign=" + signVK(payload)

		id, err := parseVKLaunchParams(header, []byte(vkTestSecret), now)
		require.NoError(t, err)
		assert.Equal(t, int64(42), id)
	})

	t.Run("failure (user id)", func(t *testing.T) {
		payload := "vk_ts=" + ts + "&vk_user_id=0"
		header := "VK " + payload + "&sign=" + signVK(payload)

		_, err := parseVKLaunchParams(header, []byte(vkTestSecret), now)
		assert.ErrorIs(t, err, errors.ErrInvalidSign)
	})

	t.Run("failure (expired)", func(t *testing.T) {
		old := strconv.FormatInt(now.Add(-vkSignMaxAge-time.Second).Unix(), 10)
		payload := "vk_ts=" + old + "&vk_user_id=42"
		header := "VK " + payload + "&sign=" + signVK(payload)

		_, err := parseVKLaunchParams(header, []byte(vkTestSecret), now)
		assert.ErrorIs(t, err, errors.ErrInvalidSign)
	})

	t.Run("failure (bad header)", func(t *testing.T) {
		_, err := parseVKLaunchParams("Bearer 1", []byte(vkTestSecret), now)
		assert.ErrorIs(t, err, errors.ErrInvalidSign)
	})
}
