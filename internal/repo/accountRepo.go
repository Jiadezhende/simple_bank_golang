package repo

// repo: 面向业务的单表操作

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Jiadezhende/simple_bank_golang/internal/db"
)

var (
	ErrNotFound         = errors.New("account not found")
	ErrInsufficientFund = errors.New("insufficient funds")
)

type AccountRepo struct {
	db      *sql.DB
	queries *db.Queries
}

func NewAccountRepo(dbConn *sql.DB) *AccountRepo {
	return &AccountRepo{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

// Create 使用 sqlc 生成的方法
func (r *AccountRepo) Create(ctx context.Context, owner, balance, currency string) (db.Account, error) {
	return r.queries.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    owner,
		Balance:  balance,
		Currency: currency,
	})
}

// Get 直接返回 sqlc 类型，并映射 sql.ErrNoRows
func (r *AccountRepo) Get(ctx context.Context, id int64) (db.Account, error) {
	acc, err := r.queries.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Account{}, ErrNotFound
		}
		return db.Account{}, err
	}
	return acc, nil
}

// List 封装分页
func (r *AccountRepo) List(ctx context.Context, limit int32, offset int32) ([]db.Account, error) {
	return r.queries.ListAccount(ctx, db.ListAccountParams{Limit: limit, Offset: offset})
}

// UpdateBalanceInTx 在给定 tx 上更新余额（使用 sqlc 的 UpdateBalance 或直接在 tx 上执行）
func (r *AccountRepo) UpdateBalanceInTx(ctx context.Context, tx *sql.Tx, id int64, balance string) (db.Account, error) {
	q := db.New(tx) // sqlc 生成的 constructor 支持 *sql.Tx
	acc, err := q.UpdateBalance(ctx, db.UpdateBalanceParams{ID: id, Balance: balance})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Account{}, ErrNotFound
		}
		return db.Account{}, err
	}
	return acc, nil
}

// DebitInTx 原子扣款：在 tx 内执行带余额检查的 UPDATE ... RETURNING
func (r *AccountRepo) DebitInTx(ctx context.Context, tx *sql.Tx, id int64, amount string) (db.Account, error) {
	var acc db.Account
	row := tx.QueryRowContext(ctx, `
        UPDATE accounts
        SET balance = balance - $1
        WHERE id = $2 AND balance >= $1
        RETURNING id, owner, balance, currency, created_at
    `, amount, id)
	if err := row.Scan(&acc.ID, &acc.Owner, &acc.Balance, &acc.Currency, &acc.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Account{}, ErrInsufficientFund
		}
		return db.Account{}, err
	}
	return acc, nil
}

// CreditInTx 原子加款
func (r *AccountRepo) CreditInTx(ctx context.Context, tx *sql.Tx, id int64, amount string) (db.Account, error) {
	var acc db.Account
	row := tx.QueryRowContext(ctx, `
        UPDATE accounts
        SET balance = balance + $1
        WHERE id = $2
        RETURNING id, owner, balance, currency, created_at
    `, amount, id)
	if err := row.Scan(&acc.ID, &acc.Owner, &acc.Balance, &acc.Currency, &acc.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Account{}, ErrNotFound
		}
		return db.Account{}, err
	}
	return acc, nil
}
