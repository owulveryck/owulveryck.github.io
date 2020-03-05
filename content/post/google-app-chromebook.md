+++
date = "2015-10-26T10:41:57Z"
draft = false
title = "Developping \"Google Apps\" on my Chromebook"
tags = [
    "chromebook",
    "configuration"
]
+++

It is a week now that I'm playing with my chromebook.
I really enjoy this little internet Terminal.

I "geeked" it a little bit and I installed my favorites dev tools eg:

* [The solarized theme for the terminal](https://gist.github.com/johnbender/5018685)
* `zsh` with [Oh-my-zsh](https://github.com/robbyrussell/oh-my-zsh)
* `tmux` (stared with `tmux -2` to get 256 colors)
* `git`
* `vim`
* a `Go` compiler
* The [`HUGO`](http://gohugo.io/overview/quickstart/) tool to write this blog.


All of it has been installed thanks to the "brew" package manager and following [those instructions](https://github.com/Homebrew/linuxbrew/wiki/Chromebook-Install-Instructions).

## Google Development Environment

I've installed the Google Development Environement as described [here](https://cloud.google.com/appengine/docs/go/gettingstarted/devenvironment).

Python 2.7 is a requirements so I `brewed it` without any noticeable issue.

When I wanted to serve locally my very first Google App developement, I ran into the following error:

```
~ go app serve $GOPATH/src/myapp
...
ImportError: No module named _sqlite3
error while running dev_appserver.py: exit status 1
```

Too bad. I've read that this module should be built with python, but a even a `find /` (I know it's evil) didn't return me any occurence.

So, I have:

* Googled 
* reinstalled sqlite with `brew reinstall sqlite`
* reinstalled python with `brew reinstall python`
* played with brew link, unlink and so
* ...

Still no luck!

I've also tried the compilation with a `verbose` option, and I the log file, there is an explicit message:

```
Python build finished, but the necessary bits to build these modules were not found:
_bsddb  _sqlite3_tkinter
...
To find the necessary bits, look in setup.py in detect_modules() for the modules name.
```


That's where I am now, stuck with a stupid python error. I'd like the folks at google to provide a pure go developement enrironement that would avoid the bootstraping problems.

I'll post an update as soon as I have solved this issue !

*EDIT*:

I've had a look in the `setup.py` file. To compile the sqlite extension, it looks into the following paths:

```
...
sqlite_incdir = sqlite_libdir = None
sqlite_inc_paths = [ '/usr/include',
                     '/usr/include/sqlite',
                     '/usr/include/sqlite3',
                     '/usr/local/include',
                     '/usr/local/include/sqlite',
                     '/usr/local/include/sqlite3',
                   ]
...
```

But in my configuration, the libraries are present in `/usr/local/linuxbrew/*`. Hence, simply linking the include and libs did the trick

I'm now facing another error when I try to run the `goapp serve` command:

```
...
AttributeError: 'module' object has no attribute 'poll'
error while running dev_appserver.py: exit status 1
```

Google told me, that on OSX the poll system call is broken and has been disabled.
As brew is mainly developped on MacOS, that may be the reason

I've recompiled the python with the `--with-poll` option and that did the trick.

## Finally

Here are my options for compiling python:

```
~ brew reinstall python --with-brewed-openssl --with-brewed-sqlite --with-poll 
...
Warning: The given option --with-poll enables a somewhat broken poll() on OS X (https://bugs.python.org/issue5154)  Formula git:(master)).
...
```

And the `goapp serve` is finally working on my Chromebook:

```
~ goapp serve /home/chronos/user/GOPROJECTS/src/github.com/owulveryck/google-app-example/
INFO     2015-10-26 15:48:04,840 devappserver2.py:763] Skipping SDK update check.
INFO     2015-10-26 15:48:04,935 api_server.py:205] Starting API server at: http://localhost:54116
INFO     2015-10-26 15:48:06,092 dispatcher.py:197] Starting module "default" running at: http://localhost:8080
INFO     2015-10-26 15:48:06,096 admin_server.py:116] Starting admin server at: http://localhost:8000
INFO     2015-10-26 15:48:16,700 shutdown.py:45] Shutting down.
INFO     2015-10-26 15:48:16,701 api_server.py:648] Applying all pending transactions and saving the datastore
INFO     2015-10-26 15:48:16,701 api_server.py:651] Saving search indexes
```
