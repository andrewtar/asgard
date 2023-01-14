package service

import (
	"asgard/common/debug/metadata"
	"asgard/common/log"
)

func Init() {
	log.Init().
		WithField("version", metadata.Version).
		WithField("build_time", metadata.BuildTime).
		Info("Build metadata")
}
