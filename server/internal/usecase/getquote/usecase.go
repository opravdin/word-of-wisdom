package getquote

import (
	"context"

	"github.com/opravdin/word-of-wisdom/internal/domain"
	"github.com/opravdin/word-of-wisdom/internal/logger"
)

type DefaultUsecase struct {
	quotesRepo quotesRepository
	logger     logger.Logger
}

func NewDefaultUsecase(quotesRepo quotesRepository, log logger.Logger) *DefaultUsecase {
	return &DefaultUsecase{
		quotesRepo: quotesRepo,
		logger:     log,
	}
}

func (u DefaultUsecase) Quote(ctx context.Context) domain.Quote {
	u.logger.Debug("Getting random quote")
	quote := u.quotesRepo.GetRandom(ctx)
	u.logger.Debug("Got random quote", "author", quote.Author)
	return quote
}
