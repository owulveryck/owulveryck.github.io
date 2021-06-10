---
title: "20210608 Datamesh Cue"
date: 2021-06-08T15:47:20+02:00
draft: true
---

In 2021, a rich set of data is the soil that empowers the business of all the Internet Giants (GAFAM, NATU, ...).

Meanwhile, traditional companies are striving to remain competitive. The mandatory acceleration of their business goes through a massive digitalization of their operations and assets.

Amongst the most valuable assets of the companies are the data. Big data's promises are attractive. However, in the wild, the data unit is commonly separated from the core business. Even if much those data units provide much effort, the business plan usually look like:

- step 1: collect
- step2: ?
- step3: profit

{{< figure src="/assets/datamesh/gnome_data.png" width=100 >}}

> If we collect and manage the data, profit will come.

In this article, I will present a way to address **step 2** of the plan. I will take some concepts from the data-mesh paradigm, and, to avoid the [cargo cult](https://en.wikipedia.org/wiki/Cargo_cult), I will eventually illustrate them with a trivial technological implementation.

## Previously in the world of _Data_

As explained in the introduction, much effort has been in technological solutions to address big data issues and extract its value.

The adage "_during a gold rush, sell shovels_" lead to various technological implementations such as data-warehouse, data-lake, and lately data-factories. Even if it they may sounds ok, those solution share a common problem: they can hardly scale.

To manage this rush by addressing the _shovels_ ecosystem issues while focusing on the _gold_ (the data),  Zhamak Deghani introduced a paradigm shift called _data mesh_.
The data-mesh is a way to exploit the data in a distributed manner.
In essence, the paradigm shit is:

- focusing on the distribution of ownership and technological architecture;
- placing the data at the center of each distributed component.

All the rest of the data mesh is about solving the problems that come with that.

The datamesh is based on four pilars:

- a federated computational governance
- a domain-driven data ownership architecture
- thinking data as a produt
- relying on a self-serve infrastructure platform

### The pilars of the data mesh in a glimpse


![](/assets/datamesh/Picture-of-the-schematic-diagram-of-a-general-communication-system-Claude-Shannon-on.png)

## Why configuration

A configuration can be defined as a human-computer interface for modifying system behavior. (chapter 14  SRE book - Configuration Design and Best Practices)

https://sre.google/workbook/configuration-specifics/

## Why a configuration language ?

A configuration language is any kind of language allows making it easier to deal with large amount of data.
Exemple is json specification, protobuf, ...

> Configuration languages create files usually read and interpreted only once, during initialization
(Certified Software Development Professional (CSDP) )

> "Configuration drift occurs when a standardized group of IT resources, be they virtual servers, standard router configurations in VNF deployments, or any other deployment group that is built from a standard template, diverge in configuration over time. … The Infrastructure as Code methodology from DevOps is designed to combat Configuration Drift and other infrastructure management problems."
[Kemp Technologies on Configuration Drift](https://kemptechnologies.com/glossary/configuration-drift/)


In the datamesh we need a configuration language that can work "at scale". There is going to be a consequent amount of dataproduct and a huge amount of data.

We need a configuration language suitable to provide automation (federated *computational* governance).

Managing the configuration of thoses products at scale is going to be an issue.

## Generic pipeline description

`collect_data | marshal_data | emit_data | validate | marshal_event | send_to channel`

`channel:=$(find channel for data)`

`read_from channel | unmarshal_event | unmarshal_data | profit`

### `emit_data`

### `validate`

![garbage in, garbage out!](/assets/datamesh/garbage_in_out.png)

## Expcted properties of the language

Some expected problems are:

- the lack of validation (weak typing for example)
- Inheritance (TODO)

Fundamental properties:

- Associative ( A & B ) & C = A & ( B & C )
- Commutative A & B == B & A
- Idempotent A & A == A 

> definig Schema is cumbersome, but keeping them sync is much worse

## Summary of the data product

![a set of data products](/assets/datamesh/set_data_products.png)

--- 

![a set of data products](/assets/datamesh/data_mesh.png)

---

![Streming Plateform](/assets/datamesh/data_streaming.png)


- Discoverable
Declared on a catalog and a search engine

- Understandable
provide semantic (meaning), syntactic (topology) and usage description (behavior)

- Addressable
must participate in a global eco-system with a unique address that helps its users to programatically find and access it

- Secure
Be able to be accessed security with global policies (role-base-access, GDPR, Info security, data sovereignty …)

- Interoperable
Be able to reuse, correlate and stitch them together across namespaces for new use-cases

- Trustworthy – Truthful 
Provided data provenance and lineage and data quality from the owner

- Natively Accessible 
Provided multimodal access like Web services, event of file interfaces

- Valuable on its own 
Designed to higher insights when combined and correlatedCommitted to SLAs

- Committed to SLOs
Must respect expected service levels, in terms of data availability

## Output data port

the event problem:

The envelop's role is to define the rules that allows the recipient to see if the message is of any interest fot it, in a glimpse.
In this sense, it is composed of metadata.

The role of the federated computational governance is to standardized the semantic of the metadata.
It participates in the *addressable* and *interoperable* needs of the data products.

