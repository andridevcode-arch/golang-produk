package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// k kode error PostgreSQL yang sering dibutuhkan
const (
	pgCodeUniqueViolation      = "23505"
	pgCodeForeignKeyViolation  = "23503"
	pgCodeCheckViolation       = "23514"
	pgCodeNotNullViolation     = "23502"
	pgCodeExclusionViolation   = "23P01"
	pgCodeSerializationFailure = "40001"
	pgCodeDeadlockDetected     = "40P01"
)

// IsUniqueViolation mengecek apakah error berasal dari pelanggaran UNIQUE constraint
func IsUniqueViolation(err error, constraint string) bool {
	return isPgViolation(err, pgCodeUniqueViolation, constraint)
}

// IsForeignKeyViolation mengecek apakah error berasal dari pelanggaran FOREIGN KEY
func IsForeignKeyViolation(err error, constraint string) bool {
	return isPgViolation(err, pgCodeForeignKeyViolation, constraint)
}

// IsCheckViolation mengecek apakah error berasal dari pelanggaran CHECK constraint
func IsCheckViolation(err error, constraint string) bool {
	return isPgViolation(err, pgCodeCheckViolation, constraint)
}

// IsNotNullViolation mengecek apakah error berasal dari pelanggaran NOT NULL
func IsNotNullViolation(err error, column string) bool {
	return isPgViolation(err, pgCodeNotNullViolation, column)
}

// IsSerializationFailure mengecek apakah error karena conflict transaksi
// sering terjadi pada retry pembayaran atau pembaruan kuota
func IsSerializationFailure(err error) bool {
	return isPgViolation(err, pgCodeSerializationFailure, "") ||
		isPgViolation(err, pgCodeDeadlockDetected, "")
}

// IsExclusionViolation mengecek apakah error berasal dari pelanggaran EXCLUDE constraint
func IsExclusionViolation(err error, constraint string) bool {
	return isPgViolation(err, pgCodeExclusionViolation, constraint)
}

// isPgViolation memeriksa pgconn.PgError berdasarkan kode dan constraint,
// constraint boleh string kosong untuk memeriksa kodenya saja
func isPgViolation(err error, code, constraint string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	if pgErr.Code != code {
		return false
	}
	if constraint == "" {
		return true
	}
	return pgErr.ConstraintName == constraint || pgErr.ColumnName == constraint
}
