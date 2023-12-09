package main

import (
	"blastoise/internal/commands"
	"blastoise/internal/services"
)

func main() {

	resultchn := make(chan []*services.RequestResult)

	service := services.NewService()

	cmd := commands.NewCommand(service, resultchn)

	cmd.Run()
}
