# Authopus

The core problem in the current AI Agent implementations is the absence of a consistent, interoperable, and developer-friendly security layer that makes AI agent identity, delegated authority, fine-grained authorization, runtime enforcement, and accountability work together.

Authopus aims to address this problem through an open-source platform that provides centralized AI agent identity, authorization, policy enforcement, delegation, and accountability.

Authopus aims to provide production-grade security infrastructure that developers and organizations can adopt, deploy, integrate, and extend.

The fundamental question Authopus seeks to answer is:

> Which agent is attempting what action, on whose behalf, under what authority, is that action permitted and how can the decision and outcome be verified?

## Core architecture

Authopus follow a modular architecture:

```mermaid
flowchart TB
    Agent["AI Agent / MCP Client"]
    Server["MCP Server / Tools"]

    subgraph Authopus["Authopus"]
        direction TB

        subgraph Enforcement["Enforcement Layer"]
            Gate["GATE<br/>MCP Interception & Enforcement<br/>Policy Enforcement Point"]
        end

        subgraph Core["Identity & Authorization"]
            direction LR
            Delegate["DELEGATE<br/>Agent Identity<br/>Delegation & Authority"]
            Guard["GUARD<br/>Native Go Policy Engine<br/>Policy Decision Point"]
        end

        subgraph Management["Policy Management"]
            Codex["CODEX<br/>Policy Definition<br/>Validation & Versioning"]
        end

        subgraph Accountability["Audit & Observability"]
            Watchtower["WATCHTOWER<br/>Decision Audit<br/>Enforcement Events"]
        end

        Gate -->|"Verify identity and authority"| Delegate
        Delegate -->|"Trusted identity context"| Guard
        Gate -->|"Authorization request"| Guard
        Codex -->|"Active validated policies"| Guard
        Guard -->|"Allow / Deny + explanation"| Gate

        Gate -.->|"Enforcement events"| Watchtower
        Guard -.->|"Decision evidence"| Watchtower
        Delegate -.->|"Identity and delegation events"| Watchtower
        Codex -.->|"Policy lifecycle events"| Watchtower
    end

    DB[("PostgreSQL<br/>Persistent Storage")]
    PolicyFiles[("Versioned YAML<br/>Policy Definitions")]
    Operator["Administrator / Developer"]

    Agent -->|"MCP request"| Gate
    Gate -->|"Authorized MCP request"| Server
    Server -->|"Tool response"| Gate
    Gate -->|"Enforced response"| Agent

    Operator -->|"Manage policies"| Codex
    PolicyFiles --> Codex

    Delegate <--> DB
    Codex <--> DB
    Watchtower --> DB

    classDef external fill:#e8eef8,stroke:#6b87b3,color:#182b48
    classDef enforcement fill:#dbeafe,stroke:#2563eb,color:#1e40af
    classDef security fill:#e0f2fe,stroke:#0284c7,color:#075985
    classDef management fill:#f3e8ff,stroke:#9333ea,color:#6b21a8
    classDef audit fill:#dcfce7,stroke:#16a34a,color:#166534
    classDef storage fill:#f1f5f9,stroke:#64748b,color:#334155

    class Agent,Server,Operator external
    class Gate enforcement
    class Delegate,Guard security
    class Codex management
    class Watchtower audit
    class DB,PolicyFiles storage
```

Refer to [architecture.md](./docs/architecture/architecuture.md) for more details

### Modules responsibility

- Gate: Intercepts requests and enforces authorization decisions
- Delegate: Establish trusted agent identity and delegated authority
- Guard: Evaluates authorization using native policy evaluation Engine
- Codex: Validates, versions and activates YAML policies
- Watchtower: Records security decisions, policy changes and enforcement outcomes

### Key architectural properties

- Single deployment: All modules within one service.
- MCP-first: Provide enforcement for MCP tools and resource interactions
- Fail-closed authorization: Unresolved, invalid or failed authorization never results in implicit permission
- Gate is the policy enforcement point while Guard is the policy decision point


