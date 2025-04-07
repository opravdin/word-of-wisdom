package quotes

import (
	"context"
	"math/rand"
	"sync"

	"github.com/opravdin/word-of-wisdom/internal/domain"
	"github.com/opravdin/word-of-wisdom/internal/logger"
)

type InMemoryRepository struct {
	quotes []domain.Quote
	mu     sync.RWMutex
	logger logger.Logger
}

func NewInMemoryRepository(log logger.Logger) *InMemoryRepository {
	repo := &InMemoryRepository{
		quotes: []domain.Quote{
			{Text: "The only true wisdom is in knowing you know nothing.", Author: "Socrates"},
			{Text: "The unexamined life is not worth living.", Author: "Socrates"},
			{Text: "Wisdom begins in wonder.", Author: "Socrates"},
			{Text: "Knowledge speaks, but wisdom listens.", Author: "Jimi Hendrix"},
			{Text: "The journey of a thousand miles begins with one step.", Author: "Lao Tzu"},
			{Text: "By three methods we may learn wisdom: by reflection, which is noblest; by imitation, which is easiest; and by experience, which is the bitterest.", Author: "Confucius"},
			{Text: "Turn your wounds into wisdom.", Author: "Oprah Winfrey"},
			{Text: "The more I read, the more I acquire, the more certain I am that I know nothing.", Author: "Voltaire"},
			{Text: "It is the mark of an educated mind to be able to entertain a thought without accepting it.", Author: "Aristotle"},
			{Text: "The fool doth think he is wise, but the wise man knows himself to be a fool.", Author: "William Shakespeare"},
		},
		logger: log,
	}

	log.Info("Initialized quotes repository", "count", len(repo.quotes))
	return repo
}

func (r *InMemoryRepository) GetRandom(ctx context.Context) domain.Quote {
	r.mu.RLock()
	defer r.mu.RUnlock()

	index := rand.Intn(len(r.quotes))
	quote := r.quotes[index]

	r.logger.Debug("Retrieved random quote", "index", index, "author", quote.Author)
	return quote
}
