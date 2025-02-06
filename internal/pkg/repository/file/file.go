package file

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"word_of_wisdom/internal/pkg/quotes/model"
)

//go:embed quotes.txt
var quotes []byte

type QuoteRepo struct {
	quotes []model.Quote
}

func NewFileRepoQuote() *QuoteRepo {
	qr := new(QuoteRepo)

	reader := bytes.NewReader(quotes)
	s := bufio.NewScanner(reader)

	// Use a regular expression to extract the quote and the author
	re := regexp.MustCompile(`^"([^"]+)"\s*-\s*(.*)$`)

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			quoteText := matches[1]
			author := matches[2]

			qr.quotes = append(qr.quotes, model.Quote{
				Quote:  quoteText,
				Author: author,
			})
		}
	}

	return qr
}

func (q *QuoteRepo) GetQuote(ctx context.Context) (model.Quote, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(q.quotes))))
	if err != nil {
		return model.Quote{}, fmt.Errorf("FileRepo: failed to generate random index: %w", err)
	}
	quote := q.quotes[n.Int64()]
	return quote, nil
}
