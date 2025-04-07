package configuration

import (
	"github.com/opravdin/word-of-wisdom/internal/configuration/env"
	"github.com/opravdin/word-of-wisdom/internal/infrastructure"
	"github.com/opravdin/word-of-wisdom/internal/logger"
	"github.com/opravdin/word-of-wisdom/internal/pow"
	powrepo "github.com/opravdin/word-of-wisdom/internal/repository/pow"
	"github.com/opravdin/word-of-wisdom/internal/repository/quotes"
	"github.com/opravdin/word-of-wisdom/internal/usecase/getquote"
)

type DefaultConfiguration struct {
	GetQuote   getquote.DefaultUsecase
	Config     *env.AppConfig
	PowService *pow.Service
}

func NewDefaultConfiguration(
	storage infrastructure.Storage,
	config *env.AppConfig,
	log logger.Logger,
) *DefaultConfiguration {
	// Repository
	quoutesRepository := quotes.NewInMemoryRepository(log)

	// Usecase
	getQuoteUsecase := getquote.NewDefaultUsecase(quoutesRepository, log)

	// Initialize PoW repository and service
	powRepository := powrepo.NewRepository(storage.Redis, config.Pow.BucketCapacity, log)
	powService := pow.NewService(powRepository, &config.Pow, log)

	return &DefaultConfiguration{
		GetQuote:   *getQuoteUsecase,
		Config:     config,
		PowService: powService,
	}
}
