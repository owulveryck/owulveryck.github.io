---
author: Olivier Wulveryck
date: 2016-08-29T21:58:17+02:00
description: A first geek interaction between my  raspberry pi 3 and my weather station
draft: false
tags:
- libusb
- udev
- linux
- weather station
- oregon RMS300
- golang
- raspberry pi
title: Getting weather data from the station to the raspberry
topics:
- topic 1
type: post
---

# Introduction

A bunch of friends/colleagues offered me a <a type=amzn>raspberry pi 3</a>.
It may become my VPN gateway, or my firewall, or the brain of my CCTV, or maybe the center of an alarm.... Maybe a spotify player...

Anyway, I have installed raspbian and I'm now playing with it.

Yesterday evening, as I was about to go to bed, I've had a very bad idea... I've linked together my <a type=amzn>rpi</a> and my <a type="amzn">Oregon</a> Weather Station.
3 hours later, I was still geeking...

As usual in the blog I will explain what I did, what did work, and what did not.

# Attaching the devices

I've plugged the device, ok! Now what does the system tells me about it:

What `dmesg` tells me is simply

<pre>
[ 2256.877522] usb 1-1.4: new low-speed USB device number 5 using dwc_otg
[ 2256.984860] usb 1-1.4: New USB device found, idVendor=0fde, idProduct=ca01
[ 2256.984881] usb 1-1.4: New USB device strings: Mfr=0, Product=1, SerialNumber=0
[ 2256.984894] usb 1-1.4: Product:  
[ 2256.992719] hid-generic 0003:0FDE:CA01.0002: hiddev0,hidraw0: USB HID v1.10 Device [ ] on usb-3f980000.usb-1.4/input0
</pre>

## Finding the device

`lsusb` gives me the list of the usb devices on my <a type=amzn>rpi</a>:

<pre>
# lsusb 
Bus 001 Device 004: ID 0fde:ca01  
Bus 001 Device 003: ID 0424:ec00 Standard Microsystems Corp. SMSC9512/9514 Fast Ethernet Adapter
Bus 001 Device 002: ID 0424:9514 Standard Microsystems Corp. 
Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
</pre>

The first one correspond to my weather station but it belongs to root:

<pre>
# ls -lrt /dev/bus/usb/001/004
crw------- 1 root root 189, 3 Aug 30 12:45 /dev/bus/usb/001/004
</pre>

## Giving access: `udev`

The first thing to do is to allow access to my usb device so I won't need to run any program as root.
By default the `pi` user belongs to a bunch of groups. One of those is called `plugdev`.
It is the one I will use for my experiment.

### Get information about my Device

<pre>
# udevadm info /dev/bus/usb/001/004

P: /devices/platform/soc/3f980000.usb/usb1/1-1/1-1.3
N: bus/usb/001/012
E: BUSNUM=001
E: DEVNAME=/dev/bus/usb/001/012
E: DEVNUM=012
E: DEVPATH=/devices/platform/soc/3f980000.usb/usb1/1-1/1-1.3
E: DEVTYPE=usb_device
E: DRIVER=usb
E: ID_BUS=usb
E: ID_MODEL_ENC=\x20
E: ID_MODEL_FROM_DATABASE=WMRS200 weather station
E: ID_MODEL_ID=ca01
E: ID_REVISION=0302
E: ID_SERIAL=0fde_
E: ID_USB_INTERFACES=:030000:
E: ID_VENDOR=0fde
E: ID_VENDOR_ENC=0fde
E: ID_VENDOR_FROM_DATABASE=Oregon Scientific
E: ID_VENDOR_ID=0fde
E: MAJOR=189
E: MINOR=11
E: PRODUCT=fde/ca01/302
E: SUBSYSTEM=__usb__
E: TYPE=0/0/0
E: USEC_INITIALIZED=5929384
</pre>

I will note the vendor ID and the product ID.
Funny stuff is that it presents itself as a WMRS200 and the model I have is a RMS300, but never mind.

Let's create the udev rule file using the previous informations about the idVendor and the idProduct and create a special file `/dev/weather-station`.
This will make the code more easy as I will be able to hard code the name, and leave the boring task of finding the device aside.

```shell
cat << EOF > /etc/udev/rules.d/50-weather-station.rules
# Weather Station
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fde", ATTRS{idProduct}=="ca01", MODE="0660", GROUP="plugdev", SYMLINK+="weather-station"
EOF
```

Once done, I can restart udev with `sudo /etc/init.d/udev restart` or reload and trigger the rules with `udevadm`

IF something goes wrong, you can check the logs by turning the log level to info, reload the rules and look into the syslog file
```
# udevadm control -l info
# udevadm control -R
# # grep -i udev /var/log/syslog 
# 
```

```
# ls -lrt /dev/weather-station                                                                                                               
lrwxrwxrwx 1 root root 15 Aug 29 21:32 /dev/weather-station -> bus/usb/001/007
# ls -lrt /dev/bus/usb/001/007                                                                                                   
crw-rw-r-- 1 root plugdev 189, 6 Aug 29 21:32 /dev/bus/usb/001/007
```

So far so good...

# Accessing the data

## The libusb
Linux has a low level library "libusb" that make the development of modules easy: [libusb-1.0](http://www.libusb.org/wiki/libusb-1.0).
On my <a type=amzn>rpi</a>, I can install the development version with a simple `sudo apt-get install libusb-1.0-0-dev`.

## Using GO: The `gousb` library

A binding for the libusb is available through the [gousb](https://github.com/truveris/gousb)

There is also a __lsusb__ version that is available as an example.
Let's grab it with a simple
`go get -v github.com/kylelemons/gousb/lsusb`

and test it 

    # ~GOPATH/bin/lsusb
    
    001.004 0fde:ca01 WMRS200 weather station (Oregon Scientific)
      Protocol: (Defined at Interface level)
      Config 01:
        --------------
        Interface 00 Setup 00
          Human Interface Device (No Subclass) None
          Endpoint 1 IN  interrupt - unsynchronized data [8 0]
        --------------
    001.003 0424:ec00 SMSC9512/9514 Fast Ethernet Adapter (Standard Microsystems Corp.)
      Protocol: Vendor Specific Class
      Config 01:
        --------------
        Interface 00 Setup 00
          Vendor Specific Class
          Endpoint 1 IN  bulk - unsynchronized data [512 0]
          Endpoint 2 OUT bulk - unsynchronized data [512 0]
          Endpoint 3 IN  interrupt - unsynchronized data [16 0]
        --------------
    001.002 0424:9514 SMC9514 Hub (Standard Microsystems Corp.)
      Protocol: Hub (Unused) TT per port
      Config 01:
        --------------
        Interface 00 Setup 00
          Hub (Unused) Single TT
          Endpoint 1 IN  interrupt - unsynchronized data [1 0]
        Interface 00 Setup 01
          Hub (Unused) TT per port
          Endpoint 1 IN  interrupt - unsynchronized data [1 0]
        --------------
    001.001 1d6b:0002 2.0 root hub (Linux Foundation)
      Protocol: Hub (Unused) Single TT
      Config 01:
        --------------
        Interface 00 Setup 00
          Hub (Unused) Full speed (or root) hub
          Endpoint 1 IN  interrupt - unsynchronized data [4 0]
      --------------

## Rawread

I want to read the raw data from the device.
The gousb package comes along with an example named "rawread". I'm using it:

<pre>
# rawread git:(master) # go run main.go -device "0fde:ca01"
2016/08/30 14:00:01 Scanning for device "0fde:ca01"...
  Protocol: (Defined at Interface level)
  Config 01:
    --------------
    Interface 00 Setup 00
      Human Interface Device (No Subclass) None
      Endpoint 1 IN  interrupt - unsynchronized data [8 0]
    --------------
2016/08/30 14:00:01 Connecting to endpoint...
2016/08/30 14:00:01 - &usb.Descriptor{Bus:0x1, Address:0x4, Spec:0x110, Device:0x302, Vendor:0xfde, Product:0xca01, Class:0x0, SubClass:0x0, Protocol:0x0, Configs:[]usb.ConfigInfo{usb.ConfigInfo{Config:0x1, Attributes:0x80, MaxPower:0x32, Interfaces:[]usb.InterfaceInfo{usb.InterfaceInfo{Number:0x0, Setups:[]usb.InterfaceSetup{usb.InterfaceSetup{Number:0x0, Alternate:0x0, IfClass:0x3, IfSubClass:0x0, IfProtocol:0x0, Endpoints:[]usb.EndpointInfo{usb.EndpointInfo{Address:0x81, Attributes:0x3, MaxPacketSize:0x8, MaxIsoPacket:0x0, PollInterval:0xa, RefreshRate:0x0, SynchAddress:0x0}}}}}}}}}
2016/08/30 14:00:01 open: usb: claim: libusb: device or resource busy [code -6]
</pre>

After digging into the documentation and forums about the libusb, it looks like the device is locked by a generic kernel driver.
So I need to detach it first.

The API call used to detach a kernel driver is `libusb_detach_kernel_driver`. Sadly it has not be bound to the golang's library.
Indeed [Joseph Poirier](https://github.com/jpoirier) maintain an active fork from the gousb library and he does implement the call.
It's a private method that is called implicitly by another call, so no need to modify the code from rawread to use it.

I've switched to his version:

<pre>
# go get github.com/jpoirier/gousb/rawread
./main -device "0fde:ca01"
2016/08/30 14:12:28 Scanning for device "0fde:ca01"...
  Protocol: (Defined at Interface level)
  Config 01:
    --------------
    Interface 00 Setup 00
      Human Interface Device (No Subclass) None
      Endpoint 1 IN  interrupt - unsynchronized data [8 0]
    --------------
2016/08/30 14:12:28 Connecting to endpoint...
2016/08/30 14:12:28 - &usb.Descriptor{Bus:0x1, Address:0x4, Spec:0x110, Device:0x302, Vendor:0xfde, Product:0xca01, Class:0x0, SubClass:0x0, Protocol:0x0, Configs:[]usb.ConfigInfo{usb.ConfigInfo{Config:0x1, Attributes:0x80, MaxPower:0x32, Interfaces:[]usb.InterfaceInfo{usb.InterfaceInfo{Number:0x0, Setups:[]usb.InterfaceSetup{usb.InterfaceSetup{Number:0x0, Alternate:0x0, IfClass:0x3, IfSubClass:0x0, IfProtocol:0x0, Endpoints:[]usb.EndpointInfo{usb.EndpointInfo{Address:0x81, Attributes:0x3, MaxPacketSize:0x8, MaxIsoPacket:0x0, PollInterval:0xa, RefreshRate:0x0, SynchAddress:0x0}}}}}}}}}
</pre>

Nothing more because the code ends by 

```go
  ep, err := dev.OpenEndpoint(uint8(*config), uint8(*iface), uint8(*setup), uint8(*endpoint)|uint8(usb.ENDPOINT_DIR_IN))
  if err != nil {
      log.Fatalf("open: %s", err)
  }
  _ = ep 
```

Cool... Now let's add some code to read from the endpoint (which is an interface and that implements a Read method as described [here](https://godoc.org/github.com/jpoirier/gousb/usb#Endpoint))

```go
  b := make([]byte, 16)
  _, err = ep.Read(b)
  if err != nil {
      log.Fatalf("read: %s", err)
  }
  log.Printf("%v", b)
  _ = ep 
```

And run the code:

<pre>
go run main.go -device "0fde:ca01"
2016/08/30 14:25:58 Scanning for device "0fde:ca01"...
  Protocol: (Defined at Interface level)
    Config 01:
    --------------
    Interface 00 Setup 00
      Human Interface Device (No Subclass) None
      Endpoint 1 IN  interrupt - unsynchronized data [8 0]
    --------------
2016/08/30 14:25:58 Connecting to endpoint...
2016/08/30 14:25:58 - &usb.Descriptor{Bus:0x1, Address:0x4, Spec:0x110, Device:0x302, Vendor:0xfde, Product:0xca01, Class:0x0, SubClass:0x0, Protocol:0x0, Configs:[]usb.ConfigInfo{usb.ConfigInfo{Config:0x1, Attributes:0x80, MaxPower:0x32, Interfaces:[]usb.InterfaceInfo{usb.InterfaceInfo{Number:0x0, Setups:[]usb.InterfaceSetup{usb.InterfaceSetup{Number:0x0, Alternate:0x0, IfClass:0x3, IfSubClass:0x0, IfProtocol:0x0, Endpoints:[]usb.EndpointInfo{usb.EndpointInfo{Address:0x81, Attributes:0x3, MaxPacketSize:0x8, MaxIsoPacket:0x0, PollInterval:0xa, RefreshRate:0x0, SynchAddress:0x0}}}}}}}}}
2016/08/30 14:25:59 [7 0 48 0 48 53 1 255 7 255 0 66 129 239 0 32]
</pre>

OK! Here are the data, now what I need to figure out, is how to interpret them!

## Decoding the Protocol

Internet is a great tool: I've found a description of the protocol [here](http://www.bashewa.com/wmr200-protocol.php)

I've read that it was mandatory to send a heartbeat sequence every 30 seconds.
I will implement the heartbeat later. For now I will send it initially to request data from the station:

```go
// This is a hearbeat request (9 bytes array)
h := []byte{0x00, 0x01, 0xd0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
log.Println("Sending heartbeat")
i, err := ep.Write(h)
if err != nil {
    log.Fatal("Cannot send heartbeat", err)
}
log.Println("%v bytes sent",i)
```

And then read the stream back. Every data payload is separate from the others by a 0xffff sequence.

### Testing the sequence initialization request

<pre>
 go run main.go -device "0fde:ca01"
2016/08/30 20:02:19 Scanning for device "0fde:ca01"...
Protocol: (Defined at Interface level)
  Config 01:
  --------------
  Interface 00 Setup 00
    Human Interface Device (No Subclass) None
    Endpoint 1 IN  interrupt - unsynchronized data [8 0]
  --------------
  2016/08/30 20:02:19 Connecting to endpoint...
2016/08/30 20:02:19 Sending heartbeat
2016/08/30 20:02:19 heartbeat failed: usb: write: not an OUT endpoint
</pre>

__What did² I do wrong?__ 
<center>
![XKCD]( http://imgs.xkcd.com/comics/debugging.png)
</center>

Easy, I didn't RTFM...
Actually, I didn't read the specification of the USB.

As described [here](http://events.linuxfoundation.org/sites/events/files/slides/elc_2014_usb_0.pdf) the USB is a __host-controlled__ bus which means that:

* Nothing on the bus happens without the host first initiating it.
* Devices cannot initiate a transaction.
* The USB is a Polled Bus
* The Host polls each device, requesting data or sending data

The possibles transaction are:

* IN : Device to Host
* OUT: Host to Device

On top of that, a device may handle 1 to N configuration which handles 1 to N endpoints which may be considered IN or OUT.

My weather station has only one endpoint which is IN.
Therefore I will not be able to send information to the station from the host. What I will be able to send is a IN token to get data on the bus.

<pre>
# lsusb -v
...
Endpoint Descriptor:
  bLength                 7
  bDescriptorType         5
  bEndpointAddress     0x81  EP 1 IN
  bmAttributes            3
    Transfer Type            Interrupt
    Synch Type               None
    Usage Type               Data
  wMaxPacketSize     0x0008  1x 8 bytes
  bInterval              10
</pre>

__Note__ I also see that the endpoint is an interrupt

# To be continued...

This blog post is quiet long, and I haven't finished my research yet. Indeed I think that there is enough information for the post to go live.
I will post a part II as soon as I will have time to continue my experiments with the USB device and the <a type=amzn>rpi</a>.


