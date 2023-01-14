package bored

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var boredApiFlag = flag.String("bored-api-endpoint", "https://www.boredapi.com/api/activity", "Bored API endpoint")

type BoredApiResponse struct {
	Activity string `json: "activity"`
}

type BoredClient struct {
	debug bool
}

func NewBoredClient(debug bool) BoredClient {
	return BoredClient{debug: debug}
}

func (this BoredClient) GetActivity() (BoredApiResponse, error) {
	response, err := http.Get(*boredApiFlag)
	if err != nil {
		return BoredApiResponse{}, fmt.Errorf("Failed to request activity: %s\n", err)
	}

	if !(response.StatusCode >= 200 && response.StatusCode < 300) {
		return BoredApiResponse{}, fmt.Errorf("Failed with code=%d to request activity: %s\n", response.StatusCode, err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if this.debug {
		log.Printf("Response from Bored API: %s", string(responseBody))
	}
	if err != nil {
		return BoredApiResponse{}, fmt.Errorf("Failed to read response from Bored API: %s\n", err)
	}

	result := BoredApiResponse{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return BoredApiResponse{}, fmt.Errorf("Failed to parse response from Bored API: %s\n", err)
	}
	return result, nil
}
