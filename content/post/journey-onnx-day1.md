---
title: "My journey with ONNX and Go - The beginning"
date: 2018-08-14T20:41:30+02:00
lastmod: 2018-08-14T20:41:30+02:00
draft: false
keywords: []
description: "This is the very begining of my journey with ONNX and Go. In this post I am describing how to decode a ONNX model from its protocol buffers serialized format to a computation graph. I am using the gonum package. Later I will switch to Gorgonia."
tags: []
categories: []
author: ""

# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: true
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

This year has started with a lot of deep thoughts about the software 2.0.
My conclusion (which is slightly different from [Andrej Karpathy's consideration](https://medium.com/@karpathy/software-2-0-a64152b37c35)) is that a software 2.0 is a combination of a Neural network model **and** its associated weights.
This is a concept; now the question is: how to materialize the idea? What artifact represents a software 2.0.

I emitted several ideas and tried one of them: to serialize the mathematical model and the weights.
The main drawback of this idea is that it is not easy to write down and to parse any mathematical equation.
The best way to express a model is, as of today, via its computation graph (this is what most ML frameworks are doing). 

Therefore, switching from a mathematical representation to the computation graph representation might lead to an excellent way to express the artifact of a software 2.0.

# A quick word about ONNX

Describing a computation graph is straightforward. A computation graph is a [Directed Acyclic Graph (DAG)](https://en.wikipedia.org/wiki/Directed_acyclic_graph). Each node  of the graph represents a tensor or an operator.
The challenge is to find a domain specific language (DSL) to describe a graph in a way that it is agnostic of its implementation.

This is the promise of [ONNX](https://onnx.ai/).
ONNX stands for Open Neural Network eXchange (format). The purpose of this project is to establish an open standard for exporting/importing ML models.

<center>
{{< figure src="https://github.com/owulveryck/onnx-go/raw/master/vignettes/imgs/ONNX_logo_main.png" >}}
</center>

In this post, I will describe the first step I have made to be able to read (and hopefully) execute an ML model encoded via ONNX into the Go ecosystem.

*I am, for now, only working to import a model from ONNX into Go. Export may come further.*

# From the _protobuf_ definition to a Go structure

ONNX is a description of a model made in _protocol buffers_ (_protobuf_).
In this section, I will dig a little bit into the _protobuf_ definition file of ONNX. Then let's create a first Go code to read and import a model.

## What are protocol buffers?

The most accurate definition comes obviously from the [official website](https://developers.google.com/protocol-buffers/):

> Protocol buffers are Google's language-neutral, platform-neutral, extensible mechanism for serializing structured data – think XML, but smaller, faster, and simpler. You define how you want your data to be structured once, then you can use special generated source code to easily write and read your structured data to and from a variety of data streams and using a variety of languages.

It is a way to serialize messages. For short, a _protobuf_ file describes an API contract. The contract is expressed via the _protobuf_ DSL. Then we use a compiler to transpile this API into a native language API (for example, you can turn the definition into a set of Java classes, or Go structures).
Once a message is serialized, it is represented in a binary format.

_Note_: I will not go deeper in the _protobuf_ description here. But, in my humble opinion, it is a perfect way to express an API when implementing a machine-to-machine communication. Better than JSON because of its simplicity, efficiency and the ability to validate a schema natively.

But let's go back to ONNX. The primary definition file for ONNX (the API contract) is hosted [here](https://github.com/onnx/onnx/blob/master/onnx/onnx.proto3) and is named `onnx.proto3`.
This is *the* file used to generate bindings to other languages.

To create a bridge between the _protobuf_ binary format and the Go ecosystem, the first thing to do is to generate the Go API. This will allow to read an ONNX file and to transpile it into a Go compatible object.

To do this, I am using the `protoc` compiler. I am also using the alternative compiler [gogoprotobuf](https://github.com/gogo/protobuf) which add some useful features (such as fast Marshaller/Unmarshaler methods).  For clarity, I will not describe how to install and use the `protoc` binary.

Merely running `protoc --gofast_out=. onnx.proto3` will generate a file [onnx.pb.go](https://github.com/owulveryck/onnx-go/blob/master/onnx.pb.go) which is usable out-of-the box.

## onnx-go 

After some discussions with the [official team](https://github.com/onnx/onnx/pull/1328), we agreed that, before the _onnx-go_ package reaches a certain level maturity, it was best to host it on my personal GitHub account. So, as of today, I am hosting the repository here: [github.com/owulveryck/onnx-go](https://godoc.org/github.com/owulveryck/onnx-go). The corresponding Godoc is hosted [here](https://godoc.org/github.com/owulveryck/onnx-go).

This package on its own is enough to read an ONNX format. 

### Testing the package

As I said before, I want to work on the import. So let's get a properly encoded ONNX file to test the generated package.
The ONNX organization has set up a [model repository](https://github.com/onnx/models) (model zoo). From this repository, let's extract the basic MNIST example to get a real life ONNX model.

```
curl https://www.cntk.ai/OnnxModels/mnist/opset_7/mnist.tar.gz | \
tar -C /tmp -xzvf -
```

Now, let's write a simple program that will read the ONNX file and decode it into an object of type [`ModelProto`](https://godoc.org/github.com/owulveryck/onnx-go#ModelProto) (which is the top-level object in the ONNX file).

Then create a very simple Go program to read and dump the model:

{{< highlight go >}}
import (
        "io/ioutil"
        "log"
        onnx "github.com/owulveryck/onnx-go"
        "github.com/y0ssar1an/q"
)

func main() {
        b, err := ioutil.ReadFile("/tmp/mnist/model.onnx")
        if err != nil {
                log.Fatal(err)
        }
        model := new(onnx.ModelProto)
        err = model.Unmarshal(b)
        if err != nil {
                log.Fatal(err)
        }
        q.Q(model)
}
{{</ highlight >}}

_Note_: I am using the [`q`](https://github.com/y0ssar1an/q) package to dump the content as the output is verbose. The result is present in the file `$TMPDIR/q`

# From the Go structure to a Graph

Now that we are able to read and decode a binary file let's dig into the functional explanation.

# Graphs

The ONNX Model document is made of several structures. One of those structures is the [GraphProto](https://godoc.org/github.com/owulveryck/onnx-go#GraphProto).
From the documentation we read that:

> A graph defines the computational logic of a model and is comprised of a parameterized list of nodes that form a directed acyclic graph based on their inputs and outputs. This is the equivalent of the "network" or "graph" in many deep learning frameworks.

As a consequence, the vertices of the graph are composed of nodes that may be Operators or Tensors. The Tensors can be a computable (learnable) element (defined by the type TensorProto) or values (defined in the type ValueInfoProto). Values are actually not computable. This means that a value is not learnable; most likely it is the input of the neural net.

As a consequence, the primary types mandatory to reconstruct the computation graph are:

* [NodeProto](https://godoc.org/github.com/owulveryck/onnx-go#NodeProto) 
* [TensorProto](https://godoc.org/github.com/owulveryck/onnx-go#TensorProto)
* [ValueInfoPRoto](https://godoc.org/github.com/owulveryck/onnx-go#ValueInfoProto)

In ONNX, all those elements are identified by their `name` which is a string. 

The [`GraphProto`](https://godoc.org/github.com/owulveryck/onnx-go#GraphProto) structure is made of:

* a list of _inputs_ of type `ValueInfoProto`
* a list _outputs_ of type `ValueInfoProto`
* a list of "_Initializers_" used to specify constant inputs of the graph of type `TensorProto`
* a list of operations of type `NodeProto`

Despite the naming of the different types, they will all be converted to node (in the sense of vertices) of the graph I am creating.

## Nodes

In ONNX, a node is a special type that holds an operator. In this section, a node will represent a `NodeProto` and I will use the term vertex to describe its implementation in the graph.
Again from the documentation, we read that: 

> Computation graphs are made up of a DAG of nodes, which represent what is commonly called a "layer" or "pipeline stage" in machine learning frameworks.

_Note_: This documentation is present "as-is" in the [GoDoc](https://godoc.org/github.com/owulveryck/onnx-go#NodeProto) and has been auto-generated from the (_protobuf_) definition; it's one of the reasons why I said the _protobuf_ is more flexible than JSON for writing API contracts.

A node knows its inputs and its output. Therefore, to generate the graph, we should:

* start by adding the inputs (which are particular vertices with an indegree of 0). For commodity, we will track the added node into a “dictionary" of vertices (a Go map). 
* add every single node reachable only from vertices presents in the dictionary. 
* add the edges 
* add the current vertex to the dictionary. 

Note: In ONNX, the NodeProto has a type and a name. The name is used as input for its successors. The type is representing the actual mathematical operator that will be applied to the inputs; we will see that later. 

## Gonum

To test and evaluate the structure in the Go environment, let's create a simple graph with the help of the Gonum's Graph package.
I will keep it simple and use the "[simple](https://godoc.org/gonum.org/v1/gonum/graph/simple#DirectedGraph)" implementation.

First, let's define a wrapper struct:

{{< highlight go >}}
type computationGraph struct {
        db      map[string]*node
        digraph *simple.DirectedGraph
}
{{</ highlight >}}

where `db` is the dictionary of nodes (vertices) as described in the previous section.

### Node 

Let's then define a simple `node` (vertex) structure that will fulfill the [Node](https://godoc.org/gonum.org/v1/gonum/graph#Node) interface.
The structure will handle various information later, but for now, let's start with its name and the operation type:

{{< highlight go >}}
type node struct {
        id        int64
        Name      string
        Operation string
}

func (n *node) ID() int64 {
        return n.id
}
{{</ highlight >}}

### Building the DAG

Let's define a wrapper struct. This will give us the flexibility to add (at least) a method to parse the graph later; this will ease the work when we switch to Gorgonia for computing the graph.

{{< highlight go >}}
type computationGraph struct {
        db      map[string]*node
        digraph *simple.DirectedGraph
}
{{</ highlight >}}

To parse the graph we will process the _Initializers_, the _Inputs_ and the _Nodes_ (let's forget the outputs for now).

{{< highlight go >}}
for _, tensorProto := range gx.Initializer {
        n := &node{
              id:    g.digraph.NewNode().ID(),
              Name:  tensorProto.GetName(),
        }
        g.digraph.AddNode(n)
        g.db[name] = n

}
for _, valueInfo := range gx.Input {
        n := &node{
              id:    g.digraph.NewNode().ID(),
              Name:  valueInfo.GetName(),
        }
        g.digraph.AddNode(n)
        g.db[name] = n
}
{{</ highlight >}}

Now a bit more tricky, let's add the Operators and the edges of the graph.
The nodes (`NodeProto`) are supposed to be in topological order in the ONNX model, 
but let's ignore this information and reconstruct the graph as explained before (by reconstructing the topology from the inputs/initializers to the output).

The algo I am using consists in removing items from the node list once it is processed and waiting for the list to be empty. There is a condition that exit the loop, just in case the graph is not a valid DAG.

_Note_: maybe a recursive algorithm would be more efficient, but efficiency is not an issue here.

For clarity, I will not copy the whole code here. Please visit [the GitHub repo](https://github.com/owulveryck/gorgonnx/blob/bfd7eea73340f63b997599b808695401d1ae6f6e/graph.go#L95-L124) for more information.
The critical point is that for each processable node we call a method of the `computationGraph` structure call `processNode`. This method evaluates the content of the node (its inputs and its name), add it to the graph and place the edges the node has with its ancestors (inputs).

## Displaying the result

Thanks to the dot encoding capability of the graph package of Gonum, it is easy to generate an output that is compatible with Graphviz.
By taking back and completing the MNIST example, gluing a little bit and adding particular methods for the `node` object (DOTID,...), we obtain this output:

{{< figure src="/assets/onnx/mnist.png" title="Representation of the MNIST Model" >}}

The graph looks good and is actually representing a convolution neural network. 
This is the end of the first part.

In the [next article](/2018/09/19/my-journey-with-onnx-and-go---running-the-graph.html), let's implement a real backend to be able to compute and evaluate the graph.

# Conclusion

We are now able to read and understand the information encoded into an ONNX model.
The next step is to be able to create a real computation graph that can process the input to produce a result.
That is what we will do in a second post.

So far, with Go, we can write small utilities to extract various pieces of information and represent the models. This is independent of any framework and can be used as a standalone tool. 
