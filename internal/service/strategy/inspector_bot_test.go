package strategy

import (
	"context"
	"main/internal/domain"
	"main/internal/entity"
	"main/internal/service/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestInspectorBot(t *testing.T) {
	suite.Run(t, new(InspectorBotSuite))
}

type InspectorBotSuite struct {
	suite.Suite
	ctx     context.Context
	pair    entity.Pair
	targets []entity.KlineType
	bot     domain.StrategyBot
}

func (su *InspectorBotSuite) SetupSuite() {
	ctx := context.Background()
	pair := entity.NewPair("BTC", "USDT")
	targets := []entity.KlineType{entity.K5m, entity.K15m, entity.K1h, entity.K4h, entity.K1d}
	kps := mock.NewKlineProvider(ctx, pair, targets...)
	bot, err := NewInspectorBot(ctx, pair, kps)
	su.Require().NoError(err)

	su.ctx = ctx
	su.pair = pair
	su.targets = targets
	su.bot = bot
}

func (su *InspectorBotSuite) Test() {
	su.T().Log("start test:", su.T().Name())
	su.Require().NoError(su.bot.Init(su.ctx))
	su.Require().NoError(su.bot.Run(su.ctx))
	time.Sleep(3 * time.Second)
	su.Require().NoError(su.bot.Shutdown(su.ctx))
}
