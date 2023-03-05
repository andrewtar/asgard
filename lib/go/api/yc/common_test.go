package yc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"asgard/common/test"

	"bou.ke/monkey"
	"github.com/golang-jwt/jwt/v5"
)

var (
	testNowTime           = time.Date(2023, time.February, 26, 11, 51, 41, 0, time.UTC)
	testIamTokenValue     = "test_token"
	testIamTokenExpiresAt = time.Date(2023, time.February, 27, 15, 4, 5, 0, time.UTC)
)

type YcCommonTestSuite struct {
	test.BaseTestSuite

	serviceAccountId string
	publicKey        []byte
	timePatch        *monkey.PatchGuard
}

func (self *YcCommonTestSuite) SetupTest() {
	self.timePatch = self.SetupTime(testNowTime)
	setupServiceAccountKey(self.ReadFile("testdata/test_key.json"))

	parsedKey := map[string]string{}
	self.Nil(json.Unmarshal([]byte(*serviceAccountKey), &parsedKey))

	self.publicKey = []byte(parsedKey["public_key"])
	self.serviceAccountId = parsedKey["service_account_id"]
}

func (self *YcCommonTestSuite) TearDownTest() {
	self.timePatch.Unpatch()
}

func (self *YcCommonTestSuite) parseJwtToken(request *http.Request) *jwt.Token {
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

func (self *YcCommonTestSuite) getRequestData(request *http.Request) map[string]interface{} {
	bodyBytes, err := io.ReadAll(request.Body)
	self.Nil(err)

	parsed := map[string]interface{}{}
	self.Nil(json.Unmarshal(bodyBytes, &parsed))
	return parsed
}

func (self *YcCommonTestSuite) jwtTime(date time.Time) float64 {
	jwtTimeByte, err := jwt.NumericDate{date}.MarshalJSON()
	self.Nil(err)
	jwtTimeFloat, err := strconv.ParseFloat(string(jwtTimeByte), 64)
	self.Nil(err)
	return jwtTimeFloat
}

func (self *YcCommonTestSuite) setupTokenExchangeResponse(responseWithToken []byte) *httptest.Server {
	testServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		token := self.parseJwtToken(request)
		self.True(token.Valid)
		self.Equal(jwt.MapClaims{
			"exp": self.jwtTime(testNowTime.Add(1 * time.Hour)),
			"iat": self.jwtTime(testNowTime),
			"iss": self.serviceAccountId,
			"aud": []interface{}{*tokenExchangeUrl},
		},
			token.Claims,
		)
		_, err := response.Write(responseWithToken)
		self.Nil(err)
	}))
	*tokenExchangeUrl = testServer.URL
	return testServer
}

func (self *YcCommonTestSuite) setupTokenExchangeResponseWithError(errorMessage string, httpCode int) *httptest.Server {
	testServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		token := self.parseJwtToken(request)
		self.True(token.Valid)
		self.Equal(jwt.MapClaims{
			"exp": self.jwtTime(testNowTime.Add(1 * time.Hour)),
			"iat": self.jwtTime(testNowTime),
			"iss": self.serviceAccountId,
			"aud": []interface{}{*tokenExchangeUrl},
		},
			token.Claims,
		)
		response.WriteHeader(httpCode)
		_, err := response.Write([]byte(errorMessage))
		self.Nil(err)
	}))
	*tokenExchangeUrl = testServer.URL
	return testServer
}

func (self *YcCommonTestSuite) setupTranslationServer(
	iamToken string,
	expectedRequest map[string]interface{},
	responseData string,
) *httptest.Server {
	testServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		self.Equal(self.getRequestData(request), expectedRequest)
		self.Equal(fmt.Sprintf("Bearer %s", iamToken), request.Header.Get("Authorization"))
		self.Equal("application/json", request.Header.Get("Content-Type"))

		_, err := response.Write([]byte(responseData))
		self.Nil(err)
	}))
	*translationUrl = testServer.URL
	return testServer
}

func setupServiceAccountKey(key string) {
	*serviceAccountKey = key
}
