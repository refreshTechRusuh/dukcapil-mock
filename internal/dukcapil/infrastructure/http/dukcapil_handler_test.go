package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/refreshTechRusuh/dukcapil-mock/internal/dukcapil/domain/dto"
	"github.com/refreshTechRusuh/dukcapil-mock/pkg/apperror"
)

type stubUC struct {
	res   dto.VerifyResult
	err   error
	got   dto.VerifyRequest
	calls int
}

func (s *stubUC) Verify(_ context.Context, req dto.VerifyRequest) (dto.VerifyResult, error) {
	s.calls++
	s.got = req
	return s.res, s.err
}

func post(h *Handler, body string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	h.Verify(rr, httptest.NewRequest(http.MethodPost, "/api/dukcapil/v1/verify", strings.NewReader(body)))
	return rr
}

const goodBody = `{"nik":"3273000000000001","nama":"Budi","tglLahir":"1990-01-01","namaIbu":"Siti"}`

func TestVerify_ShouldForwardAllFourFieldsToUsecase(t *testing.T) {
	uc := &stubUC{res: dto.VerifyResult{Valid: true}}

	rr := post(NewHandler(uc), goodBody)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, dto.VerifyRequest{NIK: "3273000000000001", Nama: "Budi", TglLahir: "1990-01-01", NamaIbu: "Siti"}, uc.got)
}

func TestVerify_ShouldReturn200_ForBothValidAndInvalid(t *testing.T) {
	for _, valid := range []bool{true, false} {
		rr := post(NewHandler(&stubUC{res: dto.VerifyResult{Valid: valid}}), goodBody)

		assert.Equal(t, http.StatusOK, rr.Code, "valid=%v: hasil pengecekan selalu 200", valid)
		var env struct {
			Status string `json:"status"`
			Data   struct {
				Valid bool `json:"valid"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &env))
		assert.Equal(t, "00000", env.Status)
		assert.Equal(t, valid, env.Data.Valid)
	}
}

func TestVerify_ShouldReturn400_WhenBodyMalformed(t *testing.T) {
	uc := &stubUC{}
	cases := map[string]string{
		"kosong":         ``,
		"bukan json":     `bukan json`,
		"json terpotong": `{"nik":`,
		"tipe salah":     `{"nik":3273000000000001}`,
		"array":          `[]`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			rr := post(NewHandler(uc), body)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})
	}
	assert.Zero(t, uc.calls, "body rusak tidak boleh sampai ke usecase")
}

func TestVerify_ShouldReturn500_WithoutLeakingCause_WhenUsecaseFails(t *testing.T) {
	uc := &stubUC{err: apperror.Internal(errors.New("password authentication failed for user postgres"))}

	rr := post(NewHandler(uc), goodBody)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.NotContains(t, rr.Body.String(), "postgres")
}

func TestVerify_ShouldIgnoreUnknownFields(t *testing.T) {
	uc := &stubUC{res: dto.VerifyResult{Valid: true}}

	rr := post(NewHandler(uc), `{"nik":"1","nama":"B","tglLahir":"t","namaIbu":"S","extra":"x"}`)

	assert.Equal(t, http.StatusOK, rr.Code)
}
