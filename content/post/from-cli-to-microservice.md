---
categories:
date: 2017-09-02T13:28:36+02:00
description: "This article explains how to turn a golang utility into a webservice using gRPC (and protobuf). I take the example of Hashicorp tools because they are often used as a leverage for the DevOps transformation. Often, the Ops use the tools for themselves, but when comes the time to provide a service around them, they are usually scared to open the engine. They prefer to make a factory around the service, which is often less reliable than a little piece of code fully tested."
draft: false
images:
- /assets/images/terraformcli.png
tags:
title: From command line tools to microservices - The example of Hashicorp tools (terraform) and gRPC
---

This post is a little different from the last ones. As usual, the introduction tries to be open, but it quickly goes deeper into a go implementation.
Some explanations may be tricky from time to times and therefore not very clear. As usual, do not hesitate to send me any comment via this blog or via twitter [@owulveryck](https://twitter.com/owulveryck).

**TL;DR**: This is a step-by-step example that turns a golang cli utility into a webservice powered by gRPC and protobuf. The code can be found [here](https://github.com/owulveryck/cli-grpc-example).

# About the cli utilities

I come from the sysadmin world... Precisely the Unix world (I have been a BSD user for years). Therefore, I have learned to use and love "_the cli utilities_". Cli utilities are all those tools that make Unix sexy and "user-friendly". 

<center>
Because, yes, Unix **is user-friendly** (it's just picky about its friends[^1]).
</center>

[^1]: This sentence is not from me. I read it once, somewhere, on the Internet. I cannot find anybody to give the credit to.

From a user perspective, cli tools remains a must nowadays because:

* there are usually developed in the pure Unix philosophy: simple enough to use for what they were made for;
* they can be easily wrapped into scripts. Therefore, it is easy to automate cli actions.

The point with cli application is that they are mainly developed for an end-user that we call "an operator". As Unix is a multi-user operating system, several operators can use the same tool, but they have to be logged onto the same host.

In case of a remote execution, it's possible to execute the cli via `ssh`, but dealing with automation, network interruption and resuming starts to be tricky.
For remote and concurrent execution web-services are more suitable.

Let's see if turning a cli tool into a webservice without re-coding the whole logic is easy in go?

## Hashicorp's cli

For the purpose of this post, and because I am using Hashicorp tools at work, I will take [@mitchellh](https://twitter.com/mitchellh)'s framework for developing command line utilities.
This package is used in all of the Hashicorp tools and is called................ "[cli](https://github.com/mitchellh/cli)"! 

This library provides a [`Command`](https://godoc.org/github.com/mitchellh/cli#Command) type that represents any action that the cli will execute.
`Command` is a go `interface` composed of three methods:

* `Help()` that returns a string describing how to use the command;
* `Run(args []string)` that takes an array of string as arguments (all cli parameters of the command) and returns an integer (the exit code);
* `Synopsis()` that returns a string describing what the command is about.
 
 _Note_: I assume that you know what an interface is (especially in go). If you don't, just google, or even better, buy the book [The Go Programming Language](https://www.amazon.com/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440) and read the _chapter 7_ :).

The main object that holds the business logic of the cli package is an implementation of [`Cli`](https://godoc.org/github.com/mitchellh/cli#CLI). 
One of the elements of the Cli structure is `Commands` which is a `map` that takes the name of the action as key. The name passed is a string and is the one that will be used on the command line. The value of the `map` is a function that returns a `Command`. This function is named [`CommandFactory`](https://godoc.org/github.com/mitchellh/cli#CommandFactory). According to the documentation, the factory is needed because _we may need to setup some state on the struct that implements the command itself_. Good idea!

## Example

First, let's create a very simple tool using the "cli" package.
The tool will have two "commands":

* hello: will display _hello args...._  on `stdout` 
* goodbye: will display _goodbye args..._ on `stderr`

{{< highlight go >}}
func main() {
      c := cli.NewCLI("server", "1.0.0")
      c.Args = os.Args[1:]
      c.Commands = map[string]cli.CommandFactory{
            "hello": func() (cli.Command, error) {
                      return &HelloCommand{}, nil
            },
            "goodbye": func() (cli.Command, error) {
                      return &GoodbyeCommand{}, nil
            },
      }
      exitStatus, err := c.Run()
      ... 
}
{{</ highlight >}}
As seen before, the first object created is a `Cli`. Then the `Commands` field is filled with the two commands "hello" and "goodbye" as keys, and an anonymous function that simply returns two structures that will implement the `Command` interface.

Now, let's create the `HelloCommand` structure that will fulfill the [`cli.Command`](https://godoc.org/github.com/mitchellh/cli#Command) interface:

{{< highlight go >}}
type HelloCommand struct{}

func (t *HelloCommand) Help() string {
      return "hello [arg0] [arg1] ... says hello to everyone"
}

func (t *HelloCommand) Run(args []string) int {
      fmt.Println("hello", args)
      return 0
}

func (t *HelloCommand) Synopsis() string {
      return "A sample command that says hello on stdout"
}
{{</ highlight >}}

The `GoodbyeCommand` is similar, and I omit it for brevity.

After a simple `go build`, here is the behavior of our new cli tool:
{{< highlight shell >}}
~ ./server help
Usage: server [--version] [--help] <command> [<args>]

Available commands are:
    goodbye    synopsis...
    hello      A sample command that says hello on stdout

~ ./server hello -help
hello [arg0] [arg1] ... says hello to everyone

~ ./server/server hello a b c
hello [a b c]
{{</ highlight >}}

So far, so good!
Now, let's see if we can turn this into a webservice.

# Micro-services

<center>_The biggest issue in changing a monolith into microservices lies in changing the communication pattern. - Martin Fowler_[^2]</center>

[^2]: from [Martin Fowler's Microservices definition](https://martinfowler.com/articles/microservices.html#SmartEndpointsAndDumbPipes).

There is, according to me, two options to consider turning our application into a webservice:

* a RESTish communication and interface;
* an RPC based communication.

SOAP is not an option anymore because it does not provide any advantage over the REST and RPC methods.

## Rest? 

I've always been a big fan of the REST "protocol". It is easy to understand and to write. On top of that, it is verbose and allows a good description of "business objects".
But, its verbosity, that is a strength, quickly become a weakness when applied to machine-to-machine communication.
The "contract" between the client and the server have to be documented manually (via something like swagger for example). And, as you only transfer objects and states, the server must handle the request, understand it, and apply it to any business logic before returning a result.
Don't get me wrong, REST remains a very good thing. But it is very good when you think about it from the beginning of your conception (and with a user experience in mind).

Indeed, it may not be a good choice for easily turning a cli into a webservice.

## RPC!

RPC, on the other hand, may be a good fit because there would be a very little modification of the code.
Actually, the principle would be to:

1. trigger a network listener
2. receive a _procedure call with arguments_,
3. execute the function
4. send back the result

The function that holds the business logic does not need any change at all.

The drawbacks of RPCs are:

* the development language need a library that supports RPC,
* the client and the server must use the same communication protocol.

Those drawbacks have been addressed by Google. They gave to the community a polyglot RPC implementation called gRPC. 

Let me quote this from the chapter "[The Production Environment at Google, from the Viewpoint of an SRE](https://landing.google.com/sre/book/chapters/production-environment.html#our-software-infrastructure-XQs4iw)" of the SRE book:

> _All of Google's services communicate using a Remote Procedure Call (RPC) infrastructure named Stubby; an open source version, gRPC, is available. Often, an RPC call is made even when a call to a subroutine in the local program needs to be performed. This makes it easier to refactor the call into a different server if more modularity is needed, or when a server's codebase grows. GSLB can load balance RPCs in the same way it load balances externally visible services._

Sounds cool! Let's dig into gRPC!

### gRPC

We will now implement a gRPC server that will trigger the `cli.Commands`.

It will receive "orders", and depending on the expected call, it will: 

* Implements a `HelloCommand` and trigger its `Run()` function;
* Implements a `GoodbyeCommand` and trigger its `Run()` function

We will also implement a gRPC client.

For the server and the client to communicate, they have to share the same protocol and understand each other with a contract.
_Protocol Buffers (a.k.a., protobuf) are Google's language-neutral, platform-neutral, extensible mechanism for serializing structured data_ 
Even if it's not mandatory, gRPC is usually used with the _Protocol Buffer_. 

So, first, let's implement the _contract_ with/in _protobuf_!

### The protobuf contract

The protocol is described in a simple text file and a specific DSL. Then there is a compiler that serializess the description and turns it into a contract that can be understood by the targeted language.

Here is a simple definition that matches our need:

{{< highlight protobuf >}}
syntax = "proto3";

package myservice;

service MyService {
    rpc Hello (Arg) returns (Output) {}
    rpc Goodbye (Arg) returns (Output) {}
}

message Arg {
    repeated string args = 1;
}

message Output {
    int32 retcode = 1;
}
{{</ highlight >}}

Here is the English description of the contract:

----
Let's take a service called _MyService_. This service provides to actions (commands) remotely:

* _Hello_ 
* _Goodbye_

Both takes as argument an object called _Arg_ that contains an infinite number of _string_ (this array is stored in a field called _args_).

Both actions return an object called _Output_ that returns an integer.

----

The specification is clear enough to code a server and a client. But the string implementation may differ from a language to another.
You may now understand why we need to "compile" the file.
Let's generate a definition suitable for the go language:

`protoc --go_out=plugins=grpc:. myservice/myservice.proto`

_Note_ the definition file has been placed into a subdirectory `myservice`

This command generates a `myservice/myservice.pb.go` file. This file is part of the `myservice` package, **as specified in the myservice.proto**.

The package myservice holds the "contract" translated in `go`. It is full of interfaces and holds helpers function to easily create a server and/or a client.
Let's see how.

### The implementation of the "contract" into the server

Let's go back to the roots and read the doc of gRPC. In the [gRPC basics -  go](https://grpc.io/docs/tutorials/basic/go.html) tutorial is written:

_To build and start a server, we:_

1. _Specify the port we want to use to listen for client requests..._
2. _Create an instance of the gRPC server using grpc.NewServer()._
3. *__Register our service implementation with the gRPC server.__*
4. _Call Serve() on the server with our port details to do a blocking wait until the process is killed or Stop() is called._

Let's decompose the third step.

#### "service implementation"
The `myservice/myservice.pb.go` file has defined an interface for our service.

{{< highlight go >}}
type MyServiceServer interface {
      // Sends a greeting
      Hello(context.Context, *Arg) (*Output, error)
      Goodbye(context.Context, *Arg) (*Output, error)
}
{{</ highlight >}}

To create a "service implementation" in our "cli" utility, we need to create any structure that implements the Hello(...) and Goodbye(...) methods.
Let's call our structure `grpcCommands`:

{{< highlight go >}}
package main

...
import "myservice"
...

type grpcCommands struct {}

func (g *grpcCommands) Hello(ctx context.Context, in *myservice.Arg) (*myservice.Output, error) {
    return &myservice.Output{int32(0)}, err
}
func (g *grpcCommands) Goodbye(ctx context.Context, in *myservice.Arg) (*myservice.Output, error) {
    return &myservice.Output{int32(0)}, err
}
{{</ highlight >}}

_Note_: *myservice.Arg is a structure that holds an array of string named Args. It corresponds to the `proto` definition exposed before.

#### "service registration"

As written in the doc, we need to register the implementation.
In the generated file `myservice.pb.go`, there is a `RegisterMyServiceServer` function.
This function is simply an autogenerated wrapper around the [`RegisterService`](https://godoc.org/google.golang.org/grpc#Server.RegisterService) method of the gRPC [`Server`](https://godoc.org/google.golang.org/grpc#Server) type.

This method takes two arguments: 

* An instance of the gRPC server
* the implementation of the contract.

The 4 steps of the documentation can be implemented like this:

{{< highlight go >}}
listener, _ := net.Listen("tcp", "127.0.0.1:1234")
grpcServer := grpc.NewServer()
myservice.RegisterMyServiceServer(grpcServer, &grpcCommands{})
grpcServer.Serve(listener)
{{</ highlight >}}

So far so good... The code compiles, but does not perform any action and always return 0.

#### Actually calling the `Run()` method

Now, let's use the `grpcCommands` structure as a bridge between the `cli.Command` and the grpc service.

What we will do is to embed the `c.Commands` object inside the structure and trigger the appropriate objects' `Run()` method from the corresponding gRPC procedures.

So first, let's embed the `c.Commands` object.

{{< highlight go >}}
type grpcCommands struct {
      commands map[string]cli.CommandFactory
}
{{</ highlight >}}

Then change the `Hello` and `Goodbye` methods of `grpcCommands` so they trigger respectively:

* `HelloCommand.Run(args)`
* `GoodbyeCommand.Run(args)`

with `args` being the array of string passed via the `in` argument of the protobuf.

as defined in `myservice.Arg.Args` (the protobuf compiler has transcribed the `repeated string args` argument into a filed `Args []string` of the type `Arg`. 

{{< highlight go >}}
func (g *grpcCommands) Hello(ctx context.Context, in *myservice.Arg) (*myservice.Output, error) {
      runner, err := g.commands["hello"]()
      if err != nil {
            return int32(0), err
      }
      ret = int32(runner.Run(in.Args))
      return &myservice.Output{int32(ret)}, err
}
func (g *grpcCommands) Goodbye(ctx context.Context, in *myservice.Arg) (*myservice.Output, error) {
      runner, err := g.commands["goodbye"]()
      if err != nil {
            return int32(0), err
      }
      ret = int32(runner.Run(in.Args))
      return &myservice.Output{int32(ret)}, err
}
{{</ highlight >}}

Let's factorize a bit and create a wrapper (that will be useful in the next section):

{{< highlight go >}}
func wrapper(cf cli.CommandFactory, args []string) (int32, error) {
      runner, err := cf()
      if err != nil {
            return int32(0), err
      }
      return int32(runner.Run(in.Args)), nil
}

func (g *grpcCommands) Hello(ctx context.Context, in *myservice.Arg) (*myservice.Output, error) {
      ret, err := wrapper(g.commands["hello"])
      return &myservice.Output{int32(ret)}, err
}
func (g *grpcCommands) Goodbye(ctx context.Context, in *myservice.Arg) (*myservice.Output, error) {
      ret, err := wrapper(g.commands["goodbye"])
      return &myservice.Output{int32(ret)}, err
}
{{</ highlight >}}

Now we have everything needed to turn our cli into a gRPC service. With a bit of plumbing, the code compiles and the service runs.
The full implementation of the service can be found [here](https://github.com/owulveryck/cli-grpc-example/blob/master/server/main.go).

## A very quick client

The principle is the same for the client. All the needed methods are auto-generated and wrapped by the `protoc` command.

The steps are:

1. create a network connection to the gRPC server (with TLS)
2. create a new instance of myservice'client
3. call a function and get a result

for example:

{{< highlight go >}}
conn, _ := grpc.Dial("127.0.0.1:1234", grpc.WithInsecure())
defer conn.Close()
client := myservice.NewMyServiceClient(conn)
output, err := client.Hello(context.Background(), &myservice.Arg{os.Args[1:]})
{{</ highlight >}}

_Note_: By default, gRPC requires some TLS. I have specified the `WithInsecure` option because I am running on the local loop and it is just an example. Don't do that in production.

# Going further

Normally, Unix tools should respect a [certain philosophy](http://www.faqs.org/docs/artu/ch01s06.html) such as:

<center>**Rule of Silence: When a program has nothing surprising to say, it should say nothing.**</center>

Anyway, we all know that tools are verbose, so let's add a feature that sends the content of stdout and stderr back to the client. (And anyway, we are implementing a service greeting. It would be useless if it was silent :))

## stdout / stderr

What we want to do is to change the output of the commands. 
Therefore, we simply add two more fields to the `Output` object in the protobuf definition:
{{< highlight protobuf >}}
message Output {
    int32 retcode = 1;
    bytes stdout = 2;
    bytes stderr = 3;
}
{{</ highlight >}}

The generated file contains the following definition for `Output`:

{{< highlight go >}}
type Output struct {
      Retcode int32  `protobuf:"varint,1,opt,name=retcode" json:"retcode,omitempty"`
      Stdout  []byte `protobuf:"bytes,2,opt,name=stdout,proto3" json:"stdout,omitempty"`
      Stderr  []byte `protobuf:"bytes,3,opt,name=stderr,proto3" json:"stderr,omitempty"`
}
{{</ highlight >}}

We have changed the Output type, but as all the fields are embedded within the structure, the "service implementation" interface (`grpcCommand`) has not changed.
We only need to change a little bit the implementation in order to return a completed `Output` object:

{{< highlight go >}}
func (g *grpcCommands) Hello(ctx context.Context, in *myservice.Arg) (*myservice.Output, error) {
    var stdout, stderr []byte
    // ...
    return &myservice.Output{ret, stdout, stderr}, err
}
{{</ highlight >}}

Now we have to change the `wrapper` function that has been defined previously to return the content of stdout and stderr:

{{< highlight go >}}
func wrapper(cf cli.CommandFactory, args []string) (int32, []byte, []byte, error) {
    // ...
}
func (g *grpcCommands) Hello(ctx context.Context, in *myservice.Arg) (*myservice.Output, error) {
    var stdout, stderr []byte
    ret, stdout, stderr, err := wrapper(g.commands["hello"], in.Args)
    return &myservice.Output{ret, stdout, stderr}, err
}
{{</ highlight >}}

All the job of capturing stdout and stderr is done within the wrapper function (This solution has been found on [StackOverflow](https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string):

* first, we backup the standard `stdout` and `stderr`
* then, we create two times, two file descriptors linked with a pipe (one for stdout and one for stderr)
* we assign the standard `stdout` and `stderr` to the input of the pipe. From now on, every interaction will be written to the pipe and will be received into the variable declared as output of the pipe
* then, we actually execute the function (the business logic)
* we get the content of the output and save it to variable
* and then we restore stdout and stderr

Here is the implementation of the `wrapper`:
{{< highlight go >}}
func wrapper(cf cli.CommandFactory, args []string) (int32, []byte, []byte, error) {
	var ret int32
	oldStdout := os.Stdout // keep backup of the real stdout
	oldStderr := os.Stderr

	// Backup the stdout
	r, w, err := os.Pipe()
        // ...
	re, we, err := os.Pipe()
        //...
	os.Stdout = w
	os.Stderr = we

	runner, err := cf()
        // ...
	ret = int32(runner.Run(args))

	outC := make(chan []byte)
	errC := make(chan []byte)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.Bytes()
	}()
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, re)
		errC <- buf.Bytes()
	}()

	// back to normal state
	w.Close()
	we.Close()
	os.Stdout = oldStdout // restoring the real stdout
	os.Stderr = oldStderr
	stdout := <-outC
	stderr := <-errC
	return ret, stdout, stderr, nil
}
{{</ highlight >}}

**Et voilà**, the cli has been transformed into a grpc webservice. The full code is available on [GitHub](https://github.com/owulveryck/cli-grpc-example).

### Side note about race conditions

The map used for cli.Command is not concurrent safe. But there is no goroutine that actually writes it so it should be ok.
Anyway, I have written a little benchmark of our function and passed it to the race detector. And it did not find any problem:

```shell
go test -race -bench=.      
goos: linux
goarch: amd64
pkg: github.com/owulveryck/cli-grpc-example/server
BenchmarkHello-2             200          10483400 ns/op
PASS
ok      github.com/owulveryck/cli-grpc-example/server   4.130s
```

The benchmark shows good result on my little chromebook, gRPC seems very efficient, but actually testing it is beyond the scope of this article.

### Interactivity

Sometimes, cli tools ask questions. Another good point with gRPC is that it is bidirectional. Therefore, it would be possible to send the question from the server to the client and get the response back. I let that for another experiment.

## Terraform ?

At the beginning of this article, I have explained that I was using this specific cli in order to derivate Hashicorp tools and turned them into webservices.
Let's take an example with the excellent terraform.

We are going to derivate terraform by changing only its cli interface, add some gRPC powered by protobuf... 

$$\frac{\partial terraform}{\partial cli} + grpc^{protobuf} = \mu service(terraform)$$ [^3]

[^3]: I know, this mathematical equation come from nowhere. But I simply like the beautifulness of this language. (I would have been damned by my math teachers because I have used the mathematical language to describe something that is not mathematical. Would you please forgive me, gentlemen :) 

### About concurrency

Terraform uses [backends](https://www.terraform.io/docs/backends/index.html) to store its states.
By default, it relies on the local filesystem, which is, obviously, not concurrent safe. It does not scale and cannot be used when dealing with webservices.
For the purpose of my article, I won't dig into the backend principle and stick to the local one.
Hence, this will only work with one and only one client. If you plan to do more work around terraform-as-a-service, changing the backend is a must!

### What will I test?

In order to narrow the exercise, I will partially implement the `plan` command.

My test case is the creation of an `EC2` instance on AWS. This example is a copy/paste of the example [Basic Two-Tier AWS Architecture](https://github.com/terraform-providers/terraform-provider-aws/tree/master/examples/two-tier).

I will not implement any kind of interactivity. Therefore, I have added some default values for the ssh key name and path.

Let's check that the basic cli is working:

{{< highlight shell >}}
localhost two-tier [master*] terraform plan | tail
      enable_classiclink_dns_support:   "<computed>"
      enable_dns_hostnames:             "<computed>"
      enable_dns_support:               "true"
      instance_tenancy:                 "<computed>"
      ipv6_association_id:              "<computed>"
      ipv6_cidr_block:                  "<computed>"
      main_route_table_id:              "<computed>"

Plan: 9 to add, 0 to change, 0 to destroy.
{{</ highlight >}}

Ok, let's "hack" terraform!

### hacking Terraform

#### Creating the protobuf contract

The contract will be placed in a `terraformservice` package.
I am using a similar approach as the one used for the greeting example described before:

{{< highlight protobuf >}}
syntax = "proto3";

package terraformservice;

service Terraform {
    rpc Plan (Arg) returns (Output) {}
}

message Arg {
    repeated string args = 1;
}

message Output {
    int32 retcode = 1;
    bytes stdout = 2;
    bytes stderr = 3;
}
{{</ highlight >}}

Then I generate the `go` version of the contract with:

`protoc --go_out=plugins=grpc:. terraformservice/terraform.proto`

### The go implementation of the interface

I am using a similar structure as the one defined in the previous example.
I only change the methods to match the new ones:

{{< highlight go >}}
type grpcCommands struct {
      commands map[string]cli.CommandFactory
}

func (g *grpcCommands) Plan(ctx context.Context, in *terraformservice.Arg) (*terraformservice.Output, error) {
      ret, stdout, stderr, err := wrapper(g.commands["plan"], in.Args)
      return &terraformservice.Output{ret, stdout, stderr}, err
}
{{</ highlight >}}

The wrapper function remains exactly the same as the one defined before because I didn't change the Output format.

### Setting a gRPC server in the main function

The only modification that has to be done is to create a listener for the grpc like the one we did before.
We place it in the main code, just before the execution of the `Cli.Run()` call: 

{{< highlight go >}}
if len(cliRunner.Args) == 0 {
        log.Println("Listening on 127.0.0.1:1234")
        listener, err := net.Listen("tcp", "127.0.0.1:1234")
        if err != nil {
                log.Fatalf("failed to listen: %v", err)
        }
        grpcServer := grpc.NewServer()
        terraformservice.RegisterTerraformServer(grpcServer, &grpcCommands{cliRunner.Commands})
        // determine whether to use TLS
        grpcServer.Serve(listener)
}
{{</ highlight >}}

### Testing it

The code compiles without any problem.
I have triggered the `terraform init` and I have a listening process waiting for a call:

```shelli
~ netstat -lntp | grep 1234
(Not all processes could be identified, non-owned process info
 will not be shown, you would have to be root to see it all.)
tcp        0      0 127.0.0.1:1234          0.0.0.0:*               LISTEN      9053/tfoliv     
```
Let's launch a client:

{{< highlight go >}}
func main() {
      conn, err := grpc.Dial("127.0.0.1:1234", grpc.WithInsecure())
      if err != nil {
            log.Fatal("Cannot reach grpc server", err)
      }
      defer conn.Close()
      client := terraformservice.NewTerraformClient(conn)
      output, err := client.Plan(context.Background(), &terraformservice.Arg{os.Args[1:]})
      stdout := bytes.NewBuffer(output.Stdout)
      stderr := bytes.NewBuffer(output.Stderr)
      io.Copy(os.Stdout, stdout)
      io.Copy(os.Stderr, stderr)
      fmt.Println(output.Retcode)
      os.Exit(output.Retcode)
}
{{</ highlight >}}

```shell
~ ./grpcclient
~ echo $?
~ 0
```

Too bad, the proper function has been called, the return code is ok, but all the output went to the console of the server... Anyway, the RPC has worked.

I can even remove the default parameters and pass them as an argument of my client:

```shell
~ ./grpcclient -var 'key_name=terraform' -var 'public_key_path=~/.ssh/terraform.pub'
~ echo $?
~ 0
```

And let's see if I give a non existent path:

```shell
~ ./grpcclient -var 'key_name=terraform' -var 'public_key_path=~/.ssh/nonexistent'
~ echo $?
~ 1
```

_about the output_: I have been a little optimistic about the stdout and stderr. 
Actually, to make it work, the best option would be to implement a custom `UI` (it should not be difficult because [`Ui is also an interface`](https://godoc.org/github.com/mitchellh/cli#Ui)).
I will try an implementation as soon as I will have enough time to do so. But for now, I have reached my first goal, and this post is long enough :)

# Conclusion

Transforming terraform into a webservice has required a very little modification of the terraform code itself which is very good for maintenance purpose:

{{< highlight diff >}}
diff --git a/main.go b/main.go
index ca4ec7c..da5215b 100644
--- a/main.go
+++ b/main.go
@@ -5,14 +5,18 @@ import (
        "io"
        "io/ioutil"
        "log"
+       "net"
        "os"
        "runtime"
        "strings"
        "sync"
 
+       "google.golang.org/grpc"
+
        "github.com/hashicorp/go-plugin"
        "github.com/hashicorp/terraform/helper/logging"
        "github.com/hashicorp/terraform/terraform"
+       "github.com/hashicorp/terraform/terraformservice"
        "github.com/mattn/go-colorable"
        "github.com/mattn/go-shellwords"
        "github.com/mitchellh/cli"
@@ -185,6 +189,18 @@ func wrappedMain() int {
        PluginOverrides.Providers = config.Providers
        PluginOverrides.Provisioners = config.Provisioners
 
+       if len(cliRunner.Args) == 0 {
+               log.Println("Listening on 127.0.0.1:1234")
+               listener, err := net.Listen("tcp", "127.0.0.1:1234")
+               if err != nil {
+                       log.Fatalf("failed to listen: %v", err)
+               }
+               grpcServer := grpc.NewServer()
+               terraformservice.RegisterTerraformServer(grpcServer, &grpcCommands{cliRunner.Commands})
+               // determine whether to use TLS
+               grpcServer.Serve(listener)
+       }
+
        exitCode, err := cliRunner.Run()
        if err != nil {
                Ui.Error(fmt.Sprintf("Error executing CLI: %s", err.Error()))
{{</ highlight >}}

Of course, there is a bit of work to setup a complete terraform-as-a-service architecture, but it looks promising.

Regarding grpc and protobuf:
gRPC is a very nice protocol, I am really looking forward an implementation in javascript to target the browser
(Meanwhile it is possible and easy to set up a grpc-to-json proxy if any web client is needed). 

But it reminds us that the main target of RPC is machine-to-machine communication. This is something that the ease-of-use-and-read of json has shadowed...

