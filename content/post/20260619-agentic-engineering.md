---
title: "Vibe Coding at Scale? Engineering Strikes Back"
date: 2026-06-19T10:00:00+02:00
lastmod: 2026-06-19T10:00:00+02:00
images: [/assets/agentic-engineering/octo_plateforme_agentique.en.svg]
draft: false
keywords: []
summary: >
  Coding "by the book" is not enough. In a large organization shipping dozens of applications, the real challenge is not building one good app — it is ensuring every app meets the organization's own standards. That cannot be decreed project by project. It must be industrialized through an agentic platform.
tags: ["AI", "agents", "architecture", "platform", "vibe-coding", "agentic-engineering"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

Generative AI has transformed how code is produced. In just a few months, we went from autocomplete to agents capable of writing, testing, and deploying entire applications. The market is now flooded with methods for framing these agents and making them produce quality code.

But this abundance raises a question few organizations are asking yet: **what happens when you are not building one app, but fifty?**

<!--more-->

The answer lies in a distinction everyone underestimates. A piece of software can be *well-crafted* without meeting *the expected standard*. Structured AI methods can achieve **a** state of the art — the generic one the industry agrees on. But an organization's state of the art is *its own*: contextual, agreed upon, living. And that is exactly the one AI does not know — and reinvents, poorly, on every project.

## "State of the art" ≠ "YOUR state of the art"

Let us start by clearing up a fundamental misconception: **there is no absolute standard of quality**.

As Christophe Thibault reminds us in ["Done with 'Technical Debt'"](https://blog.octo.com/en-finir-avec-la-dette-technique), "Quality" as a singular, universal concept does not exist. At the Plaza Athénée restaurant, you get an excellent meal for around €380. At McDonald's, where lunch costs under €8, there is also a quality department. Two radically different standards — and each perfectly legitimate in its context. Swap their criteria, and you get two absurdities nobody would pay for.

A state of the art is the set of practices a group **agrees upon** — explicitly or tacitly — to produce the best possible solution in *its* context, with *its* constraints. It is neither universal nor fixed. It is debated, tested, and renewed.

But structured AI methods carry their own state of the art: the one distilled from all the generic know-how accumulated across the internet, libraries, and frameworks. That state of the art is valuable. But it is not yours. And therein lies the challenge: an application can be flawless by generic standards and yet completely off the mark by your organization's rules.

This is why we distinguish three approaches to AI-assisted development — not by the tools they use, but by **the state of the art they guarantee**.

**Vibe Coding** is fast. You prompt, iterate, and accept whatever "looks like it works." It is the ideal approach for a prototype, a hack, an exploration. But quality depends entirely on the prompt and the developer. No guaranteed state of the art, no reproducibility.

**Structured AI assistance** (of which [BMAD](https://github.com/bmadcode/BMAD-METHOD) is a good example) goes further. It enforces templates, rules, and detailed prompts. The result is a well-built application, coded "by the book." This is what we call **contextual certainty**: for this project, in this context, the result is reliable.

**Agentic engineering** introduces **systemic certainty**. Reliability is no longer carried by an individual or a prompt — it is carried by a **platform** that makes the organization's state of the art available to every agent, every project, every time.

![Contextual Certainty vs. Systemic Certainty](/assets/agentic-engineering/octo_spectre_certitude.en.svg)
*The differentiator is not whether you use AI: it is which state of the art your approach guarantees. Inspired by Figure 3 from [^osmani2026].*

## What structured AI actually does — and what it does not

Let us be fair: structured AI is a genuine step forward. It is not the enemy. It brings a rigor that vibe coding cannot offer, and produces well-made applications.

But it has two structural blind spots.

**First blind spot: it imposes its own build chain.** Its templates, conventions, and steps are designed *outside* the organization. They ignore existing processes — the very ones that, in a large organization, guarantee the overall quality clients expect. Those processes often need updating, certainly. But they cannot be bypassed.

**Second blind spot: the certainty it produces is local.** It says nothing about the next project. The next developer, the next agent, the next project starts from zero. Best practices may be followed, but not necessarily the organization's rules. The code is well-made, but not necessarily consistent with the other applications in the same system.

This is where a dynamic Jerry Weinberg described long ago comes into play:

> *The First Law of Technology Transfer: Long-range good tends to be sacrificed to short-range good.*
> — Jerry Weinberg, *Quality Software Management*

Structured AI optimizes the short-term good (*this* application, delivered fast and clean) at the expense of the long-term good: system-wide consistency. For one project, this is invisible. For fifty, it is devastating.

## The real danger at scale: reinventing the commodity

Here is the precise mechanics of chaos.

When each agent builds its application without knowing the organization's state of the art, it **invents its own**. And more importantly, it **rebuilds components that already exist elsewhere** — without knowing it. Authentication, error handling, observability, data access, UI components: all commodity building blocks, reconstructed project after project.

The result is not N excellent applications. It is **N mediocre components**. Instead of improving a shared component and continuously raising its quality, you scatter mediocre and divergent versions of the same need. **Nothing gets productized.**

Christophe Thibault proposes a far more accurate term than "technical debt" for this phenomenon: **process conflict**. A solution is not "in debt"; it relies on processes that, having never been agreed upon together, contradict each other. Before AI, this conflict was *rare and felt*: a developer would hit friction, mention it in standup, and the conflict would surface and get resolved.

At the scale of generative AI, this conflict changes nature and becomes formidable for two reasons:

- **It becomes invisible.** Each application "looks fine" because it follows the *generic* state of the art. Nobody feels friction. The agent does not stumble on anything. The conflict with the *organization's* state of the art is never noticed — until it manifests in production, in an incident, in an inconsistency visible to the customer.
- **It multiplies.** The sheer volume of applications that AI can produce turns a sporadic conflict into systemic divergence. Fifty projects means potentially fifty incompatible standards.

And the bill always comes due. Because trust is transitive. The end user trusts a **brand**. That brand deploys **applications**. Those applications are built by **agents**. If a single link breaks, trust in the entire brand erodes. A premium brand cannot afford amateur UX. A bank cannot tolerate a security flaw in a chatbot. A telecom operator cannot let the wrong tone slip through in a customer service app.

![Transitive Trust](/assets/agentic-engineering/octo_confiance_transitive.en.svg)
*Every deployed application is a test for the brand's reputation. Trust is built slowly and lost in an instant.*

## The platform: making the state of the art executable

If the problem is that agents do not consume the organization's state of the art but reinvent it, the solution is straightforward:

> **An agentic platform is the organization's state of the art made executable and consumable by AI — and governed as a product.**

The good news is that the raw material already exists. Security policies, architecture standards, the design system, domain knowledge, conventions: it is all there, within the organization. The problem is not that it is missing — it is that it is not **structured for agents to consume**. It lives in wikis, in people's heads, in PDFs, in slide decks. Unusable by an agentic system.

The platform's role is to turn this organizational capital into consumable capabilities: to organize it, tool it, and equip it with validation guardrails. These capabilities draw from **multiple content producers** — domain experts, security teams, brand teams, architects — and make them available to **multiple consumers** — agents, projects, applications.

This is exactly the definition of a platform as Sangeet Paul Choudary describes it in ["Platform Thinking: The Future of Work"](https://medium.com/@sanguit/platform-thinking-the-future-of-work-b49aeb0c1e53): a system that connects an ecosystem of producers and consumers, and **commoditizes repeatable operations**. Choudary makes a remark that is essential to our argument: *writing code is not a repeatable operation — it is a one-off infrastructure activity, like building an assembly line.* What is repeatable are the operations that code automates. The platform industrializes precisely what is repeatable, freeing the ecosystem to focus on what is not.

In practice, the platform delivers three families of systemic capabilities:

- **Systemic context** (instructions, domain knowledge, memory, examples) is structured once and injected into every agent.
- **Systemic guardrails** (security, reliability, brand consistency, conventions) are enforced automatically, not at the developer's discretion.
- **Tooling** (MCP servers, CI/CD pipelines, evaluations) is shared and ready to use.

The specifics — business intents, product-specific guardrails — layer on top of the systemic foundation. But the foundation is there, for every project.

![Where does context come from?](/assets/agentic-engineering/octo_usine_vs_plateforme.en.svg)
*The factory is necessary. But without a platform, each project reinvents context, guardrails, and tooling. Inspired by Figure 6 from [^osmani2026].*

The developer's role is transformed in the process. They no longer produce code: they **package context**. But that context does not come from them alone: it comes from domain experts, security teams, brand teams, architects. Each produces a fragment; the developer assembles them. And if everyone does this packaging their own way, project by project, you are right back where you started. This is precisely what the platform industrializes: a framework where every source of context is structured, versioned, and automatically injected.

## Surprise with value, not with standards

Once the commodity is handled by the platform, a guiding principle emerges.

You need to **shift left** on everything that is commodity: push quality, security, and standards as early as possible, upstream of the project, into the platform, so that they are a **given** rather than a deliverable to rebuild. The standard becomes a guaranteed starting point, not a finish line to reclaim every time.

The benefit is twofold. First, because the commodity is shared and productized, it **improves continuously**: every use hardens the component, instead of scattering N mediocre versions. Second — and this is what matters most — AI's energy is focused where it creates value: on the business, on the differentiator.

Because that is the whole point:

> **An application should not surprise the customer with its standard components. It should surprise them with its business value.**

Standards are expected. They must be impeccable, and invisible. The surprise, the satisfaction, the competitive advantage — those come from what nobody else does. That is exactly what a platform enables: commoditize the mundane so the exceptional becomes possible at scale.

![The Agentic Platform](/assets/agentic-engineering/octo_plateforme_agentique.en.svg)
*Shared AI-ready capabilities for building reliable applications at scale. The specific merges with the systemic to produce the application.*

## A platform is not a dogma: it is a governed product

At this point, a serious objection arises — and it is precisely by answering it that the thesis holds.

What counts as commodity is not a stable given. What differentiates today becomes tomorrow's standard. The boundary between the common and the differentiating **shifts constantly**. A platform that froze the state of the art would risk freezing... an outdated one — exactly the trap Thibault calls out under the label of "absolute Quality."

The answer is not to give up on the platform. It is to **treat it as a governed product**. Engineering leadership — CTOs, architects, tech leads — continuously evolves what the platform treats as commodity and what it leaves to individual applications. Choudary reminds us: opening a platform to an ecosystem means accepting a loss of control, and therefore building the *checks and balances* that maintain the standard.

This is how the state of the art stays **living and agreed upon**, in Thibault's sense, rather than carved in stone. The platform is not the authority that decrees absolute quality. It is the instrument that makes executable a state of the art that humans continue to debate, experiment with, and renew.

## A conviction

Vibe coding has its place. For exploring, prototyping, learning. It is not the enemy — it is the starting point.

Structured AI assistance — BMAD and similar approaches — is a genuine step forward. It brings rigor to individual production. But let us be clear about what it does not do: **in a large organization running dozens of applications, choosing a good structured AI method is not enough.** Each project remains an island, and nothing guarantees the consistency — or the quality at the organization's standard — of the whole.

The challenge is organizational. It is about ensuring that *every* application built by agents reflects the organization's state of the art — security, reliability, brand consistency, value — and that scaling up does not mean losing control. This requires three shifts:

- **Stop believing** that a method for framing agents is sufficient at scale. "Well-made" code is not the same as "on-standard" code.
- **Treat your agentic platform as a governed product**, not a tool. This is the condition for keeping the state of the art alive.
- **Inventory and structure your own state of the art** — security, architecture, design system, domain — to make it executable and consumable by agents.

The developer no longer produces code. They package the context that produces code. And to industrialize that packaging — to make the organization's state of the art executable, at scale, without freezing it — they need a platform.

![Context Engineering](/assets/agentic-engineering/fig4_context_engineering.en.svg)
*The best systems treat context as a first-class architectural decision, reviewed and versioned like code. Adapted from Figure 4 of [^osmani2026].*

---

## Sources

[^osmani2026]: Addy Osmani, Shubham Saboo, Sokratis Kartakis — *"The New SDLC With Vibe Coding: From ad-hoc prompting to Agentic Engineering"*, Google, May 2026.

Christophe Thibault — [*"Done with 'Technical Debt'"*](https://blog.octo.com/en-finir-avec-la-dette-technique), OCTO.

Sangeet Paul Choudary — [*"Platform Thinking: The Future of Work"*](https://medium.com/@sanguit/platform-thinking-the-future-of-work-b49aeb0c1e53), Medium, 2013.

Jerry Weinberg — *Quality Software Management*.
