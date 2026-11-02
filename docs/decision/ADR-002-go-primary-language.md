# ADR-002: Use Go as Authopus's Primary Implementation Language

- **Status:** Accepted
- **Date:** 2026-11-02
- **Decision owners:** Tarang Chikhalia
- **Scope:** Core service

## Context

Authopus is a security layer that handles many network-oriented requests, including MCP interception, identity and delegation checks, policy evaluation, and audit event production. It needs predictable request handling, good concurrency support, practical deployment, clear code, and maintainable security-critical behavior. The project is open source, so contributor accessibility and operational simplicity matter alongside performance.

Go and Rust are both capable of implementing secure, performant network services. The choice is a fit-for-purpose engineering judgment, **not** a claim that Rust cannot serve this workload or that Go is universally faster or safer.

## Decision

Use **Go as the primary language** for Authopus's single deployable service and its native policy evaluation engine.

The key reasons are:

- **Network-service fit:** Go's standard networking packages, goroutines, context cancellation, and mature server ecosystem suit a high-concurrency, I/O-heavy authorization gateway.
- **Developer productivity:** A relatively small language, consistent tooling, and straightforward builds support rapid, reviewable development.
- **Maintainability:** Clear conventions and accessible code help contributors reason about security-sensitive flows.
- **Operational efficiency:** Compiled binaries, built-in profiling, and a mature observability ecosystem support straightforward deployment and performance diagnosis.

Rust is particularly attractive for low-level systems programming, memory-sensitive components, and tight control over resource usage. Those advantages do not currently outweigh Go's development and operations fit for primarily network-facing service. Rust is also viable for network services; this decision reflects project priorities rather than an inherent limitation of Rust.

## Alternatives considered

| Alternative | Strengths | Reasons not selected as primary language |
|---|---|---|
| **Go (chosen)** | Strong networking and concurrency support, simple builds, approachable maintenance, useful profiling tools | Garbage collection and runtime overhead; less fine-grained memory control |
| Rust | Memory safety without garbage collection, fine-grained resource control, strong performance | More implementation complexity and potentially slower iteration for the current network-service-focused scope |
| Java / Kotlin | Mature server ecosystem, concurrency options, enterprise integrations | More runtime/deployment considerations for the desired simple service distribution |
| TypeScript / Node.js | Rapid development, broad contributor base, strong I/O ecosystem | CPU-heavy policy evaluation needs careful event-loop isolation; dynamic-language concerns unless rigorously controlled |
| C / C++ | Fine-grained performance and systems control | Higher memory-safety burden and maintenance risk for a security-critical open-source service |

## Tradeoffs

- **Accept:** Garbage-collector pauses, runtime overhead, and a need to profile allocation-heavy policy evaluation paths.
- **Gain:** Fast development cycles, simple tooling, strong concurrency primitives, straightforward deployment, and maintainable code.
- **Mitigate:** Benchmark policy evaluation and end-to-end enforcement separately; use load tests, race detection, fuzzing, profiling, and bounded concurrency; avoid claiming performance without measurement.
- **Security constraint:** Language choice alone does not establish security. Input validation, trust-boundary enforcement, correct concurrency, dependency review, and adversarial testing remain mandatory.

## Consequences

1. Core service, modules, and policy engine are implemented in idiomatic Go.
2. Go toolchain conventions govern formatting, testing, linting, profiling, and CI.
3. APIs and module interfaces should remain language-agnostic where useful for external integrations, without requiring a polyglot core.
4. Performance work is guided by measured latency, throughput, allocations, and resource consumption.
5. A different language may be considered for a narrowly scoped component only with clear evidence and an explicit architectural decision.

## Reconsideration criteria

Revisit if profiling shows that Go runtime characteristics prevent meeting defined latency or memory targets after reasonable optimization; a new component requires low-level constraints that Go cannot adequately satisfy; maintenance or ecosystem conditions materially change; or a measured alternative offers substantial benefit without unacceptable operational complexity.

## Related decisions

- [ADR-001: Deploy Authopus as a modular monolith](ADR-001-modular-monolith.md)
- [ADR-003: Build a native policy evaluation engine](ADR-003-native-policy-engine.md)
