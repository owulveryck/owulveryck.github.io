---
author: Olivier Wulveryck
date: 2016-01-24T15:25:54+01:00
description: The setup of my new Freeebsd
draft: true
keywords:
- FreeBSD
tags:
title: Setting up my new BSD box
topics:
- topic 1
type: post
---
# About the server

I have subscribed for a now dedicated BSD box.
The provider is OVH, and this box is hosted in Canada.

It's been a few years since I first subscribed a box a OVH, and I really enjoy their service.

The main usage for this box is a `geek box`. I use it to exeperiment some stuffs, I cannot experiment into my home PC.
Therefore I usually create a bunch of `jails`, and each jail is dedicated to a task

Here are the informations of my new box:

```shell
~ uname -a
FreeBSD localhost 10.2-RELEASE-p9 FreeBSD 10.2-RELEASE-p9 #0: Thu Jan 14 01:32:46 UTC 2016
root@amd64-builder.daemonology.net:/usr/obj/usr/src/sys/GENERIC  amd64
```

## ZFS

My root is a ZFS pool named `zroot`

# Basic installation

This box is a `10.2 release` therefore it uses the "new" `pkg` tool instead of the legacy `pkg_*` tools.

## ZSH

### Installation
```shell
pkg install zsh
```
### Oh-my-zsh

Install `git` 

```shell
pkg install git
```

```shell
sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
```

# Setting up openvpn

## Installation
```shell
pkg update
pkg install openvpn
```

## Configuration
```shell
mkdir /usr/local/etc/openvpn
cp /usr/local/share/examples/openvpn/sample-config-files/server.conf /usr/local/etc/openvpn/server.conf 
```

The default options are used, and only user and groups ares set to nobody for security reasons

```shell
...
# It's a good idea to reduce the OpenVPN
# daemon's privileges after initialization.
#
# You can uncomment this out on
# non-Windows systems.
user nobody
group nobody
...
```

### Generating the keys

```shell
cp -r /usr/local/share/easy-rsa /usr/local/etc/openvpn/easy-rsa
```

### Generating a client certificate

I will generate an openvpn configuration for my chromebook

```shell
cd /usr/local/etc/openvpn/easy-rsa/
./build-key chromebook
```

Then generate the `ovpn` config file

```shell
export CLIENTNAME=chromebook
cp /usr/local/share/examples/openvpn/sample-config-files/client.conf ./client.conf
# Modifying remote...
cp client.conf $CLIENTNAME.ovpn
printf "\n<ca>\n" >> ./$CLIENTNAME.ovpn && \
cat ./ca.crt >> ./$CLIENTNAME.ovpn && \
printf "</ca>\n" >> ./$CLIENTNAME.ovpn && \
printf "\n<cert>" >> ./$CLIENTNAME.ovpn && \
grep -v '^ ' ./$CLIENTNAME.crt | grep -v 'Certificate' >> ./$CLIENTNAME.ovpn && \
printf "</cert>\n" >> ./$CLIENTNAME.ovpn && \
printf "\n<key>\n" >> ./$CLIENTNAME.ovpn && \
cat ./$CLIENTNAME.key >> ./$CLIENTNAME.ovpn && \
printf "</key>\n" >> $CLIENTNAME.ovpn
```

### Firewall

#### Firewall and routing

```shell
# default openvpn settings for the client network
vpnclients = "10.8.0.0/24"
#put your wan interface here (it will almost certainly be different)
wanint = "em0"
# put your tunnel interface here, it is usually tun0
vpnint = "tun0"
# OpenVPN by default runs on udp port 1194
udpopen = "{1194}"
icmptypes = "{echoreq, unreach}"

set skip on lo
# the essential line
nat on $wanint inet from $vpnclients to any -> $wanint

block in
pass in on $wanint proto udp from any to $wanint port $udpopen 
pass in on $wanint proto tcp from any to $wanint port 22 
# the following two lines could be made stricter if you don't trust the clients
pass out quick 
pass in on $vpnint from any to any
pass in inet proto icmp all icmp-type $icmptypes
```

```shell
~ /etc/rc.conf
...
openvpn_enable="YES"
openvpn_configfile="/usr/local/etc/openvpn/server.conf"
```

# Rescue...

Of course, I forgot one rue in my pf.conf and therefore I could not access to my box anymore

## Maganer
Boot into rescue mode

```shell
rescue-bsd# zpool import zroot
rescue-bsd# zpool list
internal error: failed to initialize ZFS library
```

That's because I did import the zroot into /

```shell
zpool import -o altroot=/mnt zroot
```
