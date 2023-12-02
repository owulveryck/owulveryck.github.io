---
title: "Simplifying Complexity: The Journey from WebSockets to HTTP Streams"
date: 2023-12-02T08:26:41+01:00
lastmod: 2023-12-02T08:26:41+01:00
draft: false
images: [/assets/crowdasleep_small.png]
videos: [/assets/present.webm]
keywords: []
summary: This article explores the transition from a WebSocket-based implementation to a simpler, more direct stream over HTTP in the context of capturing touch screen inputs on Linux. 
  
  
  It begins by introducing the main theme, encapsulated in the statement _Everything ~~is a file~~ is a stream of byte._ 
  The need to capture finger positions on a touchscreen by reading `/dev/input/events` in Linux is initially discussed, followed by a dilemma of transferring this data to a JavaScript client in a browser.
  
  
  Initially, WebSockets are chosen, leading to a discussion on how frameworks often shape our technological choices and the challenges faced in debugging WebSocket connections.
  The article then introduces an alternative about sending a stream of bytes over HTTP, drawing a parallel to Linux's approach to handling devices and files.
  
  
  Serialization, the process of encoding messages for this stream, is discussed next, highlighting the implementation specifics in GoLang and its native advantages. 
  The final section covers how to receive and decode this stream in JavaScript within a worker thread, and then send the decoded messages to the main thread using post requests.
  
  The article concludes by reflecting on the benefits of simplicity in technology, urging readers to reevaluate default choices and consider more straightforward solutions to complex problems.
tags: ["stream", "reMarkable", "http", "golang", "websocket", "javascript"]
categories: ["tools"]
author: "Olivier Wulveryck"

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: false
toc: true
autoCollapseToc: true
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

## Introduction
- Introducing the main theme: "Simple is complex."
- Key statement: "Everything ~~is a file~~ is a stream of byte."

## Background
- Explaining the need: capturing finger positions on a touchscreen by reading `/dev/input/events` in Linux.
- Technical methodology for achieving this.

## The Problem Statement
- Challenge described: Transferring information from server to a browser-based client in JavaScript.

## The Default Choice: WebSockets
- Rationale behind initially choosing WebSockets.
- Influence of frameworks on technology decision-making.
- Challenges and difficulties encountered with WebSocket debugging.

## An Alternative Approach: HTTP Streams
- Presenting HTTP streams as an alternative.
- Advantages of simplicity and directness in this approach.

## The Concept of Serialization
- Discussion on the necessity and methods of encoding messages.
- Exploration of different serialization techniques.

## Implementing HTTP Stream in Go
- Showcasing the GoLang code for streaming.
- Emphasis on the use of native language capabilities.

## Receiving and Decoding the Stream in JavaScript
- Process of stream reception in JavaScript via a worker thread.
- Details on decoding and message transfer to the main thread through post requests.

## Conclusion
- Reflecting on the journey and learnings.
- Reinforcing the core message: simpler solutions for complex problems.
- Encouraging exploration beyond default choices.

### Optional Sections
- **Case Studies or Examples**: Including real-life applications.
- **Performance Metrics**: Comparison between WebSocket and HTTP stream performance.
- **Future Implications**: Potential impact on future technological choices.

