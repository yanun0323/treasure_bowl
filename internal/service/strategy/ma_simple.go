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
	assetUpdating       *sync.RWMutex
	connecting          bool
	pair                model.Pair
	signal              chan model.Order
	asset               *model.Account
	kLineTree           map[model.KLineType]gollection.SyncBTree[uint64, model.KLine]
	supportOrderTypeMap util.SyncMap[model.OrderType, bool]
	orderMap            util.SyncMap[model.OrderType, []model.Order]
}

func NewMaSimple(pair model.Pair) (domain.StrategyServer, error) {
	return &MaSimple{
		assetUpdating:       &sync.RWMutex{},
		pair:                pair,
		signal:              make(chan model.Order, 10),
		asset:               model.NewAccount(),
		kLineTree:           map[model.KLineType]gollection.SyncBTree[uint64, model.KLine]{},
		supportOrderTypeMap: util.NewSyncMap[model.OrderType, bool](),
		orderMap:            util.NewSyncMap[model.OrderType, []model.Order](),
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
		if s.kLineTree[kLine.Type] == nil {
			s.kLineTree[kLine.Type] = gollection.NewSyncBTree[uint64, model.KLine]()
		}
		s.kLineTree[kLine.Type].Insert(kLine.Timestamp, kLine)
	}
	s.invokeStrategy(ctx)
}

func (s *MaSimple) PushAssets(ctx context.Context, accounts ...model.Account) {
	s.assetUpdating.Lock()
	defer s.assetUpdating.Unlock()
	for _, a := range accounts {
		s.asset = &a
	}
}

func (s *MaSimple) PushOrders(ctx context.Context, orders ...model.Order) {
	s.assetUpdating.Lock()
	defer s.assetUpdating.Unlock()
	for _, order := range orders {
		s.orderMap.LoadAndSet(order.Type, func(value []model.Order) []model.Order {
			return append(value, order)
		})
	}
}

func (s *MaSimple) PushSupportedOrderTypes(ctx context.Context, types ...model.OrderType) {
	s.supportOrderTypeMap.Clear()
	for _, t := range types {
		s.supportOrderTypeMap.Set(t, true)
	}
}

func (s *MaSimple) invokeStrategy(ctx context.Context) {
	if !s.connecting {
		return
	}
	s.assetUpdating.RLock()
	defer s.assetUpdating.RUnlock()

	// TODO: Implement me
}
