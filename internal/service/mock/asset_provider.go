package mock

import (
	"context"
	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"

	"github.com/yanun0323/pkg/logs"
)

var (
	/* Check Interface Implement */
	_ domain.AssetProvideServer = (*assetProvider)(nil)
)

type assetProvider struct {
	l             logs.Logger
	msg           chan model.Account
	orderConsumer <-chan model.Account
	cancel        context.CancelFunc
}

func AssetProvider(orderConsumer <-chan model.Account, acc model.Account) domain.AssetProvideServer {
	msg := make(chan model.Account, 10)
	msg <- acc
	return &assetProvider{
		l:             logs.New("mock asset provider", util.LogLevel()),
		msg:           msg,
		orderConsumer: orderConsumer,
	}
}

func (s *assetProvider) Connect(ctx context.Context) (<-chan model.Account, error) {
	c, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	go func(ctx context.Context) {
		for {
			select {
			case acc := <-s.orderConsumer:
				s.msg <- acc
			case <-ctx.Done():
				return
			}
		}
	}(c)

	return s.msg, nil
}

func (s *assetProvider) Disconnect(ctx context.Context) error {
	defer close(s.msg)
	s.cancel()
	return nil
}
