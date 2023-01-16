package auth

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	tokenExchangeUrl = "https://iam.api.cloud.yandex.net/iam/v1/tokens"
)

var serviceAccountKey = flag.String("yandex-cloud-service-account-key", "", "Telegram API endpoint")

type servicekey struct {
	KeyId           string `json:"id"`
	ServiceAcountId string `json:"service_account_id"`
	PublicKey       string `json:"public_key"`
	PrivateKey      string `json:"private_key"`
}

type jwtExchangeRequset struct {
	Jwt string `json:"jwt"`
}

type jwtExchangeResponse struct {
	Token     string `json:"iamToken"`
	ExpiresAt string `json:"expiresAt"`
}

type IAMToken struct {
	Token     string
	ExpiresAt time.Time
}

type YandexCloudAuthService struct {
}

func (YandexCloudAuthService) GetIamToken() (IAMToken, error) {
	if *serviceAccountKey == "" {
		return IAMToken{}, fmt.Errorf("service account key cannot be null")
	}

	parsedKey, err := parseKey(*serviceAccountKey)
	if err != nil {
		return IAMToken{}, fmt.Errorf("failed to parse service key: %w", err)
	}

	jwtToken, err := generateJWTToken(parsedKey)
	if err != nil {
		return IAMToken{}, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	iamResponse, err := exchangeJwtOnIam(jwtToken)
	if err != nil {
		return IAMToken{}, fmt.Errorf("failed to exchange JWT on IAM token: %w", err)
	}
	// 2006-01-02T15:04:05.999999999Z07:00
	// expiresAt, err := time.Parse(iamResponse.ExpiresAt, "2006-01-02T15:04:05.999999999Z07:00")
	expiresAt, err := time.Parse(time.RFC3339Nano, iamResponse.ExpiresAt)
	if err != nil {
		return IAMToken{}, fmt.Errorf("failed to calculate expiration time of the IAM token: %w", err)
	}
	return IAMToken{
		Token:     iamResponse.Token,
		ExpiresAt: expiresAt,
	}, nil
}

func parseKey(key string) (servicekey, error) {
	parsedKey := servicekey{}
	err := json.Unmarshal([]byte(key), &parsedKey)
	return parsedKey, err
}

func generateJWTToken(key servicekey) (string, error) {
	token := &jwt.Token{
		Method: jwt.SigningMethodPS256,
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": "PS256",
			"kid": key.KeyId,
		},
		Claims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    key.ServiceAcountId,
			Audience:  []string{tokenExchangeUrl},
		},
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key.PrivateKey))
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}

func exchangeJwtOnIam(jwtToken string) (jwtExchangeResponse, error) {
	serializeJson, err := json.Marshal(jwtExchangeRequset{Jwt: jwtToken})
	if err != nil {
		return jwtExchangeResponse{}, fmt.Errorf("failed to serialize jwt token for request: %w", err)
	}
	request, err := http.NewRequest(http.MethodPost, tokenExchangeUrl, bytes.NewBuffer(serializeJson))
	if err != nil {
		return jwtExchangeResponse{}, fmt.Errorf("failed to prepare request: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return jwtExchangeResponse{}, fmt.Errorf("failed to send request: %w", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jwtExchangeResponse{}, fmt.Errorf("failed to read response: %w", err)
	}
	response.Body.Close()

	parsedResponse := jwtExchangeResponse{}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return jwtExchangeResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}
	return parsedResponse, nil
}
