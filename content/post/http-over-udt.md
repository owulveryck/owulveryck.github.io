---
author: Olivier Wulveryck
date: 2016-10-17T20:50:18+02:00
description: How UDT (UDP-based Data Transfer Protocol) can be used with HTTP to tranfer files between two AWS instances hosted in different regions.
draft: false
keywords:
- udt
- http
- aws
tags:
- udt
- http
- golang
title: HTTP over UDT for inter-region file transfer
topics:
- UDT
type: post
---

# Introduction

Transferring files between server is no big deal with nowadays network equipments.
You use `rsync`, `scp` or even `http` to get a file from A to B.

Of course, you rely on the TCP stack so you have a decent reliability in the transport.

But TCP has its drawback, especially when it needs to go through a lot of equipments. Typically in the cloud, or over a VPN.

To prevent the drawbacks of the TCP protocol, there are several solutions:

* Use UDP, but UDP by itself is not "reliable"
* Develop another layer 4 protocol, but it cannot be done in a pure user space. You need to develop a system driver. It cannot be easily done on a large scale.
* Use UDP and another framework on top of UDP.

## Layer 4: UDP? 

Yes, UDP, but with an "extra" layer. I've had the opportunity to try three of them.

* Quic by Google
* FASP by ASPERA
* UDT by Dr GU.

### Quic


First Google, along with its [quic](https://en.wikipedia.org/wiki/QUIC) protocol, tries to enhance the user experience. Actually, a quic implementation is already present in chrome and within google web servers. I've heard about quic at the [dotGo](https://dotgo.eu); [Lucas Clemente](https://github.com/lucas-clemente) has presented its work in progress of a quic implementation in GO.

I've tried it, but it lacks a client part by now, and the [quic tools](https://www.chromium.org/quic/playing-with-quic) from chromium are far from being usable in a production environment.

### Aspera's FASP

Aspera has its own protocol. It is based on UDP. I've seen it running, and yes, it simply works!
The problem is that it is not open source and a bit expensive.

### The UDT protocol

The UDT protocol is described by ASPERA as its main competitor [here](http://asperasoft.com/fileadmin/media/Asperasoft.com/Resources/White_Papers/fasp_Critical_Technology_Comparison_AsperaWP.pdf).
It's open source and worth the try.
It's the one I will use for my tests.
The code is distributed as a C++ library, but it exists GO bindings.

## The Layer 7: HTTP

To actually transfer a file, I can use the `udtcat` tool provided in the github of go-udtwrapper. 
It is ok for a test, but I won't be able to serve multiple files, to resume a transfer etc... So I need a layer 7 protocol.
HTTP is, according to me, a good choice.

# The implementation in GO

Implementing a simple client/server http-over-udt in go is relatively easy. The HTTP is interfaced in a way that the transport can be easily changed.
Therefore, no need to reimplement a complete HTTP stack; GO has all I need in its standard library.

<center>
![/assets/images/save-princess-go.jpg](/assets/images/save-princess-go.jpg)

https://toggl.com/programming-princess
</center>

I will use this fork of [go-udtwrapper](github.com/Lupus/go-udtwrapper) which seems to be the most up-to-date.

## The server

Implementing a basic http server over UDT is very easy.

The [Serve function](https://golang.org/pkg/net/http/#Serve) from the http package takes a `net.Listener` as argument.
The `udt.Listen` function implements the [net.Listener](https://golang.org/pkg/net/#Listener) interface.

Therefore we can simply use this code to serve HTTP content via the DefaultMuxer over udt:

```go
ln, _ := udt.Listen("udt", config.Addr)
http.Serve(ln, nil)
```

A full implementation that serves local file is simply done by:

{{< gist owulveryck 6a44885c2b3527159f496c21071ab8df "server.go" >}}

## The client

The `http.Client`'s DefautTransport relies on TCP.
Therefore we must completely rewrite a Transport to use UDT.

The Transport entry of the Client implements the RoundTripper interface.

The key point is to write a client transport for UDT that implements the RoundTripper interface.

### The http.RoundTripper interface

Here is an example of an implementation:

```go 
// UdtClient implements the http.RoundTripper interface
type udtClient struct{}

func (c udtClient) RoundTrip(r *http.Request) (*http.Response, error) {
      d := udt.Dialer{}
      conn, err := d.Dial("udt", r.URL.Host)
      if err != nil {
          return nil, err
      }
      err = r.Write(conn)
      if err != nil {
          return nil, err
      }
      return http.ReadResponse(bufio.NewReader(conn), r)
}
```

### The full client code

A simple client that will perform a GET operation on our server would be:

{{< gist owulveryck 6a44885c2b3527159f496c21071ab8df "client.go" >}}

### Building the tools
As we rely on CGO, to do a static compilation, we must use the extra flags: `go build --ldflags '-extldflags "-static"'`.

# Conclusion

This is a very basic implementation of http over UDT.
I have developed a more complete tool for my client, but it cannot be published in open source.

Among the things that I have done there are:

* Gzip compression
* Partial content for resuming a broken download (with http.ServeContent)
* SHA256 checking at the end of the transport
* an HTTP middleware (Rest API) to query the transfer states, rates and efficiency via the PerfMon interface

What's not done yet:

* TLS and mutual authentication
* good benchmarks to actually measure the performances of UDT.
* Downloading chunks to optimize the speed of transfer and the bandwith usage
* maybe a POST method to upload a file in multipart
