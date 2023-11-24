package bitopro

import (
	"context"
	"testing"
	"time"

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

func (su *AssetProviderSuite) TestAssetProvider() {
	ws, err := ConnectToPrivateWs()
	su.Require().NoError(err)

	p, err := NewAssetProvider(ws)
	su.Require().NoError(err)

	ch, err := p.Connect(su.ctx)
	su.Require().NoError(err)
	defer p.Disconnect(su.ctx)

	ctx, cancel := context.WithTimeout(su.ctx, 3*time.Second)
	defer cancel()

	select {
	case acc := <-ch:
		su.NotZero(acc.Balances.Len(), "%+v", acc.Balances.Clone())
	case <-ctx.Done():
		su.Fail("consume kline timeout")
	}
}

func (su *AssetProviderSuite) TestRestfulApi() {
	c, err := ConnectToPrivateClient()
	su.Require().NoError(err)
	acc := c.GetAccountBalance()
	su.NotEmpty(acc.Data, "%+v", acc.Data)
}
