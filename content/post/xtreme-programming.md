---
date: 2017-04-14T23:00:59+02:00
description: "In this post I will describe my experience with extreme programming. I have tested it in conjunction with scrum, and I have been impressed by the results. I will try to explain why it is, according to me, a very good leverage of digital transformation for Ops team."
draft: false
images:
- /assets/images/Extreme_Programming.png
tags:
- agile
- xtreme programming
- ops
- dev
title: I have tried Extreme Programming within a sprint and I think it is an excellent agile method for the Ops!
---

# Part I: Agility

## 2003

I have discovered the notion of extreme programming more than 15 years ago. My job was to integrate *and* to develop pieces of code in Java for the IBM Websphere Business Integration server.
We were a small team with light programming skills. A part of our job was to operate the software, the other part was to develop. It was in 2003.

We were trying hard to stick to the specific framework we developed. 

Of course, in 2003, no French company I have been working for was talking about agility, The minimum viable product was not an option. The client wanted the full viable product delivered on time.

One of those morning where we were trying to find the motivation to do another step in the unknown, a colleague told us about a "new concept" he read about: Extreme Programming.

He explained that we should work in pairs. He told us that we should test every single feature even before actually implementing it, and so many other things... The ideas were good, but the milestones of the project were short. Obviously we were _too busy to innovate_ .

(By the time, as I grew up, I learned that _Good ideas don't always win!)_ 

{{< figure src="https://imgs.xkcd.com/comics/mobile_marketing.png" link="https://xkcd.com/1327/" caption="" >}}

<small>_Note:_ I intentionally put this XKCD as a reminder that a "good idea" is an abstract concept. Therefore, I insist on the fact that this blog reflects my own opinion. Even if I remain sure that it was a good idea, maybe it was not :)</small>

At the end of the project, when the [Deming wheel](https://en.wikipedia.org/wiki/PDCA) turned, I noticed that we were missing of agility.

## 2017

It is now 2017. Every IT crew thinks that agility is the way to work. 
Some of them have enough support from their management to actually implement an agile method.
Others may not be mature enough, but are pushed by the neverending decrease of the time to market to try new methods of work and delivery.

Agile and DevOps are, according the trends, the methods to use; it is seen as the holy grail.
But when it comes to agility, people are usually restricting it to an implementation of [scrum](https://en.wikipedia.org/wiki/Scrum_(software_development)) or [kanban](https://en.wikipedia.org/wiki/Kanban_(development)).

Others methods such as Extreme Programming remains rare. See for example this Google trends chart about agility in IT:

<script type="text/javascript" src="https://ssl.gstatic.com/trends_nrtr/981_RC01/embed_loader.js"></script> <script type="text/javascript"> trends.embed.renderExploreWidget("TIMESERIES", {"comparisonItem":[{"keyword":"/m/02t2n","geo":"","time":"2012-03-18 2017-03-18"},{"keyword":"/m/02zhbn","geo":"","time":"2012-03-18 2017-03-18"},{"keyword":"/m/0ck_p8","geo":"","time":"2012-03-18 2017-03-18"},{"keyword":"/m/01mwhw","geo":"","time":"2012-03-18 2017-03-18"}],"category":0,"property":""}, {"exploreQuery":"date=2012-03-18%202017-03-18&q=%2Fm%2F02t2n,%2Fm%2F02zhbn,%2Fm%2F0ck_p8,%2Fm%2F01mwhw","guestPath":"https://trends.google.com:443/trends/embed/"}); </script> 

On the other hand, [lean](https://en.wikipedia.org/wiki/Lean_software_development) concepts are usually referenced in every single agile documentation. But the echo of the lean principles is not strong enough. And still, IT crew usually refers to those methods as _only good for pure dev teams_ and _we are not devs_ (and trust me: if I had had a cent every time I faced this sentence, I would be rich). 

Ops should not be opposed to Devs. That is a fact, but why? 

Because Ops are also developers. Actually Ops are developing their business. 
In the era of "services everywhere" (XaaS), operational teams (as opposed to business team) must define contracts of services. Therefore, they must develop the services to fulfill the contracts.
They are a business team; even if their business is not related to the core business of the company they are working for.

Take a look at Amazon. AWS' core business is computer centric, but at the beginning it was only the "ops department" of a retailer.

# Part II: Using extreme programming in an "Ops" team

At the present time, I am sub-contractor for a retail company. My job is to give hints and hands to the operational teams. Their goal is to serve the business in a way efficient enough in order to follow the growth of the core business (and it is growing fast). 

## Context

Let me define the context of my job.
I am working in a team whose goal is to expose IaaS based on public cloud offers. Among other services, we want to provide to our customers a service of file transfer.
The transfer engine is an existing product. What we have decided to do is to add a RESTful API in front of the engine (this is a shortcut for clarity).

The team is composed of 4 people (mainly ops). But only one of them really knows the transfer engine. Therefore he has been designed as the legitimate implementer of the web service.

This person is my colleague [Alexandre](https://www.linkedin.com/in/alexandre-hisette-aa1076a/) from [Techsys](https://www.linkedin.com/company-beta/719121/). He is a certified system engineer. And best of all he did not tell me _I am not a dev_.

Regarding my job, I was assigned to another project that was also involving API management.

One last thing to know: the Team is also experimenting a Scrum method. We are "sprinting" for the releasing our products.

## How did it get extreme

Alex started to implement its API gateway. I managed to convince the team to use the go language (telling how is not the point of this article). We were exchanging about the implementation, the design, and the language.
At a certain point, for the past sprint, we started to work together by really sharing a screen.

When we were not sure about the design, we were instantly brainstorming around a coffee.

We decided to write the tests with a goal of 100% of code covered. 
When he was busy with something else, I wrote some tests for him, and his job, when I was by myself busy on something else, was to actually implement the code that was giving 100% of success.

Is that extreme programming? Let's recap.

## **What is Extreme Programming**

Even if this section is a vague copy/paste from [wikipedia](https://en.wikipedia.org/wiki/Extreme_programming) it is time to define some concepts of XP.
(I strongly encourage you to read the wikipedia article though)

Extreme programming is a software development methodology.

The activities of XP are:

* Coding
* Testing
* Listening
* Designing

Why is that extreme? Because all of the activities are taken to their extreme level. For example: regarding the tests, not only the business logic is tested. but every single component of the software is fully tested. (remember our goal of 100% code coverage? Yes we are extreme!)

Regarding the practices of extreme programming: there is 12 practices grouped in 4 areas (again [wikipedia](https://en.wikipedia.org/wiki/Extreme_programming_practices) is the place to go after this blog post).

* Fine-scale feedback
 * Pair programming
 * Planning game
 * Test-driven development
 * Whole team
* Continuous process
 * Continuous integration
 * Refactoring or design improvement
 * Small releases
* Shared understanding
 * Coding standards
 * Collective code ownership
 * Simple design
 * System metaphor
* Programmer welfare
 * Sustainable pace

## So are we extreme?

Yes! 

Because we are doing pair-programming. 

Because by using scrum, we used the planning poker. 

Because we are extremely testing our app.  

Because of the sprint releases, we are doing small releases

Because go impose the coding standards

Because we were proud of what we did, and this provided welfare

And probably many other things I cannot list in a single blog post.
# Conclusion

In the sprint review meeting, of course my own goals were not reached (The goals of my own project). I have passed too many times to work with my colleague to complete my own tasks.

My "product" was enhanced of 2% instead of 5%, but the other product, the one developed in pair, has increased by 30%. On the average, the quality of the whole service provided by the team made a greater gap.
Moreover we were both *proud* of what we accomplished. The code was clean, documented and tested. The product owner was very pleased of that.

* Learning was amplified
* team was empowered
* the service was more consistent as a whole (two heads are better than one)

Is that important that my product was not as advanced as it should have been? For the team, definitely not. My product is viable anyway, and the delay induced by the time "lost" is at maximum two weeks. And this delay will be filled if, by the next sprint, we decide to work together with my colleague on my project.

But that is a decision to take by the team and the product owners for the next sprint.

---
Credits:

* Image by DonWells (https://en.wikipedia.org/wiki/File:XP-feedback.gif) [CC BY-SA 3.0 (http://creativecommons.org/licenses/by-sa/3.0)], via Wikimedia Commons
