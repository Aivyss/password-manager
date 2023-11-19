package repository

import (
	"context"
	"database/sql"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/jmoiron/sqlx"
)

type TxBlock func(ctx context.Context, tx *sqlx.Tx) error
type TxBlockAuto func(ctx context.Context) error

type TxManager interface {
	Txx(ctx context.Context, txBlock TxBlockAuto) error
	TxxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlockAuto) error
	Tx(ctx context.Context, txBlock TxBlock) error
	TxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlock) error
}

func NewTxManager(writeDB *sqlx.DB) TxManager {
	return &defaultTxManager{
		writeDB: writeDB,
	}
}

type defaultTxManager struct {
	writeDB *sqlx.DB
}

func (d *defaultTxManager) Txx(ctx context.Context, txBlock TxBlockAuto) error {
	return d.TxxWithOpt(ctx, nil, txBlock)
}

func (d *defaultTxManager) TxxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlockAuto) (err error) {
	tx, dbErr := d.writeDB.BeginTxx(ctx, opts)
	if dbErr != nil {
		return dbErr
	}

	defer func() {
		rec := recover()
		if rec != nil {
			err2, ok := rec.(error)

			if ok {
				err = err2
			} else {
				err = pwmErr.Unknown
			}

			_ = tx.Rollback()
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	ctx = NewTxContext(ctx, tx)
	clientSideErr := txBlock(ctx)
	if clientSideErr != nil {
		err = clientSideErr
	}

	return err
}

func (d *defaultTxManager) Tx(ctx context.Context, txBlock TxBlock) (err error) {
	return d.TxWithOpt(ctx, nil, txBlock)
}

func (d *defaultTxManager) TxWithOpt(ctx context.Context, opts *sql.TxOptions, txBlock TxBlock) (err error) {
	tx, dbErr := d.writeDB.BeginTxx(ctx, opts)
	if dbErr != nil {
		return err
	}

	defer func() {
		rec := recover()
		if rec != nil {
			err2, ok := rec.(error)

			if ok {
				err = err2
			} else {
				err = pwmErr.Unknown
			}

			_ = tx.Rollback()
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	clientSideErr := txBlock(ctx, tx)
	if clientSideErr != nil {
		err = clientSideErr
	}

	return err
}
