package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repositiry struct {
	db *sqlx.DB
}

func NewRepositiry(db *sqlx.DB) *Repositiry {
	return &Repositiry{db: db}
}

func (r *Repositiry) AddAverageGas(blockHash string, blockNumber, averageGasPrice uint64) error {
	// TODO
	return nil
}

func (r *Repositiry) AddERC20Logs(blockNumber, logIndex uint64, blockHash, eventName, from, to, value string) error {
	// TODO
	return nil
}
