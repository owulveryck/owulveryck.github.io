---
author: author
date: 2016-01-30T11:35:36+01:00
description: description
draft: true
keywords:
- key
- words
tags:
- one
- two
title: nas for macos on freebsd
topics:
- topic 1
type: post
---

Creating the jail:

```shell
iocage create tag=nas resolver="nameserver 192.168.1.1" hostname="nas" ip4_addr="lo1|192.168.1.100/24"
```

## mDNSresponder
