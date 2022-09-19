package main

import (
	"github.com/alevinval/vendor-go/pkg/cmd"
	"github.com/alevinval/vendor-go/pkg/vending"
)

func main() {
	cmd.NewVendingCmd("vending", &vending.DefaultPreset{}).Execute()
}
