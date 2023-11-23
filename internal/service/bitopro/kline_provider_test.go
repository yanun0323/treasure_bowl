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
	su.Require().NoError(infra.Init("config"))
}

func (su *KlineProviderSuite) TestKlineProvider() {
	su.T().Log("new kline provider")
	p, err := NewKlineProvider(model.NewPair("BTC", "TWD"), model.K1m)
	su.Require().NoError(err)

	su.T().Log("connecting")
	ch, err := p.Connect(su.ctx)
	su.Require().NoError(err)
	defer p.Disconnect(su.ctx)

	su.T().Log("connected")
	ctx, cancel := context.WithTimeout(su.ctx, 5*time.Second)
	defer cancel()
	ks := make(map[int64]model.Kline, 200)

	su.T().Log("start consuming")
	func() {
		for {
			select {
			case k := <-ch:
				su.T().Logf("%+v", k)
				ks[k.Timestamp] = k
				return
			case <-ctx.Done():
				su.T().Log("end up consuming")
				return
			}
		}
	}()

	su.NotEmpty(ks)
}
