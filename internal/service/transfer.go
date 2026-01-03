// ...existing code...
package service

import (
	"context"
	"database/sql"

	"github.com/Jiadezhende/simple_bank_golang/internal/db"
	"github.com/Jiadezhende/simple_bank_golang/internal/store"
)

type TransferService struct {
	DB    *sql.DB
	Store *store.Store
}

func NewTransferService(db *sql.DB, s *store.Store) *TransferService {
	return &TransferService{DB: db, Store: s}
}

type TransferResult struct {
	From *db.Account `json:"from"`
	To   *db.Account `json:"to"`
	ID   int64       `json:"id"`
}

// Transfer: 在 service 管理事务、做幂等/权限校验后调用 store 的原子操作
func (s *TransferService) Transfer(ctx context.Context, fromID, toID int64, amount string, currency string) (*TransferResult, error) {
	// 可在此做参数校验 / 鉴权 / 幂等校验

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	fromAcc, toAcc, txID, err := s.Store.TransferInTx(ctx, tx, fromID, toID, amount, currency)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &TransferResult{
		From: fromAcc,
		To:   toAcc,
		ID:   txID,
	}, nil
}
