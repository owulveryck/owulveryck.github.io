---
title: "Think like a vertex: using Go's concurrency for graph computation"
date: 2019-10-14T22:26:42+02:00
lastmod: 2019-10-14T22:26:42+02:00
draft: false
keywords: []
summary: "In this article, I describe the prototype of a new computation machine for graph processing.  This machine takes its inspiration from the Pregel paradigm and uses Go's concurrency mechanism as a lever for a simple implementation."
tags: ["golang", "gorgonia", "pregel", "graph", "concurrency", "machine learning"]
categories: []
author: ""

# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: true
toc: false
autoCollapseToc: false
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: true
---

<!--more-->

I am often asked this question:

_Why do you do machine learning in Go?_

Of course, the main reason is that I like the language. But there are other, more generic, reasons.

In the fifth episode of the third season of [Command Line Heroes](https://www.redhat.com/en/command-line-heroes/season-3/the-infrastructure-effect),
Saron Yitbarek exposes the fact that Go's design is tidily linked to the cloud infrastructure. Indeed, the concurrency mechanism makes it super easy to write a program the
can run at scale on inexpensive machines.

And I genuinely believe that this power is underused in the data-science community.
But those are just thoughts. Only facts lead to a dispassionate debate.

In this article, I describe the prototype of a new computation machine for graph processing.
This machine takes its inspiration from the Pregel paradigm and uses Go's concurrency mechanism as a lever for a simple implementation.

# Computation model for graph processing

Efficient graph processing is key to modern computing and machine learning success.

In graph processing, [_Spark_](https://spark.apache.org/) is a reference. It is known for its ability to process large computation graphs.

Spark is base on the [Pregel paradigm](https://kowshik.github.io/JPregel/pregel_paper.pdf), which is a system for large graph processing.
Let's take a closer look at this piece of art involved behind the scene.

### About Pregel

Pregel is a computation model, not an algorithm. From the original paper, Pregel's goal is to define a program as

> (...) a sequence of iterations, in each of which a vertex can
> receive messages sent in the previous iteration, send messages to other vertices, and modify its state and that of
> its outgoing edges or mutate graph topology.
> This vertex-centric approach is flexible enough to express a broad set of algorithms

On top of this definition, those [notes from the CME 323](https://stanford.edu/~rezab/classes/cme323/S15/notes/lec8.pdf)
from Stanford University gives a useful résumé of what Pregel is:

> Pregel is essentially a message-passing interface constrained to the edges of a graph. The idea
> is to ”think like a vertex” - algorithms within the Pregel framework are algorithms in which the
> computation of state for a given node depends only on the states of its neighbors.

### Pregel in Go?

So, Pregel's goal is to solve the problem of graph processing by leveraging the power used for cloud computing.
(a cluster of inexpensive machines).

Distributed programming is, most of the time, efficient, nevertheless hard. But the original paper mention that:

> The model (...) implied synchronicity makes reasoning about programs easier.

And Go's concurrency mechanism makes it (super) easy to synchronize concurrent routines.

Let's draw an elementary computation graph.

Consider this equation (which is a patterned layer of a neural network):

$$f(X) = \sigma(W \cdot X+b)$$

Let's turn it into something more "functional":
$$f(X) = \sigma(add(mul(W,X),B))$$

Now we can express it into a graph:

<center>
![graph](/assets/pregel/graph2.png)
</center>

Now, let's think _like a vertex_:

- I am the `mul` node:
    - I am waiting for `X` and `W` to tell me their values through channel `A` and `B`
    - I am computing the value
    - I am sending the result on channel `C`
- I am the `add` node:
    - I am waiting for `mul` and `b` to tell me their values through channel `C` and `D`
    - I am computing the value
    - I am sending the result on channel `E`
...

### Trivial implementation
Implementing this in Go is fairly easy.

Consider the message as a `float64` value that flows through channels of communication.
A vertex is then a function that reads from the channels apply its content body and writes its result to the output channel:

```go
type message float64
type Vertex func(output chan<- message, input ...<-chan message)
```

The vertex implementation is straightforward:

```go
add := func(output chan<- message, input ...<-chan message) {
        a := <-input[0]
        b := <-input[1]
        output <- message(a + b)
}
```

And the original equation is encoded like this:
```go
A <- message(1.0)
B <- message(1.0)
D <- message(1.0)

mul(C, A, B)
add(E, C, D)
sigma(output, E)

fmt.Println(<-output)
```

Running this code prints `0.8807970779778823` (see the full code [here](https://gist.github.com/owulveryck/b1255077d7e1d940f9cc472bc69ef733))

### Adding concurrency

The main issue in the trivial implementation is that it is not possible to set the values after calling operators. Doing this would lead to a deadlock:

```go
mul(C, A, B)
A <- message(1.0)
```

because assignation to `A` happens after `mul`'s execution, but `mul` is waiting for a value in channel `A`.

A solution to this problem is to use go-routines:

```go
go mul(C, A, B)
go add(E, C, D)
go sigma(output, E)

A <- message(1.0)
B <- message(1.0)
D <- message(1.0)
```

Now every vertex runs in a goroutine, and the deadlock's gone. Even better, this mechanism has implicit synchronization.
Therefore, computing this graph is more efficient with this mechanism that coding it sequentially because the `mul` operations are computed ~~in parallel~~ concurrently:

<center>
![graph](/assets/pregel/graph3.png)
</center>

Let's do a simple benchmark to validate this hypothesis.
I wrote two simple bench functions:

* concurrent
* sequential

(I do not copy the code for clarity, but you can find it [on gist](https://gist.github.com/owulveryck/b1255077d7e1d940f9cc472bc69ef733#file-pregel_test-go)
or run the test directly in the [go playground](https://play.golang.org/p/0nwnzxybl71))


```text
benchcmp /tmp/sequential /tmp/concurrent
benchmark           old ns/op     new ns/op     delta
BenchmarkTest-4     693           4253          +513.71%
```

This result shows that the sequential implementation is way faster. This is understandable because the operation is trivial and highly optimized, and there is an overhead
induced by the concurrency mechanism.

Let's put bias and add a 50 microseconds sleep the `mul` and `add` operations to simulate a real-world and less trivial computation.

```go
mul = func(output chan<- message, input ...<-chan message) {
        a := <-input[0]
        b := <-input[1]
        time.Sleep(50 * time.Microsecond)
        output <- message(a * b)
}
```

With a ballast of 50 microseconds, the concurrent implementation is 22% faster.

```text
benchcmp /tmp/sequential50 /tmp/concurrent50
benchmark           old ns/op     new ns/op     delta
BenchmarkTest-4     276892        213892        -22.75%
```

Let's see how to apply this to machine learning.

# About Gorgonia

Gorgonia is a computation library written in Go.
Its goal is to facilitate machine learning in this language.

Basically, Gorgonia:

* gives the primitives to build an expression graph;
* implements "machines" to compute the expression graph and provide the result;
* does automatic differentiation, but let's put this aside for now;

## The Expression Graph

The main concept of Gorgonia is the Expression Graph (ExprGraph).
The vertices of this graph are Go structures called [`Node`](https://godoc.org/gorgonia.org/gorgonia#Node).

A node carries a [`Value`](https://godoc.org/gorgonia.org/gorgonia#Value) and an [`Op`eration](https://godoc.org/gorgonia.org/gorgonia#Op).

The Operation is an object with a particular method named Do:

```go
Do(...Value) (Value, error)
```

Some vertices are therefore directly runnable thanks to their method.

Equations, Graphs, Vertices, all of the concepts are present in Gorgonia. Therefore, it is possible to write a computation engine on the principle we've evaluated before:

* create a channel of `Value` for every edge of the graph.
* create a goroutine for every node of the graph
  * The goroutine takes the input from the channels that reach the node
  * The goroutine executes the `Do` statement of the operation
  * The goroutine write the output value to every channel issued from the current node

## How? Gorgonia's VM

### What is a VM in Gorgonia?

This paragraph explains the basic principle Gorgonia is using to compute a graph. Feel free to skip it and jump to the conclusion if you are not interested in the internal
implementation of Gorgonia

Gorgonia describes a [VM](https://godoc.org/gorgonia.org/gorgonia#VM) via a Go interface.

From the documentation:

> VM represents a structure that can execute a graph or program.

So, A VM is any Go object with three methods:

* RunAll
* Reset
* Close

_Note_: None of the methods is expecting input. Therefore, the ExprGraph shall be assigned to the VM object. For example:

```go
type myVM struct {
        g *ExprGraph
}

func newMyVM(g *ExprGraph) *myVM {
        return &myVM{
                g: g,
        }
}

func (m *myVM) RunAll() error {
        // walk and process the graph m.g
        // ...
}
```

Another important thing to note is that the machine does not return any output. The `RunAll` method modify every `*Node` of the ExprGraph and sets their Values accordingly.
Once the `RunAll` method returns, the result of the execution is accessible by extracting the value from the root nodes.

### A "Pregel" VM
So, the "Pregel" implementation we are seeking is eventually an implementation of the VM interface.

an Experimental implementation called `GoMachine` lives in the master branch of Gorgonia (under the `x` subdirectory).

The code and godoc are accessible [here](https://godoc.org/gorgonia.org/gorgonia/x/vm#GoMachine). The whole logic is carried by the `RunAll` method
(it uses the structure and the implementation described in this article.).

For clarity, the package uses two helper functions:

```go
func opWorker(n *gorgonia.Node, inputC []<-chan gorgonia.Value, outputC []chan<- gorgonia.Value)

func valueFeeder(n *gorgonia.Node, feedC []chan<- gorgonia.Value)
```

The `*Nodes` carrying an operation are embedded in the function `opWorker` and the leaves of the graph in the function `valueFeeder`.

* opWorker reads the values from the input channels and post results to the output channel
* valueFeeder assumes that the node already owns a value and post this value to the feed channels

As a consequence, the `RunAll`'s job is simply to:

* create a channel per edge of the graph
* launch a goroutine for each node of the graph by embedding it in opWorker or valueFeeder accordingly.

There are some caveats in the trivial implementation described in this post. For example, the vertex cannot read or write the IO channels
sequentially. Otherwise, it can end in a deadlock.
The GoMachine takes care of that. But the principle is not different from what is described here and the code remains simple.

But as of today, I've been able to run this machine with onnx-go, and successfully executed complex models with excellent performances.

Do not hesitate to give it a try and send some feedback.

# Conclusion

Some of the coolest features of the Go language are its simplicity. This goes from the development process to the distribution of the final binary.

This is why I started using Go for machine learning in the first place. It was the easiest way to run a neural net into production without worrying about
dependencies.

But the real power of the Go language goes far beyond those principles.
Concurrency is a crucial point. Using this lever can improve efficiency while keeping things simple.
To me, the concurrent machine is a perfect example of this.

Indeed, the full implementation is less than 200 lines of code and can run graphs such as ResNET.

I do not want to stop the experiment here. Some stuff I'd like to see are:

* The gradient computation to the machine
* The usage of CUDA for the operator supporting it
* ...
* A truly distributed graph computation over several machines, in the cloud, ...

Please let me know what you think about this.
