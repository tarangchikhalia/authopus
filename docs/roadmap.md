# Authryn — Development Roadmap

> **Status:** Proposed roadmap · **Last updated:** 2026-11-02
>
> This document describes intended development sequencing, not completed functionality or committed delivery dates. Milestones and priorities may change through documented project decisions.

## Vision

Build Authopus as an open-source, production-oriented security layer that makes AI agent identity, delegated authority, fine-grained authorization, runtime enforcement, and accountability explicit and auditable. Start with Model Context Protocol (MCP) interactions and retain architectural flexibility for other integration points.

## Architectural commitments

- **Deployment:** One service with five logically separated **modules**.
- **Gate:** MCP interception and enforcement (policy enforcement point).
- **Delegate:** Agent identity and delegated authority.
- **Guard:** Native policy evaluation and authorization decisions (policy decision point).
- **Codex:** Policy definition, validation, versioning, and activation.
- **Watchtower:** Decision and enforcement auditability.
- **Policies:** Versioned YAML is the initial authoring format; the detailed schema is not yet final.
- **Persistence:** PostgreSQL is planned; concrete schemas and operational choices remain design work.
- **Security:** Default deny, fail closed, least privilege, explicit-deny precedence, trusted identity context, tenant isolation, and explainable decisions.

## Delivery strategy

Deliver **vertical slices** that demonstrate an end-to-end capability. Avoid developing all five modules in isolation before the first working integration. Use milestones as exit gates, not calendar promises.

**Suggested release boundaries:**

- **Prototype:** A local, demonstrable MCP request that receives and enforces a native policy decision.
- **MVP / first usable alpha:** A documented and tested MCP integration with identity/delegation context, policy lifecycle basics, and auditable decisions.
- **Public beta:** Stronger security controls, operational readiness, compatibility testing, and contributor documentation.
- **Stable release:** Explicit reliability, security, compatibility, and maintenance criteria are met; no date is implied.

## Milestone 0 — Foundation and security model

**Goal:** Establish a coherent architecture and shared security contract before implementing the critical path.

**Work:**

- [ ] Establish repository layout, CI, formatting, linting, testing, and dependency update practices.
- [ ] Document module boundaries, dependency direction, and internal interfaces.
- [ ] Define the authorization request and decision data models (subject, actor, principal, action, resource, context, decision, reason, policy version, correlation ID).
- [ ] Specify trust boundaries and the threat model, including spoofed agent identity, confused deputy, privilege escalation, policy bypass, replay, and audit tampering.
- [ ] Define deny-by-default, error handling, missing-attribute semantics, and policy conflict precedence.
- [ ] Choose an initial MCP transport/interception strategy and document compatibility constraints.
- [ ] Record unresolved architectural decisions as ADRs.
- [ ] Decide on a suitable open-source license before accepting outside contributions.

**Exit criteria:** An architecture document, threat model, authorization contract, repository checks, and an executable minimal service skeleton exist. Security-critical assumptions are explicit.

## Milestone 1 — Guard: native policy engine

**Goal:** Implement the smallest useful deterministic authorization engine.

**Work:**

- [ ] Specify the initial versioned YAML policy schema and its validation rules.
- [ ] Implement policy parsing, schema validation, and typed internal representation.
- [ ] Implement policy matching for subject, action, resource, and contextual attributes.
- [ ] Implement logical composition and typed comparisons, with precise missing/unknown attribute behavior.
- [ ] Implement explicit-deny precedence, default deny, and fail-closed evaluation errors.
- [ ] Return structured decisions with reason codes and policy/version references.
- [ ] Add policy test fixtures, unit tests, property/fuzz tests, and baseline benchmarks.
- [ ] Document evaluation semantics independently of the implementation.

**Exit criteria:** Identical validated inputs and policy versions produce identical results; malformed policies cannot become active; deny/allow conflicts and evaluation failures have comprehensive tests. The evaluator is callable through an internal Go interface.

## Milestone 2 — Gate: first end-to-end MCP enforcement

**Goal:** Demonstrate a real authorization boundary around a supported MCP operation.

**Work:**

- [ ] Implement the chosen MCP integration pattern and request normalization.
- [ ] Extract action, resource, and trusted request context without treating unverified agent-supplied claims as authority.
- [ ] Call Guard for each protected operation and enforce the result.
- [ ] Block unauthorized actions **before** forwarding to the MCP server.
- [ ] Preserve supported protocol behavior, correlation IDs, error semantics, and response forwarding.
- [ ] Implement fail-closed behavior when Guard or required context is unavailable.
- [ ] Build an example MCP client/server integration and end-to-end tests.

**Exit criteria:** A local demonstration proves that a permitted request reaches a test MCP tool, a denied request never reaches it, and evaluation failures cannot bypass enforcement. Supported transports and MCP operations are documented.

**Release gate:** **Prototype** candidate.

## Milestone 3 — Delegate: identity and constrained authority

**Goal:** Establish who is acting and what authority has actually been delegated.

**Work:**

- [ ] Define agent identity, human/service principal, and delegation data models.
- [ ] Select trusted identity inputs and authentication integration boundaries.
- [ ] Represent delegation scope, issuer, audience, validity period, and revocation state.
- [ ] Validate delegation chains and prevent privilege amplification.
- [ ] Integrate trusted identity and delegation context into Gate → Guard decisions.
- [ ] Test spoofed claims, expired/revoked delegation, cross-tenant requests, and confused-deputy cases.

**Exit criteria:** An agent cannot gain authority merely by asserting identity or delegation fields; authorization respects the authenticated principal and constrained delegation. Invalid, expired, or revoked authority is denied.

## Milestone 4 — Codex: managed policy lifecycle

**Goal:** Make policy changes safe, traceable, and repeatable.

**Work:**

- [ ] Implement policy create/import, validate, test, version, and activate operations.
- [ ] Define atomic activation and rollback behavior.
- [ ] Persist policy metadata and active-version references in PostgreSQL as needed.
- [ ] Establish policy ownership, tenant boundaries, and change authorization.
- [ ] Provide a documented local workflow and a minimal management interface or CLI; select the exact interface deliberately.
- [ ] Test invalid updates, concurrent activation, rollback, and isolation between tenants.

**Exit criteria:** Operators can validate a candidate policy, activate a specific version, verify which version was used for a decision, and roll back without exposing an invalid or partially activated policy.

## Milestone 5 — Watchtower: audit and accountability

**Goal:** Capture trustworthy, useful evidence about authorization and enforcement.

**Work:**

- [ ] Define a structured audit event schema and correlation strategy.
- [ ] Record decision, policy version, identity/delegation references, and enforcement outcome.
- [ ] Clearly distinguish requested, authorized, denied, forwarded, executed, and observed states; do not claim downstream execution without evidence.
- [ ] Persist events and define retention, access controls, and sensitive-data redaction rules.
- [ ] Add query/filter capabilities sufficient for investigation.
- [ ] Test event completeness, failure handling, tampering risks, and sensitive-data minimization.

**Exit criteria:** An operator can trace a representative MCP request from trusted identity and delegation through the policy decision to the observable enforcement outcome, without exposing secrets in logs.

**Release gate:** **MVP / first usable alpha** candidate, subject to integration and security review.

## Milestone 6 — Hardening and public beta

**Goal:** Make the integrated system practical for early external adopters.

**Work:**

- [ ] Complete cross-module end-to-end, regression, and adversarial test suites.
- [ ] Define and measure policy evaluation latency and enforcement overhead under representative workloads; set thresholds from evidence, not guesses.
- [ ] Add health checks, structured logs, metrics, tracing, and operational diagnostics.
- [ ] Validate persistence migrations, backup/restore assumptions, configuration, and secret handling.
- [ ] Document supported MCP versions, transports, limitations, and upgrade behavior.
- [ ] Provide reproducible local setup, examples, deployment guidance, and troubleshooting.
- [ ] Publish README, architecture, policy language, threat model, security reporting, and contribution documentation.
- [ ] Perform focused security review and address critical findings before public beta.

**Exit criteria:** A new contributor can run the example locally; an early adopter can configure a supported integration; the project has a clear security reporting path, observable failure modes, and documented compatibility limits.

**Release gate:** **Public beta** candidate.

## Milestone 7 — Stable release readiness

**Goal:** Establish evidence-based readiness for production-oriented adoption.

**Work:**

- [ ] Publish a supported API/policy compatibility and versioning policy.
- [ ] Define upgrade and migration guarantees, release process, and maintenance commitments.
- [ ] Validate recovery, degradation, load behavior, and security controls against documented targets.
- [ ] Resolve release-blocking issues and document known limitations.
- [ ] Complete release artifacts, changelog, provenance/integrity practices, and operator documentation.

**Exit criteria:** Explicit release criteria are met and independently verifiable. Stable status is not assigned based solely on feature completeness.

## Potential later capabilities (not committed)

- Additional agent integration points beyond MCP.
- Richer policy language features, tooling, and authoring UX.
- More identity providers and delegation mechanisms.
- Advanced policy simulation and decision analysis.
- Scalable audit export and external observability integrations.
- Optional deployment topologies when operational evidence justifies them.

These are exploration areas, **not** current MVP commitments. Commercialization is not on the active roadmap.

## Cross-cutting requirements

Every milestone should address:

1. **Security:** Default deny, fail closed, least privilege, trusted inputs, and tenant isolation.
2. **Correctness:** Deterministic, typed policy evaluation and explicit edge-case semantics.
3. **Explainability:** Structured decision reasons and policy version provenance.
4. **Auditability:** Evidence appropriate to the milestone, without implying unobserved outcomes.
5. **Developer experience:** Reproducible examples, clear errors, tests, and documentation.
6. **Modularity:** Preserve boundaries between Gate, Delegate, Guard, Codex, and Watchtower within one Go service.

## Decision log and open questions

The following require design decisions rather than assumptions:

- Exact MCP interception pattern, transports, and protocol-version support.
- Policy schema, expression grammar, evaluation semantics, and obligations model.
- Identity credential types, trust anchors, and delegation representation.
- Tenant model and administrative authorization boundaries.
- PostgreSQL schema, retention, and audit durability requirements.
- Management API versus CLI priorities.
- Performance and availability objectives, deployment targets, and release support policy.
- Open-source license and external contribution governance.

Record material choices in `docs/decisions/` and update this roadmap when a decision changes sequencing or acceptance criteria.

## Milestone tracking

| Milestone | Primary modules | Status | Release relevance |
|---|---|---|---|
| 0. Foundation and security model | All | Planned | Foundation |
| 1. Native policy engine | Guard | Planned | Prototype prerequisite |
| 2. MCP enforcement | Gate + Guard | Planned | Prototype |
| 3. Identity and delegation | Delegate + Gate + Guard | Planned | MVP prerequisite |
| 4. Policy lifecycle | Codex + Guard | Planned | MVP prerequisite |
| 5. Audit and accountability | Watchtower + all | Planned | MVP / alpha |
| 6. Hardening and public beta | All | Planned | Beta |
| 7. Stable release readiness | All | Planned | Stable |

**Tracking convention:** `Planned` → `In progress` → `Blocked` (if applicable) → `Complete`. Update statuses only when supported by repository evidence and verified acceptance criteria.
