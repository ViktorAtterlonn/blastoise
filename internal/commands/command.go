package commands

import (
	"blastoise/internal/runner"
	"blastoise/internal/structs"
	"blastoise/internal/view"

	"github.com/spf13/cobra"
)

type Command struct {
	command *cobra.Command
}

func NewCommand() *Command {
	var rps int
	var duration int
	var method string
	var body string

	var command = &cobra.Command{
		Use:   "blastoise",
		Args:  cobra.ExactArgs(1),
		Short: "Blastoise is a CLI tool for initiating processes",
		Run: func(cmd *cobra.Command, args []string) {

			abortchn := make(chan bool)
			resultchn := make(chan []*structs.RequestResult)

			ctx := structs.Ctx{
				Url:        args[0],
				Rps:        rps,
				Duration:   duration,
				Method:     method,
				Body:       body,
				ResultChan: resultchn,
				AbortChan:  abortchn,
			}

			view := view.NewView(&ctx)
			runner := runner.NewHttpRequestRunner(&ctx)

			go runner.Run()
			view.Start()
		},
	}

	command.PersistentFlags().IntVarP(&rps, "rps", "r", 1, "Set the number of requests per second")
	command.PersistentFlags().IntVarP(&duration, "duration", "d", 10, "Set the duration in seconds")
	command.PersistentFlags().StringVarP(&method, "method", "m", "GET", "Set the HTTP method")
	command.PersistentFlags().StringVarP(&body, "body", "b", "", "Set the HTTP body")

	return &Command{
		command: command,
	}
}

func (c *Command) Run() error {
	return c.command.Execute()
}
