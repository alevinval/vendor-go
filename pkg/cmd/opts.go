package cmd

import "github.com/alevinval/vendor-go/pkg/vending"

// Option is used to apply customizations to the cmdBuilder.
type Option = func(cb *builder)

// WithCommandName option customizes the name of the command.
func WithCommandName(name string) Option {
	return Option(
		func(b *builder) {
			b.commandName = name
		},
	)
}

// WithPreset is used customizes the Preset that will be used.
func WithPreset(preset vending.Preset) Option {
	return Option(
		func(b *builder) {
			b.preset = preset
		},
	)
}
