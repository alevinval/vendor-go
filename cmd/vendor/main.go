package main

import (
	"github.com/alevinval/vendor-go/pkg/cli"
	"github.com/alevinval/vendor-go/pkg/govendor"
)

func main() {
	cli.NewVendorCmd("vendor", &govendor.DefaultPreset{}).Execute()
}
