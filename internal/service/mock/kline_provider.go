package mock

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"main/internal/domain"
	"main/internal/entity"
	"main/internal/util"

	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

type klineProvider struct {
	l              logs.Logger
	ch             chan entity.Kline
	pair           entity.Pair
	kts            []entity.KlineType
	cronJobManager util.CronJobManager

	cacheKlineClose map[entity.KlineType]int
}

func NewKlineProvider(ctx context.Context, pair entity.Pair, targets ...entity.KlineType) domain.KlineProvideServer {
	p := &klineProvider{
		l:               logs.Get(ctx).WithField("server", "mock kline provider"),
		pair:            pair,
		kts:             targets,
		cronJobManager:  util.NewCronJobManager(len(targets)),
		cacheKlineClose: make(map[entity.KlineType]int, len(targets)),
	}
	for _, kt := range targets {
		p.cronJobManager.AddFunc(kt.Spec(), func() {
			p.publishKline(kt)
			p.cacheKlineClose[kt] = 100
		})
	}
	return p
}

func (p *klineProvider) Connect(ctx context.Context, spans ...entity.Span) (<-chan entity.Kline, error) {
	keep := false
	if len(spans) == 0 {
		keep = true
	}

	result := make([]entity.Kline, 0, 100)
	for _, kt := range p.kts {
		interval := kt.Duration()
		for _, span := range spans {
			count := int(span.End.Sub(span.Start) / interval)
			start := span.Start.Unix()
			for i := 0; i < count; i++ {
				k := randomKline(p.pair, kt, p.cacheKlineClose[kt])
				i, _ := strconv.Atoi(k.ClosePrice.String())
				p.cacheKlineClose[kt] = i

				k.Timestamp = start
				start += int64(interval.Seconds())
				result = append(result, k)
			}
		}
	}

	if len(p.ch) == 0 {
		p.ch = make(chan entity.Kline, len(result)*2)
	}

	for i := range result {
		p.ch <- result[i]
	}

	if keep {
		p.cronJobManager.Start()
	}

	return p.ch, nil
}

func (p *klineProvider) Disconnect(ctx context.Context) error {
	p.cronJobManager.Stop()
	return nil
}

func (p *klineProvider) publishKline(kt entity.KlineType) {
	k := randomKline(p.pair, kt, p.cacheKlineClose[kt])
	i, _ := strconv.Atoi(k.ClosePrice.String())
	p.cacheKlineClose[kt] = i
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
