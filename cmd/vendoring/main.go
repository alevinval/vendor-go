package main

import (
	"github.com/alevinval/vendor-go/pkg/cmd"
	"github.com/alevinval/vendor-go/pkg/vendoring"
)

func main() {
	cmd.NewVendorCmd("vendoring", &vendoring.DefaultPreset{}).Execute()
}
