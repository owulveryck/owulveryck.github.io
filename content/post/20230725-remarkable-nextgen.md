---
title: "Evolving the Game: A clientless streaming tool for reMarkable 2"
date: 2023-07-25T15:55:21+02:00
lastmod: 2023-07-25T15:55:21+02:00
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

<!--more-->

  * Briefly describe your role as a product manager and its importance
  * Explain your goal to enhance user experience with the new version of the tool

## The Evolution from Old to New
  * Outline the previous version and its dependence on the client part
  * Discuss the vision for a more versatile tool capable of streaming content to any device
  * Highlight the removal of the need for any installation on the client side

## New Architecture
  * Detail the new structure, with the focus on rendering images directly in the browser
  * Explain the use of native instructions for writing into a canvas and rotating the image
  * Discuss the initial use of websockets for validating the proof of concept

## Moving Away from Websockets
  * Explain the challenges with the websocket approach, including device compatibility and message overhead
  * Describe the transition to a stream of raw data to address these issues

## Network Consumption Optimizations
  * Discuss the challenge of high network consumption even after moving to raw data
  * Detail the first optimization step: encoding the picture in uint4 to store 2 pixels into one byte
  * Describe the implementation of a simple compression algorithm, its trade-offs (memory/cpu), and how you managed it
  * Discuss the implementation of Run-Length Encoding (RLE) and its storage efficiency
  * Explain how storing each pixel in a count of 15 allowed for compact byte count and value storage

## Final Touches
  * Detail the move to event-driven streaming, triggered only after something was written
  * Discuss any additional features or enhancements made in the final stages

## Conclusion
  * Recap the evolution of the tool and the improvements made
  * Discuss the positive impact this will have on user experience
  * Briefly look to the future: what might be next for this tool, or similar tools in development.

