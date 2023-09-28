package main

import (
	"authx/model"
	"authx/server"
	"log"

	"github.com/joho/godotenv"
)

func init(){
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading env file")
	}
}

func main() {
	model.Setup()
	server.SetupAndListen()
}
