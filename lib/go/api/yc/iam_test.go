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
	setupServiceAccountKey("")
	_, err := GetIamToken()
	self.Contains(err.Error(), "service account key cannot be empty")
}

func (self *IamTestSuite) TestGetIamTokenReturnErrorIfInvalidToken() {
	setupServiceAccountKey("invalid_key")
	_, err := GetIamToken()
	self.Contains(err.Error(), "failed to parse service key")
}

func (self *IamTestSuite) TestGetIamTokenReturnErrorIfInvalidPrivateKey() {
	setupServiceAccountKey(string(test.ToJson(map[string]interface{}{
		"id":                 "test_id",
		"service_account_id": "test_service_account_id",
		"created_at":         "2023-01-15T13:49:17.907541493Z",
		"key_algorithm":      "RSA_2048",
		"public_key":         "invalid_public_key",
		"private_key":        "invalid_private_key",
	})))

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
