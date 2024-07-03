package app

import (
	"main/internal/server"
	"main/pkg/database/mongodb"
)

func Run(configPath string) {

	mongodb.ConnectToMongo()

	server.HandleRequest()
}
