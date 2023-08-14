package composite

import (
	"testing"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	broadcast.Client
}

type MockBroadcastFactory struct{}

func (m *MockBroadcastFactory) Create() broadcast.Client {
	return &MockClient{}
}

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
