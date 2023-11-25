package bitopro

import (
	"context"
	"testing"

	"main/pkg/infra"

	"github.com/stretchr/testify/suite"
)

func TestCommon(t *testing.T) {
	suite.Run(t, new(CommonSuite))
}

type CommonSuite struct {
	suite.Suite
	ctx context.Context
}

func (su *CommonSuite) SetupSuite() {
	su.ctx = context.Background()
	su.Require().NoError(infra.Init("config-test"))
}

func (su *CommonSuite) TestConnectToPublicClient() {
	c, err := ConnectToPublicClient()
	su.Require().NoError(err)
	su.NotNil(c)
}

func (su *CommonSuite) TestConnectToPrivateClient() {
	c, err := ConnectToPrivateClient()
	su.Require().NoError(err)
	su.NotNil(c)
}

func (su *CommonSuite) TestConnectToPublicWs() {
	c, err := ConnectToPublicWs()
	su.Require().NoError(err)
	su.NotNil(c)
}

func (su *CommonSuite) TestConnectToPrivateWs() {
	c, err := ConnectToPrivateWs()
	su.Require().NoError(err)
	su.NotNil(c)
}
