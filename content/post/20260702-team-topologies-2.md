---
title: "Codifying the Rules: Building the Platform Behind the Agentic SDLC"
slug: "sdlc-team-topologies"
date: 2026-07-02T10:00:00+02:00
images: [/assets/sdlc-team-topologies/sdlc-team-topologies.svg]
draft: false
summary: "This article explores how organizations can scale reliable, AI-driven software development by combining modern Platform Engineering with the Team Topologies framework. It introduces a new Software Delivery Lifecycle (SDLC) where stream-aligned teams focus entirely on product solutions, while AI agents handle the technical implementation. To ensure enterprise-grade reliability and trust, human enabling teams temporarily intervene to establish guardrails and governance. Once these rules are codified directly into the platform, the enabling team gracefully vanishes, creating a continuous, self-sustaining loop of innovation and automation."
tags: ["architecture", "A2A", "AP2", "agents"]
categories: ["dev"]
author: "Olivier Wulveryck"
toc: false
comment: false
mathjax: false
---

## Introduction

This article is a follow-up to our deep dive on [Team Topologies](/2026/06/24/who-does-what-team-topologies-for-the-agentic-platform.html). Building reliable software at an organizational scale—across multiple products and teams—requires a fundamental shift in how we operate.

In a modern, AI-driven Software Delivery Lifecycle (SDLC), **the stream-aligned teams** should focus entirely on the *solution*, delegating the actual *implementation* to an AI system. This agentic loop is the engine of modern software, powered directly by a robust internal platform that provides the models and the inference engine.

But how do we kickstart this engine without it derailing?

Technology alone isn't enough; we need the right human collaboration. To build applications that are functional, reliable, and trustworthy, product experts must initially team up with technical experts. This **enabling team** steps in to structure the AI, establish crucial guardrails, and define the enterprise standards.

Ultimately, this article is a story of loops. The enabling team gathers these technical guardrails and bakes them directly into the platform as automated governance rules. Once the platform absorbs this knowledge, the enabling team can effectively vanish—leaving behind an agentic system capable of producing standard, trustworthy applications out-of-the-box. They will only reappear when a brand-new technical challenge emerges.

Let’s explore how to build and master these cycles, starting with a quick reminder of how the development engine works: the agentic loop.

{{< scrollytelling svg="/assets/sdlc-team-topologies/sdlc-team-topologies.svg" >}}

{{< scrollytelling-step phase="1" id="phase-1" >}}

## The agentic loop

This part is not an innovation of any kind. It is a reminder of how the agentic loop works.
Nevertheless, it is important to understand this foundation as it is the central point of software delivery. For the rest of this article, I will heavily rely on these principles.

The agentic loop is how modern software delivery powered by AI works.
You start by expressing an intention for something you want to develop.

The LLM understands the intention, and then plans some actions. During the planning step, the system evaluates what tools it needs to call to fulfill the need.
If a simple task is "write hello in test.txt", then the plan will produce something like:

* Call the writer tool, have it open `test.txt`, and write "hello" in it.

Then the agent actually calls the tools and executes the action.

Once the action is executed, the agent analyzes the return of the tool's execution to determine whether there has been an error or if the action was successful.

{{< /scrollytelling-step >}}
{{< scrollytelling-step phase="2" id="phase-2" >}}

### Self-Correction

One of the strengths of any agentic system is its ability to recover from errors. If a tool fails to reach its goal, the agent feeds the error output back into the planning phase, analyzes what went wrong, and generates a new plan.

For instance, imagine that in our previous example, the system cannot edit the file due to permission issues. The tool call returns an error like "`permission denied`". The agent adapts by prepending a step to change the permissions using the `chmod` tool:

1. call tool `chmod u+w test.txt`
2. Call the writer tool, have it open `test.txt`, and write "hello" in it.

*(While this is a simplistic example — a production-grade system would use more robust error handling rather than a blind `chmod` — it illustrates the mechanism.)*

Once these new tool calls are executed, the loop returns to the observation phase. If the execution was successful, the AI evaluates the current state: is the final intention fulfilled, or was this just an intermediate step that provides new context for the next action?

*Note: If the error is a transient error (such as a [`429`](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes) error in an API call), the agent may decide to retry the tool call without replanning.*

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="3" id="phase-3" >}}

### The true agentic loop

An agent goes beyond simply executing basic natural language commands; it can tackle complex, ambiguous tasks that require exploration. This is where the loop mechanism truly shines.

Imagine we give the agent a broad intent: *"Translate all the files in this directory into English."*

Because the agent lacks context, it cannot translate anything right away. Its first plan must be purely investigative.
The first iteration looks like this:

1. List the files in the directory.
2. Open each file to detect its source language.

Once these tools are executed, the agent enters the **observation phase**. It analyzes the results (e.g., discovering two files, `coucou.txt` and `salut.txt`, both written in French) and uses this new context to automatically **refine the intent** into something actionable:

> *Translate the files coucou.txt and salut.txt from French to English, and save them as hi.txt and hello.txt.*

Now equipped with a precise goal, the agent triggers a **second iteration of the loop**. It generates a new plan and calls the tools required to perform the actual translation.

Finally, the loop returns to the observation phase one last time, confirming that the two new files were successfully created with the correct English content. The task is complete.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="4" id="phase-4" >}}

### Exiting the loop
A loop successfully terminates when the observation phase confirms that the original intent has been fulfilled. However, to prevent the agent from getting stuck in an infinite cycle, programmatic guardrails—such as a strict maximum of 10 iterations—will force an exit if triggered.

Once the loop concludes, the agent delivers the final output (we will dive deeper into this handoff later in the article).

While developers typically run this agentic system locally on their machines, it heavily relies on LLMs that are decoupled from the local environment.
To scale this approach, we need to shift our perspective and introduce this system as part of a centralized runtime platform.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="5" id="phase-5" >}}

## The platform

The platform is the backbone of your digital infrastructure. It provides the tools, services, and resources necessary to develop, deploy, and manage applications in a scalable, flexible, and efficient manner.

It is therefore the platform's role to provide access to the Large Language Models (LLMs) that empower the agentic loop we've just described.

The platform provides not only the model but also the inference engine. When you use an agentic system like Claude Code or Antigravity as an individual developer or *solopreneur*, the platform is directly managed by the providers of those systems (such as Anthropic or Google). In these scenarios, the agentic system and the platform are tightly coupled, making it difficult to run your own models or achieve digital sovereignty.

However, to gain more flexibility and control over privacy, costs, and data regulations, a key objective for an organization should be to encapsulate these capabilities within an internal platform.

Ultimately, the goal is to make the infrastructure completely transparent to the agentic system, allowing you to seamlessly swap out the runtime, the hardware, or the underlying model without breaking the loop.

Beyond that, the platform's job is to **streamline** both the software development lifecycle and the runtime environment. Therefore, beyond just serving as an execution engine, it must provide seamless access to the organization's broader capabilities and domain knowledge.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="6" id="phase-6" >}}

### Providing access to the Information System

When we described the agentic loop, we considered local tools in our example (like `write a file`). However, its true power comes from accessing resources outside the boundaries of the local development environment.

The agentic loop may require information from the Information System (IS) to complete its context—such as the documentation of a particular system or access to an authentication provider.

To guarantee the autonomy of the process (which is key in agentic development to meet expectations for velocity and reliability), the system must be able to access these resources as a service.

The platform should expose these capabilities in a way that can be easily consumed by the agentic system.

Historically, we used REST APIs to expose services, but the Model Context Protocol (MCP) is the current standard for tool access. Think of MCP as a transport layer with universal connectors that easily plug into agentic systems.

In my opinion, the platform will soon provide another kind of service: agentic services. I discussed this type of collaboration in a [previous article on my blog](/2026/06/25/from-isolated-agents-to-agentic-mesh-orchestrating-sdlc-with-a2a-and-ap2.html). Therefore, we must anticipate that the platform will eventually host multiple agentic systems, requiring robust Agent-to-Agent (A2A) transport protocols.

Now that we have set up all the plumbing required to develop a solution, let's return to what truly matters—and what remains, as of today, a fundamentally human task: framing the problem and designing the solution.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="7" id="phase-7" >}}

## Scoping a problem and designing a solution

Let's take a step back. We now know how the AI builds a solution, but without proper scoping, the agentic loop will just burn through tokens, yielding—at best—an expensive Proof of Concept (POC).

What matters most for a business is solving real user problems. In the accompanying diagram, **designing the solution** is represented as a single box, but it is actually a massive phase (and one where AI can also assist). I won't dive deeply into product management here, but the goal remains unchanged: design a feature that delivers true value to the end user while aligning with company strategy. **If your organization uses [Product Requirements Documents (PRDs)](https://en.wikipedia.org/wiki/Product_requirements_document), this is exactly the phase where they are drafted (and could use the agentic loop).**

This phase is an improvement of the existing design phase.

The real paradigm shift happens in the **specification phase**.

The goal of this phase is to translate that product idea into a format the implementer can perfectly understand. Before AI, the implementer was a human engineering team, and the specifications were written accordingly.

Today, the developer is an AI, and your specification literally becomes its execution context. This means the team scoping the work must understand how an agentic system "thinks" to feed it efficiently. They can no longer just describe the end goal; they must break the work down into discrete, scoped tasks that the AI can implement step by step without hallucinating.

Only when this AI-optimized specification is ready do we declare the intent to the system and let the loop run.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="7.2" id="phase-7.2" >}}

### Developing the solution

Now that we have the complete specifications and a clear intent, we can hand the work over to our developer—the AI agent.

Because this SDLC is an iterative cycle, the output of the agent's execution doesn't just go straight to production. It feeds back into the design phase to verify that the solution truly solves the underlying problem before anything is deployed or delivered.

Now that we have this lifecycle in mind, let's zoom out and look at the big picture.

{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="8" id="phase-8" >}}

### The Agentic SDLC

Our story of loops is almost complete; here is the full picture. 

Now, what is truly missing from this representation is the iteration within the design phase, as explained earlier.

The number of loops and the overall optimization are really a matter of execution. They depend heavily on the skills and ability of the human in control to correctly scope the system and react efficiently.

Let's now see how to organize this. 

So far, we have remained vague about the design phase. What about the functional design, the technical design, and the guardrails? Are they all handled and specified in this phase? Can some of them be inherited from elsewhere?

Well, it is a matter of organization. Let's explore how to structure this lifecycle to make it work efficiently.

{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="8.1" id="phase-8.1" >}}

## Organization

### The stream-aligned team

Let me assert a strong opinion that may spark some controversy here.

We are facing a massive paradigm shift.

Previously, the stream-aligned team was responsible for the actual development of the application. To me, the development of the application is no longer the responsibility of the stream-aligned team. However, the stream-aligned team remains strictly **accountable for** the result.

Of course, there are nuances to consider, and reaching this state is a goal that requires a transition phase. When bootstrapping AI-based development in an Agentic SDLC as described here, the stream-aligned team will initially own the agentic development.

But as a quick side note—before digging into the organization of this team and explaining how they can **enable agentic development** while remaining purely focused on feature development—let's consider a scenario that will arise in the near future. This shift will have a strong impact on the people and skills required in product management.

Adapting to an agentic model means the stream-aligned team must do more than just understand how to pilot the system...

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="8.2" id="phase-8.2" >}}

#### Solutions designed by stream-aligned teams will imply building agents

...they must also understand how agentic systems work in general. They need to know what escalation processes to set up, what the guardrails are, and how the AI makes its decisions.

Because tomorrow, the digital solution to a problem will likely involve building an agent.

Now let's get back to the need for enablement.

{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="9" id="phase-9" >}}

### The need for enablement

Steering an agentic system is neither easy nor straightforward. Even setting the system up is a massive undertaking in its own right (and even more difficult if it relies on elements of the platform that are niche or unusual).

We have focused heavily on designing the solution, but the generated output must also adhere to industry best practices. 

More importantly, it must respect **the organization's internal standards**, both in terms of technical architecture and aesthetic design (you don't want the AI generating a blue button if your design system strictly uses purple).

By design, an LLM can generate code that aligns with the broader market's state of the art, simply because of its training data. However, it must be explicitly instructed to follow your company's specific guidelines, conventions, and security policies.

We cannot rely solely on the stream-aligned team to manage all this context. Doing so would impose a massive cognitive load, which would ultimately penalize the delivery of the solution.

As discussed in my previous article, we must decompose the system and rely on an **enabling team** to address these cross-cutting concerns. Let's evaluate some of its tasks.

This enabling team is composed of technical experts, tech leads, and solution architects. They understand how the system works and help others use it effectively.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="10" id="phase-10" >}}

#### Enabling the stream-aligned team to use the agentic system

It is important to understand that the enabling team is not merely a support desk at the service of the stream-aligned team. Instead, both teams work in close collaboration, sharing a unified goal: delivering reliable and trustworthy solutions.

The enabling team's primary task is to set up the agentic system so that it operates with maximum autonomy—empowering the stream-aligned team to achieve the right result right from their very first prompt.

If the system requires a specific framework to support not just the application build, but the design phase as well, it is the enabling team's responsibility to integrate it. (For example, if your workflow relies on BMAD, this is the stage where it gets installed and configured.)

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="11" id="phase-11" >}}

#### Providing technical context
Another key responsibility of the enabling team is to provide the technical information needed to guide the agentic system. This technical layer complements the functional context brought by the stream-aligned team.

Specifically, this includes defining system-wide standards—such as the primary programming language, the designated authentication service, or even the specific design system to be used.

Once this technical baseline is established, it cannot remain siloed. Remember, we are scaling this Agentic SDLC across an entire organization.

Therefore, an additional responsibility of the enabling team is to share these practices, guardrails, and configurations. This cross-pollination eases the workload for other enabling teams and ensures systemic consistency across the company.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="12" id="phase-12" >}}

#### Automating governance and knowledge reuse
So far, this knowledge consolidation is entirely human-driven and is only enforced through manual governance rules. To be truly efficient at scale, this governance must be automated.

The industry often talks about "shift-left"—the practice of moving testing, quality, and security checks as early as possible in the development process, often before a single line of code is written.

In our agentic loop, human approval gates create bottlenecks. To guarantee the system's autonomy and velocity, these guardrails and standards cannot remain in a Wiki; they must be served directly by the underlying platform as consumable services.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="13" id="phase-13" >}}

## Standards by design

The platform should industrialize these practices. This means transforming them into consumable services that are packaged and managed as internal products.

These platform products must support the entire agentic loop, **not just the initial context**. In fact, the planning, tool execution, and observation phases should all inherit strict directives and pass through deterministic validation gateways to ensure trust and compliance.

By embedding these guardrails directly into the infrastructure, the platform guarantees that the agentic system remains both highly autonomous and perfectly safe.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="14" id="phase-14" >}}

### Empowering each element of the agentic loop

Now that the platform provides these capabilities as a product, we must adapt the agentic loop to leverage them at every step of the workflow. Think of these standards not as bureaucratic inertia to fight against, but as a secondary engine propelling the agentic loop forward.

What might this look like in practice? Let's take the development of a UI as an example:

* **Context:** The system is provided with specific guidelines to design a webpage according to the company's current design system.
* **Action/Tools:** The AI uses a dedicated tool to fetch and select only pre-approved UI components (e.g., buttons with the correct styling).
* **Observation:** The AI's output is evaluated by a programmatic mechanism that scans all color declarations to verify they strictly match the available corporate palette.

This goes far beyond mere context engineering.

As noted in a previous article, Markdown directives in prompts often exist simply to compensate for the current weaknesses of LLMs—they are meant to disappear as models naturally evolve. In contrast, programmatic guardrails and dedicated tooling actively *augment* the power of the agentic loop. This is the true value of treating platform capabilities as an internal product.

But what about the enabling teams? What becomes of them?

If they succeed, they vanish.
{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="15" id="phase-15" >}}

## Measuring the enabling team's success

The true measure of an enabling team's success is its ability to vanish—leaving the stream-aligned team with complete autonomy.

However, they do not vanish forever. Because industry practices evolve and products constantly change, spinning up and dissolving these teams is a natural part of the organizational lifecycle.

The crucial point is that no institutional knowledge is lost when the enabling team steps away. Because they have codified their expertise into the platform as automated guardrails and reusable services, the stream-aligned team can continue operating safely and autonomously long after the enabling team has moved on.

It is a natural phase of the organizational lifecycle.

{{< /scrollytelling-step >}}


{{< /scrollytelling >}}

## Conclusion

What we've seen on this journey is a vision for organizing around agentic development, rather than an out-of-the-box implementation plan.

As the saying goes, *Adapt is stronger than Adopt*—there is no magic recipe for this transformation. Take this as my personal conviction and a vision of how an organization could be set up to truly leverage the value of agentic development in the future.

One final disclaimer: we did not cover the **Complicated-Subsystem Team** in this picture. In *Team Topologies*, this is a specialized team that provides deep expertise at various levels to support different teams:

* They can help the stream-aligned team find the perfect algorithm for their product.
* They can assist the enabling team in customizing frameworks and tooling (such as BMAD).
* They can, and should, bring value to the platform by selecting the best models, optimizing costs, and fine-tuning inference engines.

Ultimately, it is by continuously improving the platform that the organization will be able to deliver software more efficiently and effectively.

History reminds us that raw technology doesn't scale on its own. Just as modern software delivery was defined not by Linux cgroups or Docker, but by platforms like Kubernetes, the future of AI development will belong to the platforms that orchestrate it.

Let's make AI work ~~on your machine~~ in your organization.
