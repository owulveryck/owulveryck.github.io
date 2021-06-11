---
title: "20210608 Datamesh Cue"
date: 2021-06-08T15:47:20+02:00
draft: true
---

In 2021, a rich set of data is the soil that empowers the business of all the Internet Giants (GAFAM, NATU, ...).

Meanwhile, traditional companies are striving to remain competitive. The mandatory acceleration of their business goes through a massive digitalization of their operations and assets.

Amongst the most valuable assets of the companies are the data. Big data's promises are attractive. However, in the wild, the data unit is commonly separated from the core business. Even if much those data units provide much effort, the business plan usually looks like this:

- step 1: collect
- step2: ?
- step3: profit

{{< figure src="/assets/datamesh/gnome_data.png" link="/assets/datamesh/gnome_data.png" title="The Gnome's business plan" width="300px" >}}

In this article, I will present a way to address **step 2** of the plan. I will take some concepts from the data-mesh paradigm, and, to avoid the [cargo cult](https://en.wikipedia.org/wiki/Cargo_cult), I will eventually illustrate them with a trivial technological implementation.

_Disclaimer_ the implementation described here is a skeleton that must be considered as a proof-of-value. As usual, adapt is better than adopt, and most of the
concepts described here should be adapted depending on the business context and building constraints.

## Previously in the world of _Data_

As explained in the introduction, much effort has been in technological solutions to address big data issues and extract its value.

If we apply the adage "_during a gold rush, sell shovels_" to our context, the _shovels retailers_ lead to various technological implementations such as data-warehouse, data-lake, and lately data-factories. Even if it may sound ok, those solutions share a common problem: they can hardly scale.

To manage this rush by addressing the _shovels_ ecosystem issues while focusing on the _gold_ (the data),  Zhamak Deghani introduced a paradigm shift called _data mesh_.
The data-mesh is a way to exploit the data in a distributed manner.
In essence, the paradigm shit is:

- focusing on the distribution of ownership and technological architecture;
- placing the data at the center of each distributed component.

**All the rest of the data mesh is about solving the problems that come with that.**

### The pillars of the data mesh in a glimpse

The datamesh is based on four pillars:

- Federated computational governance.
- A domain-driven data ownership architecture.
- Thinking of data as a product.
- Relying on a self-serve infrastructure platform.

Let's now go through an extract of concepts that we want to illustrate with our proof-of-value.

#### Data as product

The first pillar we need to define is to treat the data as a product.
To be profitable to the business, the data must be considered as a product.

- Discoverable: Declared on a catalog and a search engine
- Understandable: provide semantic (meaning), syntactic (topology) and usage description (behavior)
- Addressable: must participate in a global ecosystem with a unique address that helps its users to programmatically find and access it
- Secure: Be able to be accessed security with global policies (role-base-access, GDPR, Info security, data sovereignty …)
- Interoperable: Be able to reuse, correlate and stitch them together across namespaces for new use-cases
- Trustworthy – Truthful: Provided data provenance and lineage and data quality from the owner
- Natively Accessible: Provided multimodal access like Web services, event of file interfaces
- Valuable on its own: Designed to higher insights when combined and correlated
- Committed to SLOs: Must respect expected service levels, in terms of data availability and quality

#### Federated computed governance

TODO: speak about why the "computed" is important in the sentence

#### Self serve data platform

TODO: a lever to increase the TTM and to ensure the robustness of the mesh.

### Data mesh representation

To oversimplify the idea for the rest of the article, I will represent a mesh like this:

A set of autonomous products that provides values on their own (Profit) by collecting data, and that are exposing their data to others to provide a bigger value to the company:
{{< figure src="/assets/datamesh/set_data_products.png" link="/assets/datamesh/set_data_products.png" width="400px">}}

The fuel of each product is a set of data that are provided by grabbing data from the operational services, and by other products via a set of communication channels. The sum is the mesh.
{{< figure src="/assets/datamesh/data_mesh.png" link="/assets/datamesh/data_mesh.png" width="430px">}}

## Mesh: a matter of communication

In the model we've exposed, communiation is mandatory for the mesh to exist.
Without the communications, we end up with an independant set of nodes.

Managing the communication is therefore essential to build products that are _understandable_, _interoperable_ and _accessible_. On top of that, a good
communication network shall allow _discoverability_ of the products

Let's now see how to implement a robust communication system that fulfils the pre-requisistes of the data-mesh.

### Modeling the communication

In this section, we will introduce basic concepts that will help in the understanding of the technology implementation that will follow.

The basic model representing a communication system has been represented in 1948 by Claude Shannon.
Let's borrow this  representation and the explanation from the essay [_A mathematical theory of communication_](http://people.math.harvard.edu/~ctm/home/text/others/shannon/entropy/entropy.pdf):

{{< figure src="/assets/datamesh/Picture-of-the-schematic-diagram-of-a-general-communication-system-Claude-Shannon-on.png" link="/assets/datamesh/Picture-of-the-schematic-diagram-of-a-general-communication-system-Claude-Shannon-on.png" width="500px">}}

Let's put aside the noise source and focus on the other elements:

 - an _information source_ that produces a message or sequence of messages to be communicated to the receiver terminal. In our case, the message is data to be communicated to other nodes of the mesh.
- a _transmitter_ which operates on the message in some way to produce a _signal_ suitable for transmission over the channel.
- the _channel_ is merely the medium used to transmit the signal from transmitter to receiver.
- the _receiver_ ordinarily performs the inverse operation of that done by the transmitter, reconstructing the message from the _signal_.
- the destination is the thing for whom the message is intended.

Roughly, standardizing the communication network in the mesh will lead to this representation:
{{< figure src="/assets/datamesh/data_streaming.png" width="300px" >}}

## Application to our mesh

### From a communication model to a processing data pipeline

Instantiating this communication model to the data world is roughly equivalent to describe a shallow data processing pipeline.

### Making the message understandable: the semantic

The source and the destination must agree on the semantic of the message. In computer science, this is done by sharing a schema and definitions of the information.
For example, in English we can express a message like this:

_The message contains the indentity of a person. The identity of a person is composed of his first name starting with a capital letter, his last name starting with a capital lette and optionally his age which is a number less than 130_ 

On top of that, it is the role of the "**federated governance**" is to impose a common language to express the messages as well as the Schema definitions.
To be _programatically addressable_, the definition and data must be expressed in a computer friendly language such as JSON-Schema, Protobuf or CUE for example.
This is why the datamesh calls it a "**federated computational governance**".

For example here is the definition of the schema in [CUE language](cuelang.org):

```cue
message: #Identity & {
    first: "John"
    Last: "Doe"
    Age: 40
}

#Identity: {
        // first name of the person
        first: =~ "[A-Z].*"
        // Last name of the person
        Last: =~ "[A-Z].*"
        // Age of the person
        Age?: number & < 130
}
```

_Note_: more on the choice of the CUE language will come later in this article. Meanwhile, let's check that the computer understands the syntax and that is it correct by calling `cue vet identity.cue`. Changing the Last name from `Doe` do `doe` or setting an age above 130 would result in an error.

The validation prevents sending nois or garbage over the channel:

{{< figure src="/assets/datamesh/garbage_in_out.png" width="300px">}}

_As a résumé_ the role of the message emitter is to expose its semantic in a language defined by the "federated governance" and to emit a message that is syntacticly and functionally coherent with the definition (as validated by the "federated **computational** governance").

### Creating a signal: the transmitter

Once we have transformed the information into a understandable message (this operation is usually called a marshaling process), we pass it to the transmitter that will encode it into a signal and emit the data.

The role of the signal is to ensure that the information is safely transmitted over the communiation channel. 
On top of that, in our mesh, the signal format should allow multiplexing to avoid a spaghetti of channels.

Encapsulating the message into an envelop is a way to address the problem. The envelop allows creating a common structure to handle meta data such as the emitter of the message, its type, its source and so on.

Once again, it is the role of the "**federated computational governane**" to define the envelop format and standards. [Cloudevents](cloudevents.io) is one of those for example. It standardize the exchange of messages between the cloud services.

_As a résumé_, the role of the transmitter is to create a port the message over the channel by encapsulating it into an event (aka marshaling the event). The event
envelop is standardized by the governance. The transmitter is a capacity offered by an "self-serve infrastruture" (the data products should be autonomous to transmit some event)

### Broadcasting the signal: the channel

The role of the channel is to store and expose the events to the receivers.
It is the role of the channel to validate that the message, once accepted is delivered to the authorized and intended receivers. This will guarantee the security and trust of the whole infrastructure.

It is not the role of the channel to analyse the message in any way. It is, therefore, independant of the type of messages (think about the telephone, you can speak English or French in a telephone).

## Trivial implementation: a streaming platform

Now we have all the concepts let's implement the self-serve communiation infrastructure that will facilitate the developement of the product while ensuring the rules of the **federated computational governance**.

First, let's summarize the pipeline using the unix pipe (`|`) symbol:

```shell
// Send:
collect_data | marshal_message | emit_message | validate_message | marshal_event | send_to_channel

// Receive:
filter_events_from_channel | read_from_channel | unmarshal_event | unmarshal_data | profit
```

To facilitate the developpement and maintenance of the mesh, the self-serve communiation infrastructure (let's call it a streaming plateform) will provide those capabilities:

- validate_message: as said before, it is a must to ensure the data quality
- channel management
- event filtering, and more precisely event routing

On top of that, it must provide a repository of data-schemas (data catalog) to make the information `addressable`.

### A configuration language for the semantic

As explained before, the system should be generic enough to be losely coupled with the semantic of the data.

Precisely, it should be the role of the owner of the data-product to express the schema and validation rules;
We can therefore consider the validation and data cataloging capabilities as configuration of a generic instance of the streaming plateform.

_Note:_ The [Site Reliability Engineering book](https://sre.google/workbook/configuration-specifics/) defines a configuration as _human-computer interface for modifying system behavior._ 

This is why, in the example, we use a configuration language. We choose CUE because it is easy and concise in the definition. On top of that, it allows:

- data validation by design
- data composition but no inheritance
- holds a set of tooling to format, lint and validate easily the data and the schemas

Our streaming plateform is therefore a generic validation and message publishing interface configured with CUE.

### Execution/Runtime

Once we have configured our streaming plateform to handle and understand any kind of messages described in CUE, we need to go live by ingesting, validating and emitting the data.

The way CUE addresses it is by Unifying the data and Executing the result (Now you understand the name: Configure, Unify Execute).

This is what the commande `cue vet` we've issued before does under the hood. But we may want to turn it into a service. CUE has a SDK in Go that allows doing this:

For example, this simple command does the same thing and can be instanciated into a HTTP handler.

```go
type DataProduct struct {
	definition    cue.Value
	// ...
}

func (d *DataProduct) ExtractData(b []byte) (cue.Value, error) {
	data := d.definition.Context().CompileBytes(b)
	unified := d.definition.Unify(data)
	opts := []cue.Option{
		cue.Attributes(true),
		cue.Definitions(true),
		cue.Hidden(true),
	}
	return data, unified.Validate(opts...)
}
```

## Eventing/Routing

So far, we've seen that it requires very little effort to express a schema and validate the data at the entrance of the streaming platform.

Before submiting it to a communication channel (to be defined), let's ensure that we write a comprehensible envelop.
We've expressed it already, interoperability is key to success of the mesh. Using a standard envelop will guarantee that the message can move out of the ecosystem of the plateform we are building.

[Cloudvents](cloudevents.io) is a standard format of the [Cloud native conputing foundation (CNCF)](https://www.cncf.io/) that address this needs.
The dataplateform will, therefore, encaspulate the data into a cloud event.

It is the role of the plateform to expose and maintain a catalog of sources and types of events.

Example of event:

```json
{
  "specversion": "1.0",
  "id": "1234-4567-8910-1234-5678",
  "source": "MySource",
  "type": "MySource:newPerson",
  "datacontenttype": "application/json",
  "data_base64": "MyMessageInJSONEncodedInBase64=="
}
```

Of course, the plateform can handle the encoding of the event easily.

## Conclusion


