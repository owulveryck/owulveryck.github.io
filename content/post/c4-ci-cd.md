---
title: "C4 at scale: make it efficient, then make it ambient"
date: 2023-04-04T09:10:52+02:00
lastmod: 2023-04-04T09:10:52+02:00
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
In this article, we'll explore how I effectively use the C4 tool based on CUE to organize my code, handle tags, manage versioning, and integrate with Continuous Integration and Continuous Deployment (CI/CD) systems to generate and update diagrams. 
We'll also discuss the benefits of generating diagrams on push and on demand, comparing it with other tools like Mermaid.

![](/assets/images/evolution.png)

## Working in a single context: toward efficiency



### Code Organization and Packages
To make the most of the C4 model, it's essential to maintain a clear and logical structure in your codebase. 
We'll discuss how I leverage the notion of packages to group related code elements, making it easier to understand the overall architecture, reduce complexity, and facilitate collaboration.

Then, I will examine how the reusability of components, enabled by the packaging system, allows for the scalable use of the C4 model and brings additional value at the enterprise level.


### About packages

A package is a collection of related source files that are organized together within a directory. Packages provide a way to encapsulate code, which promotes reusability and modularization.

In CUE, packages are used to group related definitions, values, and constraints. CUE packages are also organized in directories, and each package has a unique name.


### Handling Tags
Tagging is a vital aspect of C4, as it helps identify and differentiate elements within diagrams. 
We'll discuss the importance of a well-defined tagging system for maintaining consistency and clarity throughout the architecture. 
We'll also share the tagging conventions and best practices I use to ensure that my diagrams are both informative and easy to navigate.

### Diagrams as code

#### Versioning
Software architectures evolve over time, and it's important to track these changes effectively. 
In this section, we'll delve into how I manage versioning in my C4 diagrams, ensuring that each update is well-documented and easily accessible to stakeholders.

#### CI/CD Integration for Diagram Generation
Automation is key to maintaining up-to-date architecture diagrams. 
Here, we'll discuss how I've integrated C4 into my CI/CD pipeline, automating the process of generating and updating diagrams whenever changes are made to the codebase.

#### Generation on Push vs. Generation on Demand:
When it comes to generating architecture diagrams, there are two main approaches: generation on push and generation on demand. 
We'll compare these methods and discuss how they differ from tools like Mermaid. 
We'll also share the pros and cons of each approach and explain why I chose the method I use in my projects.



## C4 at scale: making the model ambient

### A governance to rule them all

### About the notion of governance

Governance, in the context of using the C4 model at scale, refers to a set of guidelines, best practices, and policies that help organizations effectively adopt and apply the C4 model across various teams and projects. The primary goal of governance in this context is to enable consistency, collaboration, and maintainability while promoting a shared understanding of software architecture.

Governance should not be viewed as a controlling or punitive measure. Instead, it should be seen as a supportive framework that empowers teams to work efficiently and effectively with the C4 model. It helps create a common language and approach for architecture discussions, fostering collaboration and enhancing the overall quality of software systems.

Some key aspects of C4 model governance might include:

1. Defining and communicating standard conventions, such as naming, diagram layout, and notations, to ensure consistency and ease of understanding across teams and projects.
2. Establishing best practices for creating, maintaining, and sharing C4 diagrams, including version control, documentation, and integration with other tools.
3. Providing training and resources to help team members understand and effectively use the C4 model in their day-to-day work.
4. Encouraging and facilitating regular reviews of software architecture, promoting a culture of continuous improvement and knowledge sharing.
5. Integrating the C4 model into existing processes, such as software development life cycles, change management, and project management, to ensure that architecture considerations are an integral part of the overall development process.
6. Monitoring and evaluating the effectiveness of the C4 model's adoption within the organization, and adapting the governance framework as needed to address any challenges or opportunities for improvement.

By focusing on support, collaboration, and continuous improvement, governance can help organizations successfully adopt the C4 model at scale, fostering a culture of shared understanding and quality in software architecture.

## Conclusion:
The C4 model offers a powerful way to visualize software architecture, but harnessing its full potential requires careful organization, tagging, versioning, and automation. 
By implementing these strategies and understanding the benefits of generating diagrams on push or on demand, you'll be well-equipped to create clear,
consistent, and up-to-date representations of your software systems that facilitate collaboration and understanding among stakeholders.
