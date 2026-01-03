package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/Jiadezhende/simple_bank_golang/internal/db"
)

type EntryRepo struct {
	db      *sql.DB
	queries *db.Queries
}

func NewEntryRepo(dbConn *sql.DB) *EntryRepo {
	return &EntryRepo{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

// CreateInTx 在已给定 tx 中创建 entry；transferID 可为 sql.NullInt64，entryType 可为空表示 NULL
func (r *EntryRepo) CreateInTx(ctx context.Context, tx *sql.Tx, accountID int64, transferID sql.NullInt64, amount string, entryType string, currency string, createdAt time.Time) (db.Entry, error) {
	q := db.New(tx)
	var et sql.NullString
	if entryType != "" {
		et = sql.NullString{String: entryType, Valid: true}
	} else {
		et = sql.NullString{Valid: false}
	}
	return q.CreateEntry(ctx, db.CreateEntryParams{
		AccountID:  accountID,
		TransferID: transferID,
		Amount:     amount,
		EntryType:  et,
		Currency:   currency,
		CreatedAt:  createdAt,
	})
}
