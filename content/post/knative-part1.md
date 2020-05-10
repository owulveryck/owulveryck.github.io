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

## Coupling inside the code

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
The time-based illustration gives the feeling that coupling is risky and have a substantial impact on the velocity of new development.

To fight this problem, IT developers are applying programming paradigms, tools and methods that are evolving over the years.

[Connascence](https://en.wikipedia.org/wiki/Connascence) is such a metric to measure complexity (mainly in object-oriented paradigm).

The goal of this tool is to evaluate a level in the code. Then the developer can apply any technique to reduce the connascence level, from the most dangerous (Connascence of identity) to the less dangerous (Connascence of name). See the [Wikipedia page](https://en.wikipedia.org/wiki/Connascence) for the complete list and classification.

<center>
![Connascence](https://gist.githubusercontent.com/owulveryck/df65079edbd273d33805f00e3d5d51a6/raw/5ad99b3c3833120e780a34f9d4c2a54473d219a2/connascence.svg)
</center>

Let's now consider AI-based applications where the business logic is using compute libraries with a lifecycle that is independent of the development language.

## AI application

### Coupling between software 1.0 and 2.0

An AI model is a mathematical representation associated with some values. Let's call it a software 2.0 (this term has been introduced in late 2017 by Andrej Karpathy in a [blog post](https://medium.com/@karpathy/software-2-0-a64152b37c35), and is slowly becoming a common language in the data world).
In this paradigm, software 1.0 is the result of regular code. Its usage is to glue the interfaces, I/O; it acts as a host for the AI/ML algorithm.

To create a software 2.0, there is a need to transcribe this mathematical model into something understandable by a computer, eg. a sequence of bytes. Then, we need two different types of software 1.0:

- one will act as a helper for the training phase (the build phase of the software 2.0)
- the other one will handle the inference (the runtime of software 2.0)

The algorithm act as a link between the components. On top of that, the algorithm itself is transpiled into some code, usually using the same language as its host. This makes the software 2.0 an entity of software 1.0; therefore this is comparable to [connascence of identity](https://connascence.io/identity.html).

<center>
![Venn diagram](https://gist.githubusercontent.com/owulveryck/df65079edbd273d33805f00e3d5d51a6/raw/e39ed3495f9c4ab99f5c32a39ec91cacbabcf6cc/diagram1.svg)
</center>

#### Decoupling software 1.0 and software 2.0

Dealing with this connascence of identity is mostly dealing with this duality of the training/inference phase.

One idea is to turn the entity representing the software 2.0 into data. Therefore, we lower the connascence from _dynamic connascence_ to _static connascence_.

One way to achieve this is to think the model (software 2.0) as data.
Training and inference software must agree on the data type of the software 2.0. 

One way to do that is to use a domain-specific language (DSL) to represent the software 2.0. Using the mathematical representation is a perfect example of this.

This DSL could, eventually, be encoded into a format that would become the bablefish for AI[^6].

[^6]: I gave a lightning talk about it at dotAI in 2018: [_software 2.0 a babelfish for deep learning_](https://www.youtube.com/watch?v=Gf-pmc7Mykc)


![babel fish](/assets/babel-fish.jpg)

Developing AI applications is a combination of several processes that are loosely coupled:

- creating an algorithm and describing it with a mathematical representation;
- expressing the formulae in the IA DSL;
- creating a software 1.0 that can read the formulae and feed it with data for learning;
- creating another software 1.0 that can read the formulae and feed it with realtime data for inference;

_Note_ [Open Neural Network eXchange (ONNX)](https://onnx.ai) is an example of such a DSL; but deeping into the technical implementation is out of the scope of this article and will eventually come in a future post.

{{< blockquote author="Titus Winter">}}
software engineering is programming integrated over time.
{{< /blockquote >}}

So far, decoupling the process gives the opportunity to properly do engineering work to develop a maintainable application.[^7]

[^7]: as mentioned by Russ Cox in his blog, _It's worth seven minutes of your time to see [his presentation of this idea at CppCon 2017](https://www.youtube.com/watch?v=tISy7EJQPzI&t=8m17s), from 8:17 to 15:00 in the video._

Let's now consider the runtime, and whether we can run it smoothly "at scale".

### Running at scale

Before going further, let's define the notion of scale:

{{< blockquote author="Mark Burgess" source="Smart Spacetime (page 100)" >}} 
[_scale is_] a region of uniformity, or compatibility, a measure of compatible things, whether by distance, weight, size, color, etc.
{{< /blockquote >}}

_TODO_ 
#### Data pipeline and time reference

A data pipeline is a sseriesof data processing steps.

> time can run backward locally, as long as the process doesn't depend strongly on what happens around it.
https://youtu.be/lDFQiS9T_xk?t=2598

### Separate the infrastructure from the application

> A story cannot be written down without a medium; a process cannot exist without an enabling infrastructure - Mark Burgess

### Event-based architecture

## Conclusion