package random

import (
	"math/rand"
)

// Provider defines the interface for random number generation
type Provider interface {
	// Float32 returns a random float32 in the range [0.0, 1.0)
	Float32() float32
}

// DefaultProvider is the default implementation of Provider
type DefaultProvider struct{}

// NewProvider creates a new DefaultProvider
func NewProvider() Provider {
	return &DefaultProvider{}
}

// Float32 returns a random float32 in the range [0.0, 1.0)
func (p *DefaultProvider) Float32() float32 {
	return rand.Float32()
}
