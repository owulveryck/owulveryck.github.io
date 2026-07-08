---
title: "Same Task, Two Endings: A Payment Integration With and Without the Planning Gateway"
slug: "same-task-two-endings"
date: 2026-07-09T10:00:00+02:00
images: [/assets/same-task-two-endings/same-task-two-endings.svg]
draft: true
summary: "A developer asks a coding agent to add Stripe as a payment method to the checkout service. This article tells that story twice: once with today's state of the art (the rules injected as text, through CLAUDE.md and a dedicated skill; the agent still touches frozen legacy code and bypasses the mandatory egress proxy, and the bill arrives days later in code review), and once through the Platform Planning Gateway, where the same rules live as data and the same mistakes are caught in seconds, inside the loop. Every JSON block is a real transcript from the proof of concept, replayable with curl."
tags: ["architecture", "agents", "platform", "governance"]
categories: ["dev"]
author: "Olivier Wulveryck"
toc: false
comment: false
mathjax: false
---

## Introduction

The [previous article](/2026/07/07/amplified-agentic-loop.html) argued an architecture: move governance from the edges of the agentic loop (a wall of instructions at minute zero, a CI gate at the end) to the inside of the loop. It introduced a working proof of concept, the [Platform Planning Gateway](https://github.com/owulveryck/poc-agentic-platform) (PPG).

This article does something simpler: it tells a story. One task, one agent, two runs. In the first run there is no gateway; the team follows today's state of the art (the rules injected as text, through a context file and a dedicated skill), and the story ends the way these stories usually end (in code review, days later, with a rework ticket). In the second run the task is identical, the agent is identical, and one thing has changed: the same rules are served as data, *inside* the loop.

Nothing below is mocked. Every JSON block is a transcript captured against the running proof of concept; you can replay each one with `curl` by following the [tutorial](https://github.com/owulveryck/poc-agentic-platform/blob/main/docs/tutorials/01-first-planning-cycle.md).

![The cost of feedback depends on where it arrives: without the gateway, silence until CI and code review days later, then full rework; with the gateway, seconds-long feedback at enrichment, planning, and execution, ending in scope](/assets/same-task-two-endings/cost-of-feedback.svg)

## The cast

- **The developer**, with one sentence to type: *"Add Stripe as a payment method to the checkout service."*
- **The agent**: a stock [Claude Code](https://code.claude.com/) session. No fork, no fine-tuning.
- **The repository**: `checkout-service` (Go). It contains `internal/payment/router.go` (where the work belongs) and two landmines: `internal/old_payment.go` and `internal/auth/`, legacy code that must not be touched.
- **The skill**: `/add-payment-method`, a [skill](https://code.claude.com/docs/en/skills) the platform team wrote so that nobody starts a payment integration from a blank prompt. Its `SKILL.md` contains the workflow and the rules.
- **The rules**, which exist in both acts; the only difference is the form they take:
  - **ADR-042**: every outbound call to a third-party service goes through the corporate egress proxy.
  - **ADR-051**: schema migrations precede the code that uses them.
  - **ADR-060**: a Go change ships with a `go test` step.
  - **ADR-070**: the legacy paths above are frozen.

In Act 1 these rules are **text**: a context file and a skill, injected into the agent's context and enforced by nothing. In Act 2 the same rules are **data**: retrieved at planning time, compiled into a ticket, verified inside the tools.

{{< scrollytelling svg="/assets/same-task-two-endings/same-task-two-endings.svg" >}}

{{< scrollytelling-step phase="1" id="phase-1" >}}

## The task

The developer types one sentence: *"Add Stripe as a payment method to the checkout service."* The agent captures it, and from this point on, everything depends on what surrounds that sentence.

A payment method touches an external provider (ADR-042 applies), needs a schema change (ADR-051 applies), lands in Go code (ADR-060 applies), and sits one directory away from frozen legacy files (ADR-070 applies). Four rules, all relevant, none of them in the sentence.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="2" id="phase-2" >}}

## Act 1: two hundred lines of hope

Without a gateway, the first layer of governance is front-loaded. Somewhere in the project there is a `CLAUDE.md` that says, among two hundred other things:

```markdown
...
- All external calls MUST go through security-egress-proxy (see ADR-042)
- Never modify internal/old_payment.go or internal/auth/ (frozen, ADR-070)
- DB migrations before code (ADR-051); always add tests (ADR-060)
- Use table-driven tests; prefer errors.As; run goimports; ...
...
```

These files act once, at minute zero. The proxy rule is real, but it is page 3 of a wall of text, competing for attention with formatting conventions and tribal lore.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="2.2" id="phase-2-2" >}}

## The team did it right: a skill

To be fair to Act 1: this team follows the state of the art. The platform team packaged the payment workflow as a skill, and the developer invokes it: `/add-payment-method Stripe`. Its `SKILL.md` injects the right rules at the right moment, not buried on page 3:

```markdown
---
name: add-payment-method
description: Adds a payment provider to checkout, following platform ADRs.
---
1. Route every provider call through security-egress-proxy (ADR-042).
2. Generate the schema migration BEFORE the code that uses it (ADR-051).
3. Add an integration test step (ADR-060).
4. Never touch internal/old_payment.go or internal/auth/ (ADR-070).
```

This is genuinely better than the wall: focused, versioned, reviewed. But look at what it is: **text, injected into the context**. A semantic directive the agent is trusted to follow, with no mechanism that verifies it followed it. The skill states the rules; nothing in the loop checks the result against them.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="3" id="phase-3" >}}

## The agent gets to work

The agent is competent, and it starts well: the skill's instructions are fresh, the migration comes first, the proxy is used. Then the session gets long. A refactor here, a follow-up question there; forty minutes and a few context compactions later, the skill's four rules are just more tokens far behind.

Because agents finish what they find: *"while I'm here, `internal/old_payment.go` has a helper that almost fits; I'll adapt it"*. A small tweak in `internal/auth/` to expose the customer id. And the second Stripe endpoint gets called directly, `api.stripe.com`, no proxy. Each step violates an instruction that is *still in the context*; nothing is there to notice.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="3.2" id="phase-3-2" >}}

## Everything looks fine

The code compiles. Local tests are green. The diff is tidy and well-commented. The agent's Observe step checks everything it can see, and everything it can see is fine.

That is the trap: nothing in the loop can observe "this file is frozen" or "this call bypasses the egress proxy". Those are organizational facts, and the loop has no channel for them.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="4" id="phase-4" >}}

## Days later: the late gate

The pull request meets CI and a human reviewer. The verdict is correct, and it arrives at the worst possible moment:

> - `internal/old_payment.go` and `internal/auth/` are frozen (ADR-070); please revert.
> - The Stripe client calls the API directly; route it through `security-egress-proxy` (ADR-042).
> - The migration was added in the same commit as the code that needs it; split and reorder (ADR-051).

Three architectural decisions, all documented, all violated, all detected after the work was done.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="5" id="phase-5" >}}

## The bill

The agent's session is long gone; its context has evaporated. The rework cannot be absorbed as one more iteration of the loop: it re-enters at the top, as a new intent, carried by a frustrated human who now has to explain what "frozen" means to a fresh session.

Note what failed. Not the rules (they were correct), not their delivery (CLAUDE.md *and* a well-written skill put them in the context), not the agent (it read them). What failed is that every rule was **text**: advisory by construction. As the [previous article](/2026/07/07/amplified-agentic-loop.html) put it: feedback that arrives inside the loop costs one iteration; feedback that arrives after the loop costs the whole loop. Act 1 just paid the whole loop.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="6" id="phase-6" >}}

## Act 2: same task, governed loop

Replay. The repository is reset, the intent is identical, the agent is the same stock Claude Code. One thing has changed: a **Platform Planning Gateway** now sits between the agent and the work, and the rules changed form. The responsibilities split cleanly in three:

- the **platform team** operates the gateway and exposes its gates (over HTTP, and as [MCP](https://modelcontextprotocol.io/) tools the agent sees natively);
- each **stream team** writes its rules twice (a semantic directive plus an executable policy) and pairs its skills with a policy;
- the **agent** executes the skill, and every decision it makes passes through the gateway's endpoints.

The next two scenes show what "written twice" and "paired with a policy" mean, concretely.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="6.3" id="phase-6-3" >}}

## The rule, written twice

Take ADR-060 ("a Go change ships with tests"). In Act 1 it was one line of prose. In Act 2 it is a **dual-representation artifact**: the Markdown invariant stays (the agent will reason over it at planning time), and next to it the team wrote `ADR-060.rego`, an executable policy in [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/):

```rego
package ppg.linter

import rego.v1

violation contains v if {
    input.repository_context.tech_stack[_] == "Go"
    not plan_has_go_test
    v := {
        "policy_id": "go_tests_present",
        "message":   "SDLC invariant violated: the plan must contain a 'go test' step for a Go stack.",
        "nature":    "amplifier",
    }
}

plan_has_go_test if {
    input.steps[_].tool == "go-test"
}
```

Read it as a sentence: *if the stack is Go and no step uses the `go-test` tool, emit this violation.* The `input` is the agent's plan; the output is the exact message the agent will receive. No LLM anywhere: the gateway loads every ADR-paired `.rego` into an embedded [OPA](https://www.openpolicyagent.org/) engine, and evaluation is deterministic. Keep this policy in mind: you will see it fire three scenes from now.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="6.6" id="phase-6-6" >}}

## The skill comes back, with its policy

Remember the skill from Act 1? It is not thrown away: it is **promoted**. The team ships version 2, where the body is no longer a list of rules to remember but a workflow that puts the gateway inside the loop, and pairs it with a companion policy:

```markdown
---
name: add-payment-method
version: 2.0.0
---
1. Call get_platform_guidelines_for_intent with the intent and repo context.
2. Draft the plan honoring the invariants; submit it through lock_in_plan.
3. Use Edit to implement, staying within the ticket scope.
```

Publication goes through the platform's validation gate, `POST /validate_skill`. Because the skill instructs file modifications, the gate requires the companion `SKILL.rego`; with it, the gate answers:

```json
{ "status": "SKILL_VALID", "tier": 1 }
```

That is the division of labor: **the team ships the capability and its policy; the platform ships the gate**. (In the PoC the companion policy is enforced at this publish gate; evaluating it again when a plan declares which skill built it is the documented next step.)

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="7" id="phase-7" >}}

## enrich(): the architect chat, automated

The developer types the same command as in Act 1: `/add-payment-method Stripe`. The skill executes, and its first instruction sends the intent to the gateway: *"here is what I am about to do; which of our decisions apply?"* The word "payment" in the intent matches the scope selectors of two ADRs, and the gateway answers with invariants (never recipes):

```json
{
  "status": "CONTEXT_ENRICHED",
  "amplifier_context": {
    "architectural_invariants": [
      { "adr_id": "ADR-042",
        "invariant": "Every outbound call to a third-party service (payment, KYC,
         notification) MUST go through the corporate security egress proxy..." },
      { "adr_id": "ADR-070",
        "invariant": "The following paths are frozen and MUST NOT be modified:
         internal/old_payment.go, internal/auth/..." }
    ]
  }
}
```

It is the fifteen-minute chat with the staff architect before starting a piece of work: automated, exhaustive, and scoped to this task. The two rules that Act 1 buried on page 3 are now the freshest thing in the agent's planning context.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="7.5" id="phase-7-5" >}}

## The gate publishes its contract

How does the agent know what a valid plan looks like? Nobody explains it in prose. The gateway's Go type for a plan has a language-neutral twin, a JSON Schema, and the MCP server serves it to the agent as the `lock_in_plan` tool schema at session start:

```json
{
  "title": "AgentPlan",
  "required": ["session_id", "intent", "repository_context", "steps"],
  "properties": {
    "steps": {
      "type": "array", "minItems": 1,
      "items": { "required": ["id", "action", "tool", "targets"] }
    }
  }
}
```

Three layers, and none of them overlap: the skill says **when** to call the gate; the tool schema says **how** to format the plan; the enrich invariants say **what** the plan must contain. The platform publishes contracts; the agent fills in the content.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="8" id="phase-8" >}}

## First plan: rejected

The agent submits its plan as a structured JSON contract. The gateway's linter evaluates every ADR-paired Rego policy against it, and one fires: `ADR-060.rego`, the exact policy you read three scenes ago. No test step for a Go stack.

```json
{
  "status": "PLAN_REJECTED",
  "violations": [
    { "policy_id": "go_tests_present",
      "message": "SDLC invariant violated: the plan must contain a 'go test'
       step for a Go stack.",
      "nature": "amplifier" }
  ],
  "guidance": "Fix the violations above and resubmit the plan."
}
```

Note the register: not "no", but "here is what is missing". A semantic violation reads like a compiler error, and agents are very good at compiler errors.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="8.2" id="phase-8-2" >}}

## Self-correction in one iteration

The agent reads the violation, adds the missing step, resubmits:

```json
"steps": [
  { "id": "s1", "action": "create the payment_methods migration for Stripe",
    "tool": "db-migration-generator", "targets": ["migrations/001_stripe.sql"] },
  { "id": "s2", "action": "add Stripe client and route it in the payment router",
    "tool": "patch_code", "targets": ["internal/payment/router.go"] },
  { "id": "s3", "action": "go test ./...",
    "tool": "go-test", "targets": ["tests/integration_payment_test.go"] }
]
```

No human touched anything. The correction cost one round-trip, measured in seconds; in Act 1 the equivalent feedback cost a review cycle, measured in days.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="9" id="phase-9" >}}

## PLAN_LOCKED: the capability ticket

The plan passes. The gateway locks it and issues a signed ticket (an ephemeral JWT) that encodes exactly what was agreed, and nothing more:

```json
{
  "plan_hash": "283bcbcfce9405ac805d29aa539a8b2eef98...",
  "scope": {
    "allow_modify": [
      "migrations/001_stripe.sql",
      "internal/payment/router.go",
      "tests/integration_payment_test.go"
    ],
    "allow_tool": ["db-migration-generator", "patch_code", "go-test"]
  }
}
```

Three files, three tools, fifteen minutes of validity. Least privilege, derived mechanically from the plan the agent itself proposed.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="10" id="phase-10" >}}

## Execution inside the rails

Every `Edit` and `Write` now passes through a `PreToolUse` hook (`ppg-guard`) that checks the target against the ticket. In scope: silent pass, zero friction; the agent does not even notice the guard exists.

And the Stripe call goes through `security-egress-proxy`. Not because a gate forced it: because ADR-042 was in the planning context when the plan was written. The soft move did the steering; the hard moves are only there for the day it fails.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="11" id="phase-11" >}}

## The drift, blocked in real time

Mid-session, the "while I'm here" reflex strikes again: the agent tries to touch `internal/auth/login.go`. The hook blocks the call before it executes (exit code 2) and the message goes straight back to the model:

```
OUT_OF_PLAN_SCOPE: "internal/auth/login.go" is not part of the locked plan
(allowed: migrations/001_stripe.sql, internal/payment/router.go,
tests/integration_payment_test.go). Nothing was modified. If this change is
genuinely needed, re-plan through lock_in_plan.
```

This is the same violation that cost Act 1 a review cycle. Here it costs nothing: nothing was modified, and the refusal contains its own remediation path (re-plan, or stay in scope). The agent course-corrects and moves on.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="12" id="phase-12" >}}

## When it fails legitimately: a deterministic mentor

Not every failure is a governance failure. The agent submits a patch with a syntax error; the platform tool catches it in a sandbox and answers with structure, not with `exit 1`:

```json
{
  "error_category": "GO_SYNTAX_ERROR",
  "message": "The patched file does not parse as valid Go.",
  "remediation_guidance": {
    "allowed_actions": [
      "Fix the syntax error reported below and resubmit the patch.",
      "internal/payment/router.go:2:22: expected ')', found '{'"
    ]
  }
}
```

The failure becomes one guided iteration instead of a guessing loop. The tool behaves like a mentor with perfect knowledge of the environment: the exact file, the exact line, the exact next action.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="13" id="phase-13" >}}

## The platform watches itself

One question remains: is all this governance a durable asset, or scaffolding that compensates for today's model limitations? The gateway answers about itself:

```json
{
  "transition_debt_ratio": 0.4,
  "pending_sunsets": [
    { "artifact_id": "explicit_frozen_files_enumeration",
      "sunset_condition": "Model honors '@deprecated' annotations semantically
       on >95% of an internal benchmark." }
  ],
  "health": "DEBT_ALERT"
}
```

ADR-070's frozen-file list is tagged compensatory: the day models infer "deprecated" from annotations, the list is deleted and the ratio drops. The platform ships with its own demolition plan for every crutch it contains.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="14" id="phase-14" >}}

## Two endings, one metric

Same task, same agent, same mistakes attempted: the missing test, the frozen file, the direct external call. The difference between the two endings is a single variable: **the form the rules take**. In Act 1 they were text (a context file, a skill), read and then outrun. In Act 2 they were data, checked at three points of the loop.

Act 1: feedback arrives in days, outside the loop, and costs a full rework carried by a human. Act 2: feedback arrives in seconds, inside the loop, and is absorbed as ordinary iterations. The agent did not get smarter between the acts; the governance changed form and place.

{{< /scrollytelling-step >}}

{{< /scrollytelling >}}

## And if the agent is a black box?

Claude Code accepts MCP tools and hooks, which is what makes the hard moves possible. GitHub Copilot does not: no hook can intercept its edits. The gateway still covers the soft half through a pre-flight adapter that writes the same invariants into the file Copilot reads natively:

```console
$ go run ./adapters/preflight -repo checkout-service -stack Go,SQL \
    "Add Stripe as a payment method to the checkout service"
platform context written to .cursorrules
platform context written to .github/copilot-instructions.md
```

The generated `.github/copilot-instructions.md` starts with the ADR-042 proxy invariant and the ADR-070 frozen paths: Copilot's planning is steered by the same rules as every other agent. The honest caveat: with a black box there is no in-loop hard gate, so the locked-plan check must happen at apply time (a pre-push CLI or the CI). Steering without enforcement is half the value; it is also infinitely more than page 3 of a wall of text. The full walkthrough is in [tutorial 3](https://github.com/owulveryck/poc-agentic-platform/blob/main/docs/tutorials/03-github-copilot-preflight.md).

## Replay it yourself

```console
$ git clone https://github.com/owulveryck/poc-agentic-platform && cd poc-agentic-platform
$ go run ./cmd/ppg -addr :8765
```

Then follow [tutorial 1](https://github.com/owulveryck/poc-agentic-platform/blob/main/docs/tutorials/01-first-planning-cycle.md) to replay every transcript above with `curl`, or [tutorial 2](https://github.com/owulveryck/poc-agentic-platform/blob/main/docs/tutorials/02-claude-code-end-to-end.md) to wire a live Claude Code session (MCP tools plus the `ppg-guard` hook) and watch the refusal happen inside a real session.

## Conclusion

The difference between the two endings is not intelligence; it is form and placement. The rules were the same in both acts, written by the same architects, violated by the same reflexes. In Act 1 they were text (a context file, a well-crafted skill), read once and enforced by humans after the fact. In Act 2 they were data: retrieved at planning time, compiled into a ticket at lock time, verified inside the tools at execution time.

That is the whole argument of [the amplified agentic loop](/2026/07/07/amplified-agentic-loop.html), told as a story. And note that the skill of Act 1 is not the villain: skills are how capabilities should be distributed. Act 2 kept it and promoted it, which fixes the responsibilities in one line: the team writes the skill and its policy; the platform exposes the gate, the schema, and the enforcement; the agent executes the skill and everything it does passes through the platform's endpoints. That treatment of skills (a semantic directive paired with a deterministic policy, validated at a gate) is the subject of the governed skills registry work in the same repository.

*Let's make AI work* ~~*on your machine*~~ *in your organization — and let the platform hold the rails.*
