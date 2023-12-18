package app

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"main/internal/domain"
	"main/internal/entity"
	"main/internal/service"
	"main/internal/service/bitopro"
	"main/internal/service/mock"
	"main/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/logs"
)

const (
	/* ENV KEY */
	_keyPair                  = "PAIR"
	_keyStrategy              = "STG"
	_keyProviderKline         = "KLINE"
	_keyProviderTrade         = "TRADE"
	_keyProviderKlineDuration = "KLINE_DURATION"
	_keyBackTestingToggle     = "BT"
	_keyBackTestingStart      = "BT_START"
	_keyBackTestingEnd        = "BT_END"

	/* Strategy Value */
	_strategyMaStandup        = "standup"
	_strategyExchangeFollower = "follow"
	_strategyInspector        = "inspector"

	/* Provider Value */
	_providerMock    = "mock"
	_providerBitopro = "bitopro"
	_providerBinance = "binance"

	/* Setting */
	_timeLayout = "2006-01-02"
)

func Run() {
	l := logs.New(util.LogLevel())

	pr := strings.ToUpper(os.Getenv(_keyPair))
	if len(pr) == 0 {
		l.Fatalf("missing '%s' environment key", _keyPair)
	}

	spr := strings.Split(pr, "_")
	if len(spr) != 2 {
		l.Fatalf("unsupported pair %s, expected connecting with '/'. e.g. BTC/USDT", pr)
	}

	if len(spr[0]) == 0 || len(spr[1]) == 0 {
		l.Fatalf("empty pair base/quote %s", pr)
	}

	pair := entity.NewPair(spr[0], spr[1])

	stg := strings.ToLower(os.Getenv(_keyStrategy))
	if len(stg) == 0 {
		l.Fatalf("missing '%s' environment key", _keyStrategy)
	}

	kls := strings.ToLower(os.Getenv(_keyProviderKline))
	if len(kls) == 0 {
		l.Fatalf("missing '%s' environment key", _keyProviderKline)
	}

	tds := strings.ToLower(os.Getenv(_keyProviderTrade))
	if len(tds) == 0 {
		l.Fatalf("missing '%s' environment key", _keyProviderTrade)
	}

	dr := strings.ToLower(os.Getenv(_keyProviderKlineDuration))
	if len(dr) == 0 {
		l.Fatalf("missing '%s' environment key", _keyProviderKlineDuration)
	}

	kt := entity.KlineType(dr)
	if !kt.Validate() {
		l.Fatalf("unsupported kline duration '%s'", dr)
	}

	isBt := (strings.ToLower(os.Getenv(_keyBackTestingToggle)) == "true")
	var (
		btStart, btEnd time.Time
	)

	if isBt {
		btStartStr := os.Getenv(_keyBackTestingStart)
		bs, err := time.Parse(_timeLayout, btStartStr)
		if err != nil {
			l.WithError(err).Fatalf("back testing start time format should be: %s, but get: %s", _timeLayout, btStartStr)
		}

		btEndStr := os.Getenv(_keyBackTestingEnd)
		be, err := time.Parse(_timeLayout, btEndStr)
		if err != nil {
			l.WithError(err).Fatalf("back testing end time format should be: %s, but get: %s", _timeLayout, btEndStr)
		}

		btStart = bs
		btEnd = be
	}

	l.Info("pair: ", pr)
	l.Info("strategy: ", stg)
	l.Info("kline: ", kls)
	l.Info("trade: ", tds)

	ctx, l := l.Attach(context.Background())

	var (
		err   error
		kline []domain.KlineProvideServer
		trade []domain.TradeServer

		bot domain.StrategyBot
	)

	for _, kl := range strings.Split(kls, ",") {
		switch kl {
		case _providerMock:
			kline = append(kline, mock.NewKlineProvider(ctx, pair, kt))
		case _providerBitopro:
			k, err := bitopro.NewKlineProvider(ctx, pair, kt)
			if err != nil {
				l.WithError(err).Fatal("init bitopro kline provider")
			}
			kline = append(kline, k)
		case _providerBinance:
		default:
			l.Fatalf("unsupported kline provider '%s'", kl)
		}
	}

	for _, td := range strings.Split(tds, ",") {
		if isBt {
			td = _providerMock
		}

		switch td {
		case _providerMock:
			t, err := mock.NewTradeServer(ctx, pair, entity.OrderTypeLimit, entity.OrderTypeMarket)
			if err != nil {
				l.WithError(err).Fatal("init mock trade server")
			}
			trade = append(trade, t)
		case _providerBitopro:
			t, err := bitopro.NewTradeServer(ctx, pair)
			if err != nil {
				l.WithError(err).Fatal("init bitopro trade server")
			}
			trade = append(trade, t)
		case _providerBinance:
		default:
			l.Fatalf("unsupported trade server '%s'", td)
		}
	}

	switch stg {
	case _strategyMaStandup:
		if len(kline) == 0 {
			l.Fatal("require at least one kline provider")
		}

		if len(trade) == 0 {
			l.Fatal("require at least one trade provider")
		}

		pv := service.StdBotProvider{
			Kline: kline[0],
			Trade: trade[0],
		}
		if err := pv.Validate(); err != nil {
			l.WithError(err).Error("validate bot provider")
		}

		bot, err = service.NewMaStandupBot(ctx, pair, pv)
		if err != nil {
			l.WithError(err).Fatal("create ma standup bot")
		}
	case _strategyExchangeFollower:
		if len(kline) <= 1 {
			l.Fatal("require at least two kline provider")
		}

		if len(trade) == 0 {
			l.Fatal("require at least one trade provider")
		}

		pv := service.ExchangeFollowerProvider{
			Source: kline[0],
			Target: kline[1],
			Trade:  trade[0],
		}
		if err := pv.Validate(); err != nil {
			l.WithError(err).Error("validate bot provider")
		}

		bot, err = service.NewExchangeFollowerBot(ctx, pair, pv)
		if err != nil {
			l.WithError(err).Fatal("create exchange follower bot")
		}
	case _strategyInspector:
		if len(kline) == 0 {
			l.Fatal("require at least one kline provider")
		}
		bot, err = service.NewInspectorBot(ctx, pair, kline[0])
		if err != nil {
			l.WithError(err).Fatal("create inspector bot")
		}
	default:
		l.Fatalf("unsupported strategy '%s'", stg)
	}

	if bot == nil {
		l.Fatal("nil bot")
	}

	if isBt {
		if err := bot.BackTesting(ctx, btStart, btEnd); err != nil {
			l.WithError(err).Fatal("bot back testing")
		}
		return
	}

	if err := bot.Init(ctx); err != nil {
		l.WithError(err).Fatal("init bot")
	}

	go func() {
		if err := bot.Run(ctx); err != nil {
			l.WithError(err).Fatal("run bot")
		}
	}()

	go func() {
		l.Fatal(registerBotHttpRouter(bot))
	}()

	/* Graceful shutdown */
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	l.Info("shutdown server")
	if err := bot.Shutdown(ctx); err != nil {
		l.WithError(err).Fatal("shutdown server")
	}
}

func registerBotHttpRouter(bot domain.StrategyBot) error {
	e := echo.New()
	entry := e.Group("/strategy_bot", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))

	entry.GET("/info", bot.GetInfo)

	port := ":8080"
	if p := viper.GetString("server.port"); len(p) != 0 {
		port = p
	}

	if port[0] != ':' {
		port = ":" + port
	}

	return e.Start(port)
}
