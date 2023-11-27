package mock

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"

	"github.com/robfig/cron"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

type klineProvider struct {
	l               logs.Logger
	ch              chan model.Kline
	pair            model.Pair
	targetKlineType model.KlineType
	cronJob         *cron.Cron

	cacheKlineClose int
}

func NewKlineProvider(ctx context.Context, pair model.Pair, target model.KlineType) domain.KlineProvideServer {
	p := &klineProvider{
		l:               logs.Get(ctx).WithField("server", "mock kline provider"),
		ch:              make(chan model.Kline, 10),
		pair:            pair,
		targetKlineType: target,
		cacheKlineClose: 100,
		cronJob:         cron.New(),
	}
	p.cronJob.AddFunc(util.CronSpec(), func() {
		p.publishKline()
	})
	return p
}

func (p *klineProvider) Connect(ctx context.Context) (<-chan model.Kline, error) {
	p.cronJob.Run()
	return p.ch, nil
}

func (p *klineProvider) Disconnect(ctx context.Context) error {
	p.cronJob.Stop()
	return nil
}

func (p *klineProvider) publishKline() {
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
