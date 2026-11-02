# ADR-003: Build an Authopus-Native Policy Evaluation Engine

- **Status:** Accepted
- **Date:** 2026-11-02
- **Decision owners:** Tarang Chikhalia
- **Scope:** Authorization policy evaluation

## Context

Platform must decide whether AI-agent actions are authorized using trusted identity, delegation, requested action, resource, and execution context. The **Guard** module is the policy decision point; **Codex** manages policy definition, validation, versioning, and activation; **Gate** enforces decisions. The initial policy representation is planned as versioned YAML, with exact syntax and semantics to be specified.

A general-purpose external policy engine could accelerate development, but it would introduce an additional dependency and potentially constrain policy semantics, explanation formats, evaluation behavior, and the integration between delegation and authorization. Authopus's goal is a focused, comprehensible policy engine implementing only the capabilities its security model needs.

## Decision

**Implement a native policy evaluation engine in Go within Guard**, rather than embedding OPA or depending on another external evaluation engine.

The primary reasons are:

1. **Full control over policy semantics:** Authopus defines its own typed operators, conflict resolution, missing-attribute behavior, and decision model.
2. **Full control over security behavior:** Deny-by-default, explicit-deny precedence, deterministic evaluation, and fail-closed handling are designed and tested as core invariants.
3. **Full control over explainability:** Decisions can expose structured reasons and policy-version references aligned with Watchtower's audit model.
4. **Build only what is needed:** Keep the language and engine narrowly scoped rather than adopting an expansive general-purpose policy framework.
5. **Avoid an external evaluation-engine dependency:** Reduce integration coupling and retain control over the engine's lifecycle and compatibility.

This decision does **not** imply that implementing an engine is inherently safer or easier than using a mature one. It makes Authopus responsible for correctness, security review, language design, and long-term maintenance.

## Alternatives considered

| Alternative | Benefits | Reasons not selected |
|---|---|---|
| **Native Go engine (chosen)** | Direct control of semantics, explainability, integration, and scope | Requires significant testing, language design, and ongoing maintenance |
| OPA / Rego | Mature general-purpose engine and policy ecosystem | Adds an engine dependency and policy model that Authopus does not need to inherit |
| Cedar | Authorization-focused language and formalized semantics | Requires adopting its language/model and integration assumptions |
| Casbin | Ready-made access-control models and enforcement APIs | May not match Authopus's intended delegation-aware policy semantics and explanations |
| Hard-coded Go authorization rules | Minimal initial machinery | Poor policy lifecycle, portability, and operator-managed flexibility |

## Tradeoffs

- **Accept:** Building a parser, validator, evaluator, and explanation model; maintaining backwards compatibility; defending against malformed or adversarial policies; proving evaluation semantics through tests.
- **Gain:** A focused language, stable Authopus-specific decision contracts, tailored explanations, and direct control over performance and security properties.
- **Mitigate:** Keep the initial policy language deliberately small; define a normative specification; use typed AST or equivalent validated representation; separate parse/validation from evaluation; enforce evaluation resource limits; build comprehensive negative tests, fuzz tests, and regression suites.
- **Security constraint:** Unknown or missing attributes, invalid policies, evaluation failures, and time/resource exhaustion must never result in implicit allow. Define precise treatment of non-applicable rules and errors before release.

## Initial design boundaries

- **Policy authoring:** Versioned YAML is the planned input format, not an unrestricted general-purpose programming language.
- **Core evaluation:** Deterministic, strongly typed, contextual evaluation with explicit deny precedence and default deny.
- **Decision output:** Allow/deny with machine-readable explanations and applicable policy identifiers/versions; the final decision schema remains to be designed.
- **Separation of duties:** Codex validates and activates policy versions; Guard evaluates active policies; Gate enforces; Watchtower records evidence.
- **Not yet decided:** Exact policy grammar, AST representation, compilation strategy, caching, policy indexing, persistence format, and performance budgets.

## Consequences

1. Guard owns the evaluation engine and its security-critical semantics.
2. Codex and Guard must share a clearly versioned policy contract without collapsing module responsibilities.
3. The project must publish a policy-language specification and compatibility strategy before stable release.
4. Testing must cover operator semantics, precedence, missing values, type errors, policy-version consistency, resource limits, and adversarial inputs.
5. Audit records should identify the effective policy version and provide enough explanation to reproduce or investigate decisions, while protecting sensitive context.
6. External policy engines are not part of the planned initial architecture.

## Reconsideration criteria

Revisit if the native engine's security or maintenance burden exceeds the project's capacity; policy needs become broad enough that a mature external engine provides a demonstrably better fit; formal assurance requirements cannot be met economically; or measured interoperability and performance needs justify an alternative. Any replacement must preserve Authopus's security invariants, explainability, and policy compatibility strategy.

## Related decisions

- [ADR-001: Deploy Authopus as a modular monolith](ADR-001-modular-monolith.md)
- [ADR-002: Go as the primary implementation language](ADR-002-go-primary-language.md)
