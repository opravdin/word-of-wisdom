package getquote

import (
	"context"

	"github.com/opravdin/word-of-wisdom/internal/domain"
)

//go:generate mockgen -source=${GOFILE} -package=mocks -destination=./mocks/deps.go

type quotesRepository interface {
	GetRandom(ctx context.Context) domain.Quote
}
