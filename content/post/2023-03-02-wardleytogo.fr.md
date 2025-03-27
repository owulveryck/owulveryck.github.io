---
title: "La logique derrière wardleyToGo"
date: 2023-03-02T21:23:17+02:00
lastmod: 2023-03-02T21:23:17+02:00
draft: false
keywords: []
description: "wardleyToGo est une bibliothèque et un langage pour décrire les cartes Wardley sous forme de code/données. Voici la logique derrière cette bibliothèque"
tags: []
categories: []
author: "Olivier Wulveryck"

# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: true
toc: true
autoCollapseToc: false
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: false
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/2023-03-02-wardleytogo/).

Les cartes de Wardley sont une façon d'exprimer une entreprise, un marché ou tout autre système par le biais d'esquisses. La carte est une esquisse qui offre une conscience situationnelle sur un certain sujet.

En tant qu'esquisse, la façon évidente de dessiner une carte est avec du papier et un stylo.
Bien que ce soit un excellent point de départ, une représentation papier d'une carte a un problème : elle est statique.
Ce que je veux dire par là, c'est que l'ajustement du placement de certains composants dans la phase de conception peut être fastidieux (même avec un bon crayon à gomme).

Représenter une carte numériquement présente quelques avantages : elles sont faciles à partager et faciles à **maintenir** et à **exploiter**.

- **Maintenir** une carte est, dans mon contexte, la possibilité de faire des ajustements sans changer sa signification (pensez aux ajustements qui peuvent se produire après une discussion avec des pairs).
- **Exploiter** une carte consiste à comprendre le paysage pour orienter les décisions futures.

_Note_ : Un point important à garder à l'esprit lors de la cartographie est : une carte ne doit pas être utilisée pour illustrer l'histoire que vous voulez raconter, l'histoire doit être extraite de la carte.

Dans la phase d'orientation, pouvoir retirer certains éléments de la carte, par exemple, pour se concentrer sur un certain chemin est utile.

Bien sûr, nous pourrions enregistrer l'historique de création et faire quelques annulations/rétablissements pour avoir une représentation intermédiaire, mais cela impliquerait que les étapes de création reflétaient le chemin.
C'est orthogonal avec l'idée de chercher une stratégie par l'observation **puis** l'orientation.

La meilleure façon d'avoir une représentation partielle est donc d'avoir une représentation intermédiaire de la carte où nous pouvons facilement commenter des blocs pour cacher des éléments.

## Les cartes sous forme de code

Dans les phases de conception et de maintenance, exprimer une carte sous forme de code présente une certaine flexibilité.
Le code est la source de la carte et les outils et pratiques pour gérer le code source sont répandus (gérer le code source est une commodité).
Des outils comme `git` offrent des capacités pour :
- versionner et étiqueter la carte;
- collaborer sur la carte (même de manière asynchrone ou à distance)
- stocker nativement l'historique de la carte.

[onlinewardleymaps](https://onlinewardleymaps.com/) est un outil en ligne créé par [@damonsk](https://twitter.com/damonsk).
Il est mature et largement utilisé. Pour moi, cet outil a largement contribué à pousser les cartes-sous-forme-de-code dans l'_étape_ II de l'évolution (comme décrit par Simon Wardley).

Le langage de définition d'onlinewardleymaps s'appelle `owm`.

## Prendre ma carte "à emporter"

Au début, la seule version de l'outil était une version en ligne. C'est bon pour l'expérience utilisateur car cela ne nécessite aucune configuration pour commencer à cartographier.
Le problème est le couplage entre l'outil pour rendre la carte et le langage. Vous ne pouvez pas utiliser `owm` en dehors du navigateur.

### Versionner la carte

Le problème est que la gestion du code source était un peu fastidieuse car elle nécessitait beaucoup de copier/coller depuis et vers les outils.

L'extension Visual Studio, qui est apparue plus tard, a rendu la gestion du code source un peu plus simple. En ce sens, elle facilite grandement le versionnement de la carte.

Mais il n'y a aucun moyen d'exporter facilement la carte pour la stocker avec une version particulière du code, car vous devez exporter la carte manuellement.

Cela rend un peu difficile de prendre la carte "à emporter" (ou à emporter selon la saveur d'anglais que vous pratiquez (sic)).

### Besoin de base : CI/CD

Depuis des années maintenant, l'intégration continue (CI) et la livraison continue (CD) ont prouvé certains avantages dans le cycle de publication d'un actif numérique.
L'idée est que chaque révision qui est validée déclenche une construction et un test automatisés. Avec la livraison continue, les changements de code sont automatiquement construits, testés et préparés pour une publication en production

Comme j'étais capable de valider le code source de mes cartes, je voulais utiliser un mécanisme de CI/CD pour compiler mon code source et rendre mes cartes.

onlinewardleymaps ne fournit aucun SDK me permettant d'utiliser le moteur de rendu dans un processus de construction sans interface.
En tant que geek, j'aurais pu soulever un problème et commencer à contribuer au projet.
Mais je voulais construire le mien pour comprendre en profondeur comment construire une carte de l'intérieur et parce que je trouvais amusant d'avoir un projet parallèle de plus.

## Construire un SDK pour dessiner une carte "sous forme de code"

Je suis un grand fan du langage Go pour diverses raisons. En un coup d'œil, les raisons pour lesquelles j'utilise Go pour mon SDK sont :

- Je trouve amusant de coder en Go;
- Je maîtrise suffisamment le langage pour accélérer le développement et me concentrer sur la conception de ce que je veux réaliser;
- le langage est adapté pour construire des "outils Ops".

Nommer les choses est difficile, j'ai nommé mon SDK warldeyToGo (prenez votre carte "à emporter" avec Go).

La conception du SDK est :

- un package central qui agit comme une représentation intermédiaire de ce qu'est une carte. (en un coup d'œil, c'est un graphe orienté dans lequel les nœuds sont des composants qui sont capables de donner leur position dans un canevas euclidien)
- un ensemble de composants qui implémentent la représentation intermédiaire et qui sont capables de se représenter en SVG (par exemple, j'ai des composants Wardley et des composants Team Topologies)
- des analyseurs pour des langages de haut niveau qui transpirent la représentation dans la représentation intermédiaire.

Pour commencer, j'ai implémenté un analyseur pour le langage `owm`.

Ainsi, je peux construire ma carte avec `onlinewardleymaps`, extraire le code source (la représentation `owm`) et construire un outil avec le SDK pour transpirer `owm` dans la représentation intermédiaire et le rendre en SVG.

## Un langage de haut niveau pour exprimer la carte "sous forme de données"

Donc le SDK permet vraiment de coder des cartes Wardley. La syntaxe `owm` est donc simplement une interface utilisateur entre le besoin d'exprimer une carte et la représentation.

> Tout comme les interfaces utilisateur sont le canal entre les humains et la fonctionnalité des produits, les produits sont le canal entre les clients et l'équipe de personnes résolvant un problème particulier.

onlinewardleymaps est une solution au besoin de représenter une carte numériquement.
wardleyToGo est une solution pour coder une carte.

Mais ce qui me manquait, c'est une solution qui permet à un ordinateur de m'assister dans la conception d'une carte.

### Le problème avec la représentation euclidienne

Le problème que je rencontre lors de la conception d'une carte avec un outil basé sur la représentation euclidienne (comme owm) est qu'il me demande de penser à la position exacte d'un composant sur le canevas (en termes de coordonnées X et Y).

Le problème principal est avec l'axe vertical... qui n'est pas un axe :

![](/assets/images/wardley_axis.png)
[source](https://twitter.com/swardley/status/1237707981116055552)

De plus, l'axe d'évolution est décomposé en 4 étapes, et placer les composants sur une certaine étape nécessite l'utilisation de notre _système 2_ (comme décrit dans le livre ["_Système 1 / Système 2 : Les deux vitesses de la pensée_" par Daniel Kahneman](https://fr.wikipedia.org/wiki/Syst%C3%A8me_1_/_Syst%C3%A8me_2_:_Les_deux_vitesses_de_la_pens%C3%A9e))

### Un langage informatique pour façonner notre façon de penser

Donc, mon objectif est maintenant :

- Avoir un langage suffisamment facile pour être utilisé par mon _système 1_
- Qui façonne la façon dont je pense à ma carte et qui m'aidera donc dans la phase de conception.

Je peux également ajouter des algorithmes pour m'aider dans le placement des composants.
Pour ce faire, j'ai besoin d'implémenter un nouveau langage.
La conception du SDK wardleyToGo permet assez facilement de concevoir un langage par essais et erreurs. La représentation intermédiaire et les bibliothèques de composants facilitent le rendu et réduisent donc la boucle de rétroaction.

Voyons maintenant ce qu'il faut attendre de ce nouveau langage que nous appellerons `wtg` (pour WardleyToGo... nommer les choses est très difficile)

### Penser en termes de chaîne de valeur et de composants : introduction à wtg

Concevoir une carte est essentiellement un processus en deux étapes :

1. créer la chaîne de valeur
2. évaluer les composants sur l'axe d'évolution

Dans la phase de conception : la position du composant sur l'axe d'évolution n'a pas d'impact sur la chaîne de valeur, et la visibilité du composant (son placement vertical) n'a pas d'impact sur sa position sur l'axe d'évolution.

De plus, comme vu précédemment, la chaîne de valeur n'est pas un axe. Par conséquent, nous pouvons complètement nous débarrasser de toute représentation euclidienne dans la première étape.

#### Chaîne de valeur

La visibilité des composants est relative aux autres composants.

La façon dont nous décrivons une chaîne de valeur dans `wtg` est en utilisant un tiret (`-`) pour lier les composants entre eux. Donc `a - b` signifie que `a` dépend de `b`. Ensuite, pour placer le composant verticalement, nous pouvons ajouter plus de tirets.

Par exemple :

```
a - b 
a -- c 
a --- d
```

signifie ceci :

![](/assets/images/simplevc.svg)

_Notez_ que le placement horizontal est, pour l'instant, sans signification

Le placement vertical est calculé par un algorithme. Cela me permet donc de me concentrer sur la chaîne de valeur elle-même.

#### Axe d'évolution

Après avoir construit la chaîne de valeur, nous pouvons placer les composants sur l'axe horizontal.
Ce placement est une configuration de chaque composant et indépendant de la chaîne de valeur.

La syntaxe pour placer un composant ressemble à ceci : `|..|..x..|..|..|` où chaque intervalle entre deux barres verticales (`|`) est une étape d'évolution et le `x` est le placement du composant.
Vous pouvez influencer le placement dans une étape en ajoutant des points `.` avant ou après le `x`.

Il est également possible de typer les composants (`build, buy, outsource`) ou d'ajouter des couleurs. D'autres options peuvent être facilement ajoutées ultérieurement.

Le fameux exemple du "salon de thé" ressemble à ceci en `wtg` :

```dot
business - cup of tea
public - cup of tea
cup of tea - cup
cup of tea -- tea
cup of tea --- hot water
hot water - water
hot water -- kettle
kettle - power

business:   |...|.....|...x.|..........|
public:     |...|.....|.....|..x.......|
cup of tea: |...|.....|..x..|..........|
cup:        |...|.....|.....|.....x....|
tea:        |...|.....|.....|.....x....|
hot water:  |...|.....|.....|....x.....|
water:      |...|.....|.....|.....x....|
kettle:     |...|...x.|..>..|..........|
power:      |...|.....|....x|.....>....|
```

Vous pouvez le construire étape par étape en suivant [ce tutoriel](https://owulveryck.github.io/wardleyToGo/tutorials/helloworld/)

## Conclusion et références

Le langage `wtg` répond à mes propres besoins. J'ai réalisé plusieurs cartes avec.

Il est au-delà de la portée de cet article de décrire complètement la grammaire. J'ai commencé un site web basé sur le [cadre de documentation divio](https://documentation.divio.com/) comme compagnon du langage.
Vous pouvez trouver la référence du langage [ici](https://owulveryck.github.io/wardleyToGo/reference/wtg/).

Il existe une version en ligne qui peut être utilisée pour créer facilement des cartes avec le langage.
Même si j'ai ajouté quelques petites fonctionnalités à la démo comme la possibilité de cacher des liens dans la chaîne pour avoir une meilleure observation des composants, la version en ligne est une démo.
Considérez-la comme une preuve de valeur et non comme un outil de production.

En plus de cela, un ensemble d'outils en CLI sont présents dans le [dépôt du projet](https://github.com/owulveryck/wardleyToGo). Par exemple, il y a un outil pour surveiller un fichier wtg et le rendre dans
le navigateur en direct. Ainsi, wtg pourrait être édité avec votre éditeur de texte préféré et, vous pouvez présenter la carte dans un appel zoom.

Une fonctionnalité intéressante que j'aimerais voir à l'avenir est la possibilité de regrouper certains éléments et de laisser l'ordinateur ajouter une bordure environnante dans la phase de rendu.

En conclusion, voici une simple carte qui tente de résumer les idées exposées dans ce post :

![](/assets/images/wardleyToGo.svg)

Vous pouvez jouer avec la carte en suivant [ce lien](https://owulveryck.github.io/wardleyToGo/demo/?url=https://raw.githubusercontent.com/owulveryck/wardleyToGo/main/docs/content/en/illustration.wtg)