// SPDX-License-Identifier: GPL-2.0

#include "common.h"

// Context map: the entry program writes its u64 payload here before
// tail-calling into the trampoline.  Rewritten at load time to
// CollectTrampolineRef.CtxMap (named "demo_probe_value" by the tracer).
struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, __u64);
} demo_probe_value SEC(".maps");

// Single-entry prog array populated at load time with the kernel program ID
// of the collect trampoline (CollectTrampolineRef.TailCallDestinationID).
struct {
    __uint(type, BPF_MAP_TYPE_PROG_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, __u32);
} demo_trampoline SEC(".maps");

// modulo is set at load time; 0 disables modulo reduction.
volatile const __u32 modulo = 0;

SEC("kprobe/demoprobe")
int demoprobe(struct pt_regs *ctx)
{
    __u32 key = 0;
    __u64 val = (__u64)bpf_get_prandom_u32();
    if (modulo != 0)
        val = val % modulo;
    bpf_map_update_elem(&demo_probe_value, &key, &val, BPF_ANY);
    bpf_tail_call(ctx, &demo_trampoline, 0);
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
