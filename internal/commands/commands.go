package commands

import (
	"blastoise/internal/services"

	"github.com/spf13/cobra"
)

type Command struct {
	command *cobra.Command
	service *services.Service
}

type Ctx struct {
	url          string
	rps          int
	duration     int
	method       string
	body         string
	service      *services.Service
	requestsChan chan []*services.RequestResult
}

func NewCommand(s *services.Service, channel chan []*services.RequestResult) *Command {
	var rps int
	var duration int
	var method string
	var body string

	var command = &cobra.Command{
		Use:   "blastoise",
		Args:  cobra.ExactArgs(1),
		Short: "Blastoise is a CLI tool for initiating processes",
		Run: func(cmd *cobra.Command, args []string) {

			ctx := Ctx{
				url:          args[0],
				rps:          rps,
				duration:     duration,
				method:       method,
				body:         body,
				service:      s,
				requestsChan: channel,
			}

			model := NewModel(ctx)

			model.Start()
		},
	}

	command.PersistentFlags().IntVarP(&rps, "rps", "r", 1, "Set the number of requests per second")
	command.PersistentFlags().IntVarP(&duration, "duration", "d", 10, "Set the duration in seconds")
	command.PersistentFlags().StringVarP(&method, "method", "m", "GET", "Set the HTTP method")
	command.PersistentFlags().StringVarP(&body, "body", "b", "", "Set the HTTP body")

	return &Command{
		command: command,
		service: s,
	}
}

func (c *Command) Run() error {
	return c.command.Execute()
}
