package main

import (
	"github.com/guru-invest/guru.api.preregistry/api"
	global "github.com/guru-invest/guru.api.preregistry/configuration"
	l "github.com/guru-invest/guru.api.preregistry/logger"
)

func main() {
	global.CONFIGURATION.Load()
	l.InitLog(global.CONFIGURATION.API.Route, global.CONFIGURATION.API.Route, global.CONFIGURATION.OTHER.LogLevel)
	l.LOG.Info("Initialize application")
	api.InitializeApi()
}
