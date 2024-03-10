---
title: "After the BYOD, BYON (briging your own network): a journey from Home to the World"
date: 2024-03-09T12:15:33+01:00
lastmod: 2024-03-09T12:15:33+01:00
draft: true
keywords: []
summary: Discover how I transformed my reMarkable tablet into a portable whiteboard 📒✨, accessible from anywhere via a secure WireGuard VPN (tailscale) and cloud-based reverse proxy setup.
  
  From the comforts of WFH to the dynamic world of mobility, learn the tech behind the solution.
tags: []
categories: []
author: ""

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S.
comment: false
toc: true
autoCollapseToc: false
# You can also define another contentCopyright.
contentCopyright: false
reward: false
mathjax: false

# Uncomment to add to the homepage's dropdown menu; weight = order of article
# menu:
#   main:
#     parent: "docs"
#     weight: 1
---

In the height of the work-from-home era, I developed goMarkableStream, a tool designed to seamlessly stream content from my reMarkable tablet during video calls.
The tool moved from a proof of concept to become a part of my daily toolbox. The idea behind the tool is:

A service is running on the reMarkable device and captures the image. 
It exposes a service that serves the image over HTTP(s) with a custom implementation. 
Then, a renderer is encoded in the browser in WebGL/JS to display the content of the screen. 
During a video call, I can share a browser tab and therefore share what I am writing with the audience.

Indeed, the solution brought value for remote collaboration and sharing ideas in real-time.

I've shared the details of this journey, highlighting how little by little, it bridged the gap between physical and virtual meetings.

However, as the full-time remote work phase has declined and we've transitioned back to a more mobile lifestyle, the solution that exposes a service on a local network reaches its limits.

In a hybrid scenario, I don't know what network topology I will encounter and their limitations.
Therefore, I faced situations where I tethered my mobile connection to my tablet and was simply unable to stream the content due to limitations.

With hybrid work, and even with a work-from-anywhere situation (home, office, client sites), I need to change the paradigm and to be able to stream over the Internet.
As I cannot simply expose the streaming service hosted on my tablet to the Internet, I need to set up an infrastructure that will secure and route the traffic from the internet to the streaming service.

This article describes the journey to achieve this, from a simple reverse proxy via NGrok to a VPN solution based on WireGuard.

I will first expose the solution based on a reverse proxy powered by NGrok.
The I will explain the limitations that lead me to the solution of accessing the service through a VPN powered by Tailscale.
This part will give hints about the wireguard mechanism, and expose the basic elements of the infrastructure in place to expose the streaming service.

Before the pandemic, we used VPNs to connect to the office from home... Now, I've switched the paradigm to connect to home from the office.
I guess that this is the follow-up of the bring your own device (BYOD) evolution.

## First problem and first solution: NGrok

As [I blogged a couple of months ago](https://blog.owulveryck.info/2023/10/10/rethinking-presentations-beyond-static-slides.html), I am using my tablet as support for presentations.
This is working smoothly on my own network, but I was facing problems when I moved to a site with limitations.
I thought that I could always bring my own laptop with me, but that is not always the case. So, I needed a way to expose the streaming service to the Internet and give the address to the people in charge of presenting the content.

The first and easy step I found as a solution was to embed the NGrok service in my tool.
Actually, NGrok's promise is:

> Connect to external networks in a consistent, secure, and repeatable manner without requiring any changes to network configurations.
- Bring Your Own Cloud (BYOC) Connectivity
- IoT Connectivity

The implementation was fairly easy to be embedded in the tool. 
Actually, as there is a Go SDK for NGrok and my tool is written in Go, I simply needed to set up a custom listener to my service and the framework did the rest.

Basically, NGrok creates a custom `Listener`, and all I need to do is to switch the basic listener of the HTTP service to use this listener instead.

Here is a helper function to initialize the listener based on a configuration structure:
```go 
func setupListener(ctx context.Context, c *configuration) (net.Listener, error) {
        switch c.BindAddr {
        case "ngrok":
                l, err := ngrok.Listen(ctx,
                        config.HTTPEndpoint(),
                        ngrok.WithAuthtokenFromEnv(),
                )
                c.BindAddr = l.Addr().String()
                c.TLS = false
                return l, err
        default:
                return net.Listen("tcp", c.BindAddr)
        }
}
```

And here is its usage in the main loop (`handler` had been configured before):

```go 
l, err := setupListener(context.Background(), &c)
// ...
log.Fatal(http.Serve(l, handler))
```

When I launch the tool, with the correct environment variables, it connects to the NGrok service and displays the external URL to connect to.
And voilà: it works!

However, there are constraints and limitations:

- First of all, with the free version of NGrok, the network is limited. I will not be able to use my tool the entire month, but I could live with it.
- The second problem is that I cannot configure the DNS of the endpoint on the free version. And every time it restarts, the URL of the endpoint changes. This is annoying.

All of these problems would have been fixed by paying for the NGrok service, but it is far too expensive for my needs and indeed, would not have solved the last problem:

The solution does not handle roaming (changing networks) and long pauses (when the tablet is sleeping for a long time) well. That made the solution unreliable.

So I looked for another solution.

### Solution: setup a VPN

### Why WireGuard and Tailscale?

In traditional VPN protocols, such as IPsec or OpenVPN, the VPN connection is typically tied to the IP address of the connecting device.
If the device's IP address changes (for example, when moving from one network to another), the VPN connection would usually drop, requiring a re-establishment of the connection under the new IP address.
This process can introduce delays and interruptions in connectivity.

#### WireGuard's Approach

WireGuard, on the other hand, takes a different approach that inherently supports seamless roaming:

- **Connection Identification:** WireGuard identifies connections not by the source or destination IP addresses but through the cryptographic identity of the peers (i.e., their public keys).
This means that as long as the cryptographic identity remains the same, WireGuard does not care if the actual IP address of a device changes.
- **Session Persistence:** When a WireGuard client moves to a different network and obtains a new IP address, it simply sends authenticated packets from its new IP to the WireGuard server (or peer).
The server recognizes the client by its public key and continues the session without interruption.
The server then automatically updates its internal routing table with the client’s new IP address, maintaining the encrypted tunnel without needing to re-establish the connection.
- **Rapid Response:** This mechanism allows for almost instantaneous switching between networks.
Users typically do not notice any disruption in their VPN connection as they move across different networks, making WireGuard particularly suited for mobile devices that frequently change network environments.

#### Benefits of Roaming

The roaming feature offers several practical benefits:

- **Uninterrupted Connectivity:** Users experience seamless transitions between networks without losing VPN protection or experiencing drops in their connections.
- **Improved User Experience:** The seamless nature of roaming with WireGuard enhances the overall user experience, as there are no manual reconnections or VPN downtime during network changes.
- **Efficiency for Mobile Devices:** For devices that switch networks often, WireGuard's efficient handling of IP changes ensures minimal battery consumption and resource usage compared to protocols that require re-establishing connections.

In summary, WireGuard's roaming capability is a significant advancement in VPN technology, providing users with stable, continuous connectivity across different networks and enhancing the usability of VPNs on mobile devices and in dynamic network environments."

### Challenges and Solutions in Integration

TODO: explain how I tried to implement tailscalem but the kernel of the remarkable does not support tunneling, therefore I must choose a solution that works in userspace.

### Introduction to the tsnet Library

> tsnet is a library that lets you embed Tailscale inside of a Go program.
This uses a userspace TCP/IP networking stack and makes direct connections to your nodes over your tailnet just like any other machine on your tailnet would.
When combined with other features of Tailscale, this lets you create new and interesting ways to use computers that you would have never thought about before.


### Security Considerations

- **VPN Security:** Outline the security measures inherent in using WireGuard and Tailscale, such as end-to-end encryption.
- **Additional Security Measures:** Describe any extra steps taken to ensure data privacy and protection, like network segmentation or access controls.

### Configuration and Setup

- **Setting Up Tailscale and WireGuard:** Provide a step-by-step guide on configuring Tailscale on the reMarkable and setting up WireGuard VPN.
- **Cloud Machine as Gateway:** Explain how you configured the cloud machine to act as a gateway and the role of the reverse proxy in this setup.

### Reverse Proxy and TLS Management

- **Role of the Reverse Proxy:** Detail how the reverse proxy manages traffic and facilitates secure access to your network.
- **Automating TLS with Let's Encrypt:** Describe the process of obtaining and renewing TLS certificates automatically to keep connections secure.

### Practical Benefits and Use Cases

- **Enhanced Accessibility:** Share how accessing the reMarkable's streaming service remotely has improved your workflow or use cases.
- **Examples of Benefits:** Provide specific examples or scenarios where this setup has been particularly beneficial.

### Performance and Reliability

- **Impact on Device Performance:** Discuss any observed effects on the reMarkable's performance or battery life.
- **Reliability of the Setup:** Share insights on the setup's reliability, including any downtime or connectivity issues experienced.

### Future Improvements and Expansions

- **Planned Enhancements:** Talk about potential future improvements or expansions, such as adding more devices to the VPN network.
- **Security Enhancements:** Consider future security measures that could further protect your setup.

### Conclusion

- **Recap:** Summarize the key points discussed in the article and the benefits of your setup.
- **Encouragement:** Encourage readers to explore the possibilities of VPNs and secure networking for their own projects.





