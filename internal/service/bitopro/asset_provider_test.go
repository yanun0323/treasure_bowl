package bitopro

import (
	"context"
	"testing"
	"time"

	"main/internal/model"
	"main/pkg/infra"

	"github.com/stretchr/testify/suite"
)

func TestAssetProvider(t *testing.T) {
	suite.Run(t, new(AssetProviderSuite))
}

type AssetProviderSuite struct {
	suite.Suite
	ctx context.Context
}

func (su *AssetProviderSuite) SetupSuite() {
	su.ctx = context.Background()
	su.Require().NoError(infra.Init("config"))
}

func (su *AssetProviderSuite) Test() {
	su.T().Log("connect to private ws")
	ws, err := ConnectToPrivateWs()
	su.Require().NoError(err)

	su.T().Log("new asset provider")
	p, err := NewAssetProvider(ws)
	su.Require().NoError(err)

	su.T().Log("connecting")
	ch, err := p.Connect(su.ctx)
	su.Require().NoError(err)
	defer p.Disconnect(su.ctx)

	su.T().Log("connected")
	ctx, cancel := context.WithTimeout(su.ctx, 3*time.Second)
	defer cancel()

	su.T().Log("start consuming")
	select {
	case acc := <-ch:
		su.T().Logf("%+v", acc)
		acc.Balances.Iter(func(k string, b model.Balance) bool {
			su.T().Logf("%s: %+v", k, b)
			return true
		})
		su.NotEmpty(acc)
	case <-ctx.Done():
		su.T().Log("end up consuming")
		su.Fail("consume kline timeout")
	}
}
