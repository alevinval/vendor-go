package main

import (
	"github.com/alevinval/vendor-go/pkg/cmd"
	"github.com/alevinval/vendor-go/pkg/vending"
)

func main() {
	cmd.NewVendorCmd("vending", &vending.DefaultPreset{}).Execute()
}
