package mock

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"main/internal/domain"
	"main/internal/entity"

	"github.com/pkg/errors"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

type tradeServer struct {
	l         logs.Logger
	pair      entity.Pair
	account   entity.Account
	order     chan entity.Order
	connected bool
	supported map[entity.OrderType]bool
}

func NewTradeServer(ctx context.Context, pr entity.Pair, supportedOrderTypes ...entity.OrderType) (domain.TradeServer, error) {
	if len(supportedOrderTypes) == 0 {
		return nil, errors.New("need at least 1 supported order type")
	}

	m := make(map[entity.OrderType]bool, len(supportedOrderTypes))
	for _, ot := range supportedOrderTypes {
		m[ot] = true
	}

	return &tradeServer{
		l:    logs.Get(ctx).WithField("server", "mock trade server"),
		pair: pr,
		account: entity.NewAccount(map[string]entity.Balance{
			"BTC":  {Available: decimal.Require("1_000_000_000"), Unavailable: "0"},
			"USDT": {Available: decimal.Require("1_000_000_000"), Unavailable: "0"},
		}),
		order:     make(chan entity.Order, 100),
		supported: m,
	}, nil
}

func (s *tradeServer) Connect(context.Context) (entity.Account, <-chan entity.Order, error) {
	s.connected = true
	return s.account, s.order, nil
}

func (s *tradeServer) Connected() bool {
	return s.connected
}

func (s *tradeServer) Disconnect(context.Context) error {
	s.connected = false
	return nil
}

func (s *tradeServer) IsSupported(o entity.OrderType) bool {
	return s.supported[o]
}

func (s *tradeServer) PushOrder(ctx context.Context, order entity.Order) (entity.Account, error) {
	switch order.Action {
	case entity.OrderActionBuy, entity.OrderActionSell, entity.OrderActionCancelBuy, entity.OrderActionCancelSell:
		order.ID = strconv.FormatInt(time.Now().UnixMilli(), 10)
		order.Status = entity.OrderStatusComplete
		order.Action = entity.OrderActionNone
		s.order <- order
		return s.account, nil
	default:
		return nil, errors.New(fmt.Sprintf("unsupported order action: %s", order.Action.String()))
	}
}
