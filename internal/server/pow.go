package server

type PoWVerifier interface {
	VerifyPoW(solution string, target int) bool
}
