---
categories:
- cloud
- distributed systems
date: 2017-02-28T20:57:38+01:00
description: "This is the second part of a series on my attempt to build a deployement language on a cloud scale"
draft: false
images:
- https://upload.wikimedia.org/wikipedia/commons/9/9d/SICP_cover.jpg
tags:
- go
- zygomys
- Lisp
- linda
title: To go and touch Linda's Lisp
---

The title is not a typo nor dyslexia. I will really talk about Lisp.

In a [previous post](/2017/02/03/linda-31yo-with-5-starving-philosophers.../index.html) I explained my will to implement the dining of the philosophers with Linda in GO.

The ultimate goal is to use a distributed and abstract language to go straight from the design to the runtime of an application.

# The problem I've faced

I want to use a GO implementation for the Linda language because a go binary is a container by itself. Therefore if I build my linda language within go, I will be able to run it easily across the computes nodes without any more dependencies.
The point is that the Linda implementation may be seen as a go package. Therefore every implementation of every algorithm must be coded in go. Therefore I will lack a lot of flexibility as I will need one agent per host __and__ per algorithm. For example the binary that will solve the problem of the dining of the philosophers will only be useful for this specific problem.

What would be nice it to use an embedded scripting language. This language would implement the Linda primitives (_in, rd, eval, out_). And the go binary would be in charge to communicate with the tuple space.

## Tuple space: _I want your Sexp_

I have thought a lot about a way to encode my tuples for the tuple space.
Of course go as a lot of encoding available:

- json
- xml
- protobuf
- gob

None of them gave me entire satisfaction. The reason is that go is strongly typed. A tuple must be internally represented by an empty *interface{}* to remain flexible.
Obviously I would need to use a lot of reflexion in my code. And reflexion is not always easy. And a bad implementation leads to an unpredictable code.

To keep it simple (and idiomatic) I took a little refresh about the principles of the reflection. So I took my [book](https://books.google.fr/books/about/The_Go_Programming_Language.html?id=SJHvCgAAQBAJ) about go (I bought it when I started learning the language).

In this book there is a full example about encoding and decoding [s-expression](https://en.wikipedia.org/wiki/S-expression). And what is an s-expression? A tuple! 

__Eureka__!

## Lisp/zygomys

I started working on s-expression... I could have used the parser described in my book and that would be enough for the purpose of my test.
I could have created a package _encoding/sexpr_ and that would do the job.

But the more I was reading about s-expression, the more I was digging in list processing. 

List processing, s-expression, embedded language, functional programming.... That was it: I really needed a lisp based embedded language to implement my linda solution.

![xkcd 297](https://imgs.xkcd.com/comics/lisp_cycles.png)

I found [Zygomys](https://github.com/glycerine/zygomys). This project was a perfect fit for my needs because it seemed stable enough and easily extensible.
The main drawback is that its author decided not to use the godoc format. That is a bit annoying but the documentation exists and is in a wiki. On the other hand the author has replied to all of my solicitations. So I gave it a go.

# The POC

## Implementing the linda primitives in the REPL

The linda primitives are implemented as GO functions. I will more or less use the same structure as the one I have already used in my first attempt.
The go functions will be exported into the repl as documented in the [wiki of zygomys](https://github.com/glycerine/zygomys/wiki/Go-API)

For example the Out function is implemented with this signature:

{{< highlight go >}}
func (l *Linda) Out(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
    ...
}
{{</ highlight >}}

and it will be presented to the repl by this command:

{{< highlight go >}}
lda := &linda.Linda{}
env := zygo.NewGlisp()
env.AddFunction("out", lda.Out)
{{</ highlight >}}

I have decided to let the linda implementation in a separate package and to implement the repl ad a separate command.

## _etcd_ as a tuple space

My trivial implementation of tuple space based on channels was inaccurate. So I need to implement something more robust.
In the future the tuple space will be distributed at the scale of the cloud.

A raft-based key value store is nowadays a good choice. 
I have chosen to play with etcd by now (but I will try consul later on).

For the moment I will run a single instance locally.

### Linda and etcd
The _out_ and _eval_ statements will write the tuple as the value of a key prefixed by something fixed per session and suffixed by a uuid.

The _In_ statement is a bit trickier:

It will read all the tuples prefixed by the constant defined and try to match the tuple passed as argument.
If it succeeds it will try to delete it.
If it succeeds it will return the tuple. This is needed to avoid a race condition on a single tuple.

If no matching tuple is present in the tuple space, the function watch any PUT event. If the value associated to the event matches the arguments, it tries to delete the tuple from the kv store and returns it.

## Implementing the algorithm in _zygomys_

This tasks gave me a lot of pain.
To be honest I started (once more) to read the book [structure and interpretation of computer program](https://mitpress.mit.edu/sicp/full-text/book/book.html).

There is a beautiful and functional way to implement the algorithm in lisp. I am sure about that.

But I will figure it out later.

By now, what I did was simply to transpile the algorithm in the lisp syntax thanks to the `begin` statement of _zygo_ (see the [section Sequencing in the wiki of zygomys](https://github.com/glycerine/zygomys/wiki/Language)).

This is how it looks like:

{{< highlight lisp >}}
(defn phil [i num] (
  (begin
    (think i)
    (in "room ticket")
    (in "chopstick" i)
    (in "chopstick" (mod (+ i 1) num))
    (eat i)
    (out "chopstick" i)
    (out "chopstick" (mod (+ i 1) num))
    (out "room ticket")
    (phil i num))))
{{</ highlight >}}

## Execution

_etcd_ daemon needs to be launched first. And by now it needs a clean database to avoid any side effect (I still have a lot of TODO in my code).

Then launch the linda repl with the example:

<pre>
=> localhost cmd git:(master) # ./cmd ../example/dinner/dinner.zy
Creating chopstick 0
Creating philosopher 0
Creating room ticket
0 is thinking
Creating chopstick 1
Creating philosopher 1
Creating room ticket
1 is thinking
Creating chopstick 2
Creating philosopher 2
Creating room ticket
2 is thinking
Creating chopstick 3
Creating philosopher 3
Creating room ticket
3 is thinking
Creating chopstick 4
Creating philosopher 4
4 is thinking
/4 is thinking
4 is in the room
4 took the 4's chopstick
4 took the 0's chopstick
4 is eating
/0 is thinking
0 is in the room
/3 is thinking
3 is in the room
/2 is thinking
2 is in the room
3 took the 3's chopstick
2 took the 2's chopstick
...
</pre>

# Conclusion and future work

This implementation is a start. The code needs a lot of tweaking though. The linda implementation is still far to be complete.

But by now I won't pass much time in implementing some more feature because I have enough to play with my philosophers.
What I would like to do next is to distribute my philosophers and use the real power of etcd.

I would like to start the REPL on 5 different hosts all linked by etcd.
Then I will inject my lisp code in the tuple space and wait for the philosophers to think and eat.

If you are curious, the code is on [the github of the ditrit project](https://github.com/ditrit/go-linda)

Meanwhile, let's have a drink and relax with a good sound:

<iframe src="https://embed.spotify.com/?uri=spotify:track:4QwzVlAJSkcLeCNQ6Ug30P&theme=white" width="280" height="80" frameborder="0" allowtransparency="true"></iframe>
