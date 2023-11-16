package mock

import (
	"context"
	"fmt"
	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/yanun0323/gollection/v2"
	"github.com/yanun0323/pkg/logs"
)

type orderServer struct {
	l            logs.Logger
	accPublisher chan<- model.Account
	msg          chan model.Order

	connected        atomic.Bool
	supportedTypeSet gollection.Set[model.OrderType]

	account       model.Account
	currentOrders util.SyncMap[string, model.Order]
}

func OrderServer(accChan chan<- model.Account, acc model.Account, supportedOrderTypes ...model.OrderType) domain.OrderServer {
	set := gollection.NewSet[model.OrderType]()
	set.Insert(supportedOrderTypes...)
	return &orderServer{
		l:                logs.New("mock order server", util.LogLevel()),
		accPublisher:     accChan,
		msg:              make(chan model.Order, 10),
		supportedTypeSet: set,
		account:          acc,
		currentOrders:    util.NewSyncMap[string, model.Order](),
	}
}

func (s *orderServer) Connect(ctx context.Context) (<-chan model.Order, error) {
	s.connected.Store(true)
	return s.msg, nil
}

func (s *orderServer) DisConnect(ctx context.Context) error {
	s.connected.Store(false)
	return nil
}

func (s *orderServer) SupportOrderType(ctx context.Context) ([]model.OrderType, error) {
	if !s.connected.Load() {
		return nil, nil
	}

	return s.supportedTypeSet.ToSlice(), nil
}

func (s *orderServer) PostOrder(ctx context.Context, order model.Order) error {
	if !s.connected.Load() {
		return errors.New("order server is disconnect")
	}

	if !s.supportedTypeSet.Contain(order.Type) {
		return errors.New(fmt.Sprintf("unsupported type: %s", order.Type.String()))
	}

	switch order.Action {
	case model.BUY:
		if err := s.account.MoveToInTrade(order.Pair.Quote(), order.GetAmount()); err != nil {
			return errors.Wrapf(err, "buy %s", order.Pair.Quote())
		}
	case model.SELL:
		if err := s.account.MoveToInTrade(order.Pair.Base(), order.GetAmount()); err != nil {
			return errors.Wrapf(err, "sell %s", order.Pair.Base())
		}
	case model.BUY_CANCEL:
		if err := s.account.MoveToAvailable(order.Pair.Quote(), order.GetAmount()); err != nil {
			return errors.Wrapf(err, "buy %s", order.Pair.Quote())
		}
	case model.SELL_CANCEL:
		if err := s.account.MoveToAvailable(order.Pair.Base(), order.GetAmount()); err != nil {
			return errors.Wrapf(err, "sell %s", order.Pair.Base())
		}
	}

	s.accPublisher <- s.account

	return nil
}
