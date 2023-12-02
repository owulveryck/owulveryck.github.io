---
title: "Simplifying Complexity: The Journey from WebSockets to HTTP Streams"
date: 2023-12-02T08:26:41+01:00
lastmod: 2023-12-02T08:26:41+01:00
draft: false
images: [/assets/crowdasleep_small.png]
videos: [/assets/present.webm]
eywords: []
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

To add a new functionality to my tool, [goMarkableStream](https://github.com/owulveryck/goMarkableStream), I needed to capture gesture positions from the tablet's screen and relay them to the browser to trigger local actions.
For example, a swipe left could activate a specific function in the browser.

My approach involved capturing gestures from the device and then communicating them to the browser with a message stating: "this gesture has been made."

In the realm of message exchange between a server and a client in a browser, WebSockets naturally come to mind.
WebSockets are inherently designed to support streams of messages on top of TCP, unlike HTTP, which primarily handles streams of bytes without a built-in concept of a message.

Navigating through this journey, I realized the importance of extensive testing and learning to craft an effective solution.
The WebSocket protocol, in contrast to HTTP, introduces distinct challenges, especially in debugging and testing, due to its more complex nature.

Acknowledging that gestures are essentially a stream of bytes (I will explain this), I will write about:
- the process of evaluating the trade-off between the added complexity of WebSockets and the functionalities they offer.
- how I streamlined the development by designing my own messaging system over HTTP.

## Background

I've previously discussed using my tablet for various live presentations.
Through iterative testing, I developed a hybrid solution that combines elements of a whiteboard with static slides.
This solution features the screen's drawing as an overlay on existing slides.
The challenge now lies in switching slides directly from the tablet to streamline the presentation and minimize interactions with the laptop displaying the slides.

The slides are displayed within an iFrame on the client side of my tool.
Consequently, I needed a method to send commands to the iFrame to control slide transitions.
The [reveal.js presentation framework](https://revealjs.com/) supports native embedding and allows slide control from the top frame via an API that uses [postMessages](https://revealjs.com/postmessage/).

To transmit slide control commands from the tablet to the client, I considered various methods.
The optimal solution I identified was to utilize touch gestures on the screen of the reMarkable tablet.
By swiping on the tablet, I could send events to the client, which would then respond accordingly to switch the slides.

### Capturing the touch events on reMarkable/Linux

The reMarkable operates on a Linux-based system.
Input events (both pen and touch) are managed through [Event Devices (evdev)](https://en.wikipedia.org/wiki/Evdev).
The event exposure is as follows:
- `/dev/input/event1` captures the pen events.
- `/dev/input/event2` captures touch events.

In Unix, the philosophy that "_everything is a file_" applies.
This means I can easily access these events by opening and reading the file contents in Go.
I chose Go for the server-side language due to its self-sufficient packaging, cross-compilation capabilities, and the enjoyment I derive from using it.

> "Everything is a file" is a principle in Unix and its derivatives, where input/output interactions with resources such as documents, hard-drives, modems, keyboards, printers, and some inter-process and network communications are treated as simple byte streams accessible through the filesystem namespace - [source Wikipedia](https://en.wikipedia.org/wiki/Everything_is_a_file).

### Reding the events in Go

The "file" event is a character device, offering a binary representation of an event.
In Go, an event's set of bytes could be structured like this:

```go
type InputEvent struct {
	Time syscall.Timeval `json:"-"`
	Type uint16
	Code  uint16
	Value int32
}
```

The principle of "_everything is a file_" enables the use of basic operations from the `os` package to open the character device as an `*os.File` and `Read` the binary representation of the event.
We create an `ev` object of the `InputEvent` type to receive the information read.

The file functions as an `io.Reader`, and its content is typically loaded into a byte array.

```go 
func readEvent(inputDevice *os.File) (InputEvent, error) {
    // Size calculation: 
    // Timeval consists of two int64 (16 bytes), 
    // followed by uint16, uint16, and int32
    // (2+2+4 bytes)
    const size = 16 + 2 + 2 + 4
    eventBinary := make([]byte, size)

    _, err := inputDevice.Read(eventBinary)
    if err != nil {
        return InputEvent{}, err
    }

    var ev InputEvent
    // Assuming the binary data is in little-endian format 
    // which is the most common on Intel and ARM
    ev.Time.Sec = int64(binary.LittleEndian.Uint64(eventBinary[0:8]))
    ev.Time.Usec = int64(binary.LittleEndian.Uint64(eventBinary[8:16]))
    ev.Type = binary.LittleEndian.Uint16(eventBinary[16:18])
    ev.Code = binary.LittleEndian.Uint16(eventBinary[18:20])
    ev.Value = int32(binary.LittleEndian.Uint32(eventBinary[20:24]))

    return ev, nil
}
```
A more efficient approach could involve using an unsafe pointer to directly populate the structure, thereby bypassing Go's safety mechanisms by using the `unsafe` package:

```go
func readEvent(inputDevice *os.File) (events.InputEvent, error) {
	var ev InputEvent
    // by using (*[24]byte), we are explicitly stating that 
    // we want to treat the memory location of ev as a byte array of length 24
    // We could have used the less readable form:
    // (*(*[unsafe.Sizeof(ev)]byte)(unsafe.Pointer(&ev)))[:]
    // 
    //  Note: the trailing [:] is mandatory to convert the array to a slice
    _, err := inputDevice.Read((*[24]byte)(unsafe.Pointer(&ev))[:])
	return ev, err
}
```

## The Problem Statement

Now that I have read the events, I need to send to the client for further processing.
The current architecture is based on an HTTP server in Go and a web client in JS. Therefore, I need to find a HTTP-ish way to transfer the events.

It is beyond the scope of this article to delve into the specifics of how I publish events within the Go server.
However, for a basic understanding necessary for the rest of the article, here's a brief overview.

### Serving Structure in the Go Server
Fundamentally, I have implemented a basic [pubsub](https://github.com/owulveryck/goMarkableStream/blob/main/internal/pubsub/pubsub.go) mechanism to channel the flow of events.

The next step is to make these events accessible to the client.
This will be managed by a `http.Handler`. Below is the framework for this handler:

```go
type GestureHandler struct {
    inputEventBus *pubsub.PubSub
}

// ServeHTTP implements http.Handler
func (h *GestureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // eventC is a channel that receives all the InputEvents
    eventC := h.inputEventBus.Subscribe("eventListener")
    // ....
}
```

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

