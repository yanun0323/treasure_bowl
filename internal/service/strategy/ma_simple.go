package strategy

import (
	"context"
	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"
	"sync"

	"github.com/yanun0323/gollection/v2"
)

type MaSimple struct {
	stopLock   *sync.RWMutex
	connecting bool
	pair       string
	signal     chan model.Order
	asset      *model.Account
	kLineTree  map[model.KLineType]gollection.SyncBTree[model.KLine]
	orderMap   util.SyncMap[model.OrderType, []model.Order]
}

func NewMaSimple(pair string) (domain.StrategyServer, error) {
	return &MaSimple{
		stopLock:  &sync.RWMutex{},
		pair:      pair,
		signal:    make(chan model.Order, 10),
		asset:     model.NewAccount(),
		kLineTree: map[model.KLineType]gollection.SyncBTree[model.KLine]{},
		orderMap:  util.NewSyncMap[model.OrderType, []model.Order](),
	}, nil
}

func newKLineTree() gollection.BTree[model.KLine] {
	return gollection.NewSyncBTree[model.KLine](func(t1, t2 model.KLine) bool { return t1.Timestamp < t2.Timestamp })
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
	s.stopLock.RLock()
	for _, kLine := range kLines {
		if s.kLineTree[kLine.Type] == nil {
			s.kLineTree[kLine.Type] = newKLineTree()
		}
		s.kLineTree[kLine.Type].Insert(kLine)
	}
	s.stopLock.RUnlock()
	s.triggerCalculation()
}

func (s *MaSimple) PushAssets(ctx context.Context, accounts ...model.Account) {
	s.stopLock.Lock()
	for _, a := range accounts {
		s.asset = &a
	}
	s.stopLock.Unlock()
	s.triggerCalculation()
}

func (s *MaSimple) PushOrders(ctx context.Context, orders ...model.Order) {
	s.stopLock.Lock()
	for _, order := range orders {
		s.orderMap.LoadAndSet(order.Type, func(value []model.Order) []model.Order {
			return append(value, order)
		})
	}
	s.stopLock.Unlock()
	s.triggerCalculation()
}

func (s *MaSimple) PushSupportedOrderTypes(ctx context.Context, types ...model.OrderType) {

}

func (s *MaSimple) triggerCalculation() {

}
