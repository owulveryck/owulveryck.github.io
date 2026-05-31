---
title: "L'agentic mesh en pratique : anatomie d'un agent-produit"
date: 2026-05-31T14:00:00+02:00
lastmod: 2026-05-31T14:00:00+02:00
images: [/assets/agenticmesh/architecture-agentique.fr.svg]
draft: true
keywords: ["agentic mesh", "multi-agent", "Go", "Google Slides", "A2A"]
summary: "Comment j'ai construit un système multi-agents en Go pour générer des présentations Google Slides à partir de templates préformatés (non pas pour présenter, mais pour convaincre). Retour d'expérience et illustration concrète des principes de l'agentic mesh."
tags: ["AI", "agents", "architecture", "agentic-mesh", "Go", "retour-experience"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

Je suis consultant, et je construis régulièrement des présentations avec Google Slides. Mon équipe communication a créé des dizaines de templates préformatés (des slides pensées pour **convaincre**, pas juste pour présenter). Le problème : choisir les bons slides pour illustrer le bon discours prend du temps, et les remplir mécaniquement n'a aucune valeur ajoutée. J'ai construit un système multi-agents pour automatiser cette partie et me concentrer sur ce qui compte : le discours et l'appropriation.

Ce projet (**agentigslide**) est aussi une application concrète des principes de l'agentic mesh que j'ai décrits dans [l'article précédent](/fr/2026/05/31/lagentic-mesh-automatisation-cognitive-à-léchelle.html). Là-bas, j'ai posé un cadre conceptuel : quatre piliers, une trajectoire en trois temps, des convictions sur ce que devrait être un agent-produit. Ici, je raconte comment ces principes se sont incarnés dans du code réel, des décisions d'architecture documentées par 16 ADR, et un outil qui fonctionne en production.

> **Note.** Cet article a été co-rédigé avec une IA. Je suis le pilote : je déclare les intentions, les idées, et je revois l'ensemble du document. La rédaction à proprement parler (la plume) a été faite par un robot. Mon objectif est de partager ces idées pour ouvrir une discussion, pas d'écrire une master-piece technique qui devienne une référence stylistique. Cette version est conçue pour les humains ; si vous préférez une version adaptée aux IA, le [markdown source](https://raw.githubusercontent.com/owulveryck/owulveryck.github.io/refs/heads/master/content/post/20260601-agentic-mesh-in-practice.fr.md) est disponible.

---

## Convaincre, pas présenter

Tous les outils de création de slides résolvent le mauvais problème. Gamma, Beautiful.ai, Pitch : ils génèrent des slides visuellement correctes. Mais ils produisent des **présentations**, pas des **présentations de conviction**. La différence est fondamentale.

Convaincre, c'est d'abord un **discours**. Les slides en sont le complément, pas la redite. Un consultant qui prépare un pitch ne part pas d'un outil de génération de slides, il part de son argumentaire, de la structure de sa démonstration, des idées qu'il veut ancrer chez son interlocuteur. Les slides ne sont que le support visuel de ce raisonnement.

Dans ce contexte, mon équipe communication a produit un catalogue de ~300 slides préformatées dans Google Slides (des slides conçues par des professionnels de la communication visuelle, pensées pour illustrer des concepts spécifiques du conseil IT : slides de conviction, de cadrage, de comparaison, de processus). Ce catalogue est un **actif de marque** qui encode des conventions visuelles et rhétoriques construites sur des années d'expertise.

La **valeur du consultant** est de choisir dans ce catalogue les bons slides pour illustrer les bons concepts de son discours, puis de les compléter avec son contenu. Le choix est stratégique, le remplissage est mécanique. J'ai voulu automatiser le mécanique pour libérer le stratégique.

> **La valeur du consultant est dans le choix et l'appropriation, pas dans le remplissage mécanique.**

Dans le vocabulaire de l'agentic mesh, c'est le **pilier 2 : les domaines**. L'agent que j'ai construit appartient au domaine du conseil : son intention s'exprime dans le vocabulaire métier (structurer un pitch, choisir des slides de conviction), pas dans un vocabulaire technique générique.

---

## Le paysage stratégique : une carte Wardley

Avant de coder quoi que ce soit, j'ai cartographié le paysage avec une **carte Wardley** pour comprendre où se situe la valeur et quelles manœuvres stratégiques sont possibles.

La chaîne de valeur se lit de haut en bas : le **consultant** doit **convaincre un client**. Pour cela il produit un **discours** qu'il externalise dans un **brief structuré** (un fichier markdown). Ce brief alimente l'**orchestration agentique** qui pioche dans le **catalogue de slides préformatés** via un **index sémantique**, pour produire une présentation via les **API Google Slides et Drive**.

### Ce que la carte révèle

**Le catalogue est le fossé invisible.** Positionné en phase Custom, c'est du capital relationnel (non copiable sans l'expertise communication qui l'a produit). Un concurrent pourrait reproduire l'architecture agentique, déployer les mêmes modèles LLM, mais il ne pourrait pas copier ce catalogue sans copier les années d'expertise qui l'ont façonné.

**Le timing est critique.** L'orchestration agentique sort de la phase Genesis pour entrer en phase Custom. Les outils de génération générique (Gamma, Beautiful.ai) sont sous pression Red Queen : ils doivent évoluer en permanence pour ne pas reculer. Leur direction naturelle est d'absorber les capacités agentiques. La fenêtre pour occuper la niche "slides de conviction" est de **12 à 18 mois**.

**Le human-in-the-loop est un choix stratégique, pas une limitation.** Le consultant revoit et s'approprie les slides générées. C'est l'étape qui transforme une présentation correcte en une présentation convaincante. Ce n'est pas une rustine en attendant que l'IA soit "assez bonne" : c'est l'humain qui présente qui est responsable de ce qu'il présente.

**Le non-intrusif comme doctrine.** L'équipe communication continue d'utiliser Google Slides sans apprendre aucun nouvel outil. Le système comprend automatiquement la structure des templates et s'y adapte. Adoption sans friction = adoption réelle.

### Les manœuvres identifiées

- **Land-grab** : occuper le vocabulaire de la "conviction" avant que les outils génériques ne colonisent la niche.
- **Open-source stratégique** : publier le framework agentique (la plomberie se commoditisera) pour déplacer l'attention vers le catalogue (le vrai moat).
- **Strangler-fig sur le brief** : un agent de structuration en amont aide le consultant à structurer son discours : il remplace progressivement la production manuelle du brief, pas le discours lui-même.
- **ILC sur le catalogue** : chaque génération produit des signaux implicites (quels slides sont choisis, lesquels ne le sont jamais, quels concepts n'ont pas de slide). L'équipe communication utilise ces insights pour enrichir le catalogue (une boucle de rétroaction usage → production).

---

## Temps 1 : le prototype monolithique

Mon premier système était un **pipeline monolithique** : un seul appel à Claude recevait le catalogue complet des templates (~60 slides, ~15-20 KB de texte) et la demande de l'utilisateur, et devait en un seul passage analyser la structure, sélectionner les templates, remplir le contenu textuel de chaque champ, et maintenir la cohérence globale de la présentation.

Ça marchait. En tout cas suffisamment pour **valider l'hypothèse fondamentale** : oui, un LLM peut choisir les bons slides dans un catalogue et les remplir de manière pertinente.

C'est exactement le **Temps 1** de la trajectoire de l'agentic mesh : *"L'objectif de cette phase n'est pas de faire un bon agent. C'est de valider les hypothèses fondamentales."* Et c'est ce que le prototype a fait : il a validé que l'automatisation cognitive est possible sur ce cas d'usage.

Mais quatre **limites structurelles** ont émergé rapidement :

1. **Pas de boucle de correction.** Si le modèle choisit un slide à 6 champs pour 3 bullet points, ou dépasse une limite de caractères, l'erreur passe en production. Raisonnement en un seul passage, pas de filet.
2. **Pas de parallélisme.** Écrire le contenu du slide 3 et du slide 7 sont des tâches indépendantes, mais tout est fait dans un seul appel séquentiel.
3. **Inefficacité tokens.** Chaque sous-tâche reçoit l'intégralité du contexte (catalogue complet) alors que seule une fraction est pertinente.
4. **Couplage fort.** Un slide de couverture à 2 champs ne nécessite pas le même modèle qu'un slide de contenu à 6 champs. Mais tout passe par le même appel.

Le prototype avait rempli son rôle. Il fallait maintenant passer à l'ingénierie.

---

## Temps 2 : du prompt au produit

### L'architecture multi-agents

La première décision d'architecture (ADR 001, le 5 mai 2026) a transformé le pipeline monolithique en **quatre agents spécialisés** coordonnés par un orchestrateur en code Go pur (pas d'IA dans l'orchestration).

**L'Outliner** analyse la demande de l'utilisateur et produit un plan de présentation structuré. Point crucial : il ne reçoit **pas** le catalogue de templates. Cette isolation est délibérée : elle force le raisonnement *"de quoi a-t-on besoin ?"* avant *"qu'a-t-on à disposition ?"*, évitant le biais de disponibilité.

**Le Selector** fait le matching entre les besoins identifiés par l'Outliner et les templates disponibles dans le catalogue. Il travaille avec le contexte `itemCount` et `maxItemLength` fourni par l'Outliner pour faire des choix informés.

**Les Writers** génèrent le contenu textuel de chaque slide, en parallèle. Le modèle est sélectionné selon la complexité : **Haiku** pour les slides simples (couverture, intercalaire, ≤ 2 champs), **Sonnet** pour les slides complexes (> 2 champs). L'exécution est en goroutines Go avec un sémaphore configurable.

**Le Reviewer** valide le plan assemblé avant exécution. Il utilise **Opus avec extended thinking** pour une analyse en profondeur. S'il détecte des problèmes (overflow de texte, duplication, template inadéquat), il renvoie un feedback structuré (`ReviewIssue[]`) aux Writers concernés, qui corrigent et resoumettent (maximum 2 itérations pour borner le coût).

C'est le **Temps 2** de l'agentic mesh : *"Chaque sous-agent est spécialisé, utilise le modèle adapté à sa tâche, dispose de ses outils propres et de ses boucles de feedback."*

### Pourquoi du code natif et pas des directives

J'aurais pu utiliser Claude Code avec des serveurs MCP, ou l'Agent SDK d'Anthropic. J'ai choisi Go natif, et ce choix illustre la conviction fondatrice de l'agentic mesh : **un agent de production est un produit d'ingénierie logicielle, pas un assemblage de directives**.

| Critère | Go natif | Système sur étagère + directives |
|---------|----------|----------------------------------|
| Parallélisme fin | Goroutines + sémaphore, contrôle de concurrence | Limité, séquentiel ou parallèle simple |
| Boucles de feedback | `ReviewIssue[]` typés, retry ciblé sur sous-ensemble | Conversationnel, fragile dans le temps |
| Sorties structurées | JSON schema strict, validation programmatique | Implicite, dépendant du modèle |
| Prompt caching | Cache éphémère Vertex AI, partagé entre Writers | Non disponible ou non partagé |
| État inter-étapes | Typé, mutex, testable | Vit dans le contexte conversationnel |
| Observabilité | Métriques par agent, issue log complet | Limitée, agrégée |

L'ingénierie ne se trouve pas seulement dans les outils que l'agent utilise : elle se trouve aussi dans la façon dont l'agent lui-même est construit, testé, déployé, observé et gouverné.

### Les décisions d'ingénierie documentées par les ADR

Le projet a accumulé **16 Architecture Decision Records** en 4 semaines. Chaque ADR est une décision d'ingénierie délibérée, documentée avec son contexte, ses alternatives évaluées, et ses conséquences. Ce n'est pas de la documentation a posteriori, c'est de la **gouvernance en action**.

| ADR | Décision | Concept agentic mesh illustré |
|-----|----------|-------------------------------|
| 001 | Architecture multi-agents | Temps 2 : du prototype au produit |
| 002 | Prompt caching explicite via Vertex AI | Maîtrise des coûts d'un agent-produit |
| 004 | Externalisation des prompts (`go:embed`) | Versionnement et testabilité |
| 005-006 | Chat interactif + mode par défaut | Human-in-the-loop structurel |
| 007 | Architecture A2A (Agent-to-Agent) | Contrat d'interopérabilité, Agent Card |
| 009 | Agent Designer pour les diagrammes | Spécialisation par sous-domaine |
| 010-012 | Pipeline d'édition orchestré | Extension du domaine fonctionnel |
| 015 | Mémoire d'apprentissage par agent | Amélioration continue gouvernée |
| 016 | FormatAgent déterministe | Gouvernance automatisée sans LLM |

---

## Les 7 affordances, concrètement

Dans l'article sur l'agentic mesh, j'ai défini **7 affordances** qu'un agent doit offrir pour être un véritable produit. Voici comment elles se sont incarnées dans agentigslide.

**1. Exposer décisions et actions.** Chaque agent expose ses capacités via un schéma JSON strict imposé par le mécanisme `tool_use` de Claude. Avec l'ADR 007, chaque agent expose aussi une **Agent Card** A2A (un manifeste auto-descriptif publié sur `/.well-known/agent-card.json`).

**2. Consommer du contexte.** Les agents consomment trois types de contexte : l'index sémantique du catalogue (construit une fois, réutilisé à chaque génération), les instructions spécifiques au template (un fichier `PROMPT.md` optionnel), et les fichiers mémoire issus des exécutions précédentes.

**3. Raisonner et décider.** C'est le cœur du produit : Haiku pour les tâches simples, Sonnet pour les complexes, Opus avec extended thinking pour la review. Les prompts sont externalisés via `go:embed` et versionnés avec le code. La boucle Reviewer → Writer est une chaîne de raisonnement structurée, pas une conversation.

**4. Être découvrable.** L'Outliner est déjà déployé comme serveur A2A autonome (`cmd/outliner/main.go`). L'ADR 014 a proposé un pattern de registre d'agents pour la composition dynamique.

**5. Gérer son cycle de vie.** Les modèles sont configurables par agent via des variables d'environnement (`AGENT_OUTLINER_MODEL`, `AGENT_WRITER_MODEL`...). Le mode monolithique original est conservé comme fallback (rétro-compatibilité sans dette).

**6. Tracer les décisions.** Chaque agent rapporte ses tokens (input, output, cache read, cache creation), sa durée, et un issue log complet. L'extended thinking du Reviewer est tracé. Le FormatAgent (ADR 016) journalise chaque correction déterministe appliquée.

**7. Être gouvernable.** `MaxReviewRetries` borne les itérations de correction. `enforceMaxChars()` est un filet de sécurité programmatique (*trust but verify*). La mémoire d'agent (ADR 015) est validée par l'humain avant écriture : l'agent propose des guidelines, l'utilisateur confirme.

---

## Le human-in-the-loop : un choix, pas une rustine

Le mode par défaut d'agentigslide est le **chat interactif** (ADR 005-006). Le consultant décrit ce qu'il veut, l'Outliner propose une structure, le consultant raffine en conversation, et seulement quand le plan est validé, la génération se lance.

Ce n'est pas une béquille. C'est un choix délibéré, ancré dans une conviction : **celui qui présente est responsable de ce qu'il présente**. L'agent produit un support de qualité professionnelle, mais c'est le consultant qui l'examine, l'ajuste, et lui donne la teinte qui fera la différence entre une présentation générique et une présentation convaincante.

Dans l'agentic mesh, j'ai écrit : *"Les séquences H-A (Humain-Agent) ne sont pas des exceptions ou des garde-fous, elles sont une feature du maillage."* C'est exactement ce qui se passe ici. L'humain n'est pas dans la boucle parce que l'IA n'est pas assez bonne. Il est dans la boucle parce que c'est sa valeur ajoutée irremplaçable.

---

## Le catalogue : le fossé invisible

J'insiste sur un point que la carte Wardley rend évident : **le code n'est pas le fossé**. L'architecture multi-agents, les goroutines, le prompt caching : tout cela se commoditisera. Un concurrent pourrait reproduire l'ensemble du système technique.

Ce qu'il ne pourrait pas reproduire, c'est le **catalogue de 300 slides** conçues par des professionnels de la communication visuelle. Ce catalogue n'est pas une collection de templates, c'est un actif de marque qui encode des conventions visuelles et rhétoriques construites sur des années. C'est du capital relationnel, pas du capital technique.

Et la doctrine de **non-intrusivité** renforce ce fossé : l'équipe communication continue de travailler dans Google Slides sans apprendre aucun nouvel outil. Le système analyse automatiquement chaque slide avec Claude Vision, comprend sa structure, identifie les champs éditables, et construit un index sémantique. L'adoption est sans friction parce que l'outil s'adapte aux humains, pas l'inverse.

---

## Vers le Temps 3 : de l'agent-produit au maillage

### A2A : chaque agent expose une Agent Card

L'ADR 007 (9 mai 2026) a été le moment charnière. Le pipeline Go fonctionne bien, mais il a trois limites structurelles qu'aucune optimisation ne peut résoudre :

1. **Les agents ne peuvent pas orchestrer d'autres agents.** Le Selector ne peut pas décider dynamiquement d'appeler un agent de layout puis un agent de design.
2. **Le pipeline est fermé à l'extension externe.** Ajouter un nouvel agent requiert de modifier, recompiler et redéployer le binaire.
3. **L'unité de déploiement est le binaire entier.** Pas de cycle de vie indépendant par agent.

La solution : le protocole **A2A** (Agent-to-Agent, Google, 2025). Chaque agent expose une Agent Card et accepte des Tasks via une API REST standardisée. L'orchestrateur Go reste (déterministe, prévisible) mais il appelle des agents via A2A plutôt que des fonctions Go. C'est la convergence directe avec le **contrat d'interopérabilité** et l'**Agent Card** décrits dans l'article sur l'agentic mesh.

### Du pipeline fermé au réseau composable

Le changement de paradigme le plus profond est celui du catalogue. Aujourd'hui, le catalogue est une **contrainte dure** : si aucun template ne convient, le Selector ne peut rien faire. Avec A2A, il devient un **défaut intelligent avec fallback créatif** : quand aucun template ne convient, le Selector orchestre des sous-agents (agent de layout, agent de design, agent de validation visuelle) pour créer une slide ex nihilo à partir de primitives de design.

Ce changement a une précondition non triviale : la **charte visuelle** de l'équipe communication, aujourd'hui implicite dans les slides, doit devenir explicite, versionnée, testable. L'équipe communication ne produit plus seulement des slides, elle produit des **primitives de design** et des **règles de composition**. C'est un changement de nature dans sa contribution.

### La mémoire comme gouvernance incrémentale

L'ADR 015 a introduit la **mémoire d'apprentissage par agent**. Chaque agent dispose d'un fichier Markdown par template, stocké aux côtés du template et versionné avec git. Les guidelines sont actionnables : *"Sur le slide #42, ne jamais dépasser 120 caractères dans le champ titre (le texte déborde systématiquement)."*

En fin de pipeline, si des erreurs ont été détectées, le système synthétise des guidelines et les **propose à l'utilisateur** (pas d'écriture automatique). L'humain confirme avant que la mémoire ne soit enrichie. C'est l'affordance 7 (être gouvernable) appliquée concrètement : l'agent s'améliore à l'intérieur de limites définies par l'humain.

---

## Ce que le terrain enseigne à la théorie

Après 4 semaines de construction, 16 ADR, et un système qui fonctionne en production, voici ce que le projet agentigslide m'a confirmé (ou appris) par rapport aux principes de l'agentic mesh.

**Le vrai MVP, c'est la plateforme.** Le catalogue de slides, l'index sémantique, les API Google : c'est ça le socle qui a débloqué la valeur. L'agent n'est venu qu'après, et il n'aurait rien pu faire sans cette fondation.

**L'agent n'est pas une directive, c'est un produit d'ingénierie.** La table de comparaison du RATIONALE.md le démontre sans ambiguïté. Les systèmes sur étagère sont un bon point de départ (j'y ai eu recours pour le Temps 1). Mais dès qu'on a besoin de parallélisme contrôlé, de feedback loops typés, de validation inter-étapes et de prompt caching partagé, on fait de l'ingénierie logicielle, pas de l'assemblage de directives.

**Les ADR sont la gouvernance en action.** 16 décisions documentées en 4 semaines. Le contexte, les alternatives évaluées, les conséquences : tout est là. La gouvernance n'est pas dans un wiki, elle est dans le trail de décisions.

**Le fossé est dans le domaine, pas dans la technologie.** Le catalogue est le moat, pas le code Go. C'est le **pilier 2** de l'agentic mesh qui s'exprime : *"la responsabilité de l'agent tombe sur le domaine métier qui tire sa valeur, pas sur une équipe IA centrale."*

**La fenêtre est étroite.** Les outils génériques vont absorber les capacités agentiques. La fenêtre pour occuper la niche "slides de conviction" est de 12 à 18 mois. Le land-grab est urgent.

> **L'agentic mesh n'est pas un diagramme à accrocher au mur. C'est une trajectoire à suivre, et agentigslide en est un point de passage.**
