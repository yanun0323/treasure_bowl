// documentation
// auth: https://github.com/bitoex/bitopro-offical-api-docs/blob/master/README.md#api-security-protocol
// create: https://github.com/bitoex/bitopro-offical-api-docs/blob/master/api/v3/private/create_an_order.md
// cancel: https://github.com/bitoex/bitopro-offical-api-docs/blob/master/api/v3/private/cancel_an_order.md
// consumer: https://github.com/bitoex/bitopro-offical-api-docs/blob/master/ws/private/open_orders_stream.md
package bitopro

import (
	"context"
	"strings"
	"sync/atomic"

	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"

	"github.com/bitoex/bitopro-api-go/pkg/bitopro"
	"github.com/bitoex/bitopro-api-go/pkg/ws"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

type OrderServer struct {
	l         logs.Logger
	pair      model.Pair
	connected *atomic.Bool
	cancel    chan struct{}
	cancelFn  context.CancelFunc

	wss          *ws.Ws
	client       *bitopro.AuthAPI
	orderChannel chan model.Order
}

func NewOrderServer(pair model.Pair, wss *ws.Ws, client *bitopro.AuthAPI) (domain.OrderServer, error) {
	return &OrderServer{
		l:            logs.New("bitopro order server", util.LogLevel()),
		pair:         pair,
		wss:          wss,
		client:       client,
		orderChannel: make(chan model.Order, 100),
	}, nil
}

func (p *OrderServer) Connect(ctx context.Context) (<-chan model.Order, error) {
	p.connected.Store(true)
	ch, cancel := p.wss.RunOrdersWsConsumer(ctx)
	p.cancel = cancel

	c, cancelFn := context.WithCancel(ctx)
	p.cancelFn = cancelFn
	go p.consumeOrder(c, ch)
	return p.orderChannel, nil
}

func (p *OrderServer) DisConnect(ctx context.Context) error {
	connected := p.connected.Load()
	if !connected {
		return nil
	}
	p.connected.Store(false)
	go func() {
		defer close(p.cancel)
		p.cancel <- struct{}{}
	}()

	p.cancelFn()

	return nil
}

func (p *OrderServer) SupportOrderType(ctx context.Context) ([]model.OrderType, error) {
	return []model.OrderType{
		model.OrderTypeLimit,
		model.OrderTypeStopLimit,
		model.OrderTypeMarket,
	}, nil
}

func (p *OrderServer) PostOrder(ctx context.Context, order model.Order) error {
	// TODO: Implement me
	return nil
}

func (p *OrderServer) consumeOrder(ctx context.Context, ch <-chan ws.OrdersData) {
	for {
		select {
		case order := <-ch:
			p.l.Debugf("consume account balance: %+v", order.Data)
			for _, o := range convOrderData(p.pair, &order) {
				p.orderChannel <- o
			}
		case <-ctx.Done():
			return
		}
	}
}

func convOrderData(pair model.Pair, d *ws.OrdersData) []model.Order {
	ods := d.Data[strings.ToLower(pair.String("_"))]
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
		return model.OrderTypeStopLimit
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
