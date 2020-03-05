---
categories:
- coordination language
date: 2017-03-13T20:54:27+01:00
description: "Third article about writing a distributed linda interpreter"
draft: false
images:
- /assets/images/the_stars_look_different.jpg
tags:
- zygomys
- Linda
- golang
- etcd
title: Linda's evalc, a (tuple)space oddity
---

For a change, I will start with a good soundtrack

<iframe src="https://embed.spotify.com/?uri=spotify:track:72Z17vmmeQKAg8bptWvpVG&theme=white" width="220" height="80" frameborder="0" allowtransparency="true"></iframe>

([youtube version](https://www.youtube.com/watch?v=iYYRH4apXDo) for those who are spotify-less)

----
This is my third article about the distributed coordination language Linda.

The final target of the work is to use this coordination mechanism to deploy and maintain applications based on the description of their topology (using, for example, TOSCA as a DSL).

Last time, I introduced a lisp based language (zygomys) as an embedded programing mechanism to describe the business logic.

Today I will explain how I have implemented a new _action_ in the linda language to achieve a new step: to distribute the work among different nodes.

My test scenario remains the "dining of the philosophers".

# Introducing _evalc_

Linda is a coordination language, but the language which is more than 30 years old, has not been designed with the idea of running on multiple hosts.
The basic primitives of the language do not allow remote execution.

What I need is a sort of _eval_ function that would trigger the execution of the evaluation on another host instead of another goroutine.

I do not care about catching the result of the execution as it will be posted to the tuple space.
Indeed, if more coordination between the actors of this RPC is needed, it can be encoded using the in/out mechanism of linda.

Therefore, I have decided to introduce a new primitive called _evalc_ (for eval compute... Yeah I know, I have imagination)

# Implementing _evalc_

The evalc will not trigger a function on a new host.
Instead, each participating host will run a sort of agent (actually a clone of the zygo interpreter) that will watch a certain type of event (tainted with the evalc) and will then execute a function.

The tuple space acts like a communication channel and this implementation is like a kind of [CSP](https://en.wikipedia.org/wiki/Communicating_sequential_processes) which I like a lot.

The _evalc_ will work exactly as its equivalent _eval_. Therefore the function declaration in go will look like this:

{{< highlight go >}}
func (l *Linda) EvalC(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
    ...
    return zygo.SexpNull, nil
}
{{</ highlight >}}

## First attempt

At first I thought I could simply gob/encode the `args` which contains the `SexpFunction`, post it in the tuple space under a prefixed key. Then the worker would read an execute it in a newly created `glisp` env.

That didn't work mainly because the `SexpFunction` does not have any exported fields, therefore I cannot easily encode/decode it.

I though then that I could encode the `datastack` and post it in the tuple space. I could then decode it in the worker.

I asked for some advice to the author of zygomys [Jason E. Aten](https://www.linkedin.com/in/jason-e-aten-ph-d-45a31318) (aka [glycerine](https://github.com/glycerine))

Here is what he told me (thank you Jason btw):

_Evaluating an arbitrary expression remotely will be challenging because an expression can refer to any variable in the environment, and so would theoretically require a copying of the whole environment--the heap as well as the datastack._ 

And of course he is right!
So I will keep the idea of encoding the whole environment and send it to the workers for a later implementation.
It would need to change the zygomys implementation a lot so export and import both stack. That is too much for now.

## Second attempt

What I did as a temporary solution is a lot simpler and not elegant at all: I have posted the function and the variables in the tuple space and then I am evaluating it in a newly created env.

The main problem is that I cannot access to variables and user function defined outside of the scope of the function. But that will do the trick for now.

Regarding the problem of the philosopher, I had to change the definition of `phil` within my lisp code so it do not call `(eat)` and `(think)` functions anymore.

Here is what is posted in the tuple space when the evalc function is called:
{{< highlight lisp >}}
(defn phil [i num] ((begin (begin (printf "%v is thinking\n" i) (sleep 10000) (printf "/%v is thinking\n" i)) (in "room ticket") (printf "%v is in the room\n" i) (in "chopstick" i) (printf "%v took the %v's chopstick\n" i i) (in "chopstick" (mod (+ i 1) num)) (printf "%v took the %v's chopstick\n" i (mod (+ i 1) num)) (begin (printf "%v is eating\n" i) (sleep 10000) (printf "/%v is eating\n" i)) (printf "%v released the %v's chopstick\n" i i) (out "chopstick" i) (printf "%v released the %v's chopstick\n" i (mod (+ i 1) num)) (out "chopstick" (mod (+ i 1) num)) (printf "%v left the room\n" i) (out "room ticket") (phil i num)))) 0 5
{{</ highlight >}}

In the worker process, I am creating a new environment, loading the function (the `defn` part), and constructing an expression to be evaluated by the env. 
This is what the environment evalutates:

{{< highlight lisp >}}
(defn phil [i num] (
   (begin
    (begin
     (printf "%v is thinking\n" i)
     (sleep 10000)
     (printf "/%v is thinking\n" i))
    (in "room ticket")
    (printf "%v is in the room\n" i)
    (in "chopstick" i)
    (printf "%v took the %v's chopstick\n" i i)
    (in "chopstick" (mod (+ i 1) num))
    (printf "%v took the %v's chopstick\n" i (mod (+ i 1) num))
    (begin
     (printf "%v is eating\n" i)
     (sleep 10000)
     (printf "/%v is eating\n" i))
    (printf "%v released the %v's chopstick\n" i i)
    (out "chopstick" i)
    (printf "%v released the %v's chopstick\n" i (mod (+ i 1) num))
    (out "chopstick" (mod (+ i 1) num))
    (printf "%v left the room\n" i)
    (out "room ticket")
    (phil i num))))
(phil 0 5)
{{</ highlight >}}

# Runtime

## Running it locally: one etcd and several workers

To run it locally I need: 

* a local instance of `etcd` 
* 5 workers. 

Each worker will watch for a new event in the tuple space.
Then I can trigger the execution of the logic with a sixth worker that will read the lisp code, and execute it.

Here is a screenshot of the execution
![Runtime screenshot](https://raw.githubusercontent.com/ditrit/go-linda/master/doc/v0.3.png)

# Conclusion

Jason E. Aten also told me about _sigils_ as a way to discriminate the local variables from the variables present in the tuple space.
I haven't worked on it yet, but I think that I will use those _sigils_ to enhance my linda implementation. It can be usefull for the matching of templates and formals.

By now, I have something that is able to run a basic theorical coordination problem.

Now I think that I will go back to the application management task and see how I can encode the TOSCA workflow so it can be used by this mechanism.

Meanwhile, I will try to test this setup on a worldwide cluster (maybe based on CoreOS).

----
Credit:

The illustration has been found [here](https://www.flickr.com/photos/joebehr/23704122254)
