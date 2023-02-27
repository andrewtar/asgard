package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

var (
	testNowTime           = time.Date(2023, time.February, 26, 11, 51, 41, 0, time.UTC)
	testIamTokenValue     = "test_token"
	testIamTokenExpiresAt = time.Date(2023, time.February, 27, 15, 4, 5, 0, time.UTC)
	testIamTokenJson      = toJson(map[string]interface{}{
		"iamToken":  testIamTokenValue,
		"expiresAt": testIamTokenExpiresAt,
	})
	testIamToken = IAMToken{
		Token:     testIamTokenValue,
		ExpiresAt: testIamTokenExpiresAt,
	}
)

type IamTestSuite struct {
	suite.Suite
	service          YandexCloudAuthService
	publicKey        []byte
	serviceAccountId string
	timePatch        *monkey.PatchGuard
}

func (self *IamTestSuite) SetupTest() {
	self.timePatch = self.setupTime(testNowTime)
	self.setupServiceAccountKey(self.readFile("testdata/test_key.json"))

	parsedKey := map[string]string{}
	self.Nil(json.Unmarshal([]byte(*serviceAccountKey), &parsedKey))

	self.publicKey = []byte(parsedKey["public_key"])
	self.serviceAccountId = parsedKey["service_account_id"]
}

func (self *IamTestSuite) TearDownTest() {
	self.timePatch.Unpatch()
}

func (self *IamTestSuite) TestGetIamToken() {
	testTokenExchangeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		token := self.parseJwtToken(request)
		self.True(token.Valid)

		self.Equal(
			jwt.MapClaims{
				"exp": self.jwtTime(testNowTime.Add(1 * time.Hour)),
				"iat": self.jwtTime(testNowTime),
				"iss": self.serviceAccountId,
				"aud": []interface{}{*tokenExchangeUrl},
			},
			token.Claims,
		)
		_, err := response.Write(testIamTokenJson)
		self.Nil(err)
	}))
	defer testTokenExchangeServer.Close()
	*tokenExchangeUrl = testTokenExchangeServer.URL

	token, err := self.service.GetIamToken()
	self.Nil(err)
	self.Equal(testIamToken, token)
}

func (self *IamTestSuite) TestGetIamTokenReturnErrorIfNoToken() {
	self.setupServiceAccountKey("")

	_, err := self.service.GetIamToken()
	self.Contains(err.Error(), "service account key cannot be empty")
}

func (self *IamTestSuite) TestGetIamTokenReturnErrorIfInvalidToken() {
	self.setupServiceAccountKey("invalid_key")

	_, err := self.service.GetIamToken()
	self.Contains(err.Error(), "failed to parse service key")
}

func (self *IamTestSuite) TestGetIamTokenReturnErrorIfInvalidPrivateKey() {
	self.setupServiceAccountKey(string(toJson(map[string]interface{}{
		"id":                 "test_id",
		"service_account_id": "test_service_account_id",
		"created_at":         "2023-01-15T13:49:17.907541493Z",
		"key_algorithm":      "RSA_2048",
		"public_key":         "invalid_public_key",
		"private_key":        "invalid_private_key",
	})))

	_, err := self.service.GetIamToken()
	self.Contains(err.Error(), "failed to parse private key")
}

func (self *IamTestSuite) TestGetIamTokenReturnErrorIfReceivedHttpErrorDuringExchange() {
	testTokenExchangeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		token := self.parseJwtToken(request)
		self.True(token.Valid)

		self.Equal(
			jwt.MapClaims{
				"exp": self.jwtTime(testNowTime.Add(1 * time.Hour)),
				"iat": self.jwtTime(testNowTime),
				"iss": self.serviceAccountId,
				"aud": []interface{}{*tokenExchangeUrl},
			},
			token.Claims,
		)
		response.WriteHeader(http.StatusInternalServerError)
		_, err := response.Write([]byte("test_server_error"))
		self.Nil(err)
	}))
	defer testTokenExchangeServer.Close()
	*tokenExchangeUrl = testTokenExchangeServer.URL

	_, err := self.service.GetIamToken()
	self.Contains(err.Error(), "failed to exchange JWT on IAM token")
	self.Contains(err.Error(), fmt.Sprint(http.StatusInternalServerError))
	self.Contains(err.Error(), "test_server_error")
}

func TestIamTestSuite(t *testing.T) {
	suite.Run(t, new(IamTestSuite))
}

func (self *IamTestSuite) readFile(path string) string {
	data, err := os.ReadFile(path)
	self.Nil(err)
	return string(data)
}

func (self *IamTestSuite) jwtTime(date time.Time) float64 {
	jwtTimeByte, err := jwt.NumericDate{date}.MarshalJSON()
	self.Nil(err)
	jwtTimeFloat, err := strconv.ParseFloat(string(jwtTimeByte), 64)
	self.Nil(err)
	return jwtTimeFloat
}

func (self *IamTestSuite) setupTime(date time.Time) *monkey.PatchGuard {
	return monkey.Patch(time.Now, func() time.Time { return date })
}

func (self *IamTestSuite) setupServiceAccountKey(key string) {
	*serviceAccountKey = key
}

func (self *IamTestSuite) parseJwtToken(request *http.Request) *jwt.Token {
	bodyBytes, err := io.ReadAll(request.Body)
	self.Nil(err)

	parsed := map[string]string{}
	self.Nil(json.Unmarshal(bodyBytes, &parsed))
	self.Len(parsed, 1)
	jwtToken := parsed["jwt"]

	token, err := jwt.Parse(string(jwtToken), func(t *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(self.publicKey)
	})
	return token
}

func toJson(data map[string]interface{}) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return jsonData
}
