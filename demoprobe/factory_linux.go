//go:build linux && (amd64 || arm64)

package demoprobe // import "github.com/florianl/demo-probe/demoprobe"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

// NewFactory returns an extension.Factory for the demoprobe extension.
func NewFactory() extension.Factory {
	return extension.NewFactory(
		Type,
		func() component.Config { return &Config{} },
		createExtension,
		stability,
	)
}

func createExtension(_ context.Context, _ extension.Settings, cfg component.Config) (extension.Extension, error) {
	c := cfg.(*Config)
	return &demoProbeExtension{modulo: c.Modulo}, nil
}
