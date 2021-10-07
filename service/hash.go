package service

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var HashProvider *HashGenerator

// ServicesConfig contains configuration for service package
type ServicesConfig struct {
	// HashGenerator configs
	HashLen  int `yaml:"hash_len"`
	PoolSize int `yaml:"pool_size"`
	Duration int `yaml:"duration"`
}

// HashGenerator provides methods to generate Hash/ID for short urls
type HashGenerator struct {
	Len      int    // Len represents the expected length of the hash
	Alphabet string // Alphabet limited available symbols can be used in hash generation, which should be URL friendly
	*HashPool
}

// HashPool loads some ready-to-use hashes for consumers
type HashPool struct {
	Queue    chan string   // Queue receive hash from the provider and emit hash to consumers
	Done     chan struct{} // Done is a semaphore to end the hash generation task
	Duration time.Duration // Duration is the interval between two generation func calls
}

// NewHashGenerator returns a HashGenerator instance
func NewHashGenerator(hashLen, poolSize int, genDuration time.Duration) *HashGenerator {
	return &HashGenerator{
		Len:      hashLen,
		Alphabet: `_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`,
		HashPool: &HashPool{
			Queue:    make(chan string, poolSize),
			Done:     make(chan struct{}),
			Duration: genDuration,
		},
	}
}

// GenerateNanoID provides the hash in the manner of a generated nanoID
func (h *HashGenerator) GenerateNanoID() (string, error) {
	return gonanoid.Generate(h.Alphabet, h.Len)
}

// NanoIDProvider keeps pushing nanoIDs to HashPool
func (h *HashGenerator) NanoIDProvider() (err error) {
	var nanoID string
	defer close(h.Done)
	defer close(h.Queue)
	for {
		nanoID, err = h.GenerateNanoID()
		if err != nil {
			break
		}
		select {
		case h.HashPool.Queue <- nanoID:
		case <-h.Done:
			return nil
		}
		time.Sleep(h.Duration)
	}
	return err
}

// CloseProvider send signal to HashGenerator.Done to finish the provider task
func (h *HashGenerator) CloseProvider() {
	h.Done <- struct{}{}
}
