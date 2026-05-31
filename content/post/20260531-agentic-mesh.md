---
title: "The Agentic Mesh: Cognitive Automation at Scale"
date: 2026-05-31T10:00:00+02:00
lastmod: 2026-05-31T10:00:00+02:00
images: [/assets/agenticmesh/poster-agent-mesh.en.svg]
draft: false
keywords: []
summary: "Convictions for designing tomorrow's agentic systems — an agent mesh system capable of delivering value at scale, borrowing principles from the data mesh."
tags: ["AI", "agents", "architecture", "agentic-mesh"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

Today, we see many initiatives around the agentic paradigm. Most revolve around systems built by AI giants (Anthropic, Google, OpenAI) and often boil down to pushing natural language directives to an off-the-shelf orchestrator. One specializes an integrated system like Claude Code through markdown files, skills, and tools. Frameworks like BMAD illustrate this approach well: they transform such a system into a mixed AI/human development team. It's useful, it's fast, it's a good starting point.

But to truly unlock the full value of AI, agents will need to embed themselves in **business processes**. Their contribution goes beyond cognitive automation: it lies in the ability to **make certain decisions** that allow processes to gain in effectiveness, not just efficiency. And this can only be fully achieved when we have an **ecosystem of agents** that communicate with each other and collaborate with humans.

Not when agents do our work, but when we have **properly equipped the information system** so that agents do *their* work, and we, humans, do ours.

This article proposes an **agent mesh** system (the agentic mesh) capable of delivering value at scale. Many of its principles, if not all, are borrowed from the **data mesh**, a paradigm whose ideas were remarkable but which unfortunately did not achieve the success it deserved. In my view, this failure does not stem from the quality of its ideas or its design. Consequently, many can be recycled here (and that is what we will do).

> **Note.** This article was co-written with an AI. I am the pilot: I declare the intentions, the ideas, and I review the entire document. The actual writing (the pen) was done by a robot. My goal is to share these ideas to open a discussion, not to write a technical masterpiece that becomes a stylistic reference. This version is designed for humans; if you prefer a version suited for AI consumption, the [markdown source](https://raw.githubusercontent.com/owulveryck/owulveryck.github.io/refs/heads/master/content/post/20260531-agentic-mesh.md) is available.

---

## The agent, that fundamental building block of future information systems

Let us start by establishing some vocabulary, as the confusion between two closely related concepts is at the root of many architectural mistakes.

We define:

- a **tool** as a deterministic function: same input, same output. It programmatically executes an operation via an API or a protocol like MCP. It is idempotent, predictable, testable.
- an **agent** as a non-deterministic entity: it works and decides autonomously, based on an **intention**, depending on an **execution context**. Its result is statistical (it depends on the context at the time of the call).

The consequences of these definitions are significant:

- the tool is merely an action support (it executes, it does not interpret);
- the decision an agent makes from the same tools can vary depending on the context, the state of the world, and the intention entrusted to it.

For example: calling a weather API always returns the same data for an identical request — that's a tool. Deciding whether to postpone a delivery based on the weather, the contractual deadline, the client priority, and available resources — that's an agent decision.

> **The tool is deterministic. The agent is not: its result is statistical and depends on the execution context.**

![Agent vs Tool - poster](/assets/agenticmesh/poster-agent-vs-tool.en.svg)

This distinction is not merely academic. It has direct consequences on governance, on SLOs, on ownership, and on how to test and trust a system. You don't govern an agent the way you govern an API. The contract is not an input/output schema (it is an **intention and a decision scope**).

---

## The dominant paradigm: the agent as an off-the-shelf directive

Today, when a company integrates AI into its IS, the dominant reflex is not so much centralization as **the use of off-the-shelf agentic systems** (Claude Code equipped with frameworks like BMAD, Cursor with its rules, or platforms like n8n with AI nodes). In this paradigm, the agent execution layer is not separated: agents are seen as **simple natural language directives** applied to generic agents to specialize them.

This approach is appealing. It promises rapid implementation, without significant engineering investment. It is sufficient for prototypes, for one-off automations, for validating that a use case can be handled by AI.

However, it quickly reaches its limits when trying to move to production on complex use cases.

### What the apparent simplicity conceals

An agent is not just a prompt. An agent that works in production is a software system that must handle:

- **controlled parallelism** between sub-agents working simultaneously;
- **structured feedback loops** between a reviewer and producers, with retry counters and issue filtering;
- **programmatically validated structured outputs** between each step;
- **prompt caching** shared between invocations to control costs;
- **typed shared mutable state** protected against concurrent access;
- **fine-grained observability** (tokens consumed per agent, latency, success rate, decision traceability).

To illustrate concretely, here is a comparison from a real multi-agent orchestrator project (a presentation generation system from a natural language request):

| Criterion | Native code implementation | Off-the-shelf system + directives |
|---|---|---|
| Fine-grained parallelism | Goroutines + semaphore, concurrency control | Limited, sequential or simple parallel |
| Feedback loops | Typed issues, targeted retry on subset | Conversational, fragile over time |
| Structured outputs | Strict JSON schema, programmatic validation | Implicit, model-dependent |
| Prompt caching | Shared between sub-agents, controlled costs | Not available or not shared |
| Inter-step state | Typed, mutex, testable | Lives in conversational context |
| Observability | Per-agent metrics, complete issue log | Limited, aggregated |
| Maintainability | Typed, testable, evolvable | Fragile to model changes |

The off-the-shelf system remains relevant for a **rapid prototype** or a simple linear pipeline. It does not hold up under the load of a complex production system (not because it is bad, but because it pushes complexity into the prompt rather than into code, and the prompt is a fragile place to manage that complexity).

### Software engineering is not only in the tools

This is the founding conviction that distinguishes the agentic mesh from common approaches: **a production agent is a software engineering product, not an assembly of directives**. Engineering is not only found in the tools the agent uses (it is also found in the way the agent itself is built, tested, deployed, observed, and governed).

This conviction has a practical consequence: to build an agentic mesh, agents must be treated the same way as other critical IS components. With a dedicated platform, explicit contracts, clear ownership, and engineering rigor comparable to that applied to microservices or data-products.

---

## The four pillars of the agentic mesh

For each agent to offer adapted and validated solutions in a specific business domain, the agentic mesh rests on four pillars:

1. **The digital platform** gives agents access to action capabilities and IS information, respecting governance rules.
2. **Domains** anchor each agent in its business: it is the domain that defines the agent's intention, assumes responsibility for it, and validates the quality of its decisions.
3. **The agent as a product** (AI-as-a-product) requires treating the agent with the same rigor as a software product: contract, lifecycle, traceability, interoperability.
4. **Automated governance** acts as an enabler: it enables the construction and production deployment of agents under the security and reliability conditions required by the organization, and facilitates exchanges between agents and between agents and humans.

### Pillar 1: The digital platform: the real MVP

The first instinct of a client who wants to "add AI" to their IS is to start with the agent. This is a sequencing error.

**The real MVP is the platform.** Before deploying sophisticated agents, the IS must be able to expose cleanly, securely, and in an automatable way two types of fundamental capabilities:

- **information access capabilities**: understanding a context, querying business data, reading the system state;
- **action capabilities**: modifying a status, triggering a process, writing to a system.

These capabilities are exposed in compliance with **governance rules**: authorizations, traceability, compliance. This is where platform and automated governance converge: the platform does not expose raw access, it programmatically applies the rules that governance has defined.

This **read / write** decoupling is not new (it is a proven pattern — CQRS, event sourcing — reinterpreted for the agentic era). It has a crucial property: **it delivers value independently of agents**. A platform that properly exposes its capabilities is useful for humans, for integrations, for classic automations, and only then for agents.

This is what makes the MVP sellable: it is not a bet on AI. It is a sound architectural investment, of which agents will be the first beneficiaries.

![Platform foundations - MVP platform + POC agent](/assets/agenticmesh/poster-fondations-socle.en.svg)

**What the platform MVP validates:**
- Can we expose the right information for an agent to understand a context?
- Can we expose the right actions for an agent to act securely?
- What API granularity? What protocol (REST, MCP, events)?
- How to manage authorizations when it is an agent calling, not a human?

### Pillar 2: Domains: the agent belongs to its business

The agentic mesh borrows from Eric Evans' Domain-Driven Design the fundamental notion of **domain**. A domain is a logical division of the organization that brings together a coherent business activity, its own vocabulary (*ubiquitous language*), and the boundaries within which its models have a precise meaning (*bounded context*).

In the agentic mesh paradigm, **an agent belongs to a domain** (just as a microservice or a data-product belongs to its domain). There is no separate "agentic domain" from the business: there are business domains that own and operate their agents, just as they own their data and their services.

This belonging is not merely organizational. It structures:

- the **intention** of the agent (expressed in the domain's vocabulary);
- the **decision scope** (delimited by the domain's boundaries);
- the **ownership** (an agent product owner, member of the domain, assumes responsibility);
- the **local governance rules** (who can invoke the agent, under what conditions, with what level of autonomy).

Within a domain, two categories of participants around an agent are distinguished:

- **context producers** (source systems, sensors, APIs, other agents) that feed the agent with reliable information;
- **decision consumers** (downstream processes, humans, agents from other domains) that leverage the decisions and actions produced by the agent.

A participant can simultaneously be a context producer and a decision consumer.

What fundamentally changes with this model is that **the responsibility for the agent (its decision quality, its scope, its evolution) falls to the business domain that leverages its value**, not to a central AI team detached from the business.

![Decoupling and business domains - from monolith to agent-product](/assets/agenticmesh/poster-decouplage-domaines.en.svg)

### Pillar 3: The agent is an AI-as-a-product

To understand the agent's place in the mesh, we must first clarify two radically different ways of integrating AI into an IS.

**An AI-product** is a digital product that solves a user problem by using AI internally. The user consumes the product, not the AI. The AI is a means, hidden in the mechanics. A recommendation system for e-commerce, a customer support assistant, a content generation tool — these are AI-products. Their value is measured by the user experience of the product, not by the raw quality of the AI powering it.

**An AI-as-a-product** is the opposite: the AI *is* the product. What is exposed, contractualized, and consumed is the AI capability itself. The consumers are not end users — they are other systems (humans or agents) that invoke this capability in their own processes.

In an AI-product, the AI is an **internal means**. In an AI-as-a-product, the AI is **the contractual interface**.

![AI-product vs AI-as-a-product](/assets/agenticmesh/ia-product-vs-as-a-product.en.svg)

**An agent is an AI-as-a-product.** It is even the most accomplished form: an agent is an autonomous AI capability that decides and acts based on an intention, exposed as a consumable product. Considering the agent as a product is not a nuance — it is the condition for it to be able to join a mesh.

> *"We deployed an AI assistant for our customer support, we're ready for the agentic mesh."*
> 
> No — you have an AI-product. To join the mesh, the decision-making capability of that assistant would need to be exposed, contractualized, and invocable by other domains. This is a change in nature, not an extension.

The mesh is built with agents (therefore with AI-as-a-products). Without an exposed and contractualized agentic capability, there is nothing to mesh.

#### The 7 affordances of an agent

For an agent to truly be a product, it must offer a set of capabilities that go well beyond "the agent that answers questions." These capabilities (which we call **affordances**, in reference to Zhamak Dehghani's work on data-as-a-product) define what an agent **does**, rather than what it **is**.

**Affordance 1: Expose decisions and actions**
The agent exposes its capabilities via clearly defined interfaces. What it can decide, what it can do, under what conditions — all of this is explicit and invocable in a structured way, by humans as well as by other agents (A2A, MCP).

**Affordance 2: Consume context**
To decide, the agent consumes context from various sources: platform APIs, other agents, business data, system state. These sources are documented and are part of the product's contract.

**Affordance 3: Reason and decide**
This is the core of the product — the internal orchestration logic: the model used (or models, as a complex agent often combines several models adapted to each sub-task), the versioned prompts, the available tools, the feedback loops, the state management. It is precisely this logic that deserves to be treated as code, not as a directive.

**Affordance 4: Be discoverable and understandable**
The agent is referenced in an agent registry. Its intention, its decision scope, its invocation conditions, and its known limitations are documented and accessible — to humans as well as to automatic orchestrators.

**Affordance 5: Manage its lifecycle**
The agent must be able to evolve without interrupting its consumers. This implies prompt and model version management, deprecation strategies, and update mechanisms that do not impact the exposed interfaces.

**Affordance 6: Trace decisions**
Every decision made by the agent is traceable: input context, reasoning followed, action taken, observed result, tokens consumed, model used. This traceability serves both as an auditability tool (compliance, regulation) and as a continuous improvement tool (identifying cases where the agent makes mistakes and calibrating its scope).

**Affordance 7: Be governable**
The agent offers cross-cutting steering capabilities: access policy (who can invoke it?), autonomy level (when does it escalate to a human?), kill switch (can it be disabled in an emergency?), and management of personal data it processes.

![The 7 affordances of an agent](/assets/agenticmesh/affordances.en.svg)

#### The agent's contract: a condition for meshing

If the agent is a product, it needs a **contract**. Without an explicit contract, there is no reliable mesh — we fall back into implicit coupling, manual discovery, and interpersonal trust that does not scale.

This contract must answer specific questions: who is this agent? What can it do? How to invoke it? What are its security conditions? What formats does it accept? This is not an internal design document — it is a **public interface**, readable by humans and parseable by systems.

The ecosystem did not wait for the agentic mesh to formalize this need. The **A2A** protocol (*Agent-to-Agent*), led by Google and supported by more than fifty industry players, proposes a concrete artifact: the **Agent Card**.

#### The Agent Card of the A2A protocol

The Agent Card is a JSON document that functions as the identity card of an agent in an A2A network. It is published at a standardized URL (`/.well-known/agent-card.json`, following the RFC 8615 convention), which allows any client or agent to discover it automatically through a simple `GET`.

What makes the Agent Card interesting as a contract model is that it covers all the dimensions necessary for meshing:

**Identity and intention** (The Agent Card exposes a `name`, a `description`, and a service `url`). These are the minimum required fields. The `provider` (organization, URL) allows the agent to be attached to its ownership domain. The agent `version` and the `protocolVersion` ensure compatibility between consumer and provider.

**Technical capabilities** (The `capabilities` block declares what the agent supports beyond the simple synchronous call: `streaming` (streaming responses via SSE), `pushNotifications` (asynchronous callbacks), `stateTransitionHistory` (state transition traceability of a task)). This is an explicit capability negotiation, not discovery by trial and error.

**Skills** (The `skills` array is the functional core of the card. Each skill has an `id`, a `name`, a `description`, categorization `tags`, and optionally `examples`) — example prompts or inputs that the skill can process. Each skill can declare its own `inputModes` and `outputModes` (MIME types: `text/plain`, `application/json`, `image/png`...), or inherit the agent's default modes.

**Security** (`securitySchemes` follow the OpenAPI convention: bearer tokens, API keys, OAuth 2.0, OpenID Connect. The `security` field declares which schemes are required to invoke the agent. No surprise at execution time) — access conditions are in the contract.

**Extended discoverability** (The `supportsAuthenticatedExtendedCard` flag indicates that a more detailed version of the card is available behind authentication, for sensitive information or private capabilities).

The whole thing is constrained to **10 KB maximum**, which forces conciseness and ensures rapid discovery.

#### What the Agent Card does not cover (yet)

The Agent Card is an excellent foundation for an interoperability contract, but it does not cover everything an agent-as-a-product needs in the context of an agentic mesh:

- **The trust contract** (correct decision rate, escalation conditions to a human, decisional SLOs). The Agent Card says *what* the agent can do, not *how much* it can be trusted.
- **Governance** (data retention policy, decision traceability, regulatory compliance). These are organizational concerns that go beyond the scope of an interoperability protocol.
- **Ownership** (which business domain the agent is attached to, who assumes responsibility for its decisions).

These complementary dimensions will need to be carried by contract extensions or governance artifacts specific to each organization. The Agent Card lays the technical foundation; the complete contract of an agent-as-a-product adds the product and organizational dimensions.

### Pillar 4: Automated governance: the mesh enabler

The governance of the agentic mesh is not a framework of externally imposed constraints. It is an **enabler** that makes the mesh possible by acting on two axes.

**Technical enabler**: automated governance provides the conditions for agents to be built and deployed to production with the level of security and reliability required by the organization. This includes traceability and explainability standards, compliance obligations (AI Act, regulation on automated decision systems), interoperability protocols (A2A, MCP), and decision quality metrics. These rules are not documentary: they are **applied programmatically** by the platform. This is what justifies the term *automated*: governance is in the code, not in a wiki.

**Exchange enabler**: governance facilitates interactions between agents and between agents and humans. It defines who can invoke which agent, under what conditions, with what level of autonomy. It manages escalation mechanisms to a human, feedback loops, and data retention policies. It makes the mesh fluid by removing integration friction.

This governance applies at two levels. At the central level, it establishes the constitutional framework: the rules that apply to all agents in the mesh. At the domain level, each team applies this framework within its scope by adapting it to its operational constraints.

A concrete example: central governance will require that any agent making decisions with a financial impact above a certain threshold has an escalation mechanism to a human. The platform will apply this rule automatically. The logistics domain will then define its own policy (specific threshold, escalation delay, competent person) based on its constraints.

---

## From prototype agent to mesh: a three-phase trajectory

Setting up an agentic mesh cannot be decreed. It is built through successive iterations.

![Three-phase trajectory](/assets/agenticmesh/trajectory.en.svg)

### Phase 1: Building the foundation: MVP platform + POC agent

**The objective of this phase is not to build a good agent. It is to validate the fundamental hypotheses.**

A real business use case is selected, the minimum necessary capabilities are exposed on the platform (read + write), and an agent is built with rough tools — no need for it to be state-of-the-art. At this stage, an off-the-shelf system can absolutely suffice: the goal is to validate, not to produce.

What the **platform MVP** validates:
- Can we expose the information necessary for the decision?
- Can we expose the actions necessary for execution?
- What are the data quality, API granularity, and access governance issues?

What the **agent POC** validates:
- If we expose the right services, can we envision cognitive automation?
- Does the agent deliver measurable value on this use case?
- Under what conditions does the agent make mistakes?
- Do we need to structure information upstream?

At this stage, the agent is experimental. It is coupled to the orchestrator, its scope is unclear, its ownership is informal. This is normal.

### Phase 2: Engineering: from prototype to product

The POC validated the concept. We know that cognitive automation is possible on this scope. It is time to evolve the architecture and **move to software engineering**.

This phase consists of transforming the prototype agent into a **structured agentic architecture**:

![Structured agentic architecture](/assets/agenticmesh/architecture-agentique.en.svg)

Each sub-agent is specialized, uses the model adapted to its task (Haiku for the simple, Sonnet for the complex, for example), has its own tools, and its own feedback loops. The orchestrator coordinates without doing everything itself, and it is implemented in **code**, not as a directive.

It is in this phase that the shift from *prompt* to *product* takes place. What was an assembly of directives becomes a testable, observable, deployable software system. The questions asked are no longer "how to formulate the prompt" but "how to version prompts, how to test regressions, how to instrument decisions, how to handle retries."

It is also at this stage that certain sub-agents prove **useful beyond their initial use case**. A sub-agent that summarizes customer conversations for this process could serve other processes. This is the signal that it must change nature: be extracted from the application, given its own intent manifest, and managed as an autonomous product in the domain that leverages its value.

### Phase 3: Building the mesh: from agent-product to ecosystem

Agents are now managed as products in their respective domains. They expose their capabilities via standardized manifests. They are discoverable, invocable, governed.

The mesh can then be built: end-to-end business processes that orchestrate or choreograph sequences of agents from different domains, alternating between human decisions and agentic decisions as needed.

![Business process - cross-domain Human-Agent sequences](/assets/agenticmesh/processus-metier-ha.en.svg)

**H-A** (Human-Agent) sequences are not exceptions or safeguards — they are a **feature** of the mesh. The agentic mesh does not aim to replace humans. It aims to enable them to focus on decisions that truly require their judgment, by delegating the rest.

---

## Organization around the mesh

### Domain ownership

One of the most important organizational changes induced by the agentic mesh is the **shift in ownership**.

In a centralized model, agents belong to a cross-functional AI team. In the agentic mesh, **each agent belongs to the business domain that leverages its value**. It is this domain that:
- defines the agent's intention;
- assumes responsibility for the quality of its decisions;
- manages its lifecycle;
- is accountable for the incidents it causes.

This change in ownership is often difficult to accept, because it implies that business teams acquire a new competency — not in AI, but in **agent product management**. This is a real organizational investment, requiring time and support.

### Enablement as a cross-functional function

Giving autonomy to domains does not mean letting them reinvent the wheel. A cross-functional enablement function is necessary to:

- maintain and evolve the **digital platform** (the technical commodities);
- define and enforce **manifest and interoperability standards**;
- provide **reusable agentic architecture patterns** (orchestration, feedback loops, observability);
- support domains in the transition from **prototype agent to agent-product**.

The success of this function is measured by a simple indicator: **the time between designing an agent and the moment it is invocable by other domains**.

---

## Summary: the virtuous circle of the agentic mesh

The agentic mesh is not just a way to model and manage AI systems. It is above all a way to think about **the agent as a product** — an AI-as-a-product in its own right, designed with software engineering rigor. The product is no longer the application that uses AI internally. **The decision-making capability *is* the product**. Tools and APIs are merely a means to feed it.

Thus, service level indicators will need to be adapted:
- **relevance** and **reliability** indicators for decisions will replace availability indicators as primary metrics;
- decision **traceability** will become a full-fledged SLO, not an option;
- the **escalation rate** to humans will be an indicator of autonomy calibration, not an indicator of failure.

**The agentic mesh is not an improved version of the enterprise chatbot.** The real change it brings is the shift from the **tool-equipped human** (who uses AI as an instrument) to the **accompanied human** (who collaborates with agents within mixed teams, while retaining control of the process. Agents decide within their scope, but it is the human who remains in charge of the overall process: validating, arbitrating, and escalating. It is this collaboration under human control that produces value. This change has a strong impact on product design, on IS architecture, and on team organization).

Existing practices should not be discarded, however. If your IS already relies on principles of clean API exposure, functional read/write decoupling, and organization into coherent business domains — you are already on the right trajectory. The transition to the agentic mesh will be all the more natural as your platform already offers the foundation on which agents can decide.

---

![Agentic Mesh - the 4 pillars of the value-driven mesh](/assets/agenticmesh/poster-agent-mesh.en.svg)


To summarize, the **virtuous circle of the agentic mesh** consists of several steps:

- **build the digital platform** that properly exposes the IS's read and write capabilities;
- **build an agent** on a real use case, with genuine software engineering rigor (not as an assembly of directives, but as a product);
- **attach that agent to its business domain** which assumes ownership and responsibility for its decisions;
- that agent becomes a **node in the mesh** through interoperability standards, automated governance, and shared intent contracts;

…all with the goal of producing **high-value automated business processes**, where humans and agents collaborate naturally, each within their decision scope.

> *Agents are the communication points between domains. Value emerges from the mesh.*

---

## Appendix: Summary poster

The poster below summarizes all the convictions in a single view, organized in 4 zones: [full poster (interactive HTML)](/assets/agenticmesh/poster-architecture-agentique.en.html).

