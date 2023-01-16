package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMetadata(t *testing.T) {
	data, err := os.ReadFile("/home/andrew/workspace/andrew/github.com/asgard/key.json")
	assert.Nil(t, err)
	*serviceAccountKey = string(data)
	service := YandexCloudAuthService{}
	token, err := service.GetIamToken()
	assert.Nil(t, err)
	assert.Equal(t, IAMToken{}, token)
}
