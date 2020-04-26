---
title: "Coupling in software: divide et impera"
date: 2020-04-26T15:47:57+02:00
draft: true
---

# Introduction

TODO: Expose why coupling is a problem:

- maintenance cost
- new features

# Coupling inside the code

## Connascence

## Modern connascence: IA application

### Decoupling software 1.0 and software 2.0

_Transition_: 

# Coupling inside the software

## Splitting the application by domain design

#### Seeing the problem on a map

### Separate the infrastructure from the application

> A story cannot be written down without a medium; a process cannot exist without an enabling infrastructure - Mark Burgess

### Event-based architecture

# Conclusion





The first article would be about coupling in software (I recently learned about the word connascence). 
Besides the general introduction, I see two parts:
- Decoupling the software 1.0 and the software 2.0 part (with onnx)
- Decoupling the different domains of an application with an event-based architecture (based on the software engineering principles as exposed by Titus Winter: add time and programmers to some code);

In a second article, I would introduce the sample app of sentiment analysis. I will describe the software 1.0/software 2.0 mechanism as I did in the presentation I gave you.
Then I will add the need for decoupling the application, as seen in the previous article. This will be a lever for the knative eventing part (the one deployed on k8s).

In conclusion, I am thinking of mentioning that this decoupling mechanism has excellent benefits, but it induces some extra work on the infrastructure level. It introduces a separation of concerns between the infrastructure and the application. As shown on a Wardley Map, the infrastructure should come from a product to a commodity. This will introduce the third article:

Third article: cloud-run eventing
The introduction would be a wrap up of the first two articles.
It will show that the knative-eventing product is becoming a commodity etc. then it quickly goes into the technical setup of the application.
Ideas....
divide et impera

Introduction
From a container of data to a data processor

Object storage and cloud computing have changed the paradigm of data storage: you don't unmarshal the data to make it fit in structured storage anymore. Nowadays, you use blob storage to store the document as close as possible to its original form.

The value is not gold by the document itself anymore. What's important now is the history of changes applied to the document.

the paradigm shift

For a long time in the history of IT, people have stored data, and building an application was, more or less, designing a view on top of the data storage.

> The next generation of enterprise data platform architecture requires a paradigm shift towards ubiquitous data with a distributed data mesh. - Zhamak Dehghani

 paradigm shift

Digital world
Container of data
Process


In this article, I explain why event base architecture is relevant in some of the recent use cases I encounter.

> A story cannot be written down without a medium; a process cannot exist without an enabling infrastructure - Mark Burgess

> The notion of a document change; a document is no longer a container with a set of sentences; a document is a process with a set of changes - Jeffrey Snower

The first article would be about coupling in software (I recently learned about the word connascence). 
Besides the general introduction, I see two parts:
- Decoupling the software 1.0 and the software 2.0 part (with onnx)
- Decoupling the different domains of an application with an event-based architecture (based on the software engineering principles as exposed by Titus Winter: add time and programmers to some code);

In a second article, I would introduce the sample app of sentiment analysis. I will describe the software 1.0/software 2.0 mechanism as I did in the presentation I gave you.
Then I will add the need for decoupling the application, as seen in the previous article. This will be a lever for the knative eventing part (the one deployed on k8s).

In conclusion, I am thinking of mentioning that this decoupling mechanism has excellent benefits, but it induces some extra work on the infrastructure level. It introduces a separation of concerns between the infrastructure and the application. As shown on a Wardley Map, the infrastructure should come from a product to a commodity. This will introduce the third article:

Third article: cloud-run eventing
The introduction would be a wrap up of the first two articles.
It will show that the knative-eventing product is becoming a commodity etc. then it quickly goes into the technical setup of the application.
Ideas....
divide et impera

Introduction
From a container of data to a data processor

Object storage and cloud computing have changed the paradigm of data storage: you don't unmarshal the data to make it fit in structured storage anymore. Nowadays, you use blob storage to store the document as close as possible to its original form.

The value is not gold by the document itself anymore. What's important now is the history of changes applied to the document.

the paradigm shift

For a long time in the history of IT, people have stored data, and building an application was, more or less, designing a view on top of the data storage.

> The next generation of enterprise data platform architecture requires a paradigm shift towards ubiquitous data with a distributed data mesh. - Zhamak Dehghani

 paradigm shift

Digital world
Container of data
Process


In this article, I explain why event base architecture is relevant in some of the recent use cases I encounter.


> The notion of a document change; a document is no longer a container with a set of sentences; a document is a process with a set of changes - Jeffrey Snower
