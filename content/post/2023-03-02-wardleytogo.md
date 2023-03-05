---
title: "Rationale behind watdleyToGo"
date: 2023-03-02T21:23:17+02:00
lastmod: 2023-03-02T21:23:17+02:00
draft: false
keywords: []
description: "wardleyToGo is a library and a language to describe Wardley Maps as code/data. This is the rationale behind this library"
tags: []
categories: []
author: "Olivier Wulveryck"

# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: true
toc: true
autoCollapseToc: false
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: false
---

Wardley Maps are a way of expressing business, market, or any other system through sketching. The Map is a sketch that offers situational awareness on a certain topic.

As a sketch, the obvious way to draw a map is with paper and pen. 
While it is an excellent starting point, a paper representation of a map has a problem: it is static.  
What I mean by this is that adjusting the placement of some components in the design phase can be tedious (even with a good rubber pencil).

Representing a map digitally presents a couple of advantages: they are easy to share and easy to **maintain** and **exploit**.

- **Maintaining** a map, is, in my context, the possibility to make some adjustments without changing its meaning (think of tweaks that can occur after a discussion with some peers).
- **Exploiting** a map is about understanding the landscape to orient future decisions.

_Note_: An important point to keep in mind when mapping is: a map shall not be used to illustrate the story you want to tell, the story should be extracted from the map.

In the orientation phase, being able to withdraw some elements from the map, for example, to focus on certain a path is useful.

Of course, we could record the creation history and do a couple of undo/redo to have intermediate representation, but this would imply that the creation steps were reflecting the path.
This is orthogonal with the idea of seeking strategy through observation **and then** orientation.

The best way to have partial representation is, therefore, to have an intermediate representation of the map where we can easily comment blocks to hide elements.

## Maps as code

In the design and maintenance phases, expressing a map as code present some flexibility. 
The code is the source of the map and tooling and practices to manage source code are widespread (managing source code is a commodity). 
Tools like `git` provide capabilities to:
- version and tag the map;
- collaborate on the map (even asynchronously or remotely)
- natively store the history of the map.

[onlinewardleymaps](https://onlinewardleymaps.com/) is an online tool made by [@damonsk](https://twitter.com/damonsk).
It is mature and widely used. To me, this tool has largely contributed to push maps-as-code into the _stage_ II of the evolution (as described by Simon Wardley).

onlinewardleymaps definition language is called `owm`.

## Taking my map "to go"

A first, the only version of the tool was an online version. This is good for the user experience as it does not require any setup to start mapping.
The problem is the coupling between the tool to render the map and the language. You cannot use `owm` outside of the browser.

### Versioning the Map

The problem is that managing the source code was a bit tedious as it required a lot of copy/paste from and to the tools.

The Visual Studio extension, that appeared later, made is a bit simpler to manager the source code. In this sense, it heavily facilitates the versioning of the map.

But there is no way to easily export the map to store it with a peculiar version of the code as you need to export the map manually.

It makes it a bit difficult to take the map "to go" (or to take away depending on which flavor of English you practice (sic)).

### Basic need: CI/CD 

For years now, continuous integration (CI) and continuous delivery (CD) has proven some benefits in the release cycle of a digital asset. 
The idea is that every revision that is committed triggers an automated build and test. With continuous delivery, code changes are automatically built, tested, and prepared for a release to production

As I was able to commit the source code of my maps, I wanted to use a CI/CD mechanism to compile my source code and render my maps. 

onlinewardleymaps does not provide any SDK allowing me to use the rendering engine in a headless building process.
As a geek, I could have raised an issue an start contributing to the project.
But I wanted to build my own to deeply understand how to build a map from the inside and because I found it fun to have one more side project.

## Building an SDK to draw map "as code"

I am a big fan of the Go language for various reasons. In a glimpse the reasons I am using Go for my SDK are:

- I find it fun to code in Go;
- I master the language enough to speed the development and to focus on the design of what I want to achieve;
- the language is suitable to build "Ops tools".

Naming things is hard, I named my SDK warldeyToGo (take your map "to go" with Go).

The design of the SDK is:

- a central package that acts as an intermediate representation of what a map is. (in a glimpse, it is a directed graph in which the nodes are components that are able to give their position in an euclidean canvas )
- a set of componnets that implements the intermediate representation and that are able to represent themselves in SVG (for example, I have Wardley Component and Team Topologies components)
- parsers to high level languages that transpile the representation into the intermediate representation.

As a kickstart, I implemented a parser to the `owm` language.

So I can build my map with `onlinewardleymaps`, get the source code in `owm` and build a tool with the SDK to transpile `owm` into the intermediate representation, and render it into SVG.

## A high level language to express map "as data"

So the SDK is truly allowing to code Wardley Maps. The `owm` syntax is therefore simply a user interface between the need to express a map and the representation.

> Just as user interfaces are the conduit between humans and the functionality of products, products are the conduit between customers and the team of people solving a particular problem.

onlinewardleymaps is a solution to the need to represent a map digitally.
wardleyToGo is a solution to code a map.

But what I was missing is a solution that allows a computer to assist me in the design of a map.

### The problem with the Euclidean representation

### Wardley Map: toward system 1

### Thinking about value chain and components

### A computer language to shape the way we think
