package verifier

import (
	"fmt"
	"testing"
)

const (
	nonEmptySolution  = "non-empty" // Non-empty string for testing target=0
	invalidSolution   = "abc"       // String whose hash does not start with "0"
	validTarget       = 1           // Target value requiring one leading zero
	zeroTarget        = 0           // Target value with no leading zero requirements
	negativeTarget    = -1          // Negative target value that should cause a panic
	maxIterationCount = 1000000     // Maximum number of iterations to find a valid solution
)

// TestEmptySolution verifies that passing an empty string results in false.
func TestEmptySolution(t *testing.T) { //nolint:paralleltest
	v := NewVerifier()
	if v.VerifyPoW("", validTarget) {
		t.Error("Expected false for an empty solution string")
	}
}

// TestTargetZero verifies that for target=0, any non-empty solution is valid.
func TestTargetZero(t *testing.T) { //nolint:paralleltest
	v := NewVerifier()
	if !v.VerifyPoW(nonEmptySolution, zeroTarget) {
		t.Error("Expected true for any non-empty solution when target is 0")
	}
}

// TestInvalidSolution verifies that the string "abc" does not pass validation for target=1.
func TestInvalidSolution(t *testing.T) { //nolint:paralleltest
	v := NewVerifier()
	if v.VerifyPoW(invalidSolution, validTarget) {
		t.Error("Expected false for solution 'abc' with target=1")
	}
}

// TestValidSolutionTargetOne searches for a solution that meets the target=1 condition,
// meaning its SHA256 hash starts with "0".
func TestValidSolutionTargetOne(t *testing.T) { //nolint:paralleltest
	v := NewVerifier()
	var solution string
	found := false
	// Iterate up to maxIterationCount times until a solution meeting the condition is found
	for i := 0; i < maxIterationCount; i++ {
		solution = fmt.Sprintf("test%d", i)
		if v.VerifyPoW(solution, validTarget) {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Failed to find a valid solution for target=1 within maxIterationCount iterations")
	}
	if !v.VerifyPoW(solution, validTarget) {
		t.Errorf("Solution %s should be valid for target=1", solution)
	}
}

// TestNegativeTarget verifies that a negative target value causes a panic.
func TestNegativeTarget(t *testing.T) { //nolint:paralleltest
	v := NewVerifier()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic for a negative target value")
		}
	}()
	// Passing a negative target value, which should cause a panic.
	v.VerifyPoW(invalidSolution, negativeTarget)
}
