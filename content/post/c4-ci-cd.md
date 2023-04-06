---
title: "C4 at scale: make it efficient, then make it ambient"
date: 2023-04-05T09:10:52+02:00
lastmod: 2023-04-05T09:10:52+02:00
draft: false
keywords: []
description: ""
tags: []
categories: []
author: "Olivier Wulveryck"

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: false
toc: true
autoCollapseToc: false
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: false

---

## Introduction

C4 (Context, Containers, Components, and Code) is a powerful and flexible approach to visualizing software architectures. 

In a [previous article](http://localhost:1313/2022/03/10/ccccue-generating-c4-diagrams-with-cue.html) I explained how I created a DSL to describe my C4 as data thanks to the [CUE](https://cuelang.org) language.

In this article, we'll explore how I effectively use the C4 tool based on CUE to organize my code, handle tags, manage versioning, and integrate with Continuous Integration and Continuous Deployment (CI/CD) systems to generate and update diagrams. 
We'll also discuss the benefits of generating diagrams on push and on demand, comparing it with other tools like Mermaid.

On [Wardley's evolution axis](https://blog.owulveryck.info/2023/04/04/should-you-read-this-article-about-wardley-maps.html#the-evolution-axis-horizontal) the initial article was about _making it work_, this one is about _making it efficient_, and end on some considerations to _make c4 ambient_ within a business.

![](/assets/images/evolution.png)

## Part I - Working in a single context: toward efficiency

Most of the time, **I** use the C4 model to describe a single application within a specific domain. 
As a result, **my** C1 (Context) phase is usually light and mainly focuses on the interactions between the system and end-users and external dependencies. 
The C2 (Container) phase provides more information, and **I** often use tags and sprites to combine C2 and Deployment diagrams.
This is generally acceptable because **my** diagrams are intended to support discussions and explanations rather than being self-sufficient (they are accompanied by a narrative).

_side note_: you will understand why the usage of bold fonts later in the last part of this article.

Since I use the C1/C2 to support a narrative, I require different views and components depending on the message I want to convey. To organize these elements, I use a structured approach similar to coding.

_Ubiquitous language_: in this article I use _code_ to represent any input from the end-user to express the diagram. It can be a [Structurizr code](https://structurizr.com/), plantUML or [CUE Data](https://github.com/owulveryck/cue4puml4c4)

### Code Organization and Packages
To make the most of the C4 model, it's essential to maintain a clear and logical structure in your codebase. Let's see a practical example of what this means.

#### Organizing code and leveraging the lattice

Let's start from this illustration from the original [c4model.com](https://c4model.com) website:

> The C4 model was created as a way to help software development teams describe and communicate software architecture, both during up-front design sessions and when retrospectively documenting an existing codebase. 
It's a way to create maps of your code, at various levels of detail, in the same way you would use something like Google Maps to zoom in and out of an area you are interested in. 


In the context of diagram-as-code, the application of this principle can be implemented in two ways: 
1. by reflecting the zoom level within the code, or 
2. by adjusting the data at runtime to zoom on demand. 

I won't discuss the pros and cons of each method (I could write a complete article about it). 

Instead, I will focus on the first option, which is the one I am using with my CUE implementation. 
In fact, the goal of CUE is to provide a complete description within the code, without leaving any room for data mutation at runtime, and I appreciate this philosophy.

_Note_, the remaining of this section is mainly about how to use the CUE4Puml4C4. Feel free so skip to the [next section](#about-packages) if you are not interested in this geeky part.

By definition, a System in my DSL is something that have a description, a label and **optionnaly** containers and relations.

I first define a System without the containers and relations, and instead of repeating myself, I then define another element which is the junction of the first element **and**
the description provided.

Let's illustrate this whith an example:

{{< highlight cue >}}
mySystemBoundaryC1: c4.#System & {
	id:          "myapp"
	label:       "My Application"
	description: "My Application is fancy"
}
mySystemBoundaryC2: c4.#System & mySystemBoundaryC1 & {
	isBoundary: true
	containers: [service1, db1]
	relations: [
		{source: service1, dest: db1, description: ""},
	]
}
service1: c4.#Container & {
	id:          "service1"
	label:       "Service 1"
	description: "Service 1 is used to apply business rules",
}
{{< / highlight >}}

We have two data sets: 

- `mySystemBoundaryC1`, which represents an empty system, and 
- `mySystemBoundaryC2`, which combines `mySystemBoundaryC1` with a description of the containers and their relationships. 

It's important to note that:

- `mySystemBoundaryC2` is not a placeholder; it is a structure that always holds the containers
- `mySystemBoundaryC1` will never have any associated containers. 

However, any changes made to `mySystemBoundaryC1` will be reflected in `mySystemBoundaryC2`. 
This mechanism means that it's not possible to override the label and descriptions of `mySystemBoundaryC2` from `mySystemBoundaryC1`. 
This is because the label cannot be both `mySystemBoundaryC1`.label ('My Application') and 'another definition'

The DSL allows Systems to contains subsystems; actually, I often use this capability to group various elements that serves a similar purpose such as a database and its associated micro-service.

A good practice I have is to segregate the definition of the various services into packages to avoid naming collision and easy the maintenance of the global picture.
Let's now detail the rationale of the packaging system.

#### About packages

In this part, we'll discuss how I leverage the notion of packages to group related code elements, making it easier to understand the overall architecture, reduce complexity, and facilitate collaboration.

_Ubiquitous language_: in this article, a package is a collection of related source files that are organized together within a directory. **Packages** provide a way to **encapsulate** code, which promotes **reusability** and **modularization**.

In CUE, packages are used to group related definitions, values, and constraints. CUE packages are also organized in directories, and each package has a unique name.

I inherit my directory and package layout from my experience with Go. 
At the root of the project, there is the `main` package, that is the entrypoint that displays the global picture of the domain.

Then each subsystem is declared in its own package. I can import the C1 representation in the main model in a global C1 object, or the C2 in a global C2 object.

The resulting of this import is that the `main` package holds the definitions that I people entering the project needs to see in a glimpse. Later in this article we will see that this sentence is not only valuable for humans reading the README of the repo, but also for
the other systems depending on our context.

### Handling Tags

Tagging is a vital aspect of C4, as it helps identify and differentiate elements within diagrams. 
We'll discuss the importance of a well-defined tagging system for maintaining consistency and clarity throughout the architecture. 

I use tags for various reasons.
One of those reasons is to represent the temporal evolution of components. 
In simple terms, these tags act as markers to highlight any changes or transformations that will happen to the components over time. 
For example, you can use tags to indicate when a component is slated for removal or when it will undergo a transformation. 
By using these tags, you can effectively communicate the future state of the system or its components, providing a clearer understanding of the system's life cycle and evolution.

![](/assets/images/c4_tags.png)

In this case, tags are specifics and are relative to the context of the current diagram. Therefore they can be declared in the same codebase / packages.

But some other elements needs tags that are standardized, therefore, a good practice is to declare them in their own package and make it available to the wild.
This is an example with [Team Topologies](https://github.com/owulveryck/cue4puml4c4/tree/main/tags/teamtopologies).

![](/assets/images/c4_tt.png)

### Applying the DevOps practices

Automation is key to maintaining up-to-date architecture diagrams. 

I see two key elements to take into consideration regarding this topic:

* CI/CD Integration for Diagram Generation
* Generation on Push vs. Generation on Demand

Here, we'll discuss how I've integrated C4 into my CI/CD pipeline, automating the process of generating and updating diagrams whenever changes are made to the codebase.

When it comes to generating architecture diagrams, there are two main approaches: generation on push and generation on demand (via things like Mermaid). 

I had a discussion with a partner recently, and I'd like to take the opportunity of this post to open the discussion.



TODO


So far, I have described some patterns I use to manage my diagrams as code. Although not strictly implemented with code, these patterns draw inspiration from software development practices.

Let's explore how to take the management of diagrams one step further and treat them as an engineering asset. As Titus Winter once wrote

> Software engineering is what happens to programming when you add time and other programmers.

By treating diagrams as engineering assets, we can optimize their management and contribute to provide the value of the C4 model at scale.

## Part II - C4 at scale: making the model ambient

In the first part of this article, I used bold fonts to spot something:

> Most of the time, **I** use the C4 model to describe a single application within a specific domain. 
As a result, **my** C1 (Context) phase is usually light and mainly focuses on the interactions between the system and end-users and external dependencies. 
The C2 (Container) phase provides more information, and **I** often use tags and sprites to combine C2 and Deployment diagrams.
This is generally acceptable because **my** diagrams are intended to support discussions and explanations rather than being self-sufficient (they are accompanied by a narrative).

What's important to note is that, in the context of making things more efficient, the focus is primarily on my own usage of the DSL. 
You make things efficient for your own purpose.
An empiric consequence is that, while the DSL is a helpful tool for standardizing a model, it doesn't necessarily ensure that every user will create diagrams using the same conventions.

### A governance to rule them all

Running the **C4 at scale**	 is not only about getting a **global adoption** accros the various projects and products, but also about fostering a **culture of shared understanding and quality** in software architecture.
By focusing on support, collaboration, and continuous improvement, governance can help organizations reach these goals.

#### About the notion of governance

_Ubiquitous language_: Governance, in the context of using the C4 model at scale, refers to a set of guidelines, best practices, and policies that help organizations effectively adopt and apply the C4 model across various teams and projects. 
The primary goal of governance in this context is to enable consistency, collaboration, and maintainability while promoting a shared understanding of software architecture.

**Governance** should **not** be viewed as a **controlling or punitive** measure.
Instead, it should be seen as a **supportive framework** that **empowers teams** to work **efficiently** and **effectively** with the C4 model. 
It helps create a common language and approach for architecture discussions, increasing collaboration and enhancing the overall quality of software systems.

Some key aspects of C4 model governance might include:

1. Defining and communicating standard conventions, such as naming, diagram layout, and notations, to ensure consistency and ease of understanding across teams and projects.
2. Establishing best practices for creating, maintaining, and sharing C4 diagrams, including version control, documentation, and integration with other tools.
3. Providing training and resources to help team members understand and effectively use the C4 model in their day-to-day work.
4. Encouraging and facilitating regular reviews of software architecture, promoting a culture of continuous improvement and knowledge sharing.
5. Integrating the C4 model into existing processes, such as software development life cycles, change management, and project management, to ensure that architecture considerations are an integral part of the overall development process.
6. Monitoring and evaluating the effectiveness of the C4 model's adoption within the organization, and adapting the governance framework as needed to address any challenges or opportunities for improvement.

As it is a big area and I want to stay pragmatic in this article, let's focus on the first elements and how they can be a natural evolution of what we've seen so far.

Remember:

![](/assets/images/evolution.png)

### About Conventions

TODO

### A model of components and tags to enforce best practices

TODO 

## Conclusion and next steps
The C4 model offers a powerful way to visualize software architecture, but harnessing its full potential requires careful organization, tagging, versioning, and automation. 
By implementing these strategies and understanding the benefits of generating diagrams on push or on demand, you'll be well-equipped to create clear,
consistent, and up-to-date representations of your software systems that facilitate collaboration and understanding among stakeholders.

### What's next? 

#### Bi-temporality?


