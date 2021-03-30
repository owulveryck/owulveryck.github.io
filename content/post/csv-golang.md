---
author: Olivier Wulveryck
date: 2016-09-11T20:27:49+02:00
description: Which language gives the best performances to process a big/huge csv file between Python, Perl and Golang
draft: false
keywords:
- golang
- csv
title: Processing CSV files with golang, python and perl
topics:
- topic 1
type: post
---

# Introduction

System engineers are used to CSV files.
They may be considered as a good bridge between Excel and the CLI.
They are also a common way to format the output of a script so data can be exploited easily from command lines tools such as _sort_ , _awk_,_uniq_, _grep_, and so on.

The point is that when there is a significant amount of data, parsing it with shell may be painful and extremely slow.

This is a very simple and quick post about parsing CSV files in python, perl and golang.

# The use case

I consider that I have a CSV file with 4 fields per row.
The first field is a server name, and I may have 700 different servers.
The second field is a supposed disk size for a certain partition. The other fields are just present to discriminate the rows in my example

What I would like to know is the total disk size per server.

I will implement three versions of parsing, and I will look for the result of a certain server to see if the computation is ok.
Then I will compare the exeuction time of each implementation


## Generating the samples
I'm using a very simple shell loop to generate the samples. I'm generating a file with 600000 lines.

{{< gist owulveryck 4f9ddb952c5f1ef708b60a9907733969 "Generation.sh" >}}

I've randmly chosen to check the size of SERVER788 (but I will compute the size for all the servers).

I have a lot of entries for my server.
```bash
grep SERVER788 sample.csv| wc -l
3012
```
## The implementations
Here are the implementation in each language:

### The go implementation
The go implementation relies on the <code>encoding/csv</code> package.
The package has implemented a `Reader` method that can take the famous `io.Reader` as input. Therefore I will read a stream of data and not load the whole file in memory.

{{< gist owulveryck 4f9ddb952c5f1ef708b60a9907733969 "main.go" >}}

### The perl implementation
I did not find a csv implementation in perl that would be more efficient than the code below. Any pointer appreciated.
{{< gist owulveryck 4f9ddb952c5f1ef708b60a9907733969 "main.pl" >}}

### The python implementation
Python does have a _csv_ module. This module is optimized and seems to be as flexible as the implementation of go. It reads a stream as well.
{{< gist owulveryck 4f9ddb952c5f1ef708b60a9907733969 "main.py" >}}

## The results
I've run all the scripts through the _GNU time_ command. I didn't used the built-in time command because 
I wanted to check the memory footprint as well as the execution time.

Here are the results
{{< gist owulveryck 4f9ddb952c5f1ef708b60a9907733969 "result.sh" >}}

# Conclusion

All of the languages have very nice execution time: below 4 seconds to process the sample file. Go gives the best performances, but it's insignificant as long as the files do not exceed millions of records.
The memory footprint is low for eache implementation.

It's definitly worth a bit of work to implement a decent parser in a "modern language" 
instead of relying on a `while read` loop or a `for i in $(cat...` in shell.
I didn't write a shell implementation, but it would have take ages to run on my chromebook anyway.
