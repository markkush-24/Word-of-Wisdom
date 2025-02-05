package repo

import (
	"context"
	"word_of_wisdom/internal/pkg/quotes/model"
)

type QuoteRepo interface {
	GetQuote(ctx context.Context) (model.Quote, error)
}
