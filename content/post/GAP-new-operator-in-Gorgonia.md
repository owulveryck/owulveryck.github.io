---
title: "Implementing a new operator in Gorgonia - The Global Average Pooling example"
date: 2019-07-01T14:32:59+02:00
lastmod: 2019-07-01T14:32:59+02:00
draft: true
keywords: []
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

<!--more-->

# Introduction

## About Global Average Pooling (GAP)

# Implementation


## The exported function

```go 
// GlobalAveragePool2D consumes an input tensor X and 
// applies average pooling across the values in the same channel.
// The expected input shape is BCHW where
// B is the batch size, 
// C is the number of channels, 
// and H and W are the height and the width of the data.
func GlobalAveragePool2D(x *Node) (*Node, error) {
	return nil, nyi("operator", "Global Average Pooling")
}
```

## Expected result, a `go test` example

```go 
func TestGlobalAveragePool2D_fwdPass(t *testing.T) {
	inputT := tensor.New(
		tensor.WithShape(1, 3, 5, 5),
		tensor.WithBacking([]float64{
			1.7640524, 0.4001572, 0.978738, 2.2408931, 1.867558,
			-0.9772779, 0.95008844, -0.1513572, -0.10321885, 0.41059852,
			0.14404356, 1.4542735, 0.7610377, 0.121675014, 0.44386324,
			0.33367434, 1.4940791, -0.20515826, 0.3130677, -0.85409576,
			-2.5529897, 0.6536186, 0.8644362, -0.742165, 2.2697546,

			-1.4543657, 0.045758516, -0.18718386, 1.5327792, 1.4693588,
			0.15494743, 0.37816253, -0.88778573, -1.9807965, -0.34791216,
			0.15634897, 1.2302907, 1.2023798, -0.3873268, -0.30230275,
			-1.048553, -1.420018, -1.7062702, 1.9507754, -0.5096522,
			-0.4380743, -1.2527953, 0.7774904, -1.6138978, -0.21274029,

			-0.89546657, 0.3869025, -0.51080513, -1.1806322, -0.028182229,
			0.42833188, 0.06651722, 0.3024719, -0.6343221, -0.36274117,
			-0.67246044, -0.35955316, -0.8131463, -1.7262826, 0.17742614,
			-0.40178093, -1.6301984, 0.46278226, -0.9072984, 0.051945396,
			0.7290906, 0.12898292, 1.1394007, -1.2348258, 0.40234163}))
	expectedOutput := tensor.New(
		tensor.WithShape(1, 3, 1, 1),
		tensor.WithBacking([]float64{0.47517386, -0.1940553, -0.28326008}))
	g := NewGraph()
	assert := assert.New(t)
	x := NodeFromAny(g, inputT)
	output, err := GlobalAveragePool2D(x)

	if err != nil {
		t.Fatal(err)
	}
	m := NewTapeMachine(g)
	if err := m.RunAll(); err != nil {
		t.Fatalf("%+v", err)
	}
	defer m.Close()
	if len(output.Shape()) != len(expectedOutput.Shape()) {
		t.Fatalf("Bad output shape, expected %v, got %v", expectedOutput.Shape(), output.Shape())
	}
	for i, d := range output.Shape() {
		if expectedOutput.Shape()[i] != d {
			t.Fatalf("Bad output shape, expected %v, got %v", expectedOutput.Shape(), output.Shape())
		}
	}
	assert.InDeltaSlice(expectedOutput.Data(), output.Value().Data(), 1e-6, "the two tensors should be equal.")
}
```


```
✗  go test -run=GlobalAveragePool2D -v
=== RUN   TestGlobalAveragePool2D_fwdPass
--- FAIL: TestGlobalAveragePool2D_fwdPass (0.00s)
    nn_test.go:512: operator not yet implemented for Global Average Pooling
FAIL
exit status 1
FAIL    gorgonia.org/gorgonia   0.024s
```



## An Operator in Gorgonia

### Fuflfiling the Op interface



```go
type globalAveragePoolOp struct{}

func (g *globalAveragePoolOp) Arity() int {
	panic("not implemented")
}

func (g *globalAveragePoolOp) Type() hm.Type {
	panic("not implemented")
}

func (g *globalAveragePoolOp) InferShape(...DimSizer) (tensor.Shape, error) {
	panic("not implemented")
}

func (g *globalAveragePoolOp) Do(...Value) (Value, error) {
	panic("not implemented")
}

func (g *globalAveragePoolOp) ReturnsPtr() bool {
	panic("not implemented")
}

func (g *globalAveragePoolOp) CallsExtern() bool {
	panic("not implemented")
}

func (g *globalAveragePoolOp) OverwritesInput() int {
	panic("not implemented")
}

func (g *globalAveragePoolOp) WriteHash(h hash.Hash) {
	panic("not implemented")
}

func (g *globalAveragePoolOp) Hashcode() uint32 {
	panic("not implemented")
}

func (g *globalAveragePoolOp) String() string {
	panic("not implemented")
}
```

### The `Do` method

## Applying the Operator onto the node of the Graph (`ApplyOP`)

```go 
func ApplyOp(op Op, children ...*Node) (retVal *Node, err error)
```

```go 
// GlobalAveragePool2D consumes an input tensor X and 
// applies average pooling across the values in the same channel.
// The expected input shape is BCHW where
// B is the batch size, 
// C is the number of channels, 
// and H and W are the height and the width of the data.
func GlobalAveragePool2D(x *Node) (*Node, error) {
	return ApplyOp(&globalAveragePoolOp,x)
}
```


# Testing the Implementation

# Conclusion
