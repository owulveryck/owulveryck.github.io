---
title: "Four Pillars, One Loop: A Manifesto Is Not a Compass"
slug: "manifesto-through-the-loop"
date: 2026-07-17T10:00:00+02:00
images: [/assets/manifesto-through-the-loop/four-pillars-one-loop.svg]
draft: false
summary: "I liked the AI-Driven Development Manifesto until I saw it implemented as guiding principles. A manifesto is raw material; guiding principles are what you refine from it, and the refining is work nobody can do for you. This article does that work pillar by pillar, reading each value (Method over Model, Ownership over Delegation, Understanding over Acceptance, Outcome over Output) from inside the agentic loop, and ends with principles a platform can execute."
tags: ["agents", "platform", "manifesto", "engineering-practices"]
categories: ["dev"]
author: "Olivier Wulveryck"
toc: true
comment: false
mathjax: false
---

*A manifesto is something you sign. A guiding principle is something you execute. The whole argument of this article fits in that distinction.*

## Introduction

When I first saw the [AI-Driven Development Manifesto](https://www.ai-driven-development.org/) (AIDD), I thought it was an excellent idea. Four values, twelve commitments, written by developers against vibe coding, in the "X over Y" lineage of the Agile Manifesto. At first reading, the values are easy to agree with. I could sign them.

Then I started to see the manifesto implemented as guiding principles, and that changed my reading. In the details, some of its foundations worry me. They graft an AI onto the existing lifecycle instead of rethinking the lifecycle around the agent. I understand the original fight: in the vibe-coding era, that caution made sense. But agentic development has moved far beyond vibe coding, and the models have evolved. Applied to agentic engineering, these policies may not protect the value of the end product (a goal I share). They may just add inertia, and prevent the efficiency the loop is built to deliver.

We should not fight the obvious. The agents write the code now; the humans do not. The question worth working on is how the team keeps the knowledge, not how to verify what an AI already does better than we do.

My June and July posts built such an agentic delivery engine: the loop described in [Codifying the Rules](/2026/07/02/sdlc-team-topologies.html) and amplified in [The Amplified Agentic Loop](/2026/07/07/the-amplified-agentic-loop-guardrails-as-accelerators.html). Read from inside that loop, each pillar lands somewhere more precise than its slogan. This article reads them one by one. First, the thesis they will illustrate.

## A manifesto is not a compass

A manifesto is a public, signed, dated declaration of values. It is always written *against* something: AIDD against vibe coding, as Agile was against heavyweight process. Its register is rhetorical by design: broad values, memorable slogans. Its function is real: rally a community, mark a moment, start conversations. This article exists because the manifesto worked.

Guiding principles are a different object. They are the decision rules a team uses daily. To be usable, they must be precise, testable, owned, and revised as the tooling evolves. In the agentic era they go one step further: they are codified into the platform, where the agent can observe them.

A manifesto is raw material. Guiding principles are what you refine from it. The refining is real work, and nobody can do it for you. The mistake is not signing the manifesto. The mistake is stopping there. The four pillars below show what that refining looks like.

## The loop, briefly

An agent **captures** an intent, **plans** the steps, **acts** through tools and **observes** the results. When the observation is unsatisfactory, it re-plans. The loop is self-correcting. This inner loop is the delivery engine of modern software. Around it, humans run the outer loop: design, specification, expression of intent, acceptance, judgment of the outcome. Stripped to its essence, the collaboration is two loops and a border: intent goes down, the result comes back up.

![The mental model: two loops, one boundary, two rhythms. The human loop is one ample stroke, in days; the agent loop is a dense segmented stroke, in minutes; ownership and accountability live on the boundary.](/assets/manifesto-through-the-loop/human-agent-loops.svg)

This image is the whole model: two loops, one boundary, two rhythms. The loop is a lens, not a law. I will keep to that: the loop proposes a reading of each pillar, and you judge whether the reading holds. Each section below returns to this same frame and lights up the part where its pillar lives.

## Method over Model

*"Bet on the method, not the model: every LLM release will fade."*

![Pillar 1 on the frame: the model is a swappable part plugged into the agent loop; the next release waits outside; the loop and its guardrails endure.](/assets/manifesto-through-the-loop/pillar1-swappable-model.svg)

Half right. Seen from the loop, the model is a swappable part inside the Plan step. Every release fades, and the loop does not care which brain it runs on. Betting a team's practice on a specific model is building on sand. So far, full agreement.

The method half does not hold. We have heard this pillar before, twice. In the 1990s it read "bet on the method, not on the tool": the method was Merise, SSADM or the V-model, and the tool was the CASE workbench (Computer-Aided Software Engineering: an integrated environment meant to produce the application from the models). I never met CASE; it predates my career, and I know it only from what it left behind. The 2000s replay, I met. The method was RUP, the Rational Unified Process: phases, roles, artifacts, and a toolchain to buy. The tool was Rational Rose, the UML modeling workbench. And MDA, the Model-Driven Architecture, promised once again to generate the code from the models. The tools faded, and the methods almost died with them, superseded by agile: Merise and RUP survive in textbooks, rarely in delivery. A method co-evolves with the tooling of its era, and follows its fate. Now look at what crossed both transitions intact: testing, code review, versioning, decomposition into small reversible steps. They predate those workbenches. They outlived them. They are what we wire into the agentic loop today, as guardrails.

The line between the two groups is not accidental. A method prescribes a process, and a process is shaped by its era's tools. A practice protects a property of the software itself: correctness, reversibility, decomposability. Properties outlive processes. And if "method" is stretched until it survives every era, it stops guiding anything. What endures is not the method. It is the practices, one level lower.

## Ownership over Delegation

*"You own what you ship, even what the AI wrote."*

![Pillar 2 on the frame: only the boundary is lit, with its two crossings (intent expressed, result accepted); the commits stay dimmed inside the agent loop, out of scope.](/assets/manifesto-through-the-loop/pillar2-ownership-boundary.svg)

On the value, full agreement, and the English words carry it better than my native French. French has one word, *responsable*, for what English splits in two; the trick of this pillar hides in that split. The agent is responsible: it does the work. The human is accountable: they answer for what ships. Responsibility moves with delegation; a manager who hands work to a junior hands the doing, not the answering. Accountability does not move, because it needs someone who can answer, and an agent cannot: it cannot be blamed, sanctioned or trusted. When it fails, the failure is corrected inside the loop, and nothing that ships is ever the agent's fault. The accountability that sat on the human before the agent arrived has nowhere else to go. It stays.

The mechanism is where the pillar wobbles: "every commit is signed by a human". A signature is a ceremony, not ownership. A manager countersigning every action of the team owns nothing more than before. In the loop's frame, ownership sits at the boundary: the intent expressed on the way in, the result accepted on the way out. Not on each commit inside the loop, where commits are the agent's working memory.

To be fair, a signature does more than claim ownership. It gives traceability, non-repudiation, an audit trail. Regulators care, and they are right to. But those needs argue for provenance, not ceremony. Record who expressed the intent, which plan was locked, what was accepted and by whom. The platform captures all of this at the boundary, better than a human countersigning generated commits ever could. Keep the audit trail. Move it to where the decisions happen.

## Understanding over Acceptance

*"Don't accept what you don't understand. The AI is your collaborator, not your replacement."*

![Pillar 3 on the frame: the two crossing points of the boundary are circled (the context in, the result out); autonomy is preserved inside the agent loop.](/assets/manifesto-through-the-loop/pillar3-mutual-understanding.svg)

The first sentence deserves to survive every era: accepting what you do not understand is how debt enters a system.

The second sentence gets the naming wrong. The AI does replace the coder. In the platform this series builds, the human types none of the delivered lines. The typing is gone; the engineer is not. The engineer moves up one level: deciding what enters the context, judging what comes out. A task replaced, a role promoted. "Collaborator" hides that shift, and the hiding has a cost. A developer who believes the AI is a pair programmer ends up inside the loop, reviewing each line as it is written. That destroys the autonomy that makes the loop valuable.

Named correctly, understanding is not supervision. It is a contract at the two border crossings. On the way in, the human decides what the agent needs to know: the intent, the constraints, the invariants, the non-goals. That selection is a comprehension task, not a copy-paste task. On the way out, the human understands what was done before accepting it. Not each line: what changed, why, and against which intent it was verified. The agent has its own side of the contract: it must understand the rules it operates under. That is why guardrails live inside the loop. The platform rejects a plan that touches files outside its scope, the same way a compiler rejects a type error. Two autonomous parties, understanding each other at the border. That is where the value comes from. And that is how the team keeps the knowledge: the intent and the verification stay human, even when the code does not.

## Outcome over Output

*"Writing code is easy. Building product is not."*

![Pillar 4 on the frame: the outer loop grows into the judge (outcome, in days); the inner loop shrinks into a dimmed producer of volume (output, in minutes).](/assets/manifesto-through-the-loop/pillar4-outcome-judge.svg)

Read through the loop, this pillar holds as written. The loop only adds the structure: the SDLC is nested loops. The inner loop iterates in minutes and produces code at a volume no team has ever seen. That volume is the vanity metric the manifesto warns about. The outer loop iterates in days, and it is the only judge: intent at one end, business result at the other.

The pillar also comes with a test. Name the number the business reads. If your dashboards count merged pull requests or lines produced, you are measuring the inner loop, and the pillar is violated, whoever signed it. Output is what the blue loop maximizes. Outcome is what the green loop answers for.

## Conclusion

If I signed with my own pen, the four values would read:

1. **Engineering practices over method, method over model.** Three shelf lives: models fade in months, methods in decades, practices cross eras.
2. **Accountability at the boundary over a signature on every commit.** Answer for the intent expressed and the result accepted, not for each step in between.
3. **Mutual understanding over supervision.** Spend human attention at the border: context on the way in, verification on the way out.
4. **The outer loop over the inner loop.** Outcome is the only number that serves.

Yes, the same format. The format was never the problem; the altitude was.

That is what I would sign. Here is what I execute, the guiding principles the platform enforces:

- A plan is enforced and its execution validated deterministically, by the platform; compliance never rests on the agent's goodwill.
- Architectural invariants live as executable policy served at planning time, not as wiki prose.
- The instructions an agent loads (knowledge, directives, packaged capabilities) come through the platform, and an instruction served this way is trusted: it cannot harm the system, and it is safe to use inside the loop. If two instructions contradict each other, the platform resolves the conflict deterministically; the outcome is never left to the agent's choice.
- Every compensatory artifact (a prompt directive, a CLAUDE.md rule) carries a sunset condition and an owner.
- Feedback must arrive inside the loop: a control that fires after delivery should be a non-event, one that rarely blocks. When it does, its lesson becomes a guardrail inside the loop: the same cause never blocks again, and never escalates to a human twice.
- The human intervenes at two points: the context handed at Capture, and the result at the exit of the loop.

None of these fits on a poster. That is their value. And they are not hypothetical: the first rule runs in the proof of concept that accompanies this series. There, enforcement means that an agent modifies only the files declared in a locked plan; an edit outside it is refused before the file changes. Other implementations are possible; the determinism is not negotiable. The first version of that policy was too strict and blocked legitimate edits; it was revised. That is what "revisable" means in practice.

The loop also tells you how to measure these principles. Two numbers are enough. The first is the count of inner-loop re-plans, the *(n)* in the figure: how many times the agent corrected itself before the result was accepted. It measures efficiency, and it rises when the context or the guardrails are weak. The second is the count of human escalations: how many times a problem crossed the boundary and forced a human to re-plan. It measures effectiveness, and for known causes it should trend to zero. Neither is a business number; the outer loop still judges the product. These two judge the machine.

![The loop dashboard: two counters on the frame. Human escalations across the boundary measure effectiveness and should trend to zero; inner re-plans (n) measure efficiency and should trend down. The product is judged in the outer loop.](/assets/manifesto-through-the-loop/loop-metrics.svg)

One objection deserves a direct answer, because it turns my own argument against me. If methods fade with their tooling, and these principles live inside the platform, are they not the most perishable of all? Yes. That is the point. A guiding principle should retire when its tooling does; that is why each one carries a sunset condition and an owner. A manifesto pretends to outlive its era; a principle is honest about its mortality.

So sign the manifesto. I might. Then do the work it cannot do for you: derive your principles from your own loop, and put them where the agent can observe them.

*Let's make AI work* ~~*on your machine*~~ *in your organization.*
