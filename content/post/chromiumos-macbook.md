+++
images = ["/2016/10/image.jpg"]
description = ""
categories = ["category"]
tags = ["tag1", "tag2"]
draft = true
title = "Chromiumos Macbook"
date = 2018-01-06T09:24:39+01:00
+++


From ubuntu 14.04
first update

[ChromiumOS developer guide](https://www.chromium.org/chromium-os/developer-guide)

`sudo apt-get install git-core gitk git-gui curl lvm2 thin-provisioning-tools python-pkg-resources python-virtualenv`

# Install depot_tools

Description:

> The Chromium depot_tools suite contains many tools to assist/augment the Chromium development environment. The tools may be downloaded from here.
> There are two primary categories of tools. Some of these tools are essential to the development flow, in the sense that you cannot successfully develop Chromium without them.
> Other tools are merely helper tools. Not required, but they can substantially ease the development workflow.
> A listing of both categories of tools follows.

## Setup 

[depot_tools_tutorial(7) Manual Page](http://commondatastorage.googleapis.com/chrome-infra-docs/flat/depot_tools/docs/html/depot_tools_tutorial.html#_setting_up)

## Creating an overlay

the `make.defaults`

Kernel 4.14

USE="${USE} legacy_keyboard legacy_power_button sse kernel-4_14"


the firmwares

