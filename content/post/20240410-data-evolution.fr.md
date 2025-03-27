---
title: "Données-en-tant-que-Produit et Contrats de Données : Une approche évolutionnaire de la maturité des données"
date: 2024-04-09T12:15:33+01:00
lastmod: 2024-04-09T12:15:33+01:00
images: [/assets/data_certitude.png]
draft: false
keywords: []
summary: Cet article est une traduction automatique. L'article original a été écrit en anglais. En utilisant le modèle d'évolution de Simon Wardley, je propose un cadre pour visualiser la maturité des données dans un contexte d'entreprise, en soulignant l'importance de traiter les données comme un produit et de mettre en œuvre des contrats de données pour faciliter l'intégration et assurer la confiance. En fin de compte, je suggère que commencer par se concentrer sur les données-en-tant-que-produit est crucial pour les organisations qui se lancent dans leur parcours de maillage de données, ouvrant la voie à une transformation complète et agile.
tags: []
categories: []
author: ""

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S.
comment: false
toc: true
autoCollapseToc: false
# You can also define another contentCopyright.
contentCopyright: false
reward: false
mathjax: false

# Uncomment to add to the homepage's dropdown menu; weight = order of article
# menu:
#   main:
#     parent: "docs"
#     weight: 1
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/20240410-data-evolution/).

## Contexte

Je suis un défenseur précoce du paradigme du maillage de données (data mesh) depuis que Zhamak Dehghani l'a proposé pour la première fois.
En tant que partisan de la première heure, j'ai identifié le potentiel de cette nouvelle approche de l'organisation des données.
Quatre ans se sont écoulés, et le paradigme du maillage de données a effectivement gagné une large acceptation.
Cependant, je n'ai pas encore vu de plan de transformation concis et pratique du maillage de données au sein d'une organisation.

Quand je dis "défenseur", je veux dire que j'ai mis en évidence les avantages de ce paradigme, qui sont enracinés dans ses quatre piliers principaux :

- L'orientation de la conception vers les **domaines**
- Appliquer une **pensée produit aux données** (données-en-tant-que-produit)
- **Gouvernance** computationnelle fédérée
- Développer des **plateformes de données** en libre-service

Vient ensuite la question, **par où commencer réellement pour mettre en œuvre le maillage de données ?**

En principe, toute organisation peut commencer son voyage vers le maillage de données en se concentrant sur ces quatre piliers.

**Commencer par** une conception **orientée domaine** jette les bases d'une compréhension approfondie du maillage de données.
Cela signifie non seulement définir le maillage de données comme objectif, mais aussi **s'assurer que la décomposition** du **domaine est synchronisée avec** la **structure existante** de l'organisation.
Cependant, c'est une **approche profondément conceptuelle** qui pourrait ne pas donner de résultats immédiats, et de plus, elle manque de l'agilité qui est si bénéfique.

La **gouvernance computationnelle fédérée** et la **plateforme de données en libre-service** sont simplement des **facilitateurs** du maillage de données.
Ils partagent un objectif commun : simplifier le développement des données-en-tant-que-produit et la **création d'interconnexions**, soutenant essentiellement le maillage.
On peut essayer de les mettre en œuvre comme fondation, mais pour mailler quoi ?

Il reste donc à s'attaquer aux données-en-tant-que-produit, une pierre angulaire du maillage de données dont j'ai déjà parlé.

Il est intéressant de noter que plusieurs organisations prétendent avoir mis en œuvre le **maillage de données "par accident"**, percevant ce paradigme comme l'évolution naturelle de la gestion des données.

Dans cet article, je tente d'appliquer un modèle bien reconnu de progression évolutionnaire pour comprendre l'évolution des données.
**L'objectif** est d'aider à **visualiser la maturité des données** et d'aider les entreprises à **identifier** leur **point de basculement**,
c'est-à-dire quand elles commenceront à voir des **avantages significatifs** de la **mise en œuvre de contrats de données** et du **traitement des données comme un produit**.

## Modélisation de l'évolution

Je vais d'abord expliquer le modèle que je vais utiliser.
Ce modèle est connu sous le nom de modèle d'évolution de Simon Wardley et est implémenté avec succès dans les cartes Wardley.
Mon objectif ici n'est pas de décrire le paysage d'une entreprise spécifique, donc je n'aurai pas besoin d'une carte complète.
Au lieu de cela, j'utiliserai le modèle d'évolution et essaierai d'appliquer sa fonction générale aux données.

**Avertissement concernant le modèle :**
La théorie de l'évolution est bien adaptée pour être appliquée dans un environnement concurrentiel où tout évolue en fonction de l'offre et de la demande. Je considère les entreprises qui sont soumises à ces contraintes de concurrence et, par conséquent, leurs données suivront également ces règles. Par conséquent, le modèle s'appliquera.

**Pourquoi modéliser l'évolution ?** Comprendre l'évolution, c'est comprendre comment les composants changent au fil du temps. Modéliser l'évolution consiste à trouver un modèle pour potentiellement fournir des informations sur les trajectoires futures de ces composants.

### Le modèle en un coup d'œil

Simon Warldey avait besoin d'un moyen de représenter l'évolution des composants sur sa carte.
Il ne pouvait pas s'appuyer sur une échelle de temps de base car cela l'aurait empêché de comparer des éléments hétérogènes et aurait brisé la cohérence du mouvement.

Par exemple, sur une échelle de temps, la distance entre la genèse et la maturité d'une voiture (environ 100 ans) aurait été beaucoup plus grande que la distance entre les mêmes points pour un smartphone (environ 10 ans).
Il a découvert que **l'évolution** est une **fonction** de son **ubiquité** **et** de sa **certitude**.

Dans une économie de marché, **l'ubiquité est guidée par la demande**. Plus de demande induit plus de présence. C'est une déclinaison de la théorie de la [diffusion des innovations de Rogger](https://en.wikipedia.org/wiki/Diffusion_of_innovations).
La certitude vient de la matrice de Stacey. La matrice postule que lorsque la disponibilité des composants ou des informations clés augmente, la certitude concernant les résultats de la prise de décision augmente également, permettant une planification et une exécution plus prévisibles et éclairées.
En un sens, **la certitude est guidée par l'offre**.

Par exemple, considérons une entreprise manufacturière qui produit des gadgets électroniques. Dans ce scénario, l'un des composants critiques dont ils ont besoin sont les puces semi-conductrices. Lorsque l'approvisionnement de ces puces est faible en raison de pénuries sur le marché ou de problèmes logistiques, l'entreprise fait face à une grande incertitude concernant ses calendriers de production et sa capacité à répondre à la demande des clients.

Une analyse empirique a conduit à cette représentation :
![](https://blogger.googleusercontent.com/img/b/R29vZ2xl/AVvXsEgatdAD8t3Jp7BEjlcpxMwwUMGPbmu-zs9kwEX4KlVqZ31VwHzShmyAr1ZE0zC4YWUnTXWncgIVFPr6_-CQhKn8FO2He4qs-KGd5CrlLcW7S-ZzNxUZLAQqDQE-Vqe11g8Rt7eOaA/s1600/Screen+Shot+2014-03-15+at+18.48.03.png)
(source [Blog de Simon Wardley blog.garvediance.org](https://blog.gardeviance.org/2014/03/on-mapping-and-evolution-axis.html))

Le modèle est une sorte de courbe en S.

### Le modèle des données selon Wardley

L'analyse du modèle a permis de formaliser quatre étapes d'évolution étiquetées par défaut _Genèse, Construction personnalisée, Produit et Commodité_ :

![](https://blogger.googleusercontent.com/img/b/R29vZ2xl/AVvXsEjMFN3o1ujMDfd4y78hHCRFmPSTf9BP5C_Ej1jtEyZrmNC21aBw-18gAbVk88nKHdVa3gd_-D3z3pKKfO4Wa6XsIa1BuTkeiazqGLdu8vlUPsSaXeDgbkbvrMy3CSHlUiqk5ol1ig/s1600/Screen+Shot+2014-01-09+at+13.26.48.png)
(source [Blog de Simon Wardley blog.garvediance.org](https://blog.gardeviance.org/2014/03/on-mapping-and-evolution-axis.html))

Ce ne sont que des étiquettes communes pour une forme de capital. Pour les données, selon la théorie de Wardley, les étiquettes des quatre étapes sont : **_Non modélisé, divergent, convergent et modélisé_** :
![](https://i0.wp.com/learnwardleymapping.com/wp-content/uploads/2020/01/20200122_124314.jpg?resize=1080%2C466&ssl=1)
(source : [learnwardleymapping.com](https://learnwardleymapping.com/2020/01/22/visualizing-the-interaction-of-evolution-and-data-measurement/))

### Dérivation du modèle

#### Certitude des données

Revisitons le mécanisme de certitude pour déterminer si nous pouvons ajuster le modèle pour accommoder l'évolution des données au sein d'une entreprise. Je considère la certitude équivalente au niveau de confiance dans la décision prise sur la base de ces données. Voici les étiquettes que j'utiliserai :

- **Données brutes** : Dans mon expérience, les données commencent à l'état brut pendant la phase exploratoire. Elles manquent d'ubiquité, résidant uniquement dans la base de données et accessibles uniquement via un service et/ou une API, essentiellement un **produit de données** (un produit guidé par les données).
- **Données organisées** : Ceci marque la deuxième étape de la certitude des données. Les experts en données entrent en jeu pour assurer l'exactitude et la pertinence de la représentation des données pour l'entreprise.
- **Faisant autorité** : La dernière étape de certitude. Les données sont pertinentes, complètes, cohérentes, documentées et approuvées par des experts du domaine.

Les données **brutes** correspondent à la **première étape** de l'évolution. C'est une étape où nous définissons des preuves de concepts par exemple. Ensuite, les données **organisées** sont liées aux phases deux et trois. Et finalement, la **dernière étape** est lorsque les données font **autorité**.

#### Les étiquettes des quatre étapes d'évolution

Concernant la notion de certitude et d'ubiquité, catégorisons les 4 étapes d'évolution :

1. **POC** : Cette étape implique la validation des concepts.
2. **Application** : Dans cette étape, les données sont nettement liées à un cas d'utilisation spécifique.
3. **Domaine** : Cette étape est où cela devient intéressant : les données représentent une solution qui peut être utilisée pour aborder divers cas d'utilisation au sein du même domaine (pensez au domaine comme un espace problème, similaire au DDD).
4. **Entreprise** : Cette étape englobe tous les domaines, représentant la totalité de tous les problèmes abordés par une entreprise.

Voici la représentation de ces éléments sur un diagramme :

### La représentation

![Une courbe en S représentant l'évolution des données, le X est la certitude et Y est l'ubiquité. Il y a une division : le bas de la S contient des données brutes, le milieu est organisé, et le haut fait autorité. Le point d'inflexion est indiqué comme contrat de données.](/assets/data_certitude.svg)

## Utilisation du diagramme : Données-en-tant-que-produit et contrat de données

Maintenant, utilisons le diagramme.

Les données suivront probablement la courbe d'évolution en S. Ce qui est intéressant, c'est l'évolution des propriétés des données.
Transformer des données brutes en données organisées est maîtrisé. Il existe des processus de conception majeurs qui sont utiles dans une telle transition.

Transformer les données organisées en données faisant autorité implique que les données sont accessibles et utilisables, maintenues, précises, mais le changement est principalement que les données sont **approuvées par des parties de confiance**.
À l'échelle de l'entreprise, cela signifie que le domaine est responsable de ses données car le domaine est, par défaut, une partie de confiance dans l'organisation concernant un domaine d'activité spécifique.

La transition n'est pas si nette lorsque les données quittent leur prison : lorsqu'elles sont exposées au domaine.

C'est le point où la pensée produit appliquée aux données apporte de la valeur. Et c'est le point où un contrat de données est utile pour :

- Faciliter l'intégration dans d'autres cas d'utilisation du domaine
- Apporter la confiance dans les données

Par conséquent, penser aux données comme à un produit, comme tout autre produit, est quelque chose qui est requis dans la phase d'exploration (cela peut même être vu comme du sur-engineering), mais le modèle illustre à quel point il est important de traiter les données comme un produit pour servir un objectif général pour l'entreprise.

## Conclusion

En résumant, j'ai toujours été aux prises avec une question : _par où commencer lorsqu'on cherche à mettre en œuvre le paradigme du maillage de données ?_
À travers le voyage d'exploration de ce concept, mon aperçu le plus récent et le plus profond est : le point de départ le plus stratégique réside dans le produit de données.

Le modèle présenté souligne le rôle central du produit de données. Il est projeté comme une solution efficace à un problème impératif : son importance significative devient évidente lorsque les données migrent d'une sphère d'application unique vers le domaine plus large.
Au-delà de cela, il devient absolument critique lorsque les données sont censées délivrer une valeur tangible qui dépasse leur domaine défini d'origine.

La prochaine phase de notre voyage pour comprendre le paradigme du maillage de données implique la formalisation d'une méthode pour évaluer avec précision la maturité des données.
En examinant chaque élément de données, contrat par contrat, et domaine par domaine, nous nous rapprochons de la construction d'un maillage complet et efficace.
Tout au long de ce processus, il est crucial de se rappeler de considérer les données comme un produit. Ce faisant, cela apportera des récompenses à une organisation à mesure qu'elle évolue et mûrit dans ses stratégies de gestion des données.