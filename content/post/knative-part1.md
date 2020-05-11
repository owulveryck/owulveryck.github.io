---
title: "Divide et impera - part I: coupling"
date: 2020-04-26T15:47:57+02:00
description: "This article is about coupling in IT; divide et impera - divide and conquer"
draft: true
---

In 2020, an old debate has risen from the ashes:

_Monolith or micro-services?_

In the past years, I have been working on IT products designed to work at scale.
The micro-services pattern suits this need and is, therefore, the paradigm chosen _de facto_.

An application composed of micro-services is supposed to be scalable; the main pitfall is to think and design accurate service contracts to avoid the antipattern of a distributed monolith.
The secret is to manage the coupling between the components in the same way developers have learned to handle coupling inside the code over the years.

By mastering the coupling, the design paradigms have shown excellent results on the quality and reliability of software.

How do those paradigms apply to AI-based applications?

In this series of articles, I plan to dig into the design of code and hosting applied to AI software.
This first article is about concepts.
Then I may write a second article, more technical, to expose a concrete example and illustrate the concepts.
Eventually, a third article will show how to use managed services to transform products into commodities[^1].

[^1]: Chapter 7 of Simon Wardley's book ([Finding a new purpose](https://medium.com/wardleymaps/finding-a-new-purpose-8c60c9484d3b)) gives more explanation of the evolution of a product to a commodity.

## What is _coupling_?

In physics, two objects are said to be coupled when they are interacting with each other (source [Wikipedia](https://en.wikipedia.org/wiki/Coupling_(physics))).

Two elements are strongly coupled when a perturbation on the first element induces a disturbance on the second that impacts its behavior significantly.
A corollary is that two systems are weakly coupled if the perturbation induced on a component from the excitation of another one is low.

In chapter 12[^2] of his book [_In Search of Certainty_](http://markburgess.org/certainty.html), [Mark Burgess](https://twitter.com/markburgess_osl) discusses the notion of coupling in physics. Some of the conclusions exposed in chapter 13 are [^3]:

- strong coupling leads to precision but the fragility;
- weakly coupled components can be rearranged, taken apart and made into something else, assembled into a meaningful composition;

[^2]: Chapter 12: Molecular and Material Infrastructure: _Elastic, plastic and brittle design_, page 383
[^3]: Chapter 13: Orchestration And Creative Instability: _Or why the conductor does not promise to blow every trumpet_, page 422

A software is roughly a composition of elements (objects, services, functions, etc.) that are interacting, hence the paradigm of coupling applies to software engineering.

From my experience, the strength of coupling is roughly evaluated but not correctly measured (nevertheless, it feeds numerous discussions about the quality and reliability of the application).

Actually, in modern software, **coupling** arises at different levels, **inside the code**, *and* at the **architecture** level.

Let's first analyze the coupling of an AI base application, and then see how the notion of scalability may lead to the need for decoupling the application in different components.

### Application in the code

In the book [Reliable software through composite design](https://archive.org/details/reliablesoftware00myer) (written in 1957!), Glenford Myers exposes why designing a system to master coupling is essential. A strongly coupled system induces complexity and a high cost of maintenance. It makes the software "fragile" (as seen in the previous paragraph).

To bring an intuition of the problem, let's consider this illustration (compiled from Myers's book and [Notes on the Synthesis of form](https://en.wikipedia.org/wiki/Notes_on_the_Synthesis_of_Form) from Christopher Alexander).

{{< blockquote author="Glenford Myers" source="Reliable software through composite design" >}} 
Consider a system of 100 lamps where a lamp could represent a statement or a segment in a program. Each lamp can be in one of two states, off or on.
The lamps can be connected together in any fashion. If a lamp is on, the probability of its going off in the next second is 50% if at least one adjoining lamp is on.
If a lamp is off, it will stay off as long as all of the lamps directly connecter to it are off.
{{< /blockquote >}}

| Test case                              | Time to reach equilibrium                   |
| -------------------------------------- | ------------------------------------------- |
| no connection between any lamp         | around 2s (2<sup>1</sup>s)[^4]              |
| all lamps are connected                | around 2.4e24 years (2<sup>100</sup>s) [^5] |
| ten independent groups of fully connected lamps | around 17 minutes (2<sup>10</sup>s) |

[^4]: average time for a lamp to go off
[^5]: average time for all lamp to go off

The notion of equilibrium is more or less relative to what computer scientists call stability.
The time-based illustration gives the feeling that coupling is risky and has a substantial impact on the velocity of new development.

To fight this problem, IT developers are applying programming paradigms, tools, and methods that are evolving over the years.

[Connascence](https://en.wikipedia.org/wiki/Connascence) is such a metric to measure complexity (mainly in the object-oriented paradigm).

The goal of this tool is to evaluate a level in the code. Then to apply any technique to reduce the connascence level, from the most dangerous (Connascence of identity) to the less dangerous (Connascence of name). See the [Wikipedia page](https://en.wikipedia.org/wiki/Connascence) for the complete list and classification.

<center>
![Connascence](https://gist.githubusercontent.com/owulveryck/df65079edbd273d33805f00e3d5d51a6/raw/5ad99b3c3833120e780a34f9d4c2a54473d219a2/connascence.svg)
</center>

Let's now consider AI-based applications where the business logic is using compute libraries with a lifecycle that is independent of the development language.

## In modern applications: the AI case

### About the link between software 1.0 and 2.0

An AI model is a mathematical representation associated with some values. Let's call it a software 2.0 (this term has been introduced in late 2017 by Andrej Karpathy in a [blog post](https://medium.com/@karpathy/software-2-0-a64152b37c35), and is slowly becoming a common language in the data world).
In this paradigm, software 1.0 is the result of regular code. Its usage is to glue the interfaces, I/O; it acts as a host for the AI/ML algorithm.

Creating a software 2.0 is roughly transcribing this mathematical model into something understandable by a computer, e.g., a sequence of bytes. Then, we need two different types of software 1.0:

- one will act as a helper for the training phase (the build phase of the software 2.0)
- the other one will handle the inference (the runtime of software 2.0)

The algorithm act as a link between the components. On top of that, the algorithm itself is transpiled into some code, usually using the same language as its host. This makes software 2.0 an entity of software 1.0; therefore, this is comparable to the [connascence of identity](https://connascence.io/identity.html).

<center>
![Venn diagram](https://gist.githubusercontent.com/owulveryck/df65079edbd273d33805f00e3d5d51a6/raw/e39ed3495f9c4ab99f5c32a39ec91cacbabcf6cc/diagram1.svg)
</center>

#### Soften the link

Dealing with this connascence of identity is mostly dealing with this duality of the training/inference phase.

One idea is to turn the entity representing the software 2.0 into data. Therefore, we lower the connascence from _dynamic connascence_ to _static connascence_.

One way to achieve this is to think of the model (software 2.0) as data.
Training and inference software must agree on the data type of the software 2.0. 

One way to do that is to use a domain-specific language (DSL) to represent the software 2.0. Using the mathematical representation is a perfect example of this.

This DSL could, eventually, be encoded into a format that would become the Babelfish for AI[^6].

[^6]: I gave a lightning talk about it at dotAI in 2018: [_software 2.0 a Babelfish for deep learning_](https://www.youtube.com/watch?v=Gf-pmc7Mykc)


![babel fish](/assets/babel-fish.jpg)

Developing AI applications is a combination of several processes that are loosely coupled:

- creating an algorithm and describing it with a mathematical representation;
- expressing the formulae in the IA DSL;
- creating a software 1.0 that can read the formulae and feed it with data for learning;
- creating another software 1.0 that can read the formulae and feed it with realtime data for inference;

_Note_ [Open Neural Network eXchange (ONNX)](https://onnx.ai) is an example of such a DSL, but digging into the technical implementation is out of the scope of this article and will eventually come in a future post.

{{< blockquote author="Titus Winter">}}
Software engineering is programming integrated over time.
{{< /blockquote >}}

So far, decoupling the process allows doing engineering work and developing a maintainable application[^7].

[^7]: as mentioned by Russ Cox in his blog, _It's worth seven minutes of your time to see [Titus Winter's presentation of this idea at CppCon 2017](https://www.youtube.com/watch?v=tISy7EJQPzI&t=8m17s), from 8:17 to 15:00 in the video._

Let's now consider the runtime, and whether we can run it smoothly "at scale."

## Running AI at scale

Before going further, let's define the notion of scale:

{{< blockquote author="Mark Burgess" source="Smart Spacetime (page 100)" >}} 
[_scale is_] a region of uniformity, or compatibility, a measure of compatible things, whether by distance, weight, size, color, etc.
{{< /blockquote >}}

According to this definition, a system can run "at scale", if it's processing power is compatible with the expected usage.

Luckily, the decoupling we made before made it simpler to adapt the processing power for both the training and the inference process.

On both parts, data usually need to be pre-processed before they can feed the software 2.0. This is where the data-pipeline comes in.

### Data pipeline and time reference

A data pipeline is a series of data processing steps.

To illustrate, let's use [C4](https://c4model.com), a method to design an application (and to take care of coupling).
This method defines the concepts of _containers_ and _components_.

- a container (_not docker!_) which is _something that needs to be running in order for the overall software system to work._
- a component is _a grouping of related functionality encapsulated behind a well-defined interface_

Building a data pipeline brings the ability to split the components into the containers.

Each container defines its own scale. As a consequence, the overall system can adapt smoothly, at scale.

The different _containers_ communicate through channels of communication where they can exchange messages.

Claude Shannon explains in detail the notion of communication channel is [the mathematical theory of communication](https://en.wikipedia.org/wiki/Information_theory#Channel_capacity). For our explanation, let's classify the communication channels: _rigid_ and _robust_.

- A rigid channel is something that couple the components strongly, therefore information is received synchronously from point A to point B (such as an in-memory semaphore, for example).
- A smooth channel is an element that carries information that will eventually flow from point A to point B.

The nature of the communication channel is the next axis of work. As seen in the first paragraph, a rigid channel will lead to precision but fragility, and a smooth channel is more robust but less precise.

For example, with a smooth channel, even _time can run backward locally[^8]_, which is very useful for elements of the pipeline that deals with databases.

[^8]: _time can run backward locally, as long as the process doesn't depend strongly on what happens around it._ Mark Burgess - [Bigger, Faster, Smarter - Episode 1](https://youtu.be/lDFQiS9T_xk?t=2598)

### Event-driven architecture

The event-driven architecture is a set of models that allows weakly coupling in a process-pipeline.

For example, A process emits an event in a communication channel, and this event is received by any process listening on the same channel.

Back to our AI system 

**TODO**

#### Separate the infrastructure from the application

> A story cannot be written down without a medium; a process cannot exist without an enabling infrastructure - Mark Burgess

**TODO**

## Conclusion

AI, deep learning, and other machine learning mechanisms are continually improving. POCs, papers, startup, and academic works are opening new perspectives to business and research to the digital world.

But the diffusion of innovation is conditioned by the ability to develop and run qualitative software that can run reliably.
We saw in this article that mastering coupling is one of the keys to shit from application development to software engineering.

A lot of products exist nowadays to facilitate the application of the concepts described in this post. Some of them are even becoming commodities, thanks to cloud providers.

In upcoming articles, I will describe an example of AI application development on top of products such as Kubernetes. Then I will end this series by showing how a managed version of some service help in managing the TCO of the software.

Stay tuned!

