---
title: "Streaming the reMarkable 2"
date: 2021-03-30T17:44:10+02:00
draft: false
summary: "This is a simple article that describes the wiring of the tool I made for streaming the content of the remarkable 2 on a computer. From the proc filesystem to the gRPC implementation over HTTP/2 via the certificate generation."
tags: ["remarkable", "grpc", "protobuf", "linux", "go"]
---

I am a happy owner of a [_reMarkable 2_](https://remarkable.com/) tablet. The device is easy to use out-of-the-box. The only thing I am missing is a proper way to stream the content on my laptop to broadcast it while in visio.

Different hack exists to do so, but I wanted something easy to deploy with very few dependencies and configurations. On top of that, I am always looking for projects to code and learn new stuff. Coding a tool to fulfill my need is a perfect way to achieve both goals.

This article explains how the [goMarkableStream](https://github.com/owulveryck/goMarkableStream) tool works.

In this post, you will find:

- some material about the `/proc/` filesystem on Linux;
- gRPC client and server generation from a protobuf definition;
- A pair of embedded certificates for mutual authentication.

## Getting a picture from the tablet

The first thing to figure out is how to get a picture from the reMarkable.

The remarkable is an armv7 based device running a Linux OS. SSH access is provided, so it is pretty easy to log as `root` on the tablet.

 The casual way to grab a picture is by querying the [framebuffer](https://en.wikipedia.org/wiki/Framebuffer). The Linux kernel exposes a [framebuffer device](https://www.kernel.org/doc/Documentation/fb/framebuffer.txt) addressable via a device node (typically `/dev/fb0`). This device aims to provide an abstraction, _so the software doesn't need to know anything about the low-level (hardware register) stuff._

My first attempt failed: querying the device `/dev/fb0` does not work on the reMarkable 2. Brilliant people made some reverse engineering and provided a good explanation on this [website](https://remarkablewiki.com/tech/rm2_framebuffer). In essence:

> The rm2 does not use the embedded epdc (Electronic Paper Display Controller) of the imx7. Instead the e-Ink display is connected directly to the LCD controller. This means all stuff that the epdc would normally do is now done in software...

This means that the framebuffer is not exposed in `/dev/fb0` by the kernel but by software.

To get an image, we need to get the portion of RAM’s address containing the bitmap of the tablet’s image, and we know it is not referenced by the kernel. 

### The address of the framebuffer

To get the global framebuffer’s address in the RAM, we will query a process that knows it already. The main application of the remarkable handling the GUI is called [`xochitl`](https://remarkablewiki.com/tech/xochitl). It is a closed source software; therefore, there is no way to find what we are looking for by modifying the code.

_Note:_ This is not entirely accurate. It is possible to hack the process, but this goes far beyond my skills. See the [remarkable2-framebuffer](https://github.com/ddvk/remarkable2-framebuffer) for more info.

The Linux kernel traces the memory mapping per process and exposes it in the `proc/[pid]/maps` pseudo-file (see [man 5 procfs](https://man7.org/linux/man-pages/man5/procfs.5.html)). 

By analyzing the maps, it appears that the `xochitl` process is virtually mapping the address of the framebuffer to the pseudo-device.

```shell
grep -C1 '/dev/fb0' /proc/$(pidof xochitl)/maps
72086000-72886000 rw-p 00000000 00:00 0
72886000-74044000 rw-s a8100000 00:06 248        /dev/fb0
74044000-747d2000 rw-p 00000000 00:00 0
```

The global framebuffer is therefore located at `0x74044000` in the RAM. The RAM of the process `xochitl` is accessible through a call to `/proc/[pid]/mem` (once again see [man 5 procfs](https://man7.org/linux/man-pages/man5/procfs.5.html)).


Now, how many bytes should we extract?

The resolution of the reMarkable 2 is 1404x1872. Therefore, let's grab 2628288 bytes:

```shell
reMarkable: ~/ echo $((0x74044000))
1946435584
reMarkable: ~/ dd if=/proc/$(pidof xochitl)/mem of=image.raw count=2628288 bs=1 skip=1946435584
2628288+0 records in
2628288+0 records out
reMarkable: ~/ ls -lrth image.raw
-rw-r--r--    1 root     root        2.5M Mar 31 07:43 image.raw
```

### Our first screenshot

Let's fetch the `image.raw` and convert it to a readable format with imagemagick:

```shell
 convert -depth 8 -size 1872x1404+0 gray:image.raw image.png
 ```

 Then, we can display the image that may look like this:
 
 ![hello reMarkable](/assets/remarkable_hello.png)

 ## Building an application

Now that we are able to grab a picture, let's build an application to grab a flow in real-time.

 ### Overall architecture and principle 

The application is working in client/server mode. The server is getting the raw pictures in an infinite loop and serving them on the network.
It is then the responsibility of the client to fetch the raw pictures from the wire and to encode it into a video stream.

A trivial implementation would be to open a network connection on level 4 and use the TCP protocol as a support to the byte stream.
Nevertheless, this would induce some work to set up some delimiters between each frame and handle the bad messages.

Therefore, it is a good idea to embed each picture into a message and to rely on the capabilities of a framework to do proper encoding decoding.

So far, the widest option is to use protocol buffers as it will use a decent typing mechanism while remaining compact and easy to use.

The message represents an image, and is define like this:

```proto
message image {
    int64 width = 1;
    int64 height = 2;
    bytes image_data = 4;
}
```

Treating the flow of the messages to handle a picture one by one is part of a level 7 protocol. Instead of writing our own, let's keep on working with protobuf use [gRPC](https://grpc.io/).
gRPC is a high-performance, open-source universal RPC framework that runs on top of HTTP/2. The network overhead is therefore low, and the communication between the client and the server remains efficient.

Our streaming service will expose a `GetImage` function that will grab the picture from memory and send it on the wire:

```proto
message Input {}

service Stream {
  rpc GetImage(Input) returns (image) {}
}
```

### Implementation

The implementation of both the client and the server is made in Go.

The `protoc` tool generates the skeleton of the streaming service:

```shell
protoc --gofast_out=plugins=grpc:.  defs.proto3
```

Amongst some utility to handle the serialization and deserialization of the protobuf message (see the doc [Image](https://pkg.go.dev/github.com/owulveryck/goMarkableStream@v0.2.1/stream#Image) for more info), the gRPC framework exposes some 

The [`StreamServer`](https://pkg.go.dev/github.com/owulveryck/goMarkableStream@v0.2.1/stream#StreamServer) is an interface. It is now our responsibility to
create a structure that fulfills the interface, and that is actually implementing the `GetImage` mechanism (getting the image from the memory as exposed before)

```go
type StreamServer interface {
	GetImage(context.Context, *Input) (*Image, error)
}
```

Our server is a basic structure handling a couple of elements:

```go
// Server implementation
type Server struct {
	imagePool   sync.Pool
	r           io.ReaderAt
	pointerAddr int64
	runnable    chan struct{}
}
```

The `r` field is a pointer to the `/proc/[pid]/mem` file from where we will read the data. `pointerAddr` is the location of the framebuffer in this file (0x74044000) in our example, and `runnable` is a channel that is used to handle the requests and avoid burning the CPU of the reMarkable (TL;DR: two consecutive calls to `GetImage` will have to wait to be able to consume `runnable` and  a goroutine is putting one event every x millisecond in the runnable queue).

Basically the implementation of the `GetImage` is trivial:

```go
// GetImage input is nil
func (s *Server) GetImage(ctx context.Context, in *Input) (*Image, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-s.runnable:
		img := s.imagePool.Get().(*Image)
		_, err := s.r.ReadAt(img.ImageData, s.pointerAddr)
		if err != nil {
			s.imagePool.Put(img)
			return nil, err
		}
		return img, nil
	}
}
```

The magic is simply to read the bytes, put it in an image and return it to the caller.
Exposing the service is simply instanciating the objects and using the tools build by the gRPC framework:

```go
ln, _ := net.Listen("tcp", ":2000") // open a listener on TCP on ":2000"
s := stream.NewServer(file, addr) // create the stram object
s.Start() // start the gorouting that feeds the `runnable` channel every x ms
grpcServer := grpc.NewServer(grpc.Creds(grpcCreds)) // create the gRPC server
stream.RegisterStreamServer(grpcServer, s) // register our stream object so it is used by our server
grpcServer.Serve(ln); err != nil { // make the server listen on a TCP connection
```

The client simply dial the server and calls the `GetImage` remote procedure in an endless loop:

```go
conn, err := grpc.Dial("localhost:2000") // Dial the server
client := stream.NewStreamClient(conn)

var img image.Gray
for err == nil {
    response, err := client.GetImage(context.Background(), &stream.Input{})
```

Then it encodes the `response` into a JPEG file and adds it to an MJPEG stream.

```go
var img image.Gray
var b bytes.Buffer
img.Pix = response.ImageData
img.Stride = int(response.Width)
img.Rect = image.Rect(0, 0, int(response.Width), int(response.Height))
jpeg.Encode(&b, &img, nil)
mjpegStream.Update(b.Bytes())
```

The creation and exposition of the MJPEG stream is not detailed in this post as it is slightly out of context. Please see the code if you want more info.

### Security

Even if HTTP/2 does not require any encryption (see [here](https://http2.github.io/faq/#does-http2-require-encryption)), a lot of implementation only supports the protocol if used over an encrypted connexion. 
The Go implementation of gRPC requires by default an encryption channel (that can be bypassed with the use of an `Insecure` method, but we all know that is not a good way to Go ;)).

It is, therefore, a good practice to implement this security mechanism that will avoid sniffing of the pictures from the wifi if you use the tool on an untrusted network.

As I do not want anything difficult to maintain, I am generating a self-signed certificate that I am embedding on both the client and the server with the new `embed` command of the Go language.

I also implement a mutual authentication mechanism. Therefore, only a known client can connect to the server.
The certificate is generated per build (via a set of `go:generate` commands). Therefore, if you want to enhance security, it is your responsibility to generate new binaries, and to store them in a safe place, somewhere on your computer (as they contain the certificate).
I agree that it's not the most secure option, but it is good enough for most use cases.

### Generating the certificate

The certificate is generated in pure go code:

- An internal package is in charge of the certificate sorcery (see the [certificate doc](https://pkg.go.dev/github.com/owulveryck/goMarkableStream@v0.2.1/internal/certificate)).
- A simple CLI generates the file (see [the code](https://github.com/owulveryck/goMarkableStream/blob/v0.2.1/certs/cmd/main.go)).
- A `cert` package (see [doc here](https://pkg.go.dev/github.com/owulveryck/goMarkableStream@v0.2.1/certs) exposes a single function `GetCertificateWrapper()` returning a ready-to-use configuration based on the embeded certificate (`*certificate.CertConfigCarrier`).

Wiring the TLS into the gRPC server is straightforward:

1. For the server:

```go
cert, err := certs.GetCertificateWrapper() // Get the certificate configuration with the embeded certificate
grpcCreds := &callInfoAuthenticator{credentials.NewTLS(cert.ServerTLSConf)} // callInfoAuthenticator is fulfiling the interface https://pkg.go.dev/google.golang.org/grpc@v1.36.1/credentials#TransportCredentials and do the validation of the cerficiate of the client
grpcServer := grpc.NewServer(grpc.Creds(grpcCreds)) // creates the server with the validation mechanism
```

2. For the client:

```go
cert, err := certs.GetCertificateWrapper()
grpcCreds := credentials.NewTLS(cert.ClientTLSConf)
// Create a connection with the TLS credentials
conn, err := grpc.Dial(c.ServerAddr, grpc.WithTransportCredentials(grpcCreds), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
//...
```

## That's all folks!

The tool seems to work as expected for most users. At least it is good enough for me. I do not plan to add any fancy features. Do not hesitate to give it a try if you own a tablet:

[https://github.com/owulveryck/goMarkableStream](https://github.com/owulveryck/goMarkableStream)

The repo also contains a `goreleaser` file if you want to build you own release with your own certificates.

Here is a video of the final product:
{{< youtube c4-hJ6xRzg4 >}}