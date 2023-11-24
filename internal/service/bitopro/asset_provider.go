// documentation
// auth: https://github.com/bitoex/bitopro-offical-api-docs/blob/master/README.md#api-security-protocol
// consume: https://github.com/bitoex/bitopro-offical-api-docs/blob/master/ws/private/user_balance_stream.md
package bitopro

import (
	"context"
	"sync/atomic"

	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"

	"github.com/bitoex/bitopro-api-go/pkg/ws"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

type AssetProvider struct {
	l              logs.Logger
	connected      *atomic.Bool
	cancel         chan struct{}
	cancelFunc     context.CancelFunc
	accountChannel chan model.Account
	wss            *ws.Ws
}

func NewAssetProvider(wss *ws.Ws) (domain.AssetProvideServer, error) {
	return &AssetProvider{
		l:              logs.New("bitopro asset provider", util.LogLevel()),
		connected:      &atomic.Bool{},
		accountChannel: make(chan model.Account, 100),
		wss:            wss,
	}, nil
}

func (p *AssetProvider) Connect(ctx context.Context) (<-chan model.Account, error) {
	p.connected.Store(true)
	ch, cancel := p.wss.RunAccountBalancesWsConsumer(ctx)
	p.cancel = cancel

	c, cancelFunc := context.WithCancel(ctx)
	p.cancelFunc = cancelFunc
	go p.consumeAccountBalance(c, ch)
	return p.accountChannel, nil
}

func (p *AssetProvider) Disconnect(ctx context.Context) error {
	connected := p.connected.Load()
	if !connected {
		return nil
	}
	p.connected.Store(false)
	go func() {
		defer close(p.cancel)
		p.cancel <- struct{}{}
	}()

	p.cancelFunc()

	return nil
}

func (p *AssetProvider) consumeAccountBalance(ctx context.Context, ch <-chan ws.AccountBalanceData) {
	for {
		select {
		case acc := <-ch:
			p.l.Debugf("consume account balance: %+v", acc.Data)
			p.accountChannel <- transAccData(&acc)
		case <-ctx.Done():
			return
		}
	}
}

func transAccData(d *ws.AccountBalanceData) model.Account {
	a := model.NewAccount()
	a.Timestamp = d.Timestamp
	for _, data := range d.Data {
		amount := decimal.Require(data.Amount)
		available := decimal.Require(data.Available)
		a.Balances.Store(data.Currency, model.Balance{
			Available: available,
			InTrade:   amount.Sub(available),
		})
	}
	return a
}
