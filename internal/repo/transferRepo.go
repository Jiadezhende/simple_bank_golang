package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Jiadezhende/simple_bank_golang/internal/db"
)

var ErrTransferNotFound = errors.New("transfer not found")

type TransferRepo struct {
	db      *sql.DB
	queries *db.Queries
}

func NewTransferRepo(dbConn *sql.DB) *TransferRepo {
	return &TransferRepo{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

// CreateInTx 在已给定 tx 中创建一条 transfer（由调用方负责 Begin/Commit/Rollback）
func (r *TransferRepo) CreateInTx(ctx context.Context, tx *sql.Tx, fromID, toID int64, amount string, createdAt time.Time) (db.Transfer, error) {
	q := db.New(tx)
	tr, err := q.CreateTransfer(ctx, db.CreateTransferParams{
		From:      fromID,
		To:        toID,
		Amount:    amount,
		CreatedAt: createdAt,
	})
	return tr, err
}

// ListByAccount 返回和 account 相关的所有 transfer
func (r *TransferRepo) ListByAccount(ctx context.Context, accountID int64) ([]db.Transfer, error) {
	return r.queries.GetTransfer(ctx, accountID)
}
