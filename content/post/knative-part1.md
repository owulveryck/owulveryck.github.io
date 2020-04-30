---
title: "Divide et impera - part I: coupling"
date: 2020-04-26T15:47:57+02:00
description: "This article is about coupling in IT; divide et impera - divide and conquer"
draft: true
---

In 2020, an old debate has risen from the ashes:

_Monolith or micro-services?_

In the past years, I have been working on IT products that were designed to work at scale.
The micro-services pattern was well suited for this need and was usually the paradigm of choice.

And application composed of micro-services is supposed to be scalable; the main pitfall is to 
design the service contracts and avoid the distributed monolith.
The secret is to manage the coupling between the components in the same way developers have learned to handle coupling
inside the code over the years.

By mastering the coupling, the design paradigms have shown good results on quality and reliability of softwares.

How do those paradigms apply to IA based applications?

In this series of articles, I plan to dig in the design of code and hosting applied to IA software.
This first article is about concepts.
Then I may write a second article technical to expose a concrete example and illustrate the concepts.
Eventually, a third article will show how to use managed services to transform products into commodities[^1].

[^1]: Simon Wardley's book, chapter 7 [Finding a new purpose](https://medium.com/wardleymaps/finding-a-new-purpose-8c60c9484d3b)

## What is _coupling_?

In physics, two objects are said to be coupled when they are interacting with each other (source [wikipedia](https://en.wikipedia.org/wiki/Coupling_(physics))).

Two elements are strongly coupled when a perturbation on the first element induces a disturbance on the second that impact on its behavior strongly.
A corollary is that two systems are weakly coupled if the perturbation induced on a component from the excitation of another one is low.

In chapter 12[^2] of his book [_In search of Certainty_](http://markburgess.org/certainty.html), [Mark Burgess](https://twitter.com/markburgess_osl) discusses the notion of coupling in physics. Some of the conclusions exposed in chapter 13 are [^3]:

- strong coupling leads to precision but the fragility;
- weakly coupled components can be rearranged, taken apart and made into something else, assembled into a meaningful composition;

[^2]: Chapter 12: Molecular and Material Infrastructure: _Elastic, plastic and brittle design_, page 383
[^3]: Chapter 13: Orchestration And Creative Instability: _Or why the conductor does not promise to blow every trumpet_, page 422

A software is symbolized by a composition of elements (objects, services, functions, etc.) that are interacting, hence
the paradigm of coupling applies to software engineering.

From my experience, the strength of coupling is roughly evaluated but not correctly measured. Nevertheless, it feeds numerous discussions about the quality and reliability of the application.

Actually, in modern software, **coupling** arises at different levels, **inside the code**, *and* at the **architecture** level.

This article is about coupling in application powered by nowadays IA techniques; This article does not show any code and focuses on two principles about the separation of concern. I will eventually develop them as proof of concept in futures articles.

## Coupling inside the code

In the book [Reliable software through composite design](https://archive.org/details/reliablesoftware00myer) from 1957, Glenford Myers exposes why designing a system to master coupling is essential. A strongly coupled system induces complexity and a high cost of maintenance. It makes the software "fragile" (as seen in the previous paragraph).

Here is an illustration of this complexity in the function of time (compiled from Myers's book and [Notes on the Synthesis of form](https://en.wikipedia.org/wiki/Notes_on_the_Synthesis_of_Form) from Christopher Alexander).

{{< blockquote author="Glenford Myers" source="Reliable software through composite design" >}} 
Consider a system of 100 lamps where a lamp could represent a statement or a segment in a program. Each lamp can be in one of two states, off or on.
The lamps can be connected together in any fashion. If a lamp is on, the probability of its going off in the next second is 50% if at least one adjoining lamp is on.
If a lamp is off, it will stay off as long as all of the lamps directly connecter to it are off.
{{< /blockquote >}}

| Test case                              | Time to reach equilibrium                   |
| -------------------------------------- | ------------------------------------------- |
| no connection between any lamp         | around 2s (2<sup>1</sup>s)[^4]              |
| all lamps are connected                | around 2.4e24 years (2<sup>100</sup>s) [^5] |
| 10 independent groups of fully connected lamps | around 17 minutes (2<sup>10</sup>s) |

[^4]: average time for a lamp to go off
[^5]: average time for all lamp to go off

This time-based illustration gives the feeling that coupling can be painful and have a substantial impact on the velocity of new development.
For dozens of years, developers have switched from paradigm to paradigm to fight coupling in the code.

Object-oriented programming is, for example, such a paradigm. This paradigm introduces the notion of [connascence](https://en.wikipedia.org/wiki/Connascence) as a metric to measure complexity.

Those concepts and paradigms helped the developers to structure their code, rising the need to package, and modularized it.

The goal is to reduce the connascence level, from the most dangerous (Connascence of identity) to the less dangerous (Connascence of name). See the [wikipedia page](https://en.wikipedia.org/wiki/Connascence) for the complete list and classification.

This introduced the concept of the lifecycle of modules and thus the need for versioning.

Let's now consider IA based application where the business logic is using compute libraries with a lifecycle that is independent of the development language.

### Coupling between software 1.0 and 2.0

An IA model is a mathematical representation associated with some values. Let's call it a software 2.0 (this term has been introduced in late 2017 by Andrej Karpathy in a [blog post](https://medium.com/@karpathy/software-2-0-a64152b37c35), and is slowly becoming a common language in the data world).
In this paradigm, software 1.0 is the result of regular code. Its usage is to glue the interfaces, I/O; it acts as a host for the IA/ML algorithm.

To create a software 2.0, there is a need to transcribe this mathematical model into something understandable by a computer, eg. a sequence of bytes. Then, we need two different types of software 1.0:

- one will act as a helper for the training phase
- the other one will handle the inference (the runtime of software 2.0)

Those softwares are linked by the algorithm. The algorithm is expressed in code and therefore become an entity shared by each element.

This is [connascence of identity](https://connascence.io/identity.html) and this is a major coupling that introduce a maintenance problem.

#### Decoupling software 1.0 and software 2.0

The learning phase usually involve gradient computation. There are multiple ways to automate the computation.
One of the method is to use symbolic differentiation. This gives a new mathematical function.
It is then possible to decouple the training software and the learning software; each of them would run its own mathematical equation; only the data would migrate from an implementation to the other.
This is a [connascence of data]()

This setup introduces coupling between the elements.

> time can run backward locally, as long as the process doesn't depend strongly on what happens around it.
https://youtu.be/lDFQiS9T_xk?t=2598



However, this introduces a coupling between the algorithm and the software itself.

The regular software acts as a host for the machine learning algorithm.

Let me give an example of such a coupling:

A machine learning software's lifecycle is composed of two phases: the training phase and the inference phase (to perform prediction). 

To different types of software 1.0 are required to process a machine learning algorithm. 
One for the training phase, and a second one for inference.
The two software have different goals; They are linked together via the algorithm and its implementation.

Refactoring the software 2.0 to enhance its efficiency requires to the, the actual value of the software 1.0 is in the prediction phase.
Therefore, coupling the software 2.0 with its host 

The software 1.0 and the software 2.0 and tidily linked by the language, framework, and modules they are built with.

### a Babel fish for deep learning 


_Transition_: TODO

### Splitting the application by domain design

TODO Explain the problem to run at scale, and the need to seperate concerns

#### Seeing the problem on a map

### Separate the infrastructure from the application

> A story cannot be written down without a medium; a process cannot exist without an enabling infrastructure - Mark Burgess

### Event-based architecture

## Conclusion

> A story cannot be written down without a medium; a process cannot exist without an enabling infrastructure - Mark Burgess
> The notion of a document change; a document is no longer a container with a set of sentences; a document is a process with a set of changes - Jeffrey Snower
