package demoprobe

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target amd64,arm64 demoprobe ../bpf/demoprobe.ebpf.c -- -I../bpf/headers
