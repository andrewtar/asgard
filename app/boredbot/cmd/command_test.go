package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"asgard/common/api/bored"
	"asgard/common/debug/metadata"
)

type HaveFunCommandTestSuite struct {
	suite.Suite
}

func (self *HaveFunCommandTestSuite) TestHandle() {
	testActivity := "test_activity"
	command := HaveFunCommand{BoredClient: BoredApiStub{
		answer: bored.BoredApiResponse{Activity: testActivity},
	}}
	answer, err := command.Handle()
	self.Nil(err)
	self.Equal(testActivity, answer)
}

func (self *HaveFunCommandTestSuite) TestHandleReturnErrorFromBoredApi() {
	testError := fmt.Errorf("test_error")
	command := HaveFunCommand{BoredClient: BoredApiStub{
		err: testError,
	}}
	_, err := command.Handle()
	self.Contains(err.Error(), "failed to get activity")
	self.Contains(err.Error(), testError.Error())
}

type GetInformationCommandTestSuite struct {
	suite.Suite
}

func (self *GetInformationCommandTestSuite) TestHandle() {
	command := GetInformationCommand{}
	answer, err := command.Handle()
	self.Nil(err)
	self.Equal(
		fmt.Sprintf("Version: %s\nBuild Time: %s", metadata.Version, metadata.BuildTime),
		answer,
	)
}

func TestCommand(t *testing.T) {
	suite.Run(t, new(HaveFunCommandTestSuite))
	suite.Run(t, new(GetInformationCommandTestSuite))
}

type BoredApiStub struct {
	answer bored.BoredApiResponse
	err    error
}

func (self BoredApiStub) GetActivity() (bored.BoredApiResponse, error) {
	if self.err != nil {
		return bored.BoredApiResponse{}, self.err
	}
	return self.answer, nil
}
