---
title: "The Future of Data Management: an enabler to AI devlopment? A basic illustration with RAG, Open Standards and Data Contracts"
date: 2024-06-12T12:15:33+01:00
lastmod: 2024-06-12T12:15:33+01:00
# images: [/assets/rag/illustration.png]
draft: false
keywords: []
summary: 
tags: []
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
# You can also define another contentCopyright.
contentCopyright: false
reward: false
mathjax: false
---

## Context

In my upcoming article, I want to cover the following elements:

- **Brief Explanation of a Data Contract:** Start with a concise description of what a data contract is.
- **Importance of Open Standards and Introduction of Bitol:** Highlight the significance of open standards and introduce Bitol as a solution.
- **Validation of Contracts with CUE (Cuelang):** Provide a side explanation on how to validate a data contract using CUE (Cuelang).

For the concrete example, I'll draw inspiration from a previous article about sending queries to a book. Here, the chunks of the book and their corresponding embeddings are the data. I'll define a data contract for this data. 

I will separate them into two domains: the chunks belong to the domain of knowledge, while the embeddings will belong to another domain, perhaps the domain of semantic representation or the domain of contextual insights.

The data contract will provide confidence in the chunking process. Each chunk, considered as a data quantum or data-as-a-product, is designed to be a self-sufficient unit relevant to a particular domain. The embedding will have its own contract, which will include the version of the data-product of the chunk and the algorithm used for embedding.

This setup allows the creation of a consumer of this data, enhancing the quality of my tool that queries the data. Furthermore, by explaining various questions that my tool can answer, I can gradually build a data-mesh. 

---
**Semantic representation** is a key concept in the field of artificial intelligence and natural language processing. It aims to represent the meaning or semantics of information, such as sentences, paragraphs, or entire documents, in a way that is understandable by machines. Here is a detailed explanation:

### What is Semantic Representation?

Semantic representation involves transforming raw textual data into formal structures that capture their meaning. These structures allow machines to understand and analyze textual content in a more intelligent and contextual manner. It goes beyond the words themselves to grasp the concepts, relationships, and context they convey.

### Why is it Important?

1. **Natural Language Understanding:** Semantic representation enables machines to understand human language more accurately and nuanced. This is crucial for applications like chatbots, virtual assistants, and machine translation.

2. **Information Retrieval and Search:** By using semantic representations, search systems can understand the intent behind user queries and provide more relevant results.

3. **Sentiment and Opinion Analysis:** Understanding the nuances and context of texts allows for analyzing the sentiments and opinions expressed in texts, such as product reviews or social media comments.

### How is it Achieved?

1. **Word Embeddings:** Models like Word2Vec or GloVe represent words as vectors in a multidimensional space, where words with similar meanings are close to each other.

2. **Contextual Models:** Advanced models like BERT and GPT capture not only the relationships between words but also their context within entire sentences. This allows for a richer and more contextual understanding of texts.

3. **Knowledge Graphs:** Knowledge graphs structure information as nodes (concepts) and edges (relationships between these concepts), facilitating semantic navigation and querying of data.

### Concrete Example

Consider the sentence: "The cat sleeps on the sofa."

- **Without Semantic Representation:** A simple textual processing might just see a series of words without understanding the relationships between them.

- **With Semantic Representation:** The system understands that "cat" is an animal, "sleeps" is an action related to "cat," and "sofa" is a location. Additionally, it can understand broader relationships, such as pets often sleeping on furniture.

### Application in Data Contracts

In the context of creating data contracts, semantic representation is crucial to ensure that data is correctly interpreted and used. For example, in a data contract for book chunks and their embeddings:

- **Book Chunks (Domain of Knowledge):** Each chunk of text is analyzed for its semantic content to ensure it constitutes a coherent and self-sufficient unit of knowledge.

- **Embeddings (Domain of Semantic Representation):** The embeddings capture the contextual meaning of the text chunks, allowing for effective linking and comparison of different chunks based on their meaning, not just their lexical content.

By using semantic representations, data contracts can specify not only the structure and format of the data but also the necessary semantic and contextual relationships to ensure their utility and coherence.

---

I hope this explanation clarifies the idea of semantic representation and its importance in your project. Feel free to ask if you need more details or additional examples!
---

Let me know if you need further modifications or if you're ready to proceed with writing the article!

Great, let's break this down step-by-step to ensure we cover all the points you want in your article. Here’s a structured approach to the content:

1. **Introduction to Data Contracts**
    - Definition of a data contract
    - Importance and benefits of using data contracts

2. **The Importance of Open Standards and Introduction to Bitol**
    - Explanation of open standards in data management
    - Introduction to Bitol and its role in creating standardized data contracts

3. **Validating Data Contracts with CUE (Cuelang)**
    - Brief introduction to CUE
    - How CUE can be used to validate data contracts

4. **Concrete Example: Applying Data Contracts to Book Data**
    - Overview of the use case: querying book data
    - Detailed example:
        - Defining the data chunks (domain of knowledge)
        - Defining the embeddings (domain of insights)
        - Specifying data contracts for both chunks and embeddings

5. **Ensuring Quality and Building a Data Mesh**
    - How data contracts ensure quality in data processes
    - Creating a data consumer for improved tool quality
    - Building a data set and evolving towards a data mesh

### Draft Outline and Initial Content

---

### Introduction to Data Contracts

In the world of data management, a **data contract** is a formal agreement between data producers and data consumers that specifies the structure, quality, and business rules of the data being exchanged. Data contracts serve as a safeguard, ensuring that data meets specific criteria before being consumed, thereby increasing the reliability and trustworthiness of data-driven processes.

### The Importance of Open Standards and Introduction to Bitol

Open standards play a crucial role in the interoperability and scalability of data systems. They ensure that data can be seamlessly shared and utilized across different platforms and organizations. One such open standard is **Bitol**, which provides a framework for creating and maintaining data contracts. By adopting Bitol, organizations can establish clear, consistent, and enforceable data contracts that facilitate better data governance and quality assurance.

### Validating Data Contracts with CUE (Cuelang)

**CUE** (Configuration, Unification, and Execution) is a language for defining, generating, and validating data. It is particularly well-suited for data contracts due to its ability to enforce schema and validation rules. Using CUE, data contracts can be specified in a clear and concise manner, and compliance with these contracts can be automatically verified.

```cue
# Example CUE schema for a data contract
dataContract: {
    chunks: [...#Chunk]
    embeddings: [...#Embedding]
}

#Chunk: {
    id: string
    content: string
    domain: "knowledge"
}

#Embedding: {
    id: string
    version: string
    algorithm: string
    chunkId: string
    domain: "insights"
}
```

### Concrete Example: Applying Data Contracts to Book Data

#### Overview of the Use Case: Querying Book Data

To illustrate the concept of data contracts, let's consider a use case where we need to query data from a book. The book is divided into chunks, each representing a self-sufficient part of knowledge. Additionally, we generate embeddings for these chunks, which are used in research and insights.

#### Defining the Data Chunks (Domain of Knowledge)

The chunks are segments of the book, each covering a specific topic. These chunks form the domain of knowledge and need to be accurately defined to ensure they are self-sufficient and meaningful.

```json
{
    "chunks": [
        {
            "id": "chunk1",
            "content": "Introduction to Data Science...",
            "domain": "knowledge"
        },
        {
            "id": "chunk2",
            "content": "Data Cleaning Techniques...",
            "domain": "knowledge"
        }
        // more chunks...
    ]
}
```

#### Defining the Embeddings (Domain of Insights)

The embeddings are vector representations of the chunks, used for advanced querying and analysis. This forms the domain of insights.

```json
{
    "embeddings": [
        {
            "id": "embedding1",
            "version": "v1",
            "algorithm": "word2vec",
            "chunkId": "chunk1",
            "domain": "insights"
        },
        {
            "id": "embedding2",
            "version": "v1",
            "algorithm": "word2vec",
            "chunkId": "chunk2",
            "domain": "insights"
        }
        // more embeddings...
    ]
}
```

### Ensuring Quality and Building a Data Mesh

Data contracts play a vital role in ensuring the quality of the data chunking process. By defining clear contracts for both the chunks and the embeddings, we can ensure that each chunk is a self-sufficient piece of knowledge and each embedding accurately represents its corresponding chunk.

By creating consumers of these data contracts, we can enhance the quality and reliability of our querying tool. Additionally, by continuously adding more data and refining our contracts, we can gradually build a comprehensive **data mesh**. This decentralized approach to data management promotes greater scalability and adaptability, allowing us to meet the evolving needs of our data consumers.

---

Feel free to expand on each section with more details or examples as needed. Let me know if you need any adjustments or additional content!




```cue
// What's this data contract about?
datasetDomain:       "knowledge"     // Domain
quantumName:         "Wardley Book" // Data product name
userConsumptionMode: "operational"
version:             "1.0.0" // Version (follows semantic versioning)
status:              "test"
uuid:                "53581432-6c55-4ba2-a65f-72344a91553a"

// Lots of information
description: {
  purpose:     "Views built on top of the seller tables."
  limitations: "Data based on seller perspective, no buyer information"
  usage:       "Predict sales over time"
}

// Getting support
productDl:           "wardley-map@myorg.com"

sourcePlatform: "owulveryck's blog"
project: "Sample Data Contract and Rag"          
datasetName:    "wardley_book" 
kind:           "DataContract"
apiVersion:     "v2.2.2" // Standard version (follows semantic versioning, previously known as templateVersion)
type:           "objects"

// Physical access
driver:           "httpfs:parquet"
driverVersion:    "1.0.0"
database:         "https://blog.owulveryck.info/assets/sampledata" // Bucket name

// Dataset, schema and quality
dataset: [{
  table:          "wardleybook.parquet" // the object name
  description:    "The book from simon wardley, chunked"
  authoritativeDefinitions: [{ // NEW in v2.2.0, inspired by the column-level authoritative links
    url:  "https://blog.owulveryck.info/2024/06/12/the-future-of-data-management-an-enabler-to-ai-devlopment-a-basic-illustration-with-rag-open-standards-and-data-contracts.html"
    type: "explanation"
  }]
  dataGranularity: "Chunking manually to 1000 characters"
  columns: [{
    column:                    "chunk_id"
    isPrimaryKey:              true // NEW in v2.1.0, Optional, default value is false, indicates whether the column is primary key in the table.
    logicalType:               "int"
    physicalType:              "int"
    isNullable:                false
  }, {
    column:                    "content"
    businessName:              "Part of the book"
    logicalType:               "string"
    description:               "A chunk of the book in markdown"
  }]
}]
```

