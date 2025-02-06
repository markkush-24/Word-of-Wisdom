package repo

import (
	"context"
	"fmt"
	"word_of_wisdom/internal/pkg/quotes/model"
	"word_of_wisdom/internal/pkg/repository/file"
	"word_of_wisdom/internal/pkg/repository/jsonquote"
)

// RepoSelector for selecting the repository type
type RepoSelector string

const (
	FileRepo RepoSelector = "file"
	JSONRepo RepoSelector = "json"
)

type QuoteRepository struct {
	repo QuoteRepo
}

// NewQuoteRepository creates a repository based on the selected source
func NewQuoteRepository(source RepoSelector) (*QuoteRepository, error) {
	var repo QuoteRepo
	var err error
	switch source {
	case FileRepo:
		repo = file.NewFileRepoQuote()
	case JSONRepo:
		repo, err = jsonquote.NewJSONRepoQuote()
	default:
		return nil, fmt.Errorf("unknown source: %s", source)
	}
	if err != nil {
		return nil, fmt.Errorf("QuoteRepository: failed to create JSON repo: %w", err)
	}
	return &QuoteRepository{repo: repo}, nil
}

// GetQuote delegates the call to the specific repository
func (r *QuoteRepository) GetQuote(ctx context.Context) (model.Quote, error) {
	return r.repo.GetQuote(ctx)
}
