package cmd

import (
	"fmt"

	"asgard/common/api/bored"
	"asgard/common/api/yc"
	"asgard/common/debug/metadata"
	"asgard/common/log"
)

type boredClient interface {
	GetActivity() (bored.BoredApiResponse, error)
}

type transaltionClient interface {
	Translate(from, to yc.Language, text string) (string, error)
}

type HaveFunCommand struct {
	BoredClient       boredClient
	TransaltionClient transaltionClient
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
	translatedActivity, err := self.TransaltionClient.Translate(yc.English, yc.Russian, activity.Activity)
	if err != nil {
		log.Logger().
			WithError(err).
			WithField("activity", activity.Activity).
			Error("Failed to to transtate activity")
		return activity.Activity, nil
	}
	return translatedActivity, nil
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
