package store

// 复合操作集中于store

import (
	"database/sql"

	"github.com/Jiadezhende/simple_bank_golang/internal/repo"
)

type Store struct {
	DB          *sql.DB
	AccountRepo *repo.AccountRepo
	EntryRepo   *repo.EntryRepo
	TxRepo      *repo.TransferRepo
}

func NewStore(dbConn *sql.DB, a *repo.AccountRepo, e *repo.EntryRepo, t *repo.TransferRepo) *Store {
	return &Store{
		DB:          dbConn,
		AccountRepo: a,
		EntryRepo:   e,
		TxRepo:      t,
	}
}
