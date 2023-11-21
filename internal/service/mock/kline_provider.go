package mock

import (
	"context"
	"sync/atomic"

	"github.com/robfig/cron"
	"github.com/yanun0323/pkg/logs"

	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"
)

type KlineProvider struct {
	connected *atomic.Bool
	l         logs.Logger
	ch        chan model.Kline
	cronJob   *cron.Cron
}

func NewKlineProvider() domain.KlineProvideServer {
	return &KlineProvider{
		l:  logs.New("mock kline provider", util.LogLevel()),
		ch: make(chan model.Kline, 10),
	}
}

func (p *KlineProvider) Connect(ctx context.Context, types ...model.KlineType) (<-chan model.Kline, error) {
	p.connected.Store(true)
	if p.cronJob == nil {
		p.cronJob = cron.New()
		p.cronJob.AddFunc("*/1 * * * * *", func() {
			p.publishKline()
		})
	}
	p.cronJob.Stop()
	p.cronJob.Run()
	return p.ch, nil
}

func (p *KlineProvider) Disconnect(ctx context.Context) error {
	p.connected.Store(false)
	p.cronJob.Stop()
	return nil
}

func (p *KlineProvider) publishKline() {

}
