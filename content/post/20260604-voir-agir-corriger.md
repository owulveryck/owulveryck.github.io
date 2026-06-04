---
title: "See, Act, Correct: three levers for working with a code agent"
date: 2026-06-04T10:00:00+02:00
lastmod: 2026-06-04T10:00:00+02:00
images: [/assets/voir-agir-corriger/agent-leviers.en.svg]
draft: false
keywords: ["code agents", "Claude Code", "Copilot", "developer tools", "engineering practices"]
summary: "An out-of-the-box code agent only sees a repo and a shell. This article lays out three invariant principles for turning a gadget into a production tool."
tags: ["AI", "agents", "architecture", "engineering-practices"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

*An out-of-the-box code agent only sees a repo and a shell. For professional engineering, that is not enough. Here are the principles that make the difference between a gadget and a production tool.*

> **Foreword.** This article grew out of the talk ["Beyond the Basics with Claude Code"](https://www.youtube.com/watch?v=tuY2ChJIx48) by Daisy Holman, an engineer on the Claude Code team (May 2026). The founding ideas come from that talk, then supplemented by my own field experience and supported by documented research to cross-check them against available data.

**"If the agent can't do everything you do, it can't work with you."** — Daisy Holman, Anthropic

That sentence sums up the problem every team eventually encounters. A freshly installed code agent knows how to read code and run commands. That is enough for a prototype, a side project, a zero-to-one exercise. But real engineering work does not live in source code. It lives in Slack threads, design docs, production dashboards, review discussions, unwritten architecture decisions. The code tells *what* is done; it rarely tells *why*.

Try this: spend an entire day in your agent's terminal, without ever switching to another tool. Every alt-tab is a missing connection, a place where the agent cannot follow you. Customizing an agent is not a luxury. It is a structural necessity.

This article is aimed at senior engineers, tech leads, and architects who are deploying or considering deploying code agents in production environments. It lays out three invariant principles (valid regardless of the tool) and a method for applying them.

*Transparency note: this framework was built from extensive experience with Claude Code, then generalized. The mechanism examples (skills, hooks, MCP, prompt files) are borrowed from that ecosystem, but the principles apply to the equivalent primitives of other tools (rules files for Cursor, custom instructions for Copilot, repo maps for Aider). When a recommendation is tool-specific, it is flagged.*

---

## 1. The mental model: See, Act, Correct

Every code agent, regardless of the engine powering it, boils down to three questions:

- **SEE**: what does the agent know? What sources of information are accessible to it? Code, logs, internal documentation, conversation history, CI state.
- **ACT**: what can the agent do? Edit files, run tests, open PRs, query an API, deploy.
- **CORRECT**: what corrects it? Linters, test failures, review feedback, validation hooks, automatic feedback loops.

![The three levers of a code agent](/assets/voir-agir-corriger/agent-leviers.en.svg)

This framework aligns with several independent traditions: the **OODA** cycle (Observe, Orient, Decide, Act) in decision theory, the **Observation / Action / Reward** paradigm in reinforcement learning, feedback loops in systems engineering. This is no coincidence: these are structural constraints of any system that perceives an environment, acts on it, and adjusts.

This convergence is not perfect. Boyd's OODA cycle includes an **Orientation** phase (the step where raw data is interpreted and synthesized) that our framework subsumes under SEE. This choice is deliberate: in production engineering, Orientation is largely the product of context packaging (which information, in what order, in what form). But the cost must be acknowledged: a framework without explicit Orientation tends to treat the agent as a purely reactive system. For tasks requiring deep architectural reasoning, this dimension warrants separate treatment.

Similarly, this framework deliberately differs from the cognitive architectures proposed by academic research. Lilian Weng's widely cited architecture ("LLM Powered Autonomous Agents," 2023) divides the agent into Planning, Memory, and Tool Use [^weng]. Shunyu Yao's ReAct paradigm separates the reasoning trace from the external action. Our framework sacrifices that granularity in favor of applicability: in production, planning and memory are context packaging problems (how much history to inject, in what form, at what point). But if your agent requires self-reflection (the ability to revise its own decisions independently of external correction signals), this dimension merits separate treatment.

[^weng]: Lilian Weng, "LLM Powered Autonomous Agents," Lil'Log, 2023. https://lilianweng.github.io/posts/2023-06-23-agent/

The parallel between CORRECT and the reward signal in RL is not just a metaphor. Recent work on reinforcement learning training for code models explicitly uses composite reward functions combining functional correctness (unit tests), syntactic correctness (linters), and semantic structure (data-flow graphs) [^syncode]. Your linters and tests are not mere after-the-fact filters; they function as multidimensional reward functions that shape the agent's generation policy.

Everything else (instruction files, MCP servers, hooks, skills, prompt files) is merely the implementation of these three levers. The levers are the invariant; the mechanisms are contingent.

[^syncode]: "Domain-Adaptable Reinforcement Learning for Code Generation with Dense Rewards," arXiv:2605.21180, 2025. https://arxiv.org/abs/2605.21180

### The weakest-link principle

**An agent is never better than the weakest of its three levers.**

Imagine a bar chart. Each lever has a height. The agent's overall performance is capped by the shortest bar.

This principle is extensively documented. An internal Anthropic study of 132 engineers using Claude Code measured a 67% increase in the number of individual Pull Requests merged per day [^anthropic-study]. Yet the organization's overall delivery metrics (throughput, DORA metrics) did not change. This is the "productivity paradox": the increase in generation speed (ACT) led to a proportional increase in the review burden (code review time increased by 91%) because the validation loop (CORRECT) did not keep up. The weakest link shifted from the engineer writing the code to the engineer who had to audit a massive volume of asynchronously generated code. Observability platforms like Axify and Faros AI report the same pattern: without overhauling the test and review infrastructure, the rework rate increases [^axify] [^faros].

Customizing a single axis is pointless. Giving dozens of tools to an agent that does not understand your context is like handing it the keys to a factory without giving it the blueprint. Giving it full visibility without a correction loop is letting it produce without quality control. The balance among the three levers is what determines whether the agent is a production tool or a PR generator whose output needs constant rework.

[^anthropic-study]: "How AI Is Transforming Work at Anthropic," Anthropic Research, 2025. Study covering 53 interviews and 200,000 Claude Code session transcripts (February–August 2025). https://www.anthropic.com/research/how-ai-is-transforming-work-at-anthropic
[^axify]: "AI coding tools' impact: Metrics, ROI, and Review Signals in 2026," Axify. https://axify.io/blog/ai-coding-tools-impact
[^faros]: "Measuring Claude Code ROI," Faros AI. https://www.faros.ai/blog/how-to-measure-claude-code-roi-developer-productivity-insights-with-faros-ai

---

## 2. Context is a physical resource: package to scale

An agent's context window is a fixed, expensive resource. Most agent engines today operate around 200,000 effective tokens; some models advertise much larger windows (Gemini 2.5 claims up to 1 million), but the *usable* window for an agent in a working session remains constrained by cost and latency. Everything injected into it (tools, instructions, conversation, files read) **competes with actual work**.

As in any constrained system, every byte counts (you don't run `npm install` on an Arduino). A hidden cost compounds this: tokenizer evolution can affect the bill without any change in behavior. The token budget is not just about what you inject, but also about how the model consumes it.

This is not an artifact of any particular tool. It is a consequence of the transformer architecture (the attention mechanism and the KV cache) that underpins **all** current engines. Two rules follow from this.

### Rule 1: Don't pay for what you don't use

Every token injected into the context has a cost, both direct (billing) and indirect (it takes up space the task could have used). The classic anti-pattern: a team injects 50,000 tokens of architecture documentation on every task, including a simple variable rename. The context is saturated before the agent even starts working.

### Rule 2: Stable at the top, volatile at the bottom

To generate each token, the model recalculates the attention relationships among all preceding tokens. The **KV cache** (Key-Value cache) stores the results of these calculations to avoid redoing them on every turn. This is what keeps a long conversation fluid and affordable. But this cache works sequentially: modifying something early in the prompt invalidates everything that follows. The cost of this invalidation is publicly documented:

| Provider | Cost reduction (cache hit vs miss) | Details |
| :---- | :---- | :---- |
| Anthropic (Claude) | **90%** | Cache read at 0.1x the base price. 25% surcharge on the first write. TTL 5 min [^anthropic-cache] |
| DeepSeek | **~90%** | Cache hit at ~$0.003/M tokens vs ~$0.14/M tokens on miss [^deepseek-price] |
| Google (Gemini) | **~75%** | Configurable TTL, optimized for massive contexts [^llm-pricing] |

In agentic workflows where an agent makes dozens of round-trips for a single task, prefix caching reduces financial consumption by over 80% per session [^mindstudio-cache]. The stakes are not marginal. Hence the placement rule:

[^anthropic-cache]: "Prompt caching," Claude API documentation. https://platform.claude.com/docs/en/build-with-claude/prompt-caching
[^deepseek-price]: "DeepSeek Pricing Explained," Flowith Blog. https://flowith.io/blog/deepseek-pricing-explained-most-tokens-per-dollar/
[^llm-pricing]: "LLM API Pricing 2026," PE Collective. https://pecollective.com/blog/llm-api-pricing-comparison/
[^mindstudio-cache]: "What Is Prompt Caching in Claude Code?," MindStudio. https://www.mindstudio.ai/blog/prompt-caching-claude-code-token-savings

![Placement rule: stable at the top, volatile at the bottom](/assets/voir-agir-corriger/context-placement.en.svg)

To make this concrete, here is the difference between a well-structured context and a poorly structured one:

![Comparison of a well-structured vs poorly structured context](/assets/voir-agir-corriger/context-structure.en.svg)

The first diagram recalculates nearly the entire context on every turn. The second recalculates only the lower portion — the part that changes anyway.

**Caveat: this control does not exist everywhere.** In CLI agents (Aider, Claude Code), the engineer controls the ordering via explicit commands and root-level configuration files (CLAUDE.md). In agentic IDEs (Cursor, Copilot Workspace), the user defines intentions via rules files (.cursorrules, copilot-instructions.md), but the internal RAG engine decides which fragments are injected and in what order. In autonomous cloud agents (Devin), the packaging is an inaccessible black box [^tools-comparison]. Before investing in placement optimization, verify that your tool gives you that control.

[^tools-comparison]: "Every AI Coding Tool Compared," Developers Digest. https://www.developersdigest.tech/blog/ai-coding-tools-comparison-matrix-2026

### Packaging to scale

Scaling means treating context like embedded code. Structure it. Prioritize it. Load on demand.

For every resource you inject into the context, ask yourself: **"Will this approach hold if the project grows by a factor of 10?"** If today you bulk-load your 20 conventions, your 5 ADRs, and the descriptions of your 8 MCP tools, it works. But when the project has 200 conventions and 50 tools, that same approach will saturate the context before the agent starts working. The right abstraction loads what the agent needs, when it needs it — not everything, all the time.

The mechanisms vary by tool: instruction files scoped by directory, lazily loaded skills, tool discovery on demand rather than an exhaustive catalog. But the principle is the same everywhere: **package context like embedded code**, with the same optimization discipline.

One last point on knowledge: fine-tuning is rarely the right answer for injecting domain knowledge. Gekhman et al. (EMNLP 2024) show that LLMs acquire most of their factual knowledge during pre-training. Supervised fine-tuning primarily teaches the format and style of the interaction. When new factual knowledge does end up being assimilated through fine-tuning, it linearly increases the model's propensity to hallucinate [^gekhman].

Fine-tuning still has an edge for behavioral and stylistic skills: specific syntactic conventions, mastery of a DSL not represented in the training data, adaptation of code review tone. But for volatile factual knowledge (team conventions, architecture state, recent decisions), it is too expensive, too slow, and models evolve too fast to recoup the investment. The preferred path remains in-context learning — well-designed text files, injected at the right moment into the context.

[^gekhman]: Gekhman et al., "Does Fine-Tuning LLMs on New Knowledge Encourage Hallucinations?," EMNLP 2024. https://aclanthology.org/2024.emnlp-main.444/

---

## 3. The question that outlasts every tool

Code agents come with a variety of extension mechanisms: instruction files, prompt files, MCP tools, hooks, scripts, skills. It is tempting to evaluate them one by one. But there is a more lasting lens.

All these mechanisms fall on a single spectrum:

![The declarative / programmatic spectrum](/assets/voir-agir-corriger/audit-spectrum.en.svg)

On the left: you describe what you want. It is portable, simple to set up, but fundamentally fragile — you *hope* the model follows your instructions.

On the right: you execute code and verify the result. It is more powerful, but requires engineering effort and an engine that supports it.

### Compensatory vs amplifier

This declarative/programmatic spectrum intersects with a second, independent dimension: the **durability** of the investment in the face of model evolution.

- **Compensatory tools** work around the model's current limitations. They become **less** useful as the model improves.
- **Amplifier tools** leverage the model's capabilities. They become **more** useful as the model improves.

These two axes are independent. A tool can be declarative *and* an amplifier, or programmatic *and* compensatory:

|  | **Compensatory** (loses value as the model improves) | **Amplifier** (gains value as the model improves) |
| :---- | :---- | :---- |
| **Declarative** | A 200-line instruction listing forbidden files | Description of the system's architectural invariants; ADRs injected into context |
| **Programmatic** | Hook that forces a specific formatting the future model will understand natively | Pre-commit hook that runs tests and re-injects failures; security scanner that surfaces findings |

This distinction is formalized in human-computer interaction research under the terms "equalizer" vs "cognitive amplifier" effect. Empirical studies show that generative AI acts as a compensatory tool (equalizer) on routine tasks (leveling performance between novices and experts) but as a pure cognitive amplifier on complex tasks, where output quality depends on the human's ability to define the problem and guide iterative refinement [^equalizer]. The concept of "scaffolding," ubiquitous in frameworks like LangChain, is inherently compensatory: it is a temporary procedural structure designed to work around the model's current limitations. Recent history confirms this: the query-chaining frameworks built in 2023 to work around GPT-3.5's 4K-token windows have become useless crutches now that we have 200K+ token windows.

The question to ask for every extension you build:

**"Will it be more useful or useless when the model is twice as smart?"**

If the answer is "less useful," you are investing in transitional debt. If it is "more useful," you are building a durable advantage. This question holds regardless of the position on the declarative/programmatic axis.

The takeaway is straightforward: **prioritize amplifier tools**, whether declarative or programmatic. And keep in mind the symmetric risk: an ungoverned amplifier is an amplifier of systemic errors. If the agent generates code incorporating obsolete architectures or vulnerable libraries, and the engineer uses it without critical judgment, the tool cements and amplifies the worst practices across the organization. This is why the recommendation to invest in amplifiers is inseparable from the imperative of control (section 4).

[^equalizer]: "AI as Cognitive Amplifier: Rethinking Human Judgment in the Age of Generative AI," arXiv:2512.10961. https://arxiv.org/abs/2512.10961

---

## 4. Guardrails: enabling autonomy within the rules

Agent vendors invest heavily in the SEE lever (larger context windows, RAG, integrations) and ACT (more tools, more autonomy). The CORRECT lever remains neglected. The reasons are structural: robust verification is specific to each organization (your rules, your conventions, your regulatory constraints), and it creates tension with the commercial narrative of autonomy. Look at the product pages of major agents: the number of integrations and the size of the context window are highlighted; verification and control mechanisms are rarely mentioned.

This is your blind spot. And it is your responsibility.

### Autonomy is not a goal — it is a means

The goal is not for the agent to do whatever it wants. The goal is for it to work **within the framework you have defined**, without constant supervision. The distinction is fundamental.

Consider the contrast: the code for your VS Code extension goes through a full CI pipeline (lint, typecheck, tests, integration, security), a minimum test coverage, and weekly analytics. The code generated by your agent? Often no post-generation verification, no programmatic feedback loop, no observability (tokens, cost, quality), no audit trail.

Setting up guardrails does not hamstring the agent. It is **what enables autonomy**. An agent without rails is a risk you watch constantly. An agent with hooks that block forbidden actions, linters that correct in real time, tests that validate before every commit, and granular permissions — that agent you can let work alone.

Guardrails do not slow the agent down. They let it run faster on a marked track.

### Six dimensions to enforce

No vendor handles these questions properly, because they are commercially uncomfortable:

| Dimension | Question to ask |
| :---- | :---- |
| **Autonomy** | Does the agent propose or execute? With what reversibility? |
| **Accountability** | Who answers when the agent is wrong? |
| **Verification** | Where does the bottleneck move when production breaks? |
| **Cost / value** | Does the compute consumed justify the actual gain? |
| **Security** | What attack surface does each granted access open? |
| **Observability** | How many tokens per task? What success rate for correction loops? What progressive drift in code quality? |

The guiding principle:

**"Every additional access granted to the agent demands a proportional investment in verification."**

Access without verification is not a productivity gain. It is a security debt. Example: a team enables autonomous mode and gives the agent read access to Slack, without post-action auditing. The agent, exposed to a message containing a prompt injection, modifies a deployment file. Nobody notices. Now consider the same team, with validation hooks on critical files, tests before every commit, permissions scoped by directory. The agent works alone, fast, and within the rules.

This risk is documented: studies show success rates of 84% for indirect prompt injections on agentic code editors [^aishell], and OWASP (the reference organization for application security) classifies these vectors in its Top 10 LLM 2025 [^owasp].

[^aishell]: "Your AI, My Shell: Demystifying Prompt Injection Attacks on Agentic AI Coding Editors," arXiv:2509.22040. https://arxiv.org/abs/2509.22040
[^owasp]: OWASP Top 10 for LLM Applications v2025. https://owasp.org/www-project-top-10-for-large-language-model-applications/

### Guardrails are not infallible: defense in depth

Let's be honest: a validation hook will not catch a sophisticated prompt injection that produces a modification that is syntactically correct but semantically malicious. Programmatic guardrails are a layer of defense, not an absolute firewall.

This is why the recommended approach is **defense in depth** — multiple independent layers, each catching what the others miss:

- **Least privilege**: the agent accesses only the resources strictly necessary for its task. No broad network access, no global credentials.
- **Isolation**: the agent works in a sandbox (container, dedicated worktree) whose side effects are limited by design.
- **Immutable audit log**: every action the agent takes is logged in a non-alterable way. Even if a malicious action gets through, it is traceable and reversible.
- **Human review on critical paths**: modifications to deployment, infrastructure, or security files trigger an explicit review, regardless of the level of autonomy elsewhere.

No single layer is sufficient. It is their combination that makes the agent's autonomy reasonably safe.

---

## 5. When agents multiply: the multi-agent model

Everything above applies to a single agent in a single session. But the state of the art pushes toward **multi-agent** architectures: an orchestrator that distributes work among specialized sub-agents, each in its own context (a worktree, a container, a branch).

The See / Act / Correct model remains valid, but it becomes **distributed** — and that is where new questions arise.

### The three levers are distributed among agents

An orchestrator and its sub-agents do not share the same levers at the same level:

- **SEE** becomes segmented: each sub-agent sees only its scope (a module, a file, a task). The orchestrator must maintain an overall view, but its context is equally finite. Packaging becomes critical: what to transmit to each sub-agent, and in what form?
- **ACT** becomes parallelized: multiple agents modify code simultaneously. Isolation (worktrees, separate branches) prevents conflicts, but **reconciliation** (merging, semantic conflict resolution) becomes a problem in its own right.
- **CORRECT** becomes chained: a sub-agent's output is the orchestrator's input signal. The sub-agent's CORRECT lever (its tests, its linters) feeds the orchestrator's SEE lever. If this loop is poorly designed (for example, if the sub-agent returns a text summary instead of a structured, verifiable result), the orchestrator loses its ability to correct.

### The complexity trap

*Context Degradation* is the main risk of multi-agent architectures. Research from JetBrains presented at the NeurIPS 2025 Workshop on Deep Learning for Code ("The Complexity Trap") demonstrates that the iterative accumulation of execution traces, terminal errors, and tool outputs saturates the context window, causing a drop in effectiveness and a spike in costs [^complexity-trap].

The intuitive approach (using an LLM sub-agent to summarize semantic history) is often a complexity trap in itself. The study shows that a cruder strategy — "observation masking" (truncation and omission of old raw tool outputs) — reduces costs by over 50% while maintaining or even improving the success rate. As one analyst puts it: "Your instruction file is not a knowledge base — it is an attention budget" [^agents-md].

### The weakest link shifts

In a multi-agent architecture, the bottleneck is no longer necessarily in one lever of a single agent. It shifts to the **interfaces between agents**: the quality of the prompt the orchestrator sends (the sub-agent's SEE), the structure of the result the sub-agent returns (the orchestrator's SEE), and above all the question of **who verifies the final work** when no individual agent has seen the whole picture.

The weakest-link principle remains the same, but the unit of analysis changes: it is no longer one lever of one agent — it is the **weakest link in the chain of agents**.

On the security front, *weakest-link dynamics* apply: a single compromised MCP server providing poisoned tool definitions can contaminate the entire chain [^multi-agent-sec]. The boundaries between agents must be designed as circuit breakers, not as mere divisions of labor.

[^complexity-trap]: "The Complexity Trap," JetBrains Research, NeurIPS 2025 Workshop on Deep Learning for Code. https://arxiv.org/abs/2508.21433
[^agents-md]: "Your AGENTS.md is a Liability," paddo.dev. https://paddo.dev/blog/your-agents-md-is-a-liability/
[^multi-agent-sec]: "Open Challenges in Multi-Agent Security," arXiv:2505.02077. https://arxiv.org/abs/2505.02077

---

## 6. Four steps to evaluate any agent

The preceding principles boil down to a method applicable to any engine: Claude Code, Copilot, Cursor, Cody, Aider, or whatever comes tomorrow.

**Step 1 — Map.** Identify the agent's primitives and place them on the three levers. Which mechanisms feed SEE? Which extend ACT? Which drive CORRECT? If a lever has no mechanism, you have found your bottleneck.

**Step 2 — Identify who controls packaging.** Does the engineer control what enters the context, or does the vendor decide? This answer determines whether context economics is an actionable lever for you, or a black box you endure. The level of control varies dramatically:

| Level of control | Tools | Characteristics |
| :---- | :---- | :---- |
| **Full transparency** (CLI) | Aider, Claude Code | Surgical control: /add, /drop commands, CLAUDE.md file. The engineer decides every injected fragment |
| **Mediated declarative control** (IDE) | Cursor, Copilot Workspace | Rules files (.cursorrules, copilot-instructions.md), but the internal RAG engine orchestrates the final packaging |
| **Black box** (Cloud) | Devin, Factory | System prompt, underlying model, and context inclusion rules are proprietary and inaccessible |

**Step 3 — Classify on the spectrum.** For each extension you plan to build, ask the compensatory / amplifier question. Prioritize amplifier investments — they are the only ones whose value grows over time.

**Step 4 — Enforce the control layer.** Actively design the autonomy boundaries, cost tracking, security perimeters, and verification pipelines. No vendor will do it for you.

### The strategic trade-off: Control vs Integration

Behind these steps lies a central trade-off. Some agents give you total **control**: you decide on context packaging, you build your feedback loops, you choose your tools. Others give you **integration**: the ecosystem is preconfigured, the connections exist, but you accept the vendor's choices.

The right choice depends on two variables: the specificity of your context and the strength of your integration constraints.

|  | **Weak integration constraints** | **Strong integration constraints** |
| :---- | :---- | :---- |
| **Strong proprietary context** (specific conventions, compliance, internal tooling) | **Control**: you need to package a context that no one else understands. Example: a fintech team with compliance rules specific to its regulator. | **Hybrid**: you need both. Native integration for standard workflows, control for domain context. Highest cost, but sometimes unavoidable. |
| **Weak proprietary context** (standard stack, classic conventions) | **Lightweight**: a vanilla agent with a few instruction files is enough. Don't over-invest. Example: an early-stage startup on a standard React/Node stack. | **Integration**: your context is not special, but your tools need to communicate. Example: a full-GitHub team that wants the agent to handle PRs, issues, and CI without configuration. |

The most frequent mistake: choosing full control out of engineering instinct, when the proprietary context does not justify it. Control has an opportunity cost: every feedback loop you build is a loop you maintain.

The four-step method, however, is independent of the choice. It works in both cases.

---

## Conclusion

Three durable principles summarize this framework:

1. **An agent is never better than its weakest lever**: what it sees, what it can do, what corrects it.
2. **Invest in amplifier tools**, not compensatory ones. They are the only ones whose value grows as models improve.
3. **The control layer is your responsibility.** Guardrails do not hamstring the agent; they are the condition for its autonomy.

Models will keep improving. Some of today's context engineering will become unnecessary — transitional debt that the model will absorb. But the three-lever framework, the compensatory/amplifier test, and the control imperative will remain valid. These are not implementation details. They are structural heuristics robust enough to have survived several generations of tools. They have invalidation conditions: agents with persistent memory and online learning, multi-modal agents operating on visual interfaces, or agents capable of self-reflection without external correction signals could require an extended model. For the current state of the art, the three levers hold.

**Customizing a code agent is not about giving it everything. It is about deciding, under cost constraints, what it sees, what it can do, and what corrects it.**

---

## Appendix: the most frequent anti-patterns

| Anti-pattern | Failing lever | Symptom | Fix |
| :---- | :---- | :---- | :---- |
| 50K tokens of docs injected on every task | SEE (saturation) | Context exhausted before work begins | Lazy loading, scoping by directory |
| 15 MCP tools, no linter | CORRECT (absent) | Plausible code, breaks in CI | Pre-commit hook + automated tests |
| Autonomous mode + unaudited Slack access | CORRECT (absent) | Exposure to indirect injection | Scoped permissions, audit log, review on critical paths |
| Full control out of an engineer's reflex | Over-investment | Every loop built = a loop maintained | Assess proprietary context before investing (section 6 matrix) |
| Text summary between sub-agents | Orchestrator's SEE (degraded) | Loss of verifiability, drift | Structured, verifiable output (exit code, JSON, diff) |
| No per-task cost measurement | Observability (absent) | Unquantified risk disguised as productivity | Track tokens/task, rework rate, architectural drift |
