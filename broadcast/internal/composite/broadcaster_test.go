package composite

import (
	"testing"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/stretchr/testify/assert"
)

// MockClient to mock broadcast client
type MockClient struct {
	broadcast.Client
}

// MockBroadcastFactory to mock broadcast factory
type MockBroadcastFactory struct{}

// Create mock client
func (m *MockBroadcastFactory) Create() broadcast.Client {
	return &MockClient{}
}

// TestNewBroadcasterWithDefaultStrategy tests the NewBroadcasterWithDefaultStrategy function.
func TestNewBroadcasterWithDefaultStrategy(t *testing.T) {
	// given
	mockFactory := &MockBroadcastFactory{}

	// when
	broadcaster := NewBroadcasterWithDefaultStrategy(mockFactory.Create)

	// then
	assert.NotNil(t, broadcaster)
	_, ok := broadcaster.(*compositeBroadcaster)
	assert.True(t, ok, "Expected broadcaster to be of type *compositeBroadcaster")
}

// TestNewBroadcaster tests the NewBroadcaster function.
func TestNewBroadcaster(t *testing.T) {
	// given
	mockFactory := &MockBroadcastFactory{}

	// when
	broadcaster := NewBroadcaster(*OneByOne, mockFactory.Create)

	// then
	assert.NotNil(t, broadcaster)
	_, ok := broadcaster.(*compositeBroadcaster)
	assert.True(t, ok, "Expected broadcaster to be of type *compositeBroadcaster")
}
