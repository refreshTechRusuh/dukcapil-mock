-- migrations/001_initial_schema.sql
-- Dukcapil Mock: tabel data kependudukan tiruan.
-- verify mencocokkan {nik, nama, tgl_lahir, nama_ibu} ke salah satu baris di sini.

CREATE TABLE IF NOT EXISTS residents (
    nik        CHAR(16) PRIMARY KEY,
    nama       TEXT NOT NULL,
    tgl_lahir  DATE NOT NULL,
    nama_ibu   TEXT NOT NULL
);

-- Seed data contoh untuk dev/demo. NIK harus 16 digit.
INSERT INTO residents (nik, nama, tgl_lahir, nama_ibu) VALUES
    ('3273000000000001', 'Budi Santoso',   '1990-01-01', 'Siti Aminah'),
    ('3273000000000002', 'Ani Lestari',    '1992-05-17', 'Rina Marlina'),
    ('3273000000000003', 'Candra Wijaya',  '1988-11-30', 'Dewi Kartika'),
    ('3273000000000004', 'Dedi Pratama',   '1995-07-09', 'Eka Yuliana'),
    ('3273000000000005', 'Eka Putri',      '2000-03-25', 'Fitri Handayani')
ON CONFLICT (nik) DO NOTHING;
