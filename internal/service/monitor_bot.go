package service

import (
	"context"

	"main/internal/domain"
	"main/internal/entity"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/yanun0323/gollection/v2"
	"github.com/yanun0323/pkg/logs"
)

type inspectorBot struct {
	l      logs.Logger
	pair   entity.Pair
	kps    domain.KlineProvideServer
	ch     <-chan entity.Kline
	cancel context.CancelFunc
	tree   gollection.SyncBTree[int64, *entity.Kline]
	cap    int
}

func NewInspectorBot(ctx context.Context, pr entity.Pair, kps domain.KlineProvideServer) (domain.StrategyBot, error) {
	return &inspectorBot{
		l:    logs.Get(ctx).WithField("server", "monitor bot"),
		pair: pr,
		kps:  kps,
		cap:  200,
	}, nil
}

func (bot *inspectorBot) Init(ctx context.Context) error {
	ch, err := bot.kps.Connect(ctx, bot.cap)
	if err != nil {
		return errors.Wrap(err, "inti bot")
	}
	bot.ch = ch
	bot.tree = gollection.NewSyncBTree[int64, *entity.Kline]()
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
					if !k.IsEqual(kk) {
						bot.l.Warnf("consume: %s", k)
						bot.l.Warnf("updated: %s\n", kk)
						kk.Update(&k)
					}
					continue
				}

				if bot.tree.Len() < bot.cap {
					bot.tree.Insert(k.Timestamp, &k)
					bot.l.Info(k.String())
					continue

				}

				minTS, _, _ := bot.tree.Min()
				if k.Timestamp < minTS {
					continue
				}

				_, _, _ = bot.tree.RemoveMin()
				bot.tree.Insert(k.Timestamp, &k)

				bot.l.Infof("consume: %s", k)
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

func (bot *inspectorBot) GetInfo(c echo.Context) error {
	return nil
}
