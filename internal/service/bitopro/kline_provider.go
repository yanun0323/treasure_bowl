// documentation
// ohlc: https://github.com/bitoex/bitopro-offical-api-docs/blob/master/api/v3/public/get_ohlc_data.md
package bitopro

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

const (
	_klineUrl = "/trading-history/"
)

var (
	_klineTypeTrans = map[model.KlineType]string{
		model.K1m:  "1m",
		model.K5m:  "5m",
		model.K15m: "15m",
		model.K30m: "30m",
		model.K1h:  "1h",
		model.K4h:  "4h",
		model.K6h:  "6h",
		model.K12h: "12h",
		model.K1d:  "1d",
		model.K1w:  "1w",
		model.K1M:  "1M",
	}
)

type klineProvider struct {
	l               logs.Logger
	ch              chan model.Kline
	pair            model.Pair
	targetKlineType model.KlineType
	cronJob         *cron.Cron
	startAt         int64
}

func NewKlineProvider(ctx context.Context, pair model.Pair, target model.KlineType) (domain.KlineProvideServer, error) {
	if _, ok := _klineTypeTrans[target]; !ok {
		return nil, errors.New(fmt.Sprintf("unsupported kline type: %s", target.String()))
	}
	p := &klineProvider{
		l:               logs.Get(ctx).WithField("server", "bitopro kline"),
		pair:            pair,
		targetKlineType: target,
		cronJob:         cron.New(),
		startAt:         time.Now().Unix(),
	}

	p.cronJob.AddFunc(util.CronSpec(), func() {
		p.publishKline(ctx)
	})

	return p, nil
}

func (p *klineProvider) Connect(ctx context.Context, requiredKlineInitCount int) (<-chan model.Kline, error) {
	result, err := p.requestKline(context.Background(), requiredKlineInitCount)
	if err != nil {
		return nil, err
	}
	if len(p.ch) == 0 {
		p.ch = make(chan model.Kline, len(result)*2)
	}
	for i := range result {
		p.ch <- result[i]
	}

	p.cronJob.Start()
	return p.ch, nil
}

func (p *klineProvider) Disconnect(ctx context.Context) error {
	p.cronJob.Stop()
	return nil
}

func (p *klineProvider) requestKline(ctx context.Context, count int) ([]model.Kline, error) {
	c, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	r, err := http.NewRequestWithContext(c, http.MethodGet, p.getApiPath(count), nil)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	ohlc := &OHLC{}
	if err := json.Unmarshal(buf, ohlc); err != nil {
		errors.Wrap(err, "parsing ohlc json data")
	}

	result := make([]model.Kline, 0, len(ohlc.Data))
	for _, data := range ohlc.Data {
		result = append(result, data.Kline(p.pair, p.targetKlineType))
	}
	return result, nil
}

func (p *klineProvider) publishKline(ctx context.Context) {
	result, err := p.requestKline(ctx, 5)
	if err != nil {
		p.l.WithError(err).Error("request kline")
		return
	}

	for i := range result {
		p.ch <- result[i]
	}
}

func (p *klineProvider) getApiPath(klineCount int) string {
	to := time.Now()
	d := p.targetKlineType.Duration(to.Unix())
	from := to.Add(time.Duration(klineCount) * d)

	return fmt.Sprintf("%s%s%s?resolution=%s&from=%d&to=%d",
		_restfulHost,
		_klineUrl,
		p.pair.Uppercase("_"),
		_klineTypeTrans[p.targetKlineType],
		from.Unix(),
		to.Unix())
}

type OHLC struct {
	Data []OHLCKlineData `json:"data"`
}

type OHLCKlineData struct {
	Timestamp int64           `json:"timestamp"`
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	Volume    decimal.Decimal `json:"volume"`
}

func (d *OHLCKlineData) Kline(p model.Pair, t model.KlineType) model.Kline {
	return model.Kline{
		Pair:       p,
		Type:       t,
		Source:     model.KlineSourceBitoPro,
		OpenPrice:  d.Open,
		ClosePrice: d.Close,
		MaxPrice:   d.High,
		MinPrice:   d.Low,
		Volume:     d.Volume,
		Timestamp:  d.Timestamp / 1000,
	}
}
