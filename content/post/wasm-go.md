---
title: "Some notes about the upcoming WebAssembly support in Go"
date: 2018-06-08T21:23:17+02:00
lastmod: 2018-06-08T21:23:17+02:00
draft: false
keywords: []
description: "This is a very quick post with some notes about the support of WebAssembly (wasm) in the Go toolchain. This article is not a tutorial and as any information it contains may be obsolete soon. The Go api for Wasm is not stable yet."
tags: []
categories: []
author: "Olivier Wulveryck"

# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: true
toc: true
autoCollapseToc: false
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: false
---

<!--more-->
This is a rapid post about webassembly. Its goal is to act as a reminder for me more than a tutorial on how to use it.

The upcoming version of go 1.11 will have support for Wasm.
@neelance has done most of the work of implementation. The support for wasm can already be tested by extracting his working branch from GitHub.

See [this article](https://blog.gopheracademy.com/advent-2017/go-wasm/) for more information.


# Setup of the toolchain

To generate a wasm file from a go source, you need to get and patch the go toolset from the sources:

```
~ mkdir ~/gowasm
~ git clone https://go.googlesource.com/go ~/gowasm
~ cd ~/gowasm
~ git remote add neelance https://github.com/neelance/go
~ git fetch --all
~ git checkout wasm-wip
~ cd src
~ ./make.bash
```

Then to use this version of go, point the `GOROOT` to `~/gowasm` and use the binaries present in `~/gowasm/bin/go`.

# First sample

As usual, the first sample is a "hello world". Let's write this:

{{< highlight go >}}
package main

import "fmt"

func main() {
        fmt.Println("Hello World!")
}
{{</ highlight >}}

and compile it into a file called `example.wasm`:

`GOROOT=~/gowasm GOARCH=wasm GOOS=js ~/gowasm/bin/go build -o example.wasm main.go`

## Running the sample

Here is an extract from [The official documentation](https://webassembly.org/getting-started/js-api/):

>While there are future plans to allow WebAssembly modules to be loaded just like ES6 modules (...), WebAssembly must currently be loaded and compiled by JavaScript. For basic loading, there are three steps:
> 
> * Get the .wasm bytes into a typed array or ArrayBuffer
> * Compile the bytes into a WebAssembly.Module
> * Instantiate the WebAssembly.Module with imports to get the callable exports

Luckily, the Go authors made this task easy by providing a Javascript Loader. This loader is here `~/gowasm/misc/wasm/wasm_exec.js`. It comes with an HTML file that takes care of gluing everything in the browser.

To actually run our file, let's copy the following files in a directory and serve them by a webserver:

```
~ mkdir ~/wasmtest
~ cp ~/gowasm/misc/wasm/wasm_exec.js ~/wasmtest
~ cp ~/gowasm/misc/wasm/wasm_exec.html ~/wasmtest/index.html
~ cp example.wasm ~/wasmtest
```

Then edit the file `index.html` to run the correct sample:

{{< highlight js >}}
// ...
WebAssembly.instantiateStreaming(fetch("example.wasm"), go.importObject).then((result) => {
        mod = result.module;
        inst = result.instance;
        document.getElementById("runButton").disabled = false;
});
// ...
{{</ highlight >}}

In theory, any web server could do the job, but I had faced an issue when I tried to run it with `caddy`. The javascript loader is expecting the server to send the correct mime type for the wasm file.

Here is a quick hack to run our test: to write a go server with a particular handler for our wasm file.

{{< highlight js >}}
package main

import (
        "log"
        "net/http"
)

func wasmHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/wasm")
        http.ServeFile(w, r, "example.wasm")
}
func main() {
        mux := http.NewServeMux()
        mux.Handle("/", http.FileServer(http.Dir(".")))
        mux.HandleFunc("/example.wasm", wasmHandler)
        log.Fatal(http.ListenAndServe(":3000", mux))
}
{{</ highlight >}}

_Note_ Setting up a special router to handle all the wasm files is no big deal, but as I said, this is a POC and this post are side notes about it.

Then run the server with `go run server.go` and point your browser to [http://localhost:3000](http://localhost:3000).

Open the console, and _voilà_ !

# Interacting with the browser.

Let's interact with the world.

## Addressing the DOM

The `syscall/js` package contains the functions that allow interaction with the DOM through the javascript API.

To get the documentation about this package, just run:

`GOROOT=~/gowasm godoc -http=:6060`

and point your browser to [http://localhost:6060/pkg/syscall/js/](http://localhost:6060/pkg/syscall/js/).

Let's write a simple HTML file that displays an input field. Then, from the webassembly, let's place an event on this element and trigger an action when this event fires.

Edit the `index.html` and place this code just below the `run` button: 

{{< highlight html >}}
        <button onClick="run();" id="runButton" disabled>Run</button>
        <input type="number" id="myText" value="" />
</body>
{{</ highlight >}}

Then modify the Go file: 

{{< highlight go >}}
package main

import "fmt"

func main() {
          c := make(chan struct{}, 0)
         cb = js.NewCallback(func(args []js.Value) {
                  move := js.Global.Get("document").Call("getElementById", "myText").Get("value").Int()
                  fmt.Println(move)
          })
          js.Global.Get("document").Call("getElementById", "myText").Call("addEventListener", "input", cb)
          // The goal of the channel is to wait indefinitly
          // Otherwise, the main function ends and the wasm modules stops
          <-c
}
{{</ highlight >}}

Compile the file as you did before and refresh your browser...
Open the console and type a number in the input field.... _voilà_

## Exposing a function

This one is a bit trickier... I did not find any easy way to expose a Go function into the Javascript ecosystem.
What we need to do is to create a [`Callback`](http://localhost:6060/pkg/syscall/js/#Callback) Object in the Go file and assign it to a Javascript Object.

To get a result back, we cannot return a value to the callback and we are using a javascript object instead.

Here is the new Go code:

{{< highlight go >}}

package main
import (
        "syscall/js"
)

func main() {
        c := make(chan struct{}, 0)
        add := func(i []js.Value) {
                js.Global.Set("output", js.ValueOf(i[0].Int()+i[1].Int()))
        }
        js.Global.Set("add", js.NewCallback(add))
        <-c
}
{{</ highlight >}}

Now compile and run the code.
Open back your browser and open the console.

If you type `output` it should return a `Object not found`. Now type `add(2,3)` and type `output`... You should get `5`.

This is not the most elegant way to interact, but it is working as expected.

# Conclusion

The wasm support in Go is just starting but is in massive development. Many things that are working by now. I am even able to run a complete recurrent neural network coded thanks to Gorgonia directly in the browser. 
I will explain all of this later.

