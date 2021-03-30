---
author: Olivier Wulveryck
date: 2016-12-16T14:51:18+01:00
description: Very quick post to present a piece of code that grabs an image from a webcam and send it to amazon's image recognition service.
draft: false
tags:
- webcam
- golang
- aws
title: Image reKognition with a webcam, go and AWS.
topics:
- topic 1
type: post
---

It's been a while since I last posted something. I will fill the gap with a quick post about _rekognition_.

[rekognition](https://aws.amazon.com/rekognition/?nc1=h_ls) is a service from AWS that is described as:

> Deep learning-based image recognition

> Search, verify, and organize millions of images

In this light post, I will present a simple method to grab a picture from my webcam, send it to rekognition and display the result.

The part of the result I will focus on is the emotion. In other word, I will ask amazon: "An I happy?".

# Getting the picture from the webcam

I am using the package [github.com/blackjack/webcam](github.com/blackjack/webcam) to grab the picture.

## Capabilities of the webcam and image format

My webcam is handling the MJPEG format.
Therefore, after the creation of a _cam_ object and set the correct settings to grab mjpeg, I can read a frame in JPEG:

```go
// ...
cam, err := webcam.Open("/dev/video0") // Open webcam
// ...
// Setting the format:
_,_,_, err := cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
```

<code>format</code> is of type `uint32` and computable thanks to the informations present in [/usr/include/linux/videodev2.h](http://lxr.free-electrons.com/source/include/uapi/linux/videodev2.h)

MJPEG is: 1196444237

_Note_: To be honest, I did not evaluate the FOURCC method; I have requested the supported format of my webcam with their descriptions :)

## Grabbing the picture

In a endless `for` loop, a frame is read with a call to `ReadFrame`:

```go
for {
    timeout := uint32(5) //5 seconds
    err = cam.WaitForFrame(timeout)
    frame, err := cam.ReadFrame()
}
```

# AWS

The API to import to use the service is `github.com/aws/aws-sdk-go/service/rekognition` and is documented here: [http://docs.aws.amazon.com/sdk-for-go/api/service/rekognition/](http://docs.aws.amazon.com/sdk-for-go/api/service/rekognition/)

The operation that I am using to detect the emotion is [DetectFaces](http://docs.aws.amazon.com/sdk-for-go/api/service/rekognition/#Rekognition.DetectFaces) that takes an pointer to [DetectFacesInput](http://docs.aws.amazon.com/sdk-for-go/api/service/rekognition/#DetectFacesInput) with is composed of a pointer to an [Image](http://docs.aws.amazon.com/sdk-for-go/api/service/rekognition/#Image).

## Creating the input

The first thing that needs to be created is the [Image](http://docs.aws.amazon.com/sdk-for-go/api/service/rekognition/#Image) object from our `frame`:

```go
if len(frame) != 0 {
    image := &rekognition.Image{ // Required
        Bytes: frame,
    }
```

Then we create the DetectFacesInput:

```go
params := &rekognition.DetectFacesInput{
        Image: image,
        Attributes: []*string{
                aws.String("ALL"), 
        },
}
```

The `ALL` attributes is present, otherwise AWS does not return the complete description of what it has found.

## Sending the query

### Pricing notice and __warning__
The price of the service as of today is 1 dollar per 1000 request. That sounds cheap, but at 25 FPS, this may cost a lot.
Therefore, I have set up a read request that only process a picture if we press _enter_ 

```go
bufio.NewReader(os.Stdin).ReadBytes('\n')
```

### Session

As usual, to query AWS we need to create a session:

```go
var err error
sess, err = session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
if err != nil {
    fmt.Println("failed to create session,", err)
    return
}
svc = rekognition.New(sess)
```

_Note_: The `session` library will take care of connections informations such as environment variables like:

* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ ACCESS_KEY_ID`

### Query and result

Simply send the query 
```go
esp, err := svc.DetectFaces(params)

if err != nil {
        fmt.Println(err.Error())
        return
}
```

The result is of type [DetectFacesOutput](http://docs.aws.amazon.com/sdk-for-go/api/service/rekognition/#DetectFacesOutput).
This type is composed of a array of FaceDetails because obviously there can me more than one person per image.
So we will loop and display the emotion for each face detected:

```go
for i, fd := range resp.FaceDetails {
        fmt.Printf("The person %v is ", i)
        for _, e := range fd.Emotions {
                fmt.Printf("%v, ", *e.Type)
        }
        fmt.Printf("\n")
}
```

# Run:

<pre>
Resulting image format: MJPEG (320x240)
Press enter to process 
The person 0 is HAPPY, CONFUSED, CALM, 
</pre>

# Conclusion

That's all folks. The full code can be found [here](https://gist.github.com/owulveryck/33753125afa6284cd5dbbb1bd4d1eb54).

In the test I made, I was always happy. I've tried to be angry or sad, without success... Maybe I have a happy face.
I should try with someone else maybe.

The service is nice and opens the door to a lot of applications:
For example to monitor my home and sends an alert if someone is in my place and __not from my family__ (or not the cat :).
