---
title: "From a project to a product: the state of onnx-go"
date: 2019-04-03T20:33:42+02:00
lastmod: 2019-04-03T20:33:42+02:00
draft: false
keywords: []
description: "In this post, I am describing the state of the `onnx-go` project that I started a couple of months ago. The purpose of this post is to explain why I started the project, how I developed the idea, and what the package `onnx-go` is."
tags: []
categories: []
author: "Olivier Wulveryck"

# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: true
toc: true
autoCollapseToc: true
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: false
---
<script src="/js/fabric.min.js"></script>
<script src="/js/wasm_exec1.12.js"></script>
<script src="/js/loader_onnx.js"></script>

<!--more-->
In this post, I am describing the state of the `onnx-go` project that I started a couple of months ago. The purpose of this post is to explain why I started the project, how I developed the idea, and what the package `onnx-go` is.

## When onnx-go was a project...
On March the 25th, I was at the dotGO conference. I gave a lightning talk about `onnx-go`, a project I started a couple of months ago.

The purpose of `onnx-go` is to be able to import pre-trained neural network models (software 2.0) into the Go ecosystem.

It's been a while since machine learning started to buzz.
However, the software developers still need to learn many concepts before being able to use a neural network as a regular capacity in their code.

> if neural networks is a software 2.0, a "regular" developer should be able to use it, like any other library, easily. 

## But why Go?
A study showed that 90% of the costs of software were related to its maintenance. Therefore, the total cost of ownership of modern software is correlated mainly to three factors:

* its reliability;
* the cost for extending it (adding a new feature for example);
* its ability to efficiently run at scale.

<center>
![Mandatory XKCD](https://imgs.xkcd.com/comics/automation.png)
</center>

The way the Go authors have addressed those concerns is interesting.

Go is a simple language to learn and to use even on complex use-cases (the SDK is battery-included; no need for a third-party framework).
On top of that, the language is highly opinionated in a way that allows developers to focus on what they do instead of how to do it.

But most of all: a Go program do not increase the toil when running at scale.
Its self-contained static binary is one of the most DevOps friendly principles so far. It's probably one of the reasons why it is so widely present in the hosting and Ops ecosystem.

> Combining the capacity offered by ONNX and the Go language can make programming with neural network fun again.

## Is onnx-go a product?

If you read my blog, you may remember the first experiments I made with ONNX and Gorgonia.
The first POC was a project called Gorgonnx and was strongly linked with the Gorgonia Library.
Sadly, Gorgonnx was conceptually ok, but a weird bug was preventing the MNIST model from running successfully.

I wrote a [tool](https://github.com/owulveryck/gorgonnx/tree/master/utilities/tests_gen) to import the integration tests of the operators from the ONNX project to the Go ecosystem.
I used it to tests the bridge `onnx-go` to `Gorgonia` for a given set of operators: the set used in MNIST.
All the tests ran successfully but still no luck with the MNIST model.

As I could not figure out where the bug was, I decided to decouple onnx from Gorgonia.

### About graphs...

onnx is a description of a neural network through its computation graph. Gorgonia is a computation graph and the capacity to run the computation graph.

Therefore, decoupling onnx from Gorgonia is easy via an interface, and gonum's graph interface is the perfect fit.

`onnx-go` do not expose a raw compilation of the protobuf definition in Go anymore. Its purpose is to provide helper functions in Go to decode ana onnx binary file and create a computation graph.

So onnx-go is becoming a product and the vision statement is:

> For the Go developer 
> who needs to add a machine learning capability to his/her code
> `onnx-go` is a package 
> that facilitates the use of neural network models (software 2.0)
> and unlike any other computation library,
> this package does not require special skills in data-science.

## What's inside onnx-go now?

Several attempts to build a stable package were made while hunting the bug of the MNIST model.
I tweaked Gorgonia a lot to make it compatible with `gonum's WeightedDirectedGraphBuilder` interface. The goal was to be able to use a straightforward method similar to the JSON's `Unmarshal` function from the Go standard library:

```go
func Unmarshal(data []byte, g graph.WeightedDirectedBuilder) error
```

But this definition is not enough. 

From ONNX, I want to be able to generate simple graphs (for example to display), as well as computation graph. And a computation graph needs special features. Its nodes needs to carry _data_ and _operators_.

Therefore, I need two more interfaces: 

The nodes can be `OperationCarrier` and `DataCarrier`.


The experience told me that the graph should fulfill the `OperationCarrier`, not the nodes. 
The reason is that in my representation of a computation graph, the operands of a simple mathematical expression are the children of a node "operator." For noncommutative operations, I have decided to use weights to guarantee the order of the operands in regards to the operator.

As of today, the interface `OperationCarrier` is: 

```go
type OperationCarrier interface {
	ApplyOperation(Operation, graph.Node) error
}
```

Let's put aside what and `Operation` is by now and just consider it's an object with a name and some named attributes (for example "Convolution" with "padding, strides, dilation").

The `DataCarrier` was first a `TensorCarrier` because all I needed to address were tensors in ONNX. A talk with [Xuanyi Chew](https://twitter.com/chewxy) convinced me that, even if the carried type was a `tensor.Tensor`, technically speaking a tensor is data+operation... Well, `DataCarrier` is a good name.

By now, the interface for `DataCarrier` is: 

```go
type DataCarrier interface {
	SetTensor(t tensor.Tensor) error
}
```

A Graph can produce node that are `DataCarrier`. 

So far, I do not want to add an extra method in the interface the graph should fulfill. This would lead to a new method to generate a node and would mess up in the implementation of the backend as we would need two methods to create a node that would return the same concrete type.

For example:
```go
// NewNode required to fulfill the gonum.Graph interface
func (b *backend) NewNode() graph.Node {
	return &customnode{}
}
// NewVertex to fulfill the onnx.DataCarrier interface
func (b *backend) NewVertex() onnx.DataCarrier {
     return &customnode{}
}
```

A type asertion is performed at runtime. 

### Current state

This API has been promoted to the master branch last week, just after the conference. By now, I don't have any compatible backend, and `onnx-go` is. for now, not production ready.
All the old implementations are out-of-date. Therefore, it's still possible to run the MNIST model on the legacy code, but I have decided to freeze it in order to focus on a stable release. It's important as I am now trying to federate people and to find help for developing the product.

## What's next

I need to reimplement Gorgonia as a working backend. The goal is to be able to run the MNIST model and to quickly add some new operators to be able to run more models from the model zoo.

### Tests

I learn TDD from my friends (thanks [Jonathan](https://twitter.com/jducraft)) and I know that it's a design method more than a test method.
I am not, by now, able to write my code in TDD, but driving the design by the test is, in my humble opinion, mandatory.
This is what lead me to the current implementation which will keep on moving.

So the next step is to write helper functions to test a backend and eventually its coverage of the ONNX API.
Once all of this is ok, I can re-implement the Gorgonia operators one by one, and testing them on the fly.

On top of that, I am also considering writing a very simple and light backend based on gonum without symbolic computation. This can be useful for to run on small plateforms such as a RPI.

Stay tuned and follow [github.com/owulveryck/onnx-go](https://github.com/owulveryck/onnx-go)!

## One more thing...

Before moving on to the new implementation of ONNX, I've setup a small demo that runs entirely in the browser (I've compiled the Go code to WebAssembly).
This is a WASM port of the sample I made with Gorgonia and onnx-go to load the MNIST model.

You can load the the MNIST pre-trained model and play with it!

**Warning** Just a couple of warning. It can hang a tab of your browser (and maybe all the browser).
I don't think it works on a mobile. 
For more information on the wiring that makes this work with WASM etc., please refer to this [post](/2018/06/11/recurrent-neural-network-serverless-with-webassembly-and-s3.html).

---

Download the MNIST `model.onnx` from the [model zoo (you may need to unzip it first)](https://www.cntk.ai/OnnxModels/mnist/opset_7/mnist.tar.gz) (if not available, there is a copy [here](/assets/onnx/mnist/model.onnx)). This file contains the pre-trained neural network.


<a style="color:red">Upload it here: </a><input type="file" id="knowledgeFile" multiple size="1" style="width:250px" accept=".onnx" onChange="fileOK();">

Load the onnx-go small interpreter into the browser (the WASM file is not preloaded to spare bandwidth): <button onClick="load();" id="loadButton" style="width:125px;color:grey;" disabled>Load</button>

Wait for the file to be compiled by your browser, and then click on this button when it turns green: <button onClick="run();" id="runButton" style="width:125px;color:grey" disabled>Run</button>

<!-- Then draw a digit (around the center of the square)... and the machine will tell you what it is -->

<center>
<p id="info" style="color:green;">...</p>
<canvas style="touch-action: none;" id="canvasBox"position="relative" width="280" height="280"></canvas>

<button id="btnSubmit">Submit</button>
<button onclick="reset()">Reset</button>
<br/>

</center>

For the most curious, here is how it works internally...
<center>
{{< figure src="/assets/images/onnx-wasm-principle_small.png" link="/assets/images/onnx-wasm-principle.png" title="Schema" >}}
</center>


