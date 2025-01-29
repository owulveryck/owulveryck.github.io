---
title: "MCP Part II - Implementation: Custom Host with VertexAI and Gemini"
date: 2025-01-27T12:15:33+01:00
lastmod: 2025-01-27T12:15:33+01:00
images: [/assets/mcp/partI/Image4.png)]
draft: false
keywords: []
summary: This article details my journey in building a custom chat host for AI agents, moving away from existing solutions to gain a deeper understanding of the underlying technologies. I implement a chat engine using Google's Vertex AI and Go, focusing on compatibility with the OpenAI API to integrate with tools like Big-AGI. The article covers the core architecture, including my use of ChatSession and GenerativeModel from the Vertex AI SDK. It delves into the implementation of the /v1/chat/completions endpoint, highlighting the challenges of streaming responses and integrating function calls. I also describe a workaround for handling function calls in a streaming context and introduce the concept of a callable interface to prepare for implementing the Model Context Protocol (MCP) in future work. The goal is to move the tools outside of the agent. This will be detailes in the last part of this series.
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

In the [first part](/2025/01/27/mcp-part-i-core-concepts-past-present-and-future-of-agentic-systems.html) of this series, I explored some concepts and convictions regarding agentivity in AI and the potential of tooling to enhance agents.

The final words were about the _host_ (the application) running LLM-powered assistants (_Claude_, _ChatGPT_, ...). These applications will represent the true battleground. The companies that gain a monopoly on assistant systems will govern tomorrow's digital business.

This is one reason why, for my experimentation, I do not want to rely on an existing application. I want to understand the standards, how the wiring with the LLM works, and gain new skills in a private environment.

I have decided to implement a host "from scratch." Well, not really from scratch, because I will heavily rely on the execution engine of an LLM and its SDK.

I could have used Ollama to run a private model, but I decided to use Vertex AI on Google GCP in a private project instead. First, I have access to GCP. Second, the API is fairly stable, and I generally like how the concepts are exposed via Google's API. It helps me understand how things work.

I have also decided to use Go as a programming language because:

- I know Go.
- Go is fun (really fun).
- The Vertex AI SDK in Go exists (see doc [here at pkg.go.dev](https://pkg.go.dev/cloud.google.com/go/vertexai/genai).

## Overall architecture

I will implement a basic chat engine. I will not implement a frontend because very good open-source frontends already exist. I have chosen to use [Big-AGI](https://big-agi.com/). Since I want to use what is already existing, my chat system will be compatible with [OpenAI's v1 API](https://platform.openai.com/docs/api-reference/introduction). It appears this is a common standard in the AI world. Therefore, I will be able to easily connect my host with Big-AGI.

I have chosen to implement the `/v1/chat/completions` endpoint in a streaming manner. This will provide the best user experience when chatting and will replicate what we are used to in ChatGPT or Claude.

## Implementation of the "Chat Handler"


### Understanding how the chat session works (with VertexAI)

The entry point of the VertexAI's GenAI API is an object called the [GenerativeModel](https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerativeModel).

The documentation states:

> `GenerativeModel` is a model that can generate text.

It is therefore the entry point for a conversation.

The model can generate content (or streamed content) via calls to its methods:

- `GenerateContent`
- `GenerateContentStream`

This is suitable for a single call, but to handle a full chat session, the API offers another object: the [`ChatSession`](https://pkg.go.dev/cloud.google.com/go/vertexai/genai#ChatSession).

The `ChatSession` manages the chat history. Taking the example from Part I of the article, the `ChatSession` is responsible for maintaining the _Context Window_.

![A hand-drawn diagram illustrating the Model Context Protocol (MCP). A user interacts with an environment (e.g., a smartphone). An agent processes input and communicates with a model through a context window. The model generates responses, and the agent can modify the environment. Red and blue arrows indicate bidirectional information flow between components.](/assets/mcp/partI/Image1.png)

The `ChatSession` has two methods:

-  `SendMessage`
-  `SendMessageStream`

Both encapsulate the `GenerateContent` and `GenerateContentStream` methods mentioned earlier.

## Implementing the Chat endpoint

In our implelentation, we are encapsulating the `genai.ChatSession` and `genai.GenerativeModel` into our own structure `ChatSession`.

```go
type ChatSession struct {
	cs                 *genai.ChatSession
	model              *genai.GenerativeModel
}
```
Then we add a http handler to this structure to handle the call to `/v1/chat/completion`


```go
func main() {
    cs := NewChatSession()
    // Initialize the client
    mux := http.NewServeMux()
    mux.Handle("/v1/chat/completions", http.HandlerFunc(cs.chatCompletionHandler))
}

func (cs *ChatSession) chatCompletionHandler(w http.ResponseWriter, r *http.Request) {...}
```

### Authentication to the VertexAI service

The authentication hasn't been fully developed. I'm using the default authentication methods provided by the API, which rely on the authentication tokens generated by the `gcloud` service locally.

_Note:_ This makes it fairly portable to a well-configured GCP environment, in a service such as Cloud Run.

```go
client, err := genai.NewClient(ctx, projectID, locationID)
```

### Handling the request

The `chatCompletionHandler` is in charge of decoding the POST request and validating it.
It deserializes the JSON payload into an object `ChatCompletionRequest`

```go
type ChatCompletionRequest struct {
	Model         string                  `json:"model"`
	Messages      []ChatCompletionMessage `json:"messages"`
	MaxTokens     int                     `json:"max_tokens"`
	Temperature   float64                 `json:"temperature"`
	Stream        bool                    `json:"stream"`
	StreamOptions struct {
		IncludeUsage bool `json:"include_usage"`
	} `json:"stream_options"`
}

// ChatCompletionMessage represents a single message in the chat conversation.
type ChatCompletionMessage struct {
	Role         string                          `json:"role"`
	Content      interface{}                     `json:"content,omitempty"` // Can be string or []map[string]interface{}
	Name         string                          `json:"name,omitempty"`
	ToolCalls    []ChatCompletionMessageToolCall `json:"tool_calls,omitempty"`
	FunctionCall *ChatCompletionFunctionCall     `json:"function_call,omitempty"`
	Audio        *ChatCompletionAudio            `json:"audio,omitempty"`
}
```

It's worth noting that each completion message takes the model as a parameter. Therefore, with the OpenAI API, it's possible to use multiple models throughout a conversation.

The handler should also be responsible for validating the current session and using or creating a `genai.ChatSession` accordingly. In my proof of concept (POC), I haven't implemented this. As a consequence, I can only handle one session. This means that if I start a new conversation in Big-AGI, I will inherit the history of the previous one, and I have no way to delete it (unless I restart the host).

If the streaming is enabled in the request, we call another method:

```go
func (cs *ChatSession) streamResponse(
	w http.ResponseWriter, 
	r *http.Request, 
	request ChatCompletionRequest) {
  ...
}
```

This is where the magic happens.

### Sending/Streaming the response

This method uses the `genai.ChatSession.SendMessageStream` method to ask the model to generate content. This function takes `Message` [Part](https://pkg.go.dev/cloud.google.com/go/vertexai/genai#Part) as parameters. I've created a simple method to convert the `Messages` from the OpenAI API into Google's `genai.Part`. Note that it also handles images.

```go
func (c *ChatCompletionMessage) toGenaiPart() []genai.Part {...}
```

The `SendMessageStream` returns an iterator (Google's own implementation of an iterator, likely predating the official Go language iterator).

We need to iterate to get the full response from the model and serialize it into a `ChatCompletionStreamResponse` that will be sent back to the Big-AGI client.

```go
type ChatCompletionStreamResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []ChatCompletionStreamChoice `json:"choices"`
}

type ChatCompletionStreamChoice struct {
	Index        int         `json:"index"`
	Delta        ChatMessage `json:"delta"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}
```

The API can send multiple choices, but for the purpose of this POC, I've chosen to only send one. 

This is it, with this simple implementation, I can use the service as a backend to Big-AGI

## Implementing functions

My goal is to extend the capabilities of the host, and eventually to use the MCP protocol to do so.
The first step is to increase the capabilities of the model by providing functions.

The GenerativeModel's API has exported fields. One of them it `Tools`, and this looks like a good starting point:

```go
type GenerativeModel struct {
	GenerationConfig
	SafetySettings    []*SafetySetting
	Tools             []*Tool
	ToolConfig        *ToolConfig // configuration for tools
	SystemInstruction *Content
	// The name of the CachedContent to use.
	// Must have already been created with [Client.CreateCachedContent].
	CachedContentName string
	// contains filtered or unexported fields
}
```

In the doc, it is mentioned that a [`Tool`](https://pkg.go.dev/cloud.google.com/go/vertexai/genai#Tool), is:

> a piece of code that enables the system to interact with external systems to perform an action, or set of actions, outside of knowledge and scope of the model. A Tool object should contain exactly one type of Tool (e.g FunctionDeclaration, Retrieval or GoogleSearchRetrieval).

By implementing tools, we will be in the situation described by this diagram:

![A hand-drawn diagram illustrating the Model Context Protocol (MCP) with tool integration. A user interacts with a device (Environment ①), sending input to an agent, which communicates with a model through a context window. The model acknowledges the use of a tool, which has been programmed to read external content (e.g., a blog in Environment ②) and provide it to the model for processing."](/assets/mcp/partI/Image4.png)

Following the example of the Doc, we can create and add a tool to the model an follow the workflow:

```
// For using tools, the chat mode is useful because it provides the required
// chat context. A model needs to have tools supplied to it in the chat
// history so it can use them in subsequent conversations.
//
// The flow of message expected here is:
//
// 1. We send a question to the model
// 2. The model recognizes that it needs to use a tool to answer the question,
//    an returns a FunctionCall response asking to use the CurrentWeather
//    tool.
// 3. We send a FunctionResponse message, simulating the return value of
//    CurrentWeather for the model's query.
// 4. The model provides its text answer in response to this message.
```

### The problem with the streaming

Initially, I implemented the function as described in the example, `find_theater`. However, implementing it via streaming didn't work as expected. When I tried to send the `FunctionResponse` message, the model returned a 400 error:

`Please ensure that function call turns come immediately after a user turn or after a function response turn.`

The issue was that the iterator wasn't empty, and therefore the model was receiving the function response before the function request had been properly set.

I implemented a workaround using a stack of functions with `Push` and `Pop` functions. If the request is a function call, I `Push` the request onto the stack. When the iterator is empty, I `Pop` the function, execute it, and send its response back to the `ChatSession` with `SendMessageStream`. The reply is a new iterator that is used to complete the request sent to the end-user.

```go
type FunctionCallStack struct {
	mu    sync.Mutex
	items []genai.FunctionCall
}

// Push adds a genai.FunctionCall to the top of the stack.
func (s *FunctionCallStack) Push(call genai.FunctionCall) {...}

// Pop removes and returns the last genai.FunctionCall from the stack (FIFO).
// Returns nil if the stack is empty
func (s *FunctionCallStack) Pop() *genai.FunctionCall {...}
```

I'm aware that this isn't a bulletproof solution, but it works for my POC. 

## Conclusion and preparation of the MCP

This host is functional, and the model calls the function when it deems it useful. However, my goal is to move the tool out of the agent.

This will be implemented in the last part of this series of articles. To prepare for this, I've implemented an abstraction of the function call, so I can avoid tweaking the `chatCompletion` handlers.

I've created an interface `callable` with the following methods:

```go
type callable interface {
	GetGenaiTool() *genai.Tool
	Run(genai.FunctionCall) (*genai.FunctionResponse, error)
	Name() string
}
```

Then, I've updated my `ChatSession` structure with an inventory of `callable`:

```go
type ChatSession struct {
	cs                 *genai.ChatSession
	model              *genai.GenerativeModel
	functionsInventory map[string]callable
}
```

Now, I can register all the `callable`s as `Tools` to the `GenerativeModel` (by calling the `GetGenaiTool()` method), and within a chat session, I can detect whether the function's name is part of the inventory, call the `Run()` method, and send the reply back. 

## Final Note About the function call and Conclusion

It is interesting to note that the model decides by itself whether it needs to call the function.
As explained before, the model only deals with text. Therefore, the description of the function and its parameters are key. It is really about convincing the model that this function is useful in its reasoning.

On top of that, it is important to note that the result of the function is also injected into the context.
The fun part is that I can easily display the history and all the exchanges between the user and the model. It teaches me a lot about the mechanisms of those chatbots.

The strange part about engineering is the lack of idempotence. You can ask the same question twice and get different answers. As Luc De Brabandere wrote, we are now back to the statistical era: it works x% of the time, and we find that good enough.

This proof of concept represents the technological foundation of my POC with the Model Context Protocol.
Now it is fairly easy to add new functions. It is a matter of fulfilling the `callable` interface.
Please note that this implementation is partial because a `Tool` can expose several functions, which is not possible with my current implementation.

In the next and final part of this series, I will create an MCP client implementing the `callable` interface and a sample MCP server.
This server will register its capabilities automatically.

Therefore, providing tools to my chatbot will become easy.

If you want to try this server, the code is on [my GitHub](https://github.com/owulveryck/gomcptest/tree/main/host/openaiserver)
