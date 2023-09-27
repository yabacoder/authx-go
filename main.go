package main

import (
	"authx/model"
	"authx/server"
)

func main() {
	model.Setup()
	server.SetupAndListen()
}
