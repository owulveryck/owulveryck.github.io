---
title: "ChatGPT, Wardley et Go"
date: 2023-05-30T21:23:17+02:00
lastmod: 2023-05-30T21:23:17+02:00
draft: false
keywords: []
description: "Cet article décrit comment j'utilise le SDK wardleyToGo pour créer un plugin en Go pour ChatGPT (pour afficher des Cartes de Wardley)"
tags: []
categories: []
author: "Olivier Wulveryck"

# You can also close(false) or open(true) something for this content

# P.S. comment can only be closed

comment: true
toc: true
autoCollapseToc: false

# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."

contentCopyright: false
reward: false
mathjax: false
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/2023-05-30-chatGPTWardley/).

Dans cet article, j'explique :

- comment créer un plugin ChatGPT avec Go
- comment valider la configuration avec [CUE](cuelang.org)
- comment créer une API basique utilisable avec ChatGPT
- comment afficher des images SVG dans ChatGPT (plutôt ce qu'il faut faire et ne pas faire)

## Introduction

J'utilise ChatGPT quotidiennement comme assistant, non pas comme un dictionnaire ou une encyclopédie.
Je cherche des solutions à des problèmes et je suis conscient que je les trouverai par moi-même.
Le rôle de ChatGPT est de m'assister et d'aider les solutions à émerger de mon esprit.
Je pose des questions, et avec ses connaissances, il façonne ma façon de penser pour converger vers les solutions.

Pour la plupart des problèmes stratégiques, j'utilise les outils et techniques du Wardley Mapping, notamment :

- La chaîne de valeur (avec les besoins des utilisateurs en haut et quelques autres principes de doctrine)
- La théorie de l'évolution
- Les modèles climatiques

Construire une carte a de la valeur, et le défi de positionner les différents composants sur la carte elle-même apporte beaucoup de valeur.

Bien que ChatGPT ne soit pas conscient de la représentation spatiale, il peut fournir des justifications sur le placement.
Cependant, en tant qu'humain, une représentation visuelle est très utile.

Par conséquent, en tant que geek, j'ai commencé à réfléchir à un plugin qui permettrait à ChatGPT de dessiner une carte.

Dans un article précédent, j'ai présenté wardleyToGo, un SDK pour construire des Cartes de Wardley en utilisant du code Go.
En tant qu'abonné à ChatGPT, je peux écrire un plugin pour GPT-4.

Cet article est un voyage qui explique comment j'ai construit un plugin pour dessiner mes cartes, comment il fonctionne, ce que j'ai découvert, et plus encore.

Vous pourriez vouloir lire cet article si :

- Vous êtes un cartographe Wardley curieux.
- Vous êtes un utilisateur de ChatGPT et souhaitez être au courant des possibilités à venir.
- Vous êtes un développeur Go et souhaitez vous familiariser avec la tuyauterie nécessaire pour créer un plugin pour ChatGPT.

## Comment fonctionne un plugin

Le développement de plugins pour ChatGPT est documenté [ici](https://platform.openai.com/docs/plugins/introduction).
En un coup d'œil, un plugin est une API REST qui est appelée par ChatGPT.

Pour transformer une API en plugin, vous devez fournir deux fichiers :

- un fichier [manifeste de plugin](https://platform.openai.com/docs/plugins/getting-started/plugin-manifest) servi à `/.well-known/ai-plugin.json`
- la spécification openAPI servie via un chemin spécifié dans le fichier `ai-plugin.json`.

Le format du manifeste est important car, outre la sérialisation en JSON, le format est contraint.
Par exemple, le champ `name_for_model` ne doit pas contenir d'espace et doit être limité à 50 caractères maximum.

![](/assets/gpt/gpt.png)

### Tuyauterie Golang

Dans cette section, je décris la tuyauterie que j'ai utilisée dans le fichier Go pour créer le plugin.
Cette section n'est pas liée à la fonctionnalité principale et peut probablement être utilisée pour n'importe quel plugin Golang.

_Avertissement_ : Vous devez être familier avec les concepts de gestionnaires web et le fonctionnement de Go pour bénéficier de cette section.

Je suis actuellement dans la phase "le faire fonctionner", donc tout le code réside dans le package principal d'un fichier.
Pour faciliter la modification du code, j'ai créé deux gestionnaires distincts :

- Un pour la tuyauterie, qui sert le manifeste, la spécification, la racine et le logo.
- Un pour l'API elle-même.

La tuyauterie est gérée par la structure ChatGPTPlumbing, qui implémente l'interface [`http/handler`](https://pkg.go.dev/net/http#Handler).
Cette structure lit et génère le contenu du manifeste et de l'OpenAPI lors de la création et le met en cache en interne.

```go
type ChatGPTPlumbing struct {
 aiPlugin        *AIPlugin
 aiPluginPayload []byte
 openAPIFile     string
 openAPIContent  []byte
}

// ChatGPTPlumbing implements the http/handler interface 
func (chatgptplumbing *ChatGPTPlumbing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 mux := http.NewServeMux()
 mux.HandleFunc("/.well-known/ai-plugin.json", func(w http.ResponseWriter, _ *http.Request) { // ...
 mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, _ *http.Request) { // ...
 mux.HandleFunc("/logo.png", func(w http.ResponseWriter, _ *http.Request) { // ...
 mux.ServeHTTP(w, r)
}
```

La tuyauterie est créée dans la fonction `main` et enregistrée pour gérer les requêtes situées à `/`

```go
plumbing, _ := NewChatGPTPlumbing(aiPlugin)
mux := http.NewServeMux()
mux.Handle("/", plumbing)
```

Le manifeste et l'OpenAPI sont lus (et validés) lors de la création du plugin et mis en cache dans la structure.

#### Créer et valider le `aiplugin.json`

Le manifeste doit être instancié au moment de l'exécution pour définir le port d'écoute et l'adresse corrects.
J'ai créé une structure Golang pour gérer le contenu et utilisé le package json pour le sérialiser.

De plus, comme mentionné précédemment, le fichier aiplugin.json a des contraintes strictes.
Comme je suis actuellement dans la phase "_le faire fonctionner_" et que je change fréquemment le contenu à des fins de test, il est préférable pour moi de valider les contraintes du plugin chaque fois que je démarre le plugin.

Pour ce faire, je m'appuie sur le langage CUE.
J'ai créé un simple fichier de contraintes que je combine avec un fichier de configuration au démarrage pour générer le contenu de l'API.

- constraints.cue

```CUE
Host: string | *"AUTO"

#AIPlugin: {
 // SchemaVersion Manifest schema version - required - v1
 schema_version: string & "v1"

 // NameForHuman Human-readable name, such as the full company name. 20 character max. - required
 name_for_human: string & =~"^.{1,20}$"

 // NameForModel Name the model will use to target the plugin (no spaces allowed, only letters and numbers). 50 character max. - required
 name_for_model: string & =~"^[a-zA-Z0-9]{1,50}$"

 // DescriptionForHuman Human-readable description of the plugin. 100 character max.
 description_for_human: string & =~"^.{1,100}$"

 // DescriptionForModel Description better tailored to the model, such as token context length considerations or keyword usage for improved plugin prompting. 8,000 character max. - required
 description_for_model: string & =~"^.{20,1000}$"
 auth:                  #Auth
 api:                   #API
 logo_url:              string | *"\(Host)/logo.png"
 contact_email:         string & =~"^.*@.*$"
 legal_info_url:        string | *"\(Host)/legal"
}

#Auth: {
 type: string | *"none"
}

#API: {
 type:                  string | *"openapi"
 url:                   string | *"\(Host)/openapi.yaml"
 is_user_authenticated: bool | *false
}
```

- configuration.cue

```cue
configuration: #AIPlugin & {
 name_for_human:        "Wardley To Go"
 name_for_model:        "WardleyToGo"
 description_for_human: "This plugin draw Wardley Maps"
 description_for_model: "This plugin draw Wardley Maps"
 contact_email:         "me@address.com"
}
```

Le `Host` est complété au moment de l'exécution, et tout est combiné pour générer la structure go qui est ensuite sérialisée :

```go
host := `Host: "` + address + `"`
constraints, err := ioutil.ReadFile("constraints.cue")
configuration, err := ioutil.ReadFile("wellknown.cue")
content := append(constraints, configuration...)
content = append([]byte(host+"\n"), content...)
ctx := cuecontext.New()
v := ctx.CompileBytes(content)
v = v.Lookup("configuration")
var aiplugin AIPlugin
err = v.Decode(&aiplugin)
```

#### Créer et servir l'OpenAPI

L'OpenAPI doit également être ajusté pour définir la description et les serveurs correspondants.

- **OpenAPI.yaml : La manière facile**
La première tentative pour construire un fichier `openai.yaml` polyvalent était de créer un modèle golang et de l'analyser au moment de l'exécution.
Le problème est que mettre des modèles dans YAML mène au cauchemar YAML...
J'ai donc utilisé une méthode plus amusante et geek.

- **La manière Geek**
Je génère maintenant le openapi.json avec CUE également.
Cela offre deux avantages :
  - Je n'ai pas à écrire la spécification au format OpenAPI (ce qui signifie que je n'ai pas à me battre avec yaml ou JSON)
  - Je peux valider la charge utile envoyée par ChatGPT.

Le code nécessite quelques tests supplémentaires, mais finira par arriver dans le dépôt.

#### Tuyauterie réseau et configuration

Comme d'habitude, la configuration est gérée via des variables d'environnement.
Cela permet de définir le port d'écoute et l'adresse.
J'ai également implémenté le tunneling avec `ngrok-go` qui me donne la possibilité de tester le plugin sur un hôte distant.

_Note annexe : Le [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)_:
J'ai créé un "middleware" très simple pour gérer les requêtes préliminaires CORS :

Outre les origines autorisées `https://chat.openai.com` et `http://serveraddress:port`, ces en-têtes sont requis :

```go
w.Header().Set("Access-Control-Allow-Origin", origin)
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, openai-conversation-id, openai-ephemeral-user-id, openai-*, sentry-trace, baggage")
```

## Conception de l'API

Maintenant que nous avons toute la tuyauterie pour rendre une API REST compatible avec ChatGPT, concevons l'API.

Le premier objectif simple est de pouvoir afficher l'axe d'évolution et de laisser ChatGPT placer un composant dessus.

Nous allons donc créer un point de terminaison `/mapEvolution` ; Il gérera une requête POST avec une charge utile spécifique.

### Création d'un point de terminaison basique pour afficher l'axe d'évolution

Pour générer un simple SVG d'évolution avec un seul composant à partir du code (et avec le SDK `wardleyToGo`), ce dont j'ai besoin est :

- le nom du composant
- la position du composant sur une échelle de 0 à 100

Cela constituera la charge utile de la requête au point de terminaison `/mapEvolution`.

La conception de la documentation de l'API est essentielle, car ChatGPT "lit la doc" et génère une charge utile par lui-même.
Je dois donc guider ChatGPT et lui expliquer comment fonctionne l'axe d'évolution :

```yaml
component:
  type: string
  description: The component to add to the map
evolution:
  type: int
  description: |
    The position on the evolution axis between 0 and 100. 
    From 0 to 17 the compoenent is in stage 1 (genesis for an asset or a an activity, novel for a practice, concept for some general knowledge)
    From 18 to 40 the component is in stage 2 (custom for an asset or an activity, emerging for a practice, hypothesis for some general knowledge) 
    From 40 to 70 the component is in stage 3 (product for and asset or an activity, good for a practice, or theory for some general knowledge)
    From 70 to 99 the component is in stage 4 (commodity for an asset of an activity, best for a practice, accepted for some general knwoledge)
```

La structure Go correspondante est

```go
type EvolutionInput struct {
 Component string   `json:"component"`
 Evolution int      `json:"evolution"`
}
```

Je ne détaillerai pas l'implémentation du gestionnaire http car c'est un développement Go standard.

Dans le gestionnaire, je crée une structure [`wardleyToGo.Map`](https://pkg.go.dev/github.com/owulveryck/wardleyToGo#Map), je remplis la carte avec un [`wardley.Component`](https://pkg.go.dev/github.com/owulveryck/wardleyToGo@v0.9.1/components/wardley#Component).

Puis je crée un [encodeur SVG](https://pkg.go.dev/github.com/owulveryck/wardleyToGo@v0.9.1/encoding/svg#Encoder) sur le [`http.ResponseWriter`](https://pkg.go.dev/net/http#ResponseWriter) pour renvoyer le résultat à l'interface utilisateur de ChatGPT.

### Test du plugin

Une fois que j'ai démarré mon serveur, je peux essayer le plugin :

- installation de la version de développement

![installer le plugin](/assets/gpt/newplugin.png)

![installer le plugin](/assets/gpt/localhostplugin.png)

- envoi de la première requête :

![première requête](/assets/gpt/evolutionk8s.png)

Nous voyons que le moteur GPT-4 a compris la requête et a :

- évalué la position de kubernetes selon ses connaissances
- généré une charge utile selon l'API
- utilisé la description pour placer l'évolution sur l'échelle 0..100

Le plugin a reçu la requête, a généré la carte et envoyé le résultat :

![première réponse](/assets/gpt/badreply.png)

Le problème est que le moteur essaie d'analyser la réponse et de la formater comme un contenu markdown.
Voici l'interprétation :

```text
Here is the generated evolution map for kubernetes:

![Evolution Map](data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiBwcmVzZXJ2ZUFzcGVjdFJhdGlvPSJ4TWlkWU1pZCBtZWV0IiB4bWxuczpldj0iaHR0cDovL3d3dy53My5vcmcvMjAwMS94bWwtZXZlbnRzIiBzdHlsZT0ib3ZlcmZsb3c6IGhpZGRlbjsgIj48ZGVmcz4KICAgICAgICA8bGluZWFyR3JhZGllbnQgaWQ9IndhcmRsZXlHcmFkaWVudCIgeDE9IjAlIiB5MT0iMCUiIHgyPSIxMDAlIiB5Mj0iMCUiPgogICAgICAgICAgICA8c3RvcCBvZmZzZXQ9IjAlIiBzdG9wLWNvbG9yPSJyZ2IoMjM2LDIzNywyNDMpIj48L3N0b3A+CiAgICAgICAgICAgIDxzdG9wIG9mZnNldD0iMzAlIiBzdG9wLWNvbG9yPSJyZ2IoMjU1LDI1NSwyNTUpIj48L3N0b3A+CiAgICAgICAgICAgIDxzdG9wIG9mZnNldD0iNzAlIiBzdG9wLWNvbG9yPSJyZ2IoMjU1LDI1NSwyNTUpIj48L3N0b3A+CiAgICAgICAgICAgIDxzdG9wIG9mZnNldD0iMTAwJSIgc3RvcC1jb2xvcj0icmdiKDIzNiwyMzcsMjQzKSI+PC9zdG9wPgogICAgICAgIDwvbGluZWFyR3JhZGllbnQ+CiAgICAgICAgPG1hcmtlciBpZD0iYXJyb3ciIHJlZlg9IjE1IiByZWZZPSIwIiBtYXJrZXJXaWR0aD0iMTIiIG1hcmtlckhlaWdodD0iMTIiIHZpZXdCb3g9IjAgLTUgMTAgMTAiPgogICAgICAgICAgICA8cGF0aCBkPSJNMCwtNUwxMCwwTDAsNSIgZmlsbD0icmdiKDI1NSwwLDApIj48L3BhdGg+CiAgICAgICAgPC9tYI apologize for the confusion, but it seems there was an error in rendering the SVG image. Let's try again
```

Nous voyons ici que le moteur GPT essaie d'encoder l'image SVG dans sa représentation base64 pour pouvoir l'afficher.
C'est lent et se termine évidemment par une erreur.

### Affichage du SVG, astuces et conseils

Maintenant que je suis conscient que ChatGPT évalue la réponse, la seule façon que j'ai trouvée pour afficher l'image était d'envoyer une référence à une image.

En envoyant ce genre de réponse :

```json
{
  "ImageURL": "http://localhost:3333/api/svg/5dedf28b-f683-474f-b694-dde318cbb1cb.svg"
}
```

Le moteur GPT comprend que c'est une image et génère cette réponse :

```text
![Evolution Map](http://localhost:3333/api/svg/5dedf28b-f683-474f-b694-dde318cbb1cb.svg)
```

Qui est alors affichée correctement dans le navigateur.

Ce que j'ai fait pour l'instant, c'est que je génère et sauvegarde la carte en interne et cela répond à mes propres besoins.
Le problème est que je ne peux pas publier le plugin car je serais capable de voir toutes les cartes générées par tous les utilisateurs du plugin.

_Note_ : C'est une autre leçon du voyage des plugins : les plugins peuvent être un véritable problème de sécurité. Lorsque vous utilisez un plugin, vous consentez à partager des informations avec des tiers.

À des fins de test, j'ai développé un stockage en mémoire.
Le problème de ce stockage éphémère est que lorsque je veux consulter un ancien chat, il essaie de recharger les anciennes images, ce qui se termine par une erreur 404.
Par conséquent, j'ai également instancié un simple stockage sur disque qui sauvegarde toutes les cartes que j'ai générées.

## Conclusion

Jusqu'à présent, j'ai atteint mon objectif et je peux maintenant utiliser ChatGPT comme assistant. Il affichera les composants d'une carte, mais ce n'est qu'un début.
Maintenant, je vais continuer à travailler sur l'API pour lui donner la capacité de construire une carte complète.

Une autre idée intéressante simple à développer avec le SDK `wardleyToGo` est la possibilité pour ChatGPT d'analyser une carte sauvegardée sur onlinewardleymaps.com.

Par exemple, je pourrais demander à ChatGPT :

> que penses-tu de cette carte : https://onlinewardleymaps.com/#UtzyxpPElI1ZUjABuH

Ensuite, il enverra la requête au plugin qui :

- récupérera la représentation OWM
- construira la représentation intermédiaire (une `wardleyToGo.Map`)
- extraira du sens (par exemple : ce composant est à l'étape blabla)

donnera le résultat.

Je suis heureux d'avoir conçu wardleyToGo comme un SDK, maintenant, pour moi, le ciel est la limite !

![](/assets/gpt/gpt-evolution-result.png)

Références :

- le [SDK `wardkeyToGo`](https://github.com/owulveryck/wardleyToGo)
- le [code du plugin](https://github.com/owulveryck/wardleyToGo/tree/chatGPTPlugin/examples/chatGPT) jusqu'à présent
- un autre article sur [wardleyToGo](https://blog.owulveryck.info/2023/03/02/rationale-behind-wardleytogo.html)