---
title: "Qui fait quoi ? Team Topologies de la plateforme agentique"
date: 2026-06-22T10:00:00+02:00
lastmod: 2026-06-22T10:00:00+02:00
images: [/assets/team-topologies-agentique/octo_team_topologies_agentique.fr.svg]
draft: false
keywords: ["team topologies", "plateforme agentique", "charge cognitive", "organisation", "ingénierie agentique"]
summary: >
  La plateforme agentique dit ce qu'il faut fournir. Team Topologies dit qui le fournit, et comment les équipes interagissent pour y parvenir. Ce deuxième article applique le modèle de Skelton et Pais à un monde où les agents produisent et les humains orchestrent.
tags: ["AI", "agents", "architecture", "plateforme", "team-topologies", "ingénierie-agentique"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

> *La plateforme agentique dit ce qu'il faut fournir. Team Topologies dit qui le fournit, et comment les équipes interagissent pour y parvenir.*

---

Dans [le premier article de cette série](/fr/2026/06/19/vibe-coding-%C3%A0-l%C3%A9chelle-ling%C3%A9nierie-contre-attaque.html), nous posions la question du **quoi** : quelles capacités systémiques (contexte, garde-fous, outillage) sont nécessaires pour produire des applications fiables à l'échelle. La réponse était la plateforme agentique, et en son cœur, l'*usine agentique* : le mécanisme où les agents planifient, codent, testent et livrent.

Mais une plateforme ne se construit pas toute seule, et surtout, elle ne se consomme pas de la même manière qu'elle se construit. Il reste une question fondamentale : **qui fait quoi ?**

<!--more-->

## Le vrai problème : la charge cognitive de la production agentique

*Avant de demander qui fait quoi, il faut comprendre pourquoi la question se pose autrement.*

Produire une application, c'était orchestrer des rôles dans le temps : l'un concevait, un autre challengeait l'architecture, un troisième testait, un quatrième déployait. La complexité était réelle, mais **distribuée** : entre plusieurs personnes, et étalée sur plusieurs moments. Chaque rôle posait ses questions à son tour.

L'agent change l'équation. Il ne pose pas de questions : il produit une réponse, immédiatement. Il ne fatigue jamais, ne se repose jamais, n'attend pas. Sa vitesse est sa force, et son piège. Toutes les questions que les rôles posaient *en séquence*, l'humain qui le pilote doit désormais les anticiper *en amont, en parallèle, dans le temps court d'un prompt*. Mal cadré, l'agent ne ralentit pas : il produit vite, et à côté.

La charge cognitive ne disparaît donc pas avec l'IA : elle se transforme. Elle devient d'abord une **charge d'anticipation** : tout ce que l'humain doit prévoir avant de lancer l'agent, faute de quoi le résultat ne sera pas à la hauteur. Et parce que l'agent produit en continu, sans rythme humain, elle devient aussi un **débit cognitif** : un flux de décisions à soutenir dans le temps. Le vrai problème de la production agentique à l'échelle n'est pas que la complexité augmente, c'est qu'elle se **compresse**, sur un seul acteur et dans un temps qui ne peuvent l'absorber seuls.

C'est précisément ce que la plateforme vient régler. Elle absorbe une partie de la charge d'anticipation en se rendant **interrogeable par l'agent lui-même** : « ne te soucie pas de la sécurité » signifie que l'agent ira chercher auprès de la plateforme comment faire, et que des contrôles déterministes garantiront le résultat en aval. Elle ne supprime pas la réflexion : elle **rétrécit le champ des questions que l'humain doit porter**, pour le laisser se concentrer sur celles qui comptent vraiment, les décisions litigieuses et structurelles, là où le jugement humain reste irremplaçable.

La charge cognitive n'est donc plus seulement, comme chez Skelton et Pais, une *quantité à répartir entre équipes*. Dans le monde agentique, c'est aussi un *débit à réguler dans le temps*. Team Topologies nous dit comment répartir ; il nous reste à dire comment absorber. C'est l'objet de cet article.

## Team Topologies, une réponse à la charge

Pour répondre à ce problème de charge, nous nous appuyons sur *Team Topologies*[^skelton2019], un modèle organisationnel qui définit quatre types d'équipes et trois modes d'interaction. Ce n'est pas un hasard : **son argument fondateur est précisément la charge cognitive** : une équipe ne peut être efficace que si elle ne porte que la complexité qu'elle est capable d'absorber. Skelton et Pais raisonnent toutefois sur une charge *structurelle*, à répartir entre équipes ; nous l'étendons à la charge *dynamique* de l'acte de production agentique, décrite plus haut. Le modèle nous donne la grille pour distribuer ce qui peut l'être, et désigner ce que la plateforme doit absorber.

Précisons d'emblée notre posture : ce qui suit est une **conviction prospective**, la cible organisationnelle vers laquelle tend, selon nous, la production d'applications par des agents. La question n'est plus *de quoi l'usine a besoin*, mais *qui l'opère*.

C'est exactement ce que fait la plateforme agentique : elle **absorbe la complexité technique** pour que les équipes métier ne portent que la charge cognitive de leur domaine. Le développeur se déplace : il construit la plateforme qui permet à d'autres de produire. Symétriquement, la production d'applications s'ouvre aux équipes métier via les agents.

### Ce que nous gardons et ce que nous adaptons

Appliquer Team Topologies à l'agentique suppose d'être clair sur nos déviations. Pour éviter de vider le modèle de sa substance tout en revendiquant son nom, voici la ligne.

**Ce que nous conservons intact :** le principe directeur de la **charge cognitive** ; les **quatre types d'équipes** ; les **trois modes d'interaction** ; la **bascule de la collaboration vers le *X-as-a-Service*** à maturité.

**Ce que nous adaptons, et pourquoi :**
- les équipes *stream-aligned* peuvent être **non-techniques** (métier), parce que la plateforme absorbe la charge technique ;
- elles ne portent pas la **responsabilité opérationnelle de bout en bout** (run, incidents), absorbée par la plateforme ;
- l'*enabling* n'est pas seulement transitoire : il est **structurellement compensé** par la plateforme, le producteur n'étant plus un développeur ;
- la **charge cognitive** n'est plus seulement une quantité à répartir entre équipes, mais un **débit à réguler dans le temps** : la plateforme absorbe la charge d'anticipation en amont de la production agentique.

Ces adaptations ne trahissent pas l'esprit de Skelton et Pais : elles l'appliquent à un contexte qu'ils n'avaient pas anticipé, celui où l'agent produit, et où l'humain orchestre.

## Quatre types d'équipes, un seul objectif

L'objectif est commun : produire des applications fiables, cohérentes avec les standards de l'organisation, à l'échelle. Mais les rôles sont distincts.

![Team Topologies de la Plateforme Agentique](/assets/team-topologies-agentique/octo_team_topologies_agentique.fr.svg)
*Les quatre types d'équipes de Team Topologies appliqués à la plateforme agentique. Chaque équipe a un rôle précis dans la chaîne de production.*

### Les équipes stream-aligned : produire les applications

Les équipes stream-aligned sont les équipes produit. Ce sont elles qui pilotent l'orchestrateur IA (le moteur de l'usine agentique décrite dans le premier article), définissent les intentions métier, et fournissent le **contexte dynamique** : spécifications, garde-fous spécifiques au produit, connaissances propres à leur domaine.

La transformation est profonde : ces équipes ne sont plus nécessairement composées de développeurs. Ce sont de plus en plus des **équipes métier** (experts du domaine, product managers, analystes) qui pilotent directement la production via les agents. Cette évolution rapproche la production du besoin, mais pose un risque : ces équipes n'ont pas toujours conscience des impacts de mettre une application en production. C'est précisément pourquoi les autres types d'équipes existent.

Un point de rigueur s'impose. Dans Team Topologies classique, une équipe stream-aligned est responsable **de bout en bout** d'un flux de valeur, y compris le run et les incidents. Ici, la **plateforme absorbe la responsabilité opérationnelle** (déploiement, monitoring, rollback). L'équipe stream-aligned reste responsable du *quoi* (l'intention, la qualité métier) ; la plateforme garantit le *comment* (la mise en production fiable). Cette répartition exige une [plateforme suffisamment mature](#une-plateforme-est-mature-quand). L'**astreinte** se répartit en conséquence : l'équipe plateforme porte les **incidents systémiques** (infrastructure, garde-fous, pipelines) ; les **décisions métier** (retrait de contenu, rollback produit) restent avec l'équipe stream-aligned.

Cette frontière quoi/comment est **plus poreuse qu'il n'y paraît**. La cohérence de marque, par exemple, relève du métier (c'est le *quoi*), mais sa vérification est automatisée par la plateforme (c'est le *comment*). La plateforme garantit des *minimums* ; elle ne garantit pas l'excellence métier.

### L'équipe plateforme : industrialiser les capacités

L'équipe plateforme fournit les trois piliers systémiques en X-as-a-Service :

- Le **contexte systémique** : instructions, rôles, connaissances métier mutualisées, mémoire, exemples et patterns
- Les **garde-fous systémiques** : sécurité, fiabilité, cohérence de marque, conventions
- L'**outillage et les skills** : serveurs MCP, pipelines CI/CD, évaluations, skills partagés

Le modèle est celui du self-service : documenté, versionné, consommable sans friction. L'effort de conception est amorti une fois, puis appliqué à chaque projet.

#### Une plateforme est « mature » quand…  

Pour éviter que le mot reste un vœu pieux, voici les critères observables :

- **Couverture des garde-fous** : les dimensions critiques (sécurité, fiabilité, cohérence de marque) sont couvertes automatiquement, pas par convention orale
- **Fiabilité des pipelines** : le taux de succès des déploiements est mesurable et suivi (SLA interne)
- **Part de self-service** : la majorité des déploiements se font sans intervention de l'équipe plateforme
- **Complétude de la documentation** : chaque capacité exposée est documentée et accompagnée d'exemples
- **Traçabilité des décisions** : les garde-fous produisent un journal d'audit (pourquoi un déploiement a été bloqué, quelle règle a été appliquée, quel seuil a été franchi)

Tant que ces critères ne sont pas atteints, la plateforme ne peut pas absorber la responsabilité opérationnelle des équipes stream-aligned, et le modèle repose sur l'enabling pour compenser.

### Les équipes enabling : combler le fossé

L'équipe enabling est **temporaire par nature**. Son objectif n'est pas de devenir indispensable, mais de rendre les équipes produit autonomes. Concrètement :
- **Mise à disposition de l'environnement** : outils, accès, configuration
- **Formation** : au context packaging, aux garde-fous, au pilotage de l'orchestrateur
- **Shift-left** sur les pratiques (sécurité, tests, qualité) en attendant leur encapsulation par la plateforme

Elle comble le fossé entre l'intention métier et l'exigence de qualité. Son rôle diminue à mesure que les équipes produit montent en compétence et que la plateforme s'enrichit.

Une nuance importante : dans le monde agentique, les équipes stream-aligned restent souvent **majoritairement non-techniques**. Le fossé ne se comble donc jamais entièrement par la montée en compétence : il est **structurellement compensé par la plateforme**. Ce n'est pas un échec du modèle, c'est une adaptation à un contexte où le producteur d'applications n'est plus un développeur.

### Les équipes complicated subsystem : maîtriser la complexité technique

Les équipes *complicated subsystem* travaillent sur les aspects les plus techniques de l'infrastructure IA. Leur expertise est profonde et spécialisée : elle ne doit pas être dispersée dans les équipes produit.

**Ce type d'équipe n'est pas universel.** Une organisation qui consomme exclusivement des API de modèles n'en a pas nécessairement besoin. En revanche, dès qu'elle gère ses propres modèles, optimise ses coûts d'inférence ou a des contraintes de souveraineté, ce type d'équipe devient indispensable, de même pour l'évaluation, le red-teaming, l'ingénierie RAG avancée ou le fine-tuning.

Elles collaborent avec l'équipe plateforme sur l'efficience des modèles, le KV cache, l'inférence souveraine, l'optimisation des coûts et l'évaluation. Leur travail ne traverse jamais directement les équipes produit : il passe par la plateforme.

## Les trois modes d'interaction

Team Topologies définit trois modes d'interaction entre équipes :

**Facilitating** : l'équipe enabling *facilite* les équipes stream-aligned. Interaction temporaire, orientée vers l'autonomie : l'enabler apprend à l'équipe produit à faire elle-même.

**X-as-a-Service** : la plateforme fournit ses capacités en self-service. C'est l'interaction cible, celle qui permet le passage à l'échelle.

**Collaboration** : l'équipe complicated subsystem *collabore* avec l'équipe plateforme. Interaction profonde, justifiée en phase de construction. À maturité, elle évolue vers un X-as-a-Service : les capacités d'efficience IA deviennent des services consommables, pas des chantiers permanents.

---

## Faire vivre le modèle dans la durée

Définir les équipes et leurs interactions ne suffit pas. Un modèle organisationnel n'a de valeur que s'il survit à son lancement.

### Le voyage vers l'autonomie

Le rôle de l'équipe enabling est conçu pour diminuer : c'est le signe que le système fonctionne.

![L'évolution des modes d'interaction](/assets/team-topologies-agentique/octo_team_topologies_evolution.fr.svg)
*Les trois modes d'interaction évoluent ensemble : le facilitating s'efface, la collaboration cède la place au X-as-a-Service, qui devient le mode dominant.*

Cette évolution suit deux axes parallèles qu'il faut distinguer : la **maturité des équipes** et la **maturité de la plateforme**. Les deux se renforcent mutuellement, mais ne progressent pas au même rythme.

**Au démarrage :**
- *Côté équipes* : les équipes produit découvrent le context packaging et le pilotage de l'orchestrateur. L'enabling est omniprésent.
- *Côté plateforme* : les capacités sont de base (quelques garde-fous, une documentation partielle, un pipeline encore fragile). La plateforme n'est pas encore mature au sens des critères définis plus haut.

**En phase de maturation :**
- *Côté équipes* : les équipes produit maîtrisent les fondamentaux. L'enabling devient ciblé (un conseil sur un garde-fou spécifique, un retour sur un prompt complexe).
- *Côté plateforme* : les garde-fous s'étoffent, le self-service progresse, la documentation se structure. La collaboration avec les complicated subsystems commence à se transformer en services intégrés.

**À l'autonomie :**
- *Côté équipes* : les équipes produit sont autonomes. L'enabling est optionnel, limité à de l'expertise ponctuelle.
- *Côté plateforme* : le self-service est complet, les garde-fous couvrent les dimensions critiques, la traçabilité est en place. L'interaction dominante est le X-as-a-Service.

**L'enabling disparaît parce qu'il réussit, pas parce qu'il échoue.**

**Un exemple concret.** L'équipe marketing veut produire une landing page. Elle fournit le contexte dynamique (intention de campagne, messages clés). La plateforme injecte le contexte systémique (charte graphique, composants UI, règles d'accessibilité) et les garde-fous vérifient cohérence de marque et sécurité. L'enabling avait formé le marketing trois mois plus tôt : aujourd'hui, intervention ponctuelle. L'équipe complicated subsystem a optimisé le KV cache : le contexte de marque, déjà tokenisé, est servi depuis le cache au lieu d'être recalculé à chaque itération, réduisant le coût par génération. Résultat : une landing page conforme, sécurisée, déployée, sans que le marketing ait eu besoin de savoir ce qu'est un pipeline CI/CD.

Le revers : une équipe qui ne « voit » jamais les mécanismes de protection perd la capacité d'en juger la pertinence. Les garde-fous doivent être **transparents dans leurs décisions**, même s'ils sont opaques dans leur implémentation.

### La gouvernance des applications : éviter le shadow IT à l'échelle

Si des équipes métier non-techniques peuvent produire des applications en production, **qui décide qu'une application a le droit d'exister ?** Et qui gère son cycle de vie (dette technique, dépréciation, coût cumulé) ?

Sans gouvernance, on aboutit à un shadow IT industrialisé. La plateforme est le levier de cette gouvernance : en centralisant déploiement, monitoring et métriques d'usage, elle offre une **visibilité systémique** sur le parc applicatif. Le product owner de la plateforme peut suivre les applications actives, identifier celles qui ne sont plus maintenues, et déclencher leur dépréciation. Une application qui ne passe plus les contrôles de sécurité est signalée automatiquement, pas oubliée silencieusement.

Le principe est simple : **la facilité de production doit s'accompagner d'une facilité de suivi**. Si la plateforme rend trivial de créer une application, elle doit rendre tout aussi trivial de savoir combien existent, qui les utilise, et ce qu'elles coûtent.

### Le chemin de graduation : du spécifique au systémique

Un mécanisme souvent absent des discussions sur les plateformes : **quand un garde-fou « produit » devient-il un garde-fou « plateforme » ?**

Exemple : l'équipe marketing implémente un garde-fou de contraste WCAG. L'équipe RH rencontre le même besoin trois mois plus tard. Puis l'équipe e-commerce. C'est le signal de graduation : un garde-fou répété par au moins trois équipes devient candidat à la systémisation, la *rule of three*[^fowler] appliquée aux garde-fous. Le product owner de la plateforme généralise le candidat, le rend configurable, le documente : toutes les équipes en bénéficient.

Les équipes stream-aligned remontent les besoins récurrents. L'équipe plateforme évalue et intègre. Les enablers détectent les besoins transverses dans leur accompagnement.

L'autonomie n'est pas le silence : même à maturité, les équipes produit restent les capteurs qui nourrissent la graduation. Le self-service ne supprime pas la boucle de feedback, il la rend moins coûteuse. La plateforme est un *produit vivant*[^skelton2019], avec son backlog, son product owner, et ses cycles d'itération propres.

**Un risque doit être nommé.** Ce product owner cumule potentiellement trois charges : le backlog technique de la plateforme, la graduation des garde-fous, et la gouvernance du parc applicatif. C'est, paradoxalement, le nouveau goulot d'étranglement que le modèle prétend éviter.

**Pourquoi ce risque survient.** Plus la plateforme réussit, plus elle concentre de décisions, et un seul rôle ne peut arbitrer manuellement le flux de besoins de dizaines d'équipes produit.

**La réponse est dans l'outillage.** La plateforme doit automatiser la détection des candidats à la graduation, la priorisation par usage, et le suivi du parc. Le PO arbitre ; il ne porte pas tout.

Sans ce mécanisme de graduation, la plateforme stagne. Les équipes produit réinventent les mêmes solutions. Les garde-fous restent locaux. La fiabilité systémique s'érode.

---

## Synthèse opérationnelle

Les principes posés, il reste à les cristalliser en règles actionnables.

### Qui possède quoi ?

Chaque capacité a **un propriétaire** (redevable du résultat) et souvent **des contributeurs** (qui co-construisent sans être redevables). Sans cette clarté, les capacités tombent dans les interstices.

![Responsabilités par équipe](/assets/team-topologies-agentique/octo_team_topologies_responsabilites.fr.svg)
*Chaque capacité a un propriétaire et des contributeurs. Les cellules vides signifient « non impliqué ». Le spécifique reste avec les équipes produit, le systémique avec la plateforme.*

Le principe directeur est simple : le spécifique (contexte dynamique, garde-fous produit, connaissances métier) reste avec les équipes stream-aligned, car la connaissance métier *appartient* au métier. Le systémique (contexte systémique, garde-fous organisationnels, outillage, orchestrateur) est industrialisé par la plateforme. L'orchestrateur IA est le cas emblématique : l'outil est fourni par la plateforme, la configuration reste chez l'équipe produit.

### Une conviction

La plateforme agentique répond au *quoi*. Team Topologies répond au *qui* et au *comment*, non comme une cloison étanche, mais comme une **répartition de redevabilité** : à l'équipe produit l'intention et le contexte métier, à la plateforme la fiabilité de la mise en production.

Sans ces frontières, la plateforme devient un goulot d'étranglement ou un terrain politique.

Ce modèle n'est pas universel : il s'adresse aux organisations qui produisent **plusieurs applications agentiques en parallèle**. En pratique (c'est une heuristique empirique, pas un seuil absolu), à partir de trois à cinq équipes produit, le coût cumulé de la réinvention (contexte recréé, garde-fous réimplémentés, incohérences à corriger) dépasse l'investissement d'une plateforme partagée.

Appliqué à l'ingénierie agentique, Team Topologies donne une structure qui :

- Permet aux **équipes métier de devenir stream-aligned**, sans exiger qu'elles deviennent des développeurs
- Garantit que les **capacités systémiques sont industrialisées**, pas réinventées projet par projet
- Accepte que le **fossé entre métier et production est réel**, et le comble avec des enablers, pas avec de l'espoir
- Isole la **complexité technique** dans des équipes spécialisées qui nourrissent la plateforme

Ce qui change, c'est que produire une application **n'exige plus de coder**, mais exige d'autres compétences : formuler une intention, structurer un contexte, piloter un orchestrateur. Ce qui ne change pas, c'est qu'il faut des **garde-fous, de la cohérence et de la fiabilité** pour la mettre en production.

Ce déplacement répond au problème posé en ouverture. Aujourd'hui, c'est le développeur qui porte la charge d'anticipation : il sait quelles questions l'agent ne posera pas. Demain, à mesure que la plateforme en absorbe une part croissante, le seuil de compétence requis pour piloter baisse : un métier peut produire sans porter seul ce que la plateforme anticipe à sa place. **La trajectoire du producteur, du développeur vers le métier, n'est pas un présupposé : c'est la conséquence mesurable d'une charge d'anticipation progressivement absorbée.**

Le développeur, lui, ne disparaît pas : il se déplace. Il ne produit plus l'application, il construit la plateforme qui permet à d'autres de la produire.

---

## Sources

[^skelton2019]: Matthew Skelton, Manuel Pais — *"Team Topologies: Organizing Business and Technology Teams for Fast Flow"*, IT Revolution, 2019.

[^osmani2026]: Addy Osmani, Shubham Saboo, Sokratis Kartakis — *"The New SDLC With Vibe Coding: From ad-hoc prompting to Agentic Engineering"*, Google, Mai 2026.

[^fowler]: Don Roberts, cité par Martin Fowler — *"Refactoring: Improving the Design of Existing Code"*, Addison-Wesley, 1999. « Three strikes and you refactor. » Le principe est ici transposé du refactoring de code à la graduation des garde-fous.
