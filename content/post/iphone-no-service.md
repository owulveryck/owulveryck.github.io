---
date: 2017-08-14T09:42:12+02:00
description: "This is the story of a lambda person that had his iPhone broken after an update... And then the story of a geek, who has been told that he needed to pay 350€ for a replacement, based on assumptions and lies."
draft: true
images:
- /assets/images/iphone-blog-logo.png
title: When I apply a 350€ update on the iPhone
---

I am an apple client for years now. I was very pleased with the apple customer care until the day I really needed something from them.

This post is the "why". But as you are on a technical blog the idea behind this post is:

* to explain the reason of wrath that made me dig into the deep understanding of my phone and
* to give technical facts and hints to any geek that would like to investigate an iPhone issue based on what I learned (this is, according to me, the most interesting part)

# Part I:  The facts

I have add an iPhones, iPads and MacBook. I made this choice for different reasons amongst:

* the quality of the products
* the very user friendly interface
* the security of the phone
* the implicit contract that apple provides: "nevemind, apple will take care"

Of course, those reasons are personals and can be argued (and subject to _troll_), but this is not the point of this article.

## Updating the iOS

Apple releases updates of its iOS several times a year. As a geek, I usually read the changelog and the forums to see if it's worth applying the update as soon as it is out.
On July the 19th, Apple released iOS 10.3.3 and the changelog mentioned several security issues that were corrected with this version.
An after all, I was on holiday, and it was definitely not a good moment to check, analyse and do an update.

When I came home, I planned an update campaign at home. As usual, I start by my own phone, then I update my wife's iPhone and iPad.
This is a trick I use to avoid any inflexion in the [Wife Acceptance Factor](https://en.wikipedia.org/wiki/Wife_acceptance_factor). 
If I notice something wrong on my phone, the update campaign is stopped...

Anyway, when I updated successfuly my own phone, I proceed as planed.
I triggered the update of my wife's iPhone 6 and I left home to go to work.

At 9:30 I received a mail:

```
From: Darling
To: me 
Subject: Phone's not working

Contact me via email
```

I triggered `hangout` so we could chat and asked what was going on. She told me that the phone was searching for the carrier.
I asked her to do simple manipulations:

<center>
![IT helpdesk](/assets/images/have-you-tried.jpg)
</center>

Nothing worked.

In the evening I tried to reinstall the iOS. Then I tried to downgrade. I also tried to install a beta version of iOS11... Nothing worked.

I asked "google". I noticed that this problem was common. Moreover, it could happen after any version update. 
So I went on twitter, I found two people with the exact same problem and I asked them how they solved it:

<center>
{{< tweet 893320425534619648 >}}
</center>

Ok, so I contacted the support via twitter, and they gave me a phone appointment with an advisor. So far, so good...
We did a couple of manipulations such as _turning it off and on again_ 

They made me reinstall the software (again) and nothing worked. They say that I should go to an Apple store for a further diagnostic. I thought I could just walk into an Apple Store with my phone but NO! You need an appointment to meet an advisor in what they call "_The genius bar_". The genius bar... I was so excited to finally meet a genius!

<center>
i![genius](/assets/images/genius.jpg)
</center>

# Part II: Apple (does not) Care

I took the appointment over lunch time. I arrived on time, but due to a misunderstanding, they told me that I was 20 minutes too early... Indeed I have waited and an advisor finally arrived with a macbook. He told me:

"What I will do is to reset your phone and install the ios again"

I told him that has already been done, but he replied: this is the procedure, and it will only take a few minutes here.
10 minutes later, the iphone is restored and the problem was still there.

Then he simply told me:

"_hardware problem. Your phone is out of warranty, but I have a solution for you. For only 350€, you can have a new phone that will have a 3 months warranty!_"

I replied that this was ridiculous because my phone was working perfectly before the update. Therefore I refuse to pay this fee. So he told me that I should probably talk to a supervisor. 

15 minutes later a supervisor came... (it was already 1:30pm, and I was late for work).
He told me that AppleCare did not order a phone replacement and that there was nothing he could do. So I ask:

"_Can I call AppleCare back, so they may tell you to give me a new phone?_" (you know, I was still at the time convinced that Apple was customer centric...)

So he gave me a phone with a direct line to a supervisor at AppleCare.  This gentleman told me: 

"_Your phone is broken, and the hardware problem cannot be linked with the software update in any way. Nobody has ever seen a software that could break a hardware_"

I replied that, in my experience, it is something that can happen, and I took as the example of an [EPROM](https://en.wikipedia.org/wiki/EPROM). A bad update simply breaks it. And, even if I didn't know the iPhone architecture, I thought that similar updates were done in some components of the phone.
But anyway, nobody has given me any proof of a hardware problem. The problem can be software related.

Well, we argued for 35 minutes. At the very end he killed me by saying:

_Anyway, I won't do any exception and I will not replace your phone... But even if I wanted to (it is 3 clicks for me), by now you are in an Apple Store and only a supervisor of the store can take such a decision_

So who is lying... Actually, I didn't care, I asked to talk to another supervisor. That was more than 2 hours that I was in the shop. All I wanted was a solution or at least a real diagnose of my phone that could tell me which piece was broken and why...

I went back home, disappointed... I looked over the internet for a solution because I was not in the exception list and Apple would not do me a courtesy. 

In fact, both people from Apple Care and from the Apple Store told me that they could be exceptions. Of course, they did not want to tell me the reasons.
In a similar situation, some people reported that their [iPhone has been changed for free](https://discussions.apple.com/message/31836225#message31836225). Taking different treatment depending on who is asking heavily sounds like discrimination to me: 

<center>
![discrimination](/assets/images/discrimination.png)
</center>

Anyway, let's move on and "seek for the truth". And check whether there is really a broken part in my phone.

# Part III: Geek (at last!)

__Disclaimer__: By now, I am only playing with official tools. I am waiting for an answer from Apple. If they still refuse to change the phone, I may jailbreak it to do further investigations. I am blogging this because what took me hours to understand may help someone with a similar problem. As of today, I do not have any solution nor I have found the problem.

[Baseband device](https://www.theiphonewiki.com/wiki/Baseband_Device)

The qualcomm modem is a [MDM9625](https://www.theiphonewiki.com/wiki/MDM9635)

[https://developer.apple.com/bug-reporting/profiles-and-logs/](https://developer.apple.com/bug-reporting/profiles-and-logs/)


## Is the SIM card ok

`iPhone CommCenter[75] <Notice>: #N SIM initialization complete; all essential information available #sim-ready`

`misd[31] <Notice>: carrier service is available`


`iPhone CommCenter[75] <Notice>: #I Unbrick device was successful `
The params are displayed with the profile update
` #I Received activation info: <private>       `

`16873.516 [I] evt: Firing event 'recalculateConnectionAvailability': with params= 0, Postponement status change  `
