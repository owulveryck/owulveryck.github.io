---
title: "Data-Mesh: Data Across Domain Borders"
date: 2023-07-25T10:34:36+02:00
lastmod: 2023-07-25T10:34:36+02:00
draft: true
keywords: []
description: ""
tags: []
categories: []
author: ""

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: false
toc: false
autoCollapseToc: false
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: false

# Uncomment to add to the homepage's dropdown menu; weight = order of article
# menu:
#   main:
#     parent: "docs"
#     weight: 1
---


This article presents the impact of a Domain Driven Design (DDD) approach on data (the second pillar of the Data Mesh paradigm).

We will see how this model helps frame the data and how it brings out the need for standardization and governance to bring value at the enterprise scale.
## Introduction

Encapsulation, abstraction, C4... numerous are the **design patterns**  that propose to **simplify**  **systemic complexity**  by separating internal implementation details from external interfaces.

These boundaries are permeable, and the means to traverse them are **standardized**  to offer access to the greatest number.

The *digital revolution* has made the use of closed computer systems almost obsolete.

Nowadays, systems must interact with each other to bring value (and the more interactions, the more value).

In **computing** , these interactions generally happen **via**  a set of **services** , which expose **useful functions** .

These services are sometimes **offered services** , when looked at from within the box (**white box** ), and sometimes **rendered services**  when viewed from outside the box (**black box** ).

When all these services form a product, the notion of ***affordance***  encompasses all the actions that the product **allows**  to perform (in English, *to afford* means to permit).

A product is affordant if:
- the actions it allows to perform are visible
- the functions not intended to be used are hidden

Data does not escape this model. We speak of **data-as-a-product** , a principle that applies product thinking to data to increase its usability and utilization.

With this model, **data**  offers a **view**  of a situation (**rendered service** ), and the **product**  makes this view **usable**  (**offered service** ).

The boundary between the rendered service and the offered service is materialized by a contract (for example, an API).

This contract allows defining a level of trust in the data, making it usable beyond the domain in which it is defined.

In this article, we are interested in valuing data across the conceptual boundaries of DDD.

Starting from principles, we will see the impact of DDD on the design of data-as-a-product.

Then, we will identify the needs for **standardization** , which allow data-as-a-product to **bring value at the domain scale** ; we will then talk about ***data contracts.*** 

Finally, we will introduce the topic of data governance as a means to facilitate inter-domain communication by enabling reliable use of data across domains, at scale.
## The Domain and DDD at a Glance

Let's begin by introducing some basic concepts of DDD that are relevant to understanding this article.
### Ubiquitous Language: Domains and Bounded Contexts

Domain-Driven Design (DDD) is a software design methodology centered around the notion of the domain. This practice is well-known as a good development practice.

DDD is now widely recognized as a good software development practice.

![](https://lh5.googleusercontent.com/ApVn1qJPAoI4BUOFO5abDz56kTziFq3fLdvt838z6BghBnf0tAT05c1BeKUP1ie0xCSfn_f4fJ6ILThyQiVm9GK6SfD1v3YFYMxvmt4kqPOrSlf-yW6Dga0pFsjuqCNZtf1gNp0bSW5ihxLWRpP9510) 


In DDD: 
- **A domain**  is an aggregate of **knowledge**  or **activities**  **specific to a business: it represents/models a business area.** 
- A domain encompasses the **problems** , **concepts** , **rules** , and **relationships**  that define the business area to be modeled.

Thus, DDD enables understanding and addressing complex business problems by breaking them down into smaller, coherent units.

The **domain**  allows us to identify the **problem space**  to be addressed. The **solution**  to the problem space is then **framed**  by another DDD concept: the ***bounded context.*** 

A ***bounded context***  is a conceptual boundary in which specific business knowledge and activities are modeled to make them applicable and coherent within a software solution.
### Illustration

Let's illustrate these concepts with an example from the world of online commerce: 
- ***Domain*** : **Order Management** .

In this domain, we focus on the **problems** , **concepts** , rules, and relationships **specific to order management** .

For example, we define the notions of **customers** , **products** , **stocks** , **payments** , **shipments** , etc.
*A customer seeks to purchase a product in stock, for that, they need to pay for it, and we must ship the product to them.* 
- ***Bounded Context*** : **Order Processing** .

In this bounded context, we focus on the specific aspects of order processing. How will we model an order to manage it in a digital solution?

A good practice is to work with event streams that represent a sequence of actions.

For example:

checking the availability of products in stock -> applying discounts -> validating payment information -> generating invoices -> organizing shipment.

In our example, the *bounded context* encompasses all the knowledge and activities related to order processing.

This provides overall coherence to order processing (all the implemented solutions will share a common functional language, for example, they will have the same definition of what a customer, a product, a payment, a shipment, etc., are).

We notice that this **_bounded context is not watertight** .

Indeed, our intuition suggests that, to bring full value to the business, the elements of this bounded context must interact with other *bounded contexts* (in our example, **stock management**  or **customer management** ).

In the DDD model, interactions between bounded contexts are **normalized**  and represented in a *context map*.

This **normalization**  adapts the **behavior of digital solutions**  to **allow them to cross domain borders** .

Some of these digital solutions are fed by data (for example, solutions based on AI or ML), which also has a place in the solution space.

Now we will see what **representation**  data can have in this **space** , as well as how to **normalize**  it to **cross domain borders**  and bring value at scale.

## Domain Ownership: A Commitment of Responsibility

In DDD, ***domain ownership***  encompasses the notions of **responsibility**  and **authority**  that certain **teams**  or **individuals**  have over the **elements**  of a domain. This **ownership**  helps ensure **consistency**  and **transparency** .

As a result, it enables more effective **communication**  and **collaboration**  among different **stakeholders**  within the domain by **standardizing**  the **vocabulary**  and increasing the level of **understanding** .

This, in turn, leads to a **natural elevation**  in the level of trust in the proposed solutions.

Now let's explore the impact of this notion of ownership on data.

But before that, let's align on the definition of data.
### The Omnipresence of Data: Understanding the Problem and Supporting a Solution

When we talk about the omnipresence of data in information systems, we refer to its presence in both operational systems **as well as**  analytical systems. Data becomes "analytical" when it has a temporal dimension.

In this section, we will see that the type of data also varies depending on whether we are in the problem space or the solution space.
#### Data in the Problem Space

In the problem space, the meaning of data must be aligned with business concepts. The ubiquitous language will be applied to various notions that constitute the semantics of the data.

It is essential to note that semantics are associated with a particular domain.

![Data in the Problem Space](https://lh6.googleusercontent.com/IimlmYXP2b9DSZUv1uP9p36KiCsIbYQGT0Mok8OSeCl7QeLCprxp7fNh9PPH1GXBqQa6xJH88CDP4sODtU8K7KgIEAqS7OmvmeKmPBtGlzqWNMtyvtvhfcqrf0f4mqo7NgaprGhOMFs9ZoafZVzHTEk) 

### Data in the Solution Space

Once the problem is defined, data can help provide solutions in two different ways: 
1. **Operationally** : Data will be manipulated as facts processed by business rules to operate the business. 
2. **Analytically** : Data will be analyzed over time to understand business behavior.
#### Fact-Based Data Serving Operational Systems

Factual data is a key element of operational solutions. Digital solutions to business problems are materialized through the **application**  of business rules on a set of data.

Users interact with the results through a service, which is powered by a dataset provided by the user or collected based on the context.

**The product**  is **at the center**  of a **bounded context** . The data that supports it revolves around it in the **same bounded context** , under the **responsibility**  of the **team**  designing and maintaining the **product** .

Introducing data into this space is therefore natural.

For example, an event storming workshop that describes a process can be complemented with data representing the state of each event. Data will easily be surrounded by a bounded context.

Here's an example with **order management:** 

![Data in Operational Solutions](https://lh3.googleusercontent.com/Oi_5xd-PL4Poct1asdu5NGBlWnr315I3fqRmLoQ-t-KxBGP14BISk5RQ16znUu5oWH0ZpBTer5gIiPv9SwpbEjySIghEhsI9fP5-nGXPnHXkuA01izd2RPFtElvTspkhgjWA-sY6PgcjMC9bOIk3pcY) 

#### Analytical Data

Data Mesh is a data management paradigm that borrows operational management methods that have been successful in digital transformations.

Thus, Data Mesh proposes to apply DDD to the world of analytical data (the second pillar of Data Mesh).

As of the time this article is written, the application of DDD to analytical systems is a new practice. There is much literature describing its wonders, but there is little literature describing how to implement it.

As a result, there is no specific framework or working method that allows bounded contexts and data boundaries to emerge.

To propose the appropriate method, we will use a new model borrowed from Data Mesh: we will type the domains.

A domain can be:
- Source aligned: meaning the problems it seeks to solve will produce analytical data.
- Consumer aligned: meaning the problems it seeks to solve will consume analytical data.
- Shared: A domain whose problems consist of processing analytical data as input to provide analytical data as output.

Regardless of the domain type, note this:

"Data mesh adopts the boundary of bounded contexts to individual data products—data, its models, and its ownership." (Data Mesh - [Chapter 2 Principle of Domain Ownership](https://learning.oreilly.com/library/view/data-mesh/9781492092384/ch02.html#applying_dddapostrophes_strategic_desig)  - Zhamak Dehghani - 2022)

This means we will have data-as-a-product per bounded context.
##### Source Aligned

In the case of producing analytical data, the domains' and bounded contexts' segmentation will align with its operational counterpart's definition.

So, we'll have a gathering of a service that provides data, operational data set, and a set of analytical data sharing a common functional language.
##### Consumer Aligned and Shared

To introduce data into the solution space, a collaborative practice is to create a value chain from use cases.

For example, we can take a business domain, list the use cases that provide solutions to problems. Then we can add a level by listing the analytical data that powers these use cases.

This often involves aggregated data. Then, we add a level of dependency to bring out the source analytical data.

By analyzing the vocabulary in a workshop, we will reveal bounded contexts and, therefore, data-as-a-product.

![Data in Solution Space](https://lh4.googleusercontent.com/3GXnOzMz1E1dvd4g8O9vJoaj5duY2JhbYd35DolDDrO3Ej-H3BqLskwUu1BRl4ubg9UlgrwxmcNMXy3vDlBLQLm0rTflfBEyHdO0d9-4oK2uwGWnLOYYKrOgbSv014XkEQkb3N2Poyx-WGleCrHvXck) 


Example of a value chain for a stock management system (this is a fictional illustration).

*Note:* The horizontal positioning is arbitrary and meaningless at this point. 
- *Stockout rates* and *stock turnover rates* are consumer-aligned bounded contexts. The vocabulary of **stock**  will be consistent with business concepts. The notions of **turnover**  and **stockout**  will be documented via data semantics.
- Current stock data and delivery time data are operational data that feeds the calculations.
- Historical sales data and demographic data are analytical data that align with the sources.
- Market trends and seasonality data are shared data.

## Data Evolution

TODO: Clarify things

Reminder: Cite Wardley's evolution theory.

WIREFRAME:

Understanding data evolution requires careful consideration of several factors. First and foremost, it is essential to emphasize that domains are often in phase 4, primarily because the associated domain knowledge is mature. However, specific problems within these domains may be in phase 3, requiring a customized solution.

Domains are in phase 4 because the associated domain knowledge is mature.

The domain has a problem that is in phase 3 (this is an example).

To address the problem, we will design a custom solution (also an example).

This custom solution will rely on specific data. These specific data are known to the business, so they are in phase 4.

![](https://lh6.googleusercontent.com/ikGQfk-yU3dEOarBY-JamKecd9s4gUwXI44ZXaQkmQzSBvmWHZFjlJsjk5h9A5X_VOIJ-3qJBUcMwVC4x65bqTuAJd3ezRJoKa_K8sdKfuQt3OCO-Dq8Q42s7fpeyRdYqM7F-eL4LcmZmZErnKojX-w) 


To solve the problem, analysts will first work on extracting value from the data to make the solution work. In terms of evolution, the data is not very exploitable: the goal is to "make the solution work with the data" (make it work).

To effectively use the data in the solution space, it must be stored in a repository (the lake).

This custom solution will generally rely on specific data, which are already known to the business, placing them also in phase 4. To solve the problem, analysts will first work on extracting value from the data to make the solution work. In terms of evolution, the data is not very exploitable at the beginning: the goal is to "make the solution work with the data."

![](https://lh4.googleusercontent.com/dMP5H7AdKL1Z7GMGb9lUCmouR3JV2OHiJuzJ5kXiqLAfcH6YtTPlL2NRoiFblEabw7SQND6DvvnQuVnmSLy5MbmtzXcy_2CYqYmxDKds0-4eKvXTrj3z6kTJm92Q0AdBe3X4YN1zBkKlDcHsKhXUr1o) 


To effectively use the data in the solution space, it must be stored in a repository, commonly known as a "data lake." In the context of Solution A, applied to the domain, the evolution of the offering will push towards standardization and "productization" of the data. This is the birth of data-as-a-product, which includes aspects such as data contracts.

In the context of Solution A, applied to the domain, the evolution of the offering will push towards standardization and "productization" of the data.

This is the birth of data-as-a-product (data contract, and so on...).

![](https://lh6.googleusercontent.com/02oQW9tuTQDfqW5jnbvNqquy1MKl6C0rs21DXClHuo-adNFJsYPzYb0Ql0C2GHNMwkURtEsqRzzy6WxCfbwrRQ5WWIEQ-sPXV2ZgXMbZgdavf-1V74FRl-oeuCRYUZ4tv6tumMUt77pKOhwyCNXFOmw) 


This is when the product needs to evolve "according to demand." The demand for new solutions will require standardizing (commoditizing) the data-as-a-product that comes out of the solution space in which it was initially designed. This will lead to the need for platforms to make the product accessible to a broader audience. Data then becomes authoritative within the company.

This is when the product needs to evolve "according to demand."

Reminder:

Climatic Pattern (influence on components that are independent of what the team is doing): everything evolves according to supply and demand.

The demand for new solutions will require standardizing (commoditizing) the data-as-a-product that comes out of the solution space in which it was initially designed.

The need for platforms to make the product accessible to a broader audience will emerge.

Data then becomes authoritative within the company.

![](https://lh6.googleusercontent.com/iXg9DWI47FMabz-yGYORv7hV8GuY0wfcfnCPq1KLlMl-i3whY3JKKL-qXMtSyiBwgZhM5xLE6c-mGY4nUKNhBUrXJ6b4kT4Kd791l0ZJsVdPA_n2q3txeIFldd0EVRLIFuhYflp85Von5zQ8y2_v75E) 

## Standardization and the Emergence of Data Contracts
### The Importance of Standardization for Creating Value at the Domain Scale

TODO

Data standardization is a crucial step in creating value at the domain scale. It facilitates data exchange, exploitation, and interoperability, promoting transparency and consistency throughout the domain. Additionally, standardization harmonizes data formats, making information sharing and collaboration among different stakeholders easier.
### Definition and Role of Data Contracts in the Data Mesh

TODO

Data contracts play a crucial role in the context of the data mesh. They define the attributes, relationships, constraints, and expected behavior of the data. They act as an agreement between data producers and consumers, ensuring that the produced data is reliable, relevant, and aligned with the consumer's needs.
## Data Governance and Inter-Domain Communication
### Introduction to Data Governance and Its Role in Facilitating Inter-Domain Communication

TODO

Data governance is an essential practice aimed at establishing policies, procedures, and standards to ensure data integrity, confidentiality, accessibility, and quality. In an inter-domain context, good data governance facilitates communication and information sharing, promoting collaboration and operational efficiency.
### The Reliable Use of Inter-Domain Data and Its Importance for Systems Interoperability

The reliable use of inter-domain data is of crucial importance for systems interoperability. It enables smooth integration of systems, fostering effective collaboration and decisions based on reliable and up-to-date information.

TODO
# ![](https://lh6.googleusercontent.com/6ljD44-qXhOsuWYqji-dpG-mRnk2-7kZCTkbpJZdYo8nhuW_UppR6k3b_pmANW0ASsOway6iujwEgclEhj2_4wyvSh1vDYolvV8Qm8nQam2bl94n8VsmOvTtO58o2K7KvFrMP8yVBXx4FDNQqK8vISk) 

## Conclusion

TODO

Data management is a major challenge in today's digital world. From the omnipresence of data to its evolution as a product, it is essential to adopt robust standardization and governance practices to ensure the efficiency and interoperability
