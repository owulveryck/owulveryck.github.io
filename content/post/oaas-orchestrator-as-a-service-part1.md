---
author: Olivier Wulveryck
date: 2015-12-07T08:48:20Z
description: |
    In a previous post I have described (an implemented) and orchestrator that takes a digraph
    as input (via its adjacency matrix).
    In this post I will implement an API, so the orchestrator will be transformed into a web service.
    Say hello to OaaS
draft: true
tags:
- golang
- orchestrator
title: OaaS orchestrator as a service - part 1
topics:
- OaaS
type: post
---

In a [previous post](http://blog.owulveryck.info/2015/12/02/orchestrate-a-digraph-with-goroutine-a-concurrent-orchestrator/) I have setup and orchestrator that takes a digraph
as input (via its adjacency matrix).

In this post I will implement an API, so the orchestrator will be transformed into a web service.

# The documentation of the API

The API will be self documented with swagger (see [this post](http://blog.owulveryck.info/2015/11/11/simple-iaas-api-documentation-with-swagger/) for an "introduction").

## The verbs 

I will use only 3 verbs by now:

* POST : to send the digraph representation to the orchestrator engine
* GET : to get the actual status of the execution workflow
* DELETE : delete to obviously delete an execution task

## The JSON objects






## The Swagger interface


# The implementation


