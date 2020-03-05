+++
date = "2015-10-23T09:54:27+01:00"
draft = false
title = "Simple polling, a cloud native app - part 1"

+++

In this series of posts I'll explain how to setup a simple polling app, the cloud way.
This app, written in go, will be hosted on a PAAS, and I've chosen the [Google App Engine](https://cloud.google.com/appengine/docs) for convenience.

I will not explain in this post how to setup the Development environment as it is described [here](https://cloud.google.com/appengine/docs/go/gettingstarted/devenvironment)

# A word about the Hosting

Google Apps Engine is a cloud service aim to host applications without worrying about scalability, and technical architecture of the hosting environement.
The application is seen as a web service and proxyfied and load balanced in several compute nodes.
The storage service is provided via a schemaless NoSQL datastore, featuring:

* high availability
* consistency
* ...

and basically all the strong features you would expect from a decent production ready database provider.

To be simple: Take care of the functionnality of your app and respect the cloud principles (see for example the [12factor.net](http://12factor.net) ), and google will take care that it can run in the best conditions.

# The principles of the application

The application is composed of a stateless engine, a data bucket and some web pages for the presentation.

## The engine

The engine is the core of the application. It is the "computational element" that will take care of the inputs of the user and interact with the storage.
It is a [GO](http://golang.org) developement.

## The bucket

It is the data warehouse. It will store the participant name and its answer. It will be NoSQL based, the key will be the username and the value its answer.

## The web pages

* the question "will you participate" and a form input where you will be able to write your name and three buttons "yes", "no" and ""maybe"".
```
Will you participate 
+------------------+  +-----+ +-----+ +-------+
|  Your name       |  | YES | |  NO | | Maybe |
+------------------+  +-----+ +-----+ +-------+
```

* a simple table with two columns:

* One will hold the participant name
* The other one will display its response

```
+---------------------+-------+
|  John doe           | YES   |
+---------------------+-------+
|  Johnny Vacances    | NO    |
+---------------------+-------+
|  Foo Bar            | YES   |
+---------------------+-------+
|  Toto Titi          | NO    |
+---------------------+-------+
|  Pascal Obistro     | YES   |
+---------------------+-------+
```

# Setting up the development environment

First, we will create a directory that will host the sources of our application in our `GOPATH/src`.
_Note_: For convenience I've created a github repo named "google-app-example" to host the complete source.

```
~ mkdir -p $GOPATH/src/github.com/owulveryck/google-app-example
~ cd $GOPATH/src/github.com/owulveryck/google-app-example
~ git init
~ git remote add origin https://github.com/owulveryck/google-app-example
```


## Hello World!

Let's create the hello world first to validate the whole development chain.
As written in the doc, create the two files `hello.go` and `app.yaml`.
Obviously the `simple-polling.go` file will hold the code of the application. Let's focus a bit on the _app.yaml_ file.
The documentation of the _app.yaml_ file is [here](https://cloud.google.com/appengine/docs/go/config/appconfig). The goal of this file is to specifiy the runtime configuration of the engine.
This simple file replace the "integration" task for an application typed "born in the datacenter"

Here is my app.yaml
```yaml
application: simple-polling
version: 1
runtime: go
api_version: go1

handlers:
- url: /.*
  script: _go_app
```

And then we try our application with the command `goapp serve $GOPATH/src/github.com/owulveryck/google-app-example/`
which should display something similar to:
```
INFO     2015-10-26 21:02:10,295 devappserver2.py:763] Skipping SDK update check.
INFO     2015-10-26 21:02:10,468 api_server.py:205] Starting API server at: http://localhost:52457
INFO     2015-10-26 21:02:12,011 dispatcher.py:197] Starting module "default" running at: http://localhost:8080
INFO     2015-10-26 21:02:12,014 admin_server.py:116] Starting admin server at: http://localhost:8000
```

Then I can open my browser (or curl) and point it to http://localhost:8080 to see my brand new "Hello World!" displayed

```
~ curl http://localhost:8080
Hello, world!
```

This ends the part 1 of this serie of articles.
