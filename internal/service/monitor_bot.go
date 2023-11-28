package service

import (
	"context"

	"main/internal/domain"
	"main/internal/model"

	"github.com/pkg/errors"
	"github.com/yanun0323/gollection/v2"
	"github.com/yanun0323/pkg/logs"
)

type inspectorBot struct {
	l      logs.Logger
	pair   model.Pair
	kps    domain.KlineProvideServer
	ch     <-chan model.Kline
	cancel context.CancelFunc
	cache  gollection.SyncPriorityQueue[*model.Kline]
	tree   gollection.SyncBTree[int64, *model.Kline]
	cap    int
}

func NewInspectorBot(ctx context.Context, pr model.Pair, kps domain.KlineProvideServer) (domain.StrategyBot, error) {
	return &inspectorBot{
		l:    logs.Get(ctx).WithField("server", "monitor bot"),
		pair: pr,
		kps:  kps,
		cap:  97,
	}, nil
}

func (bot *inspectorBot) Init(ctx context.Context) error {
	ch, err := bot.kps.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "inti bot")
	}
	bot.ch = ch
	bot.cache = gollection.NewSyncPriorityQueue[*model.Kline](func(k1, k2 *model.Kline) bool {
		return k1.Timestamp < k2.Timestamp
	})
	bot.tree = gollection.NewSyncBTree[int64, *model.Kline]()
	return nil
}

func (bot *inspectorBot) Run(ctx context.Context) error {
	if bot.ch == nil {
		return errors.New("require initializing before running")
	}

	c, cancel := context.WithCancel(ctx)
	go func(c context.Context) {
		for {
			select {
			case k := <-bot.ch:
				kk, ok := bot.tree.Search(k.Timestamp)
				if ok {
					if k.IsEqual(kk) {
						continue
					}

					if k.IsTypeEqual(kk) {
						k.Update(kk)
					} else {
						ok = false
					}
				}

				if !ok {
					bot.tree.Insert(k.Timestamp, &k)
					bot.cache.Enqueue(&k)
					if bot.cache.Len() >= bot.cap {
						dk := bot.cache.Dequeue()
						bot.l.Debugf("drop kline from cache: %+v", dk)
						_, ok = bot.tree.Remove(dk.Timestamp)
						bot.l.Debugf("remove: %+v, cache: %d, tree: %d", ok, bot.cache.Len(), bot.tree.Len())
					}
				}

				bot.l.Infof("%+v", k)
			case <-c.Done():
				bot.l.Info("stop consuming kline")
				return
			}
		}
	}(c)
	bot.cancel = cancel
	return nil
}

func (bot *inspectorBot) Shutdown(ctx context.Context) error {
	if err := bot.kps.Disconnect(ctx); err != nil {
		return errors.Wrap(err, "disconnect kline provider")
	}
	if bot.cancel != nil {
		bot.cancel()
		bot.cancel = nil
	}
	bot.ch = nil
	return nil
}
