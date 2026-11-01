package vk

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	apiURL      = "https://api.vk.com/method/notifications.sendMessage"
	apiVersion  = "5.199"
	maxUserIDs  = 100
	maxMessage  = 254
	httpTimeout = 5 * time.Second

	CodeNotificationsOff = 1
)

var (
	ErrEmptyToken   = stderrors.New("notify token is empty")
	ErrEmptyMessage = stderrors.New("notify message is empty")
)

type VKClient struct {
	token      string
	httpClient *http.Client
}

func NewVKClient(token string) *VKClient {
	return &VKClient{
		token: token,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

func (c *VKClient) Send(ctx context.Context, vkUserIDs []int64, message string) error {
	if c == nil {
		return stderrors.New("vk client is nil")
	}
	if c.token == "" {
		return ErrEmptyToken
	}

	message = clipRunes(message, maxMessage)
	if message == "" {
		return ErrEmptyMessage
	}
	if len(vkUserIDs) == 0 {
		return nil
	}

	var (
		sendErr *SendError
		other   error
	)
	for start := 0; start < len(vkUserIDs); start += maxUserIDs {
		end := min(start+maxUserIDs, len(vkUserIDs))
		if err := c.sendBatch(ctx, vkUserIDs[start:end], message); err != nil {
			var batch *SendError
			if stderrors.As(err, &batch) {
				if sendErr == nil {
					sendErr = &SendError{}
				}
				sendErr.Users = append(sendErr.Users, batch.Users...)
				continue
			}
			if other == nil {
				other = err
			}
		}
	}
	switch {
	case sendErr != nil && other != nil:
		return stderrors.Join(other, sendErr)
	case sendErr != nil:
		return sendErr
	default:
		return other
	}
}

func (c *VKClient) sendBatch(ctx context.Context, vkUserIDs []int64, message string) error {
	form := url.Values{}
	form.Set("user_ids", joinIDs(vkUserIDs))
	form.Set("message", message)
	form.Set("v", apiVersion)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var parsed apiResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("vk notifications.sendMessage: status %d: %w", resp.StatusCode, err)
	}
	if parsed.Error.Code != 0 {
		return parsed.Error
	}

	failed := make([]UserFailure, 0)
	for _, item := range parsed.Response {
		if item.Status || item.Error.Code == 0 {
			continue
		}
		failed = append(failed, UserFailure{
			UserID:      item.UserID,
			Code:        item.Error.Code,
			Description: item.Error.Description,
		})
	}
	if len(failed) > 0 {
		return &SendError{Users: failed}
	}
	return nil
}

type UserFailure struct {
	UserID      int64
	Code        int
	Description string
}

type SendError struct {
	Users []UserFailure
}

func (e *SendError) Error() string {
	if e == nil || len(e.Users) == 0 {
		return "vk notifications.sendMessage"
	}
	parts := make([]string, 0, len(e.Users))
	for _, user := range e.Users {
		parts = append(parts, fmt.Sprintf("vk_user_id=%d code=%d %s", user.UserID, user.Code, user.Description))
	}
	return "vk notifications.sendMessage: " + strings.Join(parts, "; ")
}

func DisabledVKUserIDs(err error) []int64 {
	var sendErr *SendError
	if !stderrors.As(err, &sendErr) || sendErr == nil {
		return nil
	}
	ids := make([]int64, 0)
	for _, user := range sendErr.Users {
		if user.Code == CodeNotificationsOff && user.UserID > 0 {
			ids = append(ids, user.UserID)
		}
	}
	return ids
}

type apiResponse struct {
	Response []sendResult `json:"response"`
	Error    apiError     `json:"error"`
}

type sendResult struct {
	UserID int64     `json:"user_id"`
	Status bool      `json:"status"`
	Error  userError `json:"error"`
}

type userError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type apiError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

func (e apiError) Error() string {
	return fmt.Sprintf("vk api error %d: %s", e.Code, e.Message)
}

func joinIDs(ids []int64) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.FormatInt(id, 10)
	}
	return strings.Join(parts, ",")
}

func clipRunes(s string, n int) string {
	if n < 1 || utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n])
}
