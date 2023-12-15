package mock

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"main/internal/domain"
	"main/internal/entity"
	"main/internal/util"

	"github.com/robfig/cron"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

type klineProvider struct {
	l               logs.Logger
	ch              chan entity.Kline
	pair            entity.Pair
	targetKlineType entity.KlineType
	cronJob         *cron.Cron

	cacheKlineClose int
}

func NewKlineProvider(ctx context.Context, pair entity.Pair, target entity.KlineType) domain.KlineProvideServer {
	p := &klineProvider{
		l:               logs.Get(ctx).WithField("server", "mock kline provider"),
		ch:              make(chan entity.Kline, 10),
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

func (p *klineProvider) Connect(ctx context.Context, requiredKlineInitCount int) (<-chan entity.Kline, error) {
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

func randomKline(p entity.Pair, t entity.KlineType, open int) entity.Kline {
	closeThreshold := 30
	maxThreshold := 50
	minThreshold := 50
	volThreshold := 10_000

	close := open + rand.Intn(closeThreshold*2) - closeThreshold
	maxi := max(open, close) + rand.Intn(maxThreshold)
	mini := min(open, close) - rand.Intn(minThreshold)
	vol := rand.Intn(volThreshold)

	return entity.Kline{
		Pair:       p,
		HighPrice:  decimal.Require(strconv.Itoa(maxi)),
		LowPrice:   decimal.Require(strconv.Itoa(mini)),
		OpenPrice:  decimal.Require(strconv.Itoa(open)),
		ClosePrice: decimal.Require(strconv.Itoa(close)),
		Volume:     decimal.Require(strconv.Itoa(vol)),
		Type:       t,
		Timestamp:  time.Now().Unix(),
	}
}
