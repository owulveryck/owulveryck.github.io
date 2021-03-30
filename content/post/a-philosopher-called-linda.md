---
categories:
- Ideas
date: 2017-02-03T20:57:30+01:00
description: "A non geek article (for now) to talk about the dinner of the philosophers in the cloud"
draft: false
images:
- /assets/images/default-post.png
tags:
- tupleSpace
- topology
- linda
title: Linda, 31yo, with 5 starving philosophers...
---

> __The hand is the tool of tools__ - _Aristotle_.

It ain't no secret to anyone actually knowing me: I am a fan of automation. Automation and configuration management 
have come a long way since [Mark Burgess](http://markburgess.org/) wrote the first version of [cfengine](https://cfengine.com/).

But even if the landscape has changed, operators are still scripting (only the DSL has changed), and the area targeted by those scripts remains technical.

There is no real abstraction nor automation of a design. 

Let me explain that:

* You still need a human to read and understand the architecture of an application. 
* You still need another human to transpile it into a language understandable by a CM tool.
* And you need to configure/script this tool to react on some events to keep the application running and healthy. 

*Note* With a bunch of IT specialists from different major companies, we are trying to figure out the best way to achieve this goal. I will have the opportunity to talk about that in a dedicated post soon.

To describe an application I have had the opportunity to work with [TOSCA](http://docs.oasis-open.org/TOSCA/TOSCA-Simple-Profile-YAML/v1.1/cs0prd01/TOSCA-Simple-Profile-YAML-v1.1-csprd01.html) for a major bank last year (by the way, if you want to play with TOSCA, you can use my [TOSCAlib by Cisco](https://github.com/CiscoCloud/TOSCAlib).

I have really liked the idea of an independent DSL that was able to fully describe an application in a way that it can be writable and understandable by a human as well as a machine.

But it is not enough. TOSCA is based on the idea that you need an orchestrator to operate the workflow. And orchestrator is "bad". The management system must be distributed and autonomous.
(for more about that cf [Configuration management, choreography and self-aware applications](https://blog.owulveryck.info/2016/02/10/configuration-management-choreography-and-self-aware-applications/index.html)

This leads to the idea that the application is a community of elements. And every single element of the community will act regarding the information it gets from the environments and from its peers.

> __Don't communicate by sharing memory; share memory by communicating.__ - _R. Pike_

How can those elements share the information?

# [Tuple Spaces (or, Good Ideas Don't Always Win)](https://software-carpentry.org/blog/2011/03/tuple-spaces-or-good-ideas-dont-always-win.html)

The title of this section is taken from [this blog post](https://software-carpentry.org/blog/2011/03/tuple-spaces-or-good-ideas-dont-always-win.html) which is indeed a good introduction on the tuple-space and how to use them.

## First: What is a tuple

A tuple is simply a finite list of element... the element can be of any type. Therefore a tuple set could be used to describe a lot of things. Because actually we can use a tuple set to describe a vector.
And with several vectors we can describe a matrix, and with matrix...

For example, a digraph can be represented by a tuple set that discribes its adjacency matrix. Therefore, for example, it can then be possible to transpile a TOSCA description to a tuple-set (cf [Orchestrate a digraph with goroutine, a concurrent orchestrator](https://blog.owulveryck.info/2015/12/02/orchestrate-a-digraph-with-goroutine-a-concurrent-orchestrator/index.html) for the decomposition of a TOSCA lifecycle in a matrix).

Now ok, we can describe a workflow... but in a distributed application, how can the node share their states?

## Tuple space...

In short, a tuple space is a repository of tuples that can be accessed concurrently. A tuple space can be seen as a big bucket full of tuple.

The tuple space is visible and consistent through all nodes. The tuple space is the memory!

Ok, so last question: How do we access the tuples? 

# Meet Linda

Linda is a "coordination language" developed by Sudhir Ahuja at AT&T Bell Laboratories in collaboration with David Gelernter and Nicholas Carriero at Yale University in 1986 ([cf wikipedia](https://en.wikipedia.org/wiki/Linda_(coordination_language)))

Linda's principle is very simple as it relies on 4 basic operations:

* _in(t)_ is used to get a tuple from the tuple space if the tuple matches the tuple t. In blocks until a matching tuples exists in the tuple space.
* _rd(t)_ (read) is used to read a tuple from the tuple space if the tuple matches the tuple t. It does not remove it from the tuple space.
* _out(t)_ puts a tuple in the tuple space
* _eval(t)_ is a promise. It evaluates the function contained in a tuple t, immediately returns and will place the result in the tuple space later.

_Important_ A tuple can be __actual__ or __formal__. An actual tuple holds real values. Therefore the _in_ and _rd_ operations on an actual tuple succeed if every single value of the tuple matches.
A formal tuple may holds "variables". Therefore the _in_ and _rd_ operations succeed if the real values match and if the type of the formal match the actual value.

You can find a more complete description of the language and examples [here](http://www.cs.bu.edu/~best/crs/cs551/lectures/lecture-22.html).

# Think big, start small, move fast

Since my colleague [Xavier Talon](https://www.linkedin.com/in/xavier-talon-7bb5261) told me about linda and the idea of using it with TOSCA, I have thousand ideas running around.
What we would like is to use the linda language to coordinate the nodes of an application topology described by TOSCA.
As the topology is  obviously distributed the tuple space I will use/implement must exists at the scale of a cloud platform.

A raft based key/value store could be used as a tuple space. 
And of course the virtual operator that will implement the linda language and interact with the tuple space must be self-contained.
GO would be a good choice for the implementation of the communication agent because of it self-contained, static binary design (maybe RUST would be too but I don't know RUST yet).
Moreover the built-in concurrency could make the development easy (an eval can be triggered simply in a goroutine).

So __let's POC__

First of all First we need to be sure that a distributed tuple-space could work in the cloud.

As a proof of concept, I will use the philosophers dinning problem as simply described in page 452 of the paper [Linda in context](http://www.inf.ed.ac.uk/teaching/courses/ppls/linda.pdf) from Nicholas Carriero and David Gelernter.

My goals are:

* To implement a basic Linda language in go
* To run the philosopher problem locally
* To modify the code so it uses etcd as a tuple space
* To run the philosopher problem on AWS with a philosopher per region
* To use my TOSCAlib to read a topology and encode it in the tuple space
* To run a deployment at scale...

In this post I will present a basic implementation of the language that solves the dinning problem locally.

## The problem

Here is the problem as exposed in the paper:

_A round table is set with some number of plates (traditionally five); there is a single chopstick between each two plates, and a bowl of rice in the center of the table. Philosophers think, then enter the room, eat, leave the room and repeat the cycle. A philosopher can eat without two chopsticks in hand; the two he needs are the ones to the left and the right of the plate at which he is seated.  If the table is full and all philosophers simultaneously grab their left chopsticks, no right chopsticks are available and deadlock ensues. To prevent deadlock, we allow only four philosophers (or one less than the total number of plates) into the room at any one time._ 

## The implementation

I have extracted the C-Linda implementation of this problem and copied it here. 

#### The C linda implenentation
{{< highlight c >}}
Phil(i)
  int i;
{
    while(1) {
      think();
      in("room ticket");
      in("chopstick", i) ;
      in("chopstick", (i+l)%Num) ;
      eat();
      out("chopstick", i);
      out("chopstick", (i+i)%Num);
      out("room ticket");
    }
}
{{</ highlight >}}


{{< highlight c >}}
initialize()
{
  int i;
  for (i = 0; i < Hum; i++) C
    out("chopstick", i);
    eval(phil(i));
    if (i < (Num-1)) 
      out("room ticket");
  }
}
{{</ highlight >}}

### What is needed

To solve this particular problem I don't have to fully implement the linda language. There is no need for the _rd_ action. _eval_ is simply a fork that I will implement using a goroutine and _in_ and _out_ do not use formal tuples.

The actions will communicate with the tuple space via `channels`. Therefore I can create a type Linda composed of two channels for input and output. The actions will be methods of the Linda type.
both _in_ and _rd_ method will get all the tuples in a loop and decide to put them back in the space or to keep it.

{{< highlight go >}}
type Linda struct {
  Input  <-chan interface{}
  Output chan<- interface{}
}
{{</ highlight  >}}

#### The _Tuple_ type
As a tuple I will use a flat go structure. Therefore I can describe a tuple as an interface{}

{{< highlight go >}}
type Tuple interface{}
{{</ highlight  >}}

#### The _in_ action

In will read from the input channel until an object matching its argument is present. If the object read is different, It is sent back in the tuple space via the output channels:

{{< highlight go >}}
func (l *Linda) In(m Tuple) {
  for t := range l.Input {
      if match(m, t) {
        // Assign t to m
        m = t
        return
      }
      // Not for me, put the tuple back
      l.Output <- m
  }
}
{{</ highlight  >}}

### The _eval_ function

The eval function is a bit trickier because we cannot simply pass the function as it would be evaluated before the substitution of the arguments.
What I will do is to pass an array of interface{}. The first argument will hold the function as a first class citizen and the other elements are the arguments of the function.
I will use the reflection to be sure that the argument is a function and executes it in a go routine.

{{< highlight go >}}
func (l *Linda) Eval(fns []interface{}) {
	// The first argument of eval should be the function
	if reflect.ValueOf(fns[0]).Kind() == reflect.Func {
		fn := reflect.ValueOf(fns[0])
		var args []reflect.Value
		for i := 1; i < len(fns); i++ {
			args = append(args, reflect.ValueOf(fns[i]))
		}
		go fn.Call(args)
	}
}
{{</ highlight  >}}

## Back to the philosophers...

#### The Go-linda implementation
Regarding the implementation of Linda, the transcription of the algorithm is simple:

{{< highlight go >}}
for i := 0; i < num; i++ {
    ld.Out(chopstick(i))
    ld.Eval([]interface{}{phil, i})
    if i < (num - 1) {
        ld.Out(ticket{})
    }
}
{{</ highlight >}}


{{< highlight go >}}
func phil(i int) {
    p := philosopher{i}
    fmt.Printf("Philosopher %v is born\n", p.ID)
    for {
        p.think()
        fmt.Printf("[%v] is hungry\n", p.ID)
        ld.In(ticket{})
        ld.In(chopstick(i))
        ld.In(chopstick((i + 1) % num))
        p.eat()
        ld.Out(chopstick(i))
        ld.Out(chopstick((i + 1) % num))
        ld.Out(ticket{})
    }
}
{{</ highlight >}}

### The tuple space
We have Linda... that can put and read tuples via channels... But we still need to plug those channels to the tuple space.
As a first example, we won't store the information and simply pass them from output to input in an endless loop.

{{< highlight go >}}
go func() {
    for i := range output {
        input <- i
    }
}()
{{</ highlight >}}

## Execution

After compiling and executing the code, I can see my philosophers are eating and thinking... 
<pre>
Philosopher 1 is born
[1] is thinking
Philosopher 0 is born
[0] is thinking
Philosopher 3 is born
[3] is thinking
Philosopher 2 is born
[2] is thinking
Philosopher 4 is born
[4] is thinking
[2] has finished thinking
[2] is hungry
[2] is eating
[1] has finished thinking
[1] is hungry
[1] is eating
[4] has finished thinking
[4] is hungry
[4] is eating
...
</pre>
The code can be found here [github.com/owulveryck/go-linda](https://github.com/owulveryck/go-linda/releases/tag/v0.1)

# Conclusion

This is a very basic implementation of the first step. 

In my next experiment, I will try to plug etcd as a tuple space so the philosophers could be distributed around the world.
