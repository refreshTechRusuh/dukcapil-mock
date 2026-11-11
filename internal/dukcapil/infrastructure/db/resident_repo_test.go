package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fake Querier (tanpa Postgres). pgxmock v3.4.0 belum cocok dengan pgx v5.11.

type fakeQuerier struct {
	vals    []any
	err     error
	gotSQL  string
	gotArgs []any
}

func (f *fakeQuerier) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
	f.gotSQL, f.gotArgs = sql, args
	return fakeRow{vals: f.vals, err: f.err}
}

type fakeRow struct {
	vals []any
	err  error
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	for i, d := range dest {
		dv, sv := reflect.ValueOf(d).Elem(), reflect.ValueOf(r.vals[i])
		if !sv.Type().AssignableTo(dv.Type()) {
			return fmt.Errorf("scan kolom %d: %s tidak bisa ke %s", i, sv.Type(), dv.Type())
		}
		dv.Set(sv)
	}
	return nil
}

func TestFindByNIK_ShouldReturnResident_WhenFound(t *testing.T) {
	q := &fakeQuerier{vals: []any{"3273000000000001", "Budi Santoso", "1990-01-01", "Siti Aminah"}}

	got, found, err := NewResidentRepository(q).FindByNIK(context.Background(), "3273000000000001")

	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "Budi Santoso", got.Nama)
	assert.Equal(t, "1990-01-01", got.TglLahir)
	assert.Equal(t, "Siti Aminah", got.NamaIbu)
}

func TestFindByNIK_ShouldPassNIKAsParameter_NotInSQL(t *testing.T) {
	q := &fakeQuerier{err: pgx.ErrNoRows}
	nik := "3273000000000001' OR '1'='1"

	_, _, _ = NewResidentRepository(q).FindByNIK(context.Background(), nik)

	assert.Equal(t, []any{nik}, q.gotArgs, "NIK harus lewat parameter ($1)")
	assert.NotContains(t, q.gotSQL, "OR '1'='1", "tidak boleh ada penyambungan string ke SQL")
	assert.Contains(t, q.gotSQL, "$1")
}

func TestFindByNIK_ShouldReportNotFound_WithoutError_WhenNoRows(t *testing.T) {
	q := &fakeQuerier{err: pgx.ErrNoRows}

	got, found, err := NewResidentRepository(q).FindByNIK(context.Background(), "9999999999999999")

	assert.NoError(t, err, "NIK tak terdaftar bukan error teknis")
	assert.False(t, found)
	assert.Empty(t, got.NIK)
}

func TestFindByNIK_ShouldWrapTechnicalError(t *testing.T) {
	boom := errors.New("connection refused")
	q := &fakeQuerier{err: boom}

	_, found, err := NewResidentRepository(q).FindByNIK(context.Background(), "3273000000000001")

	assert.ErrorIs(t, err, boom)
	assert.False(t, found)
}

func TestFindByNIK_ShouldReturnError_WhenScanFails(t *testing.T) {
	q := &fakeQuerier{vals: []any{"n", "b", 19900101, "s"}} // tglLahir bukan string

	_, found, err := NewResidentRepository(q).FindByNIK(context.Background(), "3273000000000001")

	assert.Error(t, err)
	assert.False(t, found)
}
