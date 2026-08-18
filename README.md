# demo-probe

A demonstration of the external probe extension API for the
[OpenTelemetry eBPF profiler](https://github.com/open-telemetry/opentelemetry-ebpf-profiler),
built as a companion to [PR #1739](https://github.com/open-telemetry/opentelemetry-ebpf-profiler/pull/1739).

## What it does

`demo-probe` is an OpenTelemetry Collector extension that:

1. Attaches a kprobe to `sys_getrandom` using the profiler's `RegisterCollectTrampoline` API.
2. On each invocation, generates a random `u64` value (optionally reduced modulo a configurable
   divisor) and writes it into the trampoline context map.
3. Tail-calls into the profiler's collect trampoline, which ships the value as a profile sample.

The result is a stream of profile events emitted every time the kernel's `getrandom` syscall fires.

## Configuration

```yaml
extensions:
  demoprobe:
    modulo: 0   # 0 = disabled; N > 0 reduces each sample value to [0, N)
```

See `config/collector.yaml` for a complete working example.

## Building

### Prerequisites

- Go 1.25+
- [clang](https://clang.llvm.org/) / [llvm](https://llvm.org/) for eBPF compilation
- [bpf2go](https://github.com/cilium/ebpf/tree/main/cmd/bpf2go) (pulled automatically via `go generate`)
- [ocb](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder) — the OTel Collector Builder

### Steps

```sh
# 1. Regenerate the eBPF bindings (only needed after editing bpf/demoprobe.ebpf.c)
go generate ./...

# 2. Build a custom collector binary with the probe bundled in
ocb --config manifest.yaml
```

The resulting binary is written to `./dist/otelcol-demo-probe`.

### Running

```sh
sudo ./dist/otelcol-demo-probe \
  --config=config/collector.yaml \
  --feature-gates=+service.profilesSupport
```

## License

- Go source files: Apache License 2.0 (see `LICENSE`)
- `bpf/demoprobe.ebpf.c` and compiled `.o` files: GPL-2.0
