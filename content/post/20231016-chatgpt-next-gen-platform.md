---
title: "The Dawn of LLM-Powered personal assistants: pioneering a new platform economy"
date: 2023-10-16T10:11:52+02:00
lastmod: 2023-10-16T10:11:52+02:00
draft: true
keywords: []
description: ""
tags: []
categories: []
author: "Olivier Wulveryck	"
images: [/assets/chatgpt-platform-me.png]
summary: In the digital age, traditional PowerPoint presentations often fail to engage audiences due to their static nature. 

comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---


## Introduction
The **platform economy** is an economic environment in which digital platforms act as intermediaries, connecting various stakeholders and enabling them to transact with one another seamlessly.
These platforms harness the power of network effects, creating value by facilitating exchanges between users.

Consider a digital marketplace as an example:

- They act as a connection hub between sellers (supply) and buyers (demand).
- They leverage the network effect: as more people use the marketplace, the value of the platform increases for all users. A larger selection of products attracts more buyers, and more buyers attract more sellers.
- They derive value from technology: the platform economy is largely enabled by digital technology, and marketplaces leverage this technology to operate, scale, and improve their platforms over time.

In recent times, advancements in Language Models like **ChatGPT**, powered by Large Language Models (LLMs), have paved the way for the development of highly intuitive and capable personal assistants.
These LLM-powered assistants can become an evolution of the platform economy.
A platform where **digital storefronts**, **service providers**, and **users** **converge**, creating a **symbiotic ecosystem**.

In this article, I will explain how the _plugin_ mechanism of ChatGPT is a keystone in the construction of the platform.

In the first part, I will explain a fictional use case from the point of view of the buyer.
In the second part, I will expose how a supplier can use the plugin mechanism to meet the platform requirements.
I will insist on the communication and standardization part, discussing intermediate representation and human language.

## Example use case

Consider a scenario where one needs to prepare for a dinner engagement with a budget of **$100**.
The dinner is slated for 8pm, with a one-hour drive, and it's currently 2pm.

The aim is to find suitable attire within budget, and purchase it from a nearby store to make it to the dinner on time.

Here's where an LLM-powered assistant like ChatGPT comes into play.
The assistant, leveraging a network of digital platforms, can help identify clothing options, compare prices, and locate a nearby store, all while ensuring the user stays within budget and timeline.

![](/assets/chatgpt-platform-illustration_small.png)

### Rationale
With my colleague [Nicolas](https://www.linkedin.com/in/nicolasgutierrez/), we conducted the exercise to create a top (Wardley) map for this need.
The fundamental requirement is about finding a good.
Given Nicolas's background in knowledge graphs and semantic search, we initially explored how these elements could address the use-case.
Subsequently, we considered the potential role of personal assistants in this context.
Below is the resulting map:
![](/assets/chatgpt-kg.png)

We then asked: Does the LLM need to be in-house for the retailer, or can it operate independently using the data provided by the retailer?
The colors indicate the components owned by the retailer (note, I am not necessarily stating they should remain in-house, especially the commodities).
Components owned by the end-user are shown in green.

Surprisingly, after a discussion on the map, it became clear that the LLM instantiation might be more beneficial if kept out of the context.
This is logical, as the computing power required to run the LLM is substantial and should probably remain with an external company that specializes in building personal assistants.

So what could be the future of the personal assistant and its associated business.

My guess is that it is a next generation platform model.

Let's assess the current technical landscape and how it might evolve in the future based on the technical elements that exist today.

## Platform? 

### Platform Economy: From Pipelines to Platforms

A **platform** efficiently connects **producers** and **consumers**, allowing them to generate value through their large-scale interactions.
While this might seem like a familiar concept, **Sangeet Paul Choudary** introduces an additional perspective: the notion of a **pipeline**.

In business, a pipeline represents a linear and unidirectional transformation that enables a producer to create value and deliver it to the consumer.
Essentially, it's the **traditional system** of goods or service providers.
For instance, **Apple** provides value by creating products and selling them to consumers, transitioning from a set of components to a finished product through a series of pipelines.
These pipelines have been revolutionized by technology, impacting three main pillars:

1. **Facilitating mass production**.
2. **Encouraging mass consumption** (e.g., the television's influence on consumer behavior).
3. **Easing international exchanges and transactions**, promoting system connectivity.

However, the advent of the **Internet** and **digital technology** has further evolved these pillars:

- **Production tools** can be more easily distributed. Previously, to produce information, one needed to be a newspaper. Now, anyone can produce information on platforms like Twitter or LinkedIn.
- **Digitalization** has brought about personalized usage patterns, offering consumers tailored products.
- The Internet has influenced **mass consumption** by affecting prices, as seen with giants like **Amazon** and **Alibaba**.

Speaking of the market, one model employed in this market economy is based on **Porter's Five Forces**, which assesses:

1. **The threat of new entrants**.
2. **Supplier bargaining power**.
3. **Customer bargaining power**.
4. **The threat of substitute products or services**.
5. **Competitive rivalry**.

This model evaluates the forces that might influence a company's profitability and formulates strategies accordingly.
However, with the Internet, these forces act differently, requiring traditional "pipeline" companies to adapt.

The **platform's main idea** is to transform value creation.
Its value lies in effectively linking producers and consumers.
Taking Apple's example from the HBR paper, the **App Store platform** facilitates mass app production for a large consumer base.
It shifts from mass production linked to mass consumption to distributed production connected to personalized consumption.
The platform provides an interface easing the onboarding of new producers while also ensuring **governance** by implementing rules for both producers and consumers, ensuring, for example, that App Store applications are safe for users.


Considering our previous use-case, many suppliers are involved in the fulfillment of the users need. How to interconnect the suppliers with the buyers ?

The supplier needs to provide its service offer to the buyer, therefore acting as a **content provider**. 
the **Content Providers:** Stores and brands act as content providers.
They supply information regarding their products, including pricing, availability, and location.

![](/assets/platform_pipeline.png)
  


- **Content Production:** Plugins serve as content producers for the platform, bridging the gap between users and content providers.

- **Plugins:** These are interfaces to the platform.
They facilitate structured interactions between the platform and external systems through APIs, alongside human-readable descriptions.
By doing so, they ensure the assistant understands the user's needs and the context surrounding them.

  - **Intermediate Representation:** Plugins work by crafting an intermediate representation encompassing structured API interactions and human-centric descriptions.
This blend ensures a seamless flow of data, making it understandable both by the platform and the end-users.



## New Challenges:
As this ecosystem grows, new challenges akin to the SEO challenges faced by search engines will emerge.
The concept of **Prompt Engineering and Automatic Choice Optimization** will become crucial.
This entails optimizing how requests are processed and responded to by the LLM-powered assistants, ensuring accuracy, relevance, and efficiency in meeting the users' needs.

## Conclusion:

![](/assets/chatgpt-platform_small.webp)
The advent of LLM-powered personal assistants like ChatGPT heralds a new era in the platform economy, merging the digital and physical realms in a user-centric ecosystem.
By integrating digital storefronts, service providers, and users, a new level of value creation and exchange is realized, redefining how we interact with the digital world around us.
