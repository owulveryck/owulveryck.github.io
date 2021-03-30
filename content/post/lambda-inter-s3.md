---
categories:
date: 2017-01-06T21:59:25+01:00
description: "In this post I will explain how to setup a lambda function written in go that will be triggered when an object is dropped on a bucket on s3. The function will read the object and copy it into another bucket iin another location. It will use AWS' golang sdk."
draft: true
images:
- /assets/images/default-post.png
tags:
- aws
- lambda
- s3
- golang
title: lambda inter s3
---

In this post I will explain how to setup a lambda function written in go that will be triggered when an object is dropped on a bucket on s3. The function will read the object and copy it into another bucket iin another location. It will use AWS' golang sdk.

# The Lambda skeleton

AWS Lambda does not support the go programming language natively. Anyway, it as a go program is compiled and packed as a simgle static binary, it it easy to execute a go program whitin a wrapper (node.js or python).
You can find a lot of tutorials on the Internet to do that. That is definitly not my point here.

For my example, I will use a slighlty different approch, based on a new feature of go 1.8: the plugins. I will use this repo to code my lambda function in order to benefit from the "context" object (so I will be able to handle the timeout gracefully).

https://github.com/eawsy/aws-lambda-go-shim

