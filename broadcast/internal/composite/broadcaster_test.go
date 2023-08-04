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

func TestNewBroadcasterWithDefaultStrategy(t *testing.T) {
	// Given
	mockFactory := &MockBroadcastFactory{}

	// When
	broadcaster := NewBroadcasterWithDefaultStrategy(mockFactory.Create)

	// Then
	assert.NotNil(t, broadcaster)
	_, ok := broadcaster.(*compositeBroadcaster)
	assert.True(t, ok, "Expected broadcaster to be of type *compositeBroadcaster")
}

func TestNewBroadcaster(t *testing.T) {
	// Given
	mockFactory := &MockBroadcastFactory{}

	// When
	broadcaster := NewBroadcaster(*OneByOne, mockFactory.Create)

	// Then
	assert.NotNil(t, broadcaster)
	_, ok := broadcaster.(*compositeBroadcaster)
	assert.True(t, ok, "Expected broadcaster to be of type *compositeBroadcaster")
}
