package bitopro

import (
	"context"
	"testing"
	"time"

	"main/internal/model"
	"main/pkg/infra"

	"github.com/stretchr/testify/suite"
	"github.com/yanun0323/pkg/logs"
)

func TestTradeServer(t *testing.T) {
	suite.Run(t, new(TradeServerSuite))
}

type TradeServerSuite struct {
	suite.Suite
	ctx  context.Context
	l    logs.Logger
	pair model.Pair
}

func (su *TradeServerSuite) SetupSuite() {
	su.ctx = context.Background()
	su.Require().NoError(infra.Init("config-test"))
	su.pair = model.NewPair("usdt", "twd")
	su.l = logs.New(logs.LevelDebug)
}

func (su *TradeServerSuite) TestCreateAndCancelOrder() {
	server, err := NewTradeServer(su.ctx, su.pair)
	su.Require().NoError(err)

	acc, ch, err := server.Connect(su.ctx)
	su.Require().NoError(err)
	su.Require().NotNil(acc)
	defer server.Disconnect(su.ctx)

	acc, err = server.PushOrder(su.ctx, model.Order{
		Pair:   su.pair,
		Type:   model.OrderTypeLimit,
		Action: model.OrderActionBuy,
		Price:  "0.001",
		Amount: model.Amount{
			Total: "10",
		},
	})
	su.Require().NoError(err)
	su.Require().NotNil(acc)

	ctx, cancel := context.WithTimeout(su.ctx, 15*time.Second)
	defer cancel()

	select {
	case o := <-ch:
		o.Action = model.OrderActionCancelBuy
		su.l.Infof("consume order: %+v", o)
		acc, err = server.PushOrder(su.ctx, o)
		su.NoError(err)
		su.NotNil(acc)
	case <-ctx.Done():
		su.Fail("consume order timeout")
	}
}

func (su *TradeServerSuite) TestCancelOrderInTheBeginning() {
	server, err := NewTradeServer(su.ctx, su.pair)
	su.Require().NoError(err)

	{
		_, _, err := server.Connect(su.ctx)
		su.Require().NoError(err)

		_, err = server.PushOrder(su.ctx, model.Order{
			Pair:   su.pair,
			Type:   model.OrderTypeLimit,
			Action: model.OrderActionBuy,
			Price:  "0.001",
			Amount: model.Amount{
				Total: "10",
			},
		})
		su.Require().NoError(err)
		su.Require().NoError(server.Disconnect(su.ctx))
	}

	{
		acc, ch, err := server.Connect(su.ctx)
		su.Require().NoError(err)
		su.Require().NotNil(acc)
		defer server.Disconnect(su.ctx)

		ctx, cancel := context.WithTimeout(su.ctx, 15*time.Second)
		defer cancel()

		select {
		case o := <-ch:
			o.Action = model.OrderActionCancelBuy
			su.l.Infof("consume order: %+v", o)
			server.PushOrder(su.ctx, o)
			su.NoError(err)
			su.NotNil(acc)
		case <-ctx.Done():
			su.Fail("consume order timeout")
		}
	}
}
