package strategy

import (
	"context"
	"sync"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/yanun0323/gollection/v2"
	"github.com/yanun0323/pkg/logs"

	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"
)

type MaSimple struct {
	l                   logs.Logger
	updating            *sync.RWMutex
	pair                model.Pair
	signal              chan model.Order
	asset               model.Account
	klineTree           map[model.KlineType]gollection.SyncBTree[uint64, model.Kline]
	supportOrderTypeMap util.SyncMap[model.OrderType, bool]
	orderMap            util.SyncMap[model.OrderType, []model.Order]

	cronJob *cron.Cron
}

func NewMaSimple(pair model.Pair) (domain.StrategyServer, error) {
	return &MaSimple{
		l:                   logs.New("strategy ma simple", viper.GetUint16("log.levels")),
		updating:            &sync.RWMutex{},
		pair:                pair,
		signal:              make(chan model.Order, 10),
		asset:               model.NewAccount(),
		klineTree:           map[model.KlineType]gollection.SyncBTree[uint64, model.Kline]{},
		supportOrderTypeMap: util.NewSyncMap[model.OrderType, bool](),
		orderMap:            util.NewSyncMap[model.OrderType, []model.Order](),
	}, nil
}

func (s *MaSimple) Connect(ctx context.Context) (<-chan model.Order, error) {
	if s.cronJob == nil {
		s.cronJob = cron.New()
		s.cronJob.AddFunc("*/15 * * * * *", func() {
			s.invokeStrategy(ctx)
		})
	}
	s.cronJob.Stop()
	s.cronJob.Run()
	return s.signal, nil
}

func (s *MaSimple) Disconnect(ctx context.Context) error {
	if s.cronJob != nil {
		s.cronJob.Stop()
	}
	return nil
}

func (s *MaSimple) PushKlines(ctx context.Context, klines ...model.Kline) {
	for _, kline := range klines {
		if s.klineTree[kline.Type] == nil {
			s.klineTree[kline.Type] = gollection.NewSyncBTree[uint64, model.Kline]()
		}
		s.klineTree[kline.Type].Insert(kline.Timestamp, kline)
	}
}

func (s *MaSimple) PushAssets(ctx context.Context, accounts ...model.Account) {
	s.updating.Lock()
	defer s.updating.Unlock()
	for _, a := range accounts {
		s.asset = a
	}
}

func (s *MaSimple) PushOrders(ctx context.Context, orders ...model.Order) {
	s.updating.Lock()
	defer s.updating.Unlock()
	for _, order := range orders {
		s.orderMap.LoadAndSet(order.Type, func(value []model.Order) []model.Order {
			return append(value, order)
		})
	}
}

func (s *MaSimple) PushSupportedOrderTypes(ctx context.Context, types ...model.OrderType) {
	s.supportOrderTypeMap.Clear()
	for _, t := range types {
		s.supportOrderTypeMap.Store(t, true)
	}
}

func (s *MaSimple) invokeStrategy(ctx context.Context) {
	if locked := s.updating.TryRLock(); !locked {
		s.l.WithTime(time.Now()).Info("updating, skipped invoking strategy")
		return
	}
	defer s.updating.RUnlock()

	s.l.WithTime(time.Now()).Info("invoked strategy")

	// TODO: Implement me
}
