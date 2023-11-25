package bitopro

import (
	"context"
	"testing"

	"main/pkg/infra"

	"github.com/stretchr/testify/suite"
)

func TestOrderServer(t *testing.T) {
	suite.Run(t, new(OrderServerSuite))
}

type OrderServerSuite struct {
	suite.Suite
	ctx context.Context
}

func (su *OrderServerSuite) SetupSuite() {
	su.ctx = context.Background()
	su.Require().NoError(infra.Init("config-test"))
}

func (su *OrderServerSuite) SetupTest() {

}

func (su *OrderServerSuite) TearDownTest() {

}

func (su *OrderServerSuite) TearDownSuite() {

}

func (su *OrderServerSuite) Test() {

}

func (su *OrderServerSuite) TestWithCase() {
	testCases := []struct {
		desc string
	}{
		{},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)

		})
	}
}
