---
author: Olivier Wulveryck
date: 2016-06-23T15:32:54+02:00
description: Playing with websocket for a dynamic presentation.
draft: false
tags:
- websocket
- bootstrap
- golang
- revealjs 
- Javascript
- D3js
title: Websockets, Reveal.js, D3 and GO for a dynamic keynote
topics:
- topic 1
type: post
---

# the goal
As all my peers, I have the opportunity to talk about different technological aspects.
As all my peers, I'm asked to present a bunch of slides (powerpoint or keynote, or whatever).

In this post I won't dig into what's good or not to put in a presentation, and if that's what interest you, I 
recommend you to take a look at [Garr Reynold's tips and tricks](http://www.presentationzen.com/).

_Steve Jobs_ said:

> People who knows what they're talking about don't need PowerPoint

(actually it has been quoted in Walter Isaacson's biography see [this reference](http://blog.jgc.org/2011/11/people-who-know-what-theyre-talking.html)).

As an attendee I tend to agree; usually PowerPoints are boring and they hardly give any interest besides for the writer to say "hey look, I've worked for this presentation".

Indeed, they are a must. So for my next presentation I thought: 

wouldn't it be nice to use this wide display area to make the presentation more interactive.
One of the key point in communication is to federate people. So what if people could get represented for real in the presentation.

## how to: the architecture 

Obviously I cannot use conventional tools, such as PowerPoint, Keynote, Impress, google slides and so.
I need something that I can program; something that can interact with a server, and something that is not a console so I can get
fancy and eye-candy animations.

### The basic

[reveal.js](http://lab.hakim.se/reveal-js/) is an almost perfect candidate:

* it is a framework written in JavaScript therefore, I can easily ass code
* it's well designed
* it can be used alongside with any other JavaScript framework

### Graphs, animations, etc...

A good presentation has animations, graphs, diagrams, and stuffs that cannot be expressed simply with words.
I will interact with the audience. I will explain how later, but anyway they will send me some data.
I could process them in whatever server-side application (php, go-template base, python) but I have the feeling that's not 
the idiomatic way of doing modern web content. Actually, I would need anyway to deal with device (mobile, desktop), screen size,
browser... So what's best, I think, is to get the data on the client side and process it via Javascript.

[Data Driver Documents](https://d3js.org/) is the framework I will use to process and display the data I will get from the client.

It actually uses SVG to represent the graphs; I would have liked to use a full HTML5 to be more... 2016, but the D3 is actually very very good 
framework I wanted to use for a while.

### The attendees 

If I want the attendees to participate they need a device, to act as a client.
About all people I will talk to have a smartphone; that is what I will use. 

It has two advantages:

* it is their own device, I looks more realistic and unexpected: therefore I would get a better reception of the message I'm trying to pass.
* it usually has a Webkit based web browser with a decent Javascript engine.

I won't develop a native app, instead I will a webpage mobile compliant based on the [bootstrap](http://getbootstrap.com/) framework.

### The HUB

The point now, is how to make my clients and my presentation to exchange data.
As I said before, I would not be an easy task to setup a pure browser based peer-to-peer communication, so I will fall 
back to the traditional web server based hub.

the first idea is to use a RESTfull mechanism, but this has the major disadvantage of not being real-timed.
What I would like is a communication HUB that would broadcast events as soon as they are reveived.

I've implemented a server in go to do so. The clients will talk to the server over websockets which are now natively present in every
modern browsers.

#### the server

I've used the [Implementation from gorilla](https://github.com/gorilla/websocket) because it seemed to be the best as of today.
It implements all the RFC and the development is up-to-date.

The code heavily relies on channels to broadcast the messages between the different peers.
 I've taken the chat example present in the gorilla's package.

At first I did code all the mechanism is a simple go package. After a bunch of code, I've decided to split the code into two different
projects: The main presentation and the [gowmb](http://github.com/owulveryck/gowmb). The gowmb package is usable in others projects.

# Conclusion.

I don't go into the implementation details in this post, instead I will refer to the [github](https://github.com/owulveryck/topology-presentation)
repository where the presentation is hosted.

By now I have a good animated slideshow, and the ability to join the slides with a mobile phone.
I can also draw a topology of the attendees via D3 and interact with them.
