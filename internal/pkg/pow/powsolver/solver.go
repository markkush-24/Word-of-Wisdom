package powsolver

import (
	"crypto/sha256"
	"fmt"
	"strings"
	core "word_of_wisdom/config"
	logger "word_of_wisdom/internal/pkg/logging"
)

// SimplePoW struct for storing Proof-of-Work algorithm parameters.
type SimplePoW struct {
	Target int // Number of leading zeros required in the hash.
	logger logger.Logger
}

// NewSimplePoW creates a new instance of SimplePoW.
func NewSimplePoW(target int, log logger.Logger) *SimplePoW {
	return &SimplePoW{
		Target: target,
		logger: log,
	}
}

// calculateHash computes the SHA-256 hash for a given string.
func calculateHash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// SolvePoW solves the Proof of Work challenge using a unique challenge string.
// The candidate solution is formed as: "challenge:solutionPrefix<nonce>".
func (p *SimplePoW) SolvePoW(challenge string) string {
	for i := 0; ; i++ {
		// Формируем кандидат, объединяя уникальный вызов (challenge), фиксированный префикс и счетчик (nonce)
		candidate := fmt.Sprintf("%s:%s%d", challenge, core.SolutionPrefix, i)
		hashValue := calculateHash(candidate)

		// Если хэш начинается с требуемого количества ведущих символов, решение найдено.
		if strings.HasPrefix(hashValue, strings.Repeat(core.HashPrefixCharacter, p.Target)) {
			p.logger.Info(fmt.Sprintf("PoW solution found: %s with hash %s", candidate, hashValue))
			return candidate
		}

		// Логирование каждые 100,000 итераций
		if i%core.CheckpointInterval == 0 {
			p.logger.Info(fmt.Sprintf("Checked %d solutions for challenge %s, continuing to search...", i, challenge))
		}
	}
}
