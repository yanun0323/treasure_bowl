package app

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"main/internal/domain"
	"main/internal/model"
	"main/internal/service"
	"main/internal/service/bitopro"
	"main/internal/service/mock"
	"main/internal/util"

	"github.com/yanun0323/pkg/logs"
)

const (
	/* ENV KEY */
	_keyPair     = "PAIR"
	_keyStrategy = "STG"

	/* Strategy Value */
	_strategyMaStandup        = "STANDUP"
	_strategyExchangeFollower = "FOLLOW"

	/* Provider Key */
	_keyProviderKline = "KLINE"
	_keyProviderTrade = "TRADE"

	/* Provider Value */
	_providerMock    = "MOCK"
	_providerBitopro = "BITOPRO"
	_providerBinance = "BINANCE"
)

func Run() {
	l := logs.New(util.LogLevel())

	pr := strings.ToUpper(os.Getenv(_keyPair))
	if len(pr) == 0 {
		l.Fatalf("missing '%s' environment key", _keyPair)
	}

	spr := strings.Split(pr, "/")
	if len(spr) != 2 {
		l.Fatalf("unsupported pair %s, expected connecting with '/'. e.g. BTC/USDT", pr)
	}

	if len(spr[0]) == 0 || len(spr[1]) == 0 {
		l.Fatalf("empty pair base/quote %s", pr)
	}

	pair := model.NewPair(spr[0], spr[1])

	stg := strings.ToUpper(os.Getenv(_keyStrategy))
	if len(stg) == 0 {
		l.Fatalf("missing '%s' environment key", _keyStrategy)
	}

	kls := strings.ToUpper(os.Getenv(_keyProviderKline))
	if len(kls) == 0 {
		l.Fatalf("missing '%s' environment key", _keyProviderKline)
	}

	tds := strings.ToUpper(os.Getenv(_keyProviderTrade))
	if len(tds) == 0 {
		l.Fatalf("missing '%s' environment key", _keyProviderTrade)
	}

	l = l.WithFields(map[string]interface{}{
		"pair":     pr,
		"strategy": stg,
		"kline":    kls,
		"trade":    tds,
	})

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
			kline = append(kline, mock.NewKlineProvider(ctx, pair, model.K1m))
		case _providerBitopro:
			k, err := bitopro.NewKlineProvider(ctx, pair, model.K1m)
			if err != nil {
				l.WithError(err).Fatal("init bitopro kline provider")
			}
			kline = append(kline, k)
		case _providerBinance:
		}
	}

	for _, td := range strings.Split(tds, ",") {
		switch td {
		case _providerMock:
			t, err := mock.NewTradeServer(ctx, pair, model.OrderTypeLimit, model.OrderTypeMarket)
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

		bot, err = service.NewMaStandupBot(ctx, pair, service.MaStandupBotProvider{
			Kline: kline[0],
			Trade: trade[0],
		})
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

		bot, err = service.NewExchangeFollowerBot(ctx, pair, service.ExchangeFollowerProvider{
			Source: kline[0],
			Target: kline[1],
			Trade:  trade[0],
		})
		if err != nil {
			l.WithError(err).Fatal("create exchange follower bot")
		}
	}

	if err := bot.Init(ctx); err != nil {
		l.WithError(err).Fatal("init bot")
	}

	go func() {
		if err := bot.Run(ctx); err != nil {
			l.WithError(err).Fatal("run bot")
		}
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
