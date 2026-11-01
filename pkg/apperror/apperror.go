// Package apperror mendefinisikan tipe error domain dukcapil-mock.
// Pola diambil dari service-dashboard / event-service pkg/apperror:
// usecase menaikkan AppError{Code, Message}; lapis HTTP (apinizer) memetakan
// Code -> HTTP status lewat errors.As. Dengan begitu logika tak kenal HTTP.
package apperror

import "fmt"

// ErrorCode adalah kode error domain yang dipetakan ke HTTP status di apinizer.
type ErrorCode string

const (
	// 400
	ErrValidation ErrorCode = "VALIDATION_ERROR"

	// 500
	ErrInternal ErrorCode = "INTERNAL_ERROR"
)

// AppError adalah error domain dengan kode spesifik + pesan human-readable.
// Bungkus dengan errors.As di apinizer agar bisa di-unwrap ke AppError.
type AppError struct {
	Code    ErrorCode
	Message string
	cause   error
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code ErrorCode, message string, cause error) *AppError {
	return &AppError{Code: code, Message: message, cause: cause}
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.cause }

// ---- Constructor helpers -----------------------------------------------

func Validation(msg string) *AppError { return New(ErrValidation, msg) }
func Internal(cause error) *AppError  { return Wrap(ErrInternal, "Terjadi kesalahan pada server", cause) }
