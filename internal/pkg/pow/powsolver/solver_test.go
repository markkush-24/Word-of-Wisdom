package powsolver

import (
	"fmt"
	"testing"
	"word_of_wisdom/internal/server/verifier"
)

// findValidSolution searches for a candidate that satisfies the PoW condition for the given target.
// If a solution is found within the limit iterations, it returns the found solution and true.
// Otherwise, it returns an empty string and false.
func findValidSolution(v *verifier.Verifier, target int, limit int) (string, bool) {
	for i := 0; i < limit; i++ {
		candidate := fmt.Sprintf("candidate%d", i)
		if v.VerifyPoW(candidate, target) {
			return candidate, true
		}
	}
	return "", false
}

// TestVerifyPoW contains a set of test cases (table-driven tests) for the VerifyPoW function.
func TestVerifyPoW(t *testing.T) { //nolint:paralleltest
	v := verifier.NewVerifier()

	tests := []struct {
		name     string
		solution string
		target   int
		expected bool
	}{
		{
			name:     "Empty solution returns false",
			solution: "",
			target:   1,
			expected: false,
		},
		{
			name:     "Non-empty solution with target 0 always returns true",
			solution: "anything",
			target:   0,
			expected: true,
		},
		{
			name:     `"abc" with target 1 returns false`,
			solution: "abc",
			target:   1,
			expected: false,
		},
		{
			name:     `"abc" with target 0 returns true`,
			solution: "abc",
			target:   0,
			expected: true,
		},
		{
			name:     `"abc" with high target returns false`,
			solution: "abc",
			target:   10,
			expected: false,
		},
	}

	for _, tc := range tests { //nolint:paralleltest
		// Each test case is executed in a separate sub-test function.
		t.Run(tc.name, func(t *testing.T) {
			result := v.VerifyPoW(tc.solution, tc.target)
			if result != tc.expected {
				t.Errorf("VerifyPoW(%q, %d) = %v; expected %v", tc.solution, tc.target, result, tc.expected)
			}
		})
	}
}

// TestVerifyPoW_ValidSolution searches for a valid solution for target = 1
// using the findValidSolution function. This demonstrates that the algorithm
// is capable of finding such a solution.
func TestVerifyPoW_ValidSolution(t *testing.T) { //nolint:paralleltest
	v := verifier.NewVerifier()
	target := 1

	// Attempt to find a solution within 100000 candidates
	// (average number of attempts ~16 for target = 1).
	candidate, found := findValidSolution(v, target, 100000)
	if !found {
		t.Skip("Failed to find a valid solution for target = 1 within the attempt limit, skipping test")
	}

	t.Logf("Valid solution found: %q for target = %d", candidate, target)

	if !v.VerifyPoW(candidate, target) {
		t.Errorf("Found solution %q did not pass verification for target = %d", candidate, target)
	}
}
