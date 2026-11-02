# ADR-001: Deploy Authopus as a single service

- **Status:** Accepted
- **Date:** 2026-11-02
- **Decision owners:** Tarang Chikhalia
- **Scope:** Service architecture and deployment topology

## Context

Authopus platform has five logical modules: 

- **Gate** (enforcement), 
- **Delegate** (identity and delegation), 
- **Guard** (authorization and native policy evaluation), 
- **Codex** (policy lifecycle), and 
- **Watchtower** (audit and accountability).

The project needs a reliable authorization path, simple installation for adopters, manageable operations, and clear internal ownership boundaries. Authorization decisions may require interactions among several modules on a single request. Introducing independent services from the outset would add network calls, deployment coordination, distributed failure modes, and operational overhead before there is evidence that independent scaling is necessary.

## Decision

Build Authopus as **one deployable service containing five logically separated modules**, not five microservices. Preserve explicit module interfaces, well-defined ownership, independent tests, and dependency boundaries. Keep internal calls in-process where appropriate. Use a relational database as the planned persistent store, with schema and ownership details determined during design.

The primary motivations are **simpler deployment and operation, lower architectural complexity, and reduced latency on the authorization path**. Modular boundaries should allow a future service extraction if measured requirements justify it.

## Alternatives considered

| Alternative | Benefits | Reasons not selected now |
|---|---|---|
| Five independent microservices | Independent scaling, deployments, and fault isolation | More infrastructure, RPC contracts, distributed tracing, network latency, and failure handling; premature at current scale |
| Modular monolith (chosen) | One deployment, in-process coordination, cohesive development and testing | Shared release cadence and process-level failure domain |
| Unstructured monolith | Fastest initial coding | Erodes module boundaries and makes later evolution difficult |
| Hybrid services from day one | Targeted separation of selected components | No demonstrated operational or scaling need yet |

## Tradeoffs

- **Accept:** A shared process and release lifecycle; one module's resource exhaustion can affect others; independent scaling is unavailable initially.
- **Gain:** Fewer deployment artifacts and operational dependencies; simpler local development; easier cross-module consistency and integration tests; fewer network hops on security-critical request paths.
- **Mitigate:** Enforce interfaces and ownership; avoid cyclic module dependencies; set timeouts and resource budgets; use bounded queues for audit work where appropriate; test failure containment and authorization behavior under load.
- **Security constraint:** A failure or overload in a supporting module must not silently turn an authorization failure into an allow decision.

## Consequences

1. The initial release is packaged and deployed as one single service.
2. Gate, Delegate, Guard, Codex, and Watchtower remain separate logical modules with documented contracts and tests.
3. Internal module communication does not require network APIs solely for architectural separation.
4. Operational observability must still identify module-level latency, errors, and bottlenecks.
5. Any future extraction must preserve security invariants, decision consistency, and audit correlation.

## Reconsideration criteria

Revisit this decision if production evidence shows sustained, materially different scaling requirements; unacceptable process-level blast radius; independent deployment needs that outweigh coordination costs; organizational ownership boundaries requiring separate lifecycles; or regulatory/isolation requirements that cannot be met within one service. Base any change on measured bottlenecks and a migration plan, not hypothetical scale.

## Related decisions

- [ADR-002: Go as the primary implementation language](ADR-002-go-primary-language.md)
- [ADR-003: Build a native policy evaluation engine](ADR-003-native-policy-engine.md)
