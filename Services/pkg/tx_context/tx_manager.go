package txcontext

import (
	"context"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

type contextKey string

const txKey contextKey = "tx"

func (m *TxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) []helpers.ResponseError) []helpers.ResponseError {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		ctxWithTx := NewContextWithTx(ctx, tx)
		if errs := fn(ctxWithTx); len(errs) > 0 {
			return &txErrorWrapper{errors: errs}
		}
		return nil
	})

	if err == nil {
		return nil
	}

	if txErr, ok := err.(*txErrorWrapper); ok {
		return txErr.errors
	}

	return []helpers.ResponseError{
		{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Transaction Failed",
			Message: err.Error(),
		},
	}
}

type txErrorWrapper struct {
	errors []helpers.ResponseError
}

func (e *txErrorWrapper) Error() string {
	return "transaction failed"
}

func NewContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTxFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return nil
}
