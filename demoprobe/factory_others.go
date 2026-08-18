//go:build !(linux && (amd64 || arm64))

package demoprobe // import "github.com/florianl/demo-probe/demoprobe"

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

// NewFactory returns an extension.Factory for the demoprobe extension.
// The demoprobe extension is only functional on Linux amd64/arm64.
func NewFactory() extension.Factory {
	return extension.NewFactory(
		Type,
		func() component.Config { return &Config{} },
		func(_ context.Context, _ extension.Settings, _ component.Config) (extension.Extension, error) {
			return nil, errors.New("demoprobe extension is only supported on Linux amd64/arm64")
		},
		stability,
	)
}
