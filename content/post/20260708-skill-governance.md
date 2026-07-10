---
title: "The Governed Skills Registry: Policy-as-Code for Enterprise Agent Capabilities"
slug: "governed-skills-registry"
date: 2026-07-08T10:00:00+02:00
images: [/assets/skill-governance/skill-registry-architecture.svg]
draft: false
summary: "Skills are the distribution unit of the agentic workforce, but today's package managers are capability-blind: they version and distribute without validating governance. This article extends the Platform Planning Gateway with a skill governance linter: a deterministic OPA/Rego gate that classifies every skill by security tier before it reaches the enterprise registry."
tags: ["architecture", "agents", "platform", "governance", "skills"]
categories: ["dev"]
author: "Olivier Wulveryck"
toc: true
comment: false
mathjax: false
---

## Introduction

The [previous article](/2026/07/07/amplified-agentic-loop.html) closed the architecture of the Platform Planning Gateway: a proof of concept that governs *how* agents execute. Plans are linted against Rego policies derived from Architecture Decision Records, capability tickets constrain every tool call, and the smart tools verify the ticket before acting. The execution plane is governed.

There is a second plane that is not: the **capability plane**. Before an agent can plan or execute, it needs to know *what it can do*. In Claude Code, that answer arrives through **skills**: slash-command workflows that encapsulate reusable agentic patterns and make them invocable by any agent in the organization with a single `/skill-name` invocation.

Skills are the distribution unit of the agentic workforce. An enterprise that governs plans but leaves capabilities ungoverned has solved the wrong half of the problem: a well-governed plan built on a poorly-governed skill is still a risk.

This article asks: what makes a skill enterprise-valid? What policies must accompany it? How does a distribution mechanism (the skills equivalent of a package manager) enforce those policies? And how does the existing Platform Planning Gateway fit into the picture?

The companion repository [poc-agentic-platform](https://github.com/owulveryck/poc-agentic-platform) has been extended with a new `POST /validate_skill` endpoint and a `skill-governance/` Rego policy directory. The patterns are the same ones established for plan governance; the domain shifts from execution to capability.

---

## What is a skill?

A Claude Code skill is a `SKILL.md` file. Its front matter is deliberately minimal:

```yaml
---
name: patch-payment
description: "Applies targeted changes to the payment service, following platform ADRs for proxy and migration ordering."
version: 1.0.0
argument-hint: "[description of the change]"
---

# Skill body

Analyse the intent passed as $ARGUMENTS against the repository context.
Call `get_platform_guidelines_for_intent` to retrieve the relevant ADR invariants.
Produce a structured plan and submit it through `lock_in_plan` before using Edit.
```

Skills are discovered automatically from `~/.claude/skills/<name>/` (user-global) and `.claude/skills/<name>/` (project-local), and can be bundled inside plugins. They are not code: a skill is a **semantic directive**: natural language instructions that the agent interprets at invocation time. The model fills the business logic; the skill encodes the organizational invariants and workflow choreography.

The body is not documentation. It is the instruction set the agent executes literally when the skill is invoked. When a developer types `/patch-payment "add idempotency key"`, Claude reads the skill body and follows it step by step: it calls `get_platform_guidelines_for_intent` to retrieve the ADR invariants, it submits a structured plan through `lock_in_plan`, and it uses `Edit` to apply the changes. **The tools named in the skill body are the tools the agent will actually call at runtime.** Those tool calls are what the smart tools intercept and scope-check against the JWT capability ticket.

That semantic nature is precisely why governance matters. A skill that instructs the agent to modify files without going through `lock_in_plan` bypasses the execution governance entirely. As a reminder, `lock_in_plan` is the deterministic gate at the center of the [previous article's Platform Planning Gateway](/2026/07/07/amplified-agentic-loop.html): before touching any file, the agent submits its plan as a structured JSON contract; OPA/Rego policies reject non-conforming plans with semantic violations, and conforming plans receive a signed capability ticket that scoped tools verify in-tool. A skill distributed without a version cannot be audited or rolled back. A skill with a 10-character description is invisible in `/help`. These are not implementation bugs: they are governance gaps.

---

## Why classic package managers are not enough

A skills APM (Application Package Manager) does the right things: versioning, discovery, installation, dependency resolution. These are hard problems and the existing ecosystem solves them well.

What an APM cannot see:

- **Security tier**: does this skill instruct file modifications? Shell commands? Without a policy gate, a tier-2 skill (Bash, potentially destructive) can be published and installed with no more friction than a tier-0 read-only skill.
- **ADR alignment**: does the skill contradict an active architectural decision? A skill that instructs agents to call external services directly violates ADR-042 (all external calls go through the security proxy), but the APM has no knowledge of ADRs.
- **Companion policy**: does the skill declare the plan governance requirements it imposes? A skill that modifies files should carry a companion Rego file (`SKILL.rego`) that the plan linter can load at `lock_in_plan` time to enforce those requirements automatically.
- **Discoverability quality**: is the description long enough to allow the model to decide when to invoke the skill? A 10-character description will never trigger accurately.

These gaps are not APM failures: they are governance responsibilities the APM was never designed to carry. The solution is a **policy layer on top of the APM**: a skill governance linter that validates every skill before publication, and optionally at install time.

> **Aside: Microsoft APM — a concrete implementation**
>
> [Microsoft APM](https://microsoft.github.io/apm/) ("like npm for agent context") is the closest existing project to what this article calls "a skills APM." It handles exactly the distribution concerns described above: a single `apm.yml` manifest installs skills, MCP servers, and agent definitions across Claude Code, GitHub Copilot, Cursor, and several other harnesses. A lockfile pins content hashes for byte-for-byte reproducibility. An `apm-policy.yml` enforces enterprise → org → repo tighten-only inheritance at install time, and `apm audit --ci` runs as a mandatory CI gate on pull requests.
>
> That is a strong foundation. Where it stops is precisely where the Skill Linter starts. APM's `apm-policy.yml` is a *package-level* policy: it governs which packages are allowed or blocked, which MCP servers can be declared. It does not inspect the content of a skill's body. It cannot classify a skill by security tier (see below), require a companion Rego file for file-modifying skills, or check whether the skill's instructions contradict an active ADR. Those are semantic questions about what the skill *instructs the agent to do*: a dimension the distribution layer was never designed to reach.
>
> The two layers are complementary, not competing. APM is the right choice for distribution, integrity, and package-level policy. The PPG `POST /validate_skill` endpoint is the semantic gate APM can call (in CI before publication, or as an install hook) to enforce the governance policies that live above the package boundary. Together, they cover the full chain: APM ensures the right packages arrive intact; the Skill Linter ensures they deserve to.
>
> One more nuance: APM's git-source model (GitHub, GitLab, Azure DevOps, Bitbucket) maps to the *repository* layer described in the next section. The authoritative install source that APM points consumers at is the *registry*. The `apm-policy.yml` governance operates at the registry-install boundary, which is exactly where the Skill Linter gate belongs.

---

## Two levels of distribution: repository and registry

Before describing the governance architecture, it is worth establishing two concepts that the software ecosystem has kept distinct for decades but that the skills world often collapses into one.

A **skill repository** is any Git source: a GitHub repository, a GitLab project, an internal Bitbucket server. It is where skills are authored and versioned. Multiple repositories can coexist: a platform team's internal repo, an open-source community repo, a vendor-published repo. A repository carries no inherent authority: it is a source, not a trust boundary.

A **skill registry** is the enterprise's authoritative catalog: the single source that agents and APMs install *from*. A skill enters the registry only after passing the policy gate. The registry maintains provenance (which repository a skill came from, which version, which policy approved it) and a **dependency graph**: a skill can declare dependencies on other skills or on specific MCP servers, and the registry validates that all transitive dependencies are themselves registered and approved before the skill can be installed.

The npm analogy is exact: a GitHub repository is where you write and commit code; npmjs.com is the registry you publish to and that others install from. The `npm publish` step (and the policy checks that can gate it) is the governance boundary. The same structure applies here: `git push` lands the skill in a repository; the CI-triggered Skill Linter validates it; a passing validation promotes it to the registry.

This separation has two important consequences:

**Decentralization is a feature at the repository level.** Any team, any open-source project, any vendor can author and version skills in their own repository. The ecosystem grows without central coordination.

**The registry is the trust boundary.** Only skills that have passed the Skill Linter enter the registry. Consumers (developers, agents, CI pipelines) install from the registry, not directly from repositories. The dependency graph in the registry ensures that a skill cannot pull in an unapproved transitive dependency: every node in the graph was independently validated.

The governance gate (`POST /validate_skill`) sits precisely at the boundary between repository and registry. This is the point at which the policy questions the APM cannot answer (security tier, companion Rego requirement, ADR alignment, structural quality) are answered deterministically, before any installation can occur.

![Repository vs. registry: decentralized git sources feed a policy gate; only validated skills reach the authoritative registry](/assets/skill-governance/skill-concept.svg)

---

## Governance rules for enterprise skills

### Structural rules

These rules are always enforced, regardless of what the skill does:

| Field | Constraint | Nature |
|---|---|---|
| `name` | lowercase-kebab-case, ≤ 32 characters | amplifier |
| `description` | 50–500 characters, starts with a verb | amplifier |
| `version` | semver, required for registry publication | amplifier |
| `argument-hint` | required when body uses `$ARGUMENTS` | amplifier |
| Body length | ≤ 500 lines, no hardcoded secrets | amplifier |
| Companion Rego | required for tier ≥ 1 skills | amplifier |

All structural rules are tagged `amplifier`: they enforce durable SDLC invariants that become *more* valuable as the model improves, not less.

### Security tiers

The tier is derived from the tools the skill instructs the agent to use:

| Tier | Trigger | Publication gate |
|---|---|---|
| **0 — Read-only** | Body mentions only Read, Glob, Grep | Auto-approved by Rego |
| **1 — Standard** | Body mentions Edit or Write → file modifications | CI + Rego validation |
| **2 — Privileged** | Body mentions Bash → shell commands | Human review required |

The tier is not a prediction; it is a direct consequence of what the body instructs. Because the agent executes the body literally, a body that says `Edit` will produce file modifications and a body that says `Bash` will execute shell commands. This is why the linter classifies on tool mentions: deliberately binary, keyword-based, blind to intent. A skill that *could* be safe but mentions `Bash` is tier 2 regardless, for the same reason the plan linter is not an LLM: determinism and reproducibility take priority over nuance.

### The dual-representation skill

The insight from the previous article applies directly here. An ADR in the Platform Planning Gateway is a **dual-representation governance artifact**: its Markdown body (the semantic directive) is injected at `enrich()` time into the agent's planning context; its companion `.rego` file is evaluated deterministically at `lock_in_plan` time by the plan linter.

A skill with programmatic enforcement follows the same pattern:

| Representation | File | Consumed by | Moment |
|---|---|---|---|
| Semantic directive | `SKILL.md` body | Agent at invocation | Planning |
| Governance policy | `SKILL.rego` companion | Skill Linter at publish + Plan Linter at lock_in_plan | Validation |

![Dual representation: SKILL.md is read literally by the agent at invocation; SKILL.rego is evaluated by the Skill Linter at publish and loaded by the Plan Linter at lock_in_plan](/assets/skill-governance/skill-dual-representation.svg)

The skill's `SKILL.rego` serves two distinct purposes:
1. **At publish time**: the Skill Linter validates the skill's own structure and security tier.
2. **At runtime**: the Plan Linter can load the skill's companion Rego as an additional policy bundle when the plan originates from that skill's invocation, giving the skill author direct control over what the plan must contain.

This is the durability axis applied to capabilities: the semantic directive is a permanent amplifier (the skill instructions get more useful as the model improves); the Rego companion may be compensatory scaffolding (required today because the model does not yet reliably honor the invariant without enforcement, but removable once it does). Note that the two durability tags live at different levels and do not conflict: the structural rule *requiring* a companion Rego for tier ≥ 1 skills is an amplifier (the obligation to declare enforceable requirements ages well), while the *content* of a particular companion Rego may be compensatory and carry its own sunset condition.

---

## The architecture

Three gates validate a skill before and during use. Each addresses a distinct failure mode.

![The three gates on one lifecycle: Gate 1 validates at publish (structure, tier, companion Rego), Gate 2 at install (content hash, revalidation), Gate 3 at runtime (Plan Linter with the skill's companion Rego)](/assets/skill-governance/skill-three-gates.svg)

### Gate 1 — Publish: the Skill Linter

![Architecture overview: Skill Author → CI Pipeline → Skill Linter (OPA) → Policy Store → Registry; and consumer → Registry → PPG with Plan Linter and Smart Tools](/assets/skill-governance/skill-registry-architecture.svg)

The Skill Linter is a new endpoint on the Platform Planning Gateway: `POST /validate_skill`. A skill author pushes `SKILL.md + SKILL.rego` to a skill repository (any git host); CI triggers the linter. The endpoint receives the skill's metadata and body (and optionally its companion Rego content), evaluates them against the governance policies in `skill-governance/`, and returns either `SKILL_VALID` with the computed security tier or `SKILL_REJECTED` with structured violations. Only a passing skill is promoted to the authoritative registry.

Every violation carries a `nature` field: the same `amplifier` / `compensatory` tagging used for plan policies. This makes the violation actionable: an amplifier violation must be fixed because the rule is a durable invariant; a compensatory violation signals scaffolding that the author can accept now and schedule for removal later.

The publish sequence:

![Skill authoring and publication sequence: Author → CI → Skill Linter → Policy Store → Registry, with SKILL_REJECTED and SKILL_VALID branches](/assets/skill-governance/skill-authoring-sequence.svg)

### Gate 2 — Install: revalidation

When a consumer installs a skill from the registry (using `apm install` or equivalent), the PPG re-runs the validation. Two distinct threats live at this boundary, and they call for two distinct defenses. Tampering (a skill modified in the registry after publication) is a job for content hashes: APM's lockfile already pins each package byte-for-byte, and the registry records the approved hash at publish time; revalidating the content is redundant with verifying the hash. What revalidation uniquely catches is *policy drift*: the governance policies may have been tightened since publication, and a skill that passed last quarter's rules may fail today's. Install-time revalidation is therefore the last structural check before the skill enters a project's `.claude/skills/`.

### Gate 3 — Runtime: the Plan Linter with skill companion

The third gate is the existing `POST /lock_in_plan` mechanism, extended. When the agent submits a plan that originates from a skill invocation, it can include the `skill_id`. The Plan Linter loads the base ADR policies (the existing `adr/*.rego` files) *and* the skill's companion `SKILL.rego` from the policy store, evaluating the union of all violation rules against the plan.

This is the runtime sequence:

![Runtime sequence: Agent installs skill, invokes it, calls enrich, produces plan, submits to lock_in_plan with skill companion Rego loaded, receives capability ticket, uses Smart Tool](/assets/skill-governance/skill-runtime-sequence.svg)

The result: a skill-invoked plan is held to both the organizational ADR invariants *and* the skill-specific plan requirements. A `/patch-payment` skill can require, in its companion Rego, that any resulting plan include a `go test` step, on top of the ADR-060 requirement that already enforces this for all Go plans. The governance rules compose without coupling.

One scoping note for the honest reader: the PoC implements Gate 1 end to end (`POST /validate_skill`, replayable against the repository); install-time revalidation (Gate 2) and the `skill_id` wiring at `lock_in_plan` (Gate 3) are the documented production path, tracked in the repository's `AUDIT.md`.

---

## The POC extension: `/validate_skill`

The [poc-agentic-platform](https://github.com/owulveryck/poc-agentic-platform) now exposes:

```
POST /validate_skill
```

Request body (JSON, mirrors the `Skill` struct):

```json
{
  "name": "patch-payment",
  "description": "Applies targeted changes to the payment service, following platform ADRs for proxy and migration ordering.",
  "version": "1.0.0",
  "argument_hint": "[description of the change]",
  "body": "Analyse the intent... Use Edit to patch...",
  "rego_policy": "package ppg.skills.patch_payment\nimport rego.v1\n..."
}
```

Response on failure:

```json
{
  "status": "SKILL_REJECTED",
  "violations": [
    {
      "field": "rego_policy",
      "message": "Skills that instruct file modifications (tier ≥ 1) must include a companion SKILL.rego declaring their plan governance requirements.",
      "nature": "amplifier"
    }
  ],
  "guidance": "Fix the violations above before publishing the skill to the registry."
}
```

Response on success:

```json
{
  "status": "SKILL_VALID",
  "tier": 1
}
```

The `internal/skill` package follows exactly the same pattern as `internal/linter`: a `NewLinter(governancePolicyDir string)` constructor loads all `.rego` files from `skill-governance/` into a single OPA `PreparedEvalQuery` over `data.ppg.skills.governance.violation`; `Validate(*Skill)` evaluates the query and unmarshals violations; `Tier(*Skill)` computes the security tier in Go (keyword-based, deterministic).

The governance policies live in two Rego files that share `package ppg.skills.governance`:

- `skill-governance/structure.rego` — structural rules: name format, description length, version, argument-hint
- `skill-governance/security.rego` — security rules: companion Rego required for tier ≥ 1

All rules are `violation contains v if {...}`: the same OPA pattern as the plan linter. Adding a new governance rule is adding one Rego rule to an existing file or a new file in the directory; no Go code changes required.

---

## The debt model applied to skills

The transition-debt report (`GET /debt_report`) currently tracks plan linter policies and smart tool translators. The same model applies to skills:

- **Amplifier skills** encode durable organizational workflows that become more useful as models improve: a skill that choreographes the three-step plan → lock → execute cycle will remain valuable indefinitely.
- **Compensatory skills** compensate for current model limitations: a skill that spells out in exhaustive detail how to find deprecated files is scaffolding; it will be useless once the model infers this from `@deprecated` annotations reliably.

Every compensatory skill should carry a `sunset_condition` in its companion Rego, making the retirement debt explicit and measurable. The governance rules do not enforce this today; it is the next natural extension of the structure rules once the organization has enough skills in production to make the debt visible.

---

## Pattern generality

Nothing in this architecture is specific to Claude Code. The governance rules, the Rego policy files, and the three-gate validation model apply to any skill runtime that:

- Distributes skills as structured packages (name, description, version)
- Exposes a plan submission step where the planning gateway can be consulted
- Provides in-tool verification of a capability scope (the JWT ticket pattern)

The Claude Code specifics (the `SKILL.md` front matter, the `/skill-name` invocation, the `PreToolUse` hook) are the adapter layer. The governance model is the invariant.

---

## Honest limitations

This remains a proof of concept, and the same honesty that closed the previous article applies here. Keyword-based tier classification can be evaded by paraphrase: a body that says "run the migration from the terminal" without naming `Bash` classifies as tier 0 today. The production posture is deny-by-default: any tool mention outside the known read-only allowlist (including unknown MCP tools) escalates the tier, and a body with no recognizable tool mention is treated as unclassifiable, not as safe. There is also a version-skew window at Gate 3: the agent executes the locally installed `SKILL.md` while the Plan Linter loads the companion `SKILL.rego` from the registry. The two must be pinned together (the registry serves the Rego of the *installed* version, identified by the content hash recorded at publication), or a stale skill would run under mismatched plan rules. Neither gap invalidates the architecture; both map the path from proof of concept to production.

---

## Conclusion

The [Platform Planning Gateway](https://github.com/owulveryck/poc-agentic-platform) governs the execution plane. Skills governance closes the capability plane. Together they form a coherent answer to the question that opened the series in [See, Act, Correct](/2026/06/04/see-act-correct-three-levers-for-working-with-a-code-agent.html): how do we move from individual vibe coding to industrialized, at-scale software engineering?

The execution plane says: *every action the agent takes was planned, the plan was linted, and the scope is enforced in-tool.*

The capability plane says: *every workflow the agent can invoke was authored according to documented rules, validated by a deterministic policy engine, and published through a gated distribution channel.*

Neither plane alone is sufficient. An ungoverned capability can produce a governed plan and still cause harm: the governance happened too late. A governed capability without execution governance produces a correctly packaged skill that an agent uses to take ungoverned actions: the governance stopped too early.

The guardrail is the loop. Skills governance and plan governance are two iterations of the same pattern applied to two consecutive moments in the agentic cycle. The next iteration (observation, the third pillar explicitly deferred from the previous article) closes the loop: structured, semantic feedback from execution back into planning, making each cycle smarter than the last.

Let's make AI work in your organization.
