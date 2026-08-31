//go:build linux && (amd64 || arm64)

package demoprobe // import "github.com/florianl/demo-probe/demoprobe"

import (
	"context"
	"fmt"

	cebpf "github.com/cilium/ebpf"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/ebpf-profiler/reporter/samples"
	"go.opentelemetry.io/ebpf-profiler/tracer"
)

// demoProbeExtension implements tracer.ProbeExtension.
type demoProbeExtension struct {
	modulo uint32
}

func (e *demoProbeExtension) Start(_ context.Context, _ component.Host) error { return nil }
func (e *demoProbeExtension) Shutdown(_ context.Context) error                { return nil }

// Probe returns the extension itself as the tracer.Probe.
func (e *demoProbeExtension) Probe() tracer.Probe { return e }

func (e *demoProbeExtension) SampleType() *samples.TypeMetadata {
	return &samples.TypeMetadata{SampleType: "events", SampleUnit: "count"}
}

// Load implements tracer.Probe.
func (e *demoProbeExtension) Load(_ context.Context, probeCtx *tracer.ProbeContext) error {
	spec, err := loadDemoprobe()
	if err != nil {
		return fmt.Errorf("loading demoprobe spec: %w", err)
	}

	if v, ok := spec.Variables["modulo"]; ok {
		if err := v.Set(e.modulo); err != nil {
			return fmt.Errorf("setting modulo variable: %w", err)
		}
	}

	tailCallMap, err := probeCtx.WireTrampoline(spec, demoprobeMapDemoProbeValue, demoprobeMapDemoTrampoline)
	if err != nil {
		return err
	}
	defer tailCallMap.Close()

	coll, err := cebpf.NewCollection(spec)
	if err != nil {
		return fmt.Errorf("loading demoprobe collection: %w", err)
	}
	defer coll.Close()

	prog, ok := coll.Programs[demoprobeProgDemoprobe]
	if !ok {
		return fmt.Errorf("program %q not found after loading", demoprobeProgDemoprobe)
	}

	lnk, err := tracer.AttachProbe(prog, &tracer.ProbeSpec{
		Mode:   tracer.ProbeModeKprobe,
		Symbol: "sys_getrandom",
	})
	if err != nil {
		return fmt.Errorf("attaching probe: %w", err)
	}
	probeCtx.AddLink(lnk)

	return nil
}
