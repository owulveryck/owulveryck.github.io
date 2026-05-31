---
title: "L'agentic mesh : automatisation cognitive à l'échelle"
date: 2026-05-31T10:00:00+02:00
lastmod: 2026-05-31T10:00:00+02:00
images: [/assets/agenticmesh/poster-agent-mesh.fr.svg]
draft: false
keywords: []
summary: "Convictions pour concevoir les systèmes agentiques de demain — un système de maillage d'agents susceptible d'apporter de la valeur à l'échelle, empruntant ses principes au data mesh."
tags: ["AI", "agents", "architecture", "agentic-mesh"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

On voit fleurir aujourd'hui beaucoup d'initiatives autour de l'agentic. La plupart gravitent autour des systèmes construits par les géants de l'IA (Anthropic, Google, OpenAI) et se résument souvent à pousser des directives en langage naturel vers un orchestrateur sur étagère. On spécialise un système intégré comme Claude Code à coup de fichiers markdown, de skills et d'outils. Des frameworks comme BMAD illustrent bien cette approche : ils transforment un tel système en une équipe de développement mixte IA/humain. C'est utile, c'est rapide, c'est un bon point de départ.

Mais pour réellement tirer la pleine valeur de l'IA, les agents vont devoir s'inviter dans les **processus métier**. Leur apport ne se limite pas à l'automatisation cognitive : il est dans la capacité à **prendre certaines décisions** pour permettre au processus de gagner en efficacité, pas seulement en efficience. Et ceci ne pourra être pleinement atteint que lorsque nous disposerons d'un **écosystème d'agents** qui discutent entre eux et collaborent avec les humains.

Non pas quand les agents feront notre travail, mais quand on aura **bien outillé le système d'information** pour que les agents fassent *leur* travail, et que nous, humains, fassions le nôtre.

Cet article propose un système de **maillage d'agents** (l'agentic mesh) susceptible d'apporter de la valeur à l'échelle. Beaucoup de ses principes, si ce n'est tous, sont empruntés au **data mesh**, un paradigme dont les idées étaient remarquables mais qui n'a malheureusement pas connu le succès qu'il aurait mérité. Selon moi, cet échec ne vient pas de la qualité de ses idées ni de sa conception. Par conséquent, beaucoup peuvent être recyclées ici (et c'est ce que nous allons faire).

> **Note.** Cet article a été co-rédigé avec une IA. Je suis le pilote : je déclare les intentions, les idées, et je revois l'ensemble du document. La rédaction à proprement parler (la plume) a été faite par un robot. Mon objectif est de partager ces idées pour ouvrir une discussion, pas d'écrire une master-piece technique qui devienne une référence stylistique. D'ailleurs, en fin de document, vous trouverez un lien vers le markdown source pour le faire lire à votre IA préférée et discuter avec elle des différents concepts.

---

## L'agent, cet élément fondamental des SI du futur

Commençons par poser un peu de vocabulaire, car la confusion entre deux concepts proches est à l'origine de beaucoup d'erreurs d'architecture.

Nous définissons :

- un **outil** (*tool*) comme une fonction déterministe : même entrée, même sortie. Il exécute programmatiquement une opération via une API ou un protocole comme MCP. C'est idempotent, prévisible, testable.
- un **agent** comme une entité non-déterministe : il travaille et décide en autonomie, sur la base d'une **intention**, en fonction d'un **contexte** d'exécution. Son résultat est statistique (il dépend du contexte au moment de l'appel).

Les conséquences de ces définitions sont importantes :

- l'outil n'est qu'un support d'action (il exécute, il n'interprète pas) ;
- la décision que prend un agent à partir des mêmes outils peut varier en fonction du contexte, de l'état du monde, et de l'intention qui lui a été confiée.

Par exemple : appeler une API météo renvoie toujours la même donnée pour une requête identique, c'est un outil. Décider si l'on reporte une livraison en fonction de la météo, du délai contractuel, de la priorité client et des ressources disponibles, c'est une décision d'agent.

> **L'outil est déterministe. L'agent ne l'est pas : son résultat est statistique et dépend du contexte d'exécution.**

![Agent vs Tool - poster](/assets/agenticmesh/poster-agent-vs-tool.fr.svg)

Cette distinction n'est pas qu'académique. Elle a des conséquences directes sur la gouvernance, sur les SLO, sur l'ownership et sur la façon de tester et de faire confiance à un système. On ne gouverne pas un agent comme on gouverne une API. Le contrat n'est pas un schéma d'entrée/sortie (c'est une **intention et un périmètre de décision**).

---

## Le paradigme dominant : l'agent comme directive sur étagère

Aujourd'hui, lorsqu'une entreprise intègre de l'IA dans son SI, le réflexe dominant n'est pas tant la centralisation que **l'utilisation de systèmes agentiques sur étagère** (Claude Code équipé de frameworks comme BMAD, Cursor avec ses règles, ou des plateformes comme n8n avec des nodes IA). Dans ce paradigme, on ne sépare pas la couche d'exécution des agents : les agents sont vus comme de **simples directives en langage naturel** appliquées à des agents génériques pour les spécialiser.

Cette approche est séduisante. Elle promet une mise en œuvre rapide, sans investissement d'ingénierie significatif. Elle suffit pour des prototypes, pour des automatisations ponctuelles, pour valider qu'un cas d'usage est traitable par IA.

Elle atteint cependant rapidement ses limites dès qu'on cherche à passer en production sur des cas d'usage complexes.

### Ce que masque l'apparente simplicité

Un agent n'est pas qu'un prompt. Un agent qui fonctionne en production est un système logiciel qui doit gérer :

- du **parallélisme contrôlé** entre sous-agents qui travaillent simultanément ;
- des **boucles de feedback structurées** entre un reviewer et les producteurs, avec compteur de retries et filtrage des issues ;
- des **sorties structurées et validées** programmatiquement entre chaque étape ;
- du **prompt caching** partagé entre invocations pour maîtriser les coûts ;
- un **état mutable partagé** typé et protégé contre les accès concurrents ;
- de l'**observabilité fine** (tokens consommés par agent, latence, taux de succès, traçabilité des décisions).

Pour illustrer concrètement, voici un comparatif issu d'un projet réel d'orchestrateur multi-agents (un système de génération de présentations à partir d'une demande en langage naturel) :

| Critère | Implémentation en code natif | Système sur étagère + directives |
|---|---|---|
| Parallélisme fin | Goroutines + sémaphore, contrôle de la concurrence | Limité, séquentiel ou parallèle simple |
| Boucles de feedback | Issues typées, retry ciblé sur sous-ensemble | Conversationnel, fragile dans la durée |
| Sorties structurées | JSON schema strict, validation programmatique | Implicite, dépendante du modèle |
| Prompt caching | Partagé entre sous-agents, coûts maîtrisés | Non disponible ou non partagé |
| État inter-étapes | Typé, mutex, testable | Vit dans le contexte conversationnel |
| Observabilité | Métriques par agent, issue log complet | Limitée, agrégée |
| Maintenabilité | Typée, testable, évolutive | Fragile aux changements de modèle |

Le système sur étagère reste pertinent pour un **prototype rapide** ou un pipeline linéaire simple. Il ne tient pas la charge d'un système de production complexe (non parce qu'il est mauvais, mais parce qu'il déplace la complexité dans le prompt plutôt que dans le code, et que le prompt est un endroit fragile pour gérer cette complexité).

### L'ingénierie logicielle ne se trouve pas que dans les outils

C'est la conviction fondatrice qui distingue l'agentic mesh des approches courantes : **un agent de production est un produit d'ingénierie logicielle, pas un assemblage de directives**. L'ingénierie ne se trouve pas seulement dans les outils que l'agent utilise (elle se trouve aussi dans la façon dont l'agent lui-même est construit, testé, déployé, observé et gouverné).

Cette conviction a une conséquence pratique : pour construire un agentic mesh, il faut traiter les agents comme on traite les autres composants critiques du SI. Avec une plateforme dédiée, des contrats explicites, une ownership claire et une rigueur d'ingénierie comparable à celle qu'on applique aux microservices ou aux data-products.

---

## Les quatre piliers de l'agentic mesh

Pour que chaque agent puisse offrir des solutions adaptées et validées dans un domaine métier particulier, l'agentic mesh repose sur quatre piliers :

1. **La plateforme digitale** donne aux agents l'accès aux capacités d'action et à l'information du SI, dans le respect des règles de gouvernance.
2. **Les domaines** ancrent chaque agent dans son métier : c'est le domaine qui définit l'intention de l'agent, qui en assume la responsabilité, et qui valide la qualité de ses décisions.
3. **L'agent en tant que produit** (IA-as-a-product) impose de traiter l'agent avec la même rigueur qu'un produit logiciel : contrat, cycle de vie, traçabilité, interopérabilité.
4. **La gouvernance automatisée** agit comme un enabler : elle permet la construction et la mise en production des agents dans des conditions de sécurité et de fiabilité exigées par l'organisation, et facilite les échanges entre agents et entre agents et humains.

### Pilier 1 : La plateforme digitale : le vrai MVP

Le premier réflexe d'un client qui veut "mettre de l'IA" dans son SI est de commencer par l'agent. C'est une erreur de séquençage.

**Le vrai MVP, c'est la plateforme.** Avant de déployer des agents sophistiqués, il faut que le SI soit capable d'exposer de manière propre, sécurisée et automatisable deux types de capacités fondamentales :

- les **capacités d'accès à l'information** : comprendre un contexte, interroger des données métier, lire l'état du système ;
- les **capacités d'action** : modifier un statut, déclencher un processus, écrire dans un système.

Ces capacités sont exposées dans le respect des **règles de gouvernance** : habilitations, traçabilité, conformité. C'est ici que plateforme et gouvernance automatisée se rejoignent : la plateforme n'expose pas des accès bruts, elle applique programmatiquement les règles que la gouvernance a définies.

Ce découplage **lecture / écriture** n'est pas nouveau (c'est un pattern éprouvé, CQRS, event sourcing, réinterprété pour l'ère agentique). Il a une propriété cruciale : **il apporte de la valeur indépendamment des agents**. Une plateforme qui expose proprement ses capacités est utile pour les humains, pour les intégrations, pour les automatisations classiques, et seulement ensuite pour les agents.

C'est ce qui rend le MVP vendable : il n'est pas un pari sur l'IA. C'est un investissement d'architecture sain, dont les agents seront les premiers bénéficiaires.

![Fondations du socle - MVP plateforme + POC agent](/assets/agenticmesh/poster-fondations-socle.fr.svg)

**Ce que valide le MVP plateforme :**
- Peut-on exposer les bonnes informations pour qu'un agent comprenne un contexte ?
- Peut-on exposer les bonnes actions pour qu'un agent agisse de manière sécurisée ?
- Quelle granularité d'API ? Quel protocole (REST, MCP, events) ?
- Comment gérer les habilitations quand c'est un agent qui appelle, et non un humain ?

### Pilier 2 : Les domaines : l'agent appartient à son métier

L'agentic mesh emprunte au Domain-Driven Design d'Eric Evans la notion fondamentale de **domaine**. Un domaine est un découpage logique de l'organisation qui rassemble une activité métier cohérente, son vocabulaire propre (*ubiquitous language*), et les frontières dans lesquelles ses modèles ont un sens précis (*bounded context*).

Dans le paradigme de l'agentic mesh, **un agent appartient à un domaine** (exactement comme un microservice ou un data-product appartient à son domaine). Il n'existe pas de "domaine agentique" séparé du métier : il y a des domaines métier qui possèdent et opèrent leurs agents, au même titre qu'ils possèdent leurs données et leurs services.

Cette appartenance n'est pas qu'organisationnelle. Elle structure :

- l'**intention** de l'agent (exprimée dans le vocabulaire du domaine) ;
- le **périmètre de décision** (délimité par les frontières du domaine) ;
- l'**ownership** (un agent product owner, membre du domaine, en assume la responsabilité) ;
- les **règles de gouvernance locales** (qui peut invoquer l'agent, dans quelles conditions, avec quel niveau d'autonomie).

Dans un domaine, on distingue deux catégories de participants autour d'un agent :

- les **producteurs de contexte** (systèmes sources, capteurs, APIs, autres agents) qui alimentent l'agent en informations fiables ;
- les **consommateurs de décisions** (processus aval, humains, agents d'autres domaines) qui exploitent les décisions et actions produites par l'agent.

Un participant peut être simultanément producteur de contexte et consommateur de décisions.

Ce que change fondamentalement ce modèle, c'est que **la responsabilité de l'agent (sa qualité de décision, son périmètre, son évolution) revient au domaine métier qui en exploite la valeur**, et non à une équipe IA centrale détachée du métier.

![Découplage et domaines métier - du monolithe à l'agent-produit](/assets/agenticmesh/poster-decouplage-domaines.fr.svg)

### Pilier 3 : L'agent est un IA-as-a-product

Pour comprendre la place de l'agent dans le mesh, il faut d'abord clarifier deux façons radicalement différentes d'intégrer l'IA dans un SI.

**Un IA-product** est un produit numérique qui résout un problème utilisateur en utilisant de l'IA en interne. L'utilisateur consomme le produit, pas l'IA. L'IA est un moyen, cachée dans la mécanique. Un système de recommandation e-commerce, un assistant de support client, un outil de génération de contenu (ce sont des IA-products). Leur valeur se mesure à l'expérience utilisateur du produit, pas à la qualité brute de l'IA qui le propulse.

**Un IA-as-a-product** est l'inverse : l'IA *est* le produit. Ce qui est exposé, contractualisé et consommé, c'est la capacité d'IA elle-même. Les consommateurs ne sont pas des utilisateurs finaux (ce sont d'autres systèmes (humains ou agents) qui invoquent cette capacité dans leurs propres processus).

Dans un IA-product, l'IA est un **moyen interne**. Dans un IA-as-a-product, l'IA est **l'interface contractuelle**.

![IA-product vs IA-as-a-product](/assets/agenticmesh/ia-product-vs-as-a-product.fr.svg)

**Un agent est un IA-as-a-product.** C'est même la forme la plus aboutie : un agent est une capacité d'IA autonome, qui décide et agit sur la base d'une intention, exposée comme un produit consommable. Considérer l'agent comme un produit n'est pas une nuance (c'est la condition pour qu'il puisse entrer dans un maillage).

> *"On a déployé un assistant IA pour notre support client, on est prêts pour l'agentic mesh."*
> 
> Non (vous avez un IA-product). Pour entrer dans le mesh, il faudrait que la capacité de décision de cet assistant soit exposée, contractualisée et invocable par d'autres domaines. C'est un changement de nature, pas une extension.

Le maillage se construit avec des agents (donc avec des IA-as-a-product). Sans capacité agentique exposée et contractualisée, il n'y a rien à mailler.

#### Les 7 affordances d'un agent

Pour qu'un agent soit véritablement un produit, il doit offrir un ensemble de capacités qui vont bien au-delà de "l'agent qui répond à des questions". Ces capacités (que nous appelons **affordances**, en référence aux travaux de Zhamak Dehghani sur le data-as-a-product) définissent ce qu'un agent **fait**, plutôt que ce qu'il **est**.

**Affordance 1 : Exposer des décisions et des actions**
L'agent expose ses capacités via des interfaces clairement définies. Ce qu'il peut décider, ce qu'il peut faire, dans quelles conditions (tout cela est explicite et invocable de manière structurée, par des humains comme par d'autres agents (A2A, MCP)).

**Affordance 2 : Consommer du contexte**
Pour décider, l'agent consomme du contexte depuis des sources variées : APIs de la plateforme, autres agents, données métier, état du système. Ces sources sont documentées et font partie du contrat du produit.

**Affordance 3 : Raisonner et décider**
C'est le cœur du produit (la logique d'orchestration interne : le modèle utilisé (ou les modèles, car un agent complexe combine souvent plusieurs modèles adaptés à chaque sous-tâche), les prompts versionnés, les outils disponibles, les boucles de feedback, la gestion de l'état). C'est précisément cette logique qui mérite d'être traitée comme du code, pas comme une directive.

**Affordance 4 : Être découvrable et compréhensible**
L'agent est référencé dans un registre d'agents. Son intention, son périmètre de décision, ses conditions d'invocation et ses limites connues sont documentés et accessibles (aux humains comme aux orchestrateurs automatiques).

**Affordance 5 : Gérer son cycle de vie**
L'agent doit pouvoir évoluer sans interrompre ses consommateurs. Cela implique la gestion de versions des prompts et des modèles, des stratégies de dépréciation, et des mécanismes de mise à jour qui n'impactent pas les interfaces exposées.

**Affordance 6 : Tracer les décisions**
Toute décision prise par l'agent est traçable : contexte d'entrée, raisonnement suivi, action entreprise, résultat observé, tokens consommés, modèle utilisé. Cette traçabilité est à la fois un outil d'auditabilité (conformité, régulation) et un outil d'amélioration continue (identifier les cas où l'agent se trompe et calibrer son périmètre).

**Affordance 7 : Être gouvernable**
L'agent offre des capacités de pilotage transverse : politique d'accès (qui peut l'invoquer ?), niveau d'autonomie (quand escalade-t-il vers un humain ?), kill switch (peut-on le désactiver en urgence ?), et gestion des données personnelles qu'il traite.

![Les 7 affordances d'un agent](/assets/agenticmesh/affordances.fr.svg)

#### Le contrat de l'agent : condition du maillage

Si l'agent est un produit, il lui faut un **contrat**. Sans contrat explicite, pas de maillage fiable (on retombe dans le couplage implicite, la découverte manuelle, et la confiance interpersonnelle qui ne passe pas à l'échelle).

Ce contrat doit répondre à des questions précises : qui est cet agent ? Que sait-il faire ? Comment l'invoquer ? Quelles sont ses conditions de sécurité ? Quels formats accepte-t-il ? Ce n'est pas un document de conception interne (c'est une **interface publique**, lisible par des humains et parseable par des systèmes).

L'écosystème n'a pas attendu l'agentic mesh pour formaliser ce besoin. Le protocole **A2A** (*Agent-to-Agent*), porté par Google et soutenu par plus de cinquante acteurs de l'industrie, propose un artefact concret : l'**Agent Card**.

#### L'Agent Card du protocole A2A

L'Agent Card est un document JSON qui fonctionne comme la carte d'identité d'un agent dans un réseau A2A. Elle est publiée à une URL standardisée (`/.well-known/agent-card.json`, suivant la convention RFC 8615), ce qui permet à n'importe quel client ou agent de la découvrir automatiquement par un simple `GET`.

Ce qui rend l'Agent Card intéressante comme modèle de contrat, c'est qu'elle couvre l'ensemble des dimensions nécessaires au maillage :

**Identité et intention** (L'Agent Card expose un `name`, une `description` et une `url` de service). Ce sont les champs obligatoires minimaux. Le `provider` (organisation, URL) permet de rattacher l'agent à son domaine d'ownership. La `version` de l'agent et la `protocolVersion` garantissent la compatibilité entre consommateur et fournisseur.

**Capacités techniques** (Le bloc `capabilities` déclare ce que l'agent supporte au-delà du simple appel synchrone : `streaming` (réponses en flux via SSE), `pushNotifications` (callbacks asynchrones), `stateTransitionHistory` (traçabilité des transitions d'état d'une tâche)). C'est une négociation de capacités explicite, pas une découverte par essai-erreur.

**Compétences** (Le tableau `skills` est le cœur fonctionnel de la carte. Chaque skill a un `id`, un `name`, une `description`, des `tags` de catégorisation, et optionnellement des `examples`) des exemples de prompts ou d'entrées que le skill sait traiter. Chaque skill peut déclarer ses propres `inputModes` et `outputModes` (types MIME : `text/plain`, `application/json`, `image/png`...), ou hériter des modes par défaut de l'agent.

**Sécurité** (Les `securitySchemes` suivent la convention OpenAPI : bearer tokens, API keys, OAuth 2.0, OpenID Connect. Le champ `security` déclare quels schémas sont requis pour invoquer l'agent. Pas de surprise à l'exécution) les conditions d'accès sont dans le contrat.

**Découvrabilité étendue** (Le flag `supportsAuthenticatedExtendedCard` indique qu'une version plus détaillée de la carte est disponible derrière authentification, pour les informations sensibles ou les capacités privées).

L'ensemble est contraint à **10 Ko maximum**, ce qui force la concision et garantit une découverte rapide.

#### Ce que l'Agent Card ne couvre pas (encore)

L'Agent Card est un excellent socle de contrat d'interopérabilité, mais elle ne couvre pas tout ce dont un agent-as-a-product a besoin dans le contexte d'un agentic mesh :

- **Le contrat de confiance** (taux de décisions correctes, conditions d'escalade vers un humain, SLO décisionnels). L'Agent Card dit *ce que* l'agent sait faire, pas *à quel point* on peut lui faire confiance.
- **La gouvernance** (politique de rétention des données, traçabilité des décisions, conformité réglementaire). Ce sont des préoccupations organisationnelles qui dépassent le périmètre d'un protocole d'interopérabilité.
- **L'ownership** (à quel domaine métier l'agent est rattaché, qui assume la responsabilité de ses décisions).

Ces dimensions complémentaires devront être portées par des extensions du contrat ou par des artefacts de gouvernance propres à chaque organisation. L'Agent Card pose la fondation technique ; le contrat complet d'un agent-as-a-product y ajoute les dimensions produit et organisationnelle.

### Pilier 4 : La gouvernance automatisée : l'enabler du maillage

La gouvernance de l'agentic mesh n'est pas un cadre de contraintes imposé de l'extérieur. C'est un **enabler** qui rend le maillage possible en agissant sur deux axes.

**Enabler technique** : la gouvernance automatisée fournit les conditions pour que les agents puissent être construits et mis en production avec le niveau de sécurité et de fiabilité exigé par l'organisation. Cela inclut les standards de traçabilité et d'explicabilité, les obligations de conformité (AI Act, régulation sur les systèmes de décision automatisée), les protocoles d'interopérabilité (A2A, MCP), et les métriques de qualité de décision. Ces règles ne sont pas documentaires : elles sont **appliquées programmatiquement** par la plateforme. C'est ce qui justifie le terme *automatisée* : la gouvernance est dans le code, pas dans un wiki.

**Enabler d'échanges** : la gouvernance facilite les interactions entre agents et entre agents et humains. Elle définit qui peut invoquer quel agent, dans quelles conditions, avec quel niveau d'autonomie. Elle gère les mécanismes d'escalade vers un humain, les boucles de feedback, et les politiques de rétention des données. Elle rend le maillage fluide en supprimant les frictions d'intégration.

Cette gouvernance s'applique à deux niveaux. Au niveau central, elle pose le cadre constitutionnel : les règles qui s'appliquent à tous les agents du mesh. Au niveau du domaine, chaque équipe applique ce cadre dans son périmètre en l'adaptant à ses contraintes opérationnelles.

Un exemple concret : la gouvernance centrale imposera que tout agent qui prend des décisions ayant un impact financier supérieur à un certain seuil dispose d'un mécanisme d'escalade vers un humain. La plateforme appliquera cette règle automatiquement. Le domaine logistique définira ensuite sa propre politique (seuil précis, délai d'escalade, personne compétente) en fonction de ses contraintes.

---

## De l'agent prototype au maillage : une trajectoire en trois temps

La mise en place d'un agentic mesh ne se décrète pas. Elle se construit par itérations successives.

![Trajectoire en 3 temps](/assets/agenticmesh/trajectory.fr.svg)

### Temps 1 : Construire le socle : MVP plateforme + POC agent

**L'objectif de cette phase n'est pas de faire un bon agent. C'est de valider les hypothèses fondamentales.**

On choisit un cas d'usage métier réel, on expose les capacités minimales nécessaires sur la plateforme (lecture + écriture), et on construit un agent avec des bouts de ficelles (pas besoin qu'il soit à l'état de l'art). À ce stade, un système sur étagère peut tout à fait convenir : on cherche à valider, pas à produire.

Ce que le **MVP plateforme** valide :
- Peut-on exposer les informations nécessaires à la décision ?
- Peut-on exposer les actions nécessaires à l'exécution ?
- Quels sont les problèmes de qualité de données, de granularité d'API, de gouvernance des accès ?

Ce que le **POC agent** valide :
- Si on expose les bons services, peut-on envisager de l'automatisation cognitive ?
- L'agent apporte-t-il une valeur mesurable sur ce cas d'usage ?
- Quelles sont les conditions dans lesquelles l'agent se trompe ?
- A-t-on besoin de structurer l'information en amont ?

À ce stade, l'agent est expérimental. Il est couplé à l'orchestrateur, son périmètre est flou, son ownership est informel. C'est normal.

### Temps 2 : Faire de l'ingénierie : du prototype au produit

Le POC a validé le concept. On sait que l'automatisation cognitive est possible sur ce périmètre. Il est temps de faire évoluer l'architecture et de **passer à l'ingénierie logicielle**.

Cette phase consiste à transformer l'agent prototype en une **architecture agentique structurée** :

![Architecture agentique structurée](/assets/agenticmesh/architecture-agentique.fr.svg)

Chaque sous-agent est spécialisé, utilise le modèle adapté à sa tâche (Haiku pour le simple, Sonnet pour le complexe, par exemple), dispose de ses outils propres et de ses boucles de feedback. L'orchestrateur coordonne sans tout faire lui-même, et il est implémenté en **code**, pas en directive.

C'est dans cette phase que se joue le passage du *prompt* au *produit*. Ce qui était un assemblage de directives devient un système logiciel testable, observable, déployable. Les questions qui se posent ne sont plus "comment formuler le prompt" mais "comment versionner les prompts, comment tester les régressions, comment instrumenter les décisions, comment gérer les retries".

C'est aussi à ce stade que certains sous-agents se révèlent **utiles au-delà de leur cas d'usage initial**. Un sous-agent qui résume des conversations clients pour ce processus pourrait servir à d'autres processus. C'est le signal qu'il doit changer de nature : être extrait de l'application, doté de son propre manifeste d'intention, et géré comme un produit autonome dans le domaine qui en exploite la valeur.

### Temps 3 : Construire le maillage : de l'agent-produit à l'écosystème

Les agents sont maintenant gérés comme des produits dans leurs domaines respectifs. Ils exposent leurs capacités via des manifestes standardisés. Ils sont découvrables, invocables, gouvernés.

Le maillage peut alors se construire : des processus métier bout-en-bout qui orchestrent ou chorégraphient des séquences d'agents de domaines différents, en alternant décisions humaines et décisions agentiques selon les besoins.

![Processus métier - séquences Humain-Agent transverses](/assets/agenticmesh/processus-metier-ha.fr.svg)

Les séquences **H-A** (Humain-Agent) ne sont pas des exceptions ou des garde-fous (elles sont une **feature** du maillage). L'agentic mesh n'a pas pour ambition de remplacer les humains. Il a pour ambition de leur permettre de travailler sur les décisions qui requièrent réellement leur jugement, en déléguant les autres.

---

## L'organisation autour du maillage

### L'ownership par domaine

Un des changements organisationnels les plus importants induits par l'agentic mesh est le **déplacement de l'ownership**.

Dans un modèle centralisé, les agents appartiennent à une équipe IA transverse. Dans l'agentic mesh, **chaque agent appartient au domaine métier qui en exploite la valeur**. C'est ce domaine qui :
- définit l'intention de l'agent ;
- assume la responsabilité de la qualité de ses décisions ;
- gère son cycle de vie ;
- répond des incidents qu'il provoque.

Ce changement d'ownership est souvent difficile à faire accepter, parce qu'il implique que les équipes métier acquièrent une compétence nouvelle (non pas en IA, mais en **product management d'agents**). C'est un investissement organisationnel réel, qui demande du temps et de l'accompagnement.

### L'enablement comme fonction transverse

Donner de l'autonomie aux domaines ne signifie pas les laisser réinventer la roue. Une fonction transverse d'enablement est nécessaire pour :

- maintenir et faire évoluer la **plateforme digitale** (les commodités techniques) ;
- définir et faire respecter les **standards de manifeste et d'interopérabilité** ;
- fournir des **patterns d'architecture agentique** réutilisables (orchestration, feedback loops, observabilité) ;
- accompagner les domaines dans le passage de **l'agent prototype à l'agent-produit**.

Le succès de cette fonction se mesure à un indicateur simple : **le temps entre la conception d'un agent et le moment où il est invocable par d'autres domaines**.

---

## Synthèse : le cercle vertueux de l'agentic mesh

L'agentic mesh n'est pas seulement une façon de modéliser et gérer des systèmes d'IA. C'est avant tout une façon de penser **l'agent en tant que produit** (un IA-as-a-product à part entière, conçu avec une rigueur d'ingénierie logicielle). Le produit n'est plus l'application qui utilise l'IA en interne. **La capacité de décision *est* le produit**. Les outils et les APIs ne sont qu'un moyen de l'alimenter.

Ainsi, les indicateurs de niveau de service devront être adaptés :
- les indicateurs de **pertinence** et de **fiabilité** des décisions remplaceront les indicateurs de disponibilité comme métriques primaires ;
- la **traçabilité** des décisions deviendra un SLO à part entière, non une option ;
- le **taux d'escalade** vers des humains sera un indicateur de calibrage de l'autonomie, non un indicateur d'échec.

**L'agentic mesh n'est pas une version améliorée du chatbot d'entreprise.** Le vrai changement qu'il porte est le passage de **l'humain outillé** (qui utilise l'IA comme un instrument) à **l'humain accompagné** (qui collabore avec des agents au sein d'équipes mixtes, tout en gardant le contrôle du processus. Les agents décident dans leur périmètre, mais c'est l'humain qui reste maître du processus global : il valide, arbitre et escalade. C'est cette collaboration sous contrôle humain qui produit la valeur. Ce changement a un impact fort sur la conception de produit, sur l'architecture du SI, et sur l'organisation des équipes).

Il ne faut pas pour autant jeter les pratiques existantes. Si votre SI repose déjà sur des principes d'exposition API propres, un découplage fonctionnel lecture/écriture, et une organisation en domaines métier cohérents (vous êtes déjà sur la bonne trajectoire). La transition vers l'agentic mesh sera d'autant plus naturelle que votre plateforme propose déjà le socle sur lequel les agents peuvent décider.

---

![Agentic Mesh - les 4 piliers du maillage à valeur](/assets/agenticmesh/poster-agent-mesh.fr.svg)


Pour résumer, le **cercle vertueux de l'agentic mesh** consiste en plusieurs étapes :

- **construire la plateforme digitale** qui expose proprement les capacités de lecture et d'écriture du SI ;
- **construire un agent** sur un cas d'usage réel, avec une vraie rigueur d'ingénierie logicielle (pas comme un assemblage de directives, mais comme un produit) ;
- **rattacher cet agent à son domaine métier** qui en assume l'ownership et la responsabilité de ses décisions ;
- cet agent devient un **nœud du maillage** grâce aux standards d'interopérabilité, de gouvernance automatisée et de contrats d'intention partagés ;

…tout ceci dans le but de produire des **processus métier automatisés à forte valeur ajoutée**, où humains et agents collaborent naturellement, chacun dans son périmètre de décision.

> *Les agents sont les points de communication entre domaines. La valeur naît du maillage.*

---

## Annexe : Poster de synthèse

Le poster ci-dessous résume l'ensemble des convictions en une vue unique, organisée en 4 zones : [poster complet (HTML interactif)](/assets/agenticmesh/poster-architecture-agentique.fr.html).

---

## Source

Cet article est disponible en markdown pour le faire lire à votre IA préférée et discuter des concepts qu'il contient : [Télécharger le markdown source](https://github.com/owulveryck/owulveryck.github.io/tree/master/content/post/20260531-agentic-mesh.fr.md).
