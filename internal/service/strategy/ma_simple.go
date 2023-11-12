package strategy

import (
	"context"
	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"

	"github.com/yanun0323/gollection/v2"
)

type MaSimple struct {
	connecting bool
	pair       string
	signal     chan model.Order
	asset      *model.Account
	kLineTree  gollection.BTree[model.KLine]
	orderMap   util.SyncMap[model.OrderType, []model.Order]
}

func NewMaSimple(pair string) (domain.StrategyServer, error) {
	return &MaSimple{
		pair:      pair,
		signal:    make(chan model.Order, 10),
		asset:     model.NewAccount(),
		kLineTree: gollection.NewBTree[model.KLine](func(t1, t2 model.KLine) bool { return t1.Timestamp < t2.Timestamp }),
		orderMap:  util.NewSyncMap[model.OrderType, []model.Order](),
	}, nil
}

func (s *MaSimple) Connect(ctx context.Context) (<-chan model.Order, error) {
	s.connecting = true
	return s.signal, nil
}

func (s *MaSimple) Disconnect(ctx context.Context) error {
	s.connecting = false
	return nil
}

func (s *MaSimple) PushKLines(ctx context.Context, kLines ...model.KLine) {
	for _, kLine := range kLines {
		s.kLineTree.Insert(kLine)
	}
}

func (s *MaSimple) PushAssets(ctx context.Context, accounts ...model.Account) {
	for _, a := range accounts {
		s.asset = &a
	}
}

func (s *MaSimple) PushOrders(ctx context.Context, orders ...model.Order) {
	for _, order := range orders {
		s.orderMap.LoadAndSet(order.Type, func(value []model.Order) []model.Order {
			return append(value, order)
		})
	}
}

func (s *MaSimple) PushSupportedOrderTypes(ctx context.Context, types ...model.OrderType) {

}

func (s *MaSimple) invoke() {

}
