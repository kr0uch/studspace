package errors

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
)

type Code string

const (
	CodeBadRequest        Code = "BAD_REQUEST"
	CodeInvalidSign       Code = "INVALID_SIGN"
	CodeForbidden         Code = "FORBIDDEN"
	CodeNotFound          Code = "NOT_FOUND"
	CodeSlotConflict      Code = "SLOT_CONFLICT"
	CodeInvalidTransition Code = "INVALID_TRANSITION"
	CodeEventFull         Code = "EVENT_FULL"
	CodeValidation        Code = "VALIDATION_ERROR"
	CodeTooManyPending    Code = "TOO_MANY_PENDING"
	CodeInternal          Code = "INTERNAL"
)

type HttpError struct {
	Status  int
	Code    Code
	Message string
	Details any
	err     error
}

func NewHttpError(status int, code Code, message string) *HttpError {
	return &HttpError{Status: status, Code: code, Message: message}
}

func (e *HttpError) Error() string {
	msg := string(e.Code) + ": " + e.Message
	if e.err != nil {
		msg += ": " + e.err.Error()
	}
	return msg
}

func (e *HttpError) Unwrap() error {
	return e.err
}

func (e *HttpError) Is(target error) bool {
	t, ok := target.(*HttpError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

func (e *HttpError) WithMessage(message string) *HttpError {
	c := *e
	c.Message = message
	return &c
}

func (e *HttpError) WithDetails(details any) *HttpError {
	c := *e
	c.Details = details
	return &c
}

func (e *HttpError) Wrap(err error) *HttpError {
	c := *e
	c.err = err
	return &c
}

func ParseHttpError(err error) *HttpError {
	var val *HttpError
	if errors.As(err, &val) {
		return val
	}
	return ErrInternal.Wrap(err)
}

func (e *HttpError) SetError(err error) {
	*e = *ParseHttpError(err)
}

func (e *HttpError) HasError() bool {
	return e != nil && e.Code != "" && e.Status != 0
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type ValidationDetails struct {
	Fields map[string]string `json:"fields"`
}

func (e *HttpError) MarshalJSON() ([]byte, error) {
	httpErr := ErrInternal
	if e.HasError() {
		httpErr = e
	}

	body := ErrorBody{Code: httpErr.Code, Message: httpErr.Message}
	if httpErr.Status < http.StatusInternalServerError {
		body.Details = httpErr.Details
	}
	return json.Marshal(ErrorResponse{Error: body})
}

func MarshalError(err error) []byte {
	data, _ := json.Marshal(ParseHttpError(err))
	return data
}

func Validation(fields map[string]string) *HttpError {
	return ErrValidation.WithDetails(ValidationDetails{Fields: fields})
}

func IsTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

var (
	ErrBadRequest = NewHttpError(http.StatusBadRequest, CodeBadRequest, "Некорректный запрос")

	ErrInvalidSign = NewHttpError(http.StatusUnauthorized, CodeInvalidSign, "Неверная или устаревшая подпись параметров запуска")

	ErrForbidden = NewHttpError(http.StatusForbidden, CodeForbidden, "Недостаточно прав")

	ErrNotFound = NewHttpError(http.StatusNotFound, CodeNotFound, "Не найдено")

	ErrSlotConflict      = NewHttpError(http.StatusConflict, CodeSlotConflict, "Выбранное время уже занято")
	ErrInvalidTransition = NewHttpError(http.StatusConflict, CodeInvalidTransition, "Действие недоступно для заявки в текущем статусе")
	ErrEventFull         = NewHttpError(http.StatusConflict, CodeEventFull, "Мест нет")

	ErrValidation = NewHttpError(http.StatusUnprocessableEntity, CodeValidation, "Проверьте правильность заполнения полей")

	ErrTooManyPending = NewHttpError(http.StatusTooManyRequests, CodeTooManyPending, "Слишком много заявок на рассмотрении")

	ErrInternal = NewHttpError(http.StatusInternalServerError, CodeInternal, "Внутренняя ошибка сервера")
)
