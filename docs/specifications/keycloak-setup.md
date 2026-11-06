# Keycloak identity configuration (#11)

Use one realm per deployment. Keycloak manages users and credentials; Authopus owns authorization and assertion permissions. This milestone provides configuration loading only. Token verification is #12; containers, realm imports and real token tests are #13. Neither command loads configuration yet.

## Identity contract

Canonical identity is `(iss, sub)`. API `Actor.id` is the subject ID, never a username or client ID. Omitted API issuer resolves to the configured realm; supplied issuer must match exactly. Internally always retain both fields. Identity references in Go describe identities, not proof of authentication.

Register each agent and integration as a separate confidential client with service accounts enabled. Use client credentials; disable standard, implicit, direct access/password and device grants for these clients. Humans use ordinary user accounts and browser login through an appropriate application client. Human subjects do not need Authopus allowlist entries; delegated principal evidence is validated separately in #5.

Obtain each client's service-account user ID from Keycloak's administrative UI or operator tooling, and verify it matches the issued token's subject. Do not use the client UUID or service-account username. Authopus requires no Admin API credentials and makes no runtime directory lookup. Explicit subject/client bindings classify workloads. Roles and naming conventions cannot change classification.

## Issuer settings

Use an HTTPS realm issuer with an exact, stable public URL. Configure RS256 and a five-minute access-token lifetime. Add the `authopus` audience through an audience mapper for all relevant clients. Enable **Use at+jwt as access token header type**. Ensure issued tokens contain `iss`, `sub`, `aud`, `exp`, `iat`, `jti`, and `client_id`; configure a service-account session-note mapper for `client_id` if needed. Never substitute `azp`. Verify complete tokens against RFC 9068; header configuration alone is insufficient.

References: [Keycloak administration](https://www.keycloak.org/docs/latest/server_admin/index.html), [token header setting](https://www.keycloak.org/2025/04/keycloak-2620-released), [RFC 9068](https://www.rfc-editor.org/rfc/rfc9068.html).

## Authopus configuration

The executable example is `internal/platform/config/testdata/identity.json`. `LoadIdentity(io.Reader)` decodes, defaults and validates it; `Identity.Validate()` validates programmatically constructed values. Version is 1. Required fields are organization_id, issuer, caller/agent audience lists, explicit allowed_algorithms (RS256/ES256), and bindings. Clock skew defaults to 30 seconds on loading; explicit zero is preserved and maximum is 60. Programmatic structs must set their own default.

An optional organization_claim names a token claim that #12 must require to equal organization_id. Otherwise the configured realm defines the organization boundary; this does not require Keycloak Organizations. No body field establishes organization membership.

Each binding has subject, client_id and kind (agent/integration). An integration's optional assertable_agents list contains exact registered agent subjects. Omission means no permission. No wildcards or cross-realm references. Never put tokens, passwords or client secrets in this configuration. Loading is bounded to 1 MiB, rejects unknown/duplicate keys and nulls, and reports safe field paths and reason codes. Invalid configurations must never be used.

## Lifecycle

To add an agent, provision its client, capture its subject/client ID, and add a binding. To authorize an integration assertion, add that subject to the integration's assertable_agents. Removing permission or a binding requires deployment configuration replacement/restart; runtime reload is outside this issue. A valid access token grants no delegation or policy permission by itself.

Credential rotation preserves identity. Deleting/recreating a service account requires explicit rebinding and updating assertion references. Old subjects are never automatically transferred. Keycloak account deletion does not immediately invalidate already issued JWTs; #12 bounds signing-key cache freshness, and token expiry bounds ordinary offline validation. An agent identity represents its registered service, not an individual execution. RFC 8693 actor claims will not automatically establish delegation or assertion permission.
