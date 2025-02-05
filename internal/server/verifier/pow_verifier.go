package verifier

type PoWVerifier interface {
	VerifyPoW(solution string, target int) bool
}
