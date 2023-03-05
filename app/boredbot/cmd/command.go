package cmd

import (
	"fmt"

	"asgard/common/api/bored"
	"asgard/common/debug/metadata"
)

type Command interface {
	GetDescription() string
	GetId() string
	Handle() (string, error)
}

func GetCommands(boredClient boredClient) (map[string]Command, error) {
	commands := []Command{
		haveFunCommand{boredClient: boredClient},
		getInformationCommand{},
	}
	commandsTable := map[string]Command{}
	for _, command := range commands {
		_, alreadyExists := commandsTable[command.GetId()]
		if alreadyExists {
			fmt.Errorf("Commad with id %s already registered", command.GetId())
		}
		commandsTable[command.GetId()] = command
	}
	return commandsTable, nil
}

type boredClient interface {
	GetActivity() (bored.BoredApiResponse, error)
}

type haveFunCommand struct {
	boredClient boredClient
}

func (haveFunCommand) GetDescription() string {
	return "Get an idea for fun"
}

func (haveFunCommand) GetId() string {
	return "have_fun"
}

func (self haveFunCommand) Handle() (string, error) {
	activity, err := self.boredClient.GetActivity()
	if err != nil {
		return "", fmt.Errorf("failed to get activity: %w", err)
	}
	return activity.Activity, nil
}

type getInformationCommand struct{}

func (getInformationCommand) GetDescription() string {
	return "Get information about bot"
}

func (getInformationCommand) GetId() string {
	return "info"
}

func (getInformationCommand) Handle() (string, error) {
	return fmt.Sprintf("Version: %s\nBuild Time: %s", metadata.Version, metadata.BuildTime), nil
}
