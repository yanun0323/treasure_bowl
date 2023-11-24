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

func NewAssetProvider(orderConsumer <-chan model.Account, acc model.Account) domain.AssetProvideServer {
	msg := make(chan model.Account, 10)
	msg <- acc
	return &assetProvider{
		l:             logs.New("mock asset provider", util.LogLevel()),
		msg:           msg,
		orderConsumer: orderConsumer,
	}
}

func (p *assetProvider) Connect(ctx context.Context) (<-chan model.Account, error) {
	c, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	go func(ctx context.Context) {
		for {
			select {
			case acc := <-p.orderConsumer:
				p.msg <- acc
			case <-ctx.Done():
				return
			}
		}
	}(c)

	return p.msg, nil
}

func (p *assetProvider) Disconnect(ctx context.Context) error {
	defer close(p.msg)
	p.cancel()
	return nil
}
