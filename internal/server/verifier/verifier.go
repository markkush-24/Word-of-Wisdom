package verifier

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// Verifier implements the PoWVerifier interface
type Verifier struct{}

func NewVerifier() *Verifier {
	return &Verifier{}
}

// VerifyPoW verifies the solution to the Proof of Work problem.
// The input parameter is the solution string to check.
// target - the number of leading zeros in the hash string.
func (v *Verifier) VerifyPoW(solution string, target int) bool {
	// Check if the solution length is greater than 0
	if len(solution) == 0 {
		return false
	}

	// Hash the solution
	hash := sha256.New()
	hash.Write([]byte(solution))
	hashValue := fmt.Sprintf("%x", hash.Sum(nil))

	// Check if the hash is empty
	if len(hashValue) == 0 {
		return false
	}

	// Check if the hash starts with the required number of leading zeros
	expectedPrefix := strings.Repeat("0", target)
	return strings.HasPrefix(hashValue, expectedPrefix)
}
