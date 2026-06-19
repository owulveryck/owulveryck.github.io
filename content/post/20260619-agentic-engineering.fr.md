---
title: "Vibe Coding à l'échelle ? L'ingénierie contre-attaque"
date: 2026-06-19T10:00:00+02:00
lastmod: 2026-06-19T10:00:00+02:00
images: [/assets/agentic-engineering/octo_plateforme_agentique.fr.svg]
draft: false
keywords: []
summary: >
  Coder « dans les règles de l'art » ne suffit pas. Dans une grande organisation qui produit des dizaines d'applications, le vrai enjeu n'est pas de bien coder une app — c'est de garantir que chacune respecte l'état de l'art de l'organisation. Cela ne se décrète pas projet par projet. Cela s'industrialise dans une plateforme agentique.
tags: ["AI", "agents", "architecture", "plateforme", "vibe-coding", "ingénierie-agentique"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

L'IA générative a transformé la manière dont le code est produit. En quelques mois, nous sommes passés de l'autocomplétion à des agents capables d'écrire, tester et déployer des applications entières. Le marché est aujourd'hui inondé de méthodes pour cadrer ces agents et leur faire produire du code de qualité.

Mais cette abondance pose une question que peu d'organisations se posent encore : **que se passe-t-il quand on ne produit pas une app, mais cinquante ?**

<!--more-->

La réponse tient en une nuance que tout le monde sous-estime. Un développement peut être *qualitatif* sans être *au standard attendu*. Les méthodes d'IA structurée permettent d'atteindre **un** état de l'art — celui, générique, sur lequel l'industrie s'accorde. Mais l'état de l'art d'une organisation, lui, est *le sien* : contextuel, agréé, vivant. Et c'est précisément celui-là que l'IA ne connaît pas — et qu'elle réinvente, mal, à chaque projet.

## « Conforme à l'état de l'art » ≠ « conforme à VOTRE état de l'art »

Commençons par dissiper un malentendu de fond : **il n'y a pas de standard de qualité absolu**.

Comme le rappelle Christophe Thibault dans [« En finir avec la "dette technique" »](https://blog.octo.com/en-finir-avec-la-dette-technique), « La Qualité » au singulier, valable dans tous les contextes, n'existe pas. Au restaurant du Plaza Athénée, on fait un repas d'excellente qualité pour près de 380 €. Chez McDonald's, où l'on déjeune pour moins de 8 €, il existe aussi une direction de la qualité. Deux états de l'art radicalement différents — et chacun parfaitement légitime dans son contexte. Échangez leurs critères, et vous obtenez deux absurdités que personne n'accepterait de payer.

Un état de l'art, c'est donc l'ensemble des procédés qu'un groupe **agrée** — explicitement ou tacitement — pour produire la meilleure solution possible dans *son* contexte, avec *ses* contraintes. Il n'est ni universel, ni figé. Il se discute, s'expérimente, se renouvelle.

Or les méthodes d'IA structurée portent en elles un état de l'art : celui qui se dégage de tout le savoir-faire générique accumulé sur internet, dans les librairies et les frameworks. Cet état de l'art est précieux. Mais il n'est pas le vôtre. Et toute la difficulté est là : une application peut être impeccable au regard du standard générique, et pourtant hors-jeu au regard des règles de votre organisation.

C'est pourquoi nous distinguons trois approches du développement assisté par IA — non pas selon les outils, mais selon **l'état de l'art qu'elles garantissent**.

**Le Vibe Coding** va vite. On prompte, on itère, on accepte ce qui « a l'air de marcher ». C'est l'approche idéale pour un prototype, un hack, une exploration. Mais la qualité dépend entièrement du prompt et du développeur. Aucun état de l'art garanti, aucune reproductibilité.

**L'IA assistée structurée** (dont [BMAD](https://github.com/bmadcode/BMAD-METHOD) est un bon exemple) va plus loin. Elle impose des templates, des règles, des prompts détaillés. Le résultat est une application bien construite, codée « dans les règles de l'art ». C'est ce que nous appelons la **certitude contextuelle** : pour ce projet-là, avec ce contexte-là, le résultat est fiable.

**L'ingénierie agentique** introduit une **certitude systémique**. La fiabilité n'est plus portée par un individu ou un prompt — elle est portée par une **plateforme** qui rend l'état de l'art de l'organisation disponible à chaque agent, chaque projet, chaque fois.

![Certitude contextuelle vs. Certitude systémique](/assets/agentic-engineering/octo_spectre_certitude.fr.svg)
*Le facteur différenciant n'est pas l'utilisation de l'IA : c'est l'état de l'art que l'approche garantit. Inspiré de la Figure 3 de [^osmani2026].*

## Ce que l'IA structurée fait vraiment — et ce qu'elle ne fait pas

Soyons justes : l'IA structurée est un vrai progrès. Elle n'est pas l'ennemi. Elle apporte une rigueur que le vibe coding ne peut pas offrir, et produit des applications bien faites.

Mais elle souffre de deux angles morts structurels.

**Premier angle mort : elle impose sa propre chaîne de construction.** Ses templates, ses conventions, ses étapes sont conçus *en dehors* de l'organisation. Ils font fi des processus existants — ceux-là mêmes qui, dans une grande organisation, sont les garants de la qualité globale attendue par les clients. Ces processus méritent souvent d'être revus, certes. Mais ils ne peuvent pas être ignorés.

**Second angle mort : la certitude qu'elle produit est locale.** Elle ne dit rien sur le projet suivant. Le prochain développeur, le prochain agent, le prochain projet repart de zéro. Les règles de l'art sont respectées, mais pas forcément les règles de l'organisation. Le code est bien fait, mais pas forcément cohérent avec les autres applications du même système.

C'est ici qu'opère une mécanique que Jerry Weinberg a formulée il y a longtemps :

> *The First Law of Technology Transfer: Long-range good tends to be sacrificed to short-range good.*
> — Jerry Weinberg, *Quality Software Management*

L'IA structurée optimise le bien à court terme (*cette* application, livrée vite et proprement) au détriment du bien à long terme : la cohérence de l'ensemble. Pour un projet, c'est invisible. Pour cinquante, c'est dévastateur.

## Le vrai danger à l'échelle : on réinvente le commoditaire

Voici le mécanisme précis du chaos.

Quand chaque agent produit son application sans connaître l'état de l'art de l'organisation, il **invente le sien**. Et surtout, il **refait des sous-composants qui existent déjà ailleurs** — sans le savoir. Authentification, gestion des erreurs, observabilité, accès aux données, composants d'interface : autant de briques commoditaires reconstruites projet après projet.

Le résultat n'est pas N applications excellentes. C'est **N composants moyens**. Au lieu d'améliorer un composant partagé et d'élever continuellement son niveau de qualité, on dissémine des versions médiocres et divergentes du même besoin. **On ne productise rien.**

Christophe Thibault propose un terme bien plus juste que « dette technique » pour nommer ce phénomène : le **conflit de procédés**. Une solution n'est pas « endettée » ; elle s'appuie sur des procédés qui, faute d'avoir été agréés ensemble, se contredisent. Avant l'IA, ce conflit était *rare et ressenti* : un développeur butait sur une friction, en parlait au stand-up, le conflit se posait et se résolvait.

À l'échelle de l'IA générative, ce conflit change de nature et devient redoutable pour deux raisons :

- **Il devient invisible.** Chaque application « a l'air bonne » parce qu'elle respecte l'état de l'art *générique*. Personne ne ressent la friction. L'agent ne bute sur rien. Le conflit avec l'état de l'art *de l'organisation* n'est jamais constaté — jusqu'à ce qu'il se manifeste en production, en incident, en incohérence visible par le client.
- **Il se multiplie.** Le foisonnement et l'abondance d'applications que l'IA permet de produire transforment un conflit ponctuel en divergence systémique. Cinquante projets, c'est potentiellement cinquante états de l'art incompatibles.

Et la facture finit toujours par arriver. Car la confiance, elle, est transitive. L'utilisateur final fait confiance à une **marque**. Cette marque déploie des **applications**. Ces applications sont produites par des **agents**. Si un seul maillon cède, c'est la confiance en la marque entière qui s'érode. Une marque premium ne peut pas se permettre une UX amateur. Une banque ne peut pas tolérer une faille de sécurité dans un chatbot. Un opérateur télécom ne peut pas laisser passer un ton inadapté dans une app de service client.

![La Confiance Transitive](/assets/agentic-engineering/octo_confiance_transitive.fr.svg)
*Chaque application déployée est un test pour la réputation de la marque. La confiance se construit lentement et se perd en un instant.*

## La plateforme : l'état de l'art rendu exécutable

Si le problème, c'est que les agents ne consomment pas l'état de l'art de l'organisation mais le réinventent, la solution s'énonce simplement :

> **Une plateforme agentique, c'est l'état de l'art de l'organisation rendu exécutable et consommable par des IA — et gouverné comme un produit.**

La bonne nouvelle, c'est que la matière existe déjà. Les politiques de sécurité, les standards d'architecture, le design system, les connaissances métier, les conventions : tout cela est là, dans l'organisation. Le problème n'est pas qu'il manque : c'est qu'il n'est pas **structuré pour être exploitable par des agents**. Il vit dans des wikis, des têtes, des PDF, des slides. Inexploitable par un système agentique.

Le rôle de la plateforme est de transformer ce patrimoine en capacités consommables : l'organiser, l'outiller, et le doter de garde-fous de validation. Ces capacités s'alimentent auprès de **multiples producteurs de contenu** — experts métier, équipes sécurité, équipes marque, architectes — et les rendent disponibles à de **multiples consommateurs** — les agents, les projets, les applications.

C'est exactement la définition d'une plateforme au sens où l'entend Sangeet Paul Choudary dans [« Platform Thinking: The Future of Work »](https://medium.com/@sanguit/platform-thinking-the-future-of-work-b49aeb0c1e53) : un dispositif qui met en relation un écosystème de producteurs et de consommateurs, et qui **commoditise les opérations répétables**. Choudary fait une remarque essentielle pour notre sujet : *écrire du code n'est pas une opération répétable — c'est une activité d'infrastructure ponctuelle, comme construire une chaîne de montage.* Ce qui est répétable, ce sont les opérations que ce code automatise. La plateforme industrialise précisément ce répétable, pour libérer l'écosystème sur ce qui ne l'est pas.

Concrètement, la plateforme fournit trois familles de capacités systémiques :

- Le **contexte systémique** (instructions, connaissances métier, mémoire, exemples) est structuré une fois et injecté dans chaque agent.
- Les **garde-fous systémiques** (sécurité, fiabilité, cohérence de marque, conventions) sont appliqués automatiquement, pas au bon vouloir du développeur.
- L'**outillage** (serveurs MCP, pipelines CI/CD, évaluations) est mutualisé et prêt à l'emploi.

Le spécifique — les intentions métier, les garde-fous particuliers d'un produit — vient s'ajouter au systémique. Mais le socle est là, pour tous les projets.

![D'où vient le contexte ?](/assets/agentic-engineering/octo_usine_vs_plateforme.fr.svg)
*L'usine est nécessaire. Mais sans plateforme, chaque projet réinvente le contexte, les garde-fous, l'outillage. Inspiré de la Figure 6 de [^osmani2026].*

Le rôle du développeur s'en trouve transformé. Il ne produit plus du code : il **package du contexte**. Mais ce contexte ne vient pas de lui seul : il provient des experts métier, des équipes sécurité, des équipes marque, des architectes. Chacun produit un fragment ; le développeur les assemble. Et si chacun fait ce packaging à sa manière, projet par projet, on retombe dans le problème de départ. C'est précisément ce que la plateforme industrialise : un cadre où chaque source de contexte est structurée, versionnée, et injectée automatiquement.

## Surprendre par la valeur, pas par le standard

Une fois le commoditaire pris en charge par la plateforme, un principe directeur s'impose.

Il faut faire un **shift-left** sur tout le commoditaire : déplacer la qualité, la sécurité et les standards le plus tôt possible, en amont du projet, dans la plateforme, pour qu'ils soient un **acquis** et non un livrable à refaire. Le standard devient un point de départ garanti, pas un point d'arrivée à reconquérir à chaque fois.

Le bénéfice est double. D'abord, le commoditaire, parce qu'il est partagé et productisé, **s'améliore en continu** : chaque usage durcit le composant, au lieu de disperser N versions moyennes. Ensuite, et c'est l'essentiel, l'énergie de l'IA se concentre là où elle crée de la valeur — sur le métier, sur le différenciant.

Car c'est bien là le but :

> **Une application ne doit pas surprendre le client par ses composants standards. Elle doit le surprendre par sa valeur métier.**

Le standard, on l'attend. Il doit être irréprochable, et invisible. La surprise, la satisfaction, l'avantage concurrentiel : ils viennent de ce que personne d'autre ne fait. C'est exactement ce que permet une plateforme — commoditiser le banal pour rendre l'exceptionnel possible à l'échelle.

![La Plateforme Agentique](/assets/agentic-engineering/octo_plateforme_agentique.fr.svg)
*Des capacités AI-ready mutualisées pour produire des applications fiables à l'échelle. Le spécifique se mêle au systémique pour produire l'application.*

## Une plateforme n'est pas un dogme : c'est un produit gouverné

À ce stade, une objection sérieuse se présente — et c'est précisément en y répondant que la thèse tient debout.

Ce qui est commoditaire n'est pas une donnée stable. Ce qui différencie aujourd'hui devient le standard de demain. La frontière entre le commun et le différenciant **se déplace en permanence**. Une plateforme qui figerait l'état de l'art risquerait donc de figer un état de l'art… périmé, exactement le travers que Thibault dénonce sous le nom de « La Qualité » en absolu.

La réponse n'est pas de renoncer à la plateforme. C'est de la **traiter comme un produit, gouverné**. La direction technique — les CTO, l'architecture, les référents — fait évoluer en continu ce que la plateforme considère comme commoditaire, et ce qu'elle laisse à la main des applications. Choudary le rappelle : ouvrir une plateforme à un écosystème, c'est accepter de lâcher du contrôle, et donc bâtir les *checks and balances* qui maintiennent le niveau.

C'est ainsi que l'état de l'art reste **vivant et agréé**, au sens de Thibault, plutôt que gravé dans le marbre. La plateforme n'est pas l'autorité qui décrète la qualité absolue. Elle est l'instrument qui rend exécutable un état de l'art que des humains continuent de discuter, d'expérimenter et de renouveler.

## Une conviction

Le vibe coding a sa place. Pour explorer, prototyper, apprendre. Il n'est pas l'ennemi — il est le point de départ.

L'IA assistée structurée — BMAD et les approches similaires — est un vrai progrès. Elle apporte de la rigueur à la production individuelle. Mais soyons clairs sur ce qu'elle ne fait pas : **dans une grande organisation qui exploite des dizaines d'applications, choisir une bonne méthode d'IA structurée ne suffit pas.** Chaque projet reste un îlot, et rien ne garantit la cohérence — ni la qualité au standard de l'organisation — de l'ensemble.

L'enjeu est organisationnel. Il s'agit de garantir que *chaque* application produite par les agents reflète l'état de l'art de l'organisation — sécurité, fiabilité, cohérence de marque, valeur — et que le passage à l'échelle ne soit pas synonyme de perte de contrôle. Cela suppose trois déplacements :

- **Cesser de croire** qu'une méthode de cadrage des agents suffit à l'échelle. Le code « bien fait » n'est pas le code « au standard ».
- **Traiter sa plateforme agentique comme un produit gouverné**, et non comme un outil. C'est la condition pour que l'état de l'art reste vivant.
- **Inventorier et structurer son propre état de l'art** — sécurité, architecture, design system, métier — pour le rendre exécutable et consommable par les agents.

Le développeur ne produit plus du code. Il package le contexte qui produit le code. Et pour industrialiser ce packaging, pour rendre l'état de l'art de l'organisation exécutable, à l'échelle, sans le figer — il a besoin d'une plateforme.

![L'ingénierie du contexte](/assets/agentic-engineering/fig4_context_engineering.fr.svg)
*Les meilleurs systèmes traitent le contexte comme une décision architecturale de premier ordre, revue et versionnée comme du code. Adapté de la Figure 4 de [^osmani2026].*

---

## Sources

[^osmani2026]: Addy Osmani, Shubham Saboo, Sokratis Kartakis — *"The New SDLC With Vibe Coding: From ad-hoc prompting to Agentic Engineering"*, Google, Mai 2026.

Christophe Thibault — [*« En finir avec la "dette technique" »*](https://blog.octo.com/en-finir-avec-la-dette-technique), OCTO.

Sangeet Paul Choudary — [*« Platform Thinking: The Future of Work »*](https://medium.com/@sanguit/platform-thinking-the-future-of-work-b49aeb0c1e53), Medium, 2013.

Jerry Weinberg — *Quality Software Management*.
