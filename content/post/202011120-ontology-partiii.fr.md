---
title: "Ontologie, graphes et turtles - Partie III"
date: 2020-11-20T20:07:03+01:00
draft: false
tags: ["ontology", "taxonomy", "graph", "turtle", "rdf", "golang", "go"]
summary: "Cet article est une traduction automatique. L'article original a été écrit en anglais. Cet article traite de l'utilisation du graphe en mémoire construit dans la partie II. Dans cet article, je montrerai comment extraire les informations du graphe en utilisant un moteur de modèles. Finalement, nous construirons un serveur web de documentation qui ressemble à schema.org."
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/202011120-ontology-partiii/).

Dans [un article récent](/post/20201113-ontology/), j'ai introduit la notion d'ontologie et de graphe de connaissances.

Ensuite, j'ai exposé dans cet [article](/post/20201117-ontology-partii/) comment analyser un triplestore pour créer un graphe orienté en mémoire.

Cet article vous montrera comment manipuler le graphe. Le cas d'utilisation est une représentation des informations du graphe grâce à un moteur de modèles.

Le moteur devrait être suffisamment extensible pour générer n'importe quel rapport basé sur du texte (HTML, markdown, …)

À la fin de l'article, nous construirons finalement un serveur web basique qui présente de manière similaire les informations de l'ontologie de [schema.org](schema.org).

| schema.org | notre implémentation (localhost) |
|-|-|
| ![](/assets/ontology/schemaorg1.png) | ![](/assets/ontology/schemaorg_olwu.png) | 

**Attention:** la solution est une preuve de concept, son implémentation fonctionne mais n'est pas à toute épreuve. Certains tests la rendraient plus sûre à utiliser, et un TDD influencerait la conception du package. Considérez-le comme un code à des fins éducatives uniquement.


## Le moteur de modèles

De la documentation du package `text/template` de la bibliothèque standard de Go :

> Le package template implémente des modèles pilotés par des données pour générer une sortie textuelle.

Il fonctionne en ...

> ... appliquant [les modèles] à une structure de données. Les annotations dans le modèle se réfèrent à des éléments de la structure de données (généralement un champ d'une structure ou une clé dans une carte) pour contrôler l'exécution et dériver des valeurs à afficher. L'exécution du modèle parcourt la structure et définit le curseur, représenté par un point '.' et appelé "point", sur la valeur à l'emplacement actuel dans la structure au fur et à mesure de l'exécution.

Au premier abord, cela ressemble à un plan pour les [gnomes de southpark](https://fr.wikipedia.org/wiki/Gnomes_(South_Park)). [^1]

[^1]:  Phase 1: Collecter les sous-vêtements / Phase 2: ? / Phase 3: Profit

Rephrasons cela. En essence, le flux de travail est :

- **Collecter les données** (cela a été fait dans l'article précédent)
- créer une structure de données contenant les éléments que vous voulez représenter ;
- créer un squelette de la représentation avec des espaces réservés pour les valeurs de la structure que vous vous attendez à voir ;
- appliquer le modèle à la structure de données ;
- **profit!**

### La structure de données

Tout d'abord, nous avons besoin d'une structure de données.

La structure de données contiendra les informations que nous représenterons via l'application du modèle.

Nous nous placerons dans le contexte du nœud et [_penserons comme un sommet_](/2019/10/14/think-like-a-vertex-using-gos-concurrency-for-graph-computation.html).

Donc la structure que nous définissons commence par référencer le nœud :

{{< highlight go >}}
// Current object to apply to the template
type Current struct {
    Node  *graph.Node
}
{{</highlight>}}

Ensuite, dans le modèle, nous pouvons accéder aux champs exportés et aux méthodes de la structure `graph.Node` définie dans l'article précédent.

{{<highlight go>}}
type Node struct {
    id              int64
    Subject         rdf.Term
    PredicateObject map[rdf.Term][]rdf.Term
}
{{</highlight>}}

Par exemple, ce modèle affichera la sortie d'un appel à la méthode `RawValue()` du champ `Subject` :

{{<highlight go-text-template >}}
The subject is {{ .Node.Subject.RawValue }}
{{</highlight>}}

#### Exemple complet

Comme preuve de concept, nous pouvons écrire un exemple Go simple qui affiche une valeur à partir d'une ontologie de base :

{{<highlight go>}}
func Example() {
	const ontology = `
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix example: <http://example.com/> .

example:PostalAddress a rdfs:Class ;
    rdfs:label "PostalAddress" ;
    rdfs:comment "The mailing address." .
`
	const templ = `
The subject is {{ .Node.Subject.RawValue }}
	`
	type Current struct {
		Node *graph.Node
	}
	parser := rdf.NewParser("http://example.org")
	gr, _ := parser.Parse(strings.NewReader(ontology))
	g := graph.NewGraph(gr)

	postalAddress := g.Dict["http://example.com/PostalAddress"]
	node := g.Reference[postalAddress]

	myTemplate, _ := template.New("global").Parse(templ)
	myTemplate.Execute(os.Stdout, Current{node})

	// output:
	// The subject is http://example.com/PostalAddress
}
{{</highlight>}}

Points importants :

- À la ligne 2, nous définissons une ontologie simple comme une constante chaîne ;
- à la ligne 11, nous définissons un modèle simple pour afficher le sujet d'un nœud ;
- aux lignes 21/22, nous trouvons la référence de nœud qui contient le sujet "PostalAddress" ;
- à la ligne 25, nous appliquons le modèle à l'objet "Current" créé ad-hoc, qui contient le nœud ciblé.

Sur les mêmes principes, nous pourrions afficher le contenu des prédicats et les objets de la structure du nœud. Ce sont des tableaux (slices) et des cartes. Par conséquent, nous devons utiliser les actions intégrées du moteur de modèles pour accéder aux éléments de données (cf [le paragraphe d'action dans la documentation du modèle](https://golang.org/pkg/text/template/#hdr-Actions) pour plus d'informations)

ex:

{{<highlight go-text-template >}}
The subject is {{ .Node.Subject.RawValue }}
and the list of preficates are:
{{ range $predicate, $objects := .Node.PredicateObject -}}
* {{ $predicate }}
  - {{ range $objects }}{{ .RawValue }}{{ end }}
{{ end -}}
{{</highlight>}}

ce qui donne :

```text
The subject is http://example.com/PostalAddress
and the list of preficates are:
* <http://www.w3.org/1999/02/22-rdf-syntax-ns#type>
  - http://www.w3.org/2000/01/rdf-schema#Class
* <http://www.w3.org/2000/01/rdf-schema#label>
  - PostalAddress
* <http://www.w3.org/2000/01/rdf-schema#comment>
  - The mailing address.
```

### Aller plus loin

Ce système est suffisant pour afficher une structure simple. Mais, nous sommes dans le contexte de l'ontologie et des graphes. Il est donc essentiel de pouvoir parcourir le graphe et d'afficher les frères et les arêtes.

Pour y parvenir, nous ajouterons un pointeur vers la structure du graphe elle-même dans le fournisseur de données, c'est-à-dire l'objet "`Current`" :

{{<highlight go>}}
// Current object to apply to the template
type Current struct {
    Graph *graph.Graph
    Node  *graph.Node
}
{{</highlight>}}

Puis nous ajoutons un peu de sucre syntaxique à notre structure :

{{<highlight go>}}
// To the node with edge holding a value from the links array
func (g Current) To(links ...string) []Current { // ... }

func (g Current) From(links ...string) []Current { // ... }

func (g Current) HasPredicate(predicate, object string) bool { // ... }
{{</highlight>}}

_Note_ L'implémentation des méthodes et la documentation sont isolées dans un package de modèles. La référence de ce package peut être trouvée sur [pkg.go.dev/github.com/owulveryck/rdf2graph/template](https://pkg.go.dev/github.com/owulveryck/rdf2graph/template).

À titre d'exemple, complétons notre échantillon de triplestore :
{{<highlight turtle >}}
example:PostalAddress a rdfs:Class ;
    rdfs:label "PostalAddress" ;
	rdfs:comment "The mailing address." .
	
example:addressCountry a rdf:Property ;
	rdfs:label "addressCountry" ;
	rdfs:domain example:PostalAddress ;
	rdfs:comment "A comment" .
	
example:address a rdf:Property ;
	rdfs:label "address" ;
	rdfs:domain example:PostalAddress ;
	rdfs:comment "Physical address of the item." .
{{</highlight>}}

Et étendons le modèle avec un appel à la fonction `To` que nous avons créée.

{{<highlight go-text-template >}}
{{ range $current := .To }}
The subject is {{ .Node.Subject.RawValue }}
and the list of preficates are:
{{ range $predicate, $objects := .Node.PredicateObject -}}
- {{ $predicate }}
  - {{ range $objects}}{{.RawValue}}{{end}}
{{ end -}}
{{ end -}}
{{</highlight>}}

Ensuite, sans modifier le reste du code que nous avons exposé dans le paragraphe précédent, l'exécution de l'exemple donne le résultat suivant :

```text
The subject is http://example.com/address
and the list of preficates are:
* <http://www.w3.org/1999/02/22-rdf-syntax-ns#type>
  - http://www.w3.org/1999/02/22-rdf-syntax-ns#Property
* <http://www.w3.org/2000/01/rdf-schema#label>
  - address
* <http://www.w3.org/2000/01/rdf-schema#comment>
  - Physical address of the item.
* <http://www.w3.org/2000/01/rdf-schema#domain>
  - http://example.com/PostalAddress

The subject is http://example.com/addressCountry
and the list of preficates are:
* <http://www.w3.org/1999/02/22-rdf-syntax-ns#type>
  - http://www.w3.org/1999/02/22-rdf-syntax-ns#Property
* <http://www.w3.org/2000/01/rdf-schema#label>
  - addressCountry
* <http://www.w3.org/2000/01/rdf-schema#comment>
  - A comment
* <http://www.w3.org/2000/01/rdf-schema#domain>
  - http://example.com/PostalAddress
```

## Un service web simple

Maintenant que nous avons construit tous les outils dont nous avons besoin pour rendre le graphe, construisons un serveur web très simple pour présenter les informations du graphe de connaissances. Comme expliqué en introduction, nous utiliserons l'ontologie de schema.org comme base de données (la création du graphe de connaissances est expliquée dans l'article précédent).

### Création du modèle

Chaque représentation d'un nœud est une page html unique. On y accède par un appel à "http://serviceurl/NodeSubject".

Le squelette de la page est un modèle.
Pour simplifier les choses, nous divisons le modèle en trois sous-modèles.

- un modèle `main` qui créera la structure HTML et la table extérieure
- un modèle pour afficher une `class` comme une structure `tbody`
- un modèle de propriété pour afficher une ligne de la structure tbody

{{<highlight go-text-template >}}
{{ define "main" }}
<!DOCTYPE html>
<!-- boilerplate of the HTML file -->
    <table class="blueTable">
        <!-- header of the table -->
        {{ template "rdfs:type rdfs:Class" . }}
    </table>
</html>
{{ end }}

{{ define "rdfs:type rdfs:Property" }}
<tr>
{{ calls to display the subjects and predicates }}
</tr>
{{ end }}

{{ define "rdfs:type rdfs:Class" }}
<tbody>
    <tr>
    <!-- The rest of the table structure -->
    {{ range over the "To" nodes for the graph held in `Current` }}
        {{ for each node, if its type is "property" }}
            {{ template "rdfs:type rdfs:Property" . -}}
    {{ range over the "From" nodes for the graph held in `Current` }}
        {{ for each node, if its type is "class" }}
                {{ template "rdfs:type rdfs:Class" . -}}
    </tr>
</tbody>
{{ end }}
{{</highlight>}}

Il n'y a pas beaucoup d'intérêt à afficher tout le câblage dans cet article.
Le modèle complet est disponible [ici](https://github.com/owulveryck/rdf2graph/blob/7a6127ae4428c5501a1d07eb541a16fb4ee3ad83/examples/webview/index.tmpl#L1-L59)

Cela affichera une page bien formatée lorsqu'elle sera mélangée avec un nœud du graphe.

#### Le serveur web

Pour simplifier les choses, encapsulons tout cela dans un simple serveur web. L'objectif est de pouvoir afficher n'importe quel nœud du graphe lorsqu'on y accède via une URL.

Tout d'abord, nous créons une structure qui implémentera l'interface [http.Handler](https://golang.org/pkg/net/http/#Handler).

Par commodité, cette structure contient le graphe, le modèle et une table de hachage du rdf.IRI. Cette dernière est utilisée pour raccourcir les URL (appeler _http://myservice/blabla_ au lieu de _http://myservice/http://example.com/blabla_)

{{<highlight go >}}
type handler struct {
	namespaces map[string]*rdf.IRI
	g          graph.Graph
	tmpl       *template.Template
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
{{</highlight >}}

La méthode ServeHTTP est composée de trois parties :

- extraire le nœud de l'URL
- vérifier si le nœud existe
- appliquer le modèle et écrire le résultat sur le _ResponseWriter_

Je ne détaillerai pas tout le code pour implémenter cela. Vous pouvez regarder l'implémentation [sur GitHub](https://github.com/owulveryck/rdf2graph/blob/7a6127ae4428c5501a1d07eb541a16fb4ee3ad83/examples/webview/main.go#L60-L77).

Nous devons assembler tout le code des articles pour :

- analyser un triplestore à partir d'un `io.Reader` ;
- créer un graphe en mémoire ;
- lire le fichier de modèle et créer le modèle Go ;
- créer le gestionnaire pour référencer le graphe et le modèle ;
- enregistrer le gestionnaire pour servir une requête HTTP ;

Cela semble compliqué, mais c'est assez facile et direct si vous faites un peu de Go. Quoi qu'il en soit, une implémentation d'exemple est sur [GitHub](https://github.com/owulveryck/rdf2graph/tree/main/examples/webview). Pour la lancer, faites simplement :

```bash
curl -s https://schema.org/version/latest/schemaorg-current-http.ttl | go run main.go
```

Ensuite, pointez votre navigateur vers "http://localhost:8080/PostalAddress" et vous devriez obtenir quelque chose qui ressemble à ceci :

![](/assets/ontology/schemaorg_olwu.png)

## Conclusion

C'était le dernier article sur l'ontologie. À travers les morceaux, nous avons découvert une façon de décrire un graphe de connaissances pour qu'il soit facile à écrire pour un humain et suffisamment efficace à analyser pour une machine. Ensuite, nous avons construit un graphe en mémoire et exploité ce graphe pour représenter la connaissance. La couche de représentation peut être vue comme une couche de projection qui expose les informations requises pour un domaine fonctionnel spécifique.

Maintenant, amusons-nous avec l'ontologie pour :

- documenter les domaines fonctionnels et agir comme un langage ubiquitaire partagé
- fournir une carte d'information pour localiser et utiliser une partie spécifique de la connaissance d'un système