package main

import (
	"groupie-tracker/internal/server"
	"groupie-tracker/logger"
)

func main() {

	logger.Init()

	server.StartConfigureServer()
}
