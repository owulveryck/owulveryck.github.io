---
title: "Reading from the web offline and distraction-free"
date: 2021-10-07T10:07:41+02:00
draft: false
---

*TL;DR:* This article describes the wiring of a tool to turn a webpage into a self-sufficient epub (for reading offline). If you want to try the tool, you can grab a binary version from [GitHub](https://github.com/owulveryck/rePocketable/tags)

## The Why

To oversimplify my need, I will quote this from the _[Readability Project](https://web.archive.org/web/20150817073201/http://lab.arc90.com/2009/03/02/readability/)_

> Reading anything on the Internet has become a full-on nightmare. As media outlets attempt to eke out as much advertising revenue as possible, we’re left trying to put blinders on to mask away all the insanity that surrounds the content we’re trying to read.
>
> It’s almost like listening to talk radio, except the commercials play during the program in the background. It’s a pretty awful experience. Our friend to date has been the trusty “Print View” button. Click it and all the junk goes away. I click it all the time and rarely print. It’s really become the “Peace & Quiet” button for many.

### This article

In a recent post, I blogged about a tool I am building for my reMarkable.
In this post, I will describe a new tool that converts any webpage into an ePub file.

The goals of this tool are:

- to keep track of the articles I like without fearing any broken links
- to extract the content, and read the articles without distraction
- to be able to read the articles offline on devices such as ebook readers or my reMarkable

## Existent solutions

This feature already exists if you are using a Kobo and the getPocket service.
The problem is that it is that the offline experience is tidily linked with my Kobo device.
On top of that, getPocket does not offer any way to download the cleaned version of the articles.

We, as developers, have superpowers: we can build the tools we want.

Let's explain the features I am building step by step.

_Disclaimer_ at the time this post is written, the tool results from various experiments, but not the architecture or the code is clean and maintainable. 
Take this post as a validation of a proof of concept.

## First part: extracting the content

The most important part of this journey is the tool's ability to extract the content of a webpage.
The first idea would be to query the getPocket service that does this, but the [documentation of their API](https://getpocket.com/developer/docs/v3/article-view) mentions that:

> Pocket's Article View API will return article content and relevant meta data on any provided URL.
>
> The Pocket Article View API is currently only open to partners that are integrating Pocket specific features or full-fledged Pocket clients. For example, building a Pocket client for X platform.
>
> If you are looking for a general text parser or to provide "read now" functionality in your app - we do not currently support that. There are other companies/products that provide that type of API, for example: Diffbot.

They mention [Diffbot](https://www.diffbot.com/products/extract/), but it is a web service that requires a subscription; I'd like to build a simple tool, free of charge, for my usage, and therefore this is not an option.

### Readability / Arc90

I looked into open source initiatives that empower the reading modes of the browsers (I am/was a fan of the safari reading mode), and I found that some of them were based on an experiment made by [Arc90](https://web.archive.org/web/20150817073201/http://lab.arc90.com/2009/03/02/readability/).
This experiment led to the (discontinued) service [readability](https://en.wikipedia.org/wiki/Readability_(service)).

We can now find various implementations of the [Arc90 algorithm](https://github.com/masukomi/arc90-readability). I am using [this implementation](https://github.com/cixtor/readability) in Go for my tool.

#### Code

Feel free to skip this part if you are not interested in the code

The API of the readability library is straightforward.
First, there is a need to create a `Readability` object with an _HTML parser that reads and extracts relevant content_.

Then, calling the `Parse` method on this object, feeding it with an `io.Reader` that contains the page to analyze.

The result is an object of type `Article` that contains some metadata and the cleaned content. This content is an HTML tree and is accessible via a top-level [`html.Node`](https://pkg.go.dev/golang.org/x/net/html#Node).

{{< highlight go >}}
package main

import (
   "log"
   "net/http"
   "os"

   "github.com/cixtor/readability"
   "golang.org/x/net/html"
)

func main() {
   // create a parser
   htmlParser := readability.New()
   // Fetch a webpage
   resp, err := http.Get("https://example.com/")
   passOrDie(err)
   // Deal with errors etc...
   defer resp.Body.Close()
   // Parse the content
   article, err := htmlParser.Parse(resp.Body, "https://example.com")
   passOrDie(err)
   // Write the readable result on stdout
   html.Render(os.Stdout, article.Node)
}
{{< /highlight >}}

### The problem with reactive content and Medium articles

When the Arc90 project made this experiment, there were not many reactive contents.

On top of that, it does not handle the javascript.
This leads to images that are not correctly displayed. Let's take the [first chapter of Simon Wardley's book about maps](https://medium.com/wardleymaps/on-being-lost-2ef5f05eb1ec) to illustrate the problem.

The picture below is a screenshot of a reader view of the page with Safari:
{{< figure src="/assets/rePocketable/medium.png" title="The medium issue with Arc90" >}}

The code below is the code extracter by a curl request:
{{< highlight html >}}
<figure
   class="ja jb jc jd je jf cw cx paragraph-image">
   <div role="button" tabindex="0"
      class="jg jh ji jj aj jk">
      <div class="cw cx iz">
         <div class="jq s ji jr">
            <div
               class="js jt s">
               <div
                  class="jl jm t u v jn aj at jo jp">
                  <img alt=""
                     class="t u v jn aj ju jv jw"
                     src="https://miro.medium.com/max/60/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg?q=20"
                     width="700"
                     height="590"
                     role="presentation" />
               </div>
               <img alt=""
                  class="jl jm t u v jn aj c"
                  width="700"
                  height="590"
                  role="presentation" /><noscript><img
                     alt=""
                     class="t u v jn aj"
                     src="https://miro.medium.com/max/1400/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg"
                     width="700"
                     height="590"
                     srcSet="https://miro.medium.com/max/552/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg 276w, https://miro.medium.com/max/1104/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg 552w, https://miro.medium.com/max/1280/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg 640w, https://miro.medium.com/max/1400/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg 700w"
                     sizes="700px"
                     role="presentation" /></noscript>
            </div>
         </div>
      </div>
   </div>
</figure>
{{< /highlight >}}

Within the `<figure>` element, we can see that the first image (https://miro.medium.com/max/60/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg?q=20) is a thumbnail and it acts as a placeholder.

A couple of JavaScript routines replaces the image at rendering time in the browser.
Luckily a `<noscript>` tag is also present and exposes the complete sources of the image.

As the Arc90 library removes all the `<noscript>` elements, the only options are:

- to pre-process the HTML file before feeding the Arc90 algorithm
- to amend the Arc90 library

So far, the behavior we are addressing seems particular to articles hosted on medium. Amending the Arc90 Algo to handle this specific use-case does not seem to be a good idea.

So let's go for a pre-processing step of the document before feeding the Arc90 algo.
It is beyond the scope of this article to show and comment on the complete code to do that.

In a glimpse, the HTML content is extracted into a tree of `*html.Node` elements; then, the processing step walks the tree via a recursive function seeking `figure` elements.

{{< highlight go >}}
func preProcess(n *html.Node) error {
   if n.Type == html.ElementNode && n.Data == "figure" {
       err := processFigure(n)
      // if error, return error
   }
   for c := n.FirstChild; c != nil; c = c.NextSibling {
      err := preProcess(c)
      // if error, return error
   }
   return nil
}
{{< /highlight >}}

Then, within the `processFigure`, we once again walk through the subtree, seeking the primary `img` node, and replacing its attributes with those from `noscript/img` node.

You can find a complete code in this [_gist_](https://gist.github.com/owulveryck/5f9a07762ce40e6f6d9028e76bd798e2)

Once the HTML tree is adapted, it can go through the Arc90 Algorithm.

_Note_: as of today, the tree is rendered into HTML to match the API of Arc90. This is unoptimized. I will eventually submit a PR or fork the project to add a new API that applies the Acr90 Algo directly to an `*html.Node`.

## Second part: generating the ePub

Now that we have proper content, let's turn it into an ePub.

An ePub is a set of XHTML files carrying content, along with images and local files. All of the content is self-sufficient and packaged in a zip file.

To generate the ePub in the tool, I rely on the [`go-epub`](https://github.com/bmaupin/go-epub) library. This library is stable, and the author welcomes contributions.

The ePub generation is made in two steps:

1. building an Epub structure holding the content of the epub;
2. generating the epub file with self-sufficient content.

### First step: crafting the ePub

In the first step, we create the HTML content. The content is the HTML tree processed previously by the Arc90 algorithm.
The content is added as a single section in the ePub for commodity. A better way would be to parse the HTML tree and create a section for each `h1` tag.
But as the target is downloading a single page, there should typically be a single `h1` tag inside de page.

To be self-sufficient, there is a need to parse this tree, seeking remote content (in essence, the images) and downloading it locally.

The go-epub library provides a set of methods to handle the content to do this task smoothly. The [`AddImage`](https://pkg.go.dev/github.com/bmaupin/go-epub?utm_source=godoc#Epub.AddImage) method, for example, creates an entry in a map that references online content and provides a reference to a local file.

This code, from the doc, shows how it works:

{{< highlight go >}}
func Example() {
    e := epub.NewEpub("My title")

    // Add an image from a URL. The filename is optional
    imgPath, _ := e.AddImage("https://golang.org/doc/gopher/gophercolor16x16.png", "")

    fmt.Println(imgPath)
    // Output:
    // ../images/gophercolor16x16.png
}
{{< /highlight >}}

We need to call this method for every image element in order to populate the image map. On top of that, every `src` attribute must be changed to use the local file.
We use the same system as before and use a recusrive function applied to the root node of the HTML tree:

{{< highlight go >}}
func (d *Document) replaceImages(n *html.Node) error {
    if n.Type == html.ElementNode && n.Data == "img" {
        // find the URL of the image from the current node
        val, f, err := getURL(n.Attr)
        // error checking
        for i, a := range n.Attr {
            if a.Key == "src" {
                img, err = d.AddImage(val, "")
                // error checking
                // Add the local image name as the src attribute of the image
                n.Attr[i].Val = img
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        err := d.replaceImages(c)
        // error checking
    }
    return nil
}
{{< /highlight >}}


#### Back to _Medium_'s image problem

We addressed the JavaScript issue in the preprocessing step. Let's now address the reactive problem.
Actually, the `img` source we have set in the HTML tree relies on the `srcset` attribute.

In the `getURL` function, we will implement a logic that will set the default source value present in the `src` attribute.
If it finds a `srcset` attribute, it will parse it and sort it, so the first element holds the largest picture (we want the best possible resolution).

We implement the `sort.Sort` interface on a newly created structure `[]srcSetElements`:

{{< highlight go >}}
type srcSetElement struct {
    url            string
    intrinsicWidth string
}

type srcSetElements []srcSetElement

func (s srcSetElements) Len() int { ... }
func (s srcSetElements) Less(i int, j int) bool { ... }
func (s srcSetElements) Swap(i int, j int) { ... }
{{< /highlight >}}

I will not display the whole getURL function as its implementation is straightforward and present on the project's GitHub.


### Second step: creating the ePub

Now the structure of the epub is correct, simply call the [`Write`](https://pkg.go.dev/github.com/bmaupin/go-epub#Epub.Write) method that will:

- download the assets listed in the Epub structure;
- add some metadata;
- create the zip file.

This method ends the process and produces the expected ePub file.

## Third part: adding fancy features

Now we have an epub file, let's add some features to improve the reader experience.

### Grabbing meta information

The `Article` structure produced by the Arc90 parser references a title, an author, and a front cover for the site.
But, as explained before, Arc90 is quite old, and those pieces of informations are provided nowadays by OpenGraph elements.

Arc90 cleans those elements; therefore, we will grab them in the pre-processing step.

We rely on the [`opengraph`](github.com/dyatlov/go-opengraph/opengraph) library in Go to create a `getOpenGraph` function. 
The opengraph's entry point reads the content from an `io.Reader`. 
To optimize the memory, we will implement the `getOpenGraph` method as a middleware.

It will read the HTML file from the io.Reader, process it, and `Tee` the original into another readeri thanks to an `io.TeeReader`.
The signature of the method is:

{{< highlight go >}}
func getOpenGraph(r io.Reader) (*opengraph.OpenGraph, io.Reader) { ... }
{{< /highlight >}}

Once again, [the complete code is available on the GitHub repository of the project](https://github.com/owulveryck/rePocketable/blob/8060aa3709b89c6bdf8bf6010027dd38bccd47d7/internal/epub/opengraph.go#L130).

### Generating a cover

Now that we have some information, we can generate a cover for the ePub.
A cover is an XHTML file that references a single picture.

On the picture, we would like to see:

- the front image of the article as displayed on the social media;
- the title of the article;
- the author of the article;
- the origin of the article;

With the `image/draw` package of the standard library, we create an RGB image and compose the front cover.

The code of the cover generation is [here](https://github.com/owulveryck/rePocketable/blob/master/internal/epub/cover.go).
Then, the methods of the go-epub library add it to the ePub.

### GetPocket integration

To complete the work, we can create a GetPocket integration to grab all he elements from the GetPocket reading list and convert them to ePub.
The integration is straightforward as the API of getPocket allows retrieving a structure holding:

- the original URL
- the title of the file
- the front image
- the authors

But, a target could be to run a daemon on the eReader (for example, a reMarkable);
therefore, the internal library is handling a daemon mode to fetch the articles on a regularly (as well as when the device wakes up).

If you are curious, the mechanism is implemented in [a pocket package](https://github.com/owulveryck/rePocketable/tree/master/internal/pocket) and uses
the mechanism I implemented a while ago to hack the [remarkable_news](https://github.com/owulveryck/remarkable_news) project.

### Dealing with MathJax

Another feature that is missing from the getPocket integration on my kindle is the ability to render LaTeX formulas.
I add one more processing step to find a mathjax content, and create a png image of the formula.

To do that, I use the [github.com/go-latex/latex](github.com/go-latex/latex) package.

The principle is to find a TextNode holding a MathJax element thanks to a regular expression:

{{< highlight go >}}
var mathJax = regexp.MustCompile(`\$\$[^\$]+\$\$`)

func hasMathJax(n *html.Node) bool {
    return len(mathJax.FindAllString(n.Data, -1)) > 0
}

func preProcess(n *html.Node) error {
   // ...
    case n.Type == html.TextNode && hasMathJax(n):
        processMathTex(n)
    }
   // ...
}
{{< /highlight >}}

then the `processMathTex` function analyzes the formulae and renders them into a png encoded file. Then the file is inserted in the HTML tree in an `img` tag. The `src` attribute references an inline content of the formula, encoded with the [dataURL principle](https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/Data_URIs).

## Conclusion and future work

I don’t use the getPocket integration very often, but I use the `toEpub` tool to convert a web page daily.

The getPocket integration will be helpful once I have encoded the output file to a format suitable for the remarkable. It sounds pretty straightforward, but I have not taken the time to do it yet.

So far, my workflow is:

- grabbing the URL on my laptop
- running the toEpub locally
- sending the result to the remarkable with `rmapi` (and now gdrive)

The problem is that it requires a laptop and the tool installed on it. I am currently hacking the go-epub library, so it will no longer need a filesystem, allowing a compilation into webassembly to ease the deployment.

Stay tuned.