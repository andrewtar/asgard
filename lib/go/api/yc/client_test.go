package yc

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
}

func (self *ClientTestSuite) TestGetIamToken() {
	*serviceAccountKey = "test"

	client := YCCloud{}
	result, err := client.Translate(English, Russian, "apple")
	self.Equal("--", result)
	self.Nil(err)
	// self.Equal("--", err.Error())
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
