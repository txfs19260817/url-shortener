package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHashGenerator_GenerateNanoID(t *testing.T) {
	hashLen := 6
	h := NewHashGenerator(hashLen, 0, 0)
	nanoID, err := h.GenerateNanoID()
	assert.NoError(t, err)
	assert.Len(t, nanoID, hashLen)
	for _, c := range nanoID {
		assert.Contains(t, h.Alphabet, string(c))
	}
}

func TestHashGenerator_NanoIDProvider(t *testing.T) {
	hashLen, amount, duration := 6, 3, 100*time.Millisecond
	h := NewHashGenerator(hashLen, amount, duration)
	go func() {
		tt := t
		err := h.NanoIDProvider()
		assert.NoError(tt, err)
	}()
	time.Sleep(5 * duration)
	h.Done <- struct{}{}
	for i := 0; i < amount; i++ {
		nanoID, ok := <-h.Queue
		assert.True(t, ok)
		t.Log("received nanoID: " + nanoID)
		assert.Len(t, nanoID, hashLen)
	}
	nanoID, ok := <-h.Queue
	assert.False(t, ok)
	assert.Zero(t, nanoID)
	t.Log("Queue closed gracefully")
	done, ok := <-h.Done
	assert.False(t, ok)
	assert.Zero(t, done)
	t.Log("Done closed gracefully")
}
