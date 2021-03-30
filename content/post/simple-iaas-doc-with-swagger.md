---
author: Olivier Wulveryck
date: 2015-11-11T14:24:43+01:00
description: Experience with swagger-ui as a documentation tool for the simple iaas api
draft: false
tags:
- swagger
- api
- documentation
- JSON
- REST
- simple-iaas
title: Simple IaaS API documentation with swagger
type: post
---

In a [previous post](http://blog.owulveryck.info/2015/11/10/iaas-like-restfull-api-based-on-microservices/) I have explained how to develop a very simple API server.

Without the associated documentation, the API will be useless. Let's see how we can use [swagger-ui](https://github.com/swagger-api/swagger-ui) 
in this project to generate a beautiful documentation.

*Note* I'm blogging and experimenting, of course, in the "real" life, it's a lot better to code the API interface before implementing the middleware.

# About Swagger

Swagger is a framework. On top of the swagger project is composed of several tools.

The entry point is to write the API interface using the [Swagger Formal Specification](http://swagger.io/specification/). I will the use the [swagger-ui](https://github.com/swagger-api/swagger-ui) to display the documentation.
The swagger-ui can be modified and recompiled, but I won't do it (as I don't want to play with nodejs). Instead I will rely on the "dist" part which can be used "as-is"


# Defining the API interface with Swagger

## Header and specification version:

Swagger comes with an editor which can be used [online](http://editor.swagger.io/#/).

I will use swagger spec 2.0, as I don't see any good reason not to do so. Moreover, I will describe the API using the `YAML` format instead of the JSON format to be human-friendly.

Indeed, in my `YAML` skeleton the header of my specs will then look like this:

```yaml
swagger: '2.0'
info:
  version: 1.0.0
    title: 'Very Simple IAAS'
```

## The node creation: a POST method
Let's document the Node creation (as it is the method that we have implemented before).

The node creation is a `POST` method, that produces a JSON in output with the request ID of the node created.

The responses code may be:

* 202 : if the request has been taken in account
* 400 : when the request is not formatted correctly
* 500 : if any unhanldled exception occurred
* 502 : if the backend is not accessible (either the RPC server or the backend)

So far, the YAML spec will look like:
```yaml
paths:
  /v1/nodes:
    post:
      description: Create a node
      produces:
        - application/json
      responses:
        202:
          description: A request ID.
        400:
          description: |
            When the request is malformated or when mandatory arguments are missing
        500:
          desctiption: Unhandled error
        502:
          description: Execution backend not available
```

So far so good, let's continue with the input payload. The payload will be formatted in JSON, so I add this directive to the model:

```YAML
consumes:
  - application/json
```

I've decided in my previous post that 6 parameters were needed: 

- the kind of os 
- the size of the machine
- the initial disk size allocated
- the lease (in days)
- the environment 
- the description 

All the parameters will compose a payload and therefore will be present in the body of the request.
The YAML representation of the parameters array is:

```YAML
parameters:
  - name: kind
    in: body
    description: "The OS type"
    required: true
  - name: size 
    in: body
    description: "The size of the (virtual) Machine"
    required: true
  - name: disksize
    in: body
    description: "The initial disk capacity allocated"
    required: true
  - name: leasedays
    in: body
    description: "The lease (in days)"
    required: true
  - name: environment_type
    in: body
    description: "The target environment"
  - name: description
    in: body
    description: "The target environment"
```

Sounds ok, but when I test this implementation in the swagger editor for validation, I get this error:

```
Swagger Error
Data does not match any schemas from 'oneOf'
```

_STFWing and RTFMing..._

in the [Specifications](http://swagger.io/specification/#parameterObject), I have found this line:

<html>
If <a href="#parameterIn"><code>in</code></a> is <code>"body"</code>:</p>
<table>
<thead>
<tr>
<th>Field Name</th>
<th style="text-align: center;">Type</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td><a name="parameterSchema"></a>schema</td>
<td style="text-align: center;"><a href="#schemaObject">Schema Object</a></td>
<td><stwrong>Required.</strong> The schema defining the type used for the body parameter.</td>
</table>
</html>

Therefore, I should set a schema object for every parameter in order to define its type. In this example, I don't want to go too deeply into the swagger specification, so I won't define any type.

So I have tested the following:
```YAML
parameters:
  - name: kind
    in: body
    description: "The OS type"
    required: true
    schema:
      type: string
  - name: size 
    in: body
    description: "The size of the (virtual) Machine"
    required: true
    schema:
      type: string
    ...
```

And again, I had a validation error from the editor:

<span class="icon">✖</span> Swagger Error</h5><p class="error-description" ng-bind-html="error.description">Operation cannot have multiple body parameters</p>

_RTFMing..._

> Body - The payload that's appended to the HTTP request. 
> Since there can only be one payload, there can only be one body parameter. 
> The name of the body parameter has no effect on the parameter itself and is used for documentation purposes only. 
> Since Form parameters are also in the payload, body and form parameters cannot exist together for the same operation.

What I must do, is to create a custom type _nodeRequest_ with the input fields as properties and reference it in the body.

Here is the complete structure:

```YAML
parameters:
  - name: nodeRequest 
    in: body
    description: a node request
    required: true
    schema:
        $ref: '#/definitions/NodeRequest'
``` 

And the proper NodeRequest definition in the _definition_ area:

```YAML
definitions:
  NodeRequest:
    description: A Node Request object
    properties:
      kind:
        type: string
        description: The OS type
      size:
        type: string
        description: The size of the (virtual) machine
      disksize:
        type: integer
        format: int32
        description: The initial disk capacity size (in GB)
      leasedays:
        type: integer
        format: int32
        description: The lease
      environment_type:
        type: string
        description: the target environment
      description:
        type: string
```

OK ! The swagger file is valid... Now let's glue it together with swagger-ui and serve it from the GO API server I have developed before

# Integrating swagger-ui

As written in the README in the github of the project, swagger-ui can be used "as-is" using the files in the _dist_ folder. Let's get the files from github:
```shell
/tmp #  git clone https://github.com/swagger-api/swagger-ui.git
Cloning into 'swagger-ui'...
remote: Counting objects: 7292, done.
remote: Compressing objects: 100% (33/33), done.
remote: Total 7292 (delta 8), reused 0 (delta 0), pack-reused 7256
Receiving objects: 100% (7292/7292), 19.20 MiB | 1021.00 KiB/s, done.
Resolving deltas: 100% (3628/3628), done.
Checking connectivity... done.
```

Let's checkout our project:

```shell
/tmp # git clone https://github.com/owulveryck/example-iaas.git 
...
```

and move the `dist` folder into the project:
```
mv /tmp/swagger-ui/dist /tmp/example-iaas
```

## Adding a route to the GO server to serve the static files

I cannot simply add a route in the `routes.go` file for this very simple reason: 
The loop used in the `router.go` is using the `Path` method, and to serve the content of the directory, I need to use the `PathPrefix` method (see [The Gorilla Documentation](http://www.gorillatoolkit.org/pkg/mux#Route.PathPrefix) for more information).

To serve the content, I add this entry to the muxrouter in the `router.go` file:

```go 
router.
       Methods("GET").
       PathPrefix("/apidocs").
       Name("Apidocs").
       Handler(http.StripPrefix("/apidocs", http.FileServer(http.Dir("./dist"))))
```

Then I start the server and point my browser to http://localhost:8080/apidocs...

Wait, nothing is displayed...

# The final test

As I serve the files from the `./dist` directory, what I need to do is to move my `swagger.yaml` spec file into the dist subfolder and tell swagger to read it.

_Et voilà!_

<center>
<img class="img-square img-responsive" src="/assets/images/swagger.png" alt="Result"/>
</center>

# Final word

As you can see, there is a "Try out" button, which triggers a `curl` command... Very helpful to enter a test driven development mode.

On top of that swagger is really helpful and may be a great tool to synthesize the need of a client in term of an interface.
Once the API is fully implemented, any client binding may also be generated with the swagger framework.

No not hesitate to clone the source code from [github](https://github.com/owulveryck/example-iaas) and test the swagger.yaml file in the editor to see how the bindings are generated

You can find all the codes in the github repository [here](https://github.com/owulveryck/example-iaas) in the branch `simple-iaas-api-documentation-with-swagger`

The final YAML file can be found [here](https://github.com/owulveryck/example-iaas/blob/simple-iaas-api-documentation-with-swagger/swagger.yaml)
