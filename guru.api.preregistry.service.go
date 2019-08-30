package main

import (
	"github.com/guru-invest/guru.api.preregistry/api"
	global "github.com/guru-invest/guru.api.preregistry/configuration"
)

func main() {
	global.CONFIGURATION.Load()
	api.InitializeApi()
}
