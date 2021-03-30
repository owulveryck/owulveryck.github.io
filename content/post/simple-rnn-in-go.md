---
images: ["https://upload.wikimedia.org/wikipedia/en/6/63/Queen_A_Kind_Of_Magic.png"]
description: "You may know how enthusiast I am about machine learning. A while ago I discovered recurrent neural networks. I have read that this 'tool' allow to predict the future! Is this a kind of magic? I have read a lot of stuffs about the 'unreasonable effectiveness' of this mechanism. The litteracy that gives deep explanation exists and is excellent. There is also plehtora of examples, but most of them are using python and a calcul framework. To fully undestand how things work (as I am not a data-scientist), I needed to write my own tool 'from scratch'. This is what this post is about: a more-or-less 'from scratch' implementation of a RNN in go that can be used to applied to a lot of examples"
draft: false
title: "About Recurrent Neural Network, Shakespeare and GO"
date: 2017-10-29T07:17:33+01:00
type: post
---

# Shakespeare and I, encounter of the third type

A couple of months ago, I attended the Google Cloud Next 17 event in London.
Among the talks about SRE, and keynotes, I had the chance to listen to Martin Gorner's excellent introduction: [TensorFlow and Deep Learning without a PhD, Part 2](https://www.youtube.com/watch?v=fTUwdXUFfI8). If you don't want to look at the video, here is a quick summary:

_a 100 of lines of python are reading all Shakespeare's plays; it learns his style, and then generates a brand new play from scratch._ 

Of course, when you are not a data-scientist (and I am not), this looks pretty amazing (and a bit magical).

Back home, I told my friends how amazing it was. I downloaded the code from [github](https://github.com/martin-gorner/tensorflow-rnn-shakespeare), installed tensorflow, and played my Shakespeare to show them.
In essence, here is what they told me:

- _"Amazing, and you know how this works?_ 
- _Well..."_ let's be honest, I had only a vague idea.

It was about something called "Recurrent Neural Networks" (aka _RNN_). 
I dived into the internet... 100 lines of python shouldn't be hard to understand. And to reproduce ?
Where they? 

Actually, it took me months to be able to write this post, without any previous knowledge, it was not that easy.

This is the reason why I finally wrote this article. I want to be sure that I have understood the structure and the possibilities offered by recurrent neural networks.
I also wanted to see whether building a RNN powered tool was doable easily.

This document is divided into two parts:

* the first part is about recurrent neural networks in general;
* the second part is about a toy I made in GO to play with RNNs.

The goal of this text is not to talk about the mathematics behind the neural networks.
Of course, I may talk about vectors, but I will not talk about non-linearity or hyperbolic functions. 

I hope you will be enthusiastic, as much as I am. 

Anyway, do not hesitate to give me any feedback or suggestion that may improve my work.

# First part: The RNN and I, first episode of a time-serie 

The Web is full of resources about machine learning. You can easily find great articles, very well illustrated about neural networks.
I've read a lot...

The more I learn, the more excited I get. 

For example, I've discovered that RNN could, by nature, predict time series
(cf [how to do time series prediction using RNNs, Tensorflow and Cloud ML engine](http://dataconomy.com/2017/05/how-to-do-time-series-prediction-using-rnns-tensorflow-and-cloud-ml-engine/)).

- _"Wait, does it mean that it can predict the future?_
- Well, kind of..."

It is still in the area of "supervised learning". 

Therefore, the algorithm learns events. Based on this, the algorithm can predict what may come next; but only if it is something it has already seen. 
Let's take an example. Consider a lottery game (everybody ask me about this):

To win, your ticket's sequence of numbers must be identical to the one that will be chosen, randomly, at the next draw.
If RNN can predict the future, it should, basically, be able to predict it.

The RNN must learn about the sequences to apply its knowledge (and become a fortune-teller). So If every week the draw is made of "1 2 3 4 5 6", the RNN will learn, and tell that the next draw will be: "1 2 3 4 5 6".

Obviously this is useless; now let's consider a more complex sequence:

Week | sequence
-----|---------
1    | 1 2 3 4 5 6
2    | 2 3 4 5 6 1
3    | 3 4 5 6 1 2
4    | 4 5 6 1 2 3
5    | 5 6 1 2 3 4
6    | 6 1 2 3 4 5
7    | 1 2 3 4 5 6
8    | ? ? ? ? ? ?

Question: What will be the winning sequence of week 8? 

"2 3 4 5 6 1". Cool, you are rich! 
How did you do it? You memorized the sequence. 

RNN does exactly the same.

- "So, it **can** predict the next lottery? 
- No, because there is no sequence in the lottery, it is pure randomness".

In other word, there is no "_recurrence_" in the drawing. Therefore, "_recurrent_" neural networks cannot be applied. 
 
Anyway, beside the lottery, a lot of events are, in essence, recurrents.
The point is that the recurrency model is usually not obvious and not that easy to detect. This is the famous "feeling" of the experts. 

For example, you may have already heard:

- "Will the system crash?
- Based of what I see and what I know, I [don't] think so".

or,

- "Will the sales increase on Sunday?
- Regarding the current market situation and on my experience, it may".

This is where a RNN could shine and enhance our professional lives.

In a pure IT context, for example, you have failures "every now-and-then". Even if you don't find the root cause, it could be useful to predict the next failure. 
If you have enough data about the past failures, the RNN could learn the pattern, and tell you when the next failure will occur.

----
<center>
{{< tweet 844561153229541376 >}}
</center>

----


### Experimenting with RNN

I needed a simple tool to do experimentations.
A huge majority of articles about machine learning are using Python and a framework (such as Tensorflow).
To me, it has two major drawbacks:

* I need to fully understand how to use the framework;
* as it is Python related (and I am not fluent in Python), building **and deploying** efficient tools could take some time.

Let's be more specific about the second point. 

I have seen a lot of samples that could do very beautiful stuffs based on fake data.
Playing with everyday data usually implies to rewrite the tool, from scratch... 

Therefore, I have decided to implement a RNN engine from scratch, in GO (I am "fluent" in go, that have save me days of debugging).
The goal is simple: to understand how RNN works. 

_"Whatever is well conceived is clearly said, and the words to say it flow with ease"._

_"Ce que l'on conçoit bien s'énonce clairement, et les mots pour le dire arrivent aisément"._

<p align="right"><i>Nicolas Boileau</i></p>

### The initial example

The following example is basically an adaptation of Andrej Karpathy's post: [The Unreasonable Effectiveness of Recurrent Neural Networks](http://karpathy.github.io/2015/05/21/rnn-effectiveness/). I strongly encourage you to read it. 

Anyway, I will give you a couple of explanations of the principle.
The goal is to write and train a RNN with a certain amount of text data.

Then, once the RNN is trained, we ask the tool to generate a new text based on what it has learned, character by character.

### How does it work?

Consider the "HELLO" example as described in Karpathy's post.
The vocabulary of the example is composed of 4 letters: `H`, `E`, `L` and `O`. 

The principle is to train the RNN in order to make prediction for the next letter.

Therefore, if I give an `H` as input the fully trained RNN, it will return an `E`,

Then, the `E` will become the input, and the output will be an `L`.

This `L` will become the new input. Here is a difficulty: after an `L`, there can be:

* another `L`, 
* or an `O`. 

This is what make RNN suitable for this case: RNN have a memory!
Then, it will most probably choose a second `L`, based, not only on the last letter, but also on the previous `H` and `E` it has seen.

Correctly trained, the RNN should be able to produce an `O`.

### A classification problem

In practice, this is a [classification problem](https://en.wikipedia.org/wiki/Statistical_classification); every letter in the alphabet is a class.

Given a sequence of letters as input, the mechanism should predict which class it belongs to. This class represents the next letter to be displayed.

For example: 

- `h` belongs to class `e`
- `h e` belongs to class `l`
- `h e l` also belongs to classe `l`
- `h e l l` belongs to class `o`

The network will compute, for each class, a probability based on the input and the context.
So, every letter will be assigned a value between 0 and 1 by the algorithm.

If we formalize that in an array, the ideal situation would be:

<html>
<table border=1 align=center>
<tr>
<th>context</th><th>input</th>
<th>Probability of class H</th>
<th>Probability of class E</th>
<th>Probability of class L</th>
<th>Probability of class O</th>
</tr>
<tr><td></td><td>H</td><td>0</td><td>1</td><td>0</td><td>0</td></tr>
<tr><td>H</td><td>E</td><td>0</td><td>0</td><td>1</td><td>0</td></tr>
<tr><td>H e</td><td>L</td><td>0</td><td>0</td><td>1</td><td>0</td></tr>
<tr><td>H E L</td><td>L</td><td>0</td><td>0</td><td>0</td><td>1</td></tr>
</table>
</html>

In practice, we may have something slightly different (this is an example, do not try to interpret the values):

<html>
<table border=1 align=center>
<tr>
<th>context</th><th>input</th>
<th>Probability of class H</th>
<th>Probability of class E</th>
<th>Probability of class L</th>
<th>Probability of class O</th>
</tr>
<tr><td></td><td>H</td><td>0.1</td><td>0.8</td><td>0.05</td><td>0.05</td></tr>
<tr><td>H</td><td>E</td><td>0.1</td><td>0.07</td><td>0.8</td><td>0.03</td></tr>
<tr><td>H E</td><td>L</td><td>0.05</td><td>0.05</td><td>0.5</td><td>0.4</td></tr>
<tr><td>H E L</td><td>L</td><td>0.05</td><td>0.05</td><td>0.4</td><td>0.5</td></tr>
</table>
</html>

We have encoded the output into an array. In mathematics, such array is called a vector.

On the same principle, we can encode the input letters into a _1-of-k_ vector (1 in the cell corresponding to character, 0 elsewhere).

<pre> 
<code>  
    h e l o
h = 1 0 0 0
e = 0 1 0 0
l = 0 0 1 0
o = 0 0 0 1
</code>
</pre>

The purpose of the prediction is to apply a mathematical function to an input vector in order to produce an output vector (composed of probabilities). 
The next character should be chosen accordingly.

The RNN does not know natively an equation able to predict the correct values but a mathematical model (a composition of mathematical function). This model contains a lot of parameters or variables. 

With proper values, those parameters applied to the mathematical model should allow to compute the correct vector.

Finding the correct parameters is called _the training process_. 

The RNN, fed with a lot of data and their expected output, will adjust its internal parameters.

At each step, the difference between the output and the expected result is evaluated; it is call _the loss_. 

The purpose of the adaptation is to reduce the loss at every step.

# Second part: Let's geek

From now on, let's talk about the implementation; feel free to skip this part and jump straight to the conclusion if your are not interested in coding.

I want to create a tool able to generate a Shakespeare play as described in Karpathy's blog post.
His implementation in Python is [here](https://gist.github.com/karpathy/d4dee566867f8291f086); you can find mine [here](https://github.com/owulveryck/min-char-rnn).

**edit** At first, it was a simple transcript from Python to GO, but the tool has been enhanced. It is now a more generic tool able to use RNN as a processing unit. It's pluggable to any code able to encode and decode a sequence of bytes into a vector.

## The rnn package

I have created a separate package for two reasons:

* to fully understand what is related to the RNN;
* to see what is related to the example about character recognition.

For the same reasons, I have tried to keep parameters as private as possible within the objects.

I am using the [`mat64.Dense`](https://godoc.org/github.com/gonum/matrix/mat64) structure to represent the matrices and simple `[]float64` elements for column vectors
(for more info: [Go Slices: usage and internals](https://blog.golang.org/go-slices-usage-and-internals#clices)).

### The RNN object

The RNN structure holds the three matrices representing the weights to be adapted:

* Wxh
* Whh
* Why

On top of that, the RNN stores two "vectors" for the [bias](https://stackoverflow.com/questions/2480650/role-of-bias-in-neural-networks). One for the hidden layer, the other for the output layer.

The hidden vector is not stored within the structure. Only the last hidden vector evaluated in the process of feedforward/backpropagation is stored.

Not storing the hidden vector within the structure allows to use the same "step" function in the sampling process as well as the training process.

### RNN's step

RNN's step method is the proper implementation of the neural network as described by _Karpathy_.
As explained before, the hidden state is not part of the RNN structure. It is an output of the step function:

```go
func (rnn *RNN) step(x, hprev []float64) (y, h []float64) {
	h = tanh(
		add(
			dot(rnn.wxh, x),
			dot(rnn.whh, hprev),
			rnn.bh,
		))
	y = add(
		dot(rnn.why, h),
		rnn.by)
	return
}
```

You see here that the step function of my RNN takes two vectors as input: 

* a vector representing the currently evaluated item (remember, it is the representation of the `H`, `E`, `L` and `O` in the previous example),
* a hidden vector that is the memory of the passed elements.

It returns two vectors:

* the evaluated output in term of vector (again it is the representation of the  `H`, `E`, `L` and `O`), 
* a new and updated hidden vector. 

_Note_ : For clarity, I have declared a couple of math helpers such as `dot`, `tanh` and `add` that are out of the scope of the explanation.

### The Train method

This method is returning two channels and triggers a goroutine that does the job of training.

```go
func (rnn *RNN) Train() (chan<- TrainingSet, chan float64) {
    ...
}
```

The first channel is a feeding channel for the RNN. It receives a `TrainingSet` that is composed of:

* an input vector
* a target vector 

The goroutine will read the channel, and get all the training data.
It will evaluates the input of the training set and use the target to adapt the parameters.

The second channel is a non blocking channel. It is used to transfer the loss evaluated at each pass for information purpose.

### Forward processing

The forward processing takes a batch of inputs (an array of array) and a sequence of outputs.
It runs the step as many times as needed and stores the hidden vectors in a temporary array. Then the values are used for the back propagation.

```go
func (rnn *RNN) forwardPass(xs [][]float64, hprev []float64) (ys, hs [][]float64) {
    ...
}
```

### Back propagation through time

The back propagation is evaluating the gradient. Once the evaluation is done, parameters can be adapted according to the computed gradients. 

### Adapting the parameters via "AdaGrad"

The method used by Karpathy is the Adaptive gradient.
This one needs a memory; therefore I have declared a new object for the adagrad with a simple Apply method.
The `apply` method  of the `adagrad` object takes the neural network as a parameter and the previously evaluated gradients.

Once this process is done, the RNN is trained and usable. 

### Prediction 

I have implemented a `Predict` method that applies the same method. It starts with an empty memory (the hidden vector is zeroed), takes a sample text as input and generate the output without evaluating the gradient nor adapting the parameters.

This RNN implementation is enough to generate the Shakespeare.

## Enhancement of the tool: implementing codecs

In order to work with any character (= any symbol), the best way to _GO_ is to use the concept of [rune](https://blog.golang.org/strings).
The first implementation of the min-char-rnn I made was using this package. It was simply implementing and a couple of functions to 1-of-k encode and decode the rune, one at a time.

It was working as expected, but I was stuck within the character based neural network.

As explained before, the RNN package is working with vectors, and have no knowledge of characters, pictures, bytes or whatever.

So to continue with this level of abstraction, I have declared a codec interface. 

Therefore, the character based example is simply an implementation that fulfills this interface.

### The codec interface

The codec interface describes the required methods any object must implement in order to use the RNN.

It allows any implementation to use the RNN (imagine a log parser, an image encoder/decoder, a webservice [_insert whatever fancy idea here_]...)

The most important methods of the interface are:

```go
Decode([][]float64) io.Reader
Encode(io.Reader) [][]float64
```

Those methods are dealing with arrays of vectors on one side, and with `io.Reader` on the other side.
Therefore, it can use any input type, from a text representation to a data flow over the network (and if you are _gopher_, you know how cool `io.Reader` are!)

The other methods are simply helper functions useful to train the network. I have also chosen to add a post processing method: 

```go
ApplyDist([]float64) []float64
```

This method is a post processing of the output vector. Actually, the returned vector is made of normalized probabilities. 
In a classification mechanism, one element must be chosen. 
Obviously, the algorithm should choose the one with the best probability. 

But, in the case of the char example, we can add some randomness by selecting the output class according to a certain distribution. 
I have implemented a [Bernouilli distribution](https://godoc.org/github.com/gonum/stat/distuv#Categorical) for the char codec. It is selectable by setting ` CHAR_CODEC_CHOICE=soft` in the environment). 

This function also let the possibility to get the raw normalized probabilities by implementing a no-ops func.

_Note_ : This interface should be reworked, because _Pike_ loves the one-function-interfaces, and _Pike_ knows!

### The char implementation of the codec interface

As explained before, the char implementation consists of a couple of methods that reads a file and encode it.
It can also decode a generated output.

It is that simple. It also serve as an example for whatever new codec implementation.

## The main tool

The main tool is just the glue between all the packages. 
It can be used to train the network or to generate an output. The parameters are tweakable via environment variables (actually each package deals with their own environment variables).

you can find all the code on my [gitHub](https://github.com/owulveryck/min-char-rnn). It needs tweaking and deep testing though. 
I have also uploaded a binary distribution and a pre-trained model (I have implemented a backup and restore mechanism in order to use a pre-trained model).

# Conclusion

I have now understood the principle of RNN. It offers me a lot of opportunities for future work.

I would like to develop some sample tools that could be useful in my everyday life.

Why not writing a codec that can parse the log files of an application. The output would be an encoded status of the health of the application (red, orange or green). 

With the correct data about when warnings or failures occurred, it should able to predict the next failure... before it happens.
