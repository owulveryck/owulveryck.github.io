---
title: "MCP Partie III - Application : Serveur personnalisé pour un cas d’usage spécifique"
date: 2025-01-27T12:15:33+01:00
lastmod: 2025-01-27T12:15:33+01:00
images: [/assets/mcp/diagram.png]
draft: false
keywords: []
summary: Cet ultime article d'une série en trois parties explore la séparation des outils de l’hôte grâce au Model Context Protocol (MCP), offrant flexibilité et réutilisabilité. L’auteur développe un outil basé sur MCP en Go pour exécuter des requêtes SQL via DuckDB, permettant une interaction fluide avec un chatbot tout en préservant la confidentialité. L’implémentation couvre la gestion des appels JSON-RPC, l’encapsulation des outils et leur intégration. Le projet valide l’efficacité de MCP, avec des perspectives d’évolution incluant le remplacement de VertexAI par Ollama et l’ajout de la gestion multi-session.
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

Dans les deux premières parties (voir [partie I](/fr/2025/01/27/mcp-partie-i-concepts-fondamentaux-passé-présent-et-futur-des-systèmes-agents.html) et [partie II](/fr/2025/01/27/mcp-partie-ii-implémentation-hôte-personnalisé-avec-vertexai-et-gemini.html)) de cette série, nous avons :

- décrit le système agentique et comment il pouvait être étendu grâce à l’utilisation d’outils ;
- expliqué en quoi un standard comme MCP était essentiel pour construire un écosystème vertueux ;
- construit un hôte « from scratch » : nous avons implémenté un chatbot privé et lui avons donné la capacité d’exécuter un outil.

Dans cette dernière partie, nous allons tenir la promesse initiale : nous allons créer un outil découplé de l’hôte.

En découplant l’outil de l’hôte, nous garantissons flexibilité, interopérabilité et réutilisabilité entre différents agents d’IA. Cela s’aligne avec l’objectif plus large de création d’un système modulaire et évolutif.

Voici le schéma cible de ce que nous décrivons dans l’article :

![L’image est un organigramme illustrant le Model Context Protocol (MCP) avec des blocs interconnectés et des flèches directionnelles. Le texte et les flèches en bleu représentent un ensemble de composants et de processus, tandis que ceux en rouge en représentent un autre. En haut à gauche, un bloc avec du texte bleu est étiqueté « Role », qui est connecté à un bloc étiqueté en rouge avec le texte « Model ». Ce bloc rouge a des flèches qui s’étendent vers la droite, menant à d’autres composants étiquetés en bleu et en rouge. Sous le bloc « Role », d’autres blocs étiquetés en bleu sont connectés verticalement, indiquant des relations hiérarchiques ou séquentielles. Des flèches rouges s’entrelacent avec les composants bleus, suggérant des interactions croisées ou des dépendances entre les deux flux. Le diagramme représente visuellement les relations structurelles et les interactions au sein du cadre MCP, mettant en avant l’interaction entre les éléments bleus et rouges.](/assets/mcp/partI/Image5.png)

## Exposition d’un cas d’usage

J’ai appris SQL au siècle dernier et l’ai à peine utilisé depuis. Ainsi, chaque fois que je dois explorer un jeu de données simple, j’utilise du code (un langage de programmation ou une combinaison de `grep`, `awk`, `sed` et leurs amis).

Cependant, dans la phase d’exploration, SQL présente des avantages, et un outil comme DuckDB peut être d’une grande aide.

Actuellement, j’utilise un chatbot pour générer des requêtes SQL que je copie-colle dans l’interface de DuckDB, mais ce n’est pas optimal.  
Mon cas d’usage consiste donc à créer un outil qui sera un serveur MCP, exposant une seule fonction :

- exécuter une requête via DuckDB (l’argument étant la requête).

Mon objectif est de pouvoir utiliser un chatbot pour faire quelque chose comme ceci :

> - 👤 Peux-tu me donner le schéma des données contenues dans le fichier local en exécutant une requête SQL (DuckDB) :  
>  `/tmp/ValeursFoncieres-2024-S1.tsv`  
> - 🧠 Le fichier `/tmp/ValeursFoncieres-2024-S1.tsv` contient 43 colonnes. Voici les noms et types de données des premières colonnes :
> 
> | Nom de colonne         | Type de données |
> |----------------------|----------------|
> | Identifiant de doc…  | varchar        |
> | Référence document   | varchar        |
> | ...                  | ...            |
> | Nature culture       | varchar        |
> | Nature culture spe…  | varchar        |
> | Surface terrain      | int64          |
> 
> Je peux fournir plus de détails si tu précises les colonnes qui t’intéressent.
> - 👤 Ok, ce fichier est en français.  
> Mon code postal est 59700, quel est le coût moyen des transactions dans cette ville ?
> 
> - 🧠 **Ah, il semble que je doive caster la colonne "Valeur foncière" en un type numérique avant de calculer la moyenne. Je vais essayer un CAST en DOUBLE.**  
> 
> **Le problème vient du format de la colonne "Valeur foncière". La virgule est utilisée comme séparateur décimal, ce qui n’est pas compatible avec la fonction CAST. Je vais remplacer la virgule par un point avant de caster en DOUBLE.**  
> 
> Le coût moyen des transactions dans la ville avec le code postal 59700 est de **441 545,00**.
> 
> - 👤 Quelle est la surface moyenne des terrains dans cette ville ?
> - 🧠 La surface moyenne des terrains dans la ville avec le code postal 59700 est de **444,60 mètres carrés**.
> - 👤 Quel est le prix moyen au mètre carré ?
> - 🧠 Le prix moyen au mètre carré dans la ville avec le code postal 59700 est de **3 532,92 €**.
> - 👤 Compare avec la ville dont le code postal est 59000.
> - 🧠 Le prix moyen au mètre carré dans la ville avec le code postal 59000 est de **3 978,29 €**.  
> Comparé à la ville avec le code postal 59700, où le prix moyen au mètre carré est de **3 532,92 €**, la ville avec le code postal 59000 a un prix moyen plus élevé.
> - 👤 Et par rapport à la ville dont le code postal est 62100 ?
> - 🧠 Le prix moyen au mètre carré dans la ville avec le code postal 62100 est de **1 810,66 €**.  
> Comparé à la ville avec le code postal 59700 (**3 532,92 €/m²**) et la ville avec le code postal 59000 (**3 978,29 €/m²**), la ville avec le code postal 62100 a un prix moyen au mètre carré nettement inférieur.
> - 👤 Fais-moi un résumé de tout ça dans un tableau en français.
> - 🧠 
> 
> | Code Postal | Prix moyen au m² |
> |------------|------------------|
> | 59700      | 3 532,92 €       |
> | 59000      | 3 978,29 €       |
> | 62100      | 1 810,66 €       |
> 
> En résumé, le prix moyen au mètre carré varie considérablement entre ces trois villes. La ville avec le code postal **59000** a le prix moyen au mètre carré le plus élevé, suivie par la ville avec le code postal **59700**. La ville avec le code postal **62100** a un prix moyen nettement inférieur aux deux autres.



## L'outil

L'outil que nous développons est simple : un wrapper autour de DuckDB exposant une seule fonction `query`.  
Je ne vais pas entrer dans les détails de l’implémentation, mais en gros, c’est juste une exécution autour de `duckdb`.  
J’ai choisi de ne pas utiliser une bibliothèque embarquée pour garder les choses simples.

```go
func executeDuckDBQuery(queryStr string) (string, error) { ... }
```

### Encapsulation avec le protocole MCP

### Comprendre les appels JSON-RPC et le handshake

Le protocole MCP est basé sur JSON-RPC. L’hôte appellera l’outil via un appel RPC qui ressemblera à ceci :

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"function_name","arguments":{"query":"LA REQUÊTE SQL"}}}
```

### Gestion des requêtes en Go

Nous devons donc écrire une fonction capable de comprendre cette structure et d’appeler correctement `executeDuckDBQuery`.  
Plutôt que d’analyser manuellement les requêtes JSON-RPC, nous utilisons la bibliothèque `github/mark3labs/mcp-go`, qui simplifie la gestion des requêtes et assure la compatibilité avec le protocole MCP.

La requête JSON-RPC est encapsulée dans un objet `mcp.CallToolRequest`, et le résultat attendu est un `mcp.CallToolResult`. Il suffit d’extraire la requête et d’appeler notre fonction comme ceci :

```go
// Extrait la requête SQL de la requête et l’exécute via DuckDB
func duckDBHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	queryStr, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return mcp.NewToolResultError("query doit être une chaîne de caractères"), nil
	}
	res, err := executeDuckDBQuery(queryStr)
	if err != nil {
		return mcp.NewToolResultError("Une erreur est survenue lors de l’exécution de la requête : " + err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}
```

### Exposer l'outil au LLM hôte

La fonction doit être exposée en tant qu’outil pour être utilisable par le LLM hôte.  
Le serveur déclarera ses outils lorsque l’hôte enverra cette requête :

```json
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","clientInfo":{"name":"example-client","version":"1.0.0"},"capabilities":{}}}
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}
```

Le serveur peut répondre avec quelque chose comme :

```json
{"jsonrpc":"2.0","id":2,"result":{"tools":[{"name":"FUNCTION_NAME","description":"DESCRIPTION DE LA FONCTION","inputSchema":{"type":"object","properties":{"query":{"description":"DESCRIPTION DU PARAMÈTRE QUERY","type":"string"},"required":["query"]}}]}}
```

À noter que les descriptions fournies sont très importantes, car elles seront utilisées par le LLM pour choisir la bonne fonction et formater les paramètres correctement.

La bibliothèque `go-mcp` fournit des utilitaires pour faire cela :

```go
// Ajouter un outil
tool := mcp.NewTool("query_file",
  mcp.WithDescription("Exécute une requête SQL via DuckDB pour extraire les informations d’un fichier. Le fichier peut être local (contenant '/'), ou distant sur Hugging Face (commençant par 'hf:'). Il peut aussi contenir des caractères génériques ('*')."),
  mcp.WithString("query",
    mcp.Required(),
    mcp.Description("La requête SQL à exécuter (compatible avec DUCKDB)"),
  ),
)
```

## Exposer l'outil

Jusqu'à présent, nous avons encapsulé tous les outils dans des appels JSON-RPC.  
Nous avons maintenant besoin d'une couche de transport.

MCP propose deux options :

- Exposer et communiquer via le réseau en utilisant les Server-Sent Events (SSE).
- Exposer l'outil via un fork local et communiquer via STDIO.

Nous allons utiliser cette dernière option.

_Remarque :_ Bien qu'il soit possible d'utiliser n'importe quel langage pour l'outil, Go est un excellent choix. Ses binaires statiques auto-contenus facilitent la distribution sans se soucier des dépendances externes.

Encore une fois, la bibliothèque `go-mcp` s'occupe du code répétitif pour enregistrer l'outil avec son gestionnaire associé :

```go
// Créer un serveur MCP
s := server.NewMCPServer(
	"DuckDB 🚀",
	"1.0.0",
)
// Ajouter un gestionnaire d'outil
s.AddTool(tool, duckDBHandler)
```

## Compilation et exécution de l'outil

Une fois compilé, nous pouvons essayer d'exécuter l'outil localement.

```bash
❯ go build -o duckdbserver
❯ # Obtenir les capacités du serveur 
❯ echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./duckdbserver
{"jsonrpc":"2.0","id":2,"result":{"tools":[{"name":"query_file","description":"Exécute une requête SQL via DuckDB pour extraire les informations d’un fichier. Le fichier peut être local (contenant '/'), ou distant sur Hugging Face (commençant par 'hf:'). Il peut aussi contenir des caractères génériques ('*').","inputSchema":{"type":"object","properties":{"query":{"description":"La requête SQL à exécuter (compatible avec DUCKDB)","type":"string"}},"required":["query"]}}]}}
❯ # Exécuter une requête
❯ echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"query_file","arguments":{"query":"SELECT version() AS version;"}}}' | ./duckdbserver
{"jsonrpc":"2.0","id":3,"result":{"content":[{"type":"text","text":"┌─────────┐\n│ version │\n│ varchar │\n├─────────┤\n│ v1.1.3  │\n└─────────┘\n"}]}}
```

Tout fonctionne bien jusqu’ici... Vous remarquez que je n’ai spécifié aucun format de sortie, laissant le LLM décider de son utilisation.

## Conclusion

Le travail restant consiste à intégrer tous les composants. J’ai créé une structure générique `MCPServerTool` côté hôte, qui enregistre un outil configurable via la variable d’environnement `MCP_SERVER`.

Voici le schéma final du code :  
![Architecture diagram showing client, API gateway, MCP server, and LLM data flow](/assets/mcp/diagram.png)

Le résultat, comme vous pouvez l’imaginer, est que le dialogue initial est maintenant une véritable conversation que je mène avec l’agent.

Ainsi, ajouter un outil à un LLM tout en gardant mes informations privées devient simple. En effet, lorsqu’un serveur MCP fournit des ressources, le LLM y accède.

Les prochaines étapes incluent le remplacement de VertexAI par, par exemple, Ollama et la gestion des éléments multi-sessions.

Avec ces dernières remarques, je peux affirmer que j’ai atteint mon objectif et validé les concepts clés.

Si vous souhaitez essayer ce serveur, le code est disponible sur [mon GitHub](https://github.com/owulveryck/gomcptest).

