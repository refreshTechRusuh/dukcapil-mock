# dukcapil-mock

Mock service validasi identitas (Dukcapil) untuk **Rusuh Booking Tiket**.
Menyimulasikan pihak eksternal Dukcapil: menerima `{nik, nama, tglLahir, namaIbu}`,
mencocokkan ke data kependudukan tiruan di PostgreSQL, lalu menjawab `{valid: bool}`.

Dipakai oleh User Service & Booking Service saat validasi NIK (booking / tambah member).

- Porsi: **BE-1 (Identity & Flow)**
- Arsitektur: Hexagonal (Ports & Adapters), `net/http` + `pgx`
- Port default: **8085**

## Prasyarat

- Go 1.26+
- PostgreSQL (disediakan infra bersama di repo `tools` via docker-compose)

## Menjalankan

```bash
# 1. Nyalakan DB bersama (dari repo tools, sekali saja)
cd tools && docker compose up -d postgres

# 2. Jalankan service (dari repo ini)
make run        # = go run ./cmd  (pakai env default → localhost:5432)
```

Service listen di `:8085`. Berhenti dengan `Ctrl+C`.

## Endpoint

| Method | Path                         | Keterangan                          |
|--------|------------------------------|-------------------------------------|
| POST   | `/internal/dukcapil/verify`  | Validasi identitas (cocokkan NIK)   |
| GET    | `/sys/ping`                  | Health check                        |

### Contoh: verify

Request:
```json
POST /internal/dukcapil/verify
{ "nik": "3273000000000001", "nama": "Budi Santoso", "tglLahir": "1990-01-01", "namaIbu": "Siti Aminah" }
```

Response (cocok):
```json
{ "status": "00000", "message": "OK", "data": { "valid": true } }
```

Catatan:
- `valid: false` (NIK tidak cocok / tidak ada) tetap **HTTP 200** — itu hasil
  pengecekan yang berhasil, bukan error request.
- Body JSON rusak → **HTTP 400** (`status: "400"`).
- Pencocokan toleran spasi ujung & huruf besar/kecil.

## Konfigurasi (environment variable)

| Env            | Default         | Keterangan                    |
|----------------|-----------------|-------------------------------|
| `APP_PORT`     | `8085`          | Port HTTP                     |
| `APP_ENV`      | `local`         | `local` / `preprod` / `production` |
| `LOG_LEVEL`    | `info`          | `debug` / `info` / `warn` / `error` |
| `LOG_FORMAT`   | `json`          | `json` / `text`               |
| `DB_HOST`      | `localhost`     | Host PostgreSQL               |
| `DB_PORT`      | `5432`          | Port PostgreSQL               |
| `DB_USER`      | `postgres`      | User DB                       |
| `DB_PASSWORD`  | `password`      | Password DB                   |
| `DB_NAME`      | `rusuh_booking` | Nama database (bersama)       |
| `DB_SSL_MODE`  | `disable`       | SSL mode                      |

## Data seed

Tabel `residents` di-seed oleh init-script compose `tools` (`SELECT * FROM residents`).
Contoh NIK: `3273000000000001` (Budi Santoso, 1990-01-01, ibu Siti Aminah) dan 4 lainnya.

## Perintah Makefile

| Target                | Keterangan                          |
|-----------------------|-------------------------------------|
| `make run`            | Jalankan via `go run ./cmd`         |
| `make build`          | Compile binary ke `bin/api`         |
| `make test`           | Jalankan test                       |
| `make unit-test-coverage` | Test + laporan coverage         |
| `make golint`         | Lint (`golangci-lint`)              |
| `make docker-build`   | Build image `dukcapil-mock:latest`  |

## Struktur

```
cmd/main.go                            # entry point: startup + graceful shutdown
internal/
  config/                              # baca env → Config, SetupLogger
  delivery/server.go                   # composition root (wiring + route)
  dukcapil/
    domain/                            # entity & dto murni
    port/                              # interface (outbound)
    application/                       # usecase (logika Verify) + test
    infrastructure/
      db/                              # adapter PostgreSQL (pgx)
      http/                            # adapter HTTP (handler)
pkg/
  apinizer/                            # envelope response {status,message,data}
  db/                                  # connection pool (pgxpool)
migrations/                            # skema + seed residents
```

> Bedah kode lengkap (untuk belajar): `docs/self/DUKCAPIL-EXPLAINED.md` di repo dokumentasi.
