---
title: "Considerations about software 2.0"
date: 2018-04-16T10:54:23+02:00
lastmod: 2018-04-16T10:54:23+02:00
draft: false
keywords: []
description: "This post describes the concept of software 2.0. It evaluates the idea of implementing an equation parser (written in Unicode) to give a strict separation of the software 1.0 and the software 2.0.

In this article, I will go further in the description of this concept with the help of the famous “char RNN” example."

tags: []
categories: []
author: ""

# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: false  
toc: true
autoCollapseToc: false
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: false
---

**Disclaimer** This is a technical article about a work in progress. The primary goal is to document what I did and to clarify my ideas. A more general and complete article about software 2.0 is in development and should be published on my [company's blog](https://blog.octo.com/) later.

---

This post describes the concept of software 2.0. It evaluates an instance of the Unicode equation parser (as described [here](https://blog.owulveryck.info/2017/12/18/parsing-mathematical-equation-to-generate-computation-graphs---first-step-from-software-1.0-to-2.0-in-go.html)) to give a strict separation of the software 1.0 and the software 2.0.

In this article, I will go further in the description of this concept with the help of the famous “char RNN” example.

# Preliminary concerns

The software 1.0 is, for me, a kind of virtual machine. Its goal is to run the software 2.0.
The software 2.0 is a deep neural network represented by its mathematical equations and the weighted matrices.

The concept is similar to what we find in the Java world. A runtime environment can execute some bytecode. Our bytecode is composed of sequences of float (the weight matrices) and a bunch of equations. There should be a strict separation between the SDK (which owns training mechanism), and an _NNRE_ (Neural Network Runtime Environment) which can load and apply the software 2.0.

# Sample implementation

To illustrate this theory, let's implement a well-known use-case: a char based LSTM.
(Anderij Karpathy has promoted this example on its [blog](http://karpathy.github.io/2015/05/21/rnn-effectiveness/), please read it if you need further explanation about the concept).

The code is a Go package that can be used to build both an SDK or an NNRE. As I have developed it with the help of the Gorgonia package and that it can be compiled in pure Go, I may end with a dependency less NNRE which is kind of cool. Therefore, the neural net could be trained with CUDA on a specific type of host and exported to a mobile or even run on a serverless platform.

_Note_ The computation cost for the evaluation of a neural network is stable, as there is a fixed amount of operation needed to compute an input into an output. Therefore, a serverless platform should be a good runtime candidate. The cost involved in the execution of the software 2.0 should be linear in function of the number of users.
As AWS Lambda supports Go, experimenting should be straightforward.

## Details 

### The core implementation of the LSTM

The LSTM "kernel" that is used is from the wikipedia page

![img](https://wikimedia.org/api/rest_v1/media/math/render/svg/8a0eddfb6f592041ea04bd26526b52ba1cec192c)

Here is the implementation within the LSTM package:
![content](/assets/lstm/lstm_implem.png)
See the [code here](https://github.com/owulveryck/lstm/blob/1581884e9d2de83e1150c04fb815637351082b7a/lstm.go#L39-L46)

We can see that the equations are not transcoded into a computation graph within the source code nor at the compilation stage. Instead, the formulas are written in a proper mathematical language in Unicode. The system evaluates them at runtime.
A separation of the software 1.0 and the software 2.0 is appearing.

_Note_: This implementation is a work in progress. The mathematical formulas should not be hardcoded but should be an input of the system. Therefore, it could be saved and exported alongside the matrices.

### The input and the output of the LSTM

For the system to run, we must feed the neural network with data that it can understand: 

vectors!

Then with this input, the software 2.0 will apply its magic and produce output vectors.
For the system to be useful, we should be able to decode the output into usable data.

In its current implementation, this encoding/decoding is abstract and declared thanks to a Go interface (described in the [_datasetter_](https://github.com/owulveryck/lstm/blob/master/datasetter/definitions.go) package).
It should allow enough flexibility to implement a different use case without any modification of the core of the package.

# Enhancing the model

It is well-known that LSTM can lead to a good prediction system. Therefore such an NNRE can be useful to predict all sorts of chaotic things as long as we can encode them into vectors.
The major problem of the LSTM is the fixed size of the input and the output vectors.
Based on the concepts that I have described in this article, future work could be to implement a seq2seq kernel machine. 

The seq2seq model has an intrinsic encoder and decoder mechanism.
Then the VM should have pluggable encoder and decoder. gRPC would be an excellent candidate to expose the interface. It would allow an immutable implementation of the VM and high flexibility.

The main difference between this and the [seq2seq](https://google.github.io/seq2seq/) Tensorflow framework would be in the exploitation. You would not need to handle all the dependencies of the SDK just to run the code. And this code could be run wherever a Go binary can be compiled: from an ARM-based machine to a serverless architecture. Flexibility at all levels!

# Conclusion

In this TED talk, Lera Boroditsky exposes the way that language shapes the way we think. According to me, this will become a significant concern in the coming years. The software is eating the world, but the code 2.0 may consume the software. A giant step in computer science would be to implement a new language that could shape the way we think the software and find even better use-cases.

<div style="max-width:854px"><div style="position:relative;height:0;padding-bottom:56.25%"><iframe src="https://embed.ted.com/talks/lera_boroditsky_how_language_shapes_the_way_we_think" width="854" height="480" style="position:absolute;left:0;top:0;width:100%;height:100%" frameborder="0" scrolling="no" allowfullscreen></iframe></div></div>



