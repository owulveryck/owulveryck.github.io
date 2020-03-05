---
categories:
- aws
date: 2017-01-13T22:22:46+01:00
description: "In this post, I will explain how to extract, process and store informations from a webservice to a NoSQL database (DynamoDB)"
draft: false
images:
- /assets/images/bigdata/stones-483138_640.png
tags:
- dynamodb
- aws
- golang
title: A foot in NoSQL and a toe in big data
---

The more I work with AWS, the more I understand their models. This goes far beyond the technical principles of micro service.
As an example I recently had an opportunity to dig a bit into the billing process.
I had an explanation given by a colleague whose understanding was more advanced than mine.
In his explanation, he mentioned this blog post: [New price list API](https://aws.amazon.com/blogs/aws/new-aws-price-list-api/).

# Understanding the model
By reading this post and this [explanation](http://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/price-changes.html), I understand that the offers are categorized in families (eg AmazonS3) and that an offer is composed of a set of products.
Each product is characterized by its SKU's reference ([stock-keeping unit](https://en.wikipedia.org/wiki/Stock_keeping_unit))

## Inventory management

So finally, it is just about inventory management. In the retail, when you say "inventory management", the IT usually replies with millions dollars _ERP_.
And the more items we have, the more processing power we need and then more dollar are involved... and richer the IT specialists are (just kidding).

Moreover enhancing an item by adding some attributes can be painful and risky

![xkcd](http://imgs.xkcd.com/comics/exploits_of_a_mom.png)

## The NoSQL approach 

Due to the rise of the online shopping, inventory management must be real time.
The stock inventory is a business service. and placing it in a micro service architecture bring constraints: the request should be satisfied in micro seconds.

More over, the key/value concept allows to store "anything" in a value. Therefore, you can store a list of attributes regardless of what the attributes are.

When it comes to NoSQL, there are usually two approaches to store the data:

* simple Key/Value;
* document-oriented.

At first I did and experiment with a simple key/value store called BoltDB (which is more or less like Redis).
In this approach the value stored was a json representation... A kind of document.
Then I though that it could be a good idea to use a more document oriented service: DynamoDB

# Geek time

In this part I will explain how to get the data from AWS and to store them in the dynamoDB service. The code is written in GO and is just a proof of concept.

## The product informations

A product's technical representation is described [here](http://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/reading-an-offer.html).
We have:

{{< highlight js >}}
"Product Details": {
   "sku": {
      "sku":"The SKU of the product",
      "productFamily":"The product family of the product",
      "attributes": {
         "attributeName":"attributeValue",
      }
   }
}
{{</ highlight >}}

There are three important entries but only two are mandatories:

* *SKU*: A unique code for a product. 
* *Product Family*: The category for the type of product. For example, compute for Amazon EC2 or storage for Amazon S3.
* Attributes: A list of all of the product attributes.

## Creating the "table"

As my goal is for now to create a proof of concept and play with the data, I am creating the table manually.
DynamoDB allows the creation of two indexes per table. So I create a table _Products_ with two indexes:

* **SKU**
* **ProductFamily**

![Create Table](/assets/images/bigdata/blog-dynamo-create-table.png)

## Principle

The data is retrieved by a simple `http.Get` method. Then a `json.Decoder` takes the body (an `io.Reader`) as argument and decode it in a predefined structure.
Once the structure is filled, I will store it in the DynamoDB.


### The structures

I need three go structures. Two will be used to decode and range through the offer index. The other one will hold all the product details for a specific offer.

#### Offer Index
The offer index is composed of offers referenced in by an offer name (`map[string]offer`)
{{< highlight go >}}
type offerIndex struct {
    FormatVersion   string           `json:"formatVersion"`
    Disclaimer      string           `json:"disclaimer"`
    PublicationDate time.Time        `json:"publicationDate"`
    Offers          map[string]offer `json:"offers"`
}
{{</ highlight >}}

An offer in the index is characterized by three elements. I am catching all of them, but only `CurrrentVersionURL` is useful in my case.
{{< highlight go >}}
type offer struct {
    OfferCode         string `json:"offerCode:"`
    VersionIndexURL   string `json:"versionIndexUrl"`
    CurrentVersionURL string `json:"currentVersionUrl"`
}
{{</ highlight >}}

#### Products
I hold all the product details in a structure. The product details holds all the products in a map whose key is the SKU. Therefore a SKU field is useless.
The Attribute value is an interface{} because it can be of any type (more on this later in the post).

_Note_ : In case of massive data flow, it would probably be better to decode the stream pieces by pieces (as written in the [the go documentation](https://golang.org/pkg/encoding/json/#Decoder.Decode))

{{< highlight go >}}
type productDetails struct {
    Products map[string]struct { // the key is SKU
        ProductFamily string                 `json:"productFamily"`
        Attributes    map[string]interface{} `json:"attributes"`
    } `json:"products"`
}
{{</ highlight >}}

### Getting the data

#### Offers 
The first action is to grab the json of the offer index and put it in a object of type `offerIndex`
{{< highlight go >}}
resp, err := http.Get("https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/index.json")

var oi offerIndex
err = json.NewDecoder(resp.Body).Decode(&oi)
// oi contains all the offers
{{</ highlight >}}

Then loop for each offer and do a `GET` of every `CurrentVersionURL`
{{< highlight go >}}
for _ , o := range oi.Offers {
        resp, err := http.Get("https://pricing.us-east-1.amazonaws.com" + o.CurrentVersionURL)
{{</ highlight >}}

#### And products

The same principles applies for the products, we decode the stream in an object:

{{< highlight go >}}
var pd productDetails
err = json.NewDecoder(resp.Body).Decode(&pd)
{{</ highlight >}}

Now that we have all the informations we are ready to store them in the database.

## Storing the informations

As usual with any AWS access, you need to create a `session` and a `service` object:

{{< highlight go >}}
sess, err := session.NewSession()
svc := dynamodb.New(sess)
{{</ highlight >}}

The [session](http://docs.aws.amazon.com/sdk-for-go/api/aws/session/) will take care of the credentials by reading the appropriate files or environment variables.

the `svc` object is used to interact with the DynamoDB service. To store an object we will use the method [PutItem](http://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.PutItem) which takes as argument a reference to [PutItemInput](http://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#PutItemInput).

_Note_ All of the AWS service have the same logic and work the same way: Action takes as a parameter a reference to a type ActionInput and returns a type ActionOutput.

Let's see how to create a `PutItemInput` element from a `Product` type.

#### the Dynamodb Item

The two mandatory fields I will use for the `PutItemInput` are:

* `TableName` (which is Product in my case)
* `Item` (which obviously hold what to store)

Other fields exists, but to be honest, I don't know whether I need them by now.

The `Item` expects a map whose key is the field name (In our case it can be "SKU", "ProductFamily" or anything) and whose value is a reference to the special type [AttributeValue](http://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#AttributeValue).

From the documentation the definition is:

_AttributeValue Represents the data for an attribute. You can set one, and only one, of the elements._

The AttributeValue is _typed_ (The types are described [here](https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeValue.html))
Therefore our informations (remember the `map[string]inteface{}`) must be "convrted" to a dynamodb format.
This task has been made easy by using the package [dynamodbattribute](https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeValue.html) which does it for us:

To fill the item I need to loop for every product in the object `pd` and create an item:

{{< highlight go >}}
for k, v := range pd.Products {
      item["SKU"], err = dynamodbattribute.Marshal(k)
      item["ProductFamily"], err = dynamodbattribute.Marshal(v.ProductFamily)
      item["Attributes"], err = dynamodbattribute.Marshal(v.Attributes)
{{</ highlight >}}

Once I have an Item, I can create the parameters and send the request to the DB:

{{< highlight go >}}
Item:      item,
      TableName: aws.String(config.TableName),
}
// Now put the item, discarding the result
_ , err = svc.PutItem(params)
{{</ highlight >}}

# Execution and conclusion

Once compiled I can run the program that will take a couple of minute to execute (it can easily be optimized simply by processing each offer in a separate goroutine).
Then I can find the informations in my DB:
![Result](/assets/images/bigdata/blog-dynamo-result.png)

Now that I have the informations, on the same principle I can grab the prices and put a little web service in front of it.
And I could even code a little fronted for the service.

I am aware that if you are not an average go programmer the code may seem tricky, but I can assure you that it is not (the whole example is less than 100 lines long including the comments).
The AWS API seems strange and not idiomatic, but it has the huge advantage to be efficient and coherent.

Regarding the inventory model. it can be used for any product or even any stock and prices. It is a cheap (and yet efficient) way to manage an inventory.

# Full code

The full code of the example can be found on my [gist](https://gist.github.com/owulveryck/f9665470e8334e8609434feeeddc6071)
