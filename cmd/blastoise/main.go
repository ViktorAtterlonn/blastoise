package main

import (
	"blastoise/internal/commands"
)

func main() {
	cmd := commands.NewCommand()

	cmd.Run()
}
