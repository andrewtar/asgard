package yc

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

var (
	testWordToTranslate            = "apple"
	testTranslatedWord             = "яблоко"
	testExpectedTranslationRequest = map[string]interface{}{
		"sourceLanguageCode": "en",
		"targetLanguageCode": "ru",
		"format":             "PLAIN_TEXT",
		"texts":              []interface{}{testWordToTranslate},
	}
	testExpectedTranslationResponse = fmt.Sprintf(`{"translations":[{"text": "%s"}]}`, testTranslatedWord)
)

type ClientTestSuite struct {
	YcCommonTestSuite
	client YCCloud
}

func (self *ClientTestSuite) SetupTest() {
	self.YcCommonTestSuite.SetupTest()
	self.client = YCCloud{}
}

func (self *ClientTestSuite) TestTranslateIfNoToken() {
	testTokenExchangeServer := self.setupTokenExchangeResponse(testIamTokenResponse)
	defer testTokenExchangeServer.Close()

	testTransaltionServer := self.setupTranslationServer(
		testIamTokenValue,
		testExpectedTranslationRequest,
		testExpectedTranslationResponse,
	)
	defer testTransaltionServer.Close()

	result, err := self.client.Translate(English, Russian, testWordToTranslate)

	self.Nil(err)
	self.Equal(testTranslatedWord, result)
}

func (self *ClientTestSuite) TestTranslateIfExpiredToken() {
	self.client = YCCloud{
		token: &IAMToken{
			Token: testIamTokenValue,
			ExpiresAt: testNowTime.
				Add(30 * time.Minute).
				Add(1 * time.Second),
		},
	}

	testTransaltionServer := self.setupTranslationServer(
		testIamTokenValue,
		testExpectedTranslationRequest,
		testExpectedTranslationResponse,
	)
	defer testTransaltionServer.Close()

	result, err := self.client.Translate(English, Russian, testWordToTranslate)

	self.Nil(err)
	self.Equal(testTranslatedWord, result)
}

func (self *ClientTestSuite) TestTranslateReusePreviousToken() {
	self.client = YCCloud{
		token: &IAMToken{
			Token: testIamTokenValue,
			ExpiresAt: testIamTokenExpiresAt.
				Add(30 * time.Minute).
				Add(-1 * time.Second),
		},
	}

	testTransaltionServer := self.setupTranslationServer(
		testIamTokenValue,
		testExpectedTranslationRequest,
		testExpectedTranslationResponse,
	)
	defer testTransaltionServer.Close()

	result, err := self.client.Translate(English, Russian, testWordToTranslate)

	self.Nil(err)
	self.Equal(testTranslatedWord, result)
}

func (self *ClientTestSuite) TestTranslateReturnsErrorIfFailedToGetToken() {
	testTokenExchangeServer := self.setupTokenExchangeResponseWithError(
		testErrorMessage,
		http.StatusInternalServerError,
	)
	defer testTokenExchangeServer.Close()

	_, err := self.client.Translate(English, Russian, testWordToTranslate)

	self.Contains(err.Error(), "failed to refresh token")
	self.Contains(err.Error(), fmt.Sprint(http.StatusInternalServerError))
	self.Contains(err.Error(), testErrorMessage)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
