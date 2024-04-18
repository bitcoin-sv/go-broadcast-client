package composite

import (
	"context"
	"errors"
	"testing"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/stretchr/testify/assert"
)

type mockExecutionFunc struct {
	result Result
	err    error
}

func (m mockExecutionFunc) Execute(_ context.Context) (Result, *broadcast.SubmitFailure) {
	if m.err != nil {
		return m.result, broadcast.Failure("", m.err)
	}

	return m.result, nil
}

func TestStrategy_Execute(t *testing.T) {
	// given
	errFunc := mockExecutionFunc{
		err: errors.New("failed"),
	}
	successFunc := mockExecutionFunc{
		result: "success",
	}

	testCases := []struct {
		name     string
		funcs    []executionFunc
		expected Result
		err      error
	}{
		{
			name:     "first execution func should return success",
			funcs:    []executionFunc{successFunc.Execute, errFunc.Execute},
			expected: "success",
			err:      nil,
		},
		{
			name:     "last execution func should return success",
			funcs:    []executionFunc{errFunc.Execute, successFunc.Execute},
			expected: "success",
			err:      nil,
		},
		{
			name:     "all execution funcs should fail",
			funcs:    []executionFunc{errFunc.Execute, errFunc.Execute},
			expected: nil,
			err:      broadcast.ErrAllBroadcastersFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			strategy, err := New(OneByOneStrategy)
			assert.NoError(t, err)

			// when
			result, err := strategy.Execute(context.Background(), tc.funcs)

			// then
			if tc.err == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.err.Error())
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("should return correct strategy for OneByOneStrategy", func(t *testing.T) {
		// given
		expectedStrategyName := OneByOneStrategy

		// when
		actualStrategy, err := New(OneByOneStrategy)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedStrategyName, actualStrategy.name)
	})

	t.Run("should return error for unknown strategy name", func(t *testing.T) {
		// given
		unknownStrategyName := StrategyName("Unknown")

		// when
		_, err := New(unknownStrategyName)

		// then
		assert.Error(t, err)
	})
}
