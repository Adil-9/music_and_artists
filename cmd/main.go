package main

import (
	"groupie-tracker/internal/cache"
	"groupie-tracker/internal/server"
	"groupie-tracker/logger"
)

func main() {

	logger.Init()

	cache.RedisClientInit()

	server.StartConfigureServer()
}
