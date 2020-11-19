---
title: "Ontology, graphs and turtles - Part II"
date: 2020-11-17T17:07:03+01:00
draft: true
tags: ["ontology", "taxonomy", "graph", "turtle", "rdf", "golang", "go"]
summary: "This article is about parsing and extracting the knowledge from a triplestore to create a graph in-memory in Go."
---

In a [previous article](/2020/11/13/ontology-graphs-and-turtles-part-i.html), I introduced the notion of ontology and knowledge graph.

Let's go deeper into the concept and apply some technology to create an actual graph structure and eventually play with it.

At the end of this article, we well have parsed a triplestore in turtle format and created a graph structure in Go (based on [gonum's interface](https://pkg.go.dev/gonum.org/v1/gonum/graph))

## The triplestore

As explained in the last post, our knowledge database is a triple store. As a matter of example, I will rely on the ontology exposed by [schema.org](https://schema.org).

Schema.org is a collaborative, community activity with a mission to create, maintain, and promote schemas for structured data on the Internet, on web pages, in email messages, and beyond. Founded by Google, Microsoft, Yahoo and Yandex, Schema.org vocabularies are developed by an open community process [...].
You can grab the complete ontology with this command:

The complete definition is available in turtle format and can be downloaded easily:

```shell
curl -O -s https://schema.org/version/latest/schemaorg-current-http.ttl
```

### Parsing the store

#### Basic explanation of the parser

To parse the store, I am using the package [`gon3`](https://github.com/deiu/gon3) from [Andrei Sambra](https://twitter.com/andreisambra). Even though there is no license attached, Andrei allowed me to use it and to modify it for non-profit code.

I forked the repo to make some minor adjustments I needed for my experiments.

The entry point of the package is a `Parser` structure. Its purpose is to read a byte stream (`io.Reader`) and turn it into a functional structure called `Graph`. The `Graph` structure within the package is not representing all the edges. Still, it consists of an array of Triples (aka rdf graph):

{{<highlight go>}}
// An RDF graph is a set of RDF triples
type Graph struct {
	triples []*Triple
	uri     *IRI
}
{{</ highlight >}}

A `Triple` is a structure holding three `Term`. One is the subject, one is a predicate, and the last one is the object.

{{<highlight go>}}
// see http://www.w3.org/TR/rdf11-concepts/#dfn-rdf-triple
type Triple struct {
	Subject, Predicate, Object Term
}
{{</ highlight >}}

In the previous article, we saw that a term in RDF can be expressed in different types. As of today, the way to represent generic types in Go is to use interfaces. Therefore, a `Term` has an `interface` based definition:

{{<highlight go>}}
type Term interface {
	String() string
	Equals(Term) bool
	RawValue() string
}
{{</ highlight >}}

Two important objects are implementing the Term interface:

- IRI
- Literal

#### Generate the RDF Graph

If we glue all the concepts together we have the possibility to create a basic structure:

{{<highlight go>}}
import "github.com/owulveryck/gon3" // Other imports omited for brevity

func main() {
        baseURI := "https://example.org/foo"
        parser := gon3.NewParser(baseURI)
        gr, _ := parser.Parse(os.Stdin) // Error handling is omited for brevity
        fmt.Printf("graph contains %v triples", len(gr.Triples()))
}
{{</ highlight >}}

Then we can test the parsing with the file we downloaded from schema.org previously:

```shell
> cat schemaorg-current-http.ttl| go run main.go
graph contains 15323 triples
```

We can check that roughly the number of triple matches what's expected by counting the rdf separators from the file:

```shell
> cat schemaorg-current-http.ttl | egrep -v '^@|^$' | egrep -c ' \.$| \;$|\,$'
15329
```

The numbers are not the same but they are alike (the grep command does not evaluate the litteral and some punctuation may be wrongly counted)

### Understanding the graph structure

Schema.org is using a lot of prefixes. Amongst all of those, a couple is interesting for the sake of this article:

```turtle
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix schema: <http://schema.org/> .
```

https://schema.org/PostalAddress


## Generating a Graph

## Putting all together