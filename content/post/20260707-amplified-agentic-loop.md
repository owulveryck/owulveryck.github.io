---
title: "The Amplified Agentic Loop: Guardrails as Accelerators"
slug: "amplified-agentic-loop"
date: 2026-07-07T10:00:00+02:00
images: [/assets/amplified-agentic-loop/amplified-agentic-loop.svg]
draft: false
summary: "Today, our interaction with coding agents is asymmetric and front-loaded: all the control, context, and governance rules are concentrated in the initial prompt, and the remaining guardrails only kick in after the fact (blocking a commit, rejecting a build). This article flips that paradigm. By injecting governance, architectural context, and safety directly inside each step of the agentic loop (amplified planning, contextual in-tool execution, retroactive observation), the platform turns guardrails from brakes into accelerators. It closes with a working proof of concept: a Platform Planning Gateway that lints agent plans deterministically and issues capability tickets that smart tools verify before acting."
tags: ["architecture", "agents", "platform", "team-topologies"]
categories: ["dev"]
author: "Olivier Wulveryck"
toc: false
comment: false
mathjax: false
---

## Introduction

This article closes a trajectory started three posts ago. In [See, Act, Correct](/2026/06/04/see-act-correct-three-levers-for-working-with-a-code-agent.html), I described the three levers that turn a code agent from a gadget into a production tool, and introduced a grid to sort durable platform investments from temporary crutches. In [Who does what?](/2026/06/24/who-does-what-team-topologies-for-the-agentic-platform.html), I mapped the agentic organization onto [*Team Topologies*](https://teamtopologies.com/). In [Codifying the Rules](/2026/07/02/sdlc-team-topologies.html), the enabling team baked its guardrails into the platform and vanished.

The trajectory is simple to state: to move from individual *vibe coding* (a developer "vibing" with their LLM) to industrialized, at-scale software engineering, the **platform** (in the [Platform Engineering](https://platformengineering.org/) sense: an internal product that reduces the cognitive load of stream-aligned teams) must become the conductor of the agentic workforce.

Here is the problem this final part addresses. Today, all the control is concentrated at the *beginning* of the interaction, and the safety nets only react at the very *end*:

![The agentic lifecycle and its injection points: initial context under Intent declaration (the only lever of the front-loaded model), and the platform serving dynamic guardrails, in-tool context, and the agentic feedback loop into Planning, Execution, and Observation](/assets/amplified-agentic-loop/lifecycle-injection-points.en.svg)

The left column is where today's tooling operates: a heavy initial context, then nothing until a CI gate says "no". The three blue channels are what this article designs: governance, architectural context, and safety served by the platform **inside** each step of the agentic loop.

Several tools already cover parts of this space. [NeMo Guardrails](https://docs.nvidia.com/nemo/guardrails/latest/) distributes enforcement across input, dialog, execution, and output rails. [POLARIS](https://arxiv.org/abs/2601.11816) (January 2026) adds typed planning gates with policy-compiled guards. [Agentic JWT](https://arxiv.org/abs/2509.13597) (September 2025) proposes cryptographic action scopes that bind agent actions to user intent. Each addresses one phase. What is absent — and what a [May 2026 survey on trustworthy agentic AI](https://arxiv.org/abs/2605.23989) explicitly lists as an open challenge — is the integrated architecture: live architectural context at planning time, first-class capability tokens at execution time, and structured intent-gap feedback at observation time, combined into a single design where each phase amplifies the next. This article proposes one answer; the companion repository [poc-agentic-platform](https://github.com/owulveryck/poc-agentic-platform) is a working proof of concept, not a finished product.

{{< scrollytelling svg="/assets/amplified-agentic-loop/amplified-agentic-loop.svg" >}}

{{< scrollytelling-step phase="1" id="phase-1" >}}

## The agentic loop, a reminder

The engine has not changed since the [previous article](/2026/07/02/sdlc-team-topologies.html): an agent **captures** an intent, **plans** the steps, **acts** through tools, and **observes** the results. When the observation is unsatisfactory, the agent re-plans and iterates: the loop is self-correcting.

This loop is the delivery engine of modern software. Everything that follows in this article is about one question: *where does governance live in this loop?*

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="2" id="phase-2" >}}

## The status quo: everything is front-loaded

Look at how we steer agents today. Governance rules, architecture context, security policies: everything is injected **before** the loop starts, through the initial prompt and global instruction files ([`CLAUDE.md`](https://code.claude.com/docs/en/memory), [`.cursorrules`](https://cursor.com/docs/context/rules), [`copilot-instructions.md`](https://docs.github.com/en/copilot/customizing-copilot/adding-repository-custom-instructions-for-github-copilot), system prompts). These files are the industry's current answer to steering, and they share one property: they act once, at minute zero. The notebook on the diagram marks the spot: the initial context is assembled at **Capture**, once, and everything the stack dumps up-front lands in it.

The interaction is asymmetric: all the control is spent up-front. Once the loop starts spinning, the agent is on its own. Nothing guides it *while* it plans, *while* it executes, *while* it observes. We write two hundred lines of instructions and then hope.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="3" id="phase-3" >}}

## The late gate: control after the work is done

The remaining control arrives at the worst possible moment: the end. The commit that breaks the CI is rejected. The forbidden API is flagged in code review. The non-compliant color ships to staging and bounces back.

These late gates do their job (nothing broken reaches production), but look at the cost. The agent has already burned its tokens. The work is done and must be *re*done. Each bounce is a full trip around the loop, plus a human interruption. This is the opposite of what the industry calls [**shift-left**](https://en.wikipedia.org/wiki/Shift-left_testing) (moving verification as early as possible in the delivery process), and it frustrates the very flow the agent was supposed to accelerate.

But wait: is self-correction not precisely the loop's job? It is. The red arrow on the diagram has been doing exactly that since the first phase. The nuance is that **the loop can only correct what its Observe step can see.** The late gate delivers its verdict *after* the loop has exited, outside the agent's observation horizon. So the rejection cannot be absorbed as one more iteration; it re-enters at the top, as a new intent, usually carried by a frustrated human. Same information, wrong side of the loop boundary.

The failure is architectural, not technical: the control exists, it is just *in the wrong place*. It sits outside the loop, when the loop is where the work happens.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="4" id="phase-4" >}}

## The pivot: distribute governance inside the loop

**Here is the shift this article proposes.** Instead of concentrating everything at the start and controlling at the end, the platform injects governance *inside* each step of the loop:

* at **planning** time, the agent receives the architecture patterns, the critical dependencies, and the valid options, before it commits to a strategy;
* at **execution** time, the tools themselves carry the rules and validate state in real time;
* at **observation** time, the platform measures the gap between the result and the original intent and feeds it back, structured.

The initial context becomes light again: the intent, and little else. Everything the agent used to receive as a wall of up-front text is now delivered *just in time*, by the step that needs it.

Who builds this? The platform team, and this is Platform Engineering by the book: these mechanisms are designed, versioned, and served as [**internal products**](https://tag-app-delivery.cncf.io/whitepapers/platforms/) that stream-aligned teams and their agents consume without friction.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="5" id="phase-5" >}}

## Pillar 1 — Amplified planning

Do not let the agent plan in the void. At the exact moment it designs its strategy, the platform provides the enterprise architecture patterns, the critical dependencies, and the valid options.

The industry has a name for the shape of this guidance: the **Golden Path**, the paved road, the supported and blessed way of building something, popularized by [Spotify](https://engineering.atspotify.com/2020/08/how-we-use-golden-paths-to-solve-fragmentation-in-our-software-ecosystem/) and productized by tools like [Backstage](https://backstage.io/). But today's golden paths are documentation for humans. Here, they become a *planning input for the agent*: instead of discovering conventions by reading thousands of lines of code (or worse, by violating them), the agent receives the map of the terrain before it starts walking.

The nature of the guardrail changes: it guides toward the paved road at planning time, instead of saying "no" at the end.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="6" id="phase-6" >}}

## Pillar 2 — Contextual in-tool execution

The tools the platform hands to the agent (CLIs, APIs, refactoring scripts) must not be passive executors that return `exit 1` and a stack trace. **They must be smart**: encapsulate the business rules, validate state in real time, and, when something goes wrong, answer with an *actionable, semantic* error.

Compare the two experiences. A passive tool says:

```
Error: migration failed near line 12. SQLSTATE 42P07. Exit status 1.
```

A smart platform tool says:

```json
{
  "error_category": "DATABASE_SCHEMA_CONFLICT",
  "message": "Table 'payments' already exists in staging.",
  "remediation_guidance": {
    "allowed_actions": [
      "Use 'get_db_schema' to inspect the current structure.",
      "Add an 'IF NOT EXISTS' clause or rename the table."
    ],
    "context_update": "Current staging schema version: v2.4.1"
  }
}
```

The first output sends a weak model into a panic loop. The second one turns the failure into an immediate, guided self-correction: the tool behaves like a deterministic mentor.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="7" id="phase-7" >}}

## Pillar 3 — Retroactive observation

Start with the question this phase usually raises: *where does the observation get its yardstick?* Not from a side channel; the loop does not need one. **The context is born at Capture, and nowhere else.** That is what the capture step *is*: the moment the loop's working context is created and filled (the declared intent, the injected invariants). Every subsequent step *appends* to that same context: the plan, the tool calls, the raw results. By the time execution reaches Observation, the original intent is still sitting at the top of it. The context chip on the diagram shows exactly that: the yardstick flows from the captured context into the comparator, not from outside the loop.

Pillar 3 is about what observation *does* with that yardstick. It must not exist only to produce logs for humans. Its job is to **measure the gap between the execution result and the captured intent**, and to report it structured: *which* invariant was missed, *by how much*, *what context changed*. Not a raw "unsatisfactory result, try again".

And observation does not complete the intent by itself; nor does an intent gap short-cut back to planning. That inner arrow is reserved for **execution-level** errors: a tool was called wrong, the *how* failed, re-plan the step (pillar 2). An intent-level gap is something else entirely. It is **new knowledge** (the measured gap, the facts discovered along the way), and knowledge belongs in the context. So the feedback flows back to **Capture**, where the context gains those new elements, and the loop runs **another full cycle** (capture, plan, act, observe) on the enriched context, until the gap closes or a budget guardrail fires. Watch the diagram: the green feedback travels the loop's own path back to Capture, and the context chip lights up as it lands. The intent is completed by *cycling the loop*; pillar 3 only sharpens the signal each cycle converges on.

To be precise about what changed: none of this invents a new mechanism; it routes the existing signals to the right re-entry point. A **tool-level error** re-enters at *Plan*: fix the how, same cycle. An **intent-level gap** re-enters at *Capture*: enrich the what, new cycle. And the **late gate** of phase 3 re-enters nowhere: its verdict lands outside the loop's sight entirely. The difference between the three is not the arrow; it is *what the signal is, and where it gets to re-enter the loop*.

*(This pillar opens onto agentic telemetry: measuring guardrail efficiency across thousands of loop executions and feeding the next planning phases. That is the subject of a future article; the proof of concept below deliberately stops at pillars 1 and 2.)*

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="8" id="phase-8" >}}

## The vision: guardrails as accelerators

Put the three pillars together and the paradigm has flipped. The platform validates the agent's assumptions at every step, which gives it the confidence (the *super-powers*) to execute complex modifications safely. The result exits the loop already compliant: there is no gate at the exit because the gates did their work inside.

This is the philosophical core of the article: **a well-placed guardrail is not a brake, it is an accelerator.** Formula One cars corner faster *because* the track has barriers and the car has traction control. Remove them and every lap gets slower and more careful, not freer.

A fair objection at this point: "you have just described more scaffolding around the model — will any of it survive the next model generation?" That objection deserves a precise answer, and it needs a grid.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="9" id="phase-9" >}}

## Two orthogonal axes: durable assets vs. scaffolding

In [See, Act, Correct](/2026/06/04/see-act-correct-three-levers-for-working-with-a-code-agent.html) I introduced the distinction between **compensatory** systems, which work around the current limitations of a model and lose value as models improve, and **amplifier** systems, which build on the model's strengths and *gain* value as models improve. The test is one question: *would this artifact be more useful, or useless, if the model were twice as intelligent tomorrow?*

Cross it with the implementation axis, **declarative** (text, simple, fragile) versus **programmatic** (code, deterministic, robust), and you get the grid on the left. Read it as a decision tool, not a taxonomy: take the artifact you are about to build, ask the 2× question, route it into its column, and act on the verdict at the bottom (fund the right column; tag the left one and give it a sunset date).

The crucial insight, and the one I initially got wrong when designing what follows: **whether a mechanism blocks or guides has nothing to do with whether it is compensatory or amplifying.** A 200-line prompt enumerating forbidden files is guidance, and it is pure compensatory debt: it fights the model and rots as models improve. A hard gate that verifies a *semantic invariant* ("a schema migration precedes the code that uses it") is a blocker, and it is a durable amplifier: a smarter model satisfies it more elegantly, and never makes it useless.

So the question for every platform investment below is never "does it constrain the agent?" but "does it age well?". The bottom-right cell (programmatic amplifiers) is where the next three sections live.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="10" id="phase-10" >}}

## Materializing Pillar 1: the Platform Planning Gateway

**This is the first concrete proposal of this article.** The Platform Planning Gateway (PPG) intercepts the planning step with two moves: one soft, one hard.

**The soft move, `enrich()`: fetch the rules that apply, before planning.** Think of the fifteen-minute conversation with the staff architect before you start a piece of work (*"anything I should know before touching the payment code?"*), automated. Before planning, the agent asks the platform exactly that question: it sends its **intent** and its **repository context**, and the gateway answers with the **architectural invariants** that apply, retrieved from the organization's [Architecture Decision Records](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions). Concretely:

```json
// agent → platform
{ "intent": "Add the Seka payment method to checkout",
  "repository_context": { "name": "checkout-service", "tech_stack": ["Go"] } }

// platform → agent
{ "amplifier_context": { "architectural_invariants": [ {
    "adr_id": "ADR-042",
    "invariant": "Every outbound call to a third-party service MUST go
                  through the corporate security egress proxy." } ] } }
```

How did the gateway know ADR-042 was relevant? Each ADR declares its own *scope selectors* (here, the keyword `payment` matched the intent): keyword matching in the PoC, semantic retrieval in production. The agent injects the returned invariants into its planning context and reasons over them: its plan now includes the proxy step *before* a single line is written.

Note the two things `enrich()` deliberately does **not** do. It does not enforce anything: it *advises*. Enforcement comes next, in the hard move. And it returns no recipe ("modify file X at line Y"), only semantic invariants: the recipes are exactly what a better model derives on its own. Architects write the invariants and the selectors in the ADRs; the gateway only retrieves. That division of labor is what makes `enrich()` a declarative amplifier: a smarter model exploits the same invariants better.

**The hard move: `lock_in_plan()`.** The agent must materialize its plan as a structured document (a JSON contract, not free text) and submit it. The gateway runs a **plan linter**: deterministic code, deliberately *not* an LLM, in the spirit of policy-as-code engines like [Open Policy Agent](https://www.openpolicyagent.org/). A non-conforming plan is rejected 100% of the time, reproducibly, with semantic violations:

```json
{
  "status": "PLAN_REJECTED",
  "violations": [{
    "policy_id": "go_tests_present",
    "message": "The plan must contain a 'go test' step for a Go stack.",
    "nature": "amplifier"
  }]
}
```

The agent self-corrects, in practice in one or two iterations even with a small model: the linter redresses it like a compiler would.

When the plan passes, the gateway issues the **capability ticket**: an ephemeral signed [JWT](https://datatracker.ietf.org/doc/html/rfc7519) embedding the plan's fingerprint and the [least-privilege](https://en.wikipedia.org/wiki/Principle_of_least_privilege) scope derived from it (which files may be modified, which tools may be called). This is [capability-based security](https://en.wikipedia.org/wiki/Capability-based_security) applied to the agentic loop, and it is the key that unlocks the next phase.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="11" id="phase-11" >}}

## Materializing Pillar 2: Smart Tools honor the ticket

The ticket would be theater if the tools did not check it. So every Smart Platform Tool starts with the same three-step contract:

1. **Verify the ticket, in-tool.** Agentic drift happens *during* execution: the agent decides mid-flight to touch a file it never planned to touch. The tool decodes the ticket, compares the requested action to the locked scope, and refuses deterministically: `OUT_OF_PLAN_SCOPE`, nothing executed, no damage done, no cleanup. The check lives inside the tool because that is the last line of defense that no prompt drift can bypass.
2. **Execute in a sandbox.** Dry-run first, isolated copy, never directly on the target.
3. **Analyze semantically.** On failure, the tool returns the structured, actionable payload from pillar 2: the error category, the remediation guidance, and the *context the model cannot guess* (the staging schema version, the interface definition, the violated ADR).

Note what the ticket refusal is **not**: it is not a punishment, and, matrix in hand, it is not compensatory. It addresses the *symmetric risk* of amplifiers: an ungoverned amplifier amplifies systemic errors at the same scale as successes. Least privilege stays relevant even against a perfect model, because it protects the organization, not the model.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="12" id="phase-12" >}}

## Governing the transition debt

One piece is missing, and it is the one platform teams forget: **the exit plan for the scaffolding.**

Some of what we just built *is* compensatory, knowingly. The exhaustive enumeration of frozen legacy files exists because today's models cannot reliably infer "deprecated, do not touch" from annotations. The raw-stack-trace-to-JSON translator exists because today's models digest structured errors better than raw ones. Both will become dead weight as models improve — that is fine, *if we can find them again*.

So the proposal is to make the durability axis a **first-class attribute of every rule**: each policy, booster, and translator is tagged `amplifier` or `compensatory`, and every compensatory artifact must carry a measurable **sunset condition** ("model honors `@deprecated` semantically on >95% of our benchmark"). The platform exposes a transition-debt report: the compensatory ratio, and the list of pending sunsets. That ratio must trend toward zero across model generations: it is the health indicator of the platform investment.

For the record, the alternatives considered and rejected: putting everything in the system prompt (fragile, unverifiable, compensatory-declarative); gating only in CI (the late gate of phase 3: feedback arrives too late for self-correction); using an LLM to validate the plan (non-deterministic; we want a linter, not a judge).

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="13" id="phase-13" >}}

## The amplified loop, materialized

Here is the full picture, with everything in place.

The intent enters light. The plan is enriched by ADR invariants, linted deterministically, and locked into a capability ticket. The tools verify the ticket at their door, execute in sandboxes, and answer failures with semantic guidance. The observation measures the gap against the original intent and feeds the loop. Every artifact along the way is tagged on the durability axis, and the scaffolding has a scheduled exit.

The result flows out of the loop with no gate at the exit: not because control was removed, but because it was *relocated* to where the work happens. The stream-aligned team and its agents ship faster **because** the platform constrains them well.

{{< /scrollytelling-step >}}

{{< /scrollytelling >}}

## The proof of concept

Concepts about agent governance are cheap; deterministic behavior is not. The companion repository [**poc-agentic-platform**](https://github.com/owulveryck/poc-agentic-platform) implements pillars 1 and 2 end-to-end in Go: small enough to read in an evening, complete enough to run the [whole cycle with `curl`](https://github.com/owulveryck/poc-agentic-platform/blob/main/docs/tutorial.md).

One question the code answers clearly: *how does the agent know what plan to submit to `lock_in_plan`?* Three orthogonal layers:

| Layer | Source | Controls |
|---|---|---|
| **When** to call `lock_in_plan` | `CLAUDE.md` | Behavioral rule |
| **How** to format the plan | MCP tool schema (auto-generated from [`plan.Plan`](https://pkg.go.dev/github.com/owulveryck/poc-agentic-platform@v0.0.1/internal/plan#Plan)) | JSON structure |
| **What** the plan must contain | `enrich()` invariants | Semantic content |

The schema validates structure deterministically; the linter validates ADR compliance deterministically; the model fills in the business content from the enriched context. None of the layers overlap. The full reasoning behind each design decision is in the repository's [explanation](https://github.com/owulveryck/poc-agentic-platform/blob/main/docs/explanation.md).

The planning side implements the sequence below: enrichment from a real ADR store (four ADRs with YAML front matter, including one deliberately tagged `compensatory` with its sunset condition), a deterministic plan linter, and JWT capability tickets:

![The amplified planning sequence: enrich, lint, ticket](/assets/amplified-agentic-loop/amplified-planning-sequence.en.svg)

The execution side implements the Smart Tool contract (in-tool ticket verification, sandboxed execution, semantic analysis):

![A Smart Platform Tool consumes the ticket](/assets/amplified-agentic-loop/smart-tool-sequence.en.svg)

Two details worth noticing in the code:

* **The generic raw→JSON translator and the semantic enrichers are separate packages**, because they sit in different cells of the matrix. The day models read raw stack traces natively, the first is deleted without touching the second.
* **The debt report intentionally ships in `DEBT_ALERT`** (two compensatory artifacts out of five): the point of the mechanism is to be seen.

### What it looks like in today's tools

None of this requires a custom agent. The repository ships adapters that wire the gateway into off-the-shelf tools: a few lines of configuration, zero fork.

**Claude Code, pillar 1: planning over MCP.** The gateway is exposed as a stdio [MCP server](https://modelcontextprotocol.io/), built with the official [Go SDK](https://github.com/modelcontextprotocol/go-sdk), with two tools. One command registers it:

```bash
claude mcp add ppg -- go run ./adapters/claudecode/mcpserver
```

and one line in the project's `CLAUDE.md` states the contract: *"submit your structured plan through `lock_in_plan` before modifying anything."* From there, Claude Code calls `get_platform_guidelines_for_intent` natively, receives the ADR invariants, and locks its plan; on success the capability ticket lands in a `.ppg-ticket` file at the project root.

**Claude Code, pillar 2: gating over hooks.** This is the part I find most satisfying. Claude Code's [`PreToolUse` hooks](https://code.claude.com/docs/en/hooks) provide exactly the in-tool interception point the Smart Tool contract requires. A `ppg-guard` binary is registered on every `Edit|Write`:

```json
{ "hooks": { "PreToolUse": [ { "matcher": "Edit|Write",
    "hooks": [ { "type": "command", "command": "ppg-guard", "args": [] } ] } ] } }
```

The hook verifies the target file against the ticket scope. In scope: silent pass. Out of scope: the hook exits with code 2, which **blocks the tool call before anything executes**; the message on stderr is fed back to the model:

```
OUT_OF_PLAN_SCOPE: "internal/auth/login.go" is not part of the locked plan
(allowed: migrations/001_seka.sql, internal/payment/router.go, ...).
Nothing was modified. If this change is genuinely needed, re-plan through
lock_in_plan.
```

Read that message as the agent does: it is a deterministic refusal *and* a remediation path. In a live session, Claude reads it and either stays inside the locked plan or goes back through `lock_in_plan`: the drift is corrected by the loop itself, exactly as phase 11 of the diagram describes. An off-the-shelf agent, governed in-loop, without modifying a single line of the agent. The [adapter README](https://github.com/owulveryck/poc-agentic-platform/blob/main/adapters/claudecode/README.md) has the full setup.

**GitHub Copilot: the black-box path.** Copilot exposes no loop to intercept, so the adapter falls back to pre-flight context generation: `go run ./adapters/preflight "add the Seka payment method"` calls `/enrich` and writes the invariants into [`.github/copilot-instructions.md`](https://docs.github.com/en/copilot/customizing-copilot/adding-repository-custom-instructions-for-github-copilot), which Copilot reads natively as repository custom instructions. The same governance reaches the black box, declaratively only. The honest caveat: without an interception point, hard gating cannot happen in-loop and must fall back to apply time (a pre-push platform check). That gap is the difference between an agent with an open loop and a closed one; it is a good criterion when choosing your tooling.

The documentation follows the [Divio](https://docs.divio.com/documentation-system/)/[Diátaxis](https://diataxis.fr/) documentation system (tutorial, how-to guides, reference, explanation), so the repository doubles as a template for documenting platform products. What the PoC would need to leave PoC status is listed in its [explanation](https://github.com/owulveryck/poc-agentic-platform/blob/main/docs/explanation.md): embedding-based ADR retrieval instead of keywords, asymmetric keys behind a KMS instead of a hard-coded secret, a real copy-on-write sandbox, and the whole observation pillar.

How to judge whether such a gateway works in your organization? The success criteria I would track:

| Objective | Observable metric |
|---|---|
| Guidance before action | ADR invariants appear in the agent's plan |
| Deterministic gating | A non-conforming plan is rejected 100% of the time |
| Actionable feedback | The agent corrects its plan in ≤ 2 iterations |
| Multi-tool portability | Works via [MCP](https://modelcontextprotocol.io/) **and** via pre-flight instruction files |
| Traceability | Every locked plan is logged (hash + scope) before execution |
| Debt governance | Compensatory ratio decreases at each model upgrade |

## Where this fits in the landscape

The table below maps the closest existing systems against the three phases. Columns mark genuine coverage, not aspirational claims.

| System | Planning-time context | Execution-time gating | Capability scopes | Semantic error returns |
|---|---|---|---|---|
| [NeMo Guardrails](https://docs.nvidia.com/nemo/guardrails/latest/) | Dialog rails (scripted) | Execution rails (pre/post tool) | No | No (blocking/redirect) |
| [POLARIS](https://arxiv.org/abs/2601.11816) (Jan 2026) | Typed DAG + rubric gate | Compiled policy guards | No | Implicit only |
| [Agentic JWT](https://arxiv.org/abs/2509.13597) (Sep 2025) | No | Pre-execution scope enforcement | Yes (cryptographic) | No |
| [Harnessing Embodied Agents](https://arxiv.org/abs/2604.07833) (Apr 2026) | Capability admission | Runtime detection | Proto-tickets | No |
| **This article** | ADR invariants injected live | Capability ticket verified in-tool | Yes (governance-intent JWT) | Yes (structured, actionable) |

Two things in the bottom row have no column to the left of them in any existing system. First: the integration. Each of the four prior systems covers one phase; the mechanism that makes governance *amplifying* rather than merely *present* is the coherence across phases: the planning context shapes the ticket, the ticket constrains the tool, the tool's semantic errors re-enter the plan. Second: the framing. Every existing system treats safety and productivity as a tradeoff to be managed. The claim here is different: a guardrail placed at the right phase of the loop *removes* friction rather than adding it, because the agent acts with more confidence, fewer late-gate surprises, and less rework. The [trustworthy agentic AI survey](https://arxiv.org/abs/2605.23989) frames the residual problem as a "trust-utility tradeoff"; the argument here is that the tradeoff is an artifact of placing the governance in the wrong place, not an inherent property of governance.

This remains a proof of concept. Keyword-based ADR retrieval would need embedding-based semantic search at scale. The JWT secret is symmetric and hard-coded. The sandbox is simulated. Pillar 3 (observation) is explicitly out of scope. None of that invalidates the architecture; it maps the production path. What the PoC demonstrates is that the three-phase integration is *buildable* with off-the-shelf components (Go, JWT, MCP) and *composable* with existing agents without forking them.

## Conclusion

The thesis of this series, compressed into one sentence: **the platform is the conductor of the agentic workforce, and its guardrails are the accelerator, not the brake.**

The front-loaded model (a wall of instructions, then hope, then a late gate) treats the agent as something to be contained. The amplified model treats it as something to be *equipped*: the map at planning time, the rules inside the tools, the gap measurement at observation time. Same governance, relocated — and the relocation changes everything: feedback that arrives inside the loop costs one iteration; feedback that arrives after the loop costs the whole loop.

For the platform team deciding where to invest, the heuristic from [See, Act, Correct](/2026/06/04/see-act-correct-three-levers-for-working-with-a-code-agent.html) remains the compass: *will this be more useful, or useless, when the model is twice as intelligent?* Fund the bottom-right cell of the matrix (plan linters of semantic invariants, capability tickets, semantic feedback), tag the rest as scaffolding, and give every piece of scaffolding a sunset date.

What remains open is the third pillar at scale: agentic telemetry, measuring across thousands of loop executions which guardrails amplify and which merely compensate, and feeding that back into the ADR store. That closing of the outer loop deserves its own article.

*One side note for the future: the intent does not have to come from a human. It could arrive from another agent — an enterprise architect agent delegating a sub-task, a platform orchestrator breaking down a larger plan. [That direction](/2026/06/25/from-isolated-agents-to-agentic-mesh-orchestrating-sdlc-with-a2a-and-ap2/) is where agent-to-agent protocols (A2A, AP2) come in; the governed loop described here becomes one node in a governed mesh.*

*Let's make AI work* ~~*on your machine*~~ *in your organization — and let the platform hold the rails.*
