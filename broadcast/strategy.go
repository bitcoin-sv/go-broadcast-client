package broadcast

import (
	"context"
	"fmt"
)

type StrategyName string

const (
	OneByOneStrategy StrategyName = "OneByOneStrategy"
)

type ExecutionFunc func(context.Context, []func(context.Context) error) error

type Strategy struct {
	name          StrategyName
	executionFunc ExecutionFunc
}

func New(name StrategyName, executionFunc ExecutionFunc) *Strategy {
	return &Strategy{name: name, executionFunc: executionFunc}
}

func (s *Strategy) Execute(ctx context.Context, executionFuncs []func(context.Context) error) error {
	return s.executionFunc(ctx, executionFuncs)
}

var (
	OneByOne = New(OneByOneStrategy, func(ctx context.Context, executionFuncs []func(context.Context) error) error {
		for _, executionFunc := range executionFuncs {
			err := executionFunc(ctx)
			if err != nil {
				continue
			}
			return nil
		}
		// return factory.ErrAllBroadcastersFailed
		return fmt.Errorf("all broadcasters failed")
	})
)
