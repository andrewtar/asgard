package test

import (
	"encoding/json"
	"os"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/suite"
)

type BaseTestSuite struct {
	suite.Suite
}

func (self *BaseTestSuite) ReadFile(path string) string {
	data, err := os.ReadFile(path)
	self.Nil(err)
	return string(data)
}

func (self *BaseTestSuite) SetupTime(date time.Time) *monkey.PatchGuard {
	return monkey.Patch(time.Now, func() time.Time { return date })
}

func ToJson(data map[string]interface{}) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return jsonData
}
