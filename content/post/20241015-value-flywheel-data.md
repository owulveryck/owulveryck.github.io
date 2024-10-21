---
title: "How to Activate the Value Flywheel Effect with Your Data"
date: 2024-10-15T12:15:33+01:00
lastmod: 2024-10-15T12:15:33+01:00
images: [/assets/value-flywheel/value-flywheel-data.webp]
draft: false
keywords: []
summary:
tags: []
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
# You can also define another contentCopyright.
contentCopyright: false
reward: false
mathjax: false
---

In today's hyper-competitive world, businesses no longer rely solely on gut decisions or intuition; they depend on **data-driven insights** to stay agile and make fast, smart decisions.
However, data alone isn't the answer; it's the enabler to create momentum on a **business & technology flywheel**:
a model where data drives decisions, decisions drive actions, and those actions drive value, propelling the business forward in a self-reinforcing cycle.

In a [previous post](https://blog.owulveryck.info/2024/04/09/data-as-a-product-and-data-contract-an-evolutionary-approach-to-data-maturity.html), I used a model to explain how data could cross the borders of applications and domains to bring increasing value at the organizational level.

Here is a synthesis of the development of the post: the X-Axis is about the increase of certainty of the decision correlated to its supply, and the Y-Axis the diffusion of the data correlated to its demand.

![An S curve representing the evolution of the data, the X is the certainty and Y is the ubiquity. There is a division: the bottom of the S is raw data, the middle is curated, and the top is authoritative. The inflection point is indicated as data-contract.](/assets/data_certitude.svg)

In this post, I am using YAML (_Yet Another ModeL_ - sic) to determine how a data culture, technology and organisation can activate a value flywheel effect from a joint technological and business perspective.

By the end, you’ll gain actionable insights on analyzing and effectively implementing these principles within your own organization.

### Description of the Problem

In the world of **analytics**, data holds immense potential to drive **business value** by generating actionable insights. Yet, it often remains constrained as merely an **enabler of solutions** within a specific **business domain** (the problem space). Because data is typically treated as a tool to address isolated business challenges, it rarely crosses the boundaries of its domain. This limits its ability to contribute to broader **technological solutions** or drive value across the entire organization.

This misalignment often stems from **organizational silos**. Even worse, Data and business teams frequently operate in isolation, leading to disjointed efforts where insights, no matter how robust, are either misunderstood or underutilized.

To address this, companies need to focus on creating a **system** that fosters collaboration between these teams. This is where the analogy of a **wheel belt pulley system** becomes relevant. Just as a pulley system transforms motion and power, an aligned approach to data transforms raw information into business value.

Imagine a **flywheel** as the **driven wheel** of the system—this represents the combined **technological and business value**. **Data** serves as the **belt**, facilitating the transfer of knowledge, context, and value across the organization. The key, however, is identifying the **driver wheel**: a smaller flywheel that focuses specifically on **data maturity**, accelerating business outcomes.

![A wheel belt pulley system: a value flywheel described as "data" is the driver, a joint business and technology flywheel is driven by a belt that feed the driver with business initiatives. The driven flywheel is fed with data-as-a-product](/assets/value-flywheel/value-flywheel-data.webp)
Let's see now how companies can define and build this **data-centric driver wheel** in four steps, which will act as the engine powering the broader value flywheel of the organization.

### **The Value Flywheel Effect as a model**

#### The original concept

I have met the concept of the value flywheel in the book [The Value Flywheel Effect](https://itrevolution.com/product/the-value-flywheel-effect/) by David Anderson, with Mark McCann and Michael O’Reilly.

I love the introduction with the idea of momentum and I will simply copy here the very beginning of the book that you can find on the [website of IT Revolution](https://itrevolution.com/articles/what-is-the-value-flywheel/) for free:

> Momentum is a strange thing.
> It’s difficult to imagine what it will feel like and takes a great deal of effort to achieve.
> When we learn to ride a bicycle, for example, it feels clunky and awkward at first.
> It’s hard to get the wheels turning in the beginning, and our frustration is often evident.
> But our teacher assures us that it will pass.
> When we finally start to build momentum, the exhilaration takes our breath away.
> Every push of the pedal gets easier and takes less effort.
> Suddenly we can focus on the larger experience of gliding through a beautiful forest or tree-lined street.
> The value of our hard work is evident, and we can now continue to reap the benefits with less and less toil.

To fully understand the **value flywheel effect**, let’s first consider it from a broader business perspective.
The most well-known application of the flywheel effect comes from Amazon, where small wins in customer satisfaction lead to more traffic, which attracts third-party sellers, which in turn results in lower costs and prices, creating further momentum.
Each element strengthens the next, creating an ever-accelerating hub of growth, efficiency, and differentiated value.

The flywheel model describes a **self-reinforcing system** of momentum, growing stronger and easier to maintain as different actions compound over time.
Instead of relying on giant, one-time shifts or changes, the flywheel works by multiplying the effects of **small, consistent improvements**.
Once it's moving, it's hard to stop.

#### The Key tenets and the personas

I will refer to the model from the book *The Value Flywheel Effect*, which is the third iteration of this concept (following Amazon’s adaptation and the original idea by [Jim Collins](https://www.jimcollins.com/concepts/the-flywheel.html)). This model is thoroughly explained in the book, and I highly recommend reading it. In summary, the value flywheel is broken down into four steps:

**Clarity of Purpose** → **Challenge and Landscape** → **Next Best Actions** → **Long-Term Value**

Each of the four steps drives the next one. **Step by step**, **iteration by iteration**, we overcome inertia (in a truly agile way) and build momentum.

{{< figure src="/assets/value-flywheel/the-value-flywheel-effect.webp" alt="The value flywheel in four phases: I clarity of purpose, II challenge and landscape, III next best action, IV long-term value. The wheel connects the four phases and loops from phase IV back to phase I." title="Illustration of the Value Flywheel Effect (Adapted from the model by David Anderson)" >}}

To summarize, here are the key tenets and personas of the Value Flywheel Effect:

{{< figure src="https://itrevolution.com/wp-content/uploads/2022/10/Screen-Shot-2022-08-03-at-3.06.48-PM-1024x819.png.webp)" title="12 Key Tenets of the Value Flywheel Effect - David Anderson, Michael O’Reilly, Mark McCann - (c) ITRevolution" >}}
[Source](https://itrevolution.com/articles/12-key-tenets-of-the-value-flywheel-effect/)

Now let's see if we could adapt those key tenets to the data driver wheel. The first element is to consider the existing approach in place in the organization (Adapt > Adopt)

### The Data Value Flywheel: Organizing the Data Ecosystem

To successfully apply the **Value Flywheel Effect** within a data ecosystem, companies need a structured approach to managing data across the organization. This is usually implemented thourgh the concept of the **data factory**. A data factory functions as the central engine that powers data-driven decision-making, ensuring that data flows freely between business and technical teams, transforming it into actionable insights.

At its core, the data factory is responsible for the **end-to-end data lifecycle**: from **data collection** and **curation** to **storage**, **governance**, and **exposure** for insights. It acts as a bridge between domains, enabling the seamless transfer of data from isolated business problem spaces to broader technological and strategic applications across the organization.

However, for the data factory to be effective, the **data team** must play a pivotal role within each domain.
There are several options of interactions possible depending on governance of the organization.

One solution that some organizations have embraced is a **4-in-the-box (4ITB)** model, initially [pioneered by Walmart](http://thomasmisner.com/download/Thomas-Misner_Product-Ways-Of-Working-Presentation.pdf).
The 4ITB model brings together **cross-functional teams** involving stakeholders like Product Managers, Engineers, Data Experts, Business Representatives, and potentially others (N-ITB) to collaboratively solve business challenges.
It aims to establish both **horizontal alignment** (across departments) and **vertical alignment** (with leadership's vision) for every project.


In this context, the data team not only serves the domain’s immediate needs but also acts as a **catalyst** for driving **broader technological innovation**. Their role is to transform data from a **domain-specific enabler** into a **strategic asset** that fuels the larger flywheel of business and technological value.

Now, let’s look at how we can apply the **four phases of the value flywheel** to this data-driven approach.

#### Phase 1 - Clarity of Purpose

The first phase is about establishing **Clarity of Purpose**.
To create business value from data, the people working in the **data factory** need to understand the purpose of their work and how it supports broader company objectives.

- **Key Tenets**: Alignment on **business goals**, connecting data initiatives directly to business outcomes.
- **People in Charge**: Leadership, Executive Sponsors, and Domain Experts.
- **Goal**: Ensure that both data and business teams are working towards the same strategic goals, thus preventing misalignment in expectations or priorities.

At this stage, the central data factory needs to focus on aligning its capabilities with **business priorities**.
This clarity will guide the roles and direction for the team moving forward and ensures that data outputs lead to insights that can create real value.

The author of the original book proposes to define a data-informed North Star in this phase, and set what differentiator on the market we are seeking with the current data platform.

TODO: complete with the WHY and the difference with the why of simon sinek // Reference the linkedin post

#### Phase 2 - Challenge and Landscape

Next comes **Challenge and Landscape**, which focuses on identifying **key obstacles** and understanding the **current technological and business landscape**.
This phase might reveal that certain systems and processes need to be updated to meet modern data demands or **gaps in domain-level data ownership**.

- **Key Tenets**: Awareness of technological **debt**, bottlenecks, and organizational inefficiencies.
- **People in Charge**: Engineering: Technology Leads, Data Engineers, ...
- **Goal**: Examine the current infrastructure, tools, and systems in place and identify what challenges exist that may slow downs the spinning of the flywheel.

During this phase, it's crucial for the **data factory team** to get an accurate picture of **which challenges**—be it technical, cultural, or operational—are preventing them from scaling or accelerating their value generation.
By mapping these obstacles, I help guide my clients towards addressing the most pressing issues while also identifying new opportunities.

#### Phase 3 - Next Best Action

Once the landscape is clear, the focus moves to identifying the **Next Best Action**:
what small, **incremental** improvements can be made that will begin turning the flywheel? At this stage, it’s about **empowering teams to take ownership** and starting to decentralize responsibility without overwhelming them.

- **Key Tenets**: **Iterative progress, agile development**, unlocking quick wins, and decentralizing responsibilities.
- **People in Charge**: Domain Leaders, Product Teams, and Data Engineers.
- **Goal**: Focus on the **minimum viable actions** necessary to eliminate the identified pain points, and **prioritize pragmatic solutions** that bring value without over-complication.

Here, the **data factory** plays a pivotal role—it provides the foundational capabilities needed for domains to take charge of their **own data products**.
Small, actionable steps, like improving data accessibility or providing scalable infrastructure, can catalyze bigger leaps.

#### Phase 4 - Long-Term Value

Finally, the last phase revolves around establishing **Long-Term Value**.
This is where the flywheel effect becomes self-sustaining, where each domain confidently handles its own data products, and the **data factory** transitions to more of a **supportive, enabling role**.

- **Key Tenets**: A federated model where all teams are contributing to **shared goals** and creating sustainable value.
- **People in Charge**: Chief Data Officers (CDOs), Domain-specific Data Teams, Architects.
- **Goal**: Architect the system to work at **scale**, ensuring consistency, governance, and long-lasting value for the entire company.

At this stage, the **central data factory** should have evolved into more of a **platform provider**, focusing on empowering the individual domains while ensuring that the infrastructure is reliable, scalable, and repeatable.
By driving toward a **data-mesh**, this phase enables the system to generate value continuously, accelerating innovation for the business.

### Conclusion: Activating the Flywheel

The **value flywheel** is activated phase by phase, and the central **data factory** plays a major role throughout the journey towards a federated data management system, ensuring **cross-domain consistency** and delivering scalable solutions.
By following this structured approach, the data landscape becomes a coherent framework where each part of the organization—not just the central data team—can contribute to the company’s **long-term success**.

The shift toward a **data-mesh** architecture ensures that all domains can spin the flywheel on their own, transforming data into actionable insights and kickstarting a virtuous cycle of value, innovation, and acceleration.
