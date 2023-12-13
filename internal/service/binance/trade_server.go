// documentation
// go example: https://github.com/binance/binance-connector-go
// stream: https://github.com/binance/binance-spot-api-docs/blob/master/web-socket-api.md#account-requests
package binance

import (
	"context"
	"main/internal/domain"
	"main/internal/entity"
)

type tradeServer struct {
}

func NewTradeServer(ctx context.Context, pr entity.Pair) (domain.TradeServer, error) {
	return &tradeServer{}, nil
}

func (s *tradeServer) Connect(context.Context) (entity.Account, <-chan entity.Order, error) {
	return entity.NewAccount(map[string]entity.Balance{}), nil, nil
}

func (s *tradeServer) Connected() bool {
	return false
}

func (s *tradeServer) Disconnect(context.Context) error {
	return nil
}

func (s *tradeServer) IsSupported(entity.OrderType) bool {
	return false
}

func (s *tradeServer) PushOrder(context.Context, entity.Order) (entity.Account, error) {
	return entity.NewAccount(map[string]entity.Balance{}), nil
}
