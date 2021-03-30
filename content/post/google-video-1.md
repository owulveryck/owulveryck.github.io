---
categories:
date: 2017-06-01T22:07:56+02:00
description: "A very brief article about my first call to the new service of Google Cloud Plateform: Video Intelligence. Caution: The video used in this example is #NSFW"
draft: false  
images:
- https://cloud.google.com/images/products/video-intelligence/analysis.png
tags:
title: Analyzing a parodic trailer (NSFW) with Google Cloud Video Intelligence
---

Google has recently announced its new service called "[Google Cloud Video Intelligence](https://cloud.google.com/video-intelligence/)".
The purpose of this service is to offer tagging and annotations of digital videos.

I will try this service on a trailer of a French parody. This movie is made of several scenes taken from erotic movies of the seventies.

Why this parody?

* because it is fun
* because it is composed of a lot of different scenes
* because it is short (so it won't cost me a lot)
* because, as it is related to erotic of the seventies, I am curious about the result!

_Caution_: this video **is not a porn video**, but is indeed **not safe for work** (_#nsfw_)

# What information can the service find?

## Shot change detection

This feature will detect the different scenes and display their time range. There is no further analysis of the scene. That is to say that it won't tell, by now, that the first scene is about a sports competition. But indeed it will describe that the first scene occurs from the first microsecond until the xxx microsecond and so on.

## Label detection

The more interesting feature is the label detection.

With this operation, the service will display tags of any element found in the video, as well as the time range of the video where they can be seen.

For example, it may tell you that there is a dog in the video between x and y micro-seconds as well as between w and z micro-seconds.

# Preparing the video

I have downloaded the video, thanks to [youtube-dl](https://rg3.github.io/youtube-dl/) and I have uploaded it to [Google Cloud Storage](https://cloud.google.com/products/storage/) as the API expects the video to be here. There may be a way to post the video encoded in base64 directly, but that would have been less convenient for my tests.

![screnshot](/assets/video-intelligence/gs-trailer.png)

# Querying Google Cloud Video Intelligence 

This test is made with the simple REST API with `curl`.

## Preparing the request

To actually use the API, we need to perform a POST request. 
The payload is a simple JSON file where we specify:

* the URI of the video file to process
* an array of features to use among: Shot change detection and/or label detection

Here is my payload. I want both features for my test:

{{< highlight js >}}
{
    "inputUri": "gs://video-test-blog/trailer.mp4",
    "features": ["SHOT_CHANGE_DETECTION","LABEL_DETECTION"]
}
{{</ highlight >}}

## Launching the request

### Authorization

To actually use the service, I need an authorization token. This token is linked to a service account.
Then with the token, we can trigger the analysis by using this `curl` command:

{{< highlight shell >}}
curl -s -k -H 'Content-Type: application/json' \
      -H 'Authorization: Bearer MYTOKEN' \
      'https://videointelligence.googleapis.com/v1beta1/videos:annotate' \
      -d @demo.json
{{</ highlight >}}

The action replies with a JSON containing an `operation name`. Actually, the operation is long and asynchronous. This `operation name` can be used to get the processing status.

{{< highlight js >}}
{
   "name": "us-east1.16784866925473582660"
}
{{</ highlight >}}

### Getting the status

To request the status, we need to query the service to get the status of the operation:

{{< highlight shell >}}
curl -s -k -H 'Content-Type: application/json' \
      -H 'Authorization: Bearer MYTOKEN' \
      'https://videointelligence.googleapis.com/v1/operations/us-east1.16784866925473582660'
{{</ highlight >}}

It returns a result in `json` that in which we can find three important fields:

* `done`: a boolean that tells whether the processing of the video is complete or not
* `shotAnnotations`: an array of the shot annotations as described earlier
* `labelAnnotations`: an array of label annotations

Here is a sample output: (the full result is [here](/assets/video-intelligence/video-analysis-a-la-recherche.json))
{{< highlight js >}}
{
  "response": {
    "annotationResults": [
      {
        "shotAnnotations": [
          // ...
          {
            "endTimeOffset": "109479985",
            "startTimeOffset": "106479974"
          }
        ],
        "labelAnnotations": [
          // ... 
          {
            "locations": [
              {
                "level": "SHOT_LEVEL",
                "confidence": 0.8738658,
                "segment": {
                  "endTimeOffset": "85080015",
                  "startTimeOffset": "83840048"
                }
              }
            ],
            "languageCode": "en-us",
            "description": "Acrobatics"
          },
        ],
        "inputUri": "/video-test-blog/trailer.mp4"
      }
    ],
    //...
  },
  "done": true,
  //...
}
{{</ highlight >}}

# Interpreting the results

## Tag cloud

I will only look at the label annotations.
The API has found a lot of label described under the `description` fields and 1 to N location where such a description is found.

What I can do is to manipulate the data to list all the label with their frequency.

You can find [here](https://gist.github.com/owulveryck/70d97e1e73d664c1c927c253a862ac17) a little go code that will display labels as many times as they occur.

For example:

```
Abdomen Abdomen Abdomen Acrobatics Action figure Advertising Advertising ...
```

This allows me to generate a tag cloud with the help of [this website](https://www.jasondavies.com/wordcloud/):

So here is the visual result of what the service has found in the video:

![tag cloud](/assets/video-intelligence/wordcloud.png)

## Annotated video

To find out where the labels are, I made a little javascript that display the elements alongside of the youtube video.
Just click on the video and the tags will be displayed below.

<button id="launchyt">It is safe to watch the video, please show the result!</button>

<div id="player"></div>

<ul id="labels"></ul>

<ul>
    <li id="result1"></li>
    <li id="result2" style="color: #8A8A8A;"></li>
    <li id="result3" style="color: #9E9E9E;"></li>
    <li id="result4" style="color: #B2B2B2;"></li>
    <li id="result5" style="color: #C6C6C6;"></li>
</ul>

<script type="text/javascript" async src="/assets/video-intelligence/app.js"></script>

# Conclusion

There is a lot more to do than simply displaying the tags.
For example, We could locate an interesting tag, take a snapshot of the video, and use the photo API to find websites related to this part of the video.

For example, in this video, it can be possible to find the original movies were people are dancing for example.

I will postpone this for another geek-time.

_P.S._ The javascript has been made with _gopherjs_. It is not optimize at all (I should avoid the encoding/json package for example). If you are curious about the implementation, the code is [here](/assets/video-intelligence/main.go), [here](/assets/video-intelligence/structure.go) and [here](/assets/video-intelligence/data.go).
