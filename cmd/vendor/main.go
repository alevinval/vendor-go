package main

import (
	"github.com/alevinval/vendor-go/pkg/cmd"
	"github.com/alevinval/vendor-go/pkg/govendor"
)

func main() {
	cmd.NewVendorCmd("vendor", &govendor.DefaultPreset{}).Execute()
}
