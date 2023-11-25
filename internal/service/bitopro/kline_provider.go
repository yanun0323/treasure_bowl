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

type KlineProvider struct {
	l               logs.Logger
	ch              chan model.Kline
	pair            model.Pair
	targetKlineType model.KlineType
	cronJob         *cron.Cron
	startAt         int64
}

func NewKlineProvider(pair model.Pair, target model.KlineType) (domain.KlineProvideServer, error) {
	if _, ok := _klineTypeTrans[target]; !ok {
		return nil, errors.New(fmt.Sprintf("unsupported kline type: %s", target.String()))
	}
	p := &KlineProvider{
		l:               logs.New("bitopro kline provider", util.LogLevel()),
		ch:              make(chan model.Kline, 50),
		pair:            pair,
		targetKlineType: target,
		cronJob:         cron.New(),
		startAt:         time.Now().Unix(),
	}

	p.cronJob.AddFunc(util.CronSpec(), func() {
		p.publishKline()
	})

	return p, nil
}

func (p *KlineProvider) Connect(ctx context.Context) (<-chan model.Kline, error) {
	p.cronJob.Start()
	return p.ch, nil
}

func (p *KlineProvider) Disconnect(ctx context.Context) error {
	p.cronJob.Stop()
	return nil
}

func (p *KlineProvider) publishKline() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, p.getApiPath(), nil)
	if err != nil {
		p.l.WithError(err).Error("new request")
		return
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		p.l.WithError(err).Error("send request")
		return
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		p.l.WithError(err).Error("read response body")
		return
	}

	ohlc := &OHLC{}
	if err := json.Unmarshal(buf, ohlc); err != nil {
		p.l.WithError(err).Error("parsing ohlc json data")
	}

	for _, data := range ohlc.Data {
		p.ch <- data.Kline(p.pair, p.targetKlineType)
	}
}

func (p *KlineProvider) getApiPath() string {
	to := time.Now()
	d := p.targetKlineType.Duration(to.Unix())
	from := to.Add(99 * d)

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
		OpenPrice:  d.Open,
		ClosePrice: d.Close,
		MaxPrice:   d.High,
		MinPrice:   d.Low,
		Volume:     d.Volume,
		Timestamp:  d.Timestamp / 1000,
	}
}
