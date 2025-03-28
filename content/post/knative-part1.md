---
title: "Divide et impera - part I: coupling"
date: 2020-04-26T15:47:57+02:00
summary: "This article is about coupling in IT; divide et impera - divide and conquer"
draft: false
tags: ["architecture", "coupling", "eventing", "AI"]
---

In 2020, an old debate has risen from the ashes:

_Monolith or micro-services?_

In the past years, I have been working on IT products designed to work at scale.
The micro-services pattern suits this need and is, therefore, the paradigm chosen _de facto_.

An application composed of micro-services is supposed to be scalable; the main pitfall is to think and design accurate service contracts to avoid the antipattern of a distributed monolith.
The secret is to manage the coupling between the components in the same way developers have learned to handle coupling inside the code over the years.

The design paradigms have shown excellent results on the quality and reliability of software by mastering the coupling.

How do those paradigms apply to AI-based applications?

In this series of articles, I plan to dig into the design of code and hosting applied to AI software.
This first article is about concepts.
Then I may write a second article, more technical, to expose a concrete example and illustrate the concepts.
Eventually, a third article will show how to use managed services to transform products into commodities[^1].

[^1]: Chapter 7 of Simon Wardley's book ([Finding a new purpose](https://medium.com/wardleymaps/finding-a-new-purpose-8c60c9484d3b)) gives more explanation of the evolution of a product to a commodity.

The article' structured is:

- Understanding the notion of coupling applied to IT;
- Understanding coupling inside the code of an application and the importance of managing it;
- Understanding coupling inside a system (a composition of software) to see the strength and weakness of the event architecture pattern;

Hopefully, at the end of the article, I have all the tools mandatory tools to answer those questions: is Caesar right: _divide and rule_? Monolith or microservices? (spoiler it depends)

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

{{< blockquote author="Lord Kelvin" >}}
When you can measure what you are speaking about, and express it in numbers, you know something about it, when you cannot express it in numbers, your knowledge is of a meager and unsatisfactory kind; it may be the beginning of knowledge, but you have scarely, in your thoughts advanced to the stage of science.
{{< /blockquote >}}

Reducing the connascence level, from the most dangerous (Connascence of identity) to the less risky (Connascence of name), is a way to manage flexibility within the development.

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

The algorithm act as a link between the components. On top of that, the algorithm itself is transpiled into code, usually using the same language as its host. This makes software 2.0 an entity of software 1.0; therefore, this is comparable to the [connascence of identity](https://connascence.io/identity.html).

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
- creating another software 1.0 that can read the formulae and feed it with real-time data for inference;

_Note:_ [Open Neural Network eXchange (ONNX)](https://onnx.ai) is an example of such a DSL, but digging into the technical implementation is out of the scope of this article and will eventually come in a future post.

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

According to this definition, a system can run "at scale" if it's processing power is compatible with the expected usage.

Luckily, the decoupling we made before made it simpler to adapt the processing power for both the training and the inference process.

On both parts, data usually need to be pre-processed before they can feed the software 2.0. This is where the data-pipeline comes in.

### Data pipeline and time reference

A data pipeline is a series of data processing steps.

To illustrate, let's use [C4](https://c4model.com), a method to design an application (and to take care of coupling).
This method defines the concepts of _containers_ and _components_.

- A container (_not docker!_) is something that needs to be running for the overall software system to work.
- A component is _a grouping of related functionality encapsulated behind a well-defined interface_

Building a data pipeline brings the ability to split the components into the containers.

Each container defines its own scale. As a consequence, the overall system can adapt smoothly, and be compatible with the requests made against it.

The different _containers_ communicate through channels of communication where they can exchange messages.

Claude Shannon explains in detail the notion of communication channel in his [mathematical theory of communication](https://en.wikipedia.org/wiki/Information_theory#Channel_capacity). For our explanation, let's classify the communication channels: _rigid_ and _robust_.

- A rigid channel is something that couple the components strongly, therefore information is received synchronously from point A to point B (such as an in-memory semaphore, for example).
- A smooth channel is an element that carries information that will eventually flow from point A to point B.

The nature of the communication channel is the next axis of work. As seen in the first paragraph, a rigid channel will lead to precision but fragility, and a smooth channel is more robust but less precise.

For example, with a smooth channel, even _time can run backward locally[^8]_, which is very useful for elements of the pipeline that deals with databases.

[^8]: _time can run backward locally, as long as the process doesn't depend strongly on what happens around it._ Mark Burgess - [Bigger, Faster, Smarter - Episode 1](https://youtu.be/lDFQiS9T_xk?t=2598)

### Event-driven architecture

The event-driven architecture is a set of models that allows weakly coupling in a process-pipeline.

For example, A process emits an event in a communication channel, and this event is received by any process listening on the same channel.

An event is fired when a state changes. This architecture is suitable for an AI system.
In a data pipeline, every container processes the data; touching the data is modifying its state.

Building a data pipeline in an event driven-architecture is a way to implement a scalable, decoupled architecture. As seen before, this architecture is more robust. The main drawback is the balance that needs to be done between scalability and real-time processing.

Another concern about this design is the ordering of the events.
We saw in the last paragraph that time could run backward locally as long as the process doesn't depend strongly on what happens around it.

Giving a warranty on event order is, once again, solidifying the link between components. As a consequence, deploying an event-driven architecture is, like every design pattern, be chosen wisely with regard to the business requirements.

For more information about ordering (in channels), cf this article [Ordering messages](https://cloud.google.com/pubsub/docs/ordering) on the Google Cloud Platform's documentation.

#### Separate the infrastructure from the application

> A story cannot be written down without a medium; a process cannot exist without an enabling infrastructure - Mark Burgess

On a value chain, infrastructure stands far away from the business value. It is important to couple weakly the containers (in the C4 definition) from the infrastructure.
It is fairly easy nowadays to decouple the execution of a process from the computer thanks to containerization: Docker is the _lingua franca_. It is a sort of _connascence of type_.

More pitfall arises when we are dealing with messaging products. First of all, there is no standard in the message/event definition yet. On top of that, the need for ordering usually imposes to define a common link between the messages and the underlying infrastructure (e.g., shard keys).

Once again, thinking on how coupled the application and the infrastructure are is a decent way to manage the scalability while remaining flexible and agile in the developments.

## Conclusion

AI, deep learning, and other machine learning mechanisms are continually improving. POCs, papers, startup, and academic works are opening new perspectives to business and research to the digital world.

But the diffusion of innovation is conditioned by the ability to develop and run qualitative software that can run reliably.
We saw in this article that mastering coupling is one of the keys to shit from application development to software engineering.

**_So, Divide and Conquer?_**

Dividing the application brings flexibility, but managing the lifecycle of each part **and** the communication between the pieces is the next challenge.
This is why, in my opinion, people switch from the monolith to the microservice paradigm and vice-versa. _In fine_, it's all about the evaluation of the strength, the total cost of ownership (TCO), and risk analysis.

The popularity of the microservice architecture is at the origin of a lot of methods and tooling (service-mesh, for example). A lot of products exist nowadays to facilitate their applications. Some of them are even becoming commodities, thanks to cloud providers.

In upcoming articles, I will describe an example of AI application development on top of products such as Kubernetes. Then I will end this series by showing how a managed version of some service help in managing the TCO of the software.

Stay tuned!