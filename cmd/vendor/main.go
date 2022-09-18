package main

import (
	"github.com/alevinval/vendor-go/pkg/cmd"
	"github.com/alevinval/vendor-go/pkg/vendor"
)

func main() {
	cmd.NewVendorCmd("vendor", &vendor.DefaultPreset{}).Execute()
}
