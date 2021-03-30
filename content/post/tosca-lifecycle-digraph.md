---
author: Olivier Wulveryck
date: 2015-11-20T10:09:30Z
description: A tosca lifecycle represented as a digraph
draft: false
tags:
- TOSCA
- Digraph
- Graph Theory
- golang
title: TOSCA lifecycle as a digraph
topics:
- TOSCA
type: post
---

# About TOSCA

The [TOSCA](https://www.oasis-open.org/committees/tc_home.php?wg_abbrev=tosca) acronym stands for 
_Topology and Orchestration Specification for Cloud Applications_. It's an [OASIS](https://www.oasis-open.org) standard.

The purpose of the TOSCA project is to represent an application by its topology and formalize it using the TOSCA grammar.

The [[TOSCA-Simple-Profile-YAML-v1.0]](http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.0/csprd01/TOSCA-Simp$le-Profile-YAML-v1.0-csprd01.html) 
current specification in YAML introduces the following concepts.

> - TOSCA YAML service template: A YAML document artifact containing a (TOSCA) service template that represents a Cloud application.
> - TOSCA processor: An engine or tool that is capable of parsing and interpreting a TOSCA YAML service template for a particular purpose. For example, the purpose could be validation, translation or visual rendering.
> - TOSCA orchestrator (also called orchestration engine): A TOSCA processor that interprets a TOSCA YAML service template then instantiates and deploys the described application in a Cloud.
> - TOSCA generator: A tool that generates a TOSCA YAML service template. An example of generator is a modeling tool capable of generating or editing a TOSCA YAML service template (often such a tool would also be a TOSCA processor).
> - TOSCA archive (or TOSCA Cloud Service Archive, or “CSAR”): a package artifact that contains a TOSCA YAML service template and other artifacts usable by a TOSCA orchestrator to deploy an application.

## My work with TOSCA

I do believe that TOSCA may be a very good leverage to port a "legacy application" (aka _born in the datacenter_ application) into a cloud ready application without rewriting it completely to be cloud compliant.
To be clear, It may act on the hosting and execution plan of the application, and not on the application itself.

A single wordpress installation in a TOSCA way as written [here](http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.0/csprd01/TOSCA-Simple-Profile-YAML-v1.0-csprd01.html#_Toc430015847) is represented like that

<img class="img-square img-responsive" src="http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.0/csprd01/TOSCA-Simple-Profile-YAML-v1.0-csprd01_files/image035.png" alt="Single Wordpress representation"/>

While I was learnig GO, I have developped a [TOSCA lib](https://github.com/owulveryck/toscalib) and a [TOSCA processor](https://github.com/owulveryck/toscaviewer) which are, by far, not _idiomatic GO_.

Here are two screenshots of the rendering in a web page made with my tool (and the graphviz product):

<hr/>
*The graph representation of a _Single instance wordpress_*
<img class="img-responsive" src="/assets/images/toscaviewer_template_def.png" alt="Tosca view ofthe single instance wordpress"/>


*The graph representation of a lifecycle of _Single instance wordpress_*
<img class="img-responsive" src="/assets/images/toscaviewer_lifecycle_def.png" alt="Lifecycle representation of the single wordpress instance representation"/>
<hr/>

The TOSCA file is parsed with the help of the `TOSCALIB` and then it fills an adjacency matrix (see [FillAdjacencyMatrix](https://godoc.org/github.com/owulveryck/toscalib#ToscaDefinition.FillAdjacencyMatrix))

The [graphviz](http://graphviz.org) take care of the (di)graph representation.

What I would like to do now, is a little bit more: I would like to play with the graph and query it
Then I should perform requests on this graph. For example I could ask:

* _What are the steps to go from the state Initial of the application, to the state running_ ?
* _What are the steps to go from stop to delete_
* ...

and that would be **the premise of a TOSCA orchestrator**.

## The digraph go code

I've recently discoverd the [digraph](https://github.com/golang/tools/tree/master/cmd/digraph) tool, that I will use for querying the graphs.
The `digraph` is represented as a map with a node as a key and its immediates successors as values:

```go
// A graph maps nodes to the non-nil set of their immediate successors.
type graph map[string]nodeset

type nodeset map[string]bool
```

## From TOSCA to digraph

What I must do is to parse the adjacency matrix, get the "lifecycle action" related to the id and fill the graph g.

# Let's go 

Considering the digraph code, what I need to do is simply to override the `parse` method.

## Principle

I will fill the `graph` with a string composed of _nodename:action_ as key.
For example, if I need to do a "Configure" action of node "A" after a "Start" action on node "B", I will have the following entry in the map:

```go
g["B:Start"] = "A:Configure"
```

So What I need to do is to parse the adjjacency matrix, do a matching with the row id and the "node:action" name, and fill the `graph g` with the matching of the corresponding "node:action".

I will fill a `map` with the id of the node:action as key and the corresponding label as values:
```gloang
for node, template := range toscaTemplate.TopologyTemplate.NodeTemplates {
        ids[template.GetConfigureIndex()] = fmt.Sprintf("%v:Configure", node)
        ids[template.GetCreateIndex()] = fmt.Sprintf("%v:Create", node)
        ids[template.GetDeleteIndex()] = fmt.Sprintf("%v:Delete", node)
        ids[template.GetInitialIndex()] = fmt.Sprintf("%v:Initial", node)
        ids[template.GetPostConfigureSourceIndex()] = fmt.Sprintf("%v:PostConfigureSource", node)
        ids[template.GetPostConfigureTargetIndex()] = fmt.Sprintf("%v:PostconfigureTarget", node)
        ids[template.GetPreConfigureSourceIndex()] = fmt.Sprintf("%v:PreConfigureSource", node)
        ids[template.GetPreConfigureTargetIndex()] = fmt.Sprintf("%v:PreConfigureTarget", node)
        ids[template.GetStartIndex()] = fmt.Sprintf("%v:Start", node)
        ids[template.GetStopIndex()] = fmt.Sprintf("%v:Stop", node)
}
```

Then I can easily fill the `graph g` from the adjacency matrix:

```gloang
row, col := toscaTemplate.AdjacencyMatrix.Dims()
        for r := 1; r < row; r++ {
                for c := 1; c < col; c++ {
                        if adjacencyMatrix.At(r, c) == 1 {
                                g.addEdges(ids[r], ids[c])
                        }
                }
        }
```

That's it

# The final function

Here is the final parse function
```go
func parse(rd io.Reader) (graph, error) {
        g := make(graph)
        // Parse the input graph.
        var toscaTemplate toscalib.ToscaDefinition
        err := toscaTemplate.Parse(rd)
        if err != nil {
                return nil, err
        }
        // a map containing the ID and the corresponding action
        ids := make(map[int]string)
        // Fill in the graph with the toscaTemplate via the adjacency matrix
        for node, template := range toscaTemplate.TopologyTemplate.NodeTemplates {
                // Find the edges of the current node and add them to the graph

                ids[template.GetConfigureIndex()] = fmt.Sprintf("%v:Configure", node)
                ids[template.GetCreateIndex()] = fmt.Sprintf("%v:Create", node)
                ids[template.GetDeleteIndex()] = fmt.Sprintf("%v:Delete", node)
                ids[template.GetInitialIndex()] = fmt.Sprintf("%v:Initial", node)
                ids[template.GetPostConfigureSourceIndex()] = fmt.Sprintf("%v:PostConfigureSource", node)
                ids[template.GetPostConfigureTargetIndex()] = fmt.Sprintf("%v:PostconfigureTarget", node)
                ids[template.GetPreConfigureSourceIndex()] = fmt.Sprintf("%v:PreConfigureSource", node)
                ids[template.GetPreConfigureTargetIndex()] = fmt.Sprintf("%v:PreConfigureTarget", node)
                ids[template.GetStartIndex()] = fmt.Sprintf("%v:Start", node)
                ids[template.GetStopIndex()] = fmt.Sprintf("%v:Stop", node)
        }

        row, col := toscaTemplate.AdjacencyMatrix.Dims()
        for r := 1; r < row; r++ {
                for c := 1; c < col; c++ {
                        if adjacencyMatrix.At(r, c) == 1 {
                                g.addEdges(ids[r], ids[c])
                        }
                }
        }
        return g, nil
}

```
# Grab the source and compile it

I have a github repo with the source.
It is go-gettable 
```
go get github.com/owulveryck/digraph
cd $GOPATH/src/github.com/owulveryck/digraph && go build
```

**EDIT** As I continue to work on this tool, I have created a "blog" branch in the github which holds the version related to this post

# Example

I will use the the same example as described below: the single instance wordpress.

I've extracted the YAML and placed in in the file [tosca_single_instance_wordpress.yaml](https://github.com/owulveryck/toscaviewer/blob/master/examples/tosca_single_instance_wordpress.yaml).

Let's query the nodes first:
```sh
curl -s https://raw.githubusercontent.com/owulveryck/toscaviewer/master/examples/tosca_single_instance_wordpress.yaml | ./digraph nodes
mysql_database:Configure
mysql_database:Create
mysql_database:Start
mysql_dbms:Configure
mysql_dbms:Create
mysql_dbms:Start
server:Configure
server:Create
server:Start
webserver:Configure
webserver:Create
webserver:Start
wordpress:Configure
wordpress:Create
wordpress:Start

```

so far, so good...

Now, I can I go from a `Server:Create` to a running instance `wordpress:Start`

```
curl -s https://raw.githubusercontent.com/owulveryck/toscaviewer/master/examples/tosca_single_instance_wordpress.yaml | ./digraph somepath server:Create wordpress:Start
server:Create
server:Configure
server:Start
mysql_dbms:Create
mysql_dbms:Configure
mysql_dbms:Start
mysql_database:Create
mysql_database:Configure
mysql_database:Start
wordpress:Create
wordpress:Configure
wordpress:Start
```

Cool!

# Conclusion

The tool sounds ok. What I may add:

- a command to display the full lifecycle (finding the entry and the exit points in the matrix and call somepath with it)
- get the tosca `artifacts` and display them instead of the label to generate a deployment plan
- execute the command in `goroutines` to make them concurrent


And of course validate any other TOSCA definition to go through a bug hunting party
