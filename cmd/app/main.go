package main

import "main/internal/app"

const configPath = "configs/.env"

func main() {
	app.Run(configPath)
}
