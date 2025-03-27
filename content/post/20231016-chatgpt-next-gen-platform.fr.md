---
title: "L'aube des assistants personnels basés sur les LLM : l'émergence d'une nouvelle économie de plateforme"
date: 2023-10-16T10:11:52+02:00
lastmod: 2023-10-16T10:11:52+02:00
draft: false
keywords: []
description: ""
tags: []
categories: []
author: "Olivier Wulveryck	"
images: [/assets/chatgpt-platform-me.png]
summary: L'article examine l'évolution de **l'économie de plateforme** où des plateformes numériques, comme **ChatGPT** alimenté par des Grands Modèles de Langage (LLM), servent d'intermédiaires reliant différentes parties prenantes.
  Ces plateformes, contrairement aux **pipelines** traditionnels, exploitent la technologie numérique pour créer de la valeur grâce à des interactions personnalisées de masse.
  

  En utilisant un cas d'utilisation hypothétique, l'article démontre comment ChatGPT peut être un assistant personnel intuitif, reliant les utilisateurs avec divers fournisseurs de services.
  

  Cependant, avec la montée en puissance de telles plateformes, des défis similaires au SEO dans les moteurs de recherche sont anticipés.
  

  Des approches comme **l'Ingénierie de Prompts et l'Optimisation Automatique des Choix** deviendront essentielles.

  
  Enfin, un défi crucial pour les fournisseurs est d'être choisis par les systèmes d'IA dans un paysage dominé par quelques géants numériques.
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/20231016-chatgpt-next-gen-platform/).

## Introduction
**L'économie de plateforme** est un environnement économique dans lequel des plateformes numériques agissent comme intermédiaires, connectant diverses parties prenantes et leur permettant de transiger les unes avec les autres de manière transparente.
Ces plateformes exploitent la puissance des effets de réseau, créant de la valeur en facilitant les échanges entre utilisateurs.

Considérons une place de marché numérique comme exemple :

- Elles agissent comme un centre de connexion entre vendeurs (offre) et acheteurs (demande).
- Elles tirent parti de l'effet réseau : plus les personnes utilisent la place de marché, plus la valeur de la plateforme augmente pour tous les utilisateurs. Une plus grande sélection de produits attire plus d'acheteurs, et plus d'acheteurs attirent plus de vendeurs.
- Elles tirent leur valeur de la technologie : l'économie de plateforme est largement rendue possible par la technologie numérique, et les places de marché exploitent cette technologie pour opérer, évoluer et améliorer leurs plateformes au fil du temps.

Récemment, les avancées dans les Modèles de Langage comme **ChatGPT**, alimentés par les Grands Modèles de Langage (LLM), ont ouvert la voie au développement d'assistants personnels très intuitifs et capables.
Ces assistants alimentés par les LLM peuvent devenir une évolution de l'économie de plateforme.
Une plateforme où **vitrines numériques**, **fournisseurs de services** et **utilisateurs** **convergent**, créant un **écosystème symbiotique**.

Dans cet article, j'expliquerai comment le mécanisme de _plugin_ de ChatGPT est une pierre angulaire dans la construction de la plateforme.

Dans la première partie, j'expliquerai un cas d'utilisation fictif du point de vue de l'acheteur.
Dans la deuxième partie, après un rappel de la notion de pipeline et de plateforme, j'exposerai comment un fournisseur peut utiliser le mécanisme de plugin pour répondre aux exigences de la plateforme.
J'insisterai sur la partie communication et standardisation, discutant de la représentation intermédiaire et du langage humain.

_Note_ mon expérience est principalement axée sur ChatGPT, donc je peux utiliser ChatGPT au lieu de LLM.

## Exemple de cas d'utilisation

Considérez un scénario où l'on doit se préparer pour un dîner avec un budget de **100 $**.
Le dîner est prévu pour 20h, avec une heure de route, et il est actuellement 14h.

L'objectif est de trouver des vêtements appropriés dans le budget, et de les acheter dans un magasin à proximité pour arriver au dîner à temps.

C'est là qu'un assistant alimenté par un LLM comme ChatGPT entre en jeu.
L'assistant, exploitant un réseau de plateformes numériques, peut aider à identifier les options de vêtements, comparer les prix et localiser un magasin à proximité, tout en s'assurant que l'utilisateur reste dans les limites du budget et du calendrier.

![](/assets/chatgpt-platform-illustration_small.png)

### Raisonnement
Avec mon collègue [Nicolas](https://www.linkedin.com/in/nicolasgutierrez/), nous avons effectué l'exercice de créer une carte (Wardley) pour ce besoin.
Le besoin fondamental concerne la recherche d'un bien.
Étant donné l'expérience de Nicolas en graphes de connaissances et en recherche sémantique, nous avons d'abord exploré comment ces éléments pourraient répondre au cas d'utilisation.
Par la suite, nous avons considéré le rôle potentiel des assistants personnels dans ce contexte.
Voici la carte qui en résulte :
![](/assets/chatgpt-kg.png)

Nous nous sommes alors demandé : Le LLM doit-il être en interne pour le détaillant, ou peut-il fonctionner de manière indépendante en utilisant les données fournies par le détaillant ?
Les couleurs indiquent les composants appartenant au détaillant (notez que je ne dis pas nécessairement qu'ils doivent rester en interne, en particulier les produits de base).
Les composants appartenant à l'utilisateur final sont affichés en vert.

Étonnamment, après une discussion sur la carte, il est devenu clair que l'instanciation du LLM pourrait être plus bénéfique si elle était maintenue hors du contexte.
C'est logique, car la puissance de calcul nécessaire pour exécuter le LLM est substantielle et devrait probablement rester avec une entreprise externe qui se spécialise dans la construction d'assistants personnels.

Alors, quel pourrait être l'avenir de l'assistant personnel et de son activité associée.

Mon hypothèse est qu'il s'agit d'un modèle de plateforme de nouvelle génération.

Évaluons le paysage technique actuel et comment il pourrait évoluer à l'avenir en fonction des éléments techniques qui existent aujourd'hui.

## Plateforme ? 

Revoyons les principes fondamentaux des plateformes, avant d'aller plus loin.

### Économie de plateforme : Des pipelines aux plateformes

Une **plateforme** connecte efficacement **producteurs** et **consommateurs**, leur permettant de générer de la valeur à travers leurs interactions à grande échelle.
Bien que cela puisse sembler être un concept familier, **[Sangeet Paul Choudary](https://en.wikipedia.org/wiki/Sangeet_Paul_Choudary)** introduit une perspective supplémentaire : la notion de **pipeline**.

Son modèle est présenté dans l'article HBR : [Pipelines, Platforms, and the New Rules of Strategy](https://hbr.org/2016/04/pipelines-platforms-and-the-new-rules-of-strategy),

En un coup d'œil, dans le contexte commercial, un pipeline représente une transformation linéaire et unidirectionnelle qui permet à un producteur de créer de la valeur et de la livrer au consommateur.
Essentiellement, c'est le **système traditionnel** des fournisseurs de biens ou de services.
Par exemple, comme décrit dans le document HBR, **Apple** fournit de la valeur en créant des produits et en les vendant aux consommateurs, passant d'un ensemble de composants à un produit fini à travers une série de pipelines.
Ces pipelines ont été révolutionnés par la technologie, impactant trois piliers principaux :

1. **Faciliter la production de masse**.
2. **Encourager la consommation de masse** (par exemple, l'influence de la télévision sur le comportement des consommateurs).
3. **Faciliter les échanges et transactions internationaux**, favorisant la connectivité du système.

Cependant, l'avènement d'**Internet** et de la **technologie numérique** a fait évoluer davantage ces piliers :

- Les **outils de production** peuvent être plus facilement distribués. Auparavant, pour produire de l'information, il fallait être un journal. Maintenant, n'importe qui peut produire de l'information sur des plateformes comme Twitter ou LinkedIn.
- La **numérisation** a apporté des modèles d'utilisation personnalisés, offrant aux consommateurs des produits sur mesure.
- Internet a influencé la **consommation de masse** en affectant les prix, comme on le voit avec des géants comme **Amazon** et **Alibaba**.

**L'idée principale de la plateforme** est de transformer la **création de valeur**.
Sa valeur réside dans la liaison efficace des producteurs et des consommateurs.
Prenant l'exemple d'Apple dans l'article HBR, la **plateforme App Store** facilite la production d'applications de masse pour une large base de consommateurs.
Elle passe de la production de masse liée à la consommation de masse à la production distribuée connectée à la consommation personnalisée.
La plateforme fournit une interface facilitant l'intégration de nouveaux producteurs tout en assurant également la **gouvernance** en mettant en œuvre des règles pour les producteurs et les consommateurs, assurant, par exemple, que les applications de l'App Store sont sûres pour les utilisateurs.

![](/assets/platform_pipeline.png)

### La plateforme pour LLM ? 

Compte tenu des définitions précédentes, la notion de _pipeline_ dans le contexte des LLM peut faire référence à la création de valeur qui résulte de la transformation des données en un produit fini, comme un agent conversationnel.

Cependant, le concept de _plugin_ est quelque peu disruptif, rappelant ce que l'_App Store_ offre.
Les plugins servent d'interfaces au produit, rationalisant les interactions structurées entre utilisateurs et fournisseurs via l'agent.
Cette évolution vers la production de contenu de masse positionne ChatGPT comme une plateforme.

Dans le cas d'utilisation discuté précédemment, de nombreux fournisseurs jouent un rôle dans la satisfaction des besoins de l'utilisateur.
Les magasins et les marques fonctionnent comme des **fournisseurs de contenu**.
Ils offrent des données sur leurs marchandises, détaillant le prix, la disponibilité et la localité.

Pour illustrer la nécessité d'une implémentation double basée sur un type de **Représentation Intermédiaire**, les plugins fonctionnent en générant une représentation intermédiaire qui capture à la fois les engagements API structurés et les explications conviviales pour l'utilisateur.
Cette combinaison garantit un flux de données continu, le rendant intelligible à la fois pour la plateforme et ses utilisateurs.

En conclusion :

Les plugins permettent une **production (de contenu) de masse**, qui, combinée à la **consommation de masse** des individus s'appuyant sur un assistant personnel, donne naissance à une nouvelle **plateforme**.

## Nouveaux défis :
À mesure que cet écosystème s'étend, de nouveaux défis similaires aux défis de référencement rencontrés par les moteurs de recherche surgiront.
Le concept d'**Ingénierie de Prompts et d'Optimisation Automatique des Choix** deviendra essentiel.
Cela implique d'affiner la façon dont les questions sont traitées et répondues par les assistants alimentés par les LLM, garantissant précision, pertinence et efficacité dans la réponse aux demandes des utilisateurs.

Le défi à venir pour les fournisseurs de contenu sera d'être sélectionnés par l'intelligence artificielle pour résoudre les problèmes présentés par le client.
Cela pourrait nécessiter une compréhension approfondie du fonctionnement fondamental de la plateforme.

## Conclusion :

![](/assets/chatgpt-platform_small.webp)

L'avènement des assistants personnels alimentés par les LLM comme ChatGPT annonce une nouvelle ère dans l'économie de plateforme, fusionnant les domaines numériques et physiques dans un écosystème centré sur l'utilisateur.

En intégrant les vitrines numériques, les fournisseurs de services et les utilisateurs, un nouveau niveau de création et d'échange de valeur est réalisé, redéfinissant la façon dont nous interagissons avec le monde numérique qui nous entoure.

Pourquoi s'appuyer sur une plateforme et ne pas devenir une plateforme est un nouveau défi pour les fournisseurs ? Le problème est qu'aujourd'hui, seuls les géants du numérique possèdent les ressources pour exécuter ces modèles à grande échelle.