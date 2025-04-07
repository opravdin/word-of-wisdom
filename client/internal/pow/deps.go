package pow

import "context"

//go:generate mockgen -source=${GOFILE} -package=mocks -destination=./mocks/deps.go

// Solver defines the interface for PoW challenge solving
type Solver interface {
	// Solve solves a PoW challenge with the given parameters
	Solve(challengeID, seed string, difficultyLevel, n, r, p, keyLen int) (string, error)

	// SolveWithContext solves a PoW challenge with a context for cancellation/timeout
	SolveWithContext(ctx context.Context, challengeID, seed string, difficultyLevel, n, r, p, keyLen int) (string, error)

	// SolveWithDefaults solves a PoW challenge using default parameters
	SolveWithDefaults(challengeID, seed string) (string, error)
}

// SolverFactory defines the interface for creating PoW solvers
type SolverFactory interface {
	// NewSolver creates a new PoW solver
	NewSolver() Solver
}
