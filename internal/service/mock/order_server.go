package mock

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/yanun0323/gollection/v2"
	"github.com/yanun0323/pkg/logs"

	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"
)

type orderServer struct {
	l            logs.Logger
	accPublisher chan<- model.Account
	msg          chan model.Order

	connected        *atomic.Bool
	supportedTypeSet gollection.Set[model.OrderType]

	account       model.Account
	currentOrders util.SyncMap[string, model.Order]
}

func NewOrderServer(accChan chan<- model.Account, acc model.Account, supportedOrderTypes ...model.OrderType) domain.OrderServer {
	set := gollection.NewSet[model.OrderType]()
	set.Insert(supportedOrderTypes...)
	return &orderServer{
		l:                logs.New("mock order server", util.LogLevel()),
		accPublisher:     accChan,
		msg:              make(chan model.Order, 10),
		supportedTypeSet: set,
		account:          acc,
		currentOrders:    util.NewSyncMap[string, model.Order](),
		connected:        &atomic.Bool{},
	}
}

func (p *orderServer) Connect(ctx context.Context) (<-chan model.Order, error) {
	p.connected.Store(true)
	return p.msg, nil
}

func (p *orderServer) DisConnect(ctx context.Context) error {
	p.connected.Store(false)
	return nil
}

func (p *orderServer) SupportOrderType(ctx context.Context) ([]model.OrderType, error) {
	if !p.connected.Load() {
		return nil, nil
	}

	return p.supportedTypeSet.ToSlice(), nil
}

func (p *orderServer) PostOrder(ctx context.Context, order model.Order) error {
	if !p.connected.Load() {
		return errors.New("order server is disconnect")
	}

	if !p.supportedTypeSet.Contain(order.Type) {
		return errors.New(fmt.Sprintf("unsupported type: %s", order.Type.String()))
	}

	switch order.Action {
	case model.OrderActionBuy:
		if err := p.account.MoveToInTrade(order.Pair.Quote(), order.GetTotal()); err != nil {
			return errors.Wrapf(err, "buy %s", order.Pair.Quote())
		}
	case model.OrderActionSell:
		if err := p.account.MoveToInTrade(order.Pair.Base(), order.GetTotal()); err != nil {
			return errors.Wrapf(err, "sell %s", order.Pair.Base())
		}
	case model.OrderActionCancelBuy:
		if err := p.account.MoveToAvailable(order.Pair.Quote(), order.GetTotal()); err != nil {
			return errors.Wrapf(err, "buy %s", order.Pair.Quote())
		}
	case model.OrderActionCancelSell:
		if err := p.account.MoveToAvailable(order.Pair.Base(), order.GetTotal()); err != nil {
			return errors.Wrapf(err, "sell %s", order.Pair.Base())
		}
	}

	p.accPublisher <- p.account

	return nil
}
