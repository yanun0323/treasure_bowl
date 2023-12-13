package service

import (
	"context"
	"errors"

	"main/internal/domain"
	"main/internal/entity"

	"github.com/yanun0323/pkg/logs"
)

type exchangeFollowerBot struct {
	l         logs.Logger
	pair      entity.Pair
	kpsSource domain.KlineProvideServer
	kpsTarget domain.KlineProvideServer
	ts        domain.TradeServer
}

type ExchangeFollowerProvider struct {
	Source domain.KlineProvideServer
	Target domain.KlineProvideServer
	Trade  domain.TradeServer
}

func (p *ExchangeFollowerProvider) Validate() error {
	if p.Source == nil {
		return errors.New("nil source provider")
	}

	if p.Target == nil {
		return errors.New("nil target provider")
	}

	if p.Trade == nil {
		return errors.New("nil trade server")
	}

	return nil
}

func NewExchangeFollowerBot(ctx context.Context, pr entity.Pair, pd ExchangeFollowerProvider) (domain.StrategyBot, error) {
	return &exchangeFollowerBot{
		l:         logs.Get(ctx).WithField("server", "exchange follower bot"),
		pair:      pr,
		kpsSource: pd.Source,
		kpsTarget: pd.Target,
		ts:        pd.Trade,
	}, nil
}

func (bot *exchangeFollowerBot) Init(ctx context.Context) error {
	return nil
}

func (bot *exchangeFollowerBot) Run(ctx context.Context) error {
	return nil
}

func (bot *exchangeFollowerBot) Shutdown(ctx context.Context) error {
	return nil
}
