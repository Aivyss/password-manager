package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type TxContext struct {
	context.Context
	*sqlx.Tx
}

func NewTxContext(ctx context.Context, tx *sqlx.Tx) context.Context {
	return &TxContext{
		Context: ctx,
		Tx:      tx,
	}
}

func GetStatement(ctx context.Context, statement *sqlx.NamedStmt) *sqlx.NamedStmt {
	txContext, ok := ctx.(TxContext)

	if ok {
		return txContext.NamedStmt(statement)
	}

	return statement
}
