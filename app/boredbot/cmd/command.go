package cmd

import (
	"fmt"

	"asgard/common/api/bored"
	"asgard/common/debug/metadata"
)

type boredClient interface {
	GetActivity() (bored.BoredApiResponse, error)
}

type HaveFunCommand struct {
	BoredClient boredClient
}

func (HaveFunCommand) GetDescription() string {
	return "Get an idea for fun"
}

func (HaveFunCommand) GetId() string {
	return "have_fun"
}

func (self HaveFunCommand) Handle() (string, error) {
	activity, err := self.BoredClient.GetActivity()
	if err != nil {
		return "", fmt.Errorf("failed to get activity: %w", err)
	}
	return activity.Activity, nil
}

type GetInformationCommand struct{}

func (GetInformationCommand) GetDescription() string {
	return "Get information about bot"
}

func (GetInformationCommand) GetId() string {
	return "info"
}

func (GetInformationCommand) Handle() (string, error) {
	return fmt.Sprintf("Version: %s\nBuild Time: %s", metadata.Version, metadata.BuildTime), nil
}
