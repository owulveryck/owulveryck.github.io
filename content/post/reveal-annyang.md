---
author: Olivier Wulveryck
date: 2016-09-09T13:39:50+02:00
description: People who know what they are talking about don't need Powerpoint said Steve Jobs. Indeed, it may be useful from time to time to have a single slide to display key words. But what I understand from Jobs's sentence is that relying on the slides is a bad idea. Good stuff that machine learning can pilot the slides for us!
draft: false
tags:
- chrome
- speech recognition
- Javascript
- revealjs
- slides
title: Being a better public speaker with a little help of Speech Recognition, Javascript and Chrome
topics:
- topic 1
type: post
---

# Introduction

I usually don't like slidewares.

Actually as IT engineer working, by now, exclusively in France, I'm facing the PowerPoint problem:

* Too many boring slides,
* too much information per slide,
* a presenter dedicated to read their content.

Therefore, the audience is watching its watch while waiting for a coffee break.

I won't redo the introduction I already did in a [previous post](/2016/06/23/websockets-reveal.js-d3-and-go-for-a-dynamic-keynote/index.html) but indeed slides can,
from time to time, be a value-add to a presentation.

Is the previous post, I introduced reveal.js, this excellent javascript framework. And I already did an interactive presentation.

In this post, I will go a little bit further with the integration of Speech Recognition.

# Speech Recognition by Google

## The Google Cloud Speech API

It ain't no secret now, I'm a big fan of machine learning.
Machines learn faster than people, and they can now assist us in a lot of boring tasks.

On the base of a neuron network, Google provides an [API for speech recognition](https://cloud.google.com/speech/).
It is fairly complete and multi lingual.

## Chrome speech recognition
They "_eat their own dog food_" and use their engine for android and... Chrome.
Best of all, they provide a userland access via a Javascript API to this functionality in Chrome.

This means that you can develop a web page that will access you microphone, sends what you say to the Google cloud, get the result back and process it in your page.

You can see an introduction [here](https://developers.google.com/web/updates/2013/01/Voice-Driven-Web-Apps-Introduction-to-the-Web-Speech-API)

# What can I use that for: A case study?

I had to do a presentation recently.
This presentation was about _Agility_ and _Devops_. The main idea was to give my client a feedback about experiences I had regarding those principles in digital transformation.

I didn't want to loose my audience with slides. But I wanted to keep the key concepts alive and visible.

So what I did was a one slide presentation only with the keywords I wanted to talk about.

The day before, I though:

> "How nice it would be if as far as I speak, the key concepts would appear on screen..."

You may think: _"that's easy, learn your script and click on the right moment"_.

Ok, but there are drawbacks:

* You have to learn the script
* You cannot be spontaneous anymore
* It's a one shot, one displayed, you cannot interact with the points anymore.

What I need is "an assistant" that will _listen to me_ and _act as soon as he/she/it hear a buzz word_.
It's 2016, this assistant is a machine, and I can teach it to act correctly.

Here is a little demo of the end product (don't pay to much attention to the content, I said nonsense for the demo)

<iframe width="560" height="315" src="https://www.youtube.com/embed/MOmmBufEwPo" frameborder="0" allowfullscreen></iframe>

And another one in French.

<iframe width="560" height="315" src="https://www.youtube.com/embed/3Uyr0G0add4" frameborder="0" allowfullscreen></iframe>

# The implementation

## Annyang

I have used the [annyang](https://github.com/TalAter/annyang) which is a javascript wrapper for the Chrome speech recognition functionnality.

Instead of matching a sentence as explained in the example of annyang, I made it listen to my whole speech in a loop.
Then I passed every sentence detected by the framework to a javascript function that was applying regexp to look for keywords.

### Displaying words

For every keyword I did a match to an object of my DOM and simply changed its <code>visibility</code> style from <code>hidden</code> to <code>visible</code>

Here is the javascript code:

```javascript
if (annyang) {
  annyang.debug(true);
  annyang.setLanguage('fr-FR');
  annyang.addCallback('result', function(phrases) {
  for (s of phrases) {
     str = s.toLowerCase();
     switch (str) {
        case (str.match(/communication/) || {}).input:nnyang.start();                                                                                              
        document.getElementById('communication').style.visibility = 'visible';
        ...
    // Start listening. You can call this here, or attach this call to an event, button, etc.
    annyang.start();       
    ...
```

And the corresponding html section for the communication keyword:

```html
<div class="reveal">
  <div class="slides">
    <section>
      <h3 style="visibility: visible" id="agile">Agile</h3>
      <p> <span style="visibility: hidden;" id="communication">communication</span> </p>            
```

The speech recognition engine detects the sentence and gives a confidence note about its recognition.
All the potential results are stored in an array (here <code>phrases</code>). I've used them all so I was more confident not to miss a word.

### Making them blink

As I was not fully confident in the solution (it was late in the night and the show was the next morning), Therefore I made a fall-back solution.
All the words were displayed, and I used a little CSS tweak to make them blink when they were pronounced. 
This was done by adding and removing a css class to the concerned node.
The logic remains the same.

```css 
/* Chrome & Safari */
@-webkit-keyframes blink {
  from {
    opacity: 1;
  }
  to {
    opacity: 0;
  }
}

.blink {
  -webkit-animation: blink 750ms 2;
}
```

```Javascript
case (str.match(/communication/) || {}).input:                                                                                       
   document.getElementById("b_communication").classList.toggle('blink');
   setTimeout(function () {
     document.getElementById("b_communication").classList.remove('blink');
   }, 1500);
   break;
```

# Conclusion and TODO

The speech recognition engine is efficient and usable.
What I need to do is to code a tiny javascript library in order to get a JSON associative array between the list of spoken words that would trigger an action for an element for example:

```json
{
  "tag": "communication",
  "words": ["communication", "communicate"],
  "fuction": "blink"
}
```

Anyway, as a quick and dirty solution, the goal is achieved.

Another Idea would be to plug this with a NLP engine to perform stemming or lemmatization to do a better decoding and be even more machine learning compliant. This could be done with the help of [MITIE](https://github.com/mit-nlp/MITIE)


