package service

import (
	"context"
	"errors"

	"main/internal/domain"
	"main/internal/model"

	"github.com/yanun0323/pkg/logs"
)

type MaStandupBot struct {
	l    logs.Logger
	pair model.Pair
	kps  domain.KlineProvideServer
	ts   domain.TradeServer
}

type MaStandupBotProvider struct {
	Kline domain.KlineProvideServer
	Trade domain.TradeServer
}

func (p *MaStandupBotProvider) Validate() error {
	if p.Kline == nil {
		return errors.New("nil kline provider")
	}

	if p.Trade == nil {
		return errors.New("nil trade server")
	}

	return nil
}

func NewMaStandupBot(ctx context.Context, pr model.Pair, pd MaStandupBotProvider) (domain.StrategyBot, error) {
	return &MaStandupBot{
		l:    logs.Get(ctx).WithField("server", "ma standup bot"),
		pair: pr,
		kps:  pd.Kline,
		ts:   pd.Trade,
	}, nil
}

func (bot *MaStandupBot) Init(ctx context.Context) error {
	return nil
}

func (bot *MaStandupBot) Run(ctx context.Context) error {
	return nil
}

func (bot *MaStandupBot) Shutdown(ctx context.Context) error {
	return nil
}
