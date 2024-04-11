---
title: "Data Evolution"
date: 2024-04-09T12:15:33+01:00
lastmod: 2024-04-09T12:15:33+01:00
images: [/assets/data_certitude.png]
draft: false
keywords: []
summary: The evolution applied to the data
  
  How I changed the data-mesh approach
tags: []
categories: []
author: ""

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S.
comment: false
toc: true
autoCollapseToc: false
# You can also define another contentCopyright.
contentCopyright: false
reward: false
mathjax: false

# Uncomment to add to the homepage's dropdown menu; weight = order of article
# menu:
#   main:
#     parent: "docs"
#     weight: 1
---

## Context
I have been an early advocate for the data mesh paradigm since Zhamak Dehghani first proposed it.
As an early supporter, I identified the potential of this novel approach to data organization.
Four years have passed, and the data mesh paradigm has indeed earned widespread acceptance.
However, I have not yet seen a concise, practical data mesh transformation plan within any organization.

When I say "advocate," I mean to say that I've been highlighting the benefits of this paradigm, which are rooted in its four main pillars:

- The orientation of design towards **domains**
- Applying **product thinking to data** (data-as-a-product)
- Federated computational **governance**
- Developing self-service **data platforms**

Then comes the question, **where should one actually start to implement the data mesh?**

In principle, any organization can kick off their journey towards data mesh by giving focus to these four pillars.

**Beginning with** a design that is **domain-oriented** builds the groundwork for a thorough comprehension of the data mesh.
This means not only setting the data mesh as an objective but also **ensuring that the decomposition** of the **domain syncs with** the **existing structure** of the organization.
However, this is a **profoundly conceptual approach** that might not yield immediate results, and moreover, it lacks the agility that is so beneficial.

Both **federated computational governance** and the **self-service data platform** are simply **enablers** of the data mesh.
They share a common objective: to simplify the development of data-as-a-product and the **creation of interconnections**, essentially supporting the mesh.
One can try implementing them as a foundation, but to mesh what ?

What remains then is to tackle data-as-a-product, a cornerstone of the data mesh that I have previously discussed.

Interestingly, several organizations claim to have implemented the **data mesh "by accident,"** perceiving this paradigm as the natural evolution of data management.

In this article, I attempt to apply a well-recognized model of evolutionary progression to understand data evolution.
**The objective** is to aid in **visualizing data maturity** and assist companies in **identifying** their **tipping point**, 
i.e., when they will start seeing **significant benefits** from **implementing data contracts** and **treating data as a product**.

## Modeling the evolution

I will first explain the model I will use. 
This model is known as the evolution model by Simon Wardley and is successfully implemented in Wardley Maps. 
My goal here is not to describe a specific company's landscape, so I won't need a full map. 
Instead, I will use the evolution model and try to apply its general purpose function to the data.

**Disclaimer: Regarding the model:**
The theory of evolution is well-suited for application in a competitive environment where everything evolves based on supply and demand. I am considering businesses that are subject to those constraints of competition and, consequently, their data will also follow those rules. Therefore, the model will apply.

### TL;DR: what is this model

Simon Warldey needed a way to represent the evolution of the components on his map.
He could not rely on a basic timescale because it would prevent him from comparing heterogeneous elements and would break the consistency of movement.
For example, on a time scale, the distance between the genesis and the maturity of a car (around 100 years) would have been much greater than the distance between the same points for a smartphone (around 10 years).

The evolution was function of its ubiquity and its certainty.
In a market economy, the ubiquity is lead by _demand_. More demand induces more presence. It is a declinaison of the theory of the diffusion of innovation.
The certainty comes from the Stacey model. 

![A S curve representing the evolution of the data, the X is the certainty and Y is the ubiquity. There is a division: the botton of the S is raw data, the middle is curated and the top is authoritative. The inflexion point is indicated as data-contract](/assets/data_certitude.svg)
