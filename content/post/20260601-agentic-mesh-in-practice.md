---
title: "The Agentic Mesh in Practice: Anatomy of an Agent-Product"
date: 2026-05-31T14:00:00+02:00
lastmod: 2026-05-31T14:00:00+02:00
images: [/assets/agentigslide/pipeline-multi-agents.en.svg]
draft: false
keywords: ["agentic mesh", "multi-agent", "Go", "Google Slides", "A2A"]
summary: "How I built a multi-agent system in Go to generate Google Slides presentations from pre-formatted templates (not just to present, but to convince). A hands-on account illustrating the principles of the agentic mesh."
tags: ["AI", "agents", "architecture", "agentic-mesh", "Go", "experience-report"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

I am a consultant, and I regularly build presentations with Google Slides. My communication team has created dozens of pre-formatted templates (slides designed to **convince**, not just to present). The problem: choosing the right slides to illustrate the right narrative takes time, and filling them in mechanically adds no value. I built a multi-agent system to automate that part and focus on what matters: the narrative and making it my own.

This project (**[agentigslide](https://github.com/owulveryck/agentigslide)**) is also a concrete application of the agentic mesh principles I described in [the previous article](/2026/05/31/the-agentic-mesh-cognitive-automation-at-scale.html). There, I laid out a conceptual framework: four pillars, a three-phase trajectory, convictions about what an agent-product should be. Here, I tell the story of how these principles materialized in real code, architecture decisions documented through 16 ADRs, and a tool that works in production.

> **Note.** This article was co-written with an AI. I'm at the helm: I set the direction, the ideas, and I review the entire document. The actual writing was done by AI. My goal is to share these ideas to open a discussion, not to write a technical masterpiece that becomes a stylistic reference. This version is designed for humans; if you prefer a version suited for AI consumption, the [markdown source](https://raw.githubusercontent.com/owulveryck/owulveryck.github.io/refs/heads/master/content/post/20260601-agentic-mesh-in-practice.md) is available.

---

## Convince, not present

Every slide creation tool solves the wrong problem. Gamma, Beautiful.ai, Pitch: they generate visually correct slides. But they produce **presentations**, not **persuasion presentations**. The difference is fundamental.

Convincing starts with a **narrative**. Slides are a complement to it, not a restatement. A consultant preparing a pitch doesn't start from a slide generation tool, they start from their argument, the structure of their demonstration, the ideas they want to anchor in the listener's mind. Slides are just the visual aid for that reasoning.

In this context, my communication team produced a catalog of ~300 pre-formatted slides in Google Slides (slides designed by visual communication professionals, crafted to illustrate specific IT consulting concepts: persuasion slides, framing slides, comparison slides, process slides). This catalog is a **brand asset** that encodes visual and rhetorical conventions honed over years.

The **consultant's value** lies in choosing the right slides from this catalog to illustrate the right concepts in their argument, then filling them with their content. The choice is strategic, the filling is mechanical. I wanted to automate the mechanical to free up the strategic.

> **The consultant's value lies in the choice and the personalization, not in the mechanical filling.**

In the vocabulary of the agentic mesh, this is **Pillar 2: Domains**. The agent I built belongs to the consulting domain: its intention is expressed in business vocabulary (structuring a pitch, choosing persuasion slides), not in generic technical terms.

---

## The strategic landscape: a Wardley map

Before writing any code, I mapped the landscape with a **Wardley map** to understand where value lies and what strategic moves are possible.

The value chain reads top to bottom: the **consultant** must **convince a client**. To do this, they produce a **narrative** that they externalize into a **structured brief** (a markdown file). This brief feeds the **agentic orchestration** which draws from the **pre-formatted slide catalog** via a **semantic index**, to produce a presentation through the **Google Slides and Drive APIs**.

![Wardley map of the agentigslide landscape](/assets/agentigslide/wardley-agentigslide.en.svg)

### What the map reveals

**The catalog is the invisible moat.** Positioned in the Custom phase, it is relational capital (not replicable without the communication expertise that produced it). A competitor could reproduce the agentic architecture, deploy the same LLM models, but could not copy this catalog without copying the years of expertise that shaped it.

**Timing is critical.** Agentic orchestration is leaving the Genesis phase to enter Custom. Generic generation tools (Gamma, Beautiful.ai) are under Red Queen pressure: they must constantly evolve just to not fall behind. Their natural direction is to absorb agentic capabilities. The window to occupy the "persuasion slides" niche is **12 to 18 months**.

**Human-in-the-loop is a strategic choice, not a limitation.** The consultant reviews and appropriates the generated slides. It is the step that transforms a correct presentation into a convincing one. It is not a stopgap waiting for the AI to be "good enough": the human who presents is responsible for what they present.

**Non-intrusive by design.** The communication team continues using Google Slides without learning any new tool. The system automatically understands template structure and adapts to it. Frictionless adoption = real adoption.

### Identified strategic moves

- **Land-grab**: own the "persuasion" vocabulary before generic tools capture the niche.
- **Strategic open-source**: publish the agentic framework (the plumbing will commoditize) to shift attention to the catalog (the real moat).
- **Strangler-fig on the brief**: an upstream structuring agent helps the consultant structure their argument: it gradually replaces the manual production of the brief, not the discourse itself.
- **ILC on the catalog**: each generation produces implicit signals (which slides are chosen, which are never selected, which concepts lack a slide). The communication team uses these insights to improve the catalog (a usage → production feedback loop).

---

## Phase 1: the monolithic prototype

My first system was a **monolithic pipeline**: a single Claude call received the complete template catalog (~60 slides, ~15-20 KB of text) and the user's request, and had to analyze the structure, select templates, fill the text content of each field, and maintain overall presentation coherence in a single pass.

It worked. At least well enough to **validate the fundamental hypothesis**: yes, an LLM can choose the right slides from a catalog and fill them in relevantly.

This is exactly **Phase 1** of the agentic mesh trajectory: *"The objective of this phase is not to build a good agent. It is to validate the fundamental hypotheses."* And that is what the prototype did: it validated that cognitive automation is possible on this use case.

But four **structural limitations** emerged quickly:

1. **No feedback loop.** If the model picks a 6-field slide for 3 bullet points, or exceeds a character limit, the error goes to production. Single-pass reasoning, no safety net.
2. **No parallelism.** Writing content for slide 3 and slide 7 are independent tasks, but everything goes through a single sequential call.
3. **Token inefficiency.** Every subtask receives the entire context (full catalog) when only a fraction is relevant.
4. **Tight coupling.** A 2-field cover slide doesn't need the same model as a 6-field content slide. But everything goes through the same call.

The prototype had served its purpose. It was time to move to engineering.

---

## Phase 2: from prompt to product

### The multi-agent architecture

The first architecture decision ([ADR 001](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/001-agentic-architecture.md), May 5, 2026) transformed the monolithic pipeline into **four specialized agents** coordinated by a pure Go code orchestrator (no AI in the orchestration). This is the **coordinator-subagent** (hub-and-spoke) pattern, the standard architecture for multi-agent systems: a central coordinator decomposes the task, delegates to specialized agents with **isolated context**, then aggregates their results.

**The Outliner** analyzes the user's request and produces a structured presentation plan. Crucial point: it does **not** receive the template catalog. This **context isolation** is deliberate: it forces the reasoning *"what do we need?"* before *"what do we have?"*, avoiding availability bias. This is a fundamental principle of agentic architectures: subagents do not automatically inherit the coordinator's context, each only receives what it needs.

**The Selector** matches the needs identified by the Outliner against the available templates in the catalog. It works with the `itemCount` and `maxItemLength` context provided by the Outliner to make informed choices.

**The Writers** generate the textual content for each slide, in parallel. Each Writer receives **only one tool** (`produce_slide_content`) with a dynamic JSON schema generated from the template fields (**scoped tool access** principle: restrict tools to the bare minimum to maximize selection reliability). The model is selected based on complexity: **Haiku** for simple slides (≤ 2 fields), **Sonnet** for complex ones.

**The Reviewer** is an **independent review instance**: it does not have the Writers' reasoning context, it only receives the assembled plan and the original request. This is not a self-review (where the same model validates its own choices), it is a separate agent bringing a fresh perspective. It uses **Opus with extended thinking** for in-depth analysis. When it detects problems, it sends structured feedback (`ReviewIssue[]`) via the **retry-with-error-feedback** pattern: specific errors are injected into the affected Writer's prompt, not a vague "try again" (maximum 2 iterations to bound cost).

![agentigslide multi-agent pipeline](/assets/agentigslide/pipeline-multi-agents.en.svg)

This is **Phase 2** of the agentic mesh: *"Each sub-agent is specialized, uses the model adapted to its task, has its own tools and feedback loops."*

### Why native code and not directives

I could have used Claude Code with MCP servers, or Anthropic's Agent SDK. I chose native Go, and this choice illustrates the founding conviction of the agentic mesh: **a production agent is a software engineering product, not an assembly of directives**.

The fundamental distinction is between **programmatic enforcement** and **prompt guidance**. The former provides deterministic guarantees (validation blocks the pipeline if selected templates don't exist), the latter provides probabilistic compliance (the model will follow instructions most of the time, but not always). When non-compliance has visible consequences (an incoherent presentation delivered to a client), programmatic enforcement is not optional.

| Criterion | Native Go | Off-the-shelf system + directives |
|-----------|-----------|-----------------------------------|
| Fine-grained parallelism | Goroutines + semaphore, concurrency control | Limited, sequential or simple parallel |
| Feedback loops | Typed `ReviewIssue[]`, targeted retry on subset | Conversational, fragile over time |
| Structured outputs | Strict JSON schema, programmatic validation | Implicit, model-dependent |
| Prompt caching | Vertex AI ephemeral cache, shared between Writers | Not available or not shared |
| Inter-step state | Typed, mutex, testable | Lives in conversational context |
| Observability | Per-agent metrics, complete issue log | Limited, aggregated |

Engineering is not only found in the tools the agent uses: it is also found in how the agent itself is built, tested, deployed, observed, and governed.

### Engineering decisions documented through ADRs

The project accumulated **16 Architecture Decision Records** in 4 weeks. Each ADR is a deliberate engineering decision, documented with its context, alternatives considered, and consequences. This is not after-the-fact documentation, it is **governance in action**. Here are the most significant ones, with the agentic principle they illustrate:

| ADR | Decision | Agentic mesh concept illustrated |
|-----|----------|----------------------------------|
| [001](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/001-agentic-architecture.md) | Multi-agent architecture | Phase 2: from prototype to product |
| [002](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/002-prompt-caching.md) | Explicit prompt caching via Vertex AI | Cost control for an agent-product |
| [004](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/004-prompt-externalization.md) | Prompt externalization (`go:embed`) | Versioning and testability |
| [005](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/005-interactive-chat-mode.md)-[006](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/006-default-agent-chat-mode.md) | Interactive chat + default mode | Structural human-in-the-loop |
| [007](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/007-a2a-architecture.md) | A2A (Agent-to-Agent) architecture | Interoperability contract, Agent Card |
| [009](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/009-diagram-agent.md) | Designer agent for diagrams | Sub-domain specialization |
| [010](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/010-edit-existing-presentations.md)-[012](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/012-edit-post-processing.md) | Orchestrated edit pipeline | Functional domain extension |
| [015](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/015-agent-memory-learning.md) | Per-agent learning memory | Governed continuous improvement |
| [016](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/016-format-agent.md) | Deterministic FormatAgent | Automated governance without LLM |

---

## The 7 affordances, concretely

In the agentic mesh article, I defined **7 affordances** that an agent must offer to be a true product. Here is how they materialized in agentigslide.

**1. Expose decisions and actions.** Each agent exposes its capabilities via a strict JSON schema enforced by Claude's `tool_use` mechanism with forced `tool_choice` (the model *must* call the specified tool, eliminating parasitic text responses). With [ADR 007](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/007-a2a-architecture.md), each agent also exposes an **A2A Agent Card** (a self-descriptive manifest published at `/.well-known/agent-card.json`).

**2. Consume context.** The agents consume three types of context: the catalog's semantic index (built once, reused for every generation), template-specific instructions (an optional `PROMPT.md` file), and memory files from previous executions.

**3. Reason and decide.** This is the core of the product: Haiku for simple tasks, Sonnet for complex ones, Opus with extended thinking for review. Prompts are externalized via `go:embed` and versioned with the code. The Reviewer → Writer loop is structured **prompt chaining** (each step's output becomes the next step's input with programmatic validation between steps), not an open conversation.

**4. Be discoverable.** The Outliner is already deployed as a standalone A2A server (`cmd/outliner/main.go`). [ADR 014](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/014-agent-pipeline-registry.md) proposed an agent registry pattern for dynamic composition.

**5. Manage its lifecycle.** Models are configurable per agent via environment variables (`AGENT_OUTLINER_MODEL`, `AGENT_WRITER_MODEL`...). The original monolithic mode is preserved as a fallback (backward compatibility without debt).

**6. Trace decisions.** Each agent reports its tokens (input, output, cache read, cache creation), its duration, and a complete issue log. The Reviewer's extended thinking is traced. The FormatAgent ([ADR 016](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/016-format-agent.md)) logs every deterministic correction applied.

**7. Be governable.** `MaxReviewRetries` bounds correction iterations. `enforceMaxChars()` acts as a **post-execution hook**: it intercepts Writers' output and truncates fields that exceed the template limit, providing a deterministic guarantee that the prompt alone cannot offer (*trust but verify*). Agent memory ([ADR 015](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/015-agent-memory-learning.md)) is validated by the human before writing: the agent proposes guidelines, the user confirms.

---

## Human-in-the-loop: a choice, not a crutch

The default mode of agentigslide is **interactive chat** ([ADR 005](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/005-interactive-chat-mode.md)-[006](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/006-default-agent-chat-mode.md)). The consultant describes what they want, the Outliner proposes a structure, the consultant refines in conversation, and only when the plan is validated does the generation launch.

This is not a crutch. It is a deliberate choice, rooted in a conviction: **the person who presents is responsible for what they present**. The agent produces professional-quality material, but it is the consultant who examines it, adjusts it, and adds the personal touch that makes the difference between a generic presentation and a convincing one.

In the agentic mesh, I wrote: *"H-A (Human-Agent) sequences are not exceptions or safeguards, they are a feature of the mesh."* This is exactly what happens here. The human is not in the loop because the AI is not good enough. They are in the loop because that is their irreplaceable contribution.

---

## The catalog: the invisible moat

Let me stress a point that the Wardley map makes obvious: **the code is not the moat**. The multi-agent architecture, the goroutines, the prompt caching: all of this will commoditize. A competitor could reproduce the entire technical system.

What they could not reproduce is the **catalog of 300 slides** designed by visual communication professionals. This catalog is not a collection of templates, it is a brand asset that encodes visual and rhetorical conventions built over years. It is relational capital, not technical capital.

And the **non-intrusiveness** doctrine reinforces this moat: the communication team continues working in Google Slides without learning any new tool. The system automatically analyzes each slide with Claude Vision, understands its structure, identifies editable fields, and builds a semantic index. Adoption is frictionless because the tool adapts to humans, not the other way around.

---

## Toward Phase 3: from agent-product to mesh

### A2A: each agent exposes an Agent Card

[ADR 007](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/007-a2a-architecture.md) (May 9, 2026) was the turning point. The Go pipeline works well, but it has three structural limitations that no optimization can solve:

1. **Agents cannot orchestrate other agents.** The Selector cannot dynamically decide to call a layout agent then a design agent.
2. **The pipeline is closed to external extension.** Adding a new agent requires modifying, recompiling, and redeploying the binary.
3. **The deployment unit is the entire binary.** No independent lifecycle per agent.

The solution: the **A2A** protocol (Agent-to-Agent, Google, 2025). Each agent exposes an Agent Card and accepts Tasks via a standardized REST API. The Go orchestrator remains (deterministic, predictable) but calls agents via A2A rather than Go functions. This directly implements the **interoperability contract** and the **Agent Card** described in the agentic mesh article.

### From closed pipeline to composable network

The deepest paradigm shift concerns the catalog. Today, the catalog is a **hard constraint**: if no template fits, the Selector can do nothing. With A2A, it becomes a **smart default with creative fallback**: when no template fits, the Selector orchestrates sub-agents (layout agent, design agent, visual validation agent) to create a slide from scratch using design primitives.

This change has a non-trivial precondition: the communication team's **visual charter**, currently implicit in the slides, must become explicit, versioned, testable. The communication team no longer produces only slides, they produce **design primitives** and **composition rules**. This is a fundamental shift in their contribution.

### Memory as incremental governance

[ADR 015](https://github.com/owulveryck/agentigslide/blob/main/docs/adr/015-agent-memory-learning.md) introduced **per-agent learning memory**. Each agent has a Markdown file per template, stored alongside the template and versioned with git. Guidelines are actionable: *"On slide #42, never exceed 120 characters in the title field (text overflows systematically)."*

At the end of the pipeline, if errors were detected, the system synthesizes guidelines and **proposes them to the user** (no automatic writing). The human confirms before memory is enriched. This is affordance 7 (being governable) applied concretely: the agent improves within boundaries defined by the human.

---

## What the field teaches the theory

After 4 weeks of building, 16 ADRs, and a system that works in production, here is what the agentigslide project confirmed (or taught me) about the principles of the agentic mesh.

**The real MVP is the platform.** The slide catalog, the semantic index, the Google APIs: that is the foundation that unlocked value. The agent only came after, and it could not have done anything without this foundation.

**The agent is not a directive, it is an engineering product.** The comparison table in the [RATIONALE.md](https://github.com/owulveryck/agentigslide/blob/main/RATIONALE.md) demonstrates this unambiguously. Off-the-shelf systems are a good starting point (I used one for Phase 1). But as soon as you need controlled parallelism, typed feedback loops, inter-step validation, and shared prompt caching, you are doing software engineering, not assembling directives.

**ADRs are governance in action.** 16 documented decisions in 4 weeks. The context, the alternatives considered, the consequences: it is all there. Governance is not in a wiki, it is in the decision trail.

**The moat is in the domain, not in the technology.** The catalog is the moat, not the Go code. This is **Pillar 2** of the agentic mesh at work: *"the responsibility for the agent falls on the business domain that leverages its value, not on a central AI team."*

**The window is narrow.** Generic tools will absorb agentic capabilities. The window to occupy the "persuasion slides" niche is 12 to 18 months. The land-grab is urgent.

> **The agentic mesh is not a diagram to hang on the wall. It is a trajectory to follow, and agentigslide is a waypoint on that path.**
