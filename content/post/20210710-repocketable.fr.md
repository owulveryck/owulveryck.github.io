---
title: "Lire du contenu web hors ligne et sans distraction"
date: 2021-10-07T10:07:41+02:00
draft: false
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/20210710-repocketable/).

*TL;DR:* Cet article décrit la conception d'un outil pour transformer une page web en un epub autonome (pour lire hors ligne). Si vous voulez essayer l'outil, vous pouvez télécharger une version binaire depuis [GitHub](https://github.com/owulveryck/rePocketable/tags)

## Le Pourquoi

Pour simplifier mon besoin, je vais citer le _[Readability Project](https://web.archive.org/web/20150817073201/http://lab.arc90.com/2009/03/02/readability/)_

> Lire quoi que ce soit sur Internet est devenu un véritable cauchemar. Alors que les médias tentent de soutirer autant de revenus publicitaires que possible, nous essayons de nous mettre des œillères pour masquer toute la folie qui entoure le contenu que nous essayons de lire.
>
> C'est presque comme écouter une émission de radio, sauf que les publicités passent pendant le programme en arrière-plan. C'est une expérience assez horrible. Notre ami jusqu'à présent a été le fidèle bouton "Vue d'impression". Cliquez dessus et toutes les choses inutiles disparaissent. Je clique dessus tout le temps et j'imprime rarement. C'est vraiment devenu le bouton "Paix et Tranquillité" pour beaucoup.

### Cet article

Dans un post récent, j'ai blogué à propos d'un outil que je construis pour ma reMarkable.
Dans ce post, je vais décrire un nouvel outil qui convertit n'importe quelle page web en un fichier ePub.

Les objectifs de cet outil sont :

- garder une trace des articles que j'aime sans craindre les liens brisés
- extraire le contenu et lire les articles sans distraction
- pouvoir lire les articles hors ligne sur des appareils tels que les liseuses ou ma reMarkable

## Solutions existantes

Cette fonctionnalité existe déjà si vous utilisez une Kobo et le service getPocket.
Le problème est que l'expérience hors ligne est étroitement liée à mon appareil Kobo.
De plus, getPocket n'offre aucun moyen de télécharger la version nettoyée des articles.

Nous, en tant que développeurs, avons des super-pouvoirs : nous pouvons construire les outils que nous voulons.

Expliquons les fonctionnalités que je construis étape par étape.

_Avertissement_ au moment où cet article est écrit, l'outil résulte de diverses expériences, mais ni l'architecture ni le code n'est propre et maintenable.
Prenez ce post comme une validation d'une preuve de concept.

## Première partie : extraction du contenu

La partie la plus importante de ce voyage est la capacité de l'outil à extraire le contenu d'une page web.
La première idée serait d'interroger le service getPocket qui fait cela, mais la [documentation de leur API](https://getpocket.com/developer/docs/v3/article-view) mentionne que :

> L'API Article View de Pocket renverra le contenu de l'article et les métadonnées pertinentes sur toute URL fournie.
>
> L'API Article View de Pocket n'est actuellement ouverte qu'aux partenaires qui intègrent des fonctionnalités spécifiques à Pocket ou des clients Pocket à part entière. Par exemple, la création d'un client Pocket pour la plateforme X.
>
> Si vous recherchez un analyseur de texte général ou pour fournir une fonctionnalité "lire maintenant" dans votre application - nous ne le prenons pas en charge actuellement. Il existe d'autres entreprises/produits qui fournissent ce type d'API, par exemple : Diffbot.

Ils mentionnent [Diffbot](https://www.diffbot.com/products/extract/), mais c'est un service web qui nécessite un abonnement ; j'aimerais construire un outil simple, gratuit, pour mon usage, et donc ce n'est pas une option.

### Readability / Arc90

J'ai regardé les initiatives open source qui alimentent les modes de lecture des navigateurs (je suis/j'étais fan du mode lecture de safari), et j'ai trouvé que certaines d'entre elles étaient basées sur une expérience faite par [Arc90](https://web.archive.org/web/20150817073201/http://lab.arc90.com/2009/03/02/readability/).
Cette expérience a conduit au service (abandonné) [readability](https://en.wikipedia.org/wiki/Readability_(service)).

Nous pouvons maintenant trouver diverses implémentations de l'[algorithme Arc90](https://github.com/masukomi/arc90-readability). J'utilise [cette implémentation](https://github.com/cixtor/readability) en Go pour mon outil.

#### Code

N'hésitez pas à sauter cette partie si vous n'êtes pas intéressé par le code

L'API de la bibliothèque readability est simple.
Tout d'abord, il faut créer un objet `Readability` avec un _analyseur HTML qui lit et extrait le contenu pertinent_.

Ensuite, appeler la méthode `Parse` sur cet objet, en lui fournissant un `io.Reader` qui contient la page à analyser.

Le résultat est un objet de type `Article` qui contient des métadonnées et le contenu nettoyé. Ce contenu est un arbre HTML et est accessible via un [`html.Node`](https://pkg.go.dev/golang.org/x/net/html#Node) de haut niveau.

{{< highlight go >}}
package main

import (
   "log"
   "net/http"
   "os"

   "github.com/cixtor/readability"
   "golang.org/x/net/html"
)

func main() {
   // create a parser
   htmlParser := readability.New()
   // Fetch a webpage
   resp, err := http.Get("https://example.com/")
   passOrDie(err)
   // Deal with errors etc...
   defer resp.Body.Close()
   // Parse the content
   article, err := htmlParser.Parse(resp.Body, "https://example.com")
   passOrDie(err)
   // Write the readable result on stdout
   html.Render(os.Stdout, article.Node)
}
{{< /highlight >}}

### Le problème avec le contenu réactif et les articles Medium

Lorsque le projet Arc90 a fait cette expérience, il n'y avait pas beaucoup de contenus réactifs.

De plus, il ne gère pas le javascript.
Cela conduit à des images qui ne s'affichent pas correctement. Prenons le [premier chapitre du livre de Simon Wardley sur les cartes](https://medium.com/wardleymaps/on-being-lost-2ef5f05eb1ec) pour illustrer le problème.

L'image ci-dessous est une capture d'écran d'une vue de lecteur de la page avec Safari :
{{< figure src="/assets/rePocketable/medium.png" title="Le problème de Medium avec Arc90" >}}

Le code ci-dessous est le code extrait par une requête curl :
{{< highlight html >}}
<figure
   class="ja jb jc jd je jf cw cx paragraph-image">
   <div role="button" tabindex="0"
      class="jg jh ji jj aj jk">
      <div class="cw cx iz">
         <div class="jq s ji jr">
            <div
               class="js jt s">
               <div
                  class="jl jm t u v jn aj at jo jp">
                  <img alt=""
                     class="t u v jn aj ju jv jw"
                     src="https://miro.medium.com/max/60/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg?q=20"
                     width="700"
                     height="590"
                     role="presentation" />
               </div>
               <img alt=""
                  class="jl jm t u v jn aj c"
                  width="700"
                  height="590"
                  role="presentation" /><noscript><img
                     alt=""
                     class="t u v jn aj"
                     src="https://miro.medium.com/max/1400/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg"
                     width="700"
                     height="590"
                     srcSet="https://miro.medium.com/max/552/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg 276w, https://miro.medium.com/max/1104/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg 552w, https://miro.medium.com/max/1280/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg 640w, https://miro.medium.com/max/1400/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg 700w"
                     sizes="700px"
                     role="presentation" /></noscript>
            </div>
         </div>
      </div>
   </div>
</figure>
{{< /highlight >}}

À l'intérieur de l'élément `<figure>`, nous pouvons voir que la première image (https://miro.medium.com/max/60/1*RSH2vh_xgQtjB68Zb7oBaA.jpeg?q=20) est une miniature et elle agit comme un espace réservé.

Un ensemble de routines JavaScript remplace l'image au moment du rendu dans le navigateur.
Heureusement, une balise `<noscript>` est également présente et expose les sources complètes de l'image.

Comme la bibliothèque Arc90 supprime tous les éléments `<noscript>`, les seules options sont :

- prétraiter le fichier HTML avant d'alimenter l'algorithme Arc90
- modifier la bibliothèque Arc90

Jusqu'à présent, le comportement que nous traitons semble particulier aux articles hébergés sur medium. Modifier l'algorithme Arc90 pour gérer ce cas d'utilisation spécifique ne semble pas être une bonne idée.

Optons donc pour une étape de prétraitement du document avant d'alimenter l'algorithme Arc90.
Il dépasse le cadre de cet article de montrer et commenter le code complet pour cela.

En un coup d'œil, le contenu HTML est extrait dans un arbre d'éléments `*html.Node` ; ensuite, l'étape de traitement parcourt l'arbre via une fonction récursive cherchant des éléments `figure`.

{{< highlight go >}}
func preProcess(n *html.Node) error {
   if n.Type == html.ElementNode && n.Data == "figure" {
       err := processFigure(n)
      // if error, return error
   }
   for c := n.FirstChild; c != nil; c = c.NextSibling {
      err := preProcess(c)
      // if error, return error
   }
   return nil
}
{{< /highlight >}}

Ensuite, dans le `processFigure`, nous parcourons une fois de plus le sous-arbre, cherchant le nœud `img` principal, et remplaçant ses attributs par ceux du nœud `noscript/img`.

Vous pouvez trouver un code complet dans ce [_gist_](https://gist.github.com/owulveryck/5f9a07762ce40e6f6d9028e76bd798e2)

Une fois que l'arbre HTML est adapté, il peut passer par l'algorithme Arc90.

_Note_ : à ce jour, l'arbre est rendu en HTML pour correspondre à l'API d'Arc90. C'est non optimisé. Je soumettrai éventuellement une PR ou forkerai le projet pour ajouter une nouvelle API qui applique l'algorithme Acr90 directement à un `*html.Node`.

## Deuxième partie : génération de l'ePub

Maintenant que nous avons un contenu approprié, transformons-le en ePub.

Un ePub est un ensemble de fichiers XHTML portant du contenu, ainsi que des images et des fichiers locaux. Tout le contenu est autonome et emballé dans un fichier zip.

Pour générer l'ePub dans l'outil, je m'appuie sur la bibliothèque [`go-epub`](https://github.com/bmaupin/go-epub). Cette bibliothèque est stable, et l'auteur accueille les contributions.

La génération de l'ePub se fait en deux étapes :

1. construire une structure Epub contenant le contenu de l'epub ;
2. générer le fichier epub avec un contenu autonome.

### Première étape : élaboration de l'ePub

Dans la première étape, nous créons le contenu HTML. Le contenu est l'arbre HTML traité précédemment par l'algorithme Arc90.
Le contenu est ajouté comme une seule section dans l'ePub pour des raisons de commodité. Une meilleure façon serait d'analyser l'arbre HTML et de créer une section pour chaque balise `h1`.
Mais comme la cible est de télécharger une seule page, il devrait typiquement y avoir une seule balise `h1` à l'intérieur de la page.

Pour être autonome, il est nécessaire d'analyser cet arbre, à la recherche de contenu distant (essentiellement, les images) et de le télécharger localement.

La bibliothèque go-epub fournit un ensemble de méthodes pour gérer le contenu afin d'accomplir cette tâche en douceur. La méthode [`AddImage`](https://pkg.go.dev/github.com/bmaupin/go-epub?utm_source=godoc#Epub.AddImage), par exemple, crée une entrée dans une carte qui référence le contenu en ligne et fournit une référence à un fichier local.

Ce code, de la documentation, montre comment cela fonctionne :

{{< highlight go >}}
func Example() {
    e := epub.NewEpub("My title")

    // Add an image from a URL. The filename is optional
    imgPath, _ := e.AddImage("https://golang.org/doc/gopher/gophercolor16x16.png", "")

    fmt.Println(imgPath)
    // Output:
    // ../images/gophercolor16x16.png
}
{{< /highlight >}}

Nous devons appeler cette méthode pour chaque élément d'image afin de remplir la carte d'images. De plus, chaque attribut `src` doit être modifié pour utiliser le fichier local.
Nous utilisons le même système qu'avant et utilisons une fonction récursive appliquée au nœud racine de l'arbre HTML :

{{< highlight go >}}
func (d *Document) replaceImages(n *html.Node) error {
    if n.Type == html.ElementNode && n.Data == "img" {
        // find the URL of the image from the current node
        val, f, err := getURL(n.Attr)
        // error checking
        for i, a := range n.Attr {
            if a.Key == "src" {
                img, err = d.AddImage(val, "")
                // error checking
                // Add the local image name as the src attribute of the image
                n.Attr[i].Val = img
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        err := d.replaceImages(c)
        // error checking
    }
    return nil
}
{{< /highlight >}}


#### Retour au problème d'image de _Medium_

Nous avons traité le problème JavaScript dans l'étape de prétraitement. Abordons maintenant le problème réactif.
En fait, la source `img` que nous avons définie dans l'arbre HTML dépend de l'attribut `srcset`.

Dans la fonction `getURL`, nous implémenterons une logique qui définira la valeur source par défaut présente dans l'attribut `src`.
S'il trouve un attribut `srcset`, il l'analysera et le triera, de sorte que le premier élément contienne la plus grande image (nous voulons la meilleure résolution possible).

Nous implémentons l'interface `sort.Sort` sur une structure nouvellement créée `[]srcSetElements` :

{{< highlight go >}}
type srcSetElement struct {
    url            string
    intrinsicWidth string
}

type srcSetElements []srcSetElement

func (s srcSetElements) Len() int { ... }
func (s srcSetElements) Less(i int, j int) bool { ... }
func (s srcSetElements) Swap(i int, j int) { ... }
{{< /highlight >}}

Je n'afficherai pas la fonction getURL complète car son implémentation est simple et présente sur le GitHub du projet.


### Deuxième étape : création de l'ePub

Maintenant que la structure de l'epub est correcte, il suffit d'appeler la méthode [`Write`](https://pkg.go.dev/github.com/bmaupin/go-epub#Epub.Write) qui va :

- télécharger les ressources listées dans la structure Epub ;
- ajouter des métadonnées ;
- créer le fichier zip.

Cette méthode termine le processus et produit le fichier ePub attendu.

## Troisième partie : ajout de fonctionnalités fantaisistes

Maintenant que nous avons un fichier epub, ajoutons quelques fonctionnalités pour améliorer l'expérience du lecteur.

### Saisie des méta-informations

La structure `Article` produite par l'analyseur Arc90 référence un titre, un auteur et une couverture pour le site.
Mais, comme expliqué précédemment, Arc90 est assez ancien, et ces informations sont fournies de nos jours par des éléments OpenGraph.

Arc90 nettoie ces éléments ; par conséquent, nous les saisirons dans l'étape de prétraitement.

Nous nous appuyons sur la bibliothèque [`opengraph`](github.com/dyatlov/go-opengraph/opengraph) en Go pour créer une fonction `getOpenGraph`.
Le point d'entrée d'opengraph lit le contenu à partir d'un `io.Reader`.
Pour optimiser la mémoire, nous implémenterons la méthode `getOpenGraph` comme un middleware.

Elle lira le fichier HTML à partir de l'io.Reader, le traitera, et `Tee` l'original dans un autre reader grâce à un `io.TeeReader`.
La signature de la méthode est :

{{< highlight go >}}
func getOpenGraph(r io.Reader) (*opengraph.OpenGraph, io.Reader) { ... }
{{< /highlight >}}

Encore une fois, [le code complet est disponible sur le dépôt GitHub du projet](https://github.com/owulveryck/rePocketable/blob/8060aa3709b89c6bdf8bf6010027dd38bccd47d7/internal/epub/opengraph.go#L130).

### Génération d'une couverture

Maintenant que nous avons quelques informations, nous pouvons générer une couverture pour l'ePub.
Une couverture est un fichier XHTML qui référence une seule image.

Sur l'image, nous aimerions voir :

- l'image frontale de l'article telle qu'affichée sur les médias sociaux ;
- le titre de l'article ;
- l'auteur de l'article ;
- l'origine de l'article ;

Avec le package `image/draw` de la bibliothèque standard, nous créons une image RGB et composons la couverture frontale.

Le code de la génération de couverture est [ici](https://github.com/owulveryck/rePocketable/blob/master/internal/epub/cover.go).
Ensuite, les méthodes de la bibliothèque go-epub l'ajoutent à l'ePub.

### Intégration GetPocket

Pour compléter le travail, nous pouvons créer une intégration GetPocket pour saisir tous les éléments de la liste de lecture GetPocket et les convertir en ePub.
L'intégration est simple car l'API de getPocket permet de récupérer une structure contenant :

- l'URL originale
- le titre du fichier
- l'image frontale
- les auteurs

Mais, un objectif pourrait être d'exécuter un démon sur la liseuse (par exemple, une reMarkable) ;
par conséquent, la bibliothèque interne gère un mode démon pour récupérer les articles régulièrement (ainsi que lorsque l'appareil se réveille).

Si vous êtes curieux, le mécanisme est implémenté dans [un package pocket](https://github.com/owulveryck/rePocketable/tree/master/internal/pocket) et utilise
le mécanisme que j'ai implémenté il y a un moment pour hacker le projet [remarkable_news](https://github.com/owulveryck/remarkable_news).

### Gestion de MathJax

Une autre fonctionnalité qui manque à l'intégration getPocket sur mon kindle est la capacité à rendre les formules LaTeX.
J'ajoute une étape de traitement supplémentaire pour trouver un contenu mathjax, et créer une image png de la formule.

Pour cela, j'utilise le package [github.com/go-latex/latex](github.com/go-latex/latex).

Le principe est de trouver un TextNode contenant un élément MathJax grâce à une expression régulière :

{{< highlight go >}}
var mathJax = regexp.MustCompile(`\$\$[^\$]+\$\$`)

func hasMathJax(n *html.Node) bool {
    return len(mathJax.FindAllString(n.Data, -1)) > 0
}

func preProcess(n *html.Node) error {
   // ...
    case n.Type == html.TextNode && hasMathJax(n):
        processMathTex(n)
    }
   // ...
}
{{< /highlight >}}

ensuite, la fonction `processMathTex` analyse les formules et les rend dans un fichier encodé en png. Ensuite, le fichier est inséré dans l'arbre HTML dans une balise `img`. L'attribut `src` référence un contenu en ligne de la formule, encodé avec le [principe dataURL](https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/Data_URIs).

## Conclusion et travaux futurs

Je n'utilise pas très souvent l'intégration getPocket, mais j'utilise l'outil `toEpub` quotidiennement pour convertir une page web.

L'intégration getPocket sera utile une fois que j'aurai encodé le fichier de sortie dans un format adapté à la remarkable. Cela semble assez simple, mais je n'ai pas encore pris le temps de le faire.

Jusqu'à présent, mon flux de travail est :

- saisir l'URL sur mon ordinateur portable
- exécuter toEpub localement
- envoyer le résultat à la remarkable avec `rmapi` (et maintenant gdrive)

Le problème est que cela nécessite un ordinateur portable et l'outil installé dessus. Je suis actuellement en train de hacker la bibliothèque go-epub, afin qu'elle n'ait plus besoin d'un système de fichiers, permettant une compilation en webassembly pour faciliter le déploiement.

Restez à l'écoute.