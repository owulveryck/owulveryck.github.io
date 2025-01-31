---
title: "MCP Part III - Application: Custom Server for a Specific Use Case"
date: 2025-01-27T12:15:33+01:00
lastmod: 2025-01-27T12:15:33+01:00
images: [/assets/mcp/diagram.png]
draft: false
keywords: []
summary: This final article in a three-part series explores decoupling tools from the host using the Model Context Protocol (MCP) for flexibility and reusability. The author builds an MCP-based tool in Go to execute SQL queries via DuckDB, enabling seamless chatbot interaction while preserving privacy. The implementation covers JSON-RPC handling, tool encapsulation, and integration. The project validates MCP’s effectiveness, with future plans to replace VertexAI with Ollama and add multi-session support.
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

In the first two parts (see [part I](/2025/01/27/mcp-part-i-core-concepts-past-present-and-future-of-agentic-systems.html) and [part II](/2025/01/27/mcp-part-ii-implementation-custom-host-with-vertexai-and-gemini.html)) of this series, we have:

- described the agentic system and how it could be extended via tool usage;
- explained how a standard like MCP was key to building a virtuous ecosystem;
- built a host "from scratch": we implemented a private chatbot and gave it the ability to run a tool.

In this final part, we will fulfill the initial promise: we will create a tool that is decoupled from the host.

By decoupling the tool from the host, we ensure flexibility, interoperability, and reusability across different AI agents. This aligns with the broader goal of creating a modular, scalable system.

Here target schema of what we are describing in the article:

![The image is a flowchart illustrating the Model Context Protocol (MCP) with interconnected blocks and directional arrows. Blue text and arrows represent one set of components and processes, while red text and arrows represent another set of interactions. At the top left, a block with blue text is labeled 'Role,' which connects to a red-labeled block with the text 'Model.' The red block has arrows branching to the right, leading to additional components labeled in blue and red. Below the 'Role' block, there are additional blue-labeled blocks, connected vertically, indicating hierarchical or sequential relationships. Red arrows interweave with the blue components, suggesting cross-interactions or dependencies between the two flows. The diagram visually represents the structural relationships and interactions within the MCP framework, emphasizing the interplay between blue and red elements.](/assets/mcp/partI/Image5.png)

## Exposing a use case

I learned SQL in the past century and have hardly used it since. So, every time I need to explore a simple dataset, I use code (either a programming language or a combination of `grep`, `awk`, `sed`, and their friends).

However, in the exploration phase, SQL has its advantages, and a tool like DuckDB can provide significant help.

I currently use a chatbot to generate SQL queries that I copy and paste into the DuckDB interface, but this is not optimal.
Therefore, my use case involves creating a tool that will be an MCP server, exposing a single function:

- running a query via DuckDB (the argument is the query).

My goal is to be able to use a chatbot to do something like this:

> - 👤 Please give me the schema of the data contained in the local file by running a sql query (duckdb): 
>  "/tmp/ValeursFoncieres-2024-S1.tsv"
> - 🧠 The file `/tmp/ValeursFoncieres-2024-S1.tsv` has 43 columns. Here are the names and data types of the first few columns:
> 
> | column name          | data type |
> |----------------------|-----------|
> | Identifiant de doc…  | varchar   |
> | Reference document   | varchar   |
> | ...                  | ...       |
> | Nature culture       | varchar   |
> | Nature culture spe…  | varchar   |
> | Surface terrain      | int64     |
> 
> I can provide more details if you specify which columns you are interested in.
> - 👤 ok, this file is in french. 
> My hometown's zip code is 59700, what is the average transaction cost in this city ?
> 
> - 🧠 **Ah, it seems I need to cast the "Valeur fonciere" column to a numeric type before calculating the average. I'll try casting it to DOUBLE.**
> 
> **It seems like the issue is with the format of the "Valeur fonciere" column. The comma is used as a decimal separator, which is not compatible with the CAST function. I need to replace the comma with a dot before casting to DOUBLE.**
> 
> The average transaction cost in the city with zip code 59700 is 441545.00.
> 
> - 👤 quelle est la surface moyenne du terrain dans cette ville ?
> - 🧠 The average land surface area in the city with zip code 59700 is 444.60 square meters.
> - 👤 quelle est le prix moyen au metre carré ?
> - 🧠 The average price per square meter in the city with zip code 59700 is 3532.92.
> - 👤 compare par rapport a la ville dont le code postal est 59000
> - 🧠 The average price per square meter in the city with zip code 59000 is 3978.29.
> Compared to the city with zip code 59700, where the average price per square meter is 3532.92, the city with zip code 59000 has a higher average price per square meter.
> - 👤 et par rapport a la ville dont le code postal est 62100 ?
> - 🧠 The average price per square meter in the city with zip code 62100 is 1810.66.
> Compared to the city with zip code 59700 (3532.92 €/m²) and the city with zip code 59000 (3978.29 €/m²), the city with zip code 62100 has a significantly lower average price per square meter.
> - 👤 fais moi un résumé de tout ça dans un tableau en français
> - 🧠 
> 
> | Code Postal | Prix moyen au m² |
> |-------------|------------------|
> | 59700       | 3532.92          |
> | 59000       | 3978.29          |
> | 62100       | 1810.66          |
> 
> En résumé, le prix moyen au mètre carré varie considérablement entre ces trois villes. La ville avec le code postal 59000 a le prix moyen au mètre carré le plus élevé, suivie par la ville avec le code postal 59700. La ville avec le code postal 62100 a un prix moyen au mètre carré nettement inférieur aux deux autres.

## The tool

The tool we are developing is straightforward: a wrapper around duckdb exposing a single `query` function.
I won't go into details about the implementation, but basically it is just an exec around `duckdb`. I chose not to use an embedded library to keep things simple.

```go
func executeDuckDBQuery(queryStr string) (string, error) { ... }
```

### Encapsulation with the MCP Protocol

### Understandig the JSON-RPC calls and handshake

The protocol MCP is based on JSON-RPC. The host will call the tool via RPC call that will look like:

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"function_name","arguments":{"query":"THE SQL QUERY"}}}
```

### Handling requests in Go

Therefore, we need to write a function that can understand this structure and make the appropriate call to `executeDuckDBQuery`. 
Instead of manually parsing JSON-RPC requests, we leverage the `github/mark3labs/mcp-go` library, which simplifies request handling and ensures compatibility with the MCP protocol.

The JSON-RPC request is encapsulated in an `mcp.CallToolRequest` object, and the expected result is an `mcp.CallToolResult`. We just need to extract the query and call our function accordingly, like this:

```go
// Extracts the SQL query from the request and executes it using DuckDB
func duckDBHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	queryStr, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return mcp.NewToolResultError("query must be a string"), nil
	}
	res, err := executeDuckDBQuery(queryStr)
	if err != nil {
		return mcp.NewToolResultError("query_string encountered an error: " + err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}
```

The function needs to be exposed as a tool to be usable by the host LLM.
The server will declare its tooling when the host will send this request:

```json
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","clientInfo":{"name":"example-client","version":"1.0.0"},"capabilities":{}}}
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}
```

The server may reply something like:

```json
{"jsonrpc":"2.0","id":2,"result":{"tools":[{"name":"FUNCTION_NAME","description":"DESCRIPTION OF THE FUNCTION","inputSchema":{"type":"object","properties":{"query":{"description":"DESCRIPTION OF THE QUERY PARAMETER","type":"string"},"required":["qyuery"]}}]}}
```

Note that the descriptions you provide are very important, as they will be used by the LLM to choose the correct function and format the parameters.

The go-mcp library provides utility to do this:

```go
// Add tool
tool := mcp.NewTool("query_file",
  mcp.WithDescription("runs a SQL query through duckdb to extract the information of a file. The file can be local (containing '/'), or remote on hugginface (starting with 'hf:'). It may also contain wildcards ('*')"),
  mcp.WithString("query",
    mcp.Required(),
    mcp.Description("The SQL query to execute (compatible with DUCKDB)"),
  ),
)
```
## Exposing the tool

So far, we have encapsulated all the tooling in JSON-RPC calls.
We now need a transport layer.

MCP proposes two options:

- exposing and communicating through the network via Server-Sent Events (SSE).
- exposing via a local fork of the tool and communicating through STDIO.

We will use the latter option.

_Side Note:_ Even though it is possible to use any language for the tool, Go is a great choice. Its self-contained static binaries make it easy to distribute the tool without worrying about external dependencies.

Once again, the go-mcp library takes care of the boilerplate to register the tool with its associated handler:

```go
	// Create MCP server
	s := server.NewMCPServer(
		"DuckDB 🚀",
		"1.0.0",
	)
	// Add tool handler
	s.AddTool(tool, duckDBHandler)
```

## Compiling and running the tool

Once compiled, we can try to run the tool locally.

```bash
❯ go build -o duckdbserver
❯ # Get the capabilities of the server 
❯ echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./duckdbserver
{"jsonrpc":"2.0","id":2,"result":{"tools":[{"name":"query_file","description":"runs a SQL query through duckdb to extract the information of a file. The file can be local (containing '/'), or remote on hugginface (starting with 'hf:'). It may also contain wildcards ('*')","inputSchema":{"type":"object","properties":{"query":{"description":"The SQL query to execute (compatible with DUCKDB)","type":"string"}},"required":["query"]}}]}}
❯ # Run a query
❯ echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"query_file","arguments":{"query":"SELECT version() AS version;"}}}' | ./duckdbserver
{"jsonrpc":"2.0","id":3,"result":{"content":[{"type":"text","text":"┌─────────┐\n│ version │\n│ varchar │\n├─────────┤\n│ v1.1.3  │\n└─────────┘\n"}]}}
```

So far, so good... You notice that I did not specify any output format, I will let the LLM decide how to use it.

## Conclusion

The remaining work involves integrating all the components. I created a generic `MCPServerTool` structure on the host, which registers a tool configurable via the `MCP_SERVER` environment variable.

Here is the final schema of the code:
![](/assets/mcp/diagram.png)

The result, as you might have guessed, is that the initial dialog is now an actual conversation I conduct with the agent.

Consequently, adding a tool to an LLM while keeping my information private is now straightforward. This is because, when an MCP server provides resources, the LLM gains access to them.

Future steps include replacing VertexAI with, for example, ollama and managing multi-session elements.

With these final words, I can confidently state that I have achieved my objective and validated the core concepts.

If you want to try this server, the code is on [my GitHub](https://github.com/owulveryck/gomcptest)
