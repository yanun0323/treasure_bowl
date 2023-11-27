package bitopro

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"main/internal/domain"
	"main/internal/model"

	"github.com/bitoex/bitopro-api-go/pkg/bitopro"
	"github.com/bitoex/bitopro-api-go/pkg/ws"
	"github.com/pkg/errors"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

var (
	_supported = map[model.OrderType]bool{
		model.OrderTypeLimit:  true,
		model.OrderTypeMarket: true,
		// TODO: OrderTypeStopLimit official package didn't support 'stop limit' this type, need rewrite
		model.OrderTypeStopLimit: false,
	}
)

type tradeServer struct {
	l         logs.Logger
	pair      model.Pair
	connected bool

	wss          *ws.Ws
	client       *bitopro.AuthAPI
	clientID     int
	orderChannel chan model.Order

	cancel   chan struct{}
	cancelFn context.CancelFunc
}

func NewTradeServer(ctx context.Context, pr model.Pair) (domain.TradeServer, error) {
	return &tradeServer{
		l:    logs.Get(ctx).WithField("server", "bitopro trade server"),
		pair: pr,
	}, nil
}

func (s *tradeServer) Connect(ctx context.Context) (model.Account, <-chan model.Order, error) {
	defer func() { s.connected = true }()
	ch, cancel := s.wss.RunOrdersWsConsumer(ctx)
	s.cancel = cancel

	acc, err := s.getAccount()
	if err != nil {
		return nil, nil, errors.Wrap(err, "get account")
	}

	s.orderChannel = make(chan model.Order, len(ch)*2)

	c, cancelFn := context.WithCancel(ctx)
	s.cancelFn = cancelFn
	go s.consumeOrder(c, ch)

	return acc, s.orderChannel, nil
}

func (s *tradeServer) Connected() bool {
	return s.connected
}

func (s *tradeServer) Disconnect(context.Context) error {
	defer func() { s.connected = false }()
	if !s.connected {
		return nil
	}
	go func() {
		defer close(s.cancel)
		s.cancel <- struct{}{}
	}()

	s.cancelFn()
	return nil
}

func (s *tradeServer) IsSupported(t model.OrderType) bool {
	return _supported[t]
}

func (s *tradeServer) PushOrder(ctx context.Context, o model.Order) (model.Account, error) {
	if err := o.ValidatePushingOrder(); err != nil {
		return nil, errors.Wrap(err, "validate pushing order")
	}

	switch o.Action {
	case model.OrderActionBuy:
		return s.createOrderBuy(ctx, o)
	case model.OrderActionSell:
		return s.createOrderSell(ctx, o)
	case model.OrderActionCancelBuy, model.OrderActionCancelSell:
		return s.cancelOrder(ctx, o)
	default:
		return nil, errors.New("unsupported order action")
	}
}

func (s *tradeServer) getAccount() (model.Account, error) {
	a := s.client.GetAccountBalance()
	acc, err := convAccount(a)
	if err != nil {
		return nil, errors.Wrap(err, "convert account")
	}
	return acc, nil
}

func (s *tradeServer) createOrderBuy(ctx context.Context, o model.Order) (model.Account, error) {
	var c *bitopro.CreateOrder

	switch o.Type {
	case model.OrderTypeLimit:
		c = s.client.CreateOrderLimitBuy(s.clientID, s.pair.Lowercase("_"), o.Price.String(), o.Amount.Total.String())
	case model.OrderTypeMarket:
		c = s.client.CreateOrderMarketBuy(s.clientID, s.pair.Lowercase("_"), o.Amount.Total.String())
	case model.OrderTypeStopLimit:
		// TODO: official package didn't support 'stop limit' this type, need rewrite
		return nil, errors.New("stop limit doesn't support yet")
	default:
		return nil, errors.New("unsupported order type")
	}

	if c == nil {
		return nil, errors.New("pushing order connection error")
	}

	if len(c.Error) != 0 {
		return nil, errors.New(fmt.Sprintf("push order err: %s", c.Error))
	}

	if len(c.OrderID) == 0 {
		return nil, errors.New("empty response order ID")
	}

	acc, err := s.getAccount()
	if err != nil {
		return nil, errors.Wrap(err, "get account")
	}

	return acc, nil
}

func (s *tradeServer) createOrderSell(ctx context.Context, o model.Order) (model.Account, error) {
	var c *bitopro.CreateOrder

	switch o.Type {
	case model.OrderTypeLimit:
		c = s.client.CreateOrderLimitSell(s.clientID, s.pair.Lowercase("_"), o.Price.String(), o.Amount.Total.String())
	case model.OrderTypeMarket:
		c = s.client.CreateOrderMarketSell(s.clientID, s.pair.Lowercase("_"), o.Amount.Total.String())
	case model.OrderTypeStopLimit:
		// TODO: official package didn't support 'stop limit' this type, need rewrite
		return nil, errors.New("stop limit doesn't support yet")
	default:
		return nil, errors.New("unsupported order type")
	}

	if c == nil {
		return nil, errors.New("pushing order connection error")
	}

	if len(c.Error) != 0 {
		return nil, errors.New(fmt.Sprintf("push order err: %s", c.Error))
	}

	if len(c.OrderID) == 0 {
		return nil, errors.New("empty response order ID")
	}

	acc, err := s.getAccount()
	if err != nil {
		return nil, errors.Wrap(err, "get account")
	}

	return acc, nil
}

func (s *tradeServer) cancelOrder(ctx context.Context, o model.Order) (model.Account, error) {
	switch o.Type {
	case model.OrderTypeLimit, model.OrderTypeMarket:
		orderID, err := strconv.Atoi(o.ID)
		if err != nil {
			return nil, errors.Wrap(err, "convert order Id")
		}

		if orderID <= 0 {
			return nil, errors.New(fmt.Sprintf("invalid order ID (%d)", orderID))
		}

		c := s.client.CancelOrder(s.pair.Lowercase("_"), orderID)
		if c == nil {
			return nil, errors.New("pushing order connection error")
		}

		if len(c.Error) != 0 {
			return nil, errors.New(fmt.Sprintf("push order err: %s", c.Error))
		}

		if len(c.OrderID) == 0 {
			return nil, errors.New("empty response order ID")
		}
	case model.OrderTypeStopLimit:
		// TODO: official package didn't support 'stop limit' this type, need rewrite
		return nil, errors.New("stop limit doesn't support yet")
	default:
		return nil, errors.New("unsupported order type")
	}

	acc, err := s.getAccount()
	if err != nil {
		return nil, errors.Wrap(err, "get account")
	}

	return acc, nil
}

func (s *tradeServer) consumeOrder(ctx context.Context, ch <-chan ws.OrdersData) {
	for {
		select {
		case order := <-ch:
			s.l.Debugf("consume account balance: %+v", order.Data)
			for _, o := range convOrderData(s.pair, &order) {
				s.orderChannel <- o
			}
		case <-ctx.Done():
			return
		}
	}
}

func convOrderData(pair model.Pair, d *ws.OrdersData) []model.Order {
	ods := d.Data[strings.ToLower(pair.Uppercase("_"))]
	if len(ods) == 0 {
		return nil
	}

	result := make([]model.Order, 0, len(ods))
	for _, o := range ods {
		ot := convOrderType(o.Type)
		if ot.IsUnknown() {
			continue
		}
		result = append(result, model.Order{
			ID:        o.ID,
			Pair:      pair,
			Action:    model.OrderActionNone,
			Type:      ot,
			Status:    convStatus(o.Status),
			Timestamp: o.Timestamp / 1000,
			Price:     decimal.Require(o.Price),
			Amount: model.Amount{
				Total:  decimal.Require(o.OriginalAmount),
				Deal:   decimal.Require(o.ExecutedAmount),
				Remain: decimal.Require(o.RemainingAmount),
			},
		})

	}
	return result
}

// convOrderType converts api order type into model order type.
//
// LIMIT, Market or STOP_LIMIT
func convOrderType(s string) model.OrderType {
	switch strings.ToUpper(s) {
	case "LIMIT":
		return model.OrderTypeLimit
	case "MARKET":
		return model.OrderTypeMarket
	case "STOP_LIMIT":
		// TODO: official package didn't support 'stop limit' this type, need rewrite
		return model.OrderTypeUnknown
	default:
		return model.OrderTypeUnknown
	}
}

/*
convStatus converts api order status into model order status.

	// -1: Not Triggered
	// 0: In progress
	// 1: In progress (Partial deal)
	// 2: Completed
	// 3: Completed (Partial deal)
	// 4: Cancelled
	// 6: Post-only Cancelled
*/
func convStatus(s int) model.OrderStatus {
	switch s {
	case 0, 1:
		return model.OrderStatusCreated
	case 2:
		return model.OrderStatusComplete
	case 3:
		return model.OrderStatusPartialComplete
	case 4:
		return model.OrderStatusCanceled
	default:
		return model.OrderStatusPending
	}
}
