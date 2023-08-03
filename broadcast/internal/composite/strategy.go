package composite

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

type StrategyName string

const (
	OneByOneStrategy StrategyName = "OneByOneStrategy"
)

type Result interface{}

type executionFunc func(context.Context) (Result, error)
type StrategyExecutionFunc func(context.Context, []executionFunc) (Result, error)

type Strategy struct {
	name          StrategyName
	executionFunc StrategyExecutionFunc
}

func New(name StrategyName) (*Strategy, error) {
	switch name {
	case OneByOneStrategy:
		return &Strategy{name: name, executionFunc: OneByOne.executionFunc}, nil
	default:
		return nil, broadcast.ErrStrategyUnkown
	}
}

func (s *Strategy) Execute(ctx context.Context, executionFuncs []executionFunc) (Result, error) {
	return s.executionFunc(ctx, executionFuncs)
}

var (
	OneByOne = &Strategy{name: OneByOneStrategy, executionFunc: func(ctx context.Context, executionFuncs []executionFunc) (Result, error) {
		for _, executionFunc := range executionFuncs {
			result, err := executionFunc(ctx)
			if err != nil {
				continue
			}
			return result, nil
		}
		return nil, broadcast.ErrAllBroadcastersFailed
	}}
)
