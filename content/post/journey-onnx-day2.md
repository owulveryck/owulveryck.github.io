---
title: "My journey with ONNX and Go - Running the graph"
date: 2018-09-19T08:53:09+02:00
lastmod: 2018-09-19T08:53:09+02:00
draft: false
keywords: []
description: "This post is the second part of my experiments with ONNX and Go. In this post I am describing how to create a computation graph in Gorgonia (ExprGraph) from an ONNX Model."
description: ""
tags: []
categories: []
author: ""

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

In the [previous post](/2018/08/14/my-journey-with-onnx-and-go---the-beginning.html#building-the-dag), I made an introduction and a POC to interact with ONNX models and Go.

I have decoded the information to reconstruct a graph.
Now I propose to expand the principle and to create a proper execution backend based on Gorgonia.
This post is a bit more technical than the previous one because all the concepts needed to work should be present in the last article.

# Decoding the tensor

In machine learning, the fundamental element of a computation graph is a Tensor.
In ONNX this element is described in the structure [TensorProto](https://godoc.org/github.com/owulveryck/onnx-go#TensorProto). 
A tensor: 

* has a shape represented here by the field `Dims` (an array of `int64`)
* is holding a data type 
* is eventually handling some data.

Gorgonia also has a notion of tensor. In Gorgonia, a tensor is an interface (in the pure Go sense). Therefore, creating a Go object from TensorProto that fulfills the Tensor interface of Gorgonia should be easy.

Let's write a method that takes an `onnx.TensorProto` as input and that returns a `tensor.Tensor` as output

{{< highlight go >}}
func NewTensor(tx *onnx.TensorProto) (tensor.Tensor, error) { ... } 
{{</ highlight >}}

We need to address the three elements:

* convert the data type to something understandable by Go (and Gorgonia)
* read and process the data to write a tensor backend
* deal with tensor shape.

I will not focus much on tensor shape. Actually, ONNX has a notion of dimension which is an array of integer. Every entry represents the size of an axis of the tensor.
This can be converted out-of-the-box into a [`Shape`](https://godoc.org/gorgonia.org/tensor#Shape) element of the `tensor` package.

The data type conversion and the raw data processing is a (little) bit trickier, so let's focus on them.

### Data types

A tensor is composed of elements of certain types. The supported data types are described as constants in ONNX. They can be found [in the documentation of ONNX](https://github.com/onnx/onnx/blob/master/docs/IR.md#standard-data-types) and are represented in [Go constant values](https://godoc.org/github.com/owulveryck/onnx-go#TensorProto_DataType) of our Go API.

On the other hand, the tensor package of Gorgonia also has its own declaration of types represented by values of [`Dtypes`](https://godoc.org/gorgonia.org/tensor#Dtype). The list is a set of variables declared [here](https://godoc.org/gorgonia.org/tensor#pkg-variables).

Writing a function to return a `Dtype` from a `TensorProto_DataType` is relatively straightforward: 

{{< highlight go >}}
func Dtype(t *onnx.TensorProto_DataType) (tensor.Dtype, error) {
	switch *t {
	case onnx.TensorProto_FLOAT:
		return tensor.Float32, nil
        //...
{{</ highlight >}}

### Raw Data

ONNX has two way to encode the data of a tensor.
The first is really easy and is a straight serialization of the underlying type. For example, a tensor of type Float32 will have its data set in the `FloatData` field which is of type `[]float32`.

The second one is a bit trickier. ONNX allows serializing the "raw data" encoded in a sequence of bytes. The documentation says that:

> When this raw_data field is used to store tensor value, elements MUST
> be stored in as fixed-width, little-endian order.
> Floating-point data types MUST be stored in IEEE 754 format.
> Complex64 elements must be written as two consecutive FLOAT values, real component first.
> Complex128 elements must be written as two consecutive DOUBLE values, real component first.
> Boolean type MUST be written one byte per tensor element (00000001 for true, 00000000 for false).
>
> Note: the advantage of specific field rather than the raw_data field is
> that in some cases (e.g. int data), protobuf does a better packing via
> variable length storage, and may lead to smaller binary footprint.
> When this field is present, the data_type field MUST NOT be STRING or UNDEFINED

So our function must handle this particular case.
Let's focus on the Float32 type for now. Go has natively everything needed to read this famous `IEEE 754 format` (thanks to the binary and the math packages).

Here is how to read the informations and to transcribe it into a `[]float32`:

{{< highlight go >}}
buf := bytes.NewReader(tx.RawData)
element := make([]byte, 4)
var backing []float32
for {
        var n int
        n, err = buf.Read(element)
        if err != nil || n != 4 {
                break
        }
        uintElement := binary.LittleEndian.Uint32(element)
        backing = append(backing, math.Float32frombits(uintElement))
}
{{</ highlight >}}

## Vizualizing the tensor

With all those elements, it is easy to write the content of the `NewTensor` function. No need to paste all the code in this post, but you can find the implementation [here](https://github.com/owulveryck/gorgonnx/blob/9df285e6d96d6ad9494aeeb420fb9f42ebe7f360/vendor/gorgonia.org/tensor/tensonnx/tensor.go#L16).

To do an eye-test of the result, let's convert a 3D-tensor back into an image.

An example can be found in the MNIST example we used in the last post:
```
curl https://www.cntk.ai/OnnxModels/mnist/opset_7/mnist.tar.gz | \
tar -C /tmp -xzvf -
```

The model is delivered with three tests. The tests are made of an input tensor (3D) and the expected output tensor.
Let's take one of the input tensor, convert it to a Gorgonia tensor and create a picture from it (so see if the data, types and shapes are coherents).
I am using the `image` package of the standard Go distribution and dumping a png file on stdout for commodity:

{{< highlight go >}}
b, _ := ioutil.ReadFile("/tmp/mnist/test_data_set_0/input_0.pb")
sampleTestData := new(onnx.TensorProto)
sampleTestData.Unmarshal(b)
t, _ := NewTensor(sampleTestData)
width := t.Shape()[2]
height := t.Shape()[3]
im := image.NewGray(image.Rectangle{Max: image.Point{X: width, Y: height}})
for w := 0; w < width; w++ {
        for h := 0; h < height; h++ {
                v, _ := t.At(0, 0, w, h)
                im.Set(w, h, color.Gray{uint8(v.(float32))})
        }
}
enc := png.Encoder{}
enc.Encode(os.Stdout, im)
{{</ highlight >}}

Running the code produces a `0` as expected:

<center>
{{< figure src="/assets/onnx/0.png" height="200%" title="Representation of a zero from a tensor" >}}
</center>

# Creating an ExprGraph

Now that we are able to decode the tensors from the ONNX model let's go further and create a graph.
In the previous post, we have sliced the parsing function into three parts:

* the processing of the _Initializers_
* the processing of the _Inputs_
* the processing of the _Operators_

(cf [_Building the DAG_](/2018/08/14/my-journey-with-onnx-and-go---the-beginning.html#building-the-dag) in the previous post for more information)

I will take back the skeleton of the code I made to generate the graph with Gonum in the first article.
The main differences are:

* I am now using a pointer to `gorgonia.ExprGraph` in the `computationGraph` structure
* I am using `gorgonia.Node` instead of the `node` structure

The main loops remain similar: 

* creating a node object from the TensorProto (or ValueInfoProto) and adding them to the graph (which is transparent with gorgonia)
* processing the operator nodes

Here is an example with the _Initializers_ of the model (the tensor is generated thanks to the `NewTensor` we wrote before):

{{< highlight go >}}
type computationGraph struct {
	db map[string]*gorgonia.Node
	g  *gorgonia.ExprGraph
}
// ...
var gx *onnx.GraphProto
var cg computationGraph
// ... Various initialization such as parsing the model to fill gx ...
g := gorgonia.NewGraph(gorgonia.WithGraphName(gx.GetName()))
for _, tensorProto := range gx.Initializer {
        name := tensorProto.GetName()
        t, _ := NewTensor(tensorProto)
        n := g.AddNode(gorgonia.NewConstant(t, gorgonia.WithName(name)))
        cg.db[name] = n
}
{{</ highlight >}}

## Processing the operators

The logic is exactly the same as the one we have used in the first article.
The only modification is in the `processNode` method.

This method has a giant switch that delegates the work to other specialized methods.

{{< highlight go >}}
func (cg *computationGraph) processNode(nx *onnx.NodeProto) error {
	switch nType := *nx.OpType; nType {
	case "Add":
		return cg.addOp(nx)
        case "Relu":
                return cg.reluOp(nx)
        case "Conv":
		return cg.convOp(nx)
        //...
{{</ highlight >}}

Then each operation has its own isolated method.

_Note_: There is a better way to handle that, but refactoring will come with a certain maturity of the package. 

The purpose of each method is to analyze the `NodeProto`, extract its attributes and inputs, and create a corresponding node into the Gorgonia Graph.
The operators implemented in ONNX are very well documented in this file accessible from the ONNX repository: [Operators.md](https://github.com/onnx/onnx/blob/master/docs/Operators.md)

Here is, for example, the implementation of the _ReLU_ operator. 
{{< highlight go >}}
import 	nnops "gorgonia.org/gorgonia/ops/nn"

func (d *graph) reluOp(nx *onnx.NodeProto) error {
	n, err := nnops.Rectify(d.db[nx.Input[0]])
	if err != nil {
		return err
	}
	d.db[nx.Output[0]] = n
	return nil
}
{{</ highlight >}}

Most of the work here is to analyze the documentation of the operators from ONNX and to find a way to implement it into Gorgonia. Most of the operators already exist, but some of them may have different parameters.

### Obstacle with the broadcastable operators

A quick word about an obstacle I have faced. It is written in the ONNX documentation that the element-wise operators are broadcastable (the behavior is similar of what numpy implements). The behavior is explained [here](https://github.com/onnx/onnx/blob/master/docs/Broadcasting.md). I made a filthy hack to make my MNIST test pass, but we have an [open issue](https://github.com/gorgonia/gorgonia/issues/223) in Gorgonia to implement a proper way to apply broadcasting in a non-transparent way.

# Computing the MNIST model

My goal is to run the MNIST model and to evaluate it.
Therefore, I have implemented all the Operators used in the model:

* Add
* Conv
* Reshape
* Relu
* MaxPool
* MatMul

With all those operators, I can open the MNIST model, and create the Graph. The tensor shapes are compatible and a dump of the graph gives the following output:

{{< figure src="/assets/onnx/mnist_gorgonia.svg" title="MNIST Model with Gorgonia" >}}
(for clarity I removed the values of the tensors)


## Running the graph

To run the graph, I need to give it one input to work on.
I will use the input of the first section (the one that displayed a `0`).
The model should return a vector of 10 entries with a higher value in the first position:

{{< highlight go >}}
b, _ := ioutil.ReadFile("/tmp/mnist/model.onnx")
model := new(onnx.ModelProto)
model.Unmarshal(b)
g, _ := NewGraph(model.GetGraph())
// Open the tensorproto sample file
b, _ = ioutil.ReadFile("/tmp/mnist/test_data_set_1/input_0.pb")
sampleTestData := new(onnx.TensorProto)
sampleTestData.Unmarshal(b)
t, _ := NewTensor(sampleTestData)
gorgonia.Let(g.ByName("Input3")[0], t)
machine := gorgonia.NewTapeMachine(g)
if err = machine.RunAll(); err != nil {
        log.Fatal(err)
}
fmt.Printf("%v", GetOutputGraphNodes(g)[0].Value().Data())
{{</ highlight >}}

The expected output is also present in the MNIST package. Decoding it gives the following values:

```
[5041.8887 -3568.878 -187.82423 -1685.797 -1183.3232 -614.42926 892.6643 -373.65845 -290.2623 -111.176216]
```

sadly my computation gives the following result:

```
[55.41009 984.514 -1191.4886 -652.1293 802.4857 497.57553 -303.6763 952.77106 -233.73296 -672.92255]
```

By now, I am stuck with this bug, the the goal is reached, I have generated a computation graph that actually runs and gives me a result.
Let's now write a temporary conclusion.

# Conclusion

I am glad to be able to read, understand and compute an ONNX model. Getting the wrong result is annoying but gives a good challenge.
Finding where the problem is not trivial, and debugging a neural network is not easy, but it is a good learning experience to analyze the behavior of the operators in detail.
I have started to implement unit tests for each operator I need in the MNIST model. This is an heavy task, and sometimes I wish I did TDD for this (but this is another story). 

More recently, I have noticed that the ONNX repository was [full of simple test cases made to evaluate the backends](https://github.com/onnx/onnx/tree/master/onnx/backend/test/data/node). This is the next step to implement into the decoding package. 

I am very excited by the possibility to run an ONNX model thanks to a entirely self-sufficient runtime environment. 
With the potential to export this "_VM_" to WASM, we can imagine great applications such as running an ImagerNet network straight to the browser while capting images from the webcam or the microphone. So the challenge now is to fix the MNIST model and to implement more Operators. Then to play and have some more fun with ML! 

If you are interested in testing or contributing, I have set up a repository where you will find the sources and the MNIST example (that you can run with `go test`).
This repository is really a work in progress, and I will not provide (for now) any support around it. 

[https://github.com/owulveryck/gorgonnx](https://github.com/owulveryck/gorgonnx)
