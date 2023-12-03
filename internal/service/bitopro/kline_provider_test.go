package bitopro

import (
	"context"
	"testing"
	"time"

	"main/internal/model"
	"main/pkg/infra"

	"github.com/stretchr/testify/suite"
)

func TestKlineProvider(t *testing.T) {
	suite.Run(t, new(KlineProviderSuite))
}

type KlineProviderSuite struct {
	suite.Suite
	ctx context.Context
}

func (su *KlineProviderSuite) SetupSuite() {
	su.ctx = context.Background()
	su.Require().NoError(infra.Init("config-test"))
}

func (su *KlineProviderSuite) TestKlineProvider() {
	p, err := NewKlineProvider(su.ctx, model.NewPair("BTC", "TWD"), model.K1m)
	su.Require().NoError(err)

	ch, err := p.Connect(su.ctx, 10)
	su.Require().NoError(err)
	defer p.Disconnect(su.ctx)

	ctx, cancel := context.WithTimeout(su.ctx, 3*time.Second)
	defer cancel()

	select {
	case k := <-ch:
		su.NotEmpty(k, "%+v", k)
	case <-ctx.Done():
		su.Fail("consume kline timeout")
	}
}
