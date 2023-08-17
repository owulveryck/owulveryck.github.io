---
title: "Evolving the Game: A clientless streaming tool for reMarkable 2"
date: 2023-07-25T15:55:21+02:00
lastmod: 2023-07-25T15:55:21+02:00
draft: false
keywords: []
summary: In this article, I expose the new version of the streaming tool for my reMarkable tablet. 

  
  Originally crafted in 2021, this tool allowed me to stream sketches during video calls. 


  Aiming for more user-friendliness, I revamped the design to a clientless implementation.
  This articles describes the implementation with some code illustrations in Javascript and Go on how to

  * get the picture and to display it in a canvas

  * optimize the flow by playing with `uint4` and `RLE`
tags: ["golang", "reMarkable", "JS", "Optimization", "RLE", "ChatGPT", "hack", "WebSocket"]
categories: []
author: ""

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

In 2021, I developed a tool to stream the content of my reMarkable. 
(and I also blogged about it [here](https://blog.owulveryck.info/2021/03/30/streaming-the-remarkable-2.html)).
Given that I was predominantly working from home, this tool was very useful, allowing me to sketch elements during conference calls.

One of the primary benefits of this tool was its ability to stream content directly into a web browser tab. 
This feature was particularly useful because it meant I could exclusively share this tab during video calls, ensuring focus on the content I intended to present.

At its core, the tool consisted of two main components:
1. A server running on the device, responsible for capturing the raw image and transmitting it to the laptop.
2. A service on the laptop, which fetched the raw image from the server and processed it into a format suitable for browser viewing (I produced an MJPEG stream for this).

Being the product manager of my own tools offered a unique perspective. 
One piece of feedback I provided from my experience as a user was the slightly cumbersome nature of the tool's on-the-fly activation. 
The challenge arose from the necessity of a local service. During video calls, this meant some preparatory steps to initiate the local service, adding an additional layer of complexity to the process.

Recognizing this pain point, my goal became clear: eliminate the need for the local service. 
This article delves into the revamped implementation of the streaming tool, which now boasts a more user-friendly design and improved performances.

## The Evolution from Old to New

The code running on the device must have a low footprint. 
One way to ensure it remains lightweight is by eliminating any heavy computation on the device. 
The sole function of the code running on the server is to grab the raw picture from the memory and expose it over the network. 
This led to a three-tier design: server/client/renderer.

_Ubiquitous language:_

In this article:

- **The server** refers to the code running on the reMarkable (the device). Its main purpose is to expose the raw image of the current display on the reMarkable.
- **The client** is responsible for fetching the raw image from the server and performing additional processes to convert it into a usable format.
- **The renderer** accepts the output from the client and displays it on a PC screen.

```plain
+---------------------------+        +-----------------------------------+
|          reMarkable       |        |               Laptop              |
|                           |        |                                   |
|       +-------+           |  gRPC  |       +-------+                   |
|       |Server |<--------- |<------>|------>|Client |                   |
|       +-------+           | Fetch  |       +-------+                   |
|                           | Stream |           |                       | 
|                           |        |     HTTP MJPEG stream             |
|                           |        |           |                       |
|                           |        |           v                       |
|                           |        |       +--------+                  |
|                           |        |       |Renderer|                  |
|                           |        |       +--------+                  |
|                           |        |      (Browser/VLC)                |
+---------------------------+        +-----------------------------------+
```

To minimize CPU usage, the server extracts the picture only when the client is connected.
This functionality was achieved through gRPC communication.

The server can thus operate as a daemon on the reMarkable, awaiting an RPC call from the client.
To initiate streaming, I simply needed to activate the client.
The client retrieves the image in a loop, and each image is encoded in JPEG before being added to an MJPEG stream.
This stream is then made available as an HTTP service by the client.

The renderer is any software capable of reading the MJPEG stream via HTTP, such as VLC or a web browser.

One challenge with this setup is that it necessitates a specific network topology and configuration.
The client must not only be aware of the reMarkable's address but also possess the necessary privileges to establish a server.
Additionally, the renderer must be familiar with the client's IP address.

While this wasn't an issue in my personal setup, complications arose post-pandemic when I returned to in-person presentations.
I realized the need for a more straightforward solution.
My ultimate aim became the ability to simply input the reMarkable's address into any browser and instantly access the stream.

### New Architecture

To achieve the objective, the solution involves eliminating the client and instead establishing an HTTP server within the server component.
The client should be implemented in a format that's interpretable by a browser, such as Javascript or WASM.

My initial approach was to compile the client into WASM.
This seemed promising as it would let me leverage my expertise in Go development.
However, I encountered several limitations that would have necessitated substantial modifications.

Consequently, I opted to develop a second version of the tool, with the client written in Javascript.

_Side note:_ At this juncture, I confronted another challenge.
While I possess a broad understanding of Javascript's functionality and browser rendering processes (what we might term "architectural skills"), 
I felt less confident in my hands-on JS development abilities.
I turned to my digital assistant, ChatGPT, for guidance.
With my direction on the desired solution, it provided the necessary code fragments and explanations to bring my vision to life.
I was the developer, he was the coder.

### Validation of the "canvas" renderer

Initially, it was imperative to move away from the MJPEG stream, especially since my operations were now closely aligned with the renderer, 
and Javascript possesses the required primitives for image manipulation.

In the browser, the conventional method for handling images is via the `canvas` element.
My preliminary task was to validate that I could retrieve a raw image from the server and showcase it within a `canvas`.

I achieved this by accessing the backbone of the canvas that represents the pixel map data in RGBA format and adjusting the pixels based on their values in the raw image from the reMarkable:

```js
<canvas id="fixedCanvas" width="1872" height="1404" class="hidden"></canvas>
<script>
    // Use the fixed-size canvas context to draw on the canvas
    var fixedCanvas = document.getElementById("fixedCanvas");
    var fixedContext = fixedCanvas.getContext("2d");
    function processBinaryData(data) {

        // Assuming each pixel is represented by 4 bytes (RGBA)
        var pixels = new Uint8Array(data);
        // Create an ImageData object with the byte array length
        var imageData = fixedContext.createImageData(fixedCanvas.width, fixedCanvas.height);
        // Assign the byte array values to the ImageData data property
        for (var i = 0; i < pixels.length; i++) {
            imageData.data[i*4] = pixels[i];
            imageData.data[i*4+1] = pixels[i];
            imageData.data[i*4+2] = pixels[i];
            imageData.data[i*4+3] = 255;
        }

        // Display the ImageData on the canvas
        fixedContext.putImageData(imageData, 0, 0);
    }
```

There's also a requirement to adjust the image for responsiveness according to the browser's size, as well as for image rotation and potential colorization.

To achieve this, I maintain the `fixedCanvas` in a hidden state and transfer its contents to another canvas using the `drawImage` method.
The dimensions of the destination canvas (its width and height) undergo adjustments if a resize event is detected within the browser window.

{{< highlight js "hl_lines=4 32" >}}
var resizableCanvas = document.getElementById("canvas");
var resizableContext = resizableCanvas.getContext("2d");
function copyCanvasContent() {
    resizableContext.drawImage(fixedCanvas, 0, 0, resizableCanvas.width, resizableCanvas.height);
}

// JavaScript code for working with the canvas element
function resizeCanvas() {
    var canvas = document.getElementById("canvas");
    var container = document.getElementById("container");

    var aspectRatio = 1872 / 1404;

    var containerWidth = container.offsetWidth;
    var containerHeight = container.offsetHeight;

    var containerAspectRatio = containerWidth / containerHeight;

    if (containerAspectRatio > aspectRatio) {
        canvas.style.width = containerHeight * aspectRatio + "px";
        canvas.style.height = containerHeight + "px";
    } else {
        canvas.style.width = containerWidth + "px";
        canvas.style.height = containerWidth / aspectRatio + "px";
    }

    // Use the canvas context to draw on the canvas
    copyCanvasContent();
}

// Resize the canvas whenever the window is resized
window.addEventListener("resize", resizeCanvas);
{{< / highlight >}}

### Basic Drop-In Replacement

With the renderer in place, my next step was to implement a lightweight JS client as a substitute.
While gRPC offers robust functionality, it's not typically considered the gold standard in web development.

Thus, for communication and encapsulation, I gravitated towards the WebSocket protocol.
It was straightforward enough to be incorporated server-side, serving as a seamless drop-in replacement for the RPC server.

Messages delivered via this protocol carry the raw image.
The client remains on standby for these messages, and with each incoming message, the canvas content is updated, effectively emulating a streaming process.

One notable advantage of this approach is the enhanced control over server-side load.
Extracting the raw image demands both memory and CPU resources on the device.
By regulating the message emission frequency, I can efficiently manage device load.

## Moving Away from Websockets

The implementation with Websockets was working 
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

