package apperror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructors_ShouldSetCodeAndMessage(t *testing.T) {
	assert.Equal(t, ErrValidation, Validation("salah").Code)
	assert.Equal(t, "salah", Validation("salah").Message)
	assert.Equal(t, ErrInternal, Internal(errors.New("x")).Code)
}

func TestError_ShouldIncludeCause_WhenWrapped(t *testing.T) {
	cause := errors.New("db down")

	wrapped := Wrap(ErrInternal, "gagal", cause)

	assert.Contains(t, wrapped.Error(), "db down")
	assert.ErrorIs(t, wrapped, cause)
}

func TestError_ShouldOmitCause_WhenNotWrapped(t *testing.T) {
	assert.Equal(t, "[VALIDATION_ERROR] salah", New(ErrValidation, "salah").Error())
}

func TestInternal_ShouldNotLeakCauseToClient(t *testing.T) {
	assert.NotContains(t, Internal(errors.New("password=rahasia")).Message, "rahasia")
}
