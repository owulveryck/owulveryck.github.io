---
categories:
date: 2017-09-12T13:28:36+02:00
description: "This is a second part of the last article. I now really dig into Terraform. This article will explain how to use the Terraform sub-packages in order to create a brand new binary that acts as a gRPC server instead of a cli."
draft: false
images:
- https://nhite.github.io/images/logo.png
tags:
title: Terraform is hip... Introducing Nhite
---

In a previous post, I did some experiments with gRPC, protocol buffer and Terraform.
The idea was to transform the "Terraform" cli tool into a micro-service thanks to gRPC.

This post is the second part of the experiment. I will go deeper in the code and see if it is possible
to create a brand new utility, without hacking Terraform. The idea is to import some packages that compose the binary
and create my own service based on gRPC.

# The Terraform structure

Terraform is a binary utility written in `go`.
The `main` package resides in the root directory of the `terraform` directory.
As usual with go projects, all other subdirectories are different modules.

The whole business logic of Terraform is coded into the subpackages. The "`main`" package is simply an enveloppe for kick-starting the utility (env variables, etc.) and to initiate the command line.

### The cli implementation

The command line flags are instantiated by Mitchell Hashimoto's cli package.
As explained in the previous post, this cli package is calling a specific function for every action.

### The _command_ package

Every single action is fulfilling the `cli.Command` interface and is implemented in the [`command`](https://godoc.org/github.com/hashicorp/terraform/command) subpackage.
Therefore, every "action" of Terraform has a definition in the command package and the logic is coded into a `Run(args []string) int` method (see the [doc of the Command interface for a complete definition](https://godoc.org/github.com/mitchellh/cli#Command).

# Creating a new binary

The idea is not to hack any of the packages of Terraform to allow an easier maintenance of my code. 
In order to create a custom service, I will instead implement a new utility; therefore a new `main` package.
This package will implement a gRPC server. This server will implement wrappers around the functions declared in the `terraform.Command` package.

For the purpose of my POC, I will only implement three actions of Terraform:

* `terraform init`
* `terraform plan`
* `terraform apply`

## The gRPC contract

In order to create a gRPC server, we need a service definition.
To keep it simple, let's consider the contract defined in the previous post ([cf the section: Creating the protobuf contract](https://blog.owulveryck.info/2017/09/02/from-command-line-tools-to-microservices---the-example-of-hashicorp-tools-terraform-and-grpc.html#creating-the-protobuf-contract)).
I simply add the missing procedure calls:

{{< highlight protobuf >}}
syntax = "proto3";

package pbnhite;

service Terraform {
    rpc Init (Arg) returns (Output) {}
    rpc Plan (Arg) returns (Output) {}
    rpc Apply (Arg) returns (Output) {}
}

message Arg {
    repeated string args = 2;
}

message Output {
    int32 retcode = 1;
    bytes  stdout = 2;
    bytes stderr = 3;
}
{{</ highlight >}}

## Fulfilling the contract

As described previously, I am creating a `grpcCommand` structure that will have the required methods to fulfill the contract:

{{< highlight go >}}
type grpcCommands struct {}

func (g *grpcCommands) Init(ctx context.Context, in *pb.Arg) (*pb.Output, error) {
    ....
}
func (g *grpcCommands) Plan(ctx context.Context, in *pb.Arg) (*pb.Output, error) {
    ....
}
func (g *grpcCommands) Apply(ctx context.Context, in *pb.Arg) (*pb.Output, error) {
    ....
}
{{</ highlight >}}

In the previous post, I have filled the `grpcCommand` structure with a `map` filled with the command definition.
The idea was to keep the same CLI interface.
As we are now building a completely new binary with only a gRPC interface, we don't need that anymore.
Indeed, there is still a need to execute the `Run` method of every Terraform command.

Let's take the example of the Init command. 

Let's see the definition of the command by looking at the [godoc](https://godoc.org/github.com/hashicorp/terraform/command#InitCommand):

{{< highlight go >}}
type InitCommand struct {
    Meta
    // contains filtered or unexported fields
}
{{</ highlight >}}

This command holds a substructure called `Meta`. `Meta` is defined [here](https://godoc.org/github.com/hashicorp/terraform/command#Meta) and holds _the meta-options that are available on all or most commands_. Obviously we need a Meta definition in the command to make it work properly.

For now, let's add it to the `grpcCommand` globally, and we will see later how to implement it.

Here is the gRPC implementation of the contract:

{{< highlight go >}}
type grpcCommands struct {
    meta command.Meta
}

func (g *grpcCommands) Init(ctx context.Context, in *pb.Arg) (*pb.Output, error) {
    // ...
    tfCommand := &command.InitCommand{
        Meta: g.meta,
    }
    ret := int32(tfCommand.Run(in.Args))
    return &pb.Output{ret, stdout, stderr}, err
}
{{</ highlight >}}

## How to initialize the _grpcCommand_  object

Now that we have a proper `grpcCommand` than can be registered to the grpc server, let's see how to create an instance.
As the grpcCommand only contains one field, we simply need to create a `meta` object.

Let's simply copy/paste the code done in Terraform's main package and we now have:

{{< highlight go >}}
var PluginOverrides command.PluginOverrides
meta := command.Meta{
    Color:            false,
    GlobalPluginDirs: globalPluginDirs(),
    PluginOverrides:  &PluginOverrides,
    Ui:               &grpcUI{},
}
pb.RegisterTerraformServer(grpcServer, &grpcCommands{meta: meta})
{{</ highlight >}}

According to the comments in the code, the `globalPluginDirs()` _returns directories that should be searched for
globally-installed plugins (not specific to the current configuration)_.
I will simply copy the function into my main package

## About the UI

In the example CLI that I developed in the previous post, what I did was to redirect stdout and stderr to an array of bytes, in order to capture it and send it back to a gRPC client.
I noticed that this was not working with Terraform.
This is because of the UI!
UI is an interface whose purpose is to get the output stream and write it down to a specific io.Writer.

Our tool will need a custom UI.

### A custom UI

As UI is an interface ([see the doc here](https://godoc.org/github.com/mitchellh/cli#Ui)), it is easy to implement our own. Let's define a structure that holds two array of bytes called `stdout` and `stderr`. Then let's implement the methods that will write into these elements:

{{< highlight go >}}
type grpcUI struct {
    stdout []byte
    stderr []byte
}

func (g *grpcUI) Output(msg string) {
    g.stdout = append(g.stdout, []byte(msg)...)
}
{{</ highlight>}}

_Note 1_: I omit the methods `Info`, `Warn`, and `Error` for brevity.

_Note 2_: For now, I do not implement any logic into the `Ask` and `AskSecret` methods. Therefore, my client will not be able to ask something. But as gRPC is bidirectional, it would be possible to implement such an interaction.

Now, we can instantiate this UI for every call, and assign it to the meta-options of the command:

{{< highlight go >}}
var stdout []byte
var stderr []byte
myUI := &grpcUI{
    stdout: stdout,
    stderr: stderr,
}
tfCommand.Meta.Ui = myUI
{{</ highlight >}}

So far, so good: we now have a new Terraform binary, that is working via gRPC with a very little code.

# What did we miss?

## Multi-stack
It is fun but not usable for real purpose because the server needs to be launched from the directory where the _tf_ files are... 
Therefore the service can only be used for one single Terraform stack... Come on!

Let's change that and pass as a parameter of the RPC call the directory where the server needs to work. It is as simple as adding an extra argument to the `message Arg`:

{{< highlight protobuf >}}
message Arg {
    string workingDir = 1;
    repeated string args = 2;
}
{{</ highlight >}}

and then, simply do a `change directory` in the implementation of the command:

{{< highlight go >}}
func (g *grpcCommands) Init(ctx context.Context, in *pb.Arg) (*pb.Output, error) {
    err := os.Chdir(in.WorkingDir)
    if err != nil {
        return &pb.Output{int32(0), nil, nil}, err
    }
    tfCommand := &command.InitCommand{
        Meta: g.meta,
    }
    var stdout []byte
    var stderr []byte
    myUI := &grpcUI{
        stdout: stdout,
        stderr: stderr,
    }
    ret := int32(tfCommand.Run(in.Args))
    return &pb.Output{ret, stdout, stderr}, err
}
{{</ highlight >}}

## Implementing a new _push_ command

I have a Terraform service. Alright.
Can an "Operator" use it?

The service we have deployed is working exactly like Terraform. I have only changed the user interface.
Therefore, in order to deploy a stack, the 'tf' files must be present locally on the host.

Obviously we do not want to give access to the server that hosts Terraform. This is not how micro-services work.

Terraform has a push command that Hashicorp has implemented to communicate with Terraform enterprise.
This command is linked with their close-source product called "Atlas" and is therefore useless for us.

Let's take the same principle and implement our own _push_ command.

### Principle

The push command will zip all the `tf` files of the current directory in memory, and transfer the zip via a specific message to the server.
The server will then decompress the zip into a unique temporary directory and send back the ID of that directory.
Then every other Terraform command can use the id of the directory and use the stack (as before).

Let's implement a protobuf contract:

{{< highlight protobuf >}}
service Terraform {
    // ...
    rpc Push(stream Body) returns (Id) {}
}

message Body {
    bytes zipfile = 1;
}

message Id {
    string tmpdir = 1;
} 
{{</ highlight >}}

_Note_: By now I assume that the whole zip can fit into a single message. I will probably have to implement chunking later

Then instantiate the definition into the code of the server:

{{< highlight go >}}
func (g *grpcCommands) Push(stream pb.Terraform_PushServer) error {
    workdir, err := ioutil.TempDir("", ".terraformgrpc")
    if err != nil {
    return err
    }
    err = os.Chdir(workdir)
    if err != nil {
    return err
    }

    body, err := stream.Recv()
    if err == io.EOF || err == nil {
        // We have all the file
        // Now let's extract the zipfile
        // ...
    }
    if err != nil {
        return err
    }
    return stream.SendAndClose(&pb.Id{
            Tmpdir: workdir,
    })
}
{{</ highlight >}}

# going further...

The problem with this architecture is that it's stateful, and therefore easily scalable.

A solution would be to store the zip file in a third party service, identify it with a unique id.
And then call the Terraform commands with this ID as a parameter. 
The Terraform engine would then grab the zip file from the third party service if needed and process the file

## Implementing a micro-service of backend

I want to keep the same logic, therefore the storage service can be a gRPC microservice.
We can then have different services (such as s3, google storage, dynamodb, NAS, ...) written in different languages.

The Terraform service will act as a client of this "backend" service (take care, it is not the same backend as the one defined within Terraform).

Our Terraform-service can then be configured in runtime to call the host/port of the correct backend-service. We can even imagine the backend address being served via consul.

This is a work in progress and may be part of another blog post.

# Hip[^1] is _cooler than cool_: Introducing _Nhite_

[^1]: [hip definition on wikipedia](https://en.wikipedia.org/wiki/Hip_(slang))

I have talked to some people about all this stuff and I feel that people are interested.
Therefore, I have set up a GitHub organisation and a GitHub project to centralize what I will do around that.

The project is called Nhite.

* The GitHub organization is called [nhite](https://github.com/nhite)
* The web page is [https://nhite.github.io](https://nhite.github.io)

There is still a lot to do, but I really think that this could make sense to create a community. It may give a product by the end, or go in my attic of dead projects.
Anyway, so far I've had a lot of fun!
