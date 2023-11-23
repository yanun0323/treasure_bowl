package mock

import (
	"context"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/robfig/cron"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"

	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"
)

type KlineProvider struct {
	connected       *atomic.Bool
	l               logs.Logger
	ch              chan model.Kline
	pair            model.Pair
	cronJob         *cron.Cron
	targetKlineType model.KlineType

	cacheKlineClose int
}

func NewKlineProvider(p model.Pair, target model.KlineType) domain.KlineProvideServer {
	return &KlineProvider{
		l:               logs.New("mock kline provider", util.LogLevel()),
		ch:              make(chan model.Kline, 10),
		targetKlineType: target,
		cacheKlineClose: 100,
		pair:            p,
	}
}

func (p *KlineProvider) Connect(ctx context.Context) (<-chan model.Kline, error) {
	p.connected.Store(true)
	if p.cronJob == nil {
		p.cronJob = cron.New()
		p.cronJob.AddFunc(p.targetKlineType.CronSpec(), func() {
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
	k := randomKline(p.pair, p.targetKlineType, p.cacheKlineClose)
	i, _ := strconv.Atoi(k.ClosePrice.String())
	p.cacheKlineClose = i
	p.ch <- k
}

func randomKline(p model.Pair, t model.KlineType, open int) model.Kline {
	closeThreshold := 30
	maxThreshold := 50
	minThreshold := 50
	volThreshold := 10_000

	close := open + rand.Intn(closeThreshold*2) - closeThreshold
	maxi := max(open, close) + rand.Intn(maxThreshold)
	mini := min(open, close) - rand.Intn(minThreshold)
	vol := rand.Intn(volThreshold)

	return model.Kline{
		Pair:       p,
		MaxPrice:   decimal.Require(strconv.Itoa(maxi)),
		MinPrice:   decimal.Require(strconv.Itoa(mini)),
		OpenPrice:  decimal.Require(strconv.Itoa(open)),
		ClosePrice: decimal.Require(strconv.Itoa(close)),
		Volume:     decimal.Require(strconv.Itoa(vol)),
		Type:       t,
		Timestamp:  time.Now().Unix(),
	}
}
