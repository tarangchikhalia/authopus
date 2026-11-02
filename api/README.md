# A03 — Authopus API Contracts

**Status:** A05-reconciled specification; repository commit and automated verification pending.  
**Target:** November 9, 2026  
**Version:** `v0.1.0` (experimental)

## Files
- `openapi.yaml` — OpenAPI 3.1 contract, endpoints, typed schemas and error responses.
- `*-request.json` / `*-response.json` — concrete JSON examples.
- `contract-tests.md` — acceptance scenarios and invariant checks for automated implementation.

## Confirmed decisions
1. `POST /v1/authorize` is for trusted integrations, not arbitrary agents.
2. The bearer token authenticates the caller; typed request context describes agent, principal, delegation, action and resource, and must be verified against trusted evidence.
3. Four outcomes exist: `ALLOW`, `DENY`, `APPROVAL_REQUIRED`, `ERROR`. Only the first three are HTTP 200; `ERROR` is non-2xx.
4. Approvals are one-time, bound to a specific action, and revalidated at execution.
5. Decision retrieval is single-ID lookup; no decision search API.

## Endpoint scope
- `POST /v1/authorize`
- `GET /v1/decisions/{id}`
- `GET /v1/approvals`
- `GET /v1/approvals/{id}`
- `POST /v1/approvals/{id}/approve`
- `POST /v1/approvals/{id}/reject`
- `GET /healthz`, `GET /readyz`

## A05 reconciliation (contract revision)

- `agent_identity.method` is either `agent_oidc` (requires separately verified agent OIDC token) or `integration_assertion` (caller must be explicitly authorized to assert that agent identity). Both bind to `agent.id` and issuer.
- `authority` is a required discriminated union: `delegated` requires `principal` and `delegation_token`; `self` requires `principal` bound to the verified agent/workload and forbids delegation-token substitution.
- `delegation_token` is signed evidence; a bare `evidence_ref`, grant ID, or caller-supplied scope is not evidence. Signed-token claim profile and key management are finalized in A09.
- The `Decision` success schema excludes `ERROR`; errors use `ErrorResponse` on non-2xx.
- Never log, persist or echo raw OIDC or delegation tokens. Bind authority and request fingerprint at approval time, then revalidate at execution.
- The examples use placeholder tokens and are schema examples, not executable credentials.

## Security and architecture rules

The embedded MCP gateway invokes the Authorization module through an internal Go interface, rather than making an unnecessary network hop. The REST endpoint exposes the same decision contract to trusted integrations. The native Policy module evaluates YAML rules. The Audit & Approval module persists decision evidence and approval state in PostgreSQL. All critical audit persistence failures block execution. Git policy loading, OIDC identity, Docker Compose and single-organization deployment remain in scope. MCP resources/prompts, additional enforcement adapters, dashboards, and audit search remain excluded.

The API does not claim to prevent bypass outside configured enforcement points. Clients must not equate an authorization decision with observed execution. Approval does not itself execute an MCP tool.

## A03 acceptance and definition of done

- [x] Typed REST requests and responses cover allow, deny, approval required and errors.
- [x] OpenAPI 3.1 YAML authored.
- [x] Request/response examples authored.
- [x] Contract test cases authored.
- [ ] Schema and behavior tests implemented and passing in CI (A14/A22).
- [ ] Files committed to the Authopus repository and linked from the roadmap.

**Implementation boundary:** Exact MCP transport, OIDC claim mapping, delegation JWT profile and key management, policy operator set, approval action fingerprint, and storage transactions are detailed engineering tasks, not new release-scope decisions.
