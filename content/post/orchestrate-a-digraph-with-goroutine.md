---
author: Olivier Wulveryck
date: 2015-12-02T14:24:21Z
description: |
    A directed graph may be represented by its adjacency matrix.
    Considering each vertice as a runable element and any edge as a dependency,
    I will describe a method to "run" the graph in a concurrent way using goalang's goroutine
draft: false
tags:
- golang
- digraph
title: Orchestrate a digraph with goroutine, a concurrent orchestrator
type: post
---

I've read a lot about graph theory recently.
They have changed the world a lot. From the simple representation to Bayesian network via Markov chains, the applications are numerous.

Today I would like to imagine a graph as a workflow of execution. Every node would be considered as runnable. And every  edge would be a dependency.

It is an important framework that may be used to as an orchestrator for any model, and of course I am a lot thinkingabout __TOSCA__

# The use case 
If we consider this very simple graph (example taken from the french wikipedia page)

<img class="img-responsive" src="/assets/images/digraph1.png" alt="digraph example"/> 

its corresponding adjacency matrix is:

<img class="img-responsive" src="/assets/images/matrix1.png" alt="Adjacency matrix"/>

its dimension is 8x8

For the lab, I will consider that each node has to do a simple task which is to wait for a random number of millisecond (such as Rob Pike's _boring_ function, see references)

# Let's GO

## How will it work

Every node will be run in a `goroutine`. That is a point. But how do I deal with concurrency ?

Every single goroutine will be initially launched and then wait for an information.

It will have an input communication channel, and a _conductor_ will feed this channel with enough information for the goroutine to decides whether it should run or not. 
This information is simply the adjacency matrix up-to-date. That means that is a node is done, its value is set to zero.

Every goroutine will then check in the adjacency matrix, whether it has predecessor (that means if the corresponding vector is null, or every digit in column N is 0) and therefore will execute the step or not.

Once the execution of task is over, the goroutine will then feed another channel to tell the conductor that its job is done. and then the conductor will broadcast the information.

### Example

In our example, nodes _3_, _5_, and _7_ do not have any predecessor, so they will be able to run first.

* __(1)__ The conductors feed the nodes with the matrix

<a href="/assets/orchestrate-a-digraph-with-goroutine/digraph_step1.dot"><img class="img-responsive img-thumbnail" src="/assets/images/digraph_step1.png" alt="digraph example"/></a>

* __(2)__ Every node get the data and analyse the matrix

<a href="/assets/orchestrate-a-digraph-with-goroutine/digraph_step2.dot"><img class="img-responsive img-thumbnail" src="/assets/images/digraph_step2.png" alt="digraph example"/> </a>

* __(3)__ Nodes 3, 5 and 7 have no predecessor (their column in the matrix sums to zero): they can run

<a href="/assets/orchestrate-a-digraph-with-goroutine/digraph_step3.dot"><img class="img-responsive img-thumbnail" src="/assets/images/digraph_step3.png" alt="digraph example"/> </a>

* __(4)__ Nodes 3 and 5 are done, they informs the conductor

<a href="/assets/orchestrate-a-digraph-with-goroutine/digraph_step4.dot"><img class="img-responsive img-thumbnail" src="/assets/images/digraph_step4.png" alt="digraph example"/> </a>

* __(5)__ conductor update the matrix. It fills the rows 3 and 5 with zeros (actually rows 4 and 6, because our first node is 0)

<a href="/assets/orchestrate-a-digraph-with-goroutine/digraph_step5.dot"><img class="img-responsive img-thumbnail" src="/assets/images/digraph_step5.png" alt="digraph example"/> </a>

* __(6)__ The conductor feeds the nodes with the matrix

<a href="/assets/orchestrate-a-digraph-with-goroutine/digraph_step6.dot"><img class="img-responsive img-thumbnail" src="/assets/images/digraph_step6.png" alt="digraph example"/> </a>

* __(7)__ The nodes analyse the matrix

<a href="/assets/orchestrate-a-digraph-with-goroutine/digraph_step7.dot"><img class="img-responsive img-thumbnail" src="/assets/images/digraph_step7.png" alt="digraph example"/> </a>

* __(8)__ Node 2 can run...

<a href="/assets/orchestrate-a-digraph-with-goroutine/digraph_step8.dot"><img class="img-responsive img-thumbnail" src="/assets/images/digraph_step8.png" alt="digraph example"/> </a>

## The representation of the use case in go

### Data representation
to keep it simple, I won't use a `list` or a `slice` to represent the matrix, but instead I will rely on the [package mat64](https://godoc.org/github.com/gonum/matrix/mat64).

A slice may be more efficient, but by now it is not an issue. 

On top of that, I may need later to transpose or look for eigenvalues, and this package does implement the correct method to do so.
For clarity of the description, I didn't use a `float64` array to initialize the matrix.

```golang
// Allocate a zeroed array of size 8×8
m := mat64.NewDense(8, 8, nil)
m.Set(0, 1, 1); m.Set(0, 4, 1) // First row
m.Set(1, 6, 1); m.Set(1, 6, 1) // second row
m.Set(3, 2, 1); m.Set(3, 6, 1) // fourth row
m.Set(5, 0, 1); m.Set(5, 1, 1); m.Set(5, 2, 1) // fifth row
m.Set(7, 6, 1) // seventh row
fa := mat64.Formatted(m, mat64.Prefix("    "))
fmt.Printf("\nm = %v\n\n", fa)
```

### The node execution function (_run_)
The node execution is performed by a `run` function that takes two arguments: 

* The ID of the node
* The duration of the sleep it performs...

This function returns a channel that will be used to exchange a `Message`
```golang
func run(id int, duration time.Duration) <-chan Message { }
```

A `Message` is a structure that will holds:

* the id of the node who have issued the message
* a boolean which act as a flag that says whether it has already run
* a wait channel which take a matrix as argument. This channel acts as the communication back mechanism from the conductor to the node


```golang
type Message struct {
	id   int
	run  bool
	wait chan mat64.Dense
}
```

The run function will launch a goroutine which will remain active thanks to a loop.
It allows the run function to finish an returns the channel as soon as possible to it can be used by the conductor.

### The conductor

The conductor will be executed inside the main function in our example.

The first step is to launch as many `run` function as needed.

There is no need to launch them in separate goroutines, because, as explained before, 
the run function will returns the channel immediately because the intelligence is living in a goroutine already.

```golang
for i := 0; i < n; i++ { // n is the dimension of the matrix
    cs[i] = run(i, time.Duration(rand.Intn(1e3))*time.Millisecond)
    ...
```

Then, as we have launched our workers, and as the communication channel exists, we should launch `n` "angel" goroutines, that will take care of
sending back the matrix to all the workers.

```golang
    ...
	node := <-cs[i]
	go func() {
		for {
			node.wait <- *m
		}
	}()
}
```

Then we shall collect all the messages sent back by the goroutines to treat them and update the matrix as soon as a goroutine has finished.
I will use the `fanIn` function as described by _Rob Pike_ in the IO Takl of 2012 (see references) and then go in a `for loop` to get the results
as soon as they arrived:

```golang
c := fanIn(cs...)
timeout := time.After(5 * time.Second)
for {
    select {
    case node := <-c:
        if node.run == true {
            fmt.Printf("%v has finished\n", node.id)
            // 0 its row in the matrix
            for c := 0; c < n; c++ {
                m.Set(node.id, c, 0)
            }
        }
    case <-timeout:
        fmt.Println("Timeout")
        return
    default:
        if mat64.Sum(m) == 0 {
            fmt.Println("All done!")
            return
        }
    }
}
fmt.Println("This is the end!")
```

__Note__ I have set up a timeout, just in case ([reference](https://talks.golang.org/2012/concurrency.slide#36))...
__Note2__ I do not talk about the fanIn funtion which is described [here](https://talks.golang.org/2012/concurrency.slide#28)

## The test

Here is what I got when I launch the test:

```
go run orchestrator.go 
I am 7, and I am running
I am 3, and I am running
I am 5, and I am running
3 has finished
5 has finished
I am 2, and I am running
I am 0, and I am running
0 has finished
I am 1, and I am running
I am 4, and I am running
4 has finished
7 has finished
2 has finished
1 has finished
All done!
```

Pretty cool

The complete source can be found [here](https://github.com/owulveryck/gorchestrator).

If you want to play: download go, setup a directory and a `$GOPATH` then simply

```
go get github.com/owulveryck/gorchestrator
cd $GOPATH/src/github.com/owulveryck/gorchestrator
go run orchestrator.go
```

# Conclusions

I'm really happy about this implementation. It is clear and concise, and no too far to be idiomatic go.

What I would like to do now:

* Read a TOSCA file (again) and pass the adjacency matrix to the orchestrator. That would do a complete orchestrator for cheap.
* Re-use an old implemenation of the [toscaviewer](https://github.com/owulveryck/toscaviewer).
The idea is to implement a web server that serves the matrix as a json stream. This json will be used to update the SVG (via jquery),
and then we would be able to see the progession in a graphical way.


__STAY TUNED!!!__

### References

* [Go Concurrency Patterns (Rob Pike)](https://talks.golang.org/2012/concurrency.slide)
