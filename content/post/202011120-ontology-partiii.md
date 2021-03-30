---
title: "Ontology, graphs and turtles - Part III"
date: 2020-11-20T20:07:03+01:00
draft: false
tags: ["ontology", "taxonomy", "graph", "turtle", "rdf", "golang", "go"]
summary: "This article is about using the graph in-memory from the build-in part II. In this article, I will show how to extract the graph's information using a template engine. Eventually, we will build a documentation webserver that looks like schema.org."
---

In [a recent article](/2020/11/13/ontology-graphs-and-turtles-part-i.html), I introduced the notion of ontology and knowledge graph.

Then I exposed in this [post](/2020/11/17/ontology-graphs-and-turtles-part-ii.html) how to parse a triplestore to create a directed graph in memory.

This post will show you how to manipulate the graph. The use-case is a representation of the graph's information thanks to a template engine.

The engine should be extensible enough to generate any text-based report (HTML, markdown, …)

At the end of the post, we eventually build a basic webserver that similarly presents the information of [schema.org](schema.org)'s ontology.

| schema.org | our implementation (localhost) |
|-|-|
| ![](/assets/ontology/schemaorg1.png) | ![](/assets/ontology/schemaorg_olwu.png) | 

**Caution:** the solution is a proof of concept, its implementation works but is not bulletproof. Some tests would make it safer to use, and a TDD would influence the package's design. Consider it as a code for educational purposes only.


## The template engine

From the documentation of the `text/template` package of Go's standard library:

> Package template implements data-driven templates for generating textual output.

It works by ...

> ... applying [templates] to a data structure. Annotations in the template refer to elements of the data structure (typically a field of a struct or a key in a map) to control execution and derive values to be displayed. Execution of the template walks the structure and sets the cursor, represented by a period '.' and called "dot", to the value at the current location in the structure as execution proceeds.

At first, this sounds like a plan for the [gnomes of southpark](https://en.wikipedia.org/wiki/Gnomes_(South_Park)). [^1]

[^1]:  Phase 1: Collect underpants / Phase 2: ? / Phase 3: Profit

Let's rephrase it. In essence, the workflow is:

- **Collect the data** (this has been done in the previous post)
- create a data structure holding the elements you want to represent ;
- create a skeleton of the representation with placeholders for the values of the structure you expect to see ;
- apply the template to the data structure ;
- **profit!**

### The data structure

First, we need a data structure.

The data structure will hold the information we will represent via the application of the template.

We will place ourselves in the context of the node and [_think like a vertex_](/2019/10/14/think-like-a-vertex-using-gos-concurrency-for-graph-computation.html).

So the structure we are setting starts by referencing the node:

{{< highlight go >}}
// Current object to apply to the template
type Current struct {
    Node  *graph.Node
}
{{</highlight>}}

Then within the template we can access the exported fields and methods of the `graph.Node` structure defined in the previous article.

{{<highlight go>}}
type Node struct {
    id              int64
    Subject         rdf.Term
    PredicateObject map[rdf.Term][]rdf.Term
}
{{</highlight>}}

For example, this template will display the output to a call to the `RawValue()` method of the `Subject` field:

{{<highlight go-text-template >}}
The subject is {{ .Node.Subject.RawValue }}
{{</highlight>}}

#### Full example

As a proof of concept, we can write a simple Go example that displays a value from a basic ontology:

{{<highlight go>}}
func Example() {
	const ontology = `
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix example: <http://example.com/> .

example:PostalAddress a rdfs:Class ;
    rdfs:label "PostalAddress" ;
    rdfs:comment "The mailing address." .
`
	const templ = `
The subject is {{ .Node.Subject.RawValue }}
	`
	type Current struct {
		Node *graph.Node
	}
	parser := rdf.NewParser("http://example.org")
	gr, _ := parser.Parse(strings.NewReader(ontology))
	g := graph.NewGraph(gr)

	postalAddress := g.Dict["http://example.com/PostalAddress"]
	node := g.Reference[postalAddress]

	myTemplate, _ := template.New("global").Parse(templ)
	myTemplate.Execute(os.Stdout, Current{node})

	// output:
	// The subject is http://example.com/PostalAddress
}
{{</highlight>}}

Highlight:

- On line 2, we define a simple ontology as a string constant;
- on line 11, we define a simple template to display the subject of a node;
- on lines 21/22, we find the node reference that holds the subject "PostalAddress";
- on line 25, we apply the template on the object "Current" created ad-hoc, which contains the targeted node.

On the same principles, we could display the predicates' content and the objects of the node structure. Those are arrays (slices) and maps. Therefore, we have to use the builtin actions of the template engine to access the data elements (cf [the action paragraph in the template documentation](https://golang.org/pkg/text/template/#hdr-Actions) for more info)

ex:

{{<highlight go-text-template >}}
The subject is {{ .Node.Subject.RawValue }}
and the list of preficates are:
{{ range $predicate, $objects := .Node.PredicateObject -}}
* {{ $predicate }}
  - {{ range $objects }}{{ .RawValue }}{{ end }}
{{ end -}}
{{</highlight>}}

which gives:

```text
The subject is http://example.com/PostalAddress
and the list of preficates are:
* <http://www.w3.org/1999/02/22-rdf-syntax-ns#type>
  - http://www.w3.org/2000/01/rdf-schema#Class
* <http://www.w3.org/2000/01/rdf-schema#label>
  - PostalAddress
* <http://www.w3.org/2000/01/rdf-schema#comment>
  - The mailing address.
```

### Going further

This system is good enough to display a simple structure. But, we are in the context of ontology and graphs. Therefore it is essential to be able to walk the graph and to display siblings and edges.

To achieve this, we will add a pointer to the graph structure itself in the data provider, that is, the "`Current`" object:

{{<highlight go>}}
// Current object to apply to the template
type Current struct {
    Graph *graph.Graph
    Node  *graph.Node
}
{{</highlight>}}

The we add a some pyntactic sugar to our structure:

{{<highlight go>}}
// To the node with edge holding a value from the links array
func (g Current) To(links ...string) []Current { // ... }

func (g Current) From(links ...string) []Current { // ... }

func (g Current) HasPredicate(predicate, object string) bool { // ... }
{{</highlight>}}

_Note_ The implementation of the methods and the documentation is isolated into a template package. The reference of this package can be found on [pkg.go.dev/github.com/owulveryck/rdf2graph/template](https://pkg.go.dev/github.com/owulveryck/rdf2graph/template).

As an example, let's complete our sample triplestore:
{{<highlight turtle >}}
example:PostalAddress a rdfs:Class ;
    rdfs:label "PostalAddress" ;
	rdfs:comment "The mailing address." .
	
example:addressCountry a rdf:Property ;
	rdfs:label "addressCountry" ;
	rdfs:domain example:PostalAddress ;
	rdfs:comment "A comment" .
	
example:address a rdf:Property ;
	rdfs:label "address" ;
	rdfs:domain example:PostalAddress ;
	rdfs:comment "Physical address of the item." .
{{</highlight>}}

And extend the template with a call to the `To` function we've created.

{{<highlight go-text-template >}}
{{ range $current := .To }}
The subject is {{ .Node.Subject.RawValue }}
and the list of preficates are:
{{ range $predicate, $objects := .Node.PredicateObject -}}
- {{ $predicate }}
  - {{ range $objects}}{{.RawValue}}{{end}}
{{ end -}}
{{ end -}}
{{</highlight>}}

Then, without modifying the rest of the code we've exposed in the previous paragraph, the execution of the example gives the following result:

```text
The subject is http://example.com/address
and the list of preficates are:
* <http://www.w3.org/1999/02/22-rdf-syntax-ns#type>
  - http://www.w3.org/1999/02/22-rdf-syntax-ns#Property
* <http://www.w3.org/2000/01/rdf-schema#label>
  - address
* <http://www.w3.org/2000/01/rdf-schema#comment>
  - Physical address of the item.
* <http://www.w3.org/2000/01/rdf-schema#domain>
  - http://example.com/PostalAddress

The subject is http://example.com/addressCountry
and the list of preficates are:
* <http://www.w3.org/1999/02/22-rdf-syntax-ns#type>
  - http://www.w3.org/1999/02/22-rdf-syntax-ns#Property
* <http://www.w3.org/2000/01/rdf-schema#label>
  - addressCountry
* <http://www.w3.org/2000/01/rdf-schema#comment>
  - A comment
* <http://www.w3.org/2000/01/rdf-schema#domain>
  - http://example.com/PostalAddress
```

## A simple web service

Now that we have built of the tools we need to render the graph, let's build a very simple webserver to present the information of the knowledge graph. As explained in introduction, we will use the ontology of schema.org as a database (the creation of the knowledge graph is explained in the previous article).

### Creating the template

Each representation of a node is a single html page. It is accessed through a call to "http://serviceurl/NodeSubject".

The skeleton of the page is a template. 
To make things easier, we split the tamplate into three subtemplates.

- a `main` template that will create the HTML structure and the outside table
- a template to display a `class` as a `tbody` structure
- a property template to display a line of the tbody structure

{{<highlight go-text-template >}}
{{ define "main" }}
<!DOCTYPE html>
<!-- boilerplate of the HTML file -->
    <table class="blueTable">
        <!-- header of the table -->
        {{ template "rdfs:type rdfs:Class" . }}
    </table>
</html>
{{ end }}

{{ define "rdfs:type rdfs:Property" }}
<tr>
{{ calls to display the subjects and predicates }}
</tr>
{{ end }}

{{ define "rdfs:type rdfs:Class" }}
<tbody>
    <tr>
    <!-- The rest of the table structure -->
    {{ range over the "To" nodes for the graph held in `Current` }}
        {{ for each node, if its type is "property" }}
            {{ template "rdfs:type rdfs:Property" . -}}
    {{ range over the "From" nodes for the graph held in `Current` }}
        {{ for each node, if its type is "class" }}
                {{ template "rdfs:type rdfs:Class" . -}}
    </tr>
</tbody>
{{ end }}
{{</highlight>}}

There is not much interest in displaying all the wiring inside this article.
The full template is available [here](https://github.com/owulveryck/rdf2graph/blob/7a6127ae4428c5501a1d07eb541a16fb4ee3ad83/examples/webview/index.tmpl#L1-L59)

This will display a nice formated page when mixed with a node of the graph.

#### The web server

To make things easier, let's encapsulate all of this inside a simple webserver. The goal is to be able to display any node of the graph when accessed through a URL.

First, we create a structure that will implement the [http.Handler](https://golang.org/pkg/net/http/#Handler) interface.

For convenience, this structure carries the graph, the template, and a hashmap of the rdf.IRI. The later is used to shorten the URLs (calling _http://myservice/blabla_ instead of _http://myservice/http://example.com/blabla_)

{{<highlight go >}}
type handler struct {
	namespaces map[string]*rdf.IRI
	g          graph.Graph
	tmpl       *template.Template
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
{{</highlight >}}

The ServeHTTP method is composed of three parts:

- extract the node from the URL
- check if the node exists
- apply the template and write the result on the _ResponseWriter_

I won't detail all the code to implement this. You can look at the implementation [on GitHub](https://github.com/owulveryck/rdf2graph/blob/7a6127ae4428c5501a1d07eb541a16fb4ee3ad83/examples/webview/main.go#L60-L77).

We need to glue all the code from the articles together to:

- parse a triplestore from an `io.Reader`;
- create a graph in-memory;
- read the template file and create the Go template;
- create the handler to reference the graph and the template;
- register the handler to serve an HTTP request;

Sounds tricky, but it is reasonably easy and straightforward if you do a little bit of Go. Anyhow, a sample implementation is on [GitHub](https://github.com/owulveryck/rdf2graph/tree/main/examples/webview).To launch it, simply do:

```bash
curl -s https://schema.org/version/latest/schemaorg-current-http.ttl | go run main.go
```

Then point your browser to "http://localhost:8080/PostalAddress" and you should get something that looks like this:

![](/assets/ontology/schemaorg_olwu.png)

## Conclusion

This was the last article about ontology. Through the pieces, we’ve discovered a way to describe a knowledge graph to be easy to write for a human and efficient enough to parse for a machine. Then we’ve built a graph in-memory and exploited this graph to represent the knowledge. The representation layer can be seen as a projection layer that exposes the information required for a specific functional domain.

Now let's have fun with ontology to:

- document functional areas and act as a shared ubiquitous language
- provide a map of information to locate and use a specific part of the knowledge of a system
