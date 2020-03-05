---
author: Olivier Wulveryck
date: 2016-10-19T07:24:51+02:00
description: How to write a Single Page Application MVC without blowing your mind with Javascript and a Framework.
draft: false
keywords:
tags:
- javascript
- jQuery
- Gopherjs
- SPA
- WebUI
- Cordova
title: SPA with Gopherjs
topics:
- Wed UI
type: post
---

# Introduction

Single page application (SPA) are a standard when dealing with mobile apps.
Unfortunately, by now, JavaScript is the only programming language supported by a browser.

Therefore, to code web UI it remains a must.

## Life of an ex sysadmin who wants to code a web app: getting depressed

To make the development easier, your friend, who is "web developer" will recommend you to use a marvelous framework.
Depending on the orientation of the wind, the moon or its reading, he will encourage you to use `reactjs`, `angular`, `ember` or whatever exotic
tool.

With some recommendation from my real friends and from Google, I've started an app based on [ionic](http://ionicframework.com/) which is based on [angular](https://angularjs.org/).
As I did not know anything about angular, I've watched a (very good) [introduction](https://www.youtube.com/watch?v=i9MHigUZKEM) and followed the ionic tutorial.

So far so good...

Then I implemented a SSO with Facebook. I wrote a backend in `go` to handle the token generation and the used database connection.
I started to code it by hand, until a friend tells me about the angular module [Satellizer](https://github.com/sahat/satellizer) that was suppose to handle the logic for me.
And it did.... It was suddenly automagic:

Everything was working on my browser. I was happy, So I decided to deploy my app on my iPhone and enjoy the power of Cordova.

That's when the headache started: There was something wrong on the mobile phone version. A bug!

I tried to debug it, with Xcode, with Safari... The more I was searching, the more I had to dive into the framework. Too many magic in it for 
something that was, in fine, not a bug [^1].

I asked some help from a friend and his first reply was: "which version of angular? Because in version 2 they have changed a lot of concepts"

That was too much.
I considered that this world made of JavaScript, frameworks, grunt, bower, gulp, npm or whatever fancy tool was definitely not for me.
Too many work to learn something already outdated.

On top of that, I've never been a callback man, I hate them since my X11/Motif programming course. I do like CSP!

<center>
![Matt Holt's Tweet](/assets/images/not_my_type.png)
</center>
## Out of the depression!

Ok, I abandoned those tools. But I still want to code my app, and I'm not the kind of guy that easily give up.

Let's resume my needs:

* I need a MVC, because it's the most natural way to code web ui today
* MVC is not framework dependent
* A SPA is the good choice for a mobile app and Cordova makes things easy
* Javascript is mandatory

I've digge a little bit and I've found this blog post: [Do you really want an SPA framework?](https://mmikowski.github.io/no-frameworks/) which leads me to "the solution": 

I will code my model/view/controller from scratch.
But as I want to preserve my health and stay away from Javacript, I will code it with something fun: go.

At the last doGo.eu [Dmitri Shuralyov](https://twitter.com/shurcool) gave a very good introduction about [gopherjs](https://github.com/gopherjs/gopherjs). Gopherjs is a [transpiler](https://en.wikipedia.org/wiki/Source-to-source_compiler) that converts go code into javascript.

You can code all your logic in go and transpile it in javascript, or you can use is to access you DOM or other javascript libraries.

A bunch of bindings to famous javascript libraries such as jQuery already exists on gihtub,

Let's see an example and implement a very basic routing mechanism relying on a pure js library.

# Examples

I will code a little page based on bootstrap.

## A basic go code

The dynamic part will be coded in pure GO in a file called `main.go` and transpiled into javscript code with the help of the `gopherjs` command:

```bash
gopherjs build main.go -o js/app.js -m
```

To make things easier, I can add some directives in my go code in the form of a comment:

```go
package main
//go:generate gopherjs build main.go -o js/app.js -m
// +build ignore

import (
  "log"
)

func main() {
    log.Println("Hello World")
}
```

Therefore I will be able to generate my code directly with a simple call to `go generate` and it will produce the `js/app.js` file for me.

## The page

The structure of the main page is taken from bootstrap's [starter template](http://getbootstrap.com/examples/starter-template/#) 

I simply add my javascript file generated with gopherjs :

```html
<script src="js/app.js"></script>
```

If I launch my page, I will have a "hello world" written in the javascript browser of my console.

# the SPA

Now, I will implement a very basic SPA.
It will display three tabs accessible by their names (for demo purpose):

* [/#](/#)
* [/#about](/#about)
* [/#contact](/#contact)

I want to trigger a javascript code that could change the content of the body by clicking on the links.


## Routing

A good SPA needs a good routing system to deal with anchor refs.

There are several implementation of gopherjs based routing mechanism.

But, for the purpose of this blog post, I will use a pure Javascript routing library: [Director.js](https://github.com/flatiron/director#routing-table). It's the router used in the todoMVC example and it will allow me to show how to interact with global javascript objects.

The first thing to do is to include the js file at the end of the `index.html`

```html
<script src="js/director.min.js"></script>
<script src="js/app.js"></script>
```

Then I will create a GO type _Router_ that will correspond to the Router objet in javascript.
To do so, the Router type must be a [*js.Object](https://godoc.org/github.com/gopherjs/gopherjs/js#Object)

```go
import "github.com/gopherjs/gopherjs/js"

type Router struct {
    *js.Object
}
```

Then I define a constructor, that simply get the router object from the global scope of the javascript engine:
```go
func NewRouter() *Router {
    return &Router{js.Global.Get("Router").New()}
}
```

Then, to actually implement my [adhoc-routing](https://github.com/flatiron/director#adhoc-routing) as described in the doc of director.js,
I must implement the `on` and the `init` bindings.

Once done, I add the routes in my `main` func:

{{< gist owulveryck 3256d582ad2241eeeaf118d5bf9c1cd0 "router.go" >}}

If I launch the page, I can now click on the links and it will diplay hello in my console.

You can check the full code on [this gist](https://gist.github.com/owulveryck/3256d582ad2241eeeaf118d5bf9c1cd0)

You see that I've let the function as `notImplementedYet`, but replacing it with a jQuery call is trivial:

```go
import "github.com/gopherjs/jquery"

//convenience:
var jQuery = jquery.NewJQuery

func content() {
    jQuery("#main").SetText("Welcome to GopherJS")
}
```

# Conclusion

Gopherjs is not trivial, but it has the ability to make the web development more structured. 
I've started a web ui from scratch and reach the same goal as the one I reached in javacript in only 2 days (compared to 3 weeks).

Of course, a javascript-master-of-the-world would argue that he would implement it in 2 hours, but that's not the point here.
The point is that I can use all the "benefits" of the go principles easily to write a web ui.

You can check the development of the [Nhite fronted](https://github.com/nhite/frontend) to watch the progress I will make (or not) with this technology.

----
[^1]:
1 - Actually, I figured out what the "bug" was later, when I finished the implementation in go and there was no magic anymore in the code.
The oauth2 flow I use is "[Authorization code](https://tools.ietf.org/html/rfc6749#section-1.3.1)". In this flow, you query the authorization server (here facebook) and send it the client identifier and _a redirection URI_.
In my dev environment this redirection URI is set to "http://localhost". Once the user is logged in (on the Facebook page), the navigation window redirects him in the application at localhost.
When running on iOS with _cordova_ the files are served locally (file://,,,) and there is no way to specify a redirect URI that point to file://, therefore the redirect URI must point somewhere else... but in this case, getting the code from the application becomes tricky because of the security policies. I could do a complete blog post about this.

