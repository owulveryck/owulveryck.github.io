---
title: "Rethinking Presentations: Beyond Static Slides"
date: 2023-10-10T07:55:21+02:00
lastmod: 2023-10-10T07:55:21+02:00
draft: false
images: [/assets/crowdasleep_small.png]
videos: [/assets/present.webm]
keywords: []
summary: In the digital age, traditional PowerPoint presentations often fail to engage audiences due to their static nature. 
  

  Research suggests that audience attention dwindles after just 10 minutes without engagement. 
  
  The proposed method in this article aims to revitalize presentations by
  
  * Incorporating live doodling with the assistance of tools like the reMarkable tablet for real-time interaction.
  
  * Using a script to create PDFs from images, blending the familiar structure of slides with spontaneous on-the-spot drawing.
   

  The result is a more authentic, engaging, and impactful presentation experience, though it requires deeper preparation and understanding of the topic. 
  
  The tools and methods highlighted aim to shift the focus from aesthetics to genuine content engagement.
tags: ["keynote", "reMarkable", "powerpoint", "presentation", "talk"]
categories: ["tools"]
author: "Olivier Wulveryck"

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: false
toc: true
autoCollapseToc: true
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: false

# Uncomment to add to the homepage's dropdown menu; weight = order of article
# menu:
#   main:
#     parent: "docs"
#     weight: 1
---

## Introduction
In today's fast-paced world, traditional presentations often fall short of engaging audiences. 

Their static and pre-prepared nature lacks the interactivity and dynamism necessary to capture attention.

We no longer expose "_power points_" or "_key notes_" as we did in the days of _transparencies_ and _overhead projectors_ 
(see the impressive presentation [_growing a language_](https://www.youtube.com/watch?v=_ahvzDzKdB0) by _Guy L. Steele Jr._ as an illustration, even though the content is unrelated with this article, his use of transparencies is interesting).

**Slides** have evolved **from a supportive role** to becoming the main focus, where **they are the presentation**.

I've sat through countless dull presentations where attendees are often found dozing off in front of endless static slides.

![Illustration of a crowd dozing off while listening to a monotonous presentation](/assets/crowdasleep_small.png)

In contrast, at tech conferences where I have the pleasure of witnessing live demos, the audience is noticeably more enthusiastic. It's no surprise that the _scrum_ framework incorporates demos.

However, not every presentation can be demo-centric. So, I'm exploring a new approach to enhance interactivity and engagement while maintaining a structured format.

This article explores my approach and the tools I use. While I don't claim to have **the definitive method** (as many others exist), I'm sharing my personal techniques on this tech blog in the hopes they might be helpful to someone.

## The situation 

### The science behind engagement

According to Dr. John Medina's "Brain Rules," attention begins to wane after about 10 minutes without engagement[^1].

Thus, a talk should either be shorter than 10 minutes, or there should be efforts to re-engage the audience.

A study conducted by Carnegie Mellon University demonstrated that active learning methods significantly improve academic performance, emphasizing interactive settings over traditional lectures[^2].

Therefore, in a talk lasting longer than 10 minutes, engaging the audience through a process of active learning may help in maintaining attention.

### Applying the science

I own a reMarkable tablet, and, as previously explained on my blog, I use it for video calls.

Now that the full remote period is over and people are returning to conferences, I can extend this method of presentation to use it live, in real life.

I have already utilized it in a series of BBLs as a complement to some slides, and the feedback has been positive.

I can take it a step further by replacing traditional slides with elements stored on the reMarkable, allowing me to write on them, while leaving blank pages in between to draw complementary information.

I can prompt the audience to react and adapt the drawing, fostering active learning.

![](/assets/present_looped.webp)

### The Trade-offs

While traditional slides provide a safety net for presenters, aiding them in staying on track, live doodling necessitates thorough preparation. 
One must have a deep understanding of their topic to maintain fluidity and confidence throughout the presentation.

And let's be honest, while I am relatively fluid when I need to explain something I am proficient in to a small group, presenting in front of a crowd poses a different challenge.
I can hardly rely on feedback (whether intentional or not).

Thus, I need to maintain a certain frame provided by the tooling.

## Tooling and engineering: Bringing the Presentation to Life

Now let's explore my toolbox. I will try it tomorrow at [Cloud Nord](https://cloudnord.fr/) in my hometown. I know that people are nice and friendly here, they will forgive any mistake.

### Streaming with reMarkable

One of the key tools enabling this dynamic presentation style is the reMarkable tablet.

It's not merely a device for note-taking or reading; it's a potent streaming tool that brings the essence of live doodling to the forefront.

With [goMarkableStream](https://github.com/owulveryck/goMarkableStream), I can easily stream the content of the tablet.

Recently, I added a feature that allows streaming over the Internet by embedding the [ngrok](https://ngrok.com/) tunneling capability.
This enables streaming on a computer that may not be on the same network as the tablet, although it requires decent internet connectivity.

I attempted tethering over a poor mobile network, and the result was not reliable enough for a presentation. However, if the room has adequate wifi, it can suffice.

Nevertheless, whenever possible, I adhere to a wired connection to my laptop and stream from my laptop. It's a safer choice.

### Framing the presentation in a PDF

To blend the traditional with the new, I employ a script that converts a collection of images into a PDF.

This simulates the familiar slide structure but with a twist – these "slides" can be written on, ensuring that while there's a base for the presentation, the real-time interaction remains intact.

[This gist on GitHub](https://gist.github.com/owulveryck/1317f9b22433aa18778b673000159141) introduces a `Makefile` that takes a set of images and converts them into a format suitable for the reMarkable.
It converts the images to grayscale and to a resolution of 1872x1404.

Then, it adds a border and an annotation at the bottom of each page (just because I like the style).

Afterward, it assembles the images into a PDF.

To achieve the correct order, the script reads the images from a file named `slides.txt`.

I can also add a blank page (and the blank page is also generated by the script).

I will attempt to share the result here once I have created a complete presentation with it.

## Conclusion

The age-old saying, "Form follows function," holds true even for presentations. 
While the traditional method prioritizes aesthetics, the dynamic approach emphasizes content and genuine engagement.

It re-centers the presentation on communicating the essence of a topic instead of merely reading slides.

This approach may demand more effort but promises a more authentic and impactful presentation experience.

Another possibility is to record a presentation for broadcasting... This is actually what I did for [this illustration](/assets/present.webm).

---

### References

[^1]: [Brain Rules, 10-minute rule](https://brainrules.blogspot.com/2009/03/10-minute-rule.html)

[^2]: [New Research Shows Learning Is More Effective When Active](https://www.cmu.edu/news/stories/archives/2021/october/active-learning.htm)
