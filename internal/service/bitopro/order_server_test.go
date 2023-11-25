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

func TestOrderServer(t *testing.T) {
	suite.Run(t, new(OrderServerSuite))
}

type OrderServerSuite struct {
	suite.Suite
	ctx  context.Context
	l    logs.Logger
	pair model.Pair
}

func (su *OrderServerSuite) SetupSuite() {
	su.ctx = context.Background()
	su.Require().NoError(infra.Init("config-test"))
	su.pair = model.NewPair("usdt", "twd")
	su.l = logs.New("order server test suite", logs.LevelDebug)
}

func (su *OrderServerSuite) TestCreateAndCancelOrder() {
	wss, err := ConnectToPrivateWs()
	su.Require().NoError(err)

	client, err := ConnectToPrivateClient()
	su.Require().NoError(err)

	server, err := NewOrderServer(su.pair, wss, client)
	su.Require().NoError(err)

	ch, err := server.Connect(su.ctx)
	su.Require().NoError(err)
	defer server.DisConnect(su.ctx)

	su.Require().NoError(server.PushOrder(su.ctx, model.Order{
		Pair:   su.pair,
		Type:   model.OrderTypeLimit,
		Action: model.OrderActionBuy,
		Price:  "0.001",
		Amount: model.Amount{
			Total: "10",
		},
	}))

	ctx, cancel := context.WithTimeout(su.ctx, 15*time.Second)
	defer cancel()

	select {
	case o := <-ch:
		o.Action = model.OrderActionCancelBuy
		su.l.Infof("consume order: %+v", o)
		su.NoError(server.PushOrder(su.ctx, o))
	case <-ctx.Done():
		su.Fail("consume order timeout")
	}
}

func (su *OrderServerSuite) TestCancelOrderInTheBeginning() {
	wss, err := ConnectToPrivateWs()
	su.Require().NoError(err)

	client, err := ConnectToPrivateClient()
	su.Require().NoError(err)

	server, err := NewOrderServer(su.pair, wss, client)
	su.Require().NoError(err)

	{
		_, err := server.Connect(su.ctx)
		su.Require().NoError(err)

		su.Require().NoError(server.PushOrder(su.ctx, model.Order{
			Pair:   su.pair,
			Type:   model.OrderTypeLimit,
			Action: model.OrderActionBuy,
			Price:  "0.001",
			Amount: model.Amount{
				Total: "10",
			},
		}))

		su.Require().NoError(server.DisConnect(su.ctx))
	}

	{
		ch, err := server.Connect(su.ctx)
		su.Require().NoError(err)
		defer server.DisConnect(su.ctx)

		ctx, cancel := context.WithTimeout(su.ctx, 15*time.Second)
		defer cancel()

		select {
		case o := <-ch:
			o.Action = model.OrderActionCancelBuy
			su.l.Infof("consume order: %+v", o)
			su.NoError(server.PushOrder(su.ctx, o))
		case <-ctx.Done():
			su.Fail("consume order timeout")
		}
	}
}
