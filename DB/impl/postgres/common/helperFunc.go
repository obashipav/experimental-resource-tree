package common

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"log"
)

func HandleTransactionRollback(tx pgx.Tx, ctx context.Context, cancel context.CancelFunc) {
	if tx != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil && rbErr != pgx.ErrTxClosed {
			log.Println("failed to roll back transaction: ", rbErr)
		}
	}
	cancel()
}

func HandleCommit(tx pgx.Tx, ctx context.Context) error {
	err := tx.Commit(ctx)
	if err != nil {
		log.Println("failed to roll back transaction", zap.Error(err))
		return errors.New("failed to commit the tx")
	}
	return nil
}
