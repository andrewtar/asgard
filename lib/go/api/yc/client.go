package yc

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var translationUrl = flag.String(
	"translation-api-url",
	"https://translate.api.cloud.yandex.net/translate/v2/translate",
	"Yandex Cloud token exchange url",
)

type YCCloud struct {
	mutex sync.Mutex
	token *IAMToken
}

type Language int

const (
	Russian Language = iota
	English
)

func (self Language) String() string {
	// https://en.wikipedia.org/wiki/ISO_639-1
	switch self {
	case Russian:
		return "ru"
	case English:
		return "en"
	}
	panic(fmt.Sprintf("Unknown Language(%s)", self))
}

type translationRequest struct {
	SourceLanguageCode string   `json:"sourceLanguageCode"`
	TargetLanguageCode string   `json:"targetLanguageCode"`
	Format             string   `json:"format"`
	Texts              []string `json:"texts"`
}

type translationResponse struct {
	Translations []struct {
		Text                 string `json:"text"`
		DetectedLanguageCode string `json:"detectedLanguageCode"`
	} `json:"translations"`
}

func (self *YCCloud) Translate(from, to Language, text string) (string, error) {
	responseData, err := self.makePostApiCall(*translationUrl, translationRequest{
		SourceLanguageCode: from.String(),
		TargetLanguageCode: to.String(),
		Format:             "PLAIN_TEXT",
		Texts:              []string{text},
	})
	if err != nil {
		return "", fmt.Errorf("failed to request transaltion %w", err)
	}

	parsedResponse := translationResponse{}
	err = json.Unmarshal(responseData, &parsedResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse transaltion response: %w", err)
	}
	return parsedResponse.Translations[0].Text, nil
}

func (self *YCCloud) checkTokenOrRefresh() (string, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.token == nil || self.token.timeRefresh() {
		token, err := GetIamToken()
		if err != nil {
			return "", fmt.Errorf("failed to request a new token: %w", err)
		}
		self.token = &token
	}

	return self.token.Token, nil
}

func (self *YCCloud) makePostApiCall(url string, requestBody interface{}) ([]byte, error) {
	token, err := self.checkTokenOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	serializeJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize json body: %w", err)
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(serializeJson))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("received %d error code and message %s", response.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
