package pow

type PowSolver interface {
	SolvePoW(challenge string) string
}
