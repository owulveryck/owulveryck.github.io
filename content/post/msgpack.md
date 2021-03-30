+++
date = "2015-11-10T08:56:36+01:00"
draft = false
title = "IaaS-like RESTfull API based on microservices"
tags = [
    "Development",
    "ruby",
    "go",
    "API",
    "REST",
    "msgpack",
    "simple-iaas"
]
+++

# Abstracts

Recently, I've been looking at the principles of a middleware layer and especially on how a RESTFULL API could glue different IT services together.

I am reading more and more about the "API economy"

I've also seen this excellent video made by [Mat Ryer](https://www.youtube.com/watch?v=tIm8UkSf6RA&list=PLDWZ5uzn69ezRJYeWxYNRMYebvf8DerHd) about how to code an API in GO and why go would be the perfect language to code such a portal.

The problem I'm facing is that in the organization I'm working for, the developments are heterogeneous and therefore you can find *ruby* teams as well as *python* teams and myself as a *go* team (That will change in the future anyway)
The key point is that I would like my middleware to serve as an entry point to the services provided by the department.

We (as an "ops" team) would then be able to present the interface via, for example, a [swagger](http://swagger.io) like interface, take care of the API and do whatever RPC to any submodule.

# An example: a IAAS like interface

Let's consider a node compute lifecycle.

What I'd like to be able to do is:

* to create a node
* to update a node (maybe)
* to delete a node
* to get the status of the node

## The backend

The backend is whatever service, able to create a node, such as openstack, vmware vcac, juju, ... 
Thoses services usually provide RESTfull API.

I've seen in my experience, that usually, the API are given with a library in a so called "modern language". 
This aim to simplify the development of the clients.
Sometimes this library may also be developed by an internal team that will take care of the maintenance.

## The library

In my example, we will consider that the library is a simple _gem_ file developed in ruby. 
Therefore, our service will be a simple server that will get RPC calls, call the good method in the _gemfile_ 
and that will, _in fine_ transfer it to the backend.

## The RestFull API.

I will use the example described [here](http://thenewstack.io/make-a-restful-json-api-go/) as a basis for this post.
Of course there are many other examples and excellent go packages that may be used, but according to Mat Ryer, I will stick to the idiomatic approach.

## The glue: MSGPACK-RPC

There are several methods for RPC-ing between different languages. Ages ago, there was xml-rpc; then there has been json-rpc; 
I will use [msgpack-rpc](https://github.com/msgpack-rpc/msgpack-rpc) which is a binary, json base codec.
The communication between the Go client and the ruby server will be done over TCP via HTTP for example.

Later on, outside of the scope of this post, I may use ZMQ (as I have already blogged about 0MQ communication between those languages).

# The implementation of the Client (the go part)

I will describe here the node creation via a POST method, and consider that the other methods could be implemented in a similar way.

## The signature of the node creation

Here is the expected signature for creating a compute element:

```json
{
    "kind":"linux",
    "size":"S",
    "disksize":20,
    "leasedays":1,
    "environment_type":"dev",
    "description":"my_description",
}
```

The corresponding GO structure is:

```go
type NodeRequest struct {
    Kind string `json:"kind"` // Node kind (eg linux)
    Size string `json:"size"` // size
    Disksize         int    `json:"disksize"`
    Leasedays        int    `json:"leasedays"`
    EnvironmentType  string `json:"environment_type"`
    Description      string `json:"description"`
}
```

## The route

The Middleware is using the [gorilla mux package](http://www.gorillatoolkit.org/pkg/mux). 
According the description, I will add an entry in the routes array (into the _routes.go_ file):

```go
Route{
    "NodeCreate",
    "POST",
    "/v1/nodes",
    NodeCreate,
},
```

*Note* : I am using a prefix `/v1` for my API, for exploitation purpose.

I will then create the corresponding handler in the file with this signature

```go
func NodeCreate(w http.ResponseWriter, r *http.Request){
    var nodeRequest NodeRequest
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &nodeRequest); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }    
}
```

That's in this function that will be implemented RPC (client part). To keep it simple at the beginning, 
I will instantiate a TCP connection on every call.
Don't throw things at me, that will be changed later following the advice of Mat Ryer.

## The implementation of the handler

### The effective remote procedure call

To use _msgpack_ I need to import the go implementation `github.com/msgpack-rpc/msgpack-rpc-go/rpc`.
This library will take care of the encoding/decoding of the messages.

Let's dial the RPC server and call the `NodeCreate` method with, as argument, the information we had from the JSON input

```go
    conn, err := net.Dial("tcp", "127.0.0.1:18800")
    if err != nil {
        fmt.Println("fail to connect to server.")
        return
    }
    client := rpc.NewSession(conn, true)
    retval, err := client.Send("NodeCreate", nodeRequest.Kind, nodeRequest.Size, nodeRequest.Disksize, nodeRequest.Leasedays, nodeRequest.EnvironmentType, nodeRequest.Description)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(retval)
```
# The RPC server (the ruby part)

This part is written in ruby, and will take care of the effective node creation.
At first, we should install the GEM file with the command `gem install msgpack-rpc`.

```ruby
require 'msgpack/rpc'
class MyHandler
    def NodeCreate(kind, size, disksize, leasedays, environmenttype, description) 
        print "Creating the node with parameters: ",kind, size, disksize, leasedays, environmenttype, description
        return "ok"
    end
end
svr = MessagePack::RPC::Server.new
svr.listen('0.0.0.0', 18800, MyHandler.new)
svr.run
```

# let's test it

Launch the RPC server:
`ruby server.rb`

Then launch the API rest server

`go run *go`

Then perform a POST request

```shell
curl -X POST -H 'Content-Type:application/json' -H 'Accept:application/json' -d '{"kind":"linux","size":"S","disksize":20,"leasedays":1,"environment_type":"dev","description":"my_description"}' -k http://localhost:8080/v1/nodes
```

It should write something like this: 
```
2015/11/10 13:56:51 POST        /v1/nodes       NodeCreate      2.520673ms
ok
```

And something like this in the output of the ruby code:
```
Creating the node with parameters: linux S 20 1 dev my_description
```

That's all folks! What's left:

* To implement the other methods to be "[CRUD](https://en.wikipedia.org/wiki/Create,_read,_update_and_delete)" compliant
* To implement an authentication and accreditation mechanism (JWT, Oauth, ?)
* To change the implementation of the RPC client to use a pool instead of a single connection
* To implement the swagger interface and documentation of the API
* Whatever fancy stuff you may want from a production ready interface.

You can find all the codes in the github repository [here](https://github.com/owulveryck/example-iaas) in the branch `iaas-like-restfull-api-based-on-microservices`
