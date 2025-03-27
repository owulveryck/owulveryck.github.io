---
title: "Simplifier la complexité : Le voyage des WebSockets aux flux HTTP"
date: 2023-12-02T08:26:41+01:00
lastmod: 2023-12-02T08:26:41+01:00
draft: false
images: [/assets/crowdasleep_small.png]
videos: [/assets/present.webm]
eywords: []
summary: Cet article explore la transition d'une implémentation basée sur WebSocket vers un flux plus simple et plus direct via HTTP dans le contexte de la capture d'entrées tactiles sur Linux. 
  
  
  Il commence par introduire le thème principal, résumé dans l'affirmation _Tout ~~est un fichier~~ est un flux d'octets._ 
  Le besoin de capturer les positions des doigts sur un écran tactile en lisant `/dev/input/events` sous Linux est d'abord discuté, suivi d'un dilemme concernant le transfert de ces données vers un client JavaScript dans un navigateur.
  
  
  Initialement, les WebSockets sont choisis, ce qui conduit à une discussion sur la façon dont les frameworks façonnent souvent nos choix technologiques et les défis rencontrés lors du débogage des connexions WebSocket.
  L'article présente ensuite une alternative concernant l'envoi d'un flux d'octets via HTTP, établissant un parallèle avec l'approche de Linux pour gérer les périphériques et les fichiers.
  
  
  La sérialisation, le processus d'encodage des messages pour ce flux, est discutée ensuite, mettant en évidence les spécificités d'implémentation en GoLang et ses avantages natifs. 
  La dernière section traite de la façon de recevoir et de décoder ce flux en JavaScript au sein d'un thread worker, puis d'envoyer les messages décodés au thread principal à l'aide de requêtes post.
  
  L'article conclut en réfléchissant aux avantages de la simplicité en technologie, exhortant les lecteurs à réévaluer les choix par défaut et à envisager des solutions plus directes à des problèmes complexes.
tags: ["stream", "reMarkable", "http", "golang", "websocket", "javascript"]
categories: ["tools"]
author: "Olivier Wulveryck"

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: false
toc: true
autoCollapseToc: true
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: false

# Uncomment to add to the homepage's dropdown menu; weight = order of article
# menu:
#   main:
#     parent: "docs"
#     weight: 1
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/20231202-websocket-stream/).

## Introduction

Pour ajouter une nouvelle fonctionnalité à mon outil, [goMarkableStream](https://github.com/owulveryck/goMarkableStream),
j'avais besoin de capturer les positions des gestes sur l'écran de la tablette et de les transmettre au navigateur pour déclencher des actions locales.
Par exemple, un glissement vers la gauche pourrait activer une fonction spécifique dans le navigateur.

Mon approche consistait à capturer les gestes depuis l'appareil, puis à les communiquer au navigateur avec un message indiquant : "ce geste a été effectué".

Dans le domaine de l'échange de messages entre un serveur et un client dans un navigateur, les WebSockets viennent naturellement à l'esprit.
Les WebSockets sont intrinsèquement conçus pour prendre en charge des flux de messages sur TCP, contrairement à HTTP, qui gère principalement des flux d'octets sans concept intégré de message.

En naviguant à travers ce parcours, j'ai réalisé l'importance de tests approfondis et d'apprentissage pour élaborer une solution efficace.
Le protocole WebSocket, contrairement à HTTP, introduit des défis distincts, en particulier dans le débogage et les tests, en raison de sa nature plus complexe.

Reconnaissant que les gestes sont essentiellement un flux d'octets (je vais l'expliquer), je vais écrire sur :
- le processus d'évaluation du compromis entre la complexité ajoutée des WebSockets et les fonctionnalités qu'ils offrent.
- comment j'ai rationalisé le développement en concevant mon propre système de messagerie via HTTP.

## Contexte

J'ai déjà discuté de l'utilisation de ma tablette pour diverses présentations en direct.
Grâce à des tests itératifs, j'ai développé une solution hybride qui combine des éléments d'un tableau blanc avec des diapositives statiques.
Cette solution présente le dessin de l'écran en superposition sur des diapositives existantes.
Le défi réside maintenant dans le changement de diapositives directement depuis la tablette pour rationaliser la présentation et minimiser les interactions avec l'ordinateur portable affichant les diapositives.

Les diapositives sont affichées dans un iFrame côté client de mon outil.
Par conséquent, j'avais besoin d'une méthode pour envoyer des commandes à l'iFrame afin de contrôler les transitions de diapositives.
Le [framework de présentation reveal.js](https://revealjs.com/) prend en charge l'intégration native et permet le contrôle des diapositives depuis le cadre supérieur via une API qui utilise [postMessages](https://revealjs.com/postmessage/).

Pour transmettre les commandes de contrôle des diapositives de la tablette au client, j'ai envisagé diverses méthodes.
La solution optimale que j'ai identifiée était d'utiliser les gestes tactiles sur l'écran de la tablette reMarkable.
En glissant sur la tablette, je pourrais envoyer des événements au client, qui répondrait ensuite en conséquence pour changer les diapositives.

### Capturer les événements tactiles sur reMarkable/Linux

La reMarkable fonctionne sur un système basé sur Linux.
Les événements d'entrée (à la fois stylet et toucher) sont gérés via les [Périphériques d'événements (evdev)](https://en.wikipedia.org/wiki/Evdev).
L'exposition des événements est la suivante :
- `/dev/input/event1` capture les événements du stylet.
- `/dev/input/event2` capture les événements tactiles.

Dans Unix, la philosophie selon laquelle "_tout est fichier_" s'applique.
Cela signifie que je peux facilement accéder à ces événements en ouvrant et en lisant le contenu du fichier en Go.
J'ai choisi Go comme langage côté serveur en raison de son packaging autonome, de ses capacités de compilation croisée et du plaisir que j'en tire.

> "Tout est fichier" est un principe dans Unix et ses dérivés, où les interactions d'entrée/sortie avec des ressources telles que des documents, des disques durs, des modems, des claviers, des imprimantes et certaines communications inter-processus et réseau sont traitées comme de simples flux d'octets accessibles via l'espace de noms du système de fichiers - [source Wikipedia](https://en.wikipedia.org/wiki/Everything_is_a_file).

### Lire les événements en Go

Le "fichier" d'événement est un périphérique de caractères, offrant une représentation binaire d'un événement.
En Go, un ensemble d'octets d'événement pourrait être structuré comme ceci :

```go
type InputEvent struct {
	Time syscall.Timeval `json:"-"`
	Type uint16
	Code  uint16
	Value int32
}
```

Le principe selon lequel "_tout est fichier_" permet d'utiliser des opérations de base du package `os` pour ouvrir le périphérique de caractères en tant que `*os.File` et `Read` la représentation binaire de l'événement.
Nous créons un objet `ev` du type `InputEvent` pour recevoir les informations lues.

Le fichier fonctionne comme un `io.Reader`, et son contenu est généralement chargé dans un tableau d'octets.

```go 
func readEvent(inputDevice *os.File) (InputEvent, error) {
    // Size calculation: 
    // Timeval consists of two int64 (16 bytes), 
    // followed by uint16, uint16, and int32
    // (2+2+4 bytes)
    const size = 16 + 2 + 2 + 4
    eventBinary := make([]byte, size)

    _, err := inputDevice.Read(eventBinary)
    if err != nil {
        return InputEvent{}, err
    }

    var ev InputEvent
    // Assuming the binary data is in little-endian format 
    // which is the most common on Intel and ARM
    ev.Time.Sec = int64(binary.LittleEndian.Uint64(eventBinary[0:8]))
    ev.Time.Usec = int64(binary.LittleEndian.Uint64(eventBinary[8:16]))
    ev.Type = binary.LittleEndian.Uint16(eventBinary[16:18])
    ev.Code = binary.LittleEndian.Uint16(eventBinary[18:20])
    ev.Value = int32(binary.LittleEndian.Uint32(eventBinary[20:24]))

    return ev, nil
}
```
Une approche plus efficace pourrait impliquer l'utilisation d'un pointeur unsafe pour remplir directement la structure, contournant ainsi les mécanismes de sécurité de Go en utilisant le package `unsafe` :

```go
func readEvent(inputDevice *os.File) (events.InputEvent, error) {
	var ev InputEvent
    // by using (*[24]byte), we are explicitly stating that 
    // we want to treat the memory location of ev as a byte array of length 24
    // We could have used the less readable form:
    // (*(*[unsafe.Sizeof(ev)]byte)(unsafe.Pointer(&ev)))[:]
    // 
    //  Note: the trailing [:] is mandatory to convert the array to a slice
    _, err := inputDevice.Read((*[24]byte)(unsafe.Pointer(&ev))[:])
	return ev, err
}
```

## L'énoncé du problème

Maintenant que j'ai lu les événements, je dois les envoyer au client pour un traitement ultérieur.
L'architecture actuelle est basée sur un serveur HTTP en Go et un client web en JS. Par conséquent, je dois trouver un moyen HTTP-ish pour transférer les événements.

Il est au-delà de la portée de cet article d'entrer dans les détails de la façon dont je publie les événements au sein du serveur Go.
Cependant, pour une compréhension de base nécessaire pour le reste de l'article, voici un bref aperçu.

### Structure de service dans le serveur Go
Fondamentalement, j'ai implémenté un mécanisme [pubsub](https://github.com/owulveryck/goMarkableStream/blob/main/internal/pubsub/pubsub.go) de base pour canaliser le flux d'événements.

L'étape suivante consiste à rendre ces événements accessibles au client.
Cela sera géré par un `http.Handler`. Voici le cadre de ce gestionnaire :

```go
type GestureHandler struct {
    inputEventBus *pubsub.PubSub
}

// ServeHTTP implements http.Handler
func (h *GestureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // eventC is a channel that receives all the InputEvents
    eventC := h.inputEventBus.Subscribe("eventListener")
    // ....
}
```

## Le choix par défaut : WebSockets

Maintenant que je suis dans le gestionnaire HTTP, je dois concevoir une méthode pour transférer des données "sur le fil".
Dans ce contexte, "sur le fil" fait référence à deux flux d'octets :
- le `io.Reader` encapsulé dans le corps de la requête.
- le `io.Writer` implémenté via le ResponseWriter.

La méthode la plus familière que je connaisse pour échanger des messages entre un serveur et un client est via WebSocket.
WebSocket est un protocole de couche 7 qui permet des flux bidirectionnels de messages.
Son implémentation est relativement simple, et le côté client en JavaScript fournit toutes les primitives nécessaires pour interagir avec les flux de messages.

Du côté serveur, la situation diffère, car la bibliothèque standard de Go n'inclut pas d'implémentation de WebSockets.
Cela nécessite de s'appuyer sur des bibliothèques tierces.
Bien que ce ne soit pas intrinsèquement problématique, je préfère généralement éviter les bibliothèques tierces en raison de préoccupations concernant les éléments boîte noire et les complexités de gestion des dépendances.

Néanmoins, j'ai implémenté un échange de messages de base basé sur WebSocket pour envoyer des événements du serveur au client.

Ayant établi la capacité d'écouter les événements et de les servir via WebSockets, l'étape suivante consistait à détecter avec précision un geste avant d'envoyer l'événement.
J'ai incorporé une logique métier de base dans mon gestionnaire, en utilisant une minuterie pour identifier les mouvements continus.
Cela m'a permis de transmettre le mouvement en termes de distance parcourue par le doigt, comme 100 pixels à gauche, 130 pixels à droite, 245 pixels en haut et 234 pixels en bas.
Bien qu'il s'agisse d'une implémentation simpliste qui ne fait pas la différence entre un carré et un cercle, elle suffit à mes besoins.

Cependant, tester cette implémentation posait un défi significatif.
Étant dans la phase exploratoire du développement du produit, la stratégie la plus efficace était d'adopter une approche 'test et apprentissage', plutôt que d'établir une suite de tests complète.
Cette approche est susceptible d'évoluer à mesure que le produit mûrit, mais pour le moment, il était nécessaire de "rétro-concevoir" le flux pour comprendre les types d'événements générés par des interactions spécifiques avec l'écran.

_Note_ : La théorie de l'évolution de Simon Wardley a considérablement influencé mon approche de ce projet.
Pour une compréhension plus approfondie de cette théorie, je recommande de consulter la littérature pertinente ou de me contacter pour une discussion plus approfondie.

C'est là que réside une limitation des WebSockets : ils sont distincts du protocole HTTP, ce qui signifie que des outils comme cURL ou netcat ne peuvent pas être utilisés pour se connecter au point de terminaison et surveiller les messages.
Bien qu'il existe des outils disponibles à cette fin, ils manquent souvent de certaines fonctionnalités, comme la confiance pour un certificat auto-signé.

J'ai passé un temps considérable à essayer de comprendre comment diffuser des messages à l'écran tout en déplaçant mon doigt sur la tablette.
J'ai réalisé que l'apprentissage des subtilités des outils WebSocket pourrait ne pas être l'utilisation la plus efficace de mon énergie, surtout lorsque je recherche des résultats rapides pour la fonctionnalité de geste.

## Une approche alternative : les flux HTTP

S'en tenir à un échange HTTP pur pourrait être une option plus appropriée. Revenons en arrière pour analyser le parcours jusqu'à présent :

- Les événements tactiles sont sérialisés par le noyau Linux et exposés sous forme de flux d'octets via un fichier `/dev/input/event`.
- Ce flux est disséqué en une série d'événements discrets, qui sont ensuite introduits dans un canal.
- Ces événements sont analysés pour détecter un "geste" – une séquence d'événements correspondant au même "toucher".
- Les événements agrégés et assainis sont ensuite transmis au client à l'aide de WebSocket.

Considérant que les événements initiaux sont présentés comme un flux d'octets, et voyant l'efficacité d'avoir le client lire et segmenter ces événements, cela s'aligne bien avec la philosophie Unix.

Par conséquent, j'ai décidé d'explorer une implémentation de flux de bas niveau pour la communication entre le client et le serveur.

Internet et ChatGPT lui ont donné un nom : [Server Sent Events](https://en.wikipedia.org/wiki/Server-sent_events)

{{< figure src="/assets/websockets-sequence.png" link="/assets/websockets-sequence.png" title="Diagramme de séquence" >}}

Du point de vue du serveur, le processus implique la diffusion continue d'octets dans le canal de communication.
Ces octets sont formatés spécifiquement pour annoncer des événements.
Un type MIME spécial (`text/event-stream`) est utilisé pour signaler au client que le serveur enverra un tel flux d'octets, et le client est censé le gérer en conséquence.

Initialement, j'ai envisagé d'implémenter les événements envoyés par le serveur (SSE), mais j'ai ensuite réalisé que je pourrais d'abord explorer une approche plus simple.
Cela implique de diffuser des octets sans adopter pleinement la logique complète des SSE, d'autant plus que je gère à la fois les implémentations client et serveur.
Cette approche permet un processus de développement plus rationalisé et contrôlé.

### Implémentation du flux HTTP en Go

L'implémentation d'un flux d'octets dans un point de terminaison est assez simple en Go.
Le gestionnaire est fourni avec un `ResponseWriter`, qui est un [`io.Writer`](https://pkg.go.dev/io#Writer).
Cela signifie qu'il suffit d'invoquer la méthode `Write` dans une boucle sans fin pour la tâche en question.

L'aspect crucial est de s'assurer que le flux est alimenté avec la charge utile correcte, à savoir la tranche d'octets appropriée.

### Sérialisation du message

Le concept de [sérialisation](https://en.wikipedia.org/wiki/Serialization) est :

> le processus de traduction d'une structure de données ou d'un état d'objet dans un format qui peut être stocké (par exemple.
fichiers dans des dispositifs de stockage secondaire, tampons de données dans des dispositifs de stockage primaire) ou transmis (par exemple.
flux de données sur des réseaux informatiques) et reconstitué ultérieurement / source Wikipedia

Il est donc nécessaire de "sérialiser" les messages de geste en un tableau d'octets d'une manière qui permette de les désérialiser côté client.
Comme le client est un programme basé sur Javascript, j'utiliserai JSON.

Ainsi, le geste est implémenté comme une structure qui implémente l'interface JSON Marshaler.

```go
type gesture struct {
        leftDistance, rightDistance, upDistance, downDistance int
}

func (g *gesture) MarshalJSON() ([]byte, error) {
        return []byte(fmt.Sprintf(`{ "left": %v, "right": %v, "up": %v, "down": %v}`+"\n", g.leftDistance, g.rightDistance, g.upDistance, g.downDistance)), nil
}
```

Ce que nous avons maintenant est une collection d'événements qui sont agrégés dans une structure `gesture` et sérialisés en format binaire pour transmission au client.
Nous avons mis en place un point de terminaison `/gestures` pour servir continuellement ce flux de données de gestes.

### Réception et décodage du flux en JavaScript

Du côté client, nous récupérons les données en JavaScript, en utilisant un thread worker pour récupérer et analyser les gestes.

Le worker reçoit un ensemble de mouvements (une structure `gesture` sérialisée) et les interprète en commandes de plus haut niveau, comme une action "swipe left".

```js
const gestureWorker = new Worker('worker_gesture_processing.js');

gestureWorker.onmessage = (event) => {
    const data = event.data;
    switch (data.type) {
        case 'gesture':
            switch (data.value) {
                case 'left':
                    // Send the order to switch slide to the iFrame
                    document.getElementById('content').contentWindow.postMessage(JSON.stringify({ method: 'left' }), '*');
                    break;
                // ...
```

Dans le thread worker, nous utilisons la méthode `fetch` pour obtenir les données du point de terminaison `/gestures`. 
Nous créons ensuite un `reader` et bouclons continuellement pour lire les données entrantes.

```js
const response = await fetch('/gestures');

const reader = response.body.getReader();
const decoder = new TextDecoder('utf-8');
let buffer = '';

while (true) {
    const { value, done } = await reader.read();
    //...
    buffer += decoder.decode(value, { stream: true });

    while (buffer.includes('\n')) {
        const index = buffer.indexOf('\n');
        const jsonStr = buffer.slice(0, index);
        buffer = buffer.slice(index + 1);

        try {
            const json = JSON.parse(jsonStr);
            let swipe = checkSwipeDirection(json);
            //...
        }
//...
```

La fonction `checkSwipeDirection` analyse les données JSON, identifiant les gestes de balayage et les transmettant comme actions appropriées.

Avec cette configuration, nous avons maintenant un mécanisme complet en place pour capturer les événements, détecter les gestes de balayage et initier les actions correspondantes.

C'est tout, les amis !

## Conclusion

En conclusion, le parcours de développement pour améliorer mon outil, goMarkableStream, a été un témoignage vivant de l'adage "le simple est complexe", soulignant la valeur inhérente à embrasser la simplicité.
Bien que l'attrait des frameworks et des protocoles sophistiqués soit indéniable, ce projet illustre qu'ils ne sont pas toujours le choix optimal pour des tâches simples.
En s'en tenant aux principes de base de la philosophie Unix, où chaque interaction est traitée comme un flux d'octets, j'ai pu concevoir une solution à la fois efficace et élégante dans sa simplicité.

Dans ce voyage, j'ai également présenté ma décision de lire et de traiter les événements directement en utilisant les outils Go prêts à l'emploi, sans utiliser de bibliothèques tierces.
Conformément à la sagesse de Rob Pike selon laquelle "_un peu de copie vaut mieux qu'un peu de dépendance_",
cette méthode a non seulement assuré un processus de développement plus rationalisé, mais m'a également accordé une compréhension et un contrôle plus profonds de la fonctionnalité que je construisais.

En fin de compte, cette expérience a été une célébration de la maîtrise des octets et des joies de l'artisanat logiciel pratique.
Elle sert de rappel que parfois, les meilleures solutions ne proviennent pas de la complexité et de la sophistication des outils que nous utilisons,
mais de notre capacité à dépouiller un problème jusqu'à ses fondamentaux et à l'aborder de front.
La vieille philosophie Unix, souvent négligée, recèle encore un trésor de sagesse pour les développeurs modernes, préconisant la simplicité, la clarté et le plaisir inhérent à la manipulation directe des octets.