---
title: "The Agentic Threat"
date: 2026-10-08T10:00:00+02:00
lastmod: 2026-10-08T10:00:00+02:00
images: [/assets/menace-agentique/where-guardrails-act.en.svg]
draft: false
keywords: ["agentic engineering", "security", "guardrails", "OPA", "Rego", "Persona Guardrail"]
summary: >
  The more autonomy we give agents, the more convenient it is, and "convenient" is one of the most dangerous words in the digital world. Systems must be protected from the human threat, but also from the agent itself. These two threats do not call for the same architectural answer: Uber's Persona Guardrail on one side, a deterministic gateway (PPG) on the other, and the distinction between compensating and amplifying measures.
tags: ["AI", "agents", "architecture", "security", "agentic-engineering"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

## Autonomy is convenient

I am more than convinced that **agentic engineering** is the discipline that will unlock the full value of AI in the enterprise.
I am lucky enough to work on the evolution of the processes used to build digital assets, software and others.
In that context, one key is to provide agentic systems that act with **as much autonomy as possible** to carry out non-differentiating tasks robustly and quickly (and ideally at a controlled cost).

**Development harnesses** have opened the way to agent autonomy. It is now possible to have a conversation to create digital assets with Claude Code, Copilot or others.
You switch on an auto-mode and, after some framing, off it goes.
The harness comes with a set of tools that lets it interact with the ecosystem. In a company, interacting with the ecosystem means retrieving information, acting on tools or on processes.

The more tools we give, the more autonomously agents can work. And as I was saying to a colleague this morning: the more autonomy we give, the more convenient it is…. And **"it's convenient" is one of the most dangerous phrases** in the digital ecosystem…. Because systems that offer convenience usually do so in exchange for something else of value to themselves.
In the case of agentic systems, we **trade task delegation for control**: we hand autonomy over to the agent, and we widen the **attack surface** accordingly.

And that autonomy can be dangerous.

## Two risks, two threats

We know that integrating software carries **two risks**:
- **Data corruption** (the agent breaks everything), a database or a codebase for example.
- **Data disclosure**: the agent retrieves data and may be programmed to transfer it to third parties, or simply hand it to its malicious pilot.

And these risks are carried by **two threats**:
- The **human threat**: the agent is instructed to do things it is not allowed to do, either to sabotage systems (corruption) or for espionage (disclosure). And the pilot is not necessarily a stranger on the Internet: it can be a **malicious or manipulated employee**.
- The **threat of the agent itself**: the agent is intrinsically stupid and could corrupt or disclose data "by accident". I will set aside for now the threat of the **"double" agent**, which would have a secret, learned intention leading it to deliberately extract or corrupt data (to hide a backdoor, for example). I will just note that, seen from the outside, **a stupid agent and a double agent do the same thing**: what we control for one partly protects us from the other.

Addressing this requires agentic engineering, but in terms of architecture, **these two threats do not have the same answer**.

## The human threat: Uber's Persona Guardrail

For the human threat, one answer has been provided by Uber in its paper *Persona Guardrail: A Production-Grade Defense Framework for Agentic Systems* (ref [arxiv 2610.03434](https://arxiv.org/abs/2610.03434)). The paper is long and I admit I had an LLM ingest it to grasp its essence. The idea is to put a first LLM between the pilot and the agent to check that the request falls within the agent's **declared scope** (a list of allowed intents), and thus get a green light to start execution. The same LLM checks the **final answer** before it goes back to the pilot. However, **what happens in between (tool calls, changes) is not controlled**, and the paper says so itself.

## Protecting the agent from itself: PPG and jevgo

On the other hand, I think we also need **more deterministic validations** at the level of the **agentic loop** to protect the agent from itself. This is what I am exploring with **PPG** ([poc-agentic-platform](https://github.com/owulveryck/poc-agentic-platform)): a gateway that validates, with OPA/Rego rules, the agent's plan, each of its changes and the final diff. **No ticket, no change.**
I have also considered setting up a **"system 1"** (in Kahneman's sense: fast, statistical, intuitive) next to the deterministic loop, with **classification models** like [Jev](https://en.wikipedia.org/wiki/Jev_(AI_model)) (and my toy implementation: [jevgo](https://github.com/owulveryck/jevgo)): small models that learn to imitate the rules. **They decide nothing.** When their verdict diverges from the rules', it is a sign that a **rule may be badly written**, and we escalate when in doubt.

## Compensating and amplifying measures

To read these measures, I make a distinction:
- a **compensating** measure compensates for a weakness: it blocks, and each block requires a human to step in;
- an **amplifying** measure lets the agentic loop correct itself: the error goes back to the agent, which fixes it before a human has to step in. The human moves from **"in-the-loop"** (validating every step) to **"on-the-loop"** (supervising and handling exceptions).

## Three infographics

Here is a summary in three infographics.

**1. Where each guardrail acts.** Uber filters what is asked and what is answered; PPG controls what the agent does in between.

![Where each guardrail acts in the life of a request](/assets/menace-agentique/where-guardrails-act.en.svg)

**2. Same principle, opposite mechanisms.** Both declare a scope and deny by default. But you cannot write an exact rule for natural language, whereas you can for a plan or a diff.

![Same principle, opposite mechanisms: the shape of the input decides](/assets/menace-agentique/same-principle-opposite-mechanisms.en.svg)

**3. System 1 observes, it does not judge.** jevgo runs next to the rules. A disagreement goes up to a human, who fixes the rule once and for all.

![Adding jevgo to PPG: an observer of the rules, not a judge](/assets/menace-agentique/jevgo-observer.en.svg)

## In conclusion

Uber's architectural measure is **essentially compensating** (its improvement loop makes the policy progress, not the agent). And it **does not compensate for the agent's stupidity but for its obedience**: a smarter agent will obey a malicious request better, so the measure will not become useless as models improve. On the contrary, it may struggle to keep up. The **classifier is deliberately small** to stay under 100 ms. The better agents understand innuendo, the more a subtle request can be understood by the agent without being caught by the classifier, and the more **false negatives** the gateway will let through.
On the PPG side, there is obviously a compensating aspect, since I stated that the primary goal was to address the stupidity of models. But there is also an **amplifying aspect** for the agentic loop, which lets it self-correct when the gateway returns an error, before a human notices. The human no longer validates every step; they stay **on-the-loop**.

In any case, a real agentic architecture of the future must take these two threats into account, by **combining measures that protect and measures that make autonomy safe**.
