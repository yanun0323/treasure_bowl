package bitopro

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"main/internal/domain"
	"main/internal/entity"

	"github.com/bitoex/bitopro-api-go/pkg/bitopro"
	"github.com/bitoex/bitopro-api-go/pkg/ws"
	"github.com/pkg/errors"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

var (
	_supported = map[entity.OrderType]bool{
		entity.OrderTypeLimit:  true,
		entity.OrderTypeMarket: true,
		// TODO: OrderTypeStopLimit official package didn't support 'stop limit' this type, need rewrite
		entity.OrderTypeStopLimit: false,
	}
)

type tradeServer struct {
	l         logs.Logger
	pair      entity.Pair
	connected bool

	wss          *ws.Ws
	client       *bitopro.AuthAPI
	clientID     int
	orderChannel chan entity.Order

	cancel   chan struct{}
	cancelFn context.CancelFunc
}

func NewTradeServer(ctx context.Context, pr entity.Pair) (domain.TradeServer, error) {
	return &tradeServer{
		l:    logs.Get(ctx).WithField("server", "bitopro trade server"),
		pair: pr,
	}, nil
}

func (s *tradeServer) Connect(ctx context.Context) (entity.Account, <-chan entity.Order, error) {
	defer func() { s.connected = true }()
	ch, cancel := s.wss.RunOrdersWsConsumer(ctx)
	s.cancel = cancel

	acc, err := s.getAccount()
	if err != nil {
		return nil, nil, errors.Wrap(err, "get account")
	}

	s.orderChannel = make(chan entity.Order, len(ch)*2)

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

func (s *tradeServer) IsSupported(t entity.OrderType) bool {
	return _supported[t]
}

func (s *tradeServer) PushOrder(ctx context.Context, o entity.Order) (entity.Account, error) {
	if err := o.ValidatePushingOrder(); err != nil {
		return nil, errors.Wrap(err, "validate pushing order")
	}

	switch o.Action {
	case entity.OrderActionBuy:
		return s.createOrderBuy(ctx, o)
	case entity.OrderActionSell:
		return s.createOrderSell(ctx, o)
	case entity.OrderActionCancelBuy, entity.OrderActionCancelSell:
		return s.cancelOrder(ctx, o)
	default:
		return nil, errors.New("unsupported order action")
	}
}

func (s *tradeServer) getAccount() (entity.Account, error) {
	a := s.client.GetAccountBalance()
	acc, err := convAccount(a)
	if err != nil {
		return nil, errors.Wrap(err, "convert account")
	}
	return acc, nil
}

func (s *tradeServer) createOrderBuy(ctx context.Context, o entity.Order) (entity.Account, error) {
	var c *bitopro.CreateOrder

	switch o.Type {
	case entity.OrderTypeLimit:
		c = s.client.CreateOrderLimitBuy(s.clientID, s.pair.Lowercase("_"), o.Price.String(), o.Amount.Total.String())
	case entity.OrderTypeMarket:
		c = s.client.CreateOrderMarketBuy(s.clientID, s.pair.Lowercase("_"), o.Amount.Total.String())
	case entity.OrderTypeStopLimit:
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

func (s *tradeServer) createOrderSell(ctx context.Context, o entity.Order) (entity.Account, error) {
	var c *bitopro.CreateOrder

	switch o.Type {
	case entity.OrderTypeLimit:
		c = s.client.CreateOrderLimitSell(s.clientID, s.pair.Lowercase("_"), o.Price.String(), o.Amount.Total.String())
	case entity.OrderTypeMarket:
		c = s.client.CreateOrderMarketSell(s.clientID, s.pair.Lowercase("_"), o.Amount.Total.String())
	case entity.OrderTypeStopLimit:
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

func (s *tradeServer) cancelOrder(ctx context.Context, o entity.Order) (entity.Account, error) {
	switch o.Type {
	case entity.OrderTypeLimit, entity.OrderTypeMarket:
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
	case entity.OrderTypeStopLimit:
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

func convOrderData(pair entity.Pair, d *ws.OrdersData) []entity.Order {
	ods := d.Data[strings.ToLower(pair.Uppercase("_"))]
	if len(ods) == 0 {
		return nil
	}

	result := make([]entity.Order, 0, len(ods))
	for _, o := range ods {
		ot := convOrderType(o.Type)
		if ot.IsUnknown() {
			continue
		}
		result = append(result, entity.Order{
			ID:        o.ID,
			Pair:      pair,
			Action:    entity.OrderActionNone,
			Type:      ot,
			Status:    convStatus(o.Status),
			Timestamp: o.Timestamp / 1000,
			Price:     decimal.Require(o.Price),
			Amount: entity.Amount{
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
func convOrderType(s string) entity.OrderType {
	switch strings.ToUpper(s) {
	case "LIMIT":
		return entity.OrderTypeLimit
	case "MARKET":
		return entity.OrderTypeMarket
	case "STOP_LIMIT":
		// TODO: official package didn't support 'stop limit' this type, need rewrite
		return entity.OrderTypeUnknown
	default:
		return entity.OrderTypeUnknown
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
func convStatus(s int) entity.OrderStatus {
	switch s {
	case 0, 1:
		return entity.OrderStatusCreated
	case 2:
		return entity.OrderStatusComplete
	case 3:
		return entity.OrderStatusPartialComplete
	case 4:
		return entity.OrderStatusCanceled
	default:
		return entity.OrderStatusPending
	}
}
