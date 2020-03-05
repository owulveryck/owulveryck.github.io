---
author: Olivier Wulveryck
date: 2016-03-31T10:23:02+02:00
description: How to setup RVM on an external drive on a Chromebook
draft: false
keywords:
- ruby
- rvm
- Chromebook
tags:
- Ruby
- Rvm
- Vagrant
- Bundler
- Chromebook
title: RVM from a USB stick on a Chromebook
topics:
- topic 1
type: post
---

# Introduction

#### Opening remarks

I'm not a Ruby developer, and I'm heavily discovering the ecosystem by now.
This are my notes, and if anything seems wrong to you, do not hesitate to send me remarks.

#### The scenario

For testing purpose, I wanted to play with vagrant-aws and more generally with ruby on my Chromebook.

Vagrant does not support _rubygems_ as installation method anymore ([see Mitchell Hashimoto's post](http://mitchellh.com/abandoning-rubygems))
and of course, there is no binary distribution available for the Chromebook.

So I have to install it from the sources.

The [documentation](https://github.com/mitchellh/vagrant/wiki/Installing-Vagrant-from-Source) says:

* Do __NOT__ use the system Ruby - use a Ruby version manager like rvm, chruby, etc

Alright, anyway I don't want to mess with my system and break Homebrew, so using RVM seems to be a good idea.

## Installing RVM

The RVM installation is relatively easy; simply running `curl -sSL https://get.rvm.io | bash` does the trick.
And then those commands make ruby 2.3.0 available via rvm:

```
$ source ~/.rvm/scripts/rvm  
$ rvm install 2.3.0
```

The stupid trick here is that everything is installed in my $HOME directory, and as my Chromebook is short on disk space: FS full !

Too bad.

## Using a USB stick

So my idea is to install the RVM suite onto a USB stick (because with me I don't have any SDHC card available).

### Preparing the stick

At first, the USB stick must be formatted in extendX (ext4) in order to be able to use symlinks, correct ownership etc.

```shell
sudo mkfs.ext4 -L Lexar /dev/sda1
```

__Note__: I've found that avoiding spaces in the volume name was good for rvm.


Once connected on the Chromebook, it's automatically mounted on `/media/removable/Lexar`.
The problem are the options: 

```shell
/dev/sda1 on /media/removable/Lexar type ext4 (rw,nosuid,nodev,noexec,relatime,dirsync,data=ordered)
```

the most problematic is `noexec` because I want to install executables in it.

So what I did was simply:

`sudo mount -o remount /dev/sda1 /media/removable/Lexar`

and that did the trick.

## Installing RVM on the USB

I will install rvm into `/media/removable/Lexar/rvm`. In order to avoid any ownership and permission problem I did:

```shell
mkdir /media/removable/Lexar/rvm
chown chronos:chronos /media/removable/Lexar/rvm
```

And then I created a simple `~/.rvmrc` file as indicated in the documentation with this:

```shell
$ cat ~/.rvmrc                                          
$ export rvm_path=/media/removable/Lexar/rvm
```

I also included this in my `~/.zshrc`

```shell
if [ -s "$HOME/.rvmrc"   ]; then
    source "$HOME/.rvmrc"
fi # to have $rvm_path defined if set
if [ -s "${rvm_path-$HOME/.rvm}/scripts/rvm"   ]; then
    source "${rvm_path-$HOME/.rvm}/scripts/rvm"
fi
```

## Installing rvm

the command I executed were then:

```
$ curl -sSL https://get.rvm.io | bash
$ source /media/removable/Lexar/rvm/scripts/rvm
$ rvm autolibs enable
$ rvm get stable
$ rvm install 2.3.0
```

And that did the trick

```
$ rvm list

rvm rubies

=* ruby-2.3.0 [ x84_64 ]

# => - current
# =* - current && default
#  * - default
```

## Testing with vagrant

### Cloning the vagrant sources

```shell
$ sudo mkdir /media/removable/Lexar/tools
$ sudo chown chronos:chronos /media/removable/Lexar/tools
$ cd /media/removable/Lexar/tools
$ git clone https://github.com/mitchellh/vagrant.git
```

### Preparing the rvm file for vagrant

To use the ruby 2.3.0 (that I've installed before) with vagrant, I need to create a .rvmrc in the vagrant directory:

```
$ cd /media/removable/Lexar/tools/vagrant
$ rvm --rvmrc --create 2.3.0@vagrant
```

### Installing bundler

The bundler version that is supported by vagrant must be <= 1.5.2 as written in the `Gemfile`. So I'm installing version 
1.5.2

```shell
$ cd /media/removable/Lexar/tools/vagrant
$ source .rcmrv
$ gem install bundler -v 1.5.2
```

### Compiling vagrant

Back to the vagrant documentation, what I must do is now to "compile it". To do so, the advice is to run:

```
$ bundle _1.5.2_ install  
```

(just in case several bundler are present )

I faced this error:

```shell
NoMethodError: undefined method `spec' for nil:NilClass
Did you mean?  inspect
An error occurred while installing vagrant (1.8.2.dev), and Bundler cannot continue.
Make sure that `gem install vagrant -v '1.8.2.dev'` succeeds before bundling.
```

According to google, this may be an issue with the version of bundler I'm using.
As I cannot upgrade the bundler because of vagrant, I've decided to take a chance and use
a lower version of Ruby

```shell
$ rvm install 2.2.0
$ rvm --rvmrc --create 2.2.0@vagrant
$ source .rvmrc
# and reinstalling bundler
$ gem install bundler -v 1.5.2            
$ bundle _1.5.2_ install
...
Your bundle is complete!
Use `bundle show [gemname]` to see where a bundled gem is installed.
```

# Voilà!

I can now use vagrant installed fully on the USB stick with

```shell
$ bundle _1.5.2_ exec vagrant
Vagrant appears to be running in a Bundler environment. Your 
existing Gemfile will be used. Vagrant will not auto-load any plugins
installed with `vagrant plugin`. Vagrant will autoload any plugins in
the 'plugins' group in your Gemfile. You can force Vagrant to take over
with VAGRANT_FORCE_BUNDLER.

You appear to be running Vagrant outside of the official installers.
Note that the installers are what ensure that Vagrant has all required
dependencies, and Vagrant assumes that these dependencies exist. By
running outside of the installer environment, Vagrant may not function
properly. To remove this warning, install Vagrant using one of the
official packages from vagrantup.com.
...
```

That's it for this post; next I will try to install vagrant-aws and play a little bit with it.

stay tuned.

