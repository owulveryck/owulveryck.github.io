---
title: "Who Does What? Team Topologies for the Agentic Platform"
date: 2026-06-24T10:00:00+02:00
lastmod: 2026-06-24T10:00:00+02:00
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
This is a second version more human written. The first version was roughly a translation made by AI. This version is human-edited. Hopefully it will be easier to read by a human.

---

In [the first article of this series](/2026/06/19/vibe-coding-at-scale-engineering-strikes-back.html), we asked the **what**: which systemic capabilities (context, guardrails, tooling) are needed to produce reliable applications at scale. The answer was the agentic platform, and at its core, the *agentic factory*: the mechanism where agents plan, code, test, and ship.

But a platform does not build itself, and more importantly, it is not consumed the same way it is built. A fundamental question remains: **who does what?**


## The real problem: the compression of cognitive load

*Before asking who does what, we need to understand why the question is different this time.*

Building a software application is about orchestrating roles over time: a designer, an architect, a tester, and an operations engineer. The overall complexity is real, but it is **distributed** across several people and spread out over time. Each role asks its questions in turn.

AI agents change the equation. They move beyond answering isolated questions to continuously generating output. Because they execute instantly, the traditional sequential process of problem-solving breaks down. All those questions must now be anticipated by the human in command and fit into the short window of a prompt. If poorly framed, an agent does not slow down: it produces fast, but off-target.

Cognitive load does not disappear with AI: it transforms. It becomes an **anticipation burden** (everything the human must foresee before launching the process) combined with a **cognitive throughput problem** (sustaining a high-velocity flow of decisions). The real challenge of agentic production at scale is not that complexity grows; it is that complexity compresses onto a single person over a timeframe they cannot absorb alone.

This is exactly what the platform addresses. It absorbs part of the anticipation burden by making itself **queryable by the agent**: telling the agent "don't worry about security" means it will ask the platform how to proceed, and deterministic controls will enforce the outcome downstream. The platform does not eliminate thinking; it **narrows the set of questions the human must handle**, letting them focus on what truly matters: the contested, structural decisions where human judgment remains irreplaceable.

Cognitive load is therefore no longer just, as Skelton and Pais describe, a *quantity to distribute across teams*. In the agentic world, it is also a *throughput to regulate over time*. Team Topologies tells us how to distribute; we still need to figure out how to absorb. That is the subject of this article.

## Team Topologies: An Answer to Cognitive Load

To tackle this bottleneck, we can leverage *Team Topologies*[^skelton2019], an organizational model that defines four team types and three interaction modes.

> "All models are wrong, but some are useful." — George Box

This framework has a proven track record in engineering organizations. Its founding argument **is precisely cognitive load**: a team can only be effective if it carries no more complexity than it can absorb.

Skelton and Pais reason about *structural* load, which is distributed across teams. I want to extend this to the *dynamic* load of agentic production. The model gives us a framework for distributing what can be distributed, and identifying what the internal platform (defined in the previous article) must absorb.

To be clear: what follows is a **forward-looking conviction**. It is the organizational target toward which I believe AI-driven application production is heading. The question here is no longer *what the factory needs*, but *who operates it*.

This is exactly the purpose of the agentic platform as I envision it: it **absorbs technical complexity** so that business teams only carry the cognitive load of their specific domain. The role of the developer shifts: they build the platform that enables others to produce. The mirror effect is that application development opens up to business teams, empowered by an AI-based agentic development team.

### What We Keep and What We Adapt


Applying *Team Topologies* to an AI-driven context requires clarity on where we follow the model and where we diverge. The goal is to preserve the core mechanics, not just borrow the buzzwords.

**What stays the same:**
* The core focus on managing **cognitive load**.
* The **four team types** and **three interaction modes**.
* The evolution from close collaboration to an **X-as-a-Service** model at maturity.

**What changes, and why:**
* **Stream-aligned teams can be non-technical (business-led):** The platform handles the technical heavy lifting, allowing domain experts to drive production.
* **No end-to-end operational responsibility:** The platform—not the stream-aligned team—absorbs the "run" and incident management.
* **Enabling is permanent, not transient:** Because the primary producers are business experts rather than software engineers, platform enablement becomes a structural, ongoing requirement.
* **Cognitive load becomes a throughput metric:** It is no longer just a static volume to distribute across teams, but a continuous flow of decisions to regulate over time. The platform absorbs the upfront "anticipation burden" before the agent even starts generating.

These adaptations do not betray Skelton and Pais’s original intent. They simply project it into a reality they hadn't anticipated: a combination of work from human and AI that do not suffer from cognitive load.

## Four team types, one objective

The shared objective across these teams is clear: produce reliable, standardized applications at scale. Their roles, however, remain strictly distinct. Here is a synthesis of how they interact, presented as a "map" (apologies to Simon Wardley, I know this isn't a real map).

![Team Topologies of the Agentic Platform](/assets/team-topologies-agentique/octo_team_topologies_agentique.en.svg)
*The four Team Topologies team types applied to the agentic platform. Each team plays a specific role in the production chain.*

Now let's navigate through all the teams and their roles.

### Stream-Aligned Teams: Driving Production

Stream-aligned teams are the product teams responsible for delivering core business value. In this model, they drive the AI orchestrator (the engine of the factory described in the previous article), define the business intent, and supply the **dynamic context**: specifications, product-specific guardrails, and domain knowledge.

The introduction of AI agents fundamentally transforms these teams: they no longer need to be staffed with software engineers. Instead, they are composed of **business experts** (domain specialists, product managers, analysts) who generate applications directly via agents. The traditional translation layer between business requirements and developer implementation is eliminated; the business *builds* the application. While this brings production infinitely closer to the user need, it introduces a major risk: non-technical teams rarely understand the implications of pushing code to production. This is exactly why the other team types are necessary.

**Deviation from the original model:** In classic *Team Topologies*, a stream-aligned team owns the value stream end-to-end, including operations and incidents ("You build it, you run it"). In this agentic model, the **platform absorbs operational responsibility** (deployment, monitoring, rollbacks). The stream-aligned team owns the *what* (business intent and quality); the platform guarantees the *how* (safe, reliable execution). This split requires a highly mature platform. **On-call duty** shifts accordingly: the platform team handles systemic failures (infrastructure, pipeline crashes, guardrail breaches), while business decisions (feature rollbacks, content issues) remain with the stream-aligned team.

This "what/how" boundary is more permeable than it sounds. Brand consistency, for example, is a business concern (the *what*), but verifying it can be automated by the platform (the *how*). Ultimately, the platform enforces a *baseline of safety and standards*; it does not guarantee business excellence.

### The Platform Team: Building the Engine


The platform team must provide four core pillars as a self-service offering:
* **Global Context:** System prompts, roles, shared business knowledge, memory, examples, and architectural patterns.
* **Deterministic Guardrails:** Security policies, reliability constraints, brand consistency, and coding conventions.
* **Agentic Tooling:** MCP (Model Context Protocol) servers, CI/CD pipelines, evaluation frameworks, and shared agent "skills".
* **Execution Engine:** to actually run the AI-engines (Inference engines)

The engagement model is strictly self-service: everything is documented, versioned, and consumable without friction. The engineering effort is invested once by the platform team and leveraged continuously across all agent-driven projects.

#### When is the platform "mature"?

To ground this in reality, here are the observable criteria a platform must meet to actually empower AI-driven development:
* **Hardcoded Guardrails:** Critical dimensions (security, reliability, compliance) are enforced deterministically via standard code—not just through a stochastic LLM "verbal agreement" or prompt constraint.
* **Measurable Reliability:** Deployment success rates and pipeline health are tracked against internal SLAs.
* **High Self-Service Index:** The vast majority of deployments happen with zero intervention from the platform team.
* **Agent-Readable Documentation:** Every exposed capability is fully documented with examples, formatted not just for humans, but for the agents querying them.
* **Decision Traceability:** Every guardrail intervention leaves a clear audit trail (e.g., exactly *why* a deployment was blocked, the specific rule applied, and the threshold breached).

Until these criteria are met, the platform cannot safely absorb the operational responsibility from the stream-aligned business teams. In the interim, the organization must rely heavily on *enabling teams* to bridge the gap.

### Enabling Teams: Bridging the Gap

The enabling team exists to bridge the gap between business intent and engineering rigor. In classic *Team Topologies*, this team is strictly temporary—its goal is to upskill product teams and then move on. In the agentic transition, they focus on:
* **Environment Provisioning:** Setting up tools, access, and agent workspaces.
* **Agentic Training:** Teaching business users how to properly package context, understand guardrails, and operate the AI orchestrator.
* **Manual Shift-Left:** Enforcing security, testing, and quality practices *manually* until they can be fully hardcoded into the platform.

**A Permanent Gap?** Here the model will probably reach a limit. Traditionally, an enabling team upskills developers until they are fully autonomous. But in the agentic world, stream-aligned teams are increasingly non-technical. There is a hard ceiling to how much you can "upskill" domain experts in software engineering; the gap will never close entirely through training alone.

The platform must provide structural **compensation** for this lack of technical expertise. The enabling team's ultimate job isn't just to train the business; it is to identify exactly what the business *cannot* or *should not* learn, and mandate the platform team to automate it. This is not a flaw in the model; it is a necessary adaptation to a reality where the application producer is no longer a software developer.

### Complicated Subsystem Teams: Encapsulating Deep Tech

Complicated subsystem teams tackle the most technically demanding layers of the AI infrastructure. Their expertise is deep and highly specialized; diluting it across individual product teams would be a massive waste of talent and focus.

**This team type is optional.** If your organization simply wraps external LLM APIs, you probably don't need one. However, the moment you start hosting your own open-weights models, optimizing inference costs, or dealing with strict data sovereignty constraints, this team becomes critical. They handle the heavy lifting of advanced AI engineering: red-teaming, complex RAG architectures, fine-tuning, and custom evaluation frameworks.

They collaborate closely with the platform team on hard engineering problems like KV cache optimization, sovereign inference pipelines, and compute efficiency. Crucially, their work *never* reaches the stream-aligned product teams directly. It is entirely encapsulated and distributed as a service through the platform.

### The Three Interaction Modes

*Team Topologies* defines three core interaction modes between these teams. There is not much to adapt for the agentic era; the model is mostly suitable out-of-the-box. Here is the map:

* **Facilitating:** The enabling team *facilitates* the stream-aligned teams. This is a temporary, hands-on interaction oriented entirely toward building capability. The enabler coaches the business team on how to safely interact with the platform and the AI orchestrator.
* **X-as-a-Service:** The platform team delivers its capabilities to the stream-aligned teams strictly via self-service. This is the target state and the sole engine of scaling. It eliminates blocking tickets and manual handoffs. It provides product teams with true *autonomy* (using standardized tools effectively) rather than *independence* (building shadow IT in isolation).
* **Collaboration:** The complicated subsystem team *collaborates* closely with the platform team to integrate deep technical capabilities. This is a high-bandwidth, synchronous interaction justified during the initial build phase. However, as the system matures, this too must evolve into X-as-a-Service: AI efficiency and infrastructure must eventually become consumable APIs, not permanent construction sites.

---

## Making the model stick

Defining teams and their interactions is only the first step. A static blueprint will inevitably fail. As with any framework, adapting it to your context is always better than blindly adopting it (in my experience, blind adoption is the Dire Straits anti-pattern: money for nothing). For this organizational structure to survive contact with reality, it cannot be a one-off reorganization. It must be treated as an ongoing journey toward a specific target of maturity.

### The Journey Toward Autonomy


The enabling team's role is explicitly designed to shrink. That is the ultimate sign the system is working.
![Evolution of Interaction Modes](/assets/team-topologies-agentique/octo_team_topologies_evolution.en.svg)
*The three interaction modes evolve together: facilitating fades, collaboration gives way to X-as-a-Service, which becomes the dominant mode.*

This evolution follows two parallel axes that must be carefully distinguished: **team maturity** and **platform maturity**. They reinforce each other, but they do not progress at the exact same pace.
**In the beginning:**
* **Team side:** Product teams are just discovering context packaging and how to operate the AI orchestrator. The enabling team is omnipresent, often embedded directly within the product teams.
* **Platform side:** Capabilities are basic (a few guardrails, partial documentation, a fragile deployment pipeline). The platform does not yet meet the maturity criteria defined earlier.


**During maturation:**
* **Team side:** Product teams master the fundamentals. The enabling team starts to decouple and structure its own goals. Its interventions become highly targeted (e.g., advising on a specific guardrail or troubleshooting a complex prompt).
* **Platform side:** Guardrail coverage expands, self-service adoption grows, and documentation solidifies. Deep collaboration with complicated subsystem teams begins to crystallize into integrated, stable services.


**At full autonomy:**
* **Team side:** Product teams are completely autonomous. Enabling becomes optional, functioning more like an on-demand consultancy for edge cases.
* **Platform side:** Self-service is comprehensive, deterministic guardrails cover all critical dimensions, and full audit traceability is active. *X-as-a-Service* is the absolute dominant interaction mode.

> **The enabling team disappears because it succeeds, not because it fails.**

#### A Concrete Example: The Target State
 
Let’s look at how this plays out once the organization reaches **full autonomy**.
The marketing team wants to generate a new landing page. They supply the **dynamic context** (campaign intent, key messaging). The platform automatically injects the **systemic context** (brand guidelines, UI components, accessibility rules), and its deterministic guardrails verify brand consistency and security.

*Notice how the roles have evolved:* The enabling team, who heavily supported marketing in the early days, completed their core training three months ago. Today, they only drop in for occasional edge-case support.

Meanwhile, the complicated subsystem team's deep technical expertise operates silently under the hood. 
During the build phase, they collaborated with the platform team to design a sovereign AI architecture (optimizing GPU allocation and handling concurrent accesses). They also helped the platform identify exactly which parts of the massive brand context should be frozen. Because of this, the context is now served directly from the KV cache rather than being recomputed on every prompt iteration, drastically reducing the cost per generation.

**The result:** A compliant, secure, fully deployed landing page—without the marketing team ever needing to know what a CI/CD pipeline is.

**The Flip Side:** A team that never "sees" the underlying protection mechanisms loses the ability to judge their relevance or report false positives. Therefore, guardrails must be **transparent in their decisions** (providing clear error messages and audit trails), even if they remain entirely opaque in their implementation.

### Application Governance: Preventing Shadow IT at Scale

A solid team structure is not enough. If we truly empower business teams to build and deploy, we must address the foundational question of governance: **who decides an application has the right to exist?** And who manages its lifecycle (its technical debt, its cumulative cost, and its eventual deprecation)?

Without governance, agentic production just creates *industrialized shadow IT*. The platform acts as the enforcement engine for this governance (similar to the concept of *computational governance* in Data-Mesh). By centralizing deployment, monitoring, and usage metrics, it provides **systemic visibility** across the entire application portfolio.

Through the platform, we can track active applications, identify abandoned ones, and trigger automated deprecation. If an older application suddenly fails a newly updated security guardrail, it is immediately flagged, not silently forgotten.

**Ease of production must be matched by ease of oversight.** If the platform makes it trivial to spin up an application, it must make it equally trivial to audit how many exist, who owns them, and what they cost.

### The Graduation Path: From Specific to Systemic


Specific guardrails will move toward **an** off-the-shelf product when **commoditized**. But there is one final, essential question regarding the platform's lifecycle: **when does a local "product" guardrail become a global "platform" guardrail?**

**The Rule of Three:** Let’s say the E-commerce team implements a Personally Identifiable Information (PII) guardrail (a security filter designed to automatically scrub sensitive customer data like shipping addresses or credit card partials) before it gets sent to an external LLM (a key topic when operating in Europe, for example).

Three months later, the Loyalty Program team builds a similar filter to mask customer names and birth dates when processing feedback logs. Then, the Store Operations team does the same for email addresses in Click & Collect order issues.

This is the graduation signal. Borrowing Martin Fowler’s *Rule of Three*[^fowler], once a guardrail is duplicated across three distinct teams, it becomes a candidate for systemization. The platform team abstracts the data-scrubber, makes it configurable (e.g., ensuring standard compliance with GDPR or CCPA rules), documents it, and exposes it globally as a unified "Customer Data Anonymization" service. *(Personal note: I built a POC about A2A communication using this exact example; more on this in an upcoming article).*

In **this** context, stream-aligned product teams surface the recurring needs. Enabling teams spot cross-cutting trends during their engagements (Communities of Practice or Special Interest Group meetings are excellent venues to spot these). The platform team evaluates and integrates. Autonomy does not mean isolation: product teams act as the frontline sensors that feed the platform's evolution. The platform remains a *living product*[^skelton2019] with its own backlog, Product Owner, and iteration cycles.

**The Bottleneck Paradox:** A risk emerges here. The more successful the platform becomes, the more the Platform PO inherits three massive burdens: the technical backlog, the guardrail graduation process, and application portfolio governance. Paradoxically, if one person must manually arbitrate the flow of needs from dozens of agent-driven teams, we recreate the exact cognitive bottleneck this entire model was designed to prevent.

**The answer lies in tooling and, even more so, in cognitive automation.** A human cannot manually harvest and categorize this data at an agentic scale. The platform must automate the detection of graduation candidates, usage-based prioritization, and portfolio tracking. The PO arbitrates the data; they do not carry the burden of collecting it.

Without this continuous graduation mechanism, the platform stagnates. Product teams will endlessly reinvent the same solutions, guardrails will remain localized, and systemic reliability will silently erode.

---

## Operational synthesis

Now we have all the principles in place, it remains to turn them into actionable rules.

### Who owns what?

Every capability must have one owner (accountable for the outcome) and can have multiple contributors (who co-build without carrying accountability). Establishing this strict boundary is non-negotiable.

![Responsibilities by team](/assets/team-topologies-agentique/octo_team_topologies_responsabilites.en.svg)
*Each capability has an owner and contributors. Empty cells mean "not involved." Specifics stay with product teams, systemics with the platform.*

The guiding principle is straightforward: anything specific (dynamic context, product-level guardrails, domain knowledge) stays with the stream-aligned teams because the business owns the domain. Anything systemic (global context, organizational guardrails, core tooling) is abstracted and maintained by the platform team.

The AI orchestrator is the perfect illustration of this split: the platform provides the engine (the systemic tool), but the product team drives it (the specific configuration).

### A Conviction

The agentic platform answers the *what*. *Team Topologies* provides the framework to set up the *who* and the *how*. Treat this not as a rigid partition, but as a strict **distribution of accountability**: business intent and domain context belong to the product team, while production reliability belongs to the platform.

Other organizational models will certainly emerge. But building a platform without clearly defining these boundaries and interactions will inevitably turn that platform into a massive bottleneck.
This model is not an off-the-shelf framework; it is a mental model for organizing **a team of teams building multiple agentic applications in parallel**. As an empirical heuristic (not an absolute rule), once you reach three to five product teams, the cumulative cost of reinvention (recreating context, rewriting guardrails, fixing localized inconsistencies) far outweighs the cost of investing in a shared platform.

Applied to agentic engineering, *Team Topologies* provides a structure that:

* Enables **business teams to act as stream-aligned teams** without forcing them to learn how to code.
* Ensures that **systemic capabilities are abstracted and scaled**, rather than reinvented project by project.
* Acknowledges that the **gap between business intent and production reality is vast**, bridging it with enabling teams—not with hope.
* Isolates **deep technical complexity** into specialized subsystem teams that power the platform under the hood.


What has changed in our industry is that producing an application **no longer requires writing code**. Instead, it demands a new set of skills: articulating intent, structuring dynamic context, and steering an AI orchestrator. What *has not* changed is that you still need **guardrails, consistency, and reliability** to survive in production.


This brings us back to the core problem raised at the outset. Today, the software engineer carries the entire "anticipation burden" because they know exactly which questions the agent will fail to ask. Tomorrow, as the platform deterministically absorbs that burden, the technical barrier to entry plummets. A business team can safely deploy software because the platform anticipates the systemic risks on their behalf.

**The shift of the application producer (from software engineer to business expert) is not a futuristic assumption. It is the direct, measurable consequence of the platform absorbing the anticipation burden.**

The developer does not disappear. Their role shifts upstream. They no longer build the application: they build the engine that empowers everyone else to build it.

---

## Sources

[^skelton2019]: Matthew Skelton, Manuel Pais — *"Team Topologies: Organizing Business and Technology Teams for Fast Flow"*, IT Revolution, 2019.

[^fowler]: Don Roberts, cited by Martin Fowler — *"Refactoring: Improving the Design of Existing Code"*, Addison-Wesley, 1999. "Three strikes and you refactor." The principle is transposed here from code refactoring to guardrail graduation.
