package core

import (
	"github.com/gocroot/core/config"
	"github.com/gocroot/core/routes"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	config.SetEnv()
	functions.HTTP("WebHook", routes.HandleRoutes)
}
