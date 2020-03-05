---
author: Olivier Wulveryck
date: 2015-11-13T09:21:30Z
description: Securing the API with JWT
draft: true
keywords:
tags:
- go
- API
- simple_iaas
title: Securing the API (in GO)
type: post
---

I've created a [couple of posts](http://blog.owulveryck.info/tags/simple-iaas/) about creating an API (mainly in GO).
By now, the API is open. Now I will implement a basic authentication and accreditation mechanism with a little help from [JWT](http://jwt.io)

# What is a Json Web Token

According to the [RFC 7519](https://tools.ietf.org/html/rfc7519), it is a

> compact, URL-safe means of representing
> claims to be transferred between two parties.  The claims in a JWT
> are encoded as a JSON object that is used as the payload of a JSON
> Web Signature (JWS) structure or as the plaintext of a JSON Web
> Encryption (JWE) structure, enabling the claims to be digitally
> signed or integrity protected with a Message Authentication Code
> (MAC) and/or encrypted.

It is a [widely spread mechanism](https://www.google.fr/trends/explore#q=json%20web%20token) to access web services.

_Note:_  JWT is not an authentication framework by itself, and in this post I will assume that the authentication is done via a webservice
that will simply reply `true` or `false` to authorize the user to use the API.

# Defining the user and the authentication API


