+++
date = "2015-10-22T20:40:36+02:00"
draft = false
title = "Ruby / ZeroMQ / GO" 
description = "My attempt to make a go program talk to a ruby script via a 0MQ message"
tags = [
    "Development",
    "go",
    "ruby",
    "zmq"
]
+++

# Abtract

I really like go as a programming language. It is a good tool to develop web restful API service.

On the other hand, ruby and its framework rails has also been wildly used to achieve the same goal.

Therefore we may be facing a "legacy" ruby developpement that we would like to connect to our brand new go framework.
0MQ may be a perfect choice for intefacing the two languages.

Anyway, it is, at least, a good experience to do a little bit of code to make them communicate.

# ZeroMQ

I will use the ZeroMQ version 4 as it is the latest available one.
On top of that, I can see in the [release notes](http://zeromq.org/docs:changes-4-0-0) that there is an implementation of a strong encryption, and I may use it later on 
# Go

## Installation of the library

As written in the README file, I try a `go get` installation on my chromebook.
```
~ go get github.com/pebbe/zmq4
# pkg-config --cflags libzmq
Package libzmq was not found in the pkg-config search path.
Perhaps you should add the directory containing `libzmq.pc'
to the PKG_CONFIG_PATH environment variable
No package 'libzmq' found
pkg-config: exit status 1
```

The go binding is not a pure go implementation, and it still needs the C library of zmq.

Let's _brew installing_ it:

```
~  brew install zmq
==> Downloading http://download.zeromq.org/zeromq-4.1.3.tar.gz
######################################################################## 100.0%
==> ./configure --prefix=/usr/local/linuxbrew/Cellar/zeromq/4.1.3 --without-libsodium
==> make
==> make install
/usr/local/linuxbrew/Cellar/zeromq/4.1.3: 63 files, 3.5M, built in 73 seconds
```

Let's do the go-get again:

```
~ go get github.com/pebbe/zmq4
```

so far so good. Now let's test the installation with a "hello world" example.

_Note_: the [examples directory](https://github.com/pebbe/zmq4/blob/master/examples) contains a go implementation of all the example of the ZMQ book
I will use the [hello world client](https://github.com/pebbe/zmq4/blob/master/examples/hwclient.go) and the [hello world server](https://github.com/pebbe/zmq4/blob/master/examples/hwserver.go) for my tests

The hello world client/server is implementing a Request-Reply patternt and are communicating via a TCP socket.

* The server is the *replier* and is listening on the TCP port 5555
```go
...
func main() {
    //  Socket to talk to clients
    responder, _ := zmq.NewSocket(zmq.REP)
    defer responder.Close()
    responder.Bind("tcp://*:5555")
    ...
}
```
* The client is the *requester* and is dialing the same TCP port
```go
...
func main() {
    //  Socket to talk to server
    fmt.Println("Connecting to hello world server...")
    requester, _ := zmq.NewSocket(zmq.REQ)
    defer requester.Close()
    requester.Connect("tcp://localhost:5555")
    ...
}
```

Then, the client is sending (requesting) a _hello_ message, and the server is replying a _world_ message.

## Running the example
First, start the server:

```
~ cd $GOPATH/src/github.com/pebbe/zmq4/examples
~ go run hwserver.go
```

Then the client

```
~ cd $GOPATH/src/github.com/pebbe/zmq4/examples
~ go run hwclient.go
Connecting to hello world server...
Sending  Hello 0
Received  World
Sending  Hello 1
Received  World
Sending  Hello 2
...
```

# Ruby

Now let's implement a Ruby client.

## Installation of the library

a _gem install_ is supposed to do the trick:

```
~ gem install zmq
Building native extensions.  This could take a while...
ERROR:  Error installing zmq:
ERROR: Failed to build gem native extension.

/usr/local/linuxbrew/opt/ruby/bin/ruby -r ./siteconf20151022-23021-1ehwusq.rb extconf.rb
    checking for zmq.h... yes
    checking for zmq_init() in -lzmq... yes
    Cool, I found your zmq install...
    creating Makefile

    make "DESTDIR=" clean

    make "DESTDIR="
    compiling rbzmq.c
    rbzmq.c: In function 'socket_getsockopt':
    rbzmq.c:968:7: error: 'ZMQ_RECOVERY_IVL_MSEC' undeclared (first use in this function)
        case ZMQ_RECOVERY_IVL_MSEC:
        ...
```

Arg!, something went wrong. It looks like there is a version mismatch between th libzmq brew installed and the version expected by the gem
The _zmq_ gem seems a bit old and there is a *FFI* ruby extension with a more active developement.

Moreover, I have found []the perfect website for the ruby-and-zmq-ignorant(https://github.com/andrewvc/learn-ruby-zeromq)

As written in the doc, let's install the needed gems via `gem install ffi ffi-rzmq zmqmachine`

## Let's try the lib

Ok, it is now time to run an example

```
require 'rubygems'
require 'ffi-rzmq'
def error_check(rc)
    if ZMQ::Util.resultcode_ok?(rc)
        false
    else
        STDERR.puts "Operation failed, errno [#{ZMQ::Util.errno}] description [#{ZMQ::Util.error_string}]"
        caller(1).each { |callstack| STDERR.puts(callstack)  }
        true
    end
end

ctx = ZMQ::Context.create(1)
STDERR.puts "Failed to create a Context" unless ctx

req_sock = ctx.socket(ZMQ::REQ)
rc = req_sock.connect('tcp://127.0.0.1:5555')
STDERR.puts "Failed to connect REQ socket" unless ZMQ::Util.resultcode_ok?(rc)

2.times do
    rc = req_sock.send_string('Ruby says Hello')
    break if error_check(rc)

    rep = ''
    rc = req_sock.recv_string(rep)
    break if error_check(rc)
    puts "Received reply '#{rep}'"
end
error_check(req_sock.close)

ctx.terminate
```

Running this example with a simple `ruby client.rb` command leads to the following errors:
```
ruby client.rb
Assertion failed: check () (src/msg.cpp:248)
```

But, my GO server is receiving the messages:

```
~ go run hwserver.go
Received  Ruby says Hello
Received  Ruby says Hello
```

# End of show

That's it for now. I think I'm facing a bug in the ruby implementation of the libzmq I'm using. 
Indeed, I've found an [issue](https://github.com/chuckremes/ffi-rzmq/issues/118)... 

I will check again later, or I will try on another environement but the essential is here.
