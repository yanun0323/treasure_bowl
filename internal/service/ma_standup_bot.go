package service

import (
	"context"

	"main/internal/domain"
	"main/internal/entity"

	"github.com/yanun0323/pkg/logs"
)

type maStandupBot struct {
	l    logs.Logger
	pair entity.Pair
	kps  domain.KlineProvideServer
	ts   domain.TradeServer
}

func NewMaStandupBot(ctx context.Context, pr entity.Pair, pd StdBotProvider) (domain.StrategyBot, error) {
	return &maStandupBot{
		l:    logs.Get(ctx).WithField("server", "ma standup bot"),
		pair: pr,
		kps:  pd.Kline,
		ts:   pd.Trade,
	}, nil
}

func (bot *maStandupBot) Init(ctx context.Context) error {
	return nil
}

func (bot *maStandupBot) Run(ctx context.Context) error {
	return nil
}

func (bot *maStandupBot) Shutdown(ctx context.Context) error {
	return nil
}
