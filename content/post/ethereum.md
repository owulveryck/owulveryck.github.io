---
author: author
date: 2016-08-20T21:45:37+02:00
description: description
draft: true
keywords:
- key
- words
tags:
- one
- two
title: ethereum
topics:
- topic 1
type: post
---


The greeter: https://www.ethereum.org/greeter
# Installation

From sources:

```shell
cd $GOPATH/src/github.com/ethereum
git clone https://github.com/ethereum/go-ethereum.git
```

## installing the solidity compiler

```
sudo add-apt-repository ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get install solc
which solc
```

# Triggering the console and entering the network

## Triggering the console on the test network
`$GOPATH/src/github.com/ethereum/go-ethereum/build/bin/geth --testnet --fast --cache=512 console`



```
> eth.getCompilers()
> I0820 21:59:11.083752 common/compiler/solidity.go:114] solc, the solidity compiler commandline interface
> Version: 0.3.6-0/None-Linux/g++
>
> path: /usr/bin/solc
> ["Solidity"]²
```

The client is downloading the headers, I should wait until 512 Mb is filled:

```
$ du -sh ~/.ethereum 
234M    /home/olivier/.ethereum
```
