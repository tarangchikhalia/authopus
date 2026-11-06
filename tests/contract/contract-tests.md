# A03 Contract Test Cases — Authopus v0.1.0

These are specification-level test cases to implement in automated API and MCP integration tests. Use a real OIDC test issuer, PostgreSQL and an instrumented MCP upstream. Do not treat a mocked response as evidence of enforcement.

| ID | Scenario | Expected HTTP / body | Enforcement and persistence assertion |
|---|---|---|---|
| CT-01 | Trusted caller, valid delegation, matching allow | 200 / ALLOW | Decision durably audited before tool forwarding |
| CT-02 | No matching policy | 200 / DENY | No upstream call; denial audited |
| CT-03 | Explicit deny overrides allow | 200 / DENY | No upstream call; policy reference recorded |
| CT-04 | Approval rule matches | 200 / APPROVAL_REQUIRED with approval_id | Pending approval persisted; no upstream call |
| CT-05 | Missing bearer token | 401 / ERROR | No upstream call |
| CT-06 | Forged, expired, wrong-audience OIDC token | 401 / ERROR | No upstream call |
| CT-07 | Authenticated but unauthorized integration | 403 / ERROR | No upstream call |
| CT-08 | Forged agent, principal or delegation in JSON | 403 or 422 / ERROR | Caller claims not treated as verified; no upstream call |
| CT-09 | Malformed JSON or missing required field | 400 or 422 / ERROR | No upstream call |
| CT-10 | Policy engine error | 500 / ERROR | No upstream call; correlation ID present |
| CT-11 | Authorization service timeout | 503 / ERROR if HTTP response possible | Gateway blocks; no upstream call |
| CT-12 | PostgreSQL unavailable before required audit commit | 503 / ERROR | No upstream call; never return ALLOW |
| CT-13 | Approval requested, unauthorized reviewer | 403 / ERROR | No state change or execution |
| CT-14 | Authorized reviewer approves pending request | 200 / APPROVED | Approval persisted; no immediate execution |
| CT-15 | Approved action executed with unchanged verified context | New authorization/revalidation succeeds | Exactly one execution; approval atomically consumed |
| CT-16 | Reuse same approval | 409 / ERROR or denied execution | No second execution |
| CT-17 | Expired or rejected approval | 410 or 409 / ERROR | No execution |
| CT-18 | Policy, delegation, resource or action changes after approval | Revalidation blocks | No execution despite prior approval |
| CT-19 | Concurrent consumers of one approval | At most one succeeds | At most one upstream execution |
| CT-20 | Decision lookup by authorized integration | 200 / Decision | Returns matching durable record only |
| CT-21 | Decision lookup for another integration scope | 403 or 404 / ERROR | No information disclosure |
| CT-22 | Unknown decision or approval ID | 404 / ERROR | No sensitive details |
| CT-23 | Ready endpoint while PostgreSQL unavailable | 503 / not_ready | No false readiness |
| CT-24 | Audit event includes token or secret | Never acceptable | Assert secret redaction across logs, traces and audit |
| CT-25 | MCP transport disconnect or upstream failure | Gateway reports failure | Audit distinguishes authorization from observed execution outcome |

## Contract invariants

1. `ALLOW`, `DENY` and `APPROVAL_REQUIRED` are valid evaluation results (HTTP 200); `ERROR` is a non-2xx failure outcome.
2. An HTTP 200 `ALLOW` never asserts that a downstream tool actually executed.
3. An approval is bound to a single immutable action fingerprint and is one-time use; the fingerprint and execution handoff are internal contracts to be specified during A19.
4. Verified token identifies the caller; request-body identity/delegation values are not trusted without validation.
5. Required audit persistence is a precondition to tool execution; if the durable write fails, enforcement fails closed.
6. Access to decisions and approvals is scope-checked, even in a single-organization deployment.
7. Error messages do not leak tokens, credentials, sensitive policy internals or unrelated identities.

## Implementation notes

- Status codes for individual validation failures may be tightened in A14; do not alter the success-versus-error distinction.
- Test actual upstream invocation counts using a controllable MCP server.
- Add schema validation of request and response bodies against `openapi.yaml` in CI.
- The OpenAPI contract does not expose an approval-consumption endpoint; the embedded gateway uses a private internal Go interface to consume approvals and revalidate before forwarding.

## A05 identity and delegation reconciliation tests

| ID | Scenario | Expected result | Required assertion |
|---|---|---|---|
| CT-26 | Delegated request with valid signed token, principal/agent binding and scoped action/resource | 200 / policy outcome | Validated issuer, signature, audience, time and bindings before policy evaluation |
| CT-27 | Delegated request with only `evidence_ref` / grant ID | 422 / ERROR | Schema rejects absent signed delegation token; no upstream call |
| CT-28 | Delegated request with missing or expired signed token | 422 or 403 / ERROR | No policy allow, approval or upstream execution |
| CT-29 | Delegation wrong issuer, audience, principal, agent, action or resource | 403 / ERROR | Reject invalid evidence without trusting caller-supplied scopes |
| CT-30 | Self-authority with verified agent/workload principal and explicit policy grant | 200 / policy outcome | No delegation token required; principal bound to verified identity |
| CT-31 | Self-authority with another user's principal | 403 / ERROR | Self mode cannot impersonate a human/service principal |
| CT-32 | Direct agent OIDC proof with matching `iss` and `sub` | 200 / policy outcome | Verify agent token separately from integration bearer |
| CT-33 | Direct agent OIDC proof with forged, expired or mismatched agent token | 401 or 403 / ERROR | Agent identity cannot be overridden by JSON |
| CT-34 | Integration assertion from caller not allowlisted for agent namespace | 403 / ERROR | Reject even when integration bearer itself is valid |
| CT-35 | Missing `authority_mode`, ambiguous authority object or unsupported mode | 422 / ERROR | Never infer self from missing delegation |
| CT-36 | Delegation expires, policy changes or principal authority revoked after approval | Block execution | Revalidate before one-time consumption and forwarding |
| CT-37 | Sensitive tokens in request | No secrets in outputs or storage | Check application logs, traces, audit and error messages |
| CT-38 | Required audit persistence fails for either authority mode | 503 / ERROR | No downstream call and no successful ALLOW response |

**Error taxonomy:** Invalid JSON/schema yields 400/422; missing/invalid caller authentication yields 401; valid caller without assertion permission or invalid delegation authorization yields 403 (implementation may use 422 for structurally invalid delegation). `DENY` is a completed policy/authority decision, not an authentication failure.
