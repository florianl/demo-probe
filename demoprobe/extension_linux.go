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
	modulo      uint32
	ctxMap      *cebpf.Map
	tailCallMap *cebpf.Map
}

func (e *demoProbeExtension) Start(_ context.Context, _ component.Host) error { return nil }
func (e *demoProbeExtension) Shutdown(_ context.Context) error {
	if e.ctxMap != nil {
		e.ctxMap.Close()
	}
	if e.tailCallMap != nil {
		e.tailCallMap.Close()
	}
	return nil
}

// Probe returns the extension itself as the tracer.Probe.
func (e *demoProbeExtension) Probe() tracer.Probe { return e }

// Load implements tracer.Probe.
func (e *demoProbeExtension) Load(_ context.Context, _ tracer.ProbeRegistrar, probeCtx *tracer.ProbeContext) error {
	ref, err := probeCtx.RegisterCollectTrampoline(&samples.TypeMetadata{
		SampleType: "events",
		SampleUnit: "count",
	})
	if err != nil {
		return fmt.Errorf("registering collect trampoline: %w", err)
	}

	spec, err := loadDemoprobe()
	if err != nil {
		ref.CtxMap.Close()
		return fmt.Errorf("loading demoprobe spec: %w", err)
	}

	if v, ok := spec.Variables["modulo"]; ok {
		if err := v.Set(e.modulo); err != nil {
			ref.CtxMap.Close()
			return fmt.Errorf("setting modulo variable: %w", err)
		}
	}

	tailCallMap, err := cebpf.NewMap(spec.Maps[demoprobeMapDemoTrampoline])
	if err != nil {
		ref.CtxMap.Close()
		return fmt.Errorf("creating demo_trampoline: %w", err)
	}

	trampolineProg, err := cebpf.NewProgramFromID(cebpf.ProgramID(ref.TailCallDestinationID))
	if err != nil {
		tailCallMap.Close()
		ref.CtxMap.Close()
		return fmt.Errorf("opening collect trampoline program: %w", err)
	}
	defer trampolineProg.Close()

	if err := tailCallMap.Put(uint32(0), trampolineProg); err != nil {
		tailCallMap.Close()
		ref.CtxMap.Close()
		return fmt.Errorf("populating demo_trampoline[0]: %w", err)
	}

	coll, err := cebpf.NewCollectionWithOptions(spec, cebpf.CollectionOptions{
		MapReplacements: map[string]*cebpf.Map{
			demoprobeMapDemoProbeValue: ref.CtxMap,
			demoprobeMapDemoTrampoline: tailCallMap,
		},
	})
	if err != nil {
		tailCallMap.Close()
		ref.CtxMap.Close()
		return fmt.Errorf("loading demoprobe collection: %w", err)
	}
	defer coll.Close()

	prog, ok := coll.Programs[demoprobeProgDemoprobe]
	if !ok {
		tailCallMap.Close()
		ref.CtxMap.Close()
		return fmt.Errorf("program %q not found after loading", demoprobeProgDemoprobe)
	}

	lnk, err := tracer.AttachProbe(prog, &tracer.ProbeSpec{
		Mode:   tracer.ProbeModeKprobe,
		Symbol: "sys_getrandom",
	})
	if err != nil {
		tailCallMap.Close()
		ref.CtxMap.Close()
		return fmt.Errorf("attaching probe: %w", err)
	}
	probeCtx.AddLink(lnk)

	e.ctxMap = ref.CtxMap
	e.tailCallMap = tailCallMap
	return nil
}
