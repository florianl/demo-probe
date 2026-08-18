package demoprobe // import "github.com/florianl/demo-probe/demoprobe"

import "go.opentelemetry.io/collector/component"

var (
	Type      = component.MustNewType("demoprobe")
	stability = component.StabilityLevelDevelopment
)
