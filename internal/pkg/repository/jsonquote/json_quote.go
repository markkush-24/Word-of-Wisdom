package jsonquote

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
	"word_of_wisdom/internal/pkg/quotes/model"
)

//go:embed quotes.json
var s []byte

type JSONQuote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

type Repo struct {
	quotes []model.Quote
}

func NewJSONRepoQuote() (*Repo, error) {
	var jsonQuotes []JSONQuote
	err := json.Unmarshal(s, &jsonQuotes)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling quotes json: %w", err)
	}

	modelQuotes := make([]model.Quote, len(jsonQuotes))
	for i, j := range jsonQuotes {
		modelQuotes[i] = model.Quote{
			Quote:  j.Quote,
			Author: j.Author,
		}
	}

	return &Repo{
		quotes: modelQuotes,
	}, nil
}

func (r *Repo) GetQuote(ctx context.Context) (model.Quote, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(r.quotes))))
	if err != nil {
		return model.Quote{}, fmt.Errorf("JSONRepo: failed to generate secure random number: %w", err)
	}
	return r.quotes[n.Int64()], nil
}
