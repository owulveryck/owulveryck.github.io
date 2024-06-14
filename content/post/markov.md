---
author: Olivier Wulveryck
date: 2016-02-29T20:55:01+01:00
description: I am convinced that the description of the topology of an application IS the way to handle it.
  It can be used to deploy the application, to control it, and even to cure it.
  I'm now interrested in the ability of the application to be self-aware.
  In this post, I'm trying to organize my ideas about the Markov Model and how it can be
  applied to the concept I'm describing.
  It's not meant to be a tutorial (at all). The primary goal is to organize my own ideas and learn.
  The secondary goal is, to find any other interseted person who'd like to discuss about this idea.
draft: false
tags:
- R
- Markov model
- IA
- Machine learning
- eigenvectors
title: Is there a Markov model hidden in the choreography?
topics:
- topic 1
type: post

# mathjax: true

# Use KaTeX
# See https://github.com/KaTeX/KaTeX
katex: true

# Use Mmark
# See https://gohugo.io/content-management/formats/#mmark

---

# Introduction

In my last post I introduced the notion of choreography as a way to deploy an manage application.
It could be possible to implement self-healing, elasticity and in a certain extent
self awareness.

To do so, we must not rely on the _certainty_ and the _determinism_ of the automated tasks.
_Mark Burgess_ explains in his book [in search of certainty](http://http://www.amazon.com/gp/product/1491923075/ref=pd_lpo_sbs_dp_ss_1?pf_rd_p=1944687522&pf_rd_s=lpo-top-stripe-1&pf_rd_t=201&pf_rd_i=1492389161&pf_rd_m=ATVPDKIKX0DER&pf_rd_r=1BRFTEAZ2RRQ8M77MZ0C) that none should consider the command and control anymore.

Actually we grew up with the idea that a computer will do whatever we told him to.
The truth is that it simply don't. If that sounds astonishing to you, just consider the famous bug.
A bug is a little insect that will avoid any programmed behaviour to act as it should.

In a lot of wide spread software, we find _if-then-else_ or _try-catch_ statements.
Of course one could argue that the purpose of this conditional executionis is to deal with different scenarii, which is true, but indeed,
the keyword is _try_...

## Back to the choreography

In the choreography principle, the automation is performed by a set of dancer that acts on their own. Actually, the most logical way
to program it, is to let them know about the execution plan, and assume that everything will run as expected.

What I would like to study is simply that deployement without knowing the deployement plan.
The nodes would try to perform the task, and depending on the event they receive, they learn and enhance their probability of success.

### First problem


First, I'm considering a single node $$ A $$ which can be in three states $$\alpha$$, $$\beta$$ and $$\gamma$$.
Let's $$S$$ be the set of states such as $$ S = \left\{\alpha, \beta, \gamma\right\} $$

#### Actually knowing what's expected

The expected execution is: $$ \alpha \mapsto \beta \mapsto \gamma$$

therefore, the transition matrix should be:

$$
P=\begin{pmatrix}
0 & 1 & 0 \\\\
0 & 0 & 1 \\\\
0 & 0 & 0
\end{pmatrix}
$$

Let's represent it with GNU-R (see [this blog post](http://www.r-bloggers.com/getting-started-with-markov-chains/)
for an introduction of markov reprentation with this software)

```r
> library(expm)
> library(markovchain)
> library(diagram)
> library(pracma)
> stateNames <- c("Alpha","Beta","Gamma")
> ExecutionPlan <- matrix(c(0,1,0,0,0,1,0,0,0),nrow=3, byrow=TRUE)
> row.names(ExecutionPlan) <- stateNames; colnames(ExecutionPlan) <- stateNames
> ExecutionPlan
      Alpha Beta Gamma
      Alpha     0    1     0
      Beta      0    0     1
      Gamma     0    0     0
> svg("ExecutionPlan.svg")
> plotmat(ExecutionPlan,pos = c(1,2),
         lwd = 1, box.lwd = 2,
         cex.txt = 0.8,
         box.size = 0.1,
         box.type = "circle",
         box.prop = 0.5,
         box.col = "light yellow",
         arr.length=.1,
         arr.width=.1,
         self.cex = .4,
         self.shifty = -.01,
         self.shiftx = .13,
         main = "")
> dev.off()
```
which is represented by:

![Representation](/assets/images/ExecutionPlan.svg)

#### Knowing part of the plan...


Now let's consider a different scenario. I assume now that the only known hypothesis is that we cannot go
from $$\alpha$$ to $$\gamma$$ and vice-versa, but for the rest, the execution may refer to this transition matrix:

$$
P=\begin{pmatrix}
\frac{1}{2} & \frac{1}{2} & 0 \\\\
\frac{1}{3} & \frac{1}{3} & \frac{1}{3}  \\\\
0 & \frac{1}{2} & \frac{1}{2}
\end{pmatrix}
$$
which is represented this way ![Representation](/assets/images/ExecutionPlan2.svg)

The transition matrix is regular - we can see, for example that $$P^2$$ contains all non nil numbers:

```r
> ExecutionPlan %^% 2
                Alpha     Beta      Gamma
          Alpha 0.4166667 0.4166667 0.1666667
          Beta  0.2777778 0.4444444 0.2777778
          Gamma 0.1666667 0.4166667 0.4166667
```

Therefore, Markov theorem says that:

* as n approaches infinity, $$P^n = S$$ where $$S$$ is a matrix of the form $$[\mathbf{v}, \mathbf{v},...,\mathbf{v}]$$, where $$\mathbf{v}$$ being a constant vector
* let $$X$$ be any state vector, then we have $$\lim_{n\to \infty}P^nX = \mathbf{v}$$ where $$\mathbf{v}$$ is a fixed probability vector (the sum of its entries = 1), all whose entries are positives

So we can look for vector $$\mathbf{v}$$ (also known as the **steady-state vector of the system**) to see if there is a good chance that our _finite state machine_ would converged to the desired state $$\gamma$$.

### Evaluation of the steady-state vector

Now since $$P^{n+1}=P*P^n$$ and that both $$P^{n+1}$$ and $$P^n$$  approach $$S$$, we have $$S=P*S$$.

Note that any column of this matrix equation gives $$P\mathbf{v}=\mathbf{v}$$. Therefore, the steady-state vector of a regular Markov chain with transition matrix $$P$$ is the unique probability vector $$\mathbf{v}$$ satisfying $$P\mathbf{v}=\mathbf{v}$$.

To find the steady state vector, we must solve the equation: $$P\mathbf{v}=\mathbf{v}$$. $$\mathbf{v}$$ is actually an eigenvector for an eigenvalue $$\lambda = 1$$.

_Note from [wikipedia](https://en.wikipedia.org/wiki/Eigenvalues_and_eigenvectors)_

> In linear algebra, an eigenvector or characteristic vector of a square matrix is a vector that does not change its direction under the associated linear transformation.
> In other words: if $$v$$ is a vector that is not zero, then it is an eigenvector of a square matrix $$A$$ if $$Av$$ is a scalar multiple of $$v$$. i
> This condition could be written as the equation: $$ Av = \lambda v$$, where $$\lambda$$ is a scalar known as the eigenvalue or characteristic
> value associated with the eigenvector $$v$$

To compute the eigenvector, we should find the solution to the equation $$det(A-\lambda I)=0$$ where $$I$$ is the identity matrix. Actually
I don't know how to do it anymore, and I will simply use _R_'s _eigen_ function:

```r
> eigen(ExecutionPlan)
$values
[1]  1.0000000  0.5000000 -0.1666667

$vectors
          [,1]          [,2]       [,3]
          [1,] 0.5773503  7.071068e-01  0.5144958
          [2,] 0.5773503  1.107461e-16 -0.6859943
          [3,] 0.5773503 -7.071068e-01  0.5144958

> ExecutionPlan %^% 15
        Alpha      Beta     Gamma
Alpha 0.2857295 0.4285714 0.2856990
Beta  0.2857143 0.4285714 0.2857143
Gamma 0.2856990 0.4285714 0.2857295
```

Wait, it has found 3 eigenvalues, and one of those equals 1 which is coherent.
But the eigen vector is not coherent at all with the evaluation of the matrix at step 15.

According to [stackoverflow](http://stackoverflow.com/questions/14912279/how-to-obtain-right-eigenvectors-of-matrix-in-r)
that's because it computes the _right_ eigenvector and what I need is the _left_ eigenvector.

Here is how to evaluate it.

```r
> lefteigen  <-  function(A){
       return(t(eigen(t(A))$vectors))
}
> lefteigen(ExecutionPlan)
               [,1]          [,2]       [,3]
          [1,] 0.4850713  7.276069e-01  0.4850713
          [2,] 0.7071068 -3.016795e-16 -0.7071068
          [3,] 0.4082483 -8.164966e-01  0.4082483
```

We now have the steady vector : $$\mathbf{v} = \begin{pmatrix}0.48 \\\\ 0.70 \\\\ 0.40\end{pmatrix}$$

which simply means that according to our theory, our finite state machin will most likely end in state $$\beta$$.

### Analysis

What did I learn ?
Not that much actually. I've learned that given a transition matrix (a model) I could easily compute the probability of success.
If I consider the finte state machine as the whole automator of deploiement, given the pobability of failure, I can predict
if it's worth continuing the deploiement or not.

Cool, but far away from my goal: I want a distributed application to learn how to deploy, cure, and take care of itself with a single information:
its topology.

Back to real life, the model I've described in this post could be the observable states of the application (eg: $$\alpha = initial$$,$$\beta = configured$$, $$\gamma=started$$...)

Hence, the states of the components of the application are hidden from the model (and they must remain hidden, as I don't care observing them)

And this is the proper definition of a __hidden markov model (HMM)__.
So yes, there is a Markov model hidden in the choreography!

I shall continue the study and learn how the signals sent from the compenent gives _evidences_ and do influence the Markov Model of my application.

It's a matter of inference, I-maps, Bayesian networks, HMM.... It's about machine learning which is fascinating !



