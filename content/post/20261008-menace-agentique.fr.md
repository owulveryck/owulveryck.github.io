---
title: "La menace agentique"
date: 2026-10-08T10:00:00+02:00
lastmod: 2026-10-08T10:00:00+02:00
images: [/assets/menace-agentique/where-guardrails-act.fr.svg]
draft: false
keywords: ["ingénierie agentique", "sécurité", "garde-fous", "OPA", "Rego", "Persona Guardrail"]
summary: >
  Plus on donne d'autonomie aux agents, plus c'est pratique, et « c'est pratique » est un des mots les plus dangereux du numérique. Il faut protéger les systèmes de la menace humaine, mais aussi de l'agent lui-même. Ces deux menaces n'ont pas la même réponse architecturale : le Persona Guardrail d'Uber d'un côté, une gateway déterministe (PPG) de l'autre, et la distinction entre mesures compensatoires et amplificatrices.
tags: ["AI", "agents", "architecture", "sécurité", "ingénierie-agentique"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

## L'autonomie, c'est pratique

Je suis plus que convaincu que l'**ingénierie agentique** est la discipline qui va permettre de sublimer la valeur de l'IA en entreprise.
J'ai la chance de travailler sur les évolutions des procédés de développement d'assets numériques, logiciels et autres.
Dans ce cadre, une clé est de mettre à disposition des systèmes agentiques qui vont agir avec le **plus d'autonomie possible** pour accomplir des tâches non différenciantes avec robustesse et rapidité (et idéalement à un coût maîtrisé).

La mise en place de **harnais de développement** a ouvert la voie à l'autonomie des agents. En effet, il est désormais possible de discuter pour créer des assets numériques avec Claude Code, Copilot ou autre.
On active un auto-mode et après un cadrage c'est parti.
Le harnais est équipé d'un ensemble d'outils qui permet d'interagir avec l'écosystème. Dans le cas d'une entreprise, interagir avec l'écosystème veut dire récupérer des informations, agir sur des outils ou des processus.

Plus on donne d'outils, plus les agents peuvent travailler en autonomie. Et comme je disais ce matin avec un collègue : plus on donne d'autonomie, plus c'est pratique…. Et **« c'est pratique » est un des mots les plus dangereux** dans l'écosystème numérique…. Car les systèmes qui offrent de la praticité le font en général en échange d'autre chose qui a de la valeur pour eux-mêmes.
Dans le cas des systèmes agentiques, on **échange la délégation de tâches contre du contrôle** : on cède de l'autonomie à l'agent, et on agrandit d'autant la **surface d'attaque**.

Et cette autonomie peut être dangereuse.

## Deux risques, deux menaces

On sait qu'intégrer du logiciel pose **deux risques** :
- La **corruption de données** (l'agent casse tout), une base de données ou une base de code par exemple.
- La **divulgation de données** : l'agent récupère des données et il peut être programmé pour les transférer à des tiers ou simplement les donner à son pilote malicieux.

Et ces risques sont portés par **deux menaces** :
- La **menace humaine** : on instruit l'agent en lui demandant de faire des choses qu'il n'a pas le droit de faire, soit pour saboter des systèmes (corruption), soit en vue d'espionnage (divulgation). Et le pilote n'est pas forcément un inconnu sur Internet : ça peut être un **collaborateur malveillant ou manipulé**.
- La **menace de l'agent lui-même** : intrinsèquement l'agent est stupide et pourrait corrompre ou divulguer les données par « accident ». Je passerai pour l'instant la menace de l'**agent « double »**, qui aurait une intention secrète et apprise pour laquelle il ferait de l'extraction ou de la corruption volontaire (pour cacher une backdoor par exemple). Je note juste que, vu de l'extérieur, **un agent stupide et un agent double font la même chose** : ce qu'on contrôle pour l'un protège en partie de l'autre.

Pour adresser ceci, il faut mettre en place de l'ingénierie agentique, mais en termes d'architecture, **ces deux menaces n'ont pas la même réponse**.

## La menace humaine : le Persona Guardrail d'Uber

En ce qui concerne la menace humaine, une des réponses a été apportée par Uber dans son papier *Persona Guardrail: A Production-Grade Defense Framework for Agentic Systems* (ref [arxiv 2610.03434](https://arxiv.org/abs/2610.03434)). Le papier est long et j'avoue l'avoir fait ingérer par un LLM pour en comprendre l'essence. L'idée est que l'on met en place un premier LLM entre le pilote et l'agent pour vérifier que la demande entre bien dans le **périmètre déclaré** de l'agent (une liste d'intentions autorisées), et ainsi avoir un feu vert pour lancer l'exécution. Le même LLM vérifie la **réponse finale** avant qu'elle reparte vers le pilote. En revanche, **ce qui se passe entre les deux (appels d'outils, modifications) n'est pas contrôlé**, et le papier le dit lui-même.

## Protéger l'agent de lui-même : PPG et jevgo

D'un autre côté, je pense qu'il faut aussi mettre des **validations plus déterministes** au niveau de la **boucle agentique** pour protéger l'agent de lui-même. C'est ce que j'explore avec **PPG** ([poc-agentic-platform](https://github.com/owulveryck/poc-agentic-platform)) : une gateway qui valide avec des règles OPA/Rego le plan de l'agent, chacune de ses modifications et le diff final. **Pas de ticket, pas de modification.**
J'ai aussi envisagé la mise en place d'un **« système 1 »** (au sens de Kahneman : rapide, statistique, intuitif) à côté de la boucle déterministe, avec des **modèles de classification** de type [Jev](https://fr.wikipedia.org/wiki/Jev) (et mon implémentation jouet : [jevgo](https://github.com/owulveryck/jevgo)) : de petits modèles qui apprennent à imiter les règles. **Ils ne décident rien.** Quand leur verdict diverge de celui des règles, c'est le signe qu'une **règle est peut-être mal écrite**, et on fait une escalade en cas de doute.

## Mesures compensatoires et mesures amplificatrices

Pour lire ces mesures, je fais une distinction :
- une mesure **compensatoire** compense une faiblesse : elle bloque, et chaque blocage demande l'intervention d'un humain ;
- une mesure **amplificatrice** permet à la boucle agentique de se corriger seule : l'erreur revient à l'agent, qui corrige avant que l'humain ait à intervenir. L'humain passe de **« in-the-loop »** (il valide chaque étape) à **« on-the-loop »** (il supervise et traite les exceptions).

## Trois infographies

Voici un résumé dans trois infographies.

**1. Où agit chaque garde-fou.** Uber filtre ce qu'on demande et ce qu'on répond, PPG contrôle ce que l'agent fait entre les deux.

![Où agit chaque garde-fou dans le cycle d'une requête](/assets/menace-agentique/where-guardrails-act.fr.svg)

**2. Même principe, mécanismes opposés.** Les deux déclarent un périmètre et refusent par défaut. Mais on ne peut pas écrire de règle exacte pour du langage naturel, alors qu'on le peut pour un plan ou un diff.

![Même principe, mécanismes opposés : c'est la forme des entrées qui tranche](/assets/menace-agentique/same-principle-opposite-mechanisms.fr.svg)

**3. Le système 1 observe, il ne juge pas.** jevgo tourne à côté des règles. Un désaccord remonte à un humain, qui corrige la règle une fois pour toutes.

![Ajouter jevgo à PPG : un observateur des règles, pas un juge](/assets/menace-agentique/jevgo-observer.fr.svg)

## En conclusion

La mesure architecturale d'Uber est **essentiellement compensatoire** (sa boucle d'amélioration fait progresser la politique, pas l'agent). Et elle ne compense **pas la stupidité de l'agent mais son obéissance** : un agent plus intelligent obéira mieux à une demande malveillante, donc elle ne deviendra pas inutile quand les modèles progresseront. Elle risque au contraire d'avoir du mal à suivre. Le **classifieur est volontairement petit** pour tenir sous les 100 ms. Plus les agents comprendront les insinuations, plus une demande subtile pourra être comprise par l'agent sans être repérée par le classifieur, et plus la gateway laissera passer de **faux négatifs**.
Du côté de PPG, il y a évidemment un côté compensatoire, car j'ai affiché que le premier but était d'adresser la stupidité des modèles. Mais il y a aussi un **côté amplificateur** de la boucle agentique, qui lui permet de s'auto-corriger quand la gateway lui renvoie une erreur, avant que l'humain s'en rende compte. L'humain ne valide plus chaque étape, il reste **on-the-loop**.

En tout état de cause, une vraie architecture agentique du futur doit prendre en compte ces deux menaces, en **combinant des mesures qui protègent et des mesures qui rendent l'autonomie sûre**.
