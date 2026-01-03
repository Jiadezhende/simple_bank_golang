package store

// 转账相关复合业务

import (
	"context"
	"database/sql"
	"time"

	"github.com/Jiadezhende/simple_bank_golang/internal/db"
)

// TransferInTx 在已有 tx 上执行转账的复合原子操作（不自行 Begin/Commit）
// 返回更新后的 from/to 账户和 transfer id
func (s *Store) TransferInTx(ctx context.Context, tx *sql.Tx, fromID, toID int64, amount string, currency string) (*db.Account, *db.Account, int64, error) {
	createdAt := time.Now().UTC()

	// 1) 创建 transfer（使用传入的 createdAt，若后续出错，tx 会回滚）
	tr, err := s.TxRepo.CreateInTx(ctx, tx, fromID, toID, amount, createdAt)
	if err != nil {
		return nil, nil, 0, err
	}

	// 2) 扣款（原子检查余额，repo 在 tx 内执行）
	fromAcc, err := s.AccountRepo.DebitInTx(ctx, tx, fromID, amount)
	if err != nil {
		return nil, nil, 0, err
	}

	// 3) 加款
	toAcc, err := s.AccountRepo.CreditInTx(ctx, tx, toID, amount)
	if err != nil {
		return nil, nil, 0, err
	}

	// 4) 创建 entries（使用 transfer 的 id 与统一的 created_at）
	transferID := sql.NullInt64{Int64: tr.ID, Valid: true}
	if _, err := s.EntryRepo.CreateInTx(ctx, tx, fromID, transferID, amount, "debit", currency, createdAt); err != nil {
		return nil, nil, 0, err
	}
	if _, err := s.EntryRepo.CreateInTx(ctx, tx, toID, transferID, amount, "credit", currency, createdAt); err != nil {
		return nil, nil, 0, err
	}

	return &fromAcc, &toAcc, tr.ID, nil
}
