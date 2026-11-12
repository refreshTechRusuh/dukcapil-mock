package delivery

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// residentDB = Querier tiruan berisi satu penduduk; menjalankan route → handler →
// usecase → repo ASLI, hanya DB-nya yang palsu.
type residentDB struct {
	nik, nama, tgl, ibu string
	err                 error
}

func (d residentDB) QueryRow(_ context.Context, _ string, args ...any) pgx.Row {
	if d.err != nil {
		return row{err: d.err}
	}
	if args[0] != d.nik {
		return row{err: pgx.ErrNoRows}
	}
	return row{vals: []string{d.nik, d.nama, d.tgl, d.ibu}}
}

type row struct {
	vals []string
	err  error
}

func (r row) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	for i, d := range dest {
		*(d.(*string)) = r.vals[i]
	}
	return nil
}

var seeded = residentDB{nik: "3273000000000001", nama: "Budi Santoso", tgl: "1990-01-01", ibu: "Siti Aminah"}

func verify(h http.Handler, body string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/dukcapil/v1/verify", strings.NewReader(body)))
	return rr
}

func TestVerifyRoute_ShouldReturnValidTrue_WhenAllFieldsMatch(t *testing.T) {
	rr := verify(New(seeded), `{"nik":"3273000000000001","nama":"Budi Santoso","tglLahir":"1990-01-01","namaIbu":"Siti Aminah"}`)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"valid":true`)
}

func TestVerifyRoute_ShouldBeCaseAndSpaceInsensitive(t *testing.T) {
	rr := verify(New(seeded), `{"nik":"3273000000000001","nama":"  budi SANTOSO ","tglLahir":"1990-01-01","namaIbu":"siti aminah "}`)

	assert.Contains(t, rr.Body.String(), `"valid":true`)
}

func TestVerifyRoute_ShouldReturnValidFalse_WhenAnyFieldDiffers(t *testing.T) {
	cases := map[string]string{
		"nama beda":     `{"nik":"3273000000000001","nama":"Orang Lain","tglLahir":"1990-01-01","namaIbu":"Siti Aminah"}`,
		"tgl beda":      `{"nik":"3273000000000001","nama":"Budi Santoso","tglLahir":"1991-01-01","namaIbu":"Siti Aminah"}`,
		"nama ibu beda": `{"nik":"3273000000000001","nama":"Budi Santoso","tglLahir":"1990-01-01","namaIbu":"Ibu Lain"}`,
		"nik tak ada":   `{"nik":"3273999999999999","nama":"Budi Santoso","tglLahir":"1990-01-01","namaIbu":"Siti Aminah"}`,
		"field kosong":  `{"nik":"3273000000000001"}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			rr := verify(New(seeded), body)

			assert.Equal(t, http.StatusOK, rr.Code, "tidak cocok bukan error: tetap 200")
			assert.Contains(t, rr.Body.String(), `"valid":false`)
		})
	}
}

func TestVerifyRoute_ShouldReturn400_WhenBodyMalformed(t *testing.T) {
	assert.Equal(t, http.StatusBadRequest, verify(New(seeded), `{`).Code)
}

func TestVerifyRoute_ShouldReturn500_WithoutLeakingDetails_WhenDatabaseFails(t *testing.T) {
	db := residentDB{err: errors.New("password authentication failed for user postgres")}

	rr := verify(New(db), `{"nik":"3273000000000001","nama":"B","tglLahir":"t","namaIbu":"S"}`)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.NotContains(t, rr.Body.String(), "postgres")
}

func TestHealthCheck(t *testing.T) {
	rr := httptest.NewRecorder()

	New(seeded).ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/sys/ping", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "dukcapil-mock ok")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestRouting_ShouldReject_UnknownRoutesAndWrongMethods(t *testing.T) {
	h := New(seeded)
	do := func(method, path string) int {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(method, path, nil))
		return rr.Code
	}

	assert.Equal(t, http.StatusNotFound, do(http.MethodGet, "/tidak/ada"))
	assert.Equal(t, http.StatusMethodNotAllowed, do(http.MethodGet, "/api/dukcapil/v1/verify"))
	assert.Equal(t, http.StatusMethodNotAllowed, do(http.MethodPost, "/sys/ping"))
}

func TestVerifyRoute_ShouldContinueCallersTrace_ButNotTraceHealthCheck(t *testing.T) {
	rec := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	prevTP, prevProp := otel.GetTracerProvider(), otel.GetTextMapPropagator()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() {
		otel.SetTracerProvider(prevTP)
		otel.SetTextMapPropagator(prevProp)
		_ = tp.Shutdown(context.Background())
	})
	h := New(seeded)

	req := httptest.NewRequest(http.MethodPost, "/api/dukcapil/v1/verify",
		strings.NewReader(`{"nik":"3273000000000001","nama":"Budi Santoso","tglLahir":"1990-01-01","namaIbu":"Siti Aminah"}`))
	req.Header.Set("traceparent", "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
	h.ServeHTTP(httptest.NewRecorder(), req)
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/sys/ping", nil))

	spans := rec.Ended()
	require.Len(t, spans, 1, "hanya verify yang di-trace")
	assert.Equal(t, "0af7651916cd43dd8448eb211c80319c", spans[0].SpanContext().TraceID().String(),
		"trace user-service -> dukcapil-mock harus tersambung")
}

// ---------------------------------------------------------------- keamanan

func TestSecurityHeaders_ShouldBePresent_OnVerifyAndHealthAndNotFound(t *testing.T) {
	h := New(seeded)

	for name, rr := range map[string]*httptest.ResponseRecorder{
		"verify 200": verify(h, `{"nik":"3273000000000001","nama":"Budi Santoso","tglLahir":"1990-01-01","namaIbu":"Siti Aminah"}`),
		"verify 400": verify(h, `{`),
		"health": func() *httptest.ResponseRecorder {
			r := httptest.NewRecorder()
			h.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/sys/ping", nil))
			return r
		}(),
		"404": func() *httptest.ResponseRecorder {
			r := httptest.NewRecorder()
			h.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/tidak/ada", nil))
			return r
		}(),
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
			assert.Equal(t, "default-src 'none'", rr.Header().Get("Content-Security-Policy"))
			assert.Equal(t, "no-store", rr.Header().Get("Cache-Control"))
			assert.Empty(t, rr.Header().Get("X-Frame-Options"))
		})
	}
}

func TestVerifyRoute_ShouldReturn413_WhenBodyOverOneMB(t *testing.T) {
	big := `{"nik":"3273000000000001","nama":"` + strings.Repeat("a", (1<<20)+1) + `"}`

	rr := verify(New(seeded), big)

	assert.Equal(t, http.StatusRequestEntityTooLarge, rr.Code)
}
