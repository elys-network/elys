package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MsgsTestSuite struct {
	suite.Suite
	govModule string
}

func TestMsgsTestSuite(t *testing.T) {
	suite.Run(t, new(MsgsTestSuite))
}

func (suite *MsgsTestSuite) SetupTest() {
	suite.govModule = "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3"
}
