---
date: 2017-07-07T21:06:46+02:00
description: "Imagine a CCTV at home that would trigger an alert when it detects a movement. Ok, this is easy. Imagine a CCTV that would trigger an alert when it detects a human. A little bit trickier. Now imagine a CCTV that would trigger an alert when it sees someone who is not from the family."
draft: false
images:
- /assets/images/tensorflowserving-4.png
title: A "Smart" CCTV with Tensorflow, and Inception? On a rapsberry pi?
---

Imagine a CCTV at home that would trigger an alert when it detects a movement. 

Ok, this is easy. 

Imagine a CCTV that would trigger an alert when it detects a human (and not the cat). 

A little bit trickier. 

Now imagine a CCTV that would trigger an alert when it sees someone who is not from the family...

__Disclaimer__: This article will not cover everything. I may post a second article later (or not). As you may now, I am doing those experiments during the night as all of this is not part of my job. I hope that I will find time to actually conclude the experiment. If you are a geek and you want to test that with me, feel free and welcome to contact me via the comments or via twitter [@owulveryck](https://twitter.com/owulveryck).
In this article, I will describe the method. I will also retrain a neural network to detect people. I will also use a GO static binary to run it live and evaluate the performances. By the end, I will try a static cross compilation to run it on a raspberry pi, but as my rpi is by now out-of-order, I will test it on qemu.

# Machine learning and computer vision

Machine learning and tooling around it has increased and gained in efficiency in the past years. it now "easy" to code a model that can be trained to detect and classify elements from a picture. 

Cloud providers are offering services that can instantly tag and label elements from an image. To achieve the goal of the CCTV, it would be really easy to use, for example, [AWS rekognition](https://aws.amazon.com/rekognition/), train the model, and post a request for each image seen.

This solution presents a couple of problems:

* The network bandwidth: you need a reliable network and a bandwidth big enough to upload the flow of images

* The cost: these services are cheap for thousand images, but consider about 1 fps to process (I don't even dream of 24fps), it is 86400 images a day and 2.6 million images a month... and considering that 1 million images are 1000 dollar...

I don't even talk about network latency because my CCTV would be pseudo-real-time and the ms of latency can be neglected.

The best solution would be to run the computer vision locally. There are several methods to detect people. The most up-to-date-and-accurate one is based on machine learning and precisely on neural network.

__The very simplified principle of neural network and computer vision:__

There is a lot of literacy on the web around that, but here a very small explanation to understand the rest of this post:

Imagine a picture as a sequence of numbers from 0 to 1 (0 for black 1 for white for example). Imagine a mathematical equation `f`.
You do not know what the content of `f` is. You only tell the computer: "guess a value of `f` such as `f(pictures of a man) = man`.
Then you feed him with a lot of pictures of men, and for every picture, the computer not only guess a function `f` but it adapts it so it can detect every man in every picture.

Sounds magical?

Actually, the computer does not start with a void `f` function. You provide it with a kind of skeleton that you call the neural network.
A neural network is a network of tiny function (neuron) that are applied on the decomposed values of the input (such as a pixel in a photo). 

Depending on the mathematical function coded in the neuron, it is activated by its inputs (there can be several inputs for a single neuron) or not.

You can use several layers of neurons. Each layer is composed of neurons feed by the outputs of the neurons of the previous layer.

The pictures used to feed the model is called the training set. 
You also use a test set (same kind of pictures), that is used to check whether your model generalized well and actually converge to your goal.

I won't dig any further into this description. You can read papers about the [perceptron](https://en.wikipedia.org/wiki/Perceptron) for more accuracy in the description. I heavily recommend this article if you know a little bit of go: [Perceptrons - the most basic form of a neural network](https://appliedgo.net/perceptron/)

<center>
{{< figure src="https://imgs.xkcd.com/comics/machine_learning.png" link="https://xkcd.com/1838/" caption="XKCD 1838" >}}
</center>

## Tools

### Tensorflow

I have already blogged about tensorflow. Tensorflow is not a ML library. It is a mathematical library. It is self-described as  _an open source software library for numerical computation using data flow graphs. Nodes in the graph represent mathematical operations, while the graph edges represent the multidimensional data arrays (tensors) communicated between them._

It is, therefore, an excellent tool, suitable for machine learning and especially for dealing with neural networks. It is brilliant with computer vision as the pictures are arrays of pixels and if you add the colour, you can represent every picture by a tensor.

Even better, the models generated by tensorflow can be saved once learned and transferred to another device. For example, you can train your models on very powerful machines and simply copy the resulted graph to your client (for example a phone). The graph can then be applied to an input taken from the device such as a photo.

### Inception

"[Inception](https://research.google.com/pubs/pub43022.html)" is a deep convolutional neural network architecture used to classify images originally developed by Google.

Inception is exceptionally accurate for computer vision. It can reach 78% accuracy in "Top-1" and 93.9% in "Top-5". That means that if you feed the model with a picture of sunglasses, you have 93.9% chance that the algorithm detects sunglasses amongst the top 5 results.

On top of that, Inception is implemented with Tensorflow, and well documented. Therefore, it easy "easy" to use it, to train it and "to retrain it".

here is a graphical representation of the inception v3 model. You can see the different layers of the model as explained earlier.

<center>
{{< figure src="https://raw.githubusercontent.com/tensorflow/models/master/inception/g3doc/inception_v3_architecture.png" link="https://github.com/tensorflow/models/tree/master/inception" caption="Inception v3 architecture" >}}
</center>

Actually, training the model is a very long process (several days on very efficient machines with GPU). But some folks (at google?) have discovered that retraining only the last layer of the neural network for new classes of pictures was giving good results.

# Geek

I am using the excellent blog post [How to Retrain Inception's Final Layer for New Categories](https://www.tensorflow.org/tutorials/image_retraining).
The purpose of the article is to retrain the network in order to give it the ability to categorize (recognize) different kind of flowers.
I will use exactly the same principle to recognize a class "people".

I will perform the task on a spot instance on AWS (to get it cheap), and download the model to use it locally from a go code.

## Phase 1: recognizing usual people

To keep it simple, I've created a "class" people with the flowers classes. It means that I simply added a directory "people" to my "flowers" for now.

```
[~/flower_photos]$ ls -lrt
total 696
-rw-r----- 1 ubuntu ubuntu 418049 Feb  9  2016 LICENSE.txt
drwx------ 2 ubuntu ubuntu  45056 Feb 10  2016 tulips
drwx------ 2 ubuntu ubuntu  36864 Feb 10  2016 sunflowers
drwx------ 2 ubuntu ubuntu  36864 Feb 10  2016 roses
drwx------ 2 ubuntu ubuntu  57344 Feb 10  2016 dandelion
drwx------ 2 ubuntu ubuntu  36864 Feb 10  2016 daisy
drwxrwxr-x 2 ubuntu ubuntu  77824 Jul  7 14:26 people
```

### Getting a training set full of people

I need a training set of people. That means that I need a certain amount of pictures actually representing some people.
Nowadays it is easy to get training set for free (as in free speech) on the web. 

_Note_ You can see that, by offering services, the GAFA is increasing its training set to make their service more powerful than ever.
<center>
< tweet 857609299731791872 >
</center>

Let's get back to my experiment:
I download pictures of people from [http://www.image-net.org/api/text/imagenet.synset.geturls?wnid=n07942152](http://www.image-net.org/api/text/imagenet.synset.geturls?wnid=n07942152)

{{< highlight shell >}}
curl -s  "http://www.image-net.org/api/text/imagenet.synset.geturls?wnid=n07942152" | \
sed 's/^M//' | \
while read file
do
  curl -m 3 -O $file
done
{{</ highlight >}}

Then I remove all "non-image" files:

{{< highlight shell >}}
for i in $(ls *jpg)
do
    file $i | egrep -qi "jpeg|png" || rm $i 
done
{{</ highlight >}}

## Learning phase

I've had one "issue" during the learning phase. When I executed:

`bazel-bin/tensorflow/examples/image_retraining/retrain --image_dir ~/flower_photos/` 

it failed with a message about `ModuleNotFoundError: No module named 'backports'`. I Googled and found the solution in this [issue](https://github.com/tensorflow/serving/issues/489#issuecomment-313671459). It is because I am using python3 and the tutorial has been written for python 2. No big deal.

At the end of the training (which took 12 minutes on a c4.2xlarge spot instance on AWS) I have two files that hold the previous information.

```
...
2017-07-07 19:22:53.667219: Step 3990: Cross entropy = 0.111931
2017-07-07 19:22:53.728059: Step 3990: Validation accuracy = 93.0% (N=100)
2017-07-07 19:22:54.287266: Step 3999: Train accuracy = 98.0%
2017-07-07 19:22:54.287365: Step 3999: Cross entropy = 0.148188
2017-07-07 19:22:54.348603: Step 3999: Validation accuracy = 91.0% (N=100)
Final test accuracy = 92.7% (N=492)
Converted 2 variables to const ops.
...
```

And a trained graph with a label file that I can export and use elsewhere.

```
(customenv) *[r1.2][~/sources/tensorflow]$ ls -lrth /tmp/output_*
-rw-rw-r-- 1 ubuntu ubuntu  47 Jul  7 19:22 /tmp/output_labels.txt
-rw-rw-r-- 1 ubuntu ubuntu 84M Jul  7 19:22 /tmp/output_graph.pb
```
I have followed the tutorial to [Use the retrained model](https://www.tensorflow.org/tutorials/image_retraining#using_the_retrained_model) to make sure that everything was ok before using it with my own code.

## Using the model with go

Tensorflow is coded in C++, but has some bindings for different languages. The most up-to-date is python, in which a lot of helper libraries are developed (see [tflearn](http://tflearn.org/getting_started/) for example.
A binding for go exists, but it is only implementing the core library of tensorflow. Anyway, it is an excellent choice for applying a model.

The workflow is:

- read the exported model from the disk and create a new graph
- read the label files and set the labels in an array of string
- grab jpeg pictures from the webcam in jpeg (via v4l) in an endless for loop
- Normalize the picture (see below) and create a tensor from the jpeg file.
- Apply the inception model onto the Tensor and getting the `final_result`
- Extract the most important value from the output vector (the better probability) and display the corresponding label.

I will only expose the trickiest parts.

### Getting the pictures

I use a wrapper around `v4l` in go called [go-webcam](https://github.com/blackjack/webcam). As my webcam has MJPEG capabilities, each frame is already in JPEG format.

I am applying the tensorflow model sequentially within the for loop. The problem is that it takes some time to process. And while it is processing the driver may buffer some pictures. Therefore I am totally losing the synchronism. My code may warn me that it has found a person too late.
To avoid this, I am using a non-blocking tick in a go channel within the loop. Therefore I do not process every single frame, but I process a frame every x milliseconds and I discard the rest.
I could have used a pool, but that would have add complexity for the example.

{{< gist owulveryck 1f3fc2366e5a35ab119633d57ad074b6 "tick.go" >}}

### Normalizing the picture

The [example described on the go package](https://godoc.org/github.com/tensorflow/tensorflow/tensorflow/go) is using an old inception implementation (actually version 5h which is older than the v3). Therefore it needs some adjustments. The function that produces a Tensorflow graph that will be used to normalize the picture didn't have the correct normalization values (those defined by the author of the inception v3 model) 

Here is an extract from [Image Recognition](https://www.tensorflow.org/tutorials/image_recognition):

_The model expects to get square 299x299 RGB images, so those are the `input_width` and `input_height` flags. We also need to scale the pixel values from integers that are between 0 and 255 to the floating point values that the graph operates on. We control the scaling with the `input_mean` and `input_std` flags: we first subtract `input_mean` from each pixel value, then divide it by `input_std`._
 
_These values probably look somewhat magical, but they are just defined by the original model author based on what he/she wanted to use as input images for training. If you have a graph that you've trained yourself, you'll just need to adjust the values to match whatever you used during your training process._

{{< gist owulveryck 1f3fc2366e5a35ab119633d57ad074b6 "normalizationgraph.go" >}}

Apart from that the rest of the code remains similar.

#  Conclusion

## Running it on a laptop

The program runs as expected at the rate of 2 images per seconds without overheating on a modern laptop. I have used it to monitor my house while I was on vacation. Every success was sent on an s3 bucket, so in case of the intrusion in my house, I would still have the pictures. I say that it has worked because the only pictures it has recorded were:

* me, leaving the house
* me, entering the house 2 weeks later.

You can find the full code on [my github](https://github.com/owulveryck/smarcctv)

## Further work

### Running on ARM
I want to test it on a raspberry pi, so I have cross compiled the code for ARM with those commands but I didn't have time to test it yet:

{{< highlight bash >}}
# Download a tensorflow release for rpi:
$ wget https://github.com/meinside/libtensorflow.so-raspberrypi/releases/download/v1.2.0/libtensorflow_v1.2.0_20170619.tgz
# Install the toolchain
$ sudo apt install gcc-arm-linux-gnueabihf
# Compile it
$ export CC=arm-linux-gnueabihf-gcc
$ CC=arm-linux-gnueabihf-gcc-5 GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build  -o myprogram -ldflags="-extld=$CC"
{{</ highlight >}}

#### Performances

Inception is very good. But it requires a decent CPU (or even better a GPU). I could use another model called [MobileNet](https://github.com/tensorflow/models/blob/master/slim/nets/mobilenet_v1.md) which is a _low latency, low power_ model. 
It has been [opensourced](https://research.googleblog.com/2017/06/mobilenets-open-source-models-for.html) in June 2017. The tensorflow team has added the ability to retrain it the same way inception is (by retraining the last layer). It's worth a look.

### Detecting only the family
As I explained in the beginning of the post, I want this system to trigger only if it detects someone that is not part of the family.
To do that I need to train the neuron network to classify classes such as: 

* people 
* me 
* my wife 
* kid1
* kid2 

To do so, I need training sets (labeled pictures) of my family. The best way to get it is to write a "memory cortex" to use it with my [cortical](https://github.com/owulveryck/cortical) project as explained in my previous post: [Chrome, the eye of the cloud - Computer vision with deep learning and only 2Gb of RAM](https://blog.owulveryck.info/2017/05/16/chrome-the-eye-of-the-cloud---computer-vision-with-deep-learning-and-only-2gb-of-ram.html).

