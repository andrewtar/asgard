package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"asgard/common/api/bored"
	"asgard/common/api/yc"
	"asgard/common/debug/metadata"
)

type HaveFunCommandTestSuite struct {
	suite.Suite
}

func (self *HaveFunCommandTestSuite) TestHandle() {
	testActivity := "test_activity"
	testTransaltedActivity := "test_activity_russian"
	command := HaveFunCommand{
		BoredClient: BoredApiStub{
			answer: bored.BoredApiResponse{Activity: testActivity},
		},
		TransaltionClient: self.setupTransaltion(yc.English, yc.Russian, testActivity, testTransaltedActivity),
	}
	answer, err := command.Handle()
	self.Nil(err)
	self.Equal(testTransaltedActivity, answer)
}

func (self *HaveFunCommandTestSuite) TestHandleIgnoreErrorFromTranslation() {
	testActivity := "test_activity"
	command := HaveFunCommand{
		BoredClient: BoredApiStub{
			answer: bored.BoredApiResponse{Activity: testActivity},
		},
		TransaltionClient: self.setupTransaltionErrorResponse(yc.English, yc.Russian, testActivity, fmt.Errorf("test_error")),
	}
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

func (self *HaveFunCommandTestSuite) setupTransaltion(
	expectedFrom, expectedTo yc.Language,
	expectedText string,
	answer string,
) TransaltionApiStub {
	return TransaltionApiStub{
		Assertions:   assert.New(self.T()),
		expectedFrom: expectedFrom,
		expectedTo:   expectedTo,
		expectedText: expectedText,
		answer:       answer,
	}
}

func (self *HaveFunCommandTestSuite) setupTransaltionErrorResponse(
	expectedFrom, expectedTo yc.Language,
	expectedText string,
	err error,
) TransaltionApiStub {
	return TransaltionApiStub{
		Assertions:   assert.New(self.T()),
		expectedFrom: expectedFrom,
		expectedTo:   expectedTo,
		expectedText: expectedText,
		err:          err,
	}
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

type TransaltionApiStub struct {
	*assert.Assertions
	expectedFrom yc.Language
	expectedTo   yc.Language
	expectedText string

	answer string
	err    error
}

func (self TransaltionApiStub) Translate(from, to yc.Language, text string) (string, error) {
	self.Equal(self.expectedFrom, from)
	self.Equal(self.expectedTo, to)
	self.Equal(self.expectedText, text)
	return self.answer, self.err
}
