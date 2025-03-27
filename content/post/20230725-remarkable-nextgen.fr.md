---
title: "Faire évoluer le jeu : Un outil de streaming sans client pour reMarkable 2"
date: 2023-07-25T15:55:21+02:00
lastmod: 2023-07-25T15:55:21+02:00
draft: false
videos: [/assets/goMarkableStream2.webp]
images: [/assets/goMarkableStream2.png]
keywords: []
summary: Dans cet article, je présente la nouvelle version de l'outil de streaming pour ma tablette reMarkable. 

  
  Conçu à l'origine en 2021, cet outil me permettait de diffuser des croquis pendant les appels vidéo. 


  Visant une plus grande convivialité, j'ai repensé la conception pour une implémentation sans client.
  Cet article décrit l'implémentation avec quelques illustrations de code en Javascript et Go sur comment

  * récupérer l'image et l'afficher dans un canvas

  * optimiser le flux en jouant avec `uint4` et `RLE`
tags: ["golang", "reMarkable", "JS", "Optimization", "RLE", "ChatGPT", "hack", "WebSocket"]
categories: ["dev"]
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

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/20230725-remarkable-nextgen/).

En 2021, j'ai développé un outil pour diffuser le contenu de ma reMarkable 
(et j'ai également écrit un article à ce sujet [ici](https://blog.owulveryck.info/2021/03/30/streaming-the-remarkable-2.html)).
Étant donné que je travaillais principalement de chez moi, cet outil était très utile, me permettant d'esquisser des éléments pendant les appels en conférence.

L'un des principaux avantages de cet outil était sa capacité à diffuser du contenu directement dans un onglet de navigateur web. 
Cette fonctionnalité était particulièrement utile car elle signifiait que je pouvais partager exclusivement cet onglet pendant les appels vidéo, assurant ainsi une concentration sur le contenu que je souhaitais présenter.

À la base, l'outil était composé de deux composants principaux :
1. Un serveur fonctionnant sur l'appareil, responsable de la capture de l'image brute et de sa transmission au portable.
2. Un service sur le portable, qui récupérait l'image brute du serveur et la traitait dans un format adapté à la visualisation dans le navigateur (j'ai produit un flux MJPEG pour cela).

Être le chef de produit de mes propres outils offrait une perspective unique. 
Un retour d'expérience que j'ai fourni était la nature légèrement encombrante de l'activation à la volée de l'outil. 
Le défi provenait de la nécessité d'un service local. Pendant les appels vidéo, cela signifiait quelques étapes préparatoires pour initier le service local, ajoutant une couche supplémentaire de complexité au processus.

Reconnaissant ce point douloureux, mon objectif est devenu clair : éliminer le besoin du service local. 
Cet article se penche sur l'implémentation remaniée de l'outil de streaming, qui présente désormais une conception plus conviviale et des performances améliorées.

## Faire fonctionner : L'évolution de l'ancien au nouveau

Le code fonctionnant sur l'appareil doit avoir une empreinte légère. 
Une façon de s'assurer qu'il reste léger est d'éliminer tout calcul lourd sur l'appareil. 
La seule fonction du code fonctionnant sur le serveur est de saisir l'image brute de la mémoire et de l'exposer sur le réseau. 
Cela a conduit à une conception à trois niveaux : serveur/client/rendeur.

_Langage ubiquitaire :_

Dans cet article :

- **Le serveur** fait référence au code fonctionnant sur la reMarkable (l'appareil). Son objectif principal est d'exposer l'image brute de l'affichage actuel sur la reMarkable.
- **Le client** est responsable de la récupération de l'image brute du serveur et de l'exécution de processus supplémentaires pour la convertir dans un format utilisable.
- **Le rendeur** accepte la sortie du client et l'affiche sur un écran PC.

```plain
+---------------------------+        +-----------------------------------+
|          reMarkable       |        |               Laptop              |
|                           |        |                                   |
|       +-------+           |  gRPC  |       +-------+                   |
|       |Server |<--------- |<------>|------>|Client |                   |
|       +-------+           | Fetch  |       +-------+                   |
|                           | Stream |           |                       | 
|                           |        |     HTTP MJPEG stream             |
|                           |        |           |                       |
|                           |        |           v                       |
|                           |        |       +--------+                  |
|                           |        |       |Renderer|                  |
|                           |        |       +--------+                  |
|                           |        |      (Browser/VLC)                |
+---------------------------+        +-----------------------------------+
```

Pour minimiser l'utilisation du CPU, le serveur extrait l'image uniquement lorsque le client est connecté.
Cette fonctionnalité a été réalisée grâce à la communication gRPC.

Le serveur peut ainsi fonctionner comme un démon sur la reMarkable, attendant un appel RPC du client.
Pour initier le streaming, j'avais simplement besoin d'activer le client.
Le client récupère l'image en boucle, et chaque image est encodée en JPEG avant d'être ajoutée à un flux MJPEG.
Ce flux est ensuite rendu disponible comme service HTTP par le client.

Le rendeur est n'importe quel logiciel capable de lire le flux MJPEG via HTTP, comme VLC ou un navigateur web.

Un défi avec cette configuration est qu'elle nécessite une topologie et une configuration réseau spécifiques.
Le client doit non seulement être au courant de l'adresse de la reMarkable, mais aussi posséder les privilèges nécessaires pour établir un serveur.
De plus, le rendeur doit connaître l'adresse IP du client.

Bien que cela n'ait pas été un problème dans ma configuration personnelle, des complications sont apparues après la pandémie lorsque je suis revenu aux présentations en personne.
J'ai réalisé le besoin d'une solution plus simple.
Mon objectif ultime est devenu la possibilité de simplement saisir l'adresse de la reMarkable dans n'importe quel navigateur et d'accéder instantanément au flux.

### Nouvelle Architecture

Pour atteindre l'objectif, la solution implique d'éliminer le client et d'établir plutôt un serveur HTTP au sein du composant serveur.
Le client devrait être implémenté dans un format interprétable par un navigateur, comme Javascript ou WASM.

Ma première approche était de compiler le client en WASM.
Cela semblait prometteur car cela me permettrait d'exploiter mon expertise en développement Go.
Cependant, j'ai rencontré plusieurs limitations qui auraient nécessité des modifications substantielles.

Par conséquent, j'ai choisi de développer une deuxième version de l'outil, avec le client écrit en Javascript.

_Note annexe :_ À ce stade, j'ai été confronté à un autre défi.
Bien que je possède une compréhension large du fonctionnement de Javascript et des processus de rendu du navigateur (ce que nous pourrions appeler des "compétences architecturales"), 
je me sentais moins confiant dans mes capacités pratiques de développement JS.
Je me suis tourné vers mon assistant numérique, ChatGPT, pour obtenir des conseils.
Avec ma direction sur la solution souhaitée, il a fourni les fragments de code nécessaires et les explications pour donner vie à ma vision.
J'étais le développeur, il était le codeur.

### Validation du rendeur "canvas"

Initialement, il était impératif de s'éloigner du flux MJPEG, d'autant plus que mes opérations étaient désormais étroitement alignées avec le rendeur, 
et Javascript possède les primitives requises pour la manipulation d'images.

Dans le navigateur, la méthode conventionnelle pour gérer les images est via l'élément `canvas`.
Ma tâche préliminaire était de valider que je pouvais récupérer une image brute du serveur et la présenter dans un `canvas`.

J'ai réalisé cela en accédant à la colonne vertébrale du canvas qui représente les données de la carte de pixels au format RGBA et en ajustant les pixels en fonction de leurs valeurs dans l'image brute de la reMarkable :

```js
<canvas id="fixedCanvas" width="1872" height="1404" class="hidden"></canvas>
<script>
    // Use the fixed-size canvas context to draw on the canvas
    var fixedCanvas = document.getElementById("fixedCanvas");
    var fixedContext = fixedCanvas.getContext("2d");
    function processBinaryData(data) {

        // Assuming each pixel is represented by 4 bytes (RGBA)
        var pixels = new Uint8Array(data);
        // Create an ImageData object with the byte array length
        var imageData = fixedContext.createImageData(fixedCanvas.width, fixedCanvas.height);
        // Assign the byte array values to the ImageData data property
        for (var i = 0; i < pixels.length; i++) {
            imageData.data[i*4] = pixels[i];
            imageData.data[i*4+1] = pixels[i];
            imageData.data[i*4+2] = pixels[i];
            imageData.data[i*4+3] = 255;
        }

        // Display the ImageData on the canvas
        fixedContext.putImageData(imageData, 0, 0);
    }
```

Il y a également une exigence d'ajuster l'image pour la rendre responsive selon la taille du navigateur, ainsi que pour la rotation de l'image et la colorisation potentielle.

Pour y parvenir, je maintiens le `fixedCanvas` dans un état caché et transfère son contenu vers un autre canvas en utilisant la méthode `drawImage`.
Les dimensions du canvas de destination (sa largeur et sa hauteur) subissent des ajustements si un événement de redimensionnement est détecté dans la fenêtre du navigateur.

{{< highlight js >}}
var resizableCanvas = document.getElementById("canvas");
var resizableContext = resizableCanvas.getContext("2d");
function copyCanvasContent() {
    resizableContext.drawImage(fixedCanvas, 0, 0, resizableCanvas.width, resizableCanvas.height);
}

// JavaScript code for working with the canvas element
function resizeCanvas() {
    var canvas = document.getElementById("canvas");
    var container = document.getElementById("container");

    var aspectRatio = 1872 / 1404;

    var containerWidth = container.offsetWidth;
    var containerHeight = container.offsetHeight;

    var containerAspectRatio = containerWidth / containerHeight;

    if (containerAspectRatio > aspectRatio) {
        canvas.style.width = containerHeight * aspectRatio + "px";
        canvas.style.height = containerHeight + "px";
    } else {
        canvas.style.width = containerWidth + "px";
        canvas.style.height = containerWidth / aspectRatio + "px";
    }

    // Use the canvas context to draw on the canvas
    copyCanvasContent();
}

// Resize the canvas whenever the window is resized
window.addEventListener("resize", resizeCanvas);
{{< / highlight >}}

### Remplacement de base

Avec le rendeur en place, ma prochaine étape était d'implémenter un client JS léger en remplacement.
Bien que gRPC offre une fonctionnalité robuste, il n'est généralement pas considéré comme la référence en matière de développement web.

Ainsi, pour la communication et l'encapsulation, j'ai gravité vers le protocole WebSocket.
Il était suffisamment simple pour être incorporé côté serveur, servant de remplacement transparent pour le serveur RPC.

Les messages délivrés via ce protocole transportent l'image brute.
Le client reste en veille pour ces messages, et avec chaque message entrant, le contenu du canvas est mis à jour, émulant efficacement un processus de streaming.

Un avantage notable de cette approche est le contrôle accru sur la charge côté serveur.
L'extraction de l'image brute demande à la fois des ressources mémoire et CPU sur l'appareil.
En régulant la fréquence d'émission des messages, je peux gérer efficacement la charge de l'appareil.

## Faire correctement : Changer l'architecture de streaming

La solution basée sur les Websockets était opérationnelle, mais elle introduisait des défis, particulièrement sur iOS.
De plus, l'implémentation du Websocket côté serveur introduisait une certaine surcharge, sur laquelle je manquais de contrôle.
En conséquence, j'ai poursuivi une stratégie différente pour éliminer l'utilisation des websockets.

> On pourrait se demander : Pourquoi ai-je même besoin d'une méthode d'encapsulation ? Ne puis-je pas simplement envoyer directement le flux de données ?

En effet, _la simplicité est complexe_.


Soulever cette question de simplicité m'a fait passer à une approche rudimentaire : transmettre des images brutes sur le réseau sans encapsulation supplémentaire.

C'était faisable car je connaissais la taille de l'image, qui reste constante en raison de la résolution de la reMarkable.

J'ai conçu un point de terminaison Go qui écrivait continuellement des images sur le fil (spécifiquement sur le `http.ResponseWriter`), en utilisant une méthode `Write` basique.

```go 
func (h *StreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    //ctx, cancel := context.WithTimeout(r.Context(), 1*time.Hour)
    ctx, cancel := context.WithTimeout(r.Context(), 1*time.Hour)
    defer cancel()
    tick := time.NewTicker(rate * time.Millisecond),
    defer tick.Stop()

    imageData := imagePool.Get().([]uint8)
    defer imagePool.Put(imageData) // Return the slice to the pool when done
    for {
        select {
        case <-ctx.Done():
            return
        case <-tick.C:
            // Read the content of the image
            _, err := file.ReadAt(imageData, h.pointerAddr)
            if err != nil {
                log.Fatal(err)
            }
            // Write the image
            w.Write(imageData)
        }
    }
}
```

Du point de vue du client, une méthode `fetch` capture les données et alimente la colonne vertébrale du canvas.

```js 
// Create a new ReadableStream instance from a fetch request
const response = await fetch('/stream');
const stream = response.body;

// Create a reader for the ReadableStream
const reader = stream.getReader();
// Create an ImageData object with the byte array length
var imageData = fixedContext.createImageData(fixedCanvas.width, fixedCanvas.height);

// Define a function to process the chunks of data as they arrive
const processData = async ({ done, value }) => {
        // Process the received data chunk
        // Assuming each pixel is represented by 4 bytes (RGBA)
        var uint8Array = new Uint8Array(value);
        for (let i = 0; i < uint8Array.length; i++) {
                // process data to feed the backbone of the canvas (imageData)
                // ...
                copyCanvasContent(); // copy the hidden canvas to the responsive one
            }
        }

        // Read the next chunk
        const nextChunk = await reader.read();
        processData(nextChunk);
};

// Start reading the initial chunk of data
const initialChunk = await reader.read();
processData(initialChunk);
```


## Rendre rapide : Optimisations de la consommation réseau

Avec une architecture robuste en place, il est temps d'affiner l'efficacité de l'outil. Un défi significatif provient de la taille de l'image brute—environ 2,5 Mo (avec une résolution de 1872x1404 pour la reMarkable 2). Ce volume de données doit être transféré avec chaque image.

### Empaquetage des valeurs

La reMarkable affiche "16" couleurs distinctes. Depuis la sortie du FW 3.3, cette palette de couleurs peut être représentée comme un tableau de uint4 au lieu de uint8, comme discuté dans [issue 36](https://github.com/owulveryck/goMarkableStream/issues/36). L'adoption de cette représentation peut donner une réduction de 50% du volume de données.

Cependant, Go et JavaScript manquent de support natif pour le type uint4.

Une solution viable implique de stocker les données pour deux pixels dans un seul octet (uint8). Cette approche nécessite la création de deux fonctions dédiées—une pour l'empaquetage en Go et une autre pour le désempaquetage en JavaScript.

**Fonction d'empaquetage Go :**
```go 
// Packing algorithm to encode two uint4 values into a single uint8 
// Assumes arguments as uint4 and omits verification for efficiency
func pack(value1, value2 uint8) uint8 {
        // Shift the first value by 4 bits and combine it with the second using a bitwise OR
        encodedValue := (value1 << 4) | value2;
        return encodedValue;
}
```

**Fonction de désempaquetage JavaScript :**
```js
// Unpack the uint4 values
function unpackValues(packedValue) {
        // Extract the upper 4 bits to obtain the first value
        const value1 = (packedValue >> 4) & 0x0F;

        // Isolate the lower 4 bits to get the second value
        const value2 = packedValue & 0x0F;

        return [value1, value2];
}
```

### Compression RLE pour une efficacité accrue

Ayant réalisé une réduction de 1,2 Mo par image grâce à notre technique d'empaquetage, l'étape suivante consiste à minimiser davantage le transfert de données.
Pour y parvenir, nous pouvons nous tourner vers des algorithmes de compression plus sophistiqués sans solliciter le CPU de la reMarkable ou compliquer notre implémentation.

Après consultation avec des pairs, l'algorithme [Run Length Encoding (RLE)](https://en.wikipedia.org/wiki/Run-length_encoding) a émergé comme une option recommandée en raison de sa simplicité et de son efficacité.
Sans entrer dans une explication détaillée, le principe derrière RLE est relativement simple : il s'agit de comptabiliser les occurrences consécutives de la même valeur de pixel et ensuite de transmettre ce comptage aux côtés de la valeur de pixel elle-même.

Par exemple, considérons une séquence d'échantillon :

```text
0 0 0 0 0 0 1 1 1 0 0 0 0
```

En utilisant RLE, cette séquence se transforme en :

```text
6 0 3 1 4 0
```

L'implémentation de RLE est assez directe.
Cependant, un défi se pose lorsqu'on considère les valeurs de comptage potentielles, qui peuvent s'élever jusqu'à 1872*1404.
Représenter de tels grands nombres nécessiterait un type de données comme uint64.
Cela pose un risque : dans certains scénarios, la séquence "compressée" pourrait finir par dépasser l'image non compressée en taille.

Pour éviter cela, j'ai choisi de plafonner la longueur de comptage à 15.
Ce choix ouvre la voie pour représenter à la fois le comptage et la valeur de pixel dans un seul octet, trouvant un équilibre entre simplicité et efficacité.

Un avantage supplémentaire de notre implémentation RLE en Go, qui imite un `io.Writer`, est sa réutilisabilité.
Si une situation l'exige, je peux appliquer la compression RLE deux fois, bien que les circonstances actuelles n'aient pas exigé une telle mesure.

Jusqu'à présent, le transfert est d'environ 200 Ko en moyenne.

### Envoi d'images uniquement lors de modifications

L'optimisation finale consiste à transmettre de nouvelles images uniquement lorsqu'il y a un changement.
Déterminer si une image a été modifiée nécessiterait généralement le calcul d'une somme de contrôle, ce qui peut être intensif en CPU.

Cependant, l'appareil reMarkable fonctionne sur un système Linux.
Ainsi, toute interaction avec l'écran, que ce soit par stylet ou tactile, est acheminée via `/dev/input/event*`.
J'ai introduit une goroutine qui surveille ces événements et envoie les images selon les besoins.

En conséquence, en l'absence d'événements, l'utilisation du CPU tombe à zéro, même si un client reste connecté.
Pendant les opérations d'écriture, l'utilisation du CPU oscille autour de 10 % — un niveau que je considère comme efficace.

## Notes finales

Cette application est basée sur un hack.
Le défi principal consiste à découpler efficacement l'interface, qui récupère l'image, du client/rendeur.

Dans l'implémentation précédente, il y avait un découplage complet facilité par la définition protobuf entre le client et le serveur.
Historiquement, lorsque reMarkable a introduit son firmware 3.3, cela a perturbé l'outil, comme souligné par [issue 36](https://github.com/owulveryck/goMarkableStream/issues/36) sur GitHub.
Cependant, les ajustements nécessaires pour corriger ce problème n'ont touché que le composant client.

Il semble que la version 3.6 du firmware pourrait également introduire un changement radical, comme indiqué par [issue 58](https://github.com/owulveryck/goMarkableStream/issues/58).
Je prévois que résoudre cela impliquera des modifications plus larges.
En revanche, la nature autonome de l'application (avec le client intégré au serveur) devrait simplifier les mises à jour sur l'appareil.

Le domaine de l'IT est truffé de compromis ; il n'y a pas de solution universelle.
Ce dynamisme et cette adaptabilité sont ce qui insuffle de l'excitation dans le domaine.

L'application et son code source sont accessibles sur [github.com/owulveryck/goMarkableStream](https://github.com/owulveryck/goMarkableStream).

![video](/assets/goMarkableStream2.webp)