package yc

import (
	"fmt"
	"net/http"
	"testing"

	"asgard/common/test"

	"github.com/stretchr/testify/suite"
)

var (
	testErrorMessage     = "test_server_error"
	testIamTokenResponse = test.ToJson(map[string]interface{}{
		"iamToken":  testIamTokenValue,
		"expiresAt": testIamTokenExpiresAt,
	})
	testIamTokenParsed = IAMToken{
		Token:     testIamTokenValue,
		ExpiresAt: testIamTokenExpiresAt,
	}
)

type IamTestSuite struct {
	YcCommonTestSuite
}

func (self *IamTestSuite) TestGetIamToken() {
	testTokenExchangeServer := self.setupTokenExchangeResponse(testIamTokenResponse)
	defer testTokenExchangeServer.Close()

	token, err := GetIamToken()

	self.Nil(err)
	self.Equal(testIamTokenParsed, token)
}

func (self *IamTestSuite) TestGetIamTokenReturnErrorIfNoToken() {
	*serviceAccountKeyPath = ""
	_, err := GetIamToken()
	self.Contains(err.Error(), "service account key cannot be empty")
}

func (self *IamTestSuite) TestGetIamTokenReturnErrorIfInvalidToken() {
	*serviceAccountKeyPath = "testdata/test_key_invalid.json"
	_, err := GetIamToken()
	self.Contains(err.Error(), "failed to parse service key")
}

func (self *IamTestSuite) TestGetIamTokenReturnErrorIfInvalidPrivateKey() {
	self.setupServiceAccountKey("testdata/test_key_invalid_private_key.json")
	_, err := GetIamToken()
	self.Contains(err.Error(), "failed to parse private key")
}

func (self *IamTestSuite) TestGetIamTokenReturnErrorIfReceivedHttpErrorDuringExchange() {
	testTokenExchangeServer := self.setupTokenExchangeResponseWithError(
		testErrorMessage,
		http.StatusInternalServerError,
	)
	defer testTokenExchangeServer.Close()

	_, err := GetIamToken()

	self.Contains(err.Error(), "failed to exchange JWT on IAM token")
	self.Contains(err.Error(), fmt.Sprint(http.StatusInternalServerError))
	self.Contains(err.Error(), testErrorMessage)
}

func TestIamTestSuite(t *testing.T) {
	suite.Run(t, new(IamTestSuite))
}
