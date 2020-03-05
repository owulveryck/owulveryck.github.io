---
categories:
- category
date: 2017-03-22T09:15:35+01:00
description: "My first experience with GraphQL. I will try to see how it fits the pricing model of AWS as described in an earlier post."
draft: false
images:
- /assets/images/graphqllogo.png
tags:
- golang
- aws
- graphql
title: Playing with Facebook's GraphQL (applied to AWS products and offers management)
---

# About GraphQL

GraphQL has been invented by Facebook for the purpose of refactoring their mobile application. Facebook had reached the limits of the standard REST API mainly because:

* Getting that much information was requiring a huge amount of API endpoints
* The versioning of the API was counter-productive regarding Facebook's frequents deployements.

But graphql is not only a query language related to Facebook. GraphQL is not only applicable to social data. 

Of course it is about graphs and graphs represents relationships. But you can represent relationships in all of your business objects.

Actually, GraphQL is all about your application data.

In this post I will try to take a concrete use case. I will first describe the business objects as a graph, then I will try to implement a schema with GraphQL. At the very end I will develop a small GraphQL endpoint to test the use case.

__Caution__ _I am discovering GraphQL on my own. This post reflects my own work and some stuff may be inaccurate or not idiomatic._

## The use case: AWS billing

Let's take a concrete example of a graph representation. Let's imagine that we are selling products related to Infrastructre as a Service (_IaaS_). 

For the purpose of this post, I will use the AWS data model because it is publicly available and I have already blogged about it.
We are dealing with products families, products, offers and prices.

In (a relative) proper english, let's write down a description of the relationships:

* Products
  * A product family is composed of several products
  * A product belongs to a product family
  * A product owns a set of attributes (for example its location, its operating system type, its type...)
  * A product and all its attributes are identified by a stock keeping unit (SKU)
  * A SKU has a set of offers
* Offers
  * An offer represents a selling contract
  * An offer is specific to a SKU
  * An offer is characterized by the term of the offer
  * A term is typed as either "Reserved" or "OnDemand"
  * A term has attributes
* Prices
  * An offer has at least one price dimension
  * A price dimension is characterized by its currency, its unit of measure, its price per unit, its description and eventually per a range of application (start and end)

Regarding those elements, I have extracted and represented a "t2.micro/linux in virginia" with 3 of its offers and all the prices associated.

Here is the graphical representation generated thanks to [graphviz' fdp](http://www.graphviz.org/)
![Graph Representation](/assets/graphql/graph.svg)


The goal of GraphQL is to extract a subtree of this graph to get part or all information.
As an example, here is a tree representation of the same graph:

![Graph Representation](/assets/graphql/graph_tree.svg)

_Note_: I wrote a very quick'n'dirty parser to get the information which can be found [here](https://gist.github.com/owulveryck/bac700e2f5e5b1af0fffda4e7adb9eed). I wrote an idiomatic one but it is the property of the company I made it for.

# Defining the GraphQL schema

The first thing that needs to be done is to write the [schema](http://graphql.org/learn/schema/) that will define the _query_ type.

I will not go into deep details in here. I will simple refer to this excellent document which is a _résumé_ of the language:
[Graphql shorthand notation cheat sheet](https://github.com/sogko/graphql-schema-language-cheat-sheet/raw/master/graphql-shorthand-notation-cheat-sheet.png)

We can define a product that must contains a list of offers this way and a product family like this:

{{< highlight graphql >}}
# Product definition
type Product {
  offers: [Offer]!
  location: String
  instanceType: String
  sku: String!
  operatingSystem: String
}

# Definition of the product family
type ProductFamily {
  products: [Product]!
}
{{</ highlight >}}

One offer is composed of a mandatory price list. An offer must be of a pre-defined type: _OnDemand_ or _Reserved_.
Let's define this:
{{< highlight graphql >}}
# Definition of an offer
type Offer {
  type: OFFER_TYPE!
  code: String!
  LeaseContractLength: String
  PurchaseOption: String
  OfferingClass: String
  prices: [Price]!
}

# All possible offer types
enum OFFER_TYPE {
  OnDemand
  Reserved
}

# Definition of a price
type Price {
  description: String
  unit: String
  currency: String
  price: Float
}
{{</ highlight >}}

At the very end we define the _queries_ 
Let's start by defining a single query. To make it simple for the purpose of the post, Let's assume that we will try to get a whole _product family_.
If we query the entire product family, we will be able to display all informations of all product in the family. But let's also consider that we want to limit the family and extract only a certain product identified by its SKU.

The Query definition is therefore:
{{< highlight graphql >}}
# root Query type
type Query {
    products(sku: String): [Product]
}
{{</ highlight >}}

We will query products (`{products}`) and it will return a ProductFamily.

## Query

Let's see now how a typical query would look like. To understand the structure of a query, I advise you to read this excellent blog post: [The Anatomy of a GraphQL Query](https://dev-blog.apollodata.com/the-anatomy-of-a-graphql-query-6dffa9e9e747#.jbklz6h17).

{{< highlight graphql >}}
{
  ProductFamily {
    products {
      location
      type
    }
  }
}
{{</ highlight >}}

This query should normally return all the products of the family and display their location and their type.
Let's try to implement this

# Geek time: let's go!

I will use the `go` implementation of GraphQL which is a "simple" translation in go of the [javascript's reference implementation](https://github.com/graphql/graphql-js).

To use it: 

{{< highlight go >}}
import "github.com/graphql-go/graphql"
{{</ highlight >}}

To keep it simple, I will load all the products and offers in memory. In the real life, we should implement an access to whatever database. But that is a strength of the GraphQL model: The flexibility. The backend can be changed later without breaking the model or the API.

## First pass: Only the products

### Defining the schema and the query in go

Most of the work has already been done and documented in a series of blog posts [here](http://mycodesmells.com/post/building-graphql-api-in-go)

First we must define a couple of things:

* A _Schema_ as returned by the function `graphql.NewSchema` that takes as argument a `graphql.SchemaConfig`
* The `graphql.SchemaConfig` is a structure composed of a `Query`, a `Mutation` and other alike fields which are pointers to `graphql.Object`
* The rootQuery is created by the structure `graphql.ObjectConfig` in which we pass an object of type `graphql.Fields` (which is a `map[string]*Field`)

The code to create the schema is the following:
{{< highlight go >}}
fields := graphql.Fields{}
rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
schemaConfig := graphql.SchemaConfig{
       Query: graphql.NewObject(rootQuery),
}
schema, err := graphql.NewSchema(schemaConfig)
{{</ highlight >}}

### Defining the fields

Our shema is created but nearly empty because we did not filled the "fields" variable.
the fields variable will contain what the user can request.

As seen before, fields is a map of `*Field`. The key of the map is the root query. In our definition of the Query, we declared that the query would be "products". So "products" is the key of the map.
The graphql.Field that is returned is a list type composed of productTypes.

{{< highlight go >}}
fields := graphql.Fields{
        "products": &graphql.Field{
             Type: graphql.NewList(productType),
        ...
{{</ highlight >}}

We will see in a minute how to define the _productType_. Before, we must provide a way to seek for the product in the database.
This is done by implementing the `Resolve` function:

{{< highlight go >}}
fields := graphql.Fields{
        "products": &graphql.Field{
             Type: graphql.NewList(productType),
             Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        ...
{{</ highlight >}}

The resolv function will return all the products in our database.

But wait... In the Query definition, we said that we wanted to be able to limit the product by setting a sku in the query.

To inform our schema that it can handle a we add the `Args` field to the `graphql.Field` structure:


{{< highlight go >}}
fields := graphql.Fields{
        "products": &graphql.Field{
              Type: graphql.NewList(productType),
              Args: graphql.FieldConfigArgument{
                      "sku": &graphql.ArgumentConfig{
                                Type: graphql.String,
                      },
              },
              Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        ...
{{</ highlight >}}

as the argument is not mandatory, we will use an if statement in the Resolve function to check whether we have a sku or not:

{{< highlight go >}}
if sku, skuok := p.Args["sku"].(string); skuok {
{{</ highlight >}}

### Defining the _productType_

To be able to query display the information of the product (and query the fields), we must define the productType as a graphql object.
This is done like this:

{{< highlight go >}}
var productType = graphql.NewObject(graphql.ObjectConfig{
        Name: "Product",
        Fields: graphql.Fields{
                "location": &graphql.Field{
                        Type: graphql.String,
                },
                "sku": &graphql.Field{
                        Type: graphql.String,
                },
                "operatingSystem": &graphql.Field{
                        Type: graphql.String,
                },
                "instanceType": &graphql.Field{
                        Type: graphql.String,
                },
        },
})
{{</ highlight >}}

A productType is a graphql object composed of the 4 fields. Those fields will be returned as string in the graphql.

### Querying

I will not implement a webservice to query my schema by now. This can easily be done with some handlers that are part of the project.
I will use the same technique as found on internet: I will put the query as argument to my cli.

Assuming that `query` actually holds my my graphql request, I can query my schema by doing:

{{< highlight go >}}
params := graphql.Params{Schema: schema, RequestString: query}
r := graphql.Do(params)
if r.HasErrors() {
    log.Fatalf("Failed due to errors: %v\n", r.Errors)
}
{{</ highlight >}}

### A couple of tests...
    ./pricing -db bla -query "{products(sku:\"HZC9FAP4F9Y8JW67\"){location}}" | jq "."
{{< highlight json >}}
{
  "data": {
    "products": [
      {
        "location": "US East (N. Virginia)"
      }
    ]
  }
}
{{</ highlight >}}
     
    ./pricing -db bla -query "{products(sku:\"HZC9FAP4F9Y8JW67\"){location,instanceType}}" | jq "."
{{< highlight json >}}
{
  "data": {
    "products": [
      {
        "location": "US East (N. Virginia)",
        "instanceType": "t2.micro"
      }
    ]
  }
}
{{</ highlight >}}

    ./pricing -db bla -query "{products{location}}" | jq "." | head -15
{{< highlight json >}}
{
  "data": {
    "products": [
      {
        "location": "US East (Ohio)"
      },
      {
        "location": "EU (Frankfurt)"
      },
      {
        "location": "EU (Frankfurt)"
      },
      {
        "location": "Asia Pacific (Sydney)"
      },

{{</ highlight >}}

    ./pricing -db bla -query "{products{location,operatingSystem}}" | jq "." | head -20
{{< highlight json >}}
{
  "data": {
    "products": [
      {
        "operatingSystem": "Windows",
        "location": "Asia Pacific (Sydney)"
      },
      {
        "operatingSystem": "Windows",
        "location": "AWS GovCloud (US)"
      },
      {
        "operatingSystem": "Windows",
        "location": "Asia Pacific (Mumbai)"
      },
      {
        "operatingSystem": "SUSE",
        "location": "US East (N. Virginia)"
      },
{{</ highlight >}}

## Adding the Offers

To add the offer, we should first define a new offerType

{{< highlight go >}}
var offerType = graphql.NewObject(graphql.ObjectConfig{
        Name: "Offer",
        Fields: graphql.Fields{
                "type": &graphql.Field{
                        Type: graphql.String,
                },
                "code": &graphql.Field{
                        Type: graphql.String,
                },
                "LeaseContractLenght": &graphql.Field{
                        Type: graphql.String,
                },
                "PurchaseOption": &graphql.Field{
                        Type: graphql.String,
                },
                "OfferingClass": &graphql.Field{
                        Type: graphql.String,
                },
        },
})
{{</ highlight >}}

And then make the productType aware of this new type:

{{< highlight go >}}
var productType = graphql.NewObject(graphql.ObjectConfig{
        Name: "Product",
        Fields: graphql.Fields{
                "location": &graphql.Field{
                        Type: graphql.String,
                },
                "sku": &graphql.Field{
                        Type: graphql.String,
                },
                "operatingSystem": &graphql.Field{
                        Type: graphql.String,
                },
                "instanceType": &graphql.Field{
                        Type: graphql.String,
                },
                "offers": &graphql.Field{
                        Type: graphql.NewList(offerType),
                },
        },
})
{{</ highlight >}}

Then, make sure that the resolv function is able to fill the structure of the product with the correct offer.

### Testing:

    ./pricing -db bla -query "{products(sku:\"HZC9FAP4F9Y8JW67\"){location,instanceType,offers{type,code}}}" | jq "."
{{< highlight json >}}
{
  "data": {
    "products": [
      {
        "offers": [
          {
            "type": "OnDemand",
            "code": "JRTCKXETXF"
          }
        ],
        "location": "US East (N. Virginia)",
        "instanceType": "t2.micro"
      }
    ]
  }
}
{{</ highlight >}}

This is it!

# Conclusion

I didn't document the prices, but it can be done following the same principles.

Graphql seems really powerful. Now that I have this little utility, I may try (once more) to develop a little react frontend or a _GraphiQL_ UI.
What I like most is that it has forced me to think in graph instead of the traditional relational model.

The piece of code is on [github](https://github.com/owulveryck/graphql-test)

**edit**: I have included a graphiql interpreter for testing. It works great. Everything is on github:

![GraphiQL](/assets/images/graphiql.png)
