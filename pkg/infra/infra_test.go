package infra

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

func TestInfra(t *testing.T) {
	suite.Run(t, new(InfraSuite))
}

type InfraSuite struct {
	suite.Suite
	ctx context.Context
}

func (su *InfraSuite) SetupSuite() {
	su.ctx = context.Background()
}

func (su *InfraSuite) TestInitConfig() {
	su.Require().NoError(Init("config-test"))
	su.Require().NotEmpty(viper.GetString("log.level"))
}
