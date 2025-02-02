---
title: "MCP Part II - Implementation: Custom Host with VertexAI and Gemini"
date: 2025-01-27T12:15:33+01:00
lastmod: 2025-01-27T12:15:33+01:00
images: [/assets/mcp/partI/Image4.png)]
draft: true
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

Voici la traduction en français :

---

Dans la [première partie](/2025/01/27/mcp-part-i-core-concepts-past-present-and-future-of-agentic-systems.html) de cette série, j'ai exploré certains concepts et convictions concernant l'agentivité dans l'IA et le potentiel des outils pour améliorer les agents.

Les derniers mots portaient sur l’_hôte_ (l'application) exécutant des assistants alimentés par des LLM (_Claude_, _ChatGPT_, ...). Ces applications représenteront le véritable champ de bataille. Les entreprises qui obtiendront un monopole sur les systèmes d'assistants gouverneront le business numérique de demain.

C’est l’une des raisons pour lesquelles, dans mes expérimentations, je ne veux pas dépendre d’une application existante. Je veux comprendre les standards, comment fonctionne le câblage avec le LLM et acquérir de nouvelles compétences dans un environnement privé.

J’ai décidé d’implémenter un hôte "from scratch". Enfin, pas tout à fait, car je vais fortement m’appuyer sur le moteur d’exécution d’un LLM et son SDK.

J’aurais pu utiliser Ollama pour exécuter un modèle privé, mais j’ai décidé d’utiliser Vertex AI sur Google GCP dans un projet privé à la place. Premièrement, j’ai accès à GCP. Deuxièmement, l’API est relativement stable et j’aime généralement la façon dont les concepts sont exposés via l’API de Google. Cela m’aide à comprendre comment les choses fonctionnent.

J’ai également décidé d’utiliser Go comme langage de programmation parce que :

- Je connais Go.
- Go est fun (vraiment fun).
- Le SDK Vertex AI en Go existe (voir la doc [ici sur pkg.go.dev](https://pkg.go.dev/cloud.google.com/go/vertexai/genai)).

## Architecture générale

Je vais implémenter un moteur de chat basique. Je ne vais pas implémenter de frontend, car il existe déjà de très bons frontends open source. J’ai choisi d’utiliser [Big-AGI](https://big-agi.com/). Puisque je veux utiliser ce qui existe déjà, mon système de chat sera compatible avec l'[API v1 d’OpenAI](https://platform.openai.com/docs/api-reference/introduction). Il semble que ce soit un standard courant dans le monde de l’IA. Ainsi, je pourrai facilement connecter mon hôte avec Big-AGI.

J’ai choisi d’implémenter le point d’entrée `/v1/chat/completions` de manière _streaming_. Cela offrira la meilleure expérience utilisateur lorsqu’on discute et reproduira ce à quoi nous sommes habitués avec ChatGPT ou Claude.

## Implémentation du "Chat Handler"

### Comprendre le fonctionnement de la session de chat (avec VertexAI)

Le point d’entrée de l’API GenAI de VertexAI est un objet appelé [GenerativeModel](https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerativeModel).

La documentation indique :

> `GenerativeModel` est un modèle capable de générer du texte.

C’est donc le point d’entrée d’une conversation.

Le modèle peut générer du contenu (ou du contenu en streaming) via les appels à ses méthodes :

- `GenerateContent`
- `GenerateContentStream`

Cela convient à un appel unique, mais pour gérer une session de chat complète, l’API propose un autre objet : [`ChatSession`](https://pkg.go.dev/cloud.google.com/go/vertexai/genai#ChatSession).

La `ChatSession` gère l’historique du chat. En reprenant l’exemple de la Partie I de l’article, la `ChatSession` est responsable du maintien de la _fenêtre de contexte_.

![Un schéma dessiné à la main illustrant le Model Context Protocol (MCP). Un utilisateur interagit avec un environnement (par ex., un smartphone). Un agent traite l’entrée et communique avec un modèle via une fenêtre de contexte. Le modèle génère des réponses et l’agent peut modifier l’environnement. Des flèches rouges et bleues indiquent le flux bidirectionnel d’informations entre les composants.](/assets/mcp/partI/Image1.png)

La `ChatSession` possède deux méthodes :

- `SendMessage`
- `SendMessageStream`

Toutes deux encapsulent les méthodes `GenerateContent` et `GenerateContentStream` mentionnées précédemment.

## Implémentation du point d’entrée Chat

Dans notre implémentation, nous encapsulons `genai.ChatSession` et `genai.GenerativeModel` dans notre propre structure `ChatSession`.

```go
type ChatSession struct {
        cs                 *genai.ChatSession
        model              *genai.GenerativeModel
}
```

Puis, nous ajoutons un _handler_ HTTP à cette structure pour gérer l’appel à `/v1/chat/completion`.

```go
func main() {
    cs := NewChatSession()
    // Initialisation du client
    mux := http.NewServeMux()
    mux.Handle("/v1/chat/completions", http.HandlerFunc(cs.chatCompletionHandler))
}

func (cs *ChatSession) chatCompletionHandler(w http.ResponseWriter, r *http.Request) {...}
```

### Authentification au service VertexAI

L’authentification n’a pas encore été entièrement développée. J’utilise les méthodes d’authentification par défaut fournies par l’API, qui reposent sur les jetons d’authentification générés localement par le service `gcloud`.

_Note :_ Cela le rend assez portable vers un environnement GCP bien configuré, tel que Cloud Run.

```go
client, err := genai.NewClient(ctx, projectID, locationID)
```

### Gestion de la requête

Le `chatCompletionHandler` est chargé de décoder la requête POST et de la valider.
Il désérialise le JSON en un objet `ChatCompletionRequest`.

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

// ChatCompletionMessage représente un message unique dans la conversation.
type ChatCompletionMessage struct {
        Role         string                          `json:"role"`
        Content      interface{}                     `json:"content,omitempty"`
        Name         string                          `json:"name,omitempty"`
        ToolCalls    []ChatCompletionMessageToolCall `json:"tool_calls,omitempty"`
        FunctionCall *ChatCompletionFunctionCall     `json:"function_call,omitempty"`
        Audio        *ChatCompletionAudio            `json:"audio,omitempty"`
}
```

À noter que chaque message de complétion prend le modèle comme paramètre. Avec l’API OpenAI, il est donc possible d’utiliser plusieurs modèles au cours d’une conversation.

Le _handler_ doit également valider la session en cours et utiliser ou créer une `genai.ChatSession` en conséquence. Dans mon POC, cette partie n’a pas été implémentée. Par conséquent, je ne peux gérer qu’une seule session. Cela signifie que si je démarre une nouvelle conversation dans Big-AGI, j’hérite de l’historique de la précédente sans moyen de le supprimer (sauf en redémarrant l’hôte).

Si le mode _streaming_ est activé, on appelle une autre méthode :

```go
func (cs *ChatSession) streamResponse(
        w http.ResponseWriter, 
        r *http.Request, 
        request ChatCompletionRequest) {
  ...
}
```

C’est là que la magie opère.


