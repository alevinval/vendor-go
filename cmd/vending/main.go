package main

import (
	"github.com/alevinval/vendor-go/pkg/cmd"
	"github.com/alevinval/vendor-go/pkg/vending"
)

func main() {
	cmd.NewCobraCommand(
		cmd.WithCommandName("vending"),
		cmd.WithPreset(&vending.DefaultPreset{}),
	).Execute()
}
