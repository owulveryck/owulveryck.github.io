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
images: [/assets/chatgpt-platform.webp]
summary: In the digital age, traditional PowerPoint presentations often fail to engage audiences due to their static nature. 

comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---


## Introduction:
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

## Use Case:

Consider a scenario where one needs to prepare for a dinner engagement with a budget of **$100**.
The dinner is slated for 8pm, with a one-hour drive, and it's currently 2pm.

The aim is to find suitable attire within budget, and purchase it from a nearby store to make it to the dinner on time.

![](/assets/chatgpt-kg.png)

Here's where an LLM-powered assistant like ChatGPT comes into play.
The assistant, leveraging a network of digital platforms, can help identify clothing options, compare prices, and locate a nearby store, all while ensuring the user stays within budget and timeline.

## The Platform Architecture:

![](/assets/chatgpt-platform.webp)
- **Content Providers:** Stores and brands act as content providers.
They supply information regarding their products, including pricing, availability, and location.
  
- **Plugins:** These are interfaces to the platform.
They facilitate structured interactions between the platform and external systems through APIs, alongside human-readable descriptions.
By doing so, they ensure the assistant understands the user's needs and the context surrounding them.

  - **Intermediate Representation:** Plugins work by crafting an intermediate representation encompassing structured API interactions and human-centric descriptions.
This blend ensures a seamless flow of data, making it understandable both by the platform and the end-users.

- **Content Production:** Plugins serve as content producers for the platform, bridging the gap between users and content providers.

## New Challenges:
As this ecosystem grows, new challenges akin to the SEO challenges faced by search engines will emerge.
The concept of **Prompt Engineering and Automatic Choice Optimization** will become crucial.
This entails optimizing how requests are processed and responded to by the LLM-powered assistants, ensuring accuracy, relevance, and efficiency in meeting the users' needs.

## Conclusion:
The advent of LLM-powered personal assistants like ChatGPT heralds a new era in the platform economy, merging the digital and physical realms in a user-centric ecosystem.
By integrating digital storefronts, service providers, and users, a new level of value creation and exchange is realized, redefining how we interact with the digital world around us.
