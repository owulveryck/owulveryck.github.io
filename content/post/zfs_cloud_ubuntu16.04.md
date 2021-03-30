---
categories:
- category
date: 2017-02-14T10:17:15+01:00
description: ""
draft: true
images:
- /assets/images/default-post.png
tags:
- tag1
- tag2
title: zfs_cloud_ubuntu16.04
---


I am using the device entries `/dev/xvd*` for testing purpose. Other choices may be best suitable for production.
Please refer to [ZFS on linux wiki](https://github.com/zfsonlinux/zfs/wiki/faq#selecting-dev-names-when-creating-a-pool) for more options.

```shell
sudo apt-get install zfsutils-linux

parted /dev/xvdb mklabel GPT
parted /dev/xvdc mklabel GPT
parted /dev/xvdd mklabel GPT

zpool create -m none -o ashift=12 tank raidz /dev/xvdc /dev/xvdb /dev/xvdd
```




RAIDZ is a little like RAID-5. I'm using RAID-Z1, meaning that from a 3-disk pool, I can lose one disk while maintaining the data access.

NOTE: Unlike RAID, once you build your RAIDZ, you cannot add new individual disks. It's a long story.

The -m none means that we do want to specify a mount point for this pool yet.

The -o ashift=12 forces ZFS to use 4K sectors instead of 512 byte sectors. Many new drives use 4K sectors, but lie to the OS about it for "compatability" reasons. My first ZFS filesystem used the 512-byte sectors in the beginning, and I had shocking performance (~10Mb/s write).


zpool create tank raidz /dev/xvdc /dev/xvdb /dev/xvdd 


