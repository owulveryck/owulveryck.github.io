---
title: "Who Does What? Team Topologies for the Agentic Platform"
date: 2026-06-22T10:00:00+02:00
lastmod: 2026-06-22T10:00:00+02:00
images: [/assets/team-topologies-agentique/octo_team_topologies_agentique.en.svg]
draft: false
keywords: ["team topologies", "agentic platform", "cognitive load", "organization", "agentic engineering"]
summary: >
  The agentic platform defines what needs to be provided. Team Topologies defines who provides it, and how teams interact to make it happen. This second article applies Skelton and Pais's model to a world where agents produce and humans orchestrate.
tags: ["AI", "agents", "architecture", "platform", "team-topologies", "agentic-engineering"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

> *The agentic platform defines what needs to be provided. Team Topologies defines who provides it, and how teams interact to make it happen.*

---

In [the first article of this series](/2026/06/19/vibe-coding-at-scale-engineering-strikes-back.html), we asked the **what**: which systemic capabilities (context, guardrails, tooling) are needed to produce reliable applications at scale. The answer was the agentic platform, and at its core, the *agentic factory*: the mechanism where agents plan, code, test, and ship.

But a platform does not build itself, and more importantly, it is not consumed the same way it is built. A fundamental question remains: **who does what?**

<!--more-->

## The real problem: the cognitive load of agentic production

*Before asking who does what, we need to understand why the question is different this time.*

Building an application used to mean orchestrating roles over time: one person designed, another challenged the architecture, a third tested, a fourth deployed. The complexity was real, but **distributed** — across several people, and spread over time. Each role asked its questions in turn.

Agents change the equation. They do not ask questions: they produce answers, immediately. They never tire, never rest, never wait. Their speed is their strength, and their trap. All the questions that roles used to raise *sequentially*, the human steering the agent must now anticipate *upfront, in parallel, in the short window of a prompt*. Poorly framed, an agent does not slow down: it produces fast, and off target.

Cognitive load does not disappear with AI — it transforms. It first becomes an **anticipation burden**: everything the human must foresee before launching the agent, or the output will fall short. And because the agent produces continuously, without human rhythm, it also becomes a **cognitive throughput** problem: a sustained flow of decisions over time. The real challenge of agentic production at scale is not that complexity grows — it is that complexity **compresses**, onto a single person and into a timeframe they cannot absorb alone.

This is exactly what the platform addresses. It absorbs part of the anticipation burden by making itself **queryable by the agent**: "don't worry about security" means the agent will ask the platform how to proceed, and deterministic controls will enforce the outcome downstream. The platform does not eliminate thinking — it **narrows the set of questions the human must carry**, letting them focus on what truly matters: the contested, structural decisions where human judgment remains irreplaceable.

Cognitive load is therefore no longer just, as Skelton and Pais describe, a *quantity to distribute across teams*. In the agentic world, it is also a *throughput to regulate over time*. Team Topologies tells us how to distribute; we still need to say how to absorb. That is the subject of this article.

## Team Topologies, an answer to the load

To tackle this load problem, we draw on *Team Topologies*[^skelton2019], an organizational model that defines four team types and three interaction modes. This is no coincidence: **its founding argument is precisely cognitive load** — a team can only be effective if it carries no more complexity than it can absorb. Skelton and Pais reason about *structural* load, to be distributed across teams; we extend it to the *dynamic* load of the agentic production act, described above. The model gives us the grid for distributing what can be distributed, and identifying what the platform must absorb.

Let us state our position upfront: what follows is a **forward-looking conviction**, the organizational target toward which, in our view, application production by agents is heading. The question is no longer *what the factory needs*, but *who operates it*.

This is exactly what the agentic platform does: it **absorbs technical complexity** so that business teams carry only the cognitive load of their domain. The developer shifts: they build the platform that enables others to produce. Symmetrically, application production opens up to business teams through agents.

### What we keep and what we adapt

Applying Team Topologies to the agentic context requires clarity about our departures. To avoid hollowing out the model while claiming its name, here is the line.

**What we keep intact:** the guiding principle of **cognitive load**; the **four team types**; the **three interaction modes**; the **shift from collaboration to *X-as-a-Service*** at maturity.

**What we adapt, and why:**
- *stream-aligned* teams can be **non-technical** (business), because the platform absorbs the technical load;
- they do not carry **end-to-end operational responsibility** (run, incidents), which the platform absorbs;
- *enabling* is not merely transient: it is **structurally compensated** by the platform, since the producer is no longer a developer;
- **cognitive load** is no longer just a quantity to distribute across teams, but a **throughput to regulate over time**: the platform absorbs the anticipation burden upstream of agentic production.

These adaptations do not betray Skelton and Pais's intent: they apply it to a context they had not anticipated — one where agents produce, and humans orchestrate.

## Four team types, one objective

The objective is shared: produce reliable applications, consistent with organizational standards, at scale. But the roles are distinct.

![Team Topologies of the Agentic Platform](/assets/team-topologies-agentique/octo_team_topologies_agentique.en.svg)
*The four Team Topologies team types applied to the agentic platform. Each team plays a specific role in the production chain.*

### Stream-aligned teams: producing applications

Stream-aligned teams are product teams. They drive the AI orchestrator (the engine of the agentic factory described in the first article), define business intent, and supply the **dynamic context**: specifications, product-specific guardrails, domain knowledge.

The transformation runs deep: these teams no longer need to be composed of developers. Increasingly, they are **business teams** (domain experts, product managers, analysts) who directly drive production through agents. This shift brings production closer to the need, but it carries a risk: these teams may not always appreciate the implications of putting an application into production. This is precisely why the other team types exist.

A point of rigor is in order. In classic Team Topologies, a stream-aligned team is responsible **end to end** for a value stream, including operations and incidents. Here, the **platform absorbs operational responsibility** (deployment, monitoring, rollback). The stream-aligned team remains responsible for the *what* (intent, business quality); the platform guarantees the *how* (reliable production deployment). This split requires a [sufficiently mature platform](#a-platform-is-mature-when). **On-call duty** is distributed accordingly: the platform team handles **systemic incidents** (infrastructure, guardrails, pipelines); **business decisions** (content removal, product rollback) stay with the stream-aligned team.

This what/how boundary is **more porous than it appears**. Brand consistency, for example, is a business concern (the *what*), but its verification is automated by the platform (the *how*). The platform guarantees *minimums*; it does not guarantee business excellence.

### The platform team: industrializing capabilities

The platform team provides the three systemic pillars as X-as-a-Service:

- **Systemic context**: instructions, roles, shared business knowledge, memory, examples and patterns
- **Systemic guardrails**: security, reliability, brand consistency, conventions
- **Tooling and skills**: MCP servers, CI/CD pipelines, evaluations, shared skills

The model is self-service: documented, versioned, consumable without friction. Design effort is invested once, then applied to every project.

#### A platform is "mature" when…  

To keep the word from remaining an empty promise, here are the observable criteria:

- **Guardrail coverage**: critical dimensions (security, reliability, brand consistency) are covered automatically, not by verbal agreement
- **Pipeline reliability**: deployment success rate is measurable and tracked (internal SLA)
- **Self-service share**: the majority of deployments happen without platform team intervention
- **Documentation completeness**: every exposed capability is documented and accompanied by examples
- **Decision traceability**: guardrails produce an audit trail (why a deployment was blocked, which rule was applied, what threshold was breached)

Until these criteria are met, the platform cannot absorb the operational responsibility of stream-aligned teams, and the model relies on enabling to compensate.

### Enabling teams: bridging the gap

The enabling team is **temporary by nature**. Its goal is not to become indispensable, but to make product teams autonomous. In practice:
- **Environment provisioning**: tools, access, configuration
- **Training**: on context packaging, guardrails, orchestrator operation
- **Shift-left** on practices (security, testing, quality) until they are encapsulated by the platform

It bridges the gap between business intent and quality requirements. Its role diminishes as product teams gain proficiency and the platform matures.

An important nuance: in the agentic world, stream-aligned teams often remain **predominantly non-technical**. The gap therefore never closes entirely through upskilling alone — it is **structurally compensated by the platform**. This is not a failure of the model; it is an adaptation to a context where the application producer is no longer a developer.

### Complicated subsystem teams: mastering technical complexity

*Complicated subsystem* teams work on the most technically demanding aspects of the AI infrastructure. Their expertise is deep and specialized: it must not be diluted across product teams.

**This team type is not universal.** An organization that exclusively consumes model APIs may not need one. However, as soon as it manages its own models, optimizes inference costs, or faces sovereignty constraints, this team type becomes essential — as it does for evaluation, red-teaming, advanced RAG engineering, or fine-tuning.

They collaborate with the platform team on model efficiency, KV cache, sovereign inference, cost optimization, and evaluation. Their work never reaches product teams directly: it flows through the platform.

## The three interaction modes

Team Topologies defines three interaction modes between teams:

**Facilitating**: the enabling team *facilitates* stream-aligned teams. A temporary interaction, oriented toward autonomy: the enabler teaches the product team to do it themselves.

**X-as-a-Service**: the platform delivers its capabilities as self-service. This is the target interaction — the one that enables scaling.

**Collaboration**: the complicated subsystem team *collaborates* with the platform team. A deep interaction, justified during the construction phase. At maturity, it evolves into X-as-a-Service: AI efficiency capabilities become consumable services, not permanent construction sites.

---

## Making the model last

Defining teams and their interactions is not enough. An organizational model has value only if it outlives its launch.

### The journey toward autonomy

The enabling team's role is designed to shrink: that is the sign the system is working.

![Evolution of Interaction Modes](/assets/team-topologies-agentique/octo_team_topologies_evolution.en.svg)
*The three interaction modes evolve together: facilitating fades, collaboration gives way to X-as-a-Service, which becomes the dominant mode.*

This evolution follows two parallel axes that must be distinguished: **team maturity** and **platform maturity**. They reinforce each other, but do not progress at the same pace.

**At startup:**
- *Team side*: product teams are discovering context packaging and orchestrator operation. Enabling is omnipresent.
- *Platform side*: capabilities are basic (a few guardrails, partial documentation, a still-fragile pipeline). The platform is not yet mature by the criteria defined above.

**During maturation:**
- *Team side*: product teams have mastered the fundamentals. Enabling becomes targeted (advice on a specific guardrail, feedback on a complex prompt).
- *Platform side*: guardrails expand, self-service progresses, documentation takes shape. Collaboration with complicated subsystems starts transforming into integrated services.

**At autonomy:**
- *Team side*: product teams are autonomous. Enabling is optional, limited to occasional expertise.
- *Platform side*: self-service is complete, guardrails cover critical dimensions, traceability is in place. The dominant interaction is X-as-a-Service.

**Enabling disappears because it succeeds, not because it fails.**

**A concrete example.** The marketing team wants to produce a landing page. They supply the dynamic context (campaign intent, key messages). The platform injects the systemic context (brand guidelines, UI components, accessibility rules) and guardrails verify brand consistency and security. Enabling had trained marketing three months earlier — today, just an occasional check-in. The complicated subsystem team has optimized the KV cache: the brand context, already tokenized, is served from cache instead of being recomputed on each iteration, reducing cost per generation. Result: a compliant, secure, deployed landing page — without marketing needing to know what a CI/CD pipeline is.

The flip side: a team that never "sees" the protection mechanisms loses the ability to judge their relevance. Guardrails must be **transparent in their decisions**, even if they are opaque in their implementation.

### Application governance: preventing shadow IT at scale

If non-technical business teams can produce applications in production, **who decides that an application has the right to exist?** And who manages its lifecycle (technical debt, deprecation, cumulative cost)?

Without governance, you end up with industrialized shadow IT. The platform is the lever for this governance: by centralizing deployment, monitoring, and usage metrics, it provides **systemic visibility** over the application portfolio. The platform's product owner can track active applications, identify those no longer maintained, and trigger their deprecation. An application that no longer passes security checks is flagged automatically, not silently forgotten.

The principle is simple: **ease of production must be matched by ease of oversight**. If the platform makes it trivial to create an application, it must make it equally trivial to know how many exist, who uses them, and what they cost.

### The graduation path: from specific to systemic

A mechanism often missing from platform discussions: **when does a "product" guardrail become a "platform" guardrail?**

Example: the marketing team implements a WCAG contrast guardrail. The HR team encounters the same need three months later. Then the e-commerce team. This is the graduation signal: a guardrail repeated by at least three teams becomes a candidate for systemization — the *rule of three*[^fowler] applied to guardrails. The platform's product owner generalizes the candidate, makes it configurable, documents it: all teams benefit.

Stream-aligned teams surface recurring needs. The platform team evaluates and integrates. Enablers spot cross-cutting needs during their engagements.

Autonomy is not silence: even at maturity, product teams remain the sensors that feed graduation. Self-service does not eliminate the feedback loop — it makes it cheaper. The platform is a *living product*[^skelton2019], with its own backlog, product owner, and iteration cycles.

**A risk must be named.** This product owner potentially carries three burdens: the platform's technical backlog, guardrail graduation, and application portfolio governance. Paradoxically, this is the very bottleneck the model claims to prevent.

**Why this risk emerges.** The more the platform succeeds, the more decisions it concentrates, and a single role cannot manually arbitrate the flow of needs from dozens of product teams.

**The answer lies in tooling.** The platform must automate graduation candidate detection, usage-based prioritization, and portfolio tracking. The PO arbitrates; they do not carry everything.

Without this graduation mechanism, the platform stagnates. Product teams reinvent the same solutions. Guardrails remain local. Systemic reliability erodes.

---

## Operational synthesis

With the principles in place, it remains to crystallize them into actionable rules.

### Who owns what?

Each capability has **one owner** (accountable for the outcome) and often **contributors** (who co-build without being accountable). Without this clarity, capabilities fall through the cracks.

![Responsibilities by team](/assets/team-topologies-agentique/octo_team_topologies_responsabilites.en.svg)
*Each capability has an owner and contributors. Empty cells mean "not involved." Specifics stay with product teams, systemics with the platform.*

The guiding principle is simple: the specific (dynamic context, product guardrails, business knowledge) stays with stream-aligned teams, because domain knowledge *belongs* to the business. The systemic (systemic context, organizational guardrails, tooling, orchestrator) is industrialized by the platform. The AI orchestrator is the emblematic case: the tool is provided by the platform; the configuration stays with the product team.

### A conviction

The agentic platform answers the *what*. Team Topologies answers the *who* and the *how* — not as a rigid partition, but as a **distribution of accountability**: business intent and domain context to the product team, production reliability to the platform.

Without these boundaries, the platform becomes either a bottleneck or a political battleground.

This model is not universal: it addresses organizations that produce **multiple agentic applications in parallel**. In practice (an empirical heuristic, not an absolute threshold), starting from three to five product teams, the cumulative cost of reinvention (context recreated, guardrails reimplemented, inconsistencies to fix) exceeds the investment in a shared platform.

Applied to agentic engineering, Team Topologies provides a structure that:

- Enables **business teams to become stream-aligned**, without requiring them to become developers
- Ensures that **systemic capabilities are industrialized**, not reinvented project by project
- Acknowledges that the **gap between business and production is real**, and bridges it with enablers, not with hope
- Isolates **technical complexity** in specialized teams that feed the platform

What changes is that producing an application **no longer requires coding**, but demands other skills: articulating intent, structuring context, operating an orchestrator. What does not change is that you need **guardrails, consistency, and reliability** to put it into production.

This shift addresses the problem raised at the outset. Today, it is the developer who carries the anticipation burden: they know which questions the agent will not ask. Tomorrow, as the platform absorbs an increasing share, the skill threshold required to operate lowers: a business team can produce without bearing alone what the platform anticipates on their behalf. **The producer's trajectory, from developer to business, is not an assumption — it is the measurable consequence of an anticipation burden progressively absorbed.**

The developer does not disappear: they shift. They no longer produce the application — they build the platform that enables others to produce it.

---

## Sources

[^skelton2019]: Matthew Skelton, Manuel Pais — *"Team Topologies: Organizing Business and Technology Teams for Fast Flow"*, IT Revolution, 2019.

[^osmani2026]: Addy Osmani, Shubham Saboo, Sokratis Kartakis — *"The New SDLC With Vibe Coding: From ad-hoc prompting to Agentic Engineering"*, Google, May 2026.

[^fowler]: Don Roberts, cited by Martin Fowler — *"Refactoring: Improving the Design of Existing Code"*, Addison-Wesley, 1999. "Three strikes and you refactor." The principle is transposed here from code refactoring to guardrail graduation.
