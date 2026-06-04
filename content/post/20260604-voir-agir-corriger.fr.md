---
title: "Voir, Agir, Corriger : les trois leviers du développement avec un agent de code"
date: 2026-06-04T10:00:00+02:00
lastmod: 2026-06-04T10:00:00+02:00
images: [/assets/voir-agir-corriger/agent-leviers.svg]
draft: false
keywords: ["code agents", "Claude Code", "Copilot", "developer tools", "engineering practices"]
summary: "Un agent de code out-of-the-box ne voit qu'un repo et un shell. Cet article pose trois principes invariants pour passer d'un gadget à un outil de production."
tags: ["AI", "agents", "architecture", "engineering-practices"]
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
contentCopyright: false
reward: false
mathjax: false
---

*Un agent de code out-of-the-box ne voit qu'un repo et un shell. Pour l'ingénierie professionnelle, cela ne suffit pas. Voici les principes qui font la différence entre un gadget et un outil de production.*

> **Note liminaire.** Cet article est né du talk [« Beyond the Basics with Claude Code »](https://www.youtube.com/watch?v=tuY2ChJIx48) de Daisy Holman, ingénieure sur l'équipe Claude Code (mai 2026). Les idées fondatrices en sont issues, puis complétées par ma propre expérience de terrain et étayées par une recherche documentée pour les confronter aux données disponibles.
---

**« Si l'agent ne peut pas faire tout ce que vous faites, il ne peut pas travailler avec vous. »** — Daisy Holman, Anthropic

Cette phrase résume le problème que toute équipe finit par rencontrer. Un agent de code fraîchement installé sait lire du code et exécuter des commandes. C'est suffisant pour un prototype, un side project, un exercice zero-to-one. Mais le vrai travail d'ingénierie ne vit pas dans le code source. Il vit dans les threads Slack, les design docs, les dashboards de production, les discussions de revue, les décisions d'architecture non écrites. Le code dit *ce qui* est fait ; il dit rarement *pourquoi*.

Faites le test : passez une journée entière dans le terminal de votre agent, sans jamais basculer vers un autre outil. Chaque alt-tab est une connexion manquante, un endroit où l'agent ne peut pas vous suivre. Customiser un agent n'est pas un confort. C'est une nécessité structurelle.

Cet article s'adresse aux ingénieurs seniors, tech leads et architectes qui déploient ou envisagent de déployer des agents de code en environnement de production. Il pose trois principes invariants (valables quel que soit l'outil) et une méthode pour les appliquer.

*Note de transparence : ce cadre a été élaboré à partir d'une expérience approfondie avec Claude Code, puis généralisé. Les exemples de mécanismes (skills, hooks, MCP, prompt files) sont empruntés à cet écosystème, mais les principes s'appliquent aux primitives équivalentes des autres outils (rules files pour Cursor, custom instructions pour Copilot, repo maps pour Aider). Quand une recommandation est spécifique à un outil, elle est signalée.*

---

## 1. Le modèle mental : Voir, Agir, Corriger

Tout agent de code, quel que soit le moteur qui l'anime, se ramène à trois questions :

- **VOIR** : que sait l'agent ? Quelles sources d'information lui sont accessibles ? Code, logs, documentation interne, historique de conversation, état de la CI.
- **AGIR** : que peut faire l'agent ? Éditer des fichiers, lancer des tests, ouvrir des PRs, interroger une API, déployer.
- **CORRIGER** : qu'est-ce qui le corrige ? Linters, échecs de tests, retours de revue, hooks de validation, boucles de feedback automatiques.

![Les trois leviers d'un agent de code](/assets/voir-agir-corriger/agent-leviers.svg)

Ce cadre converge avec plusieurs traditions indépendantes : le cycle **OODA** (Observe, Orient, Decide, Act) en théorie de la décision, le paradigme **Observation / Action / Récompense** en apprentissage par renforcement, les boucles de feedback en ingénierie des systèmes. Ce n'est pas un hasard : ce sont des contraintes structurelles de tout système qui perçoit un environnement, agit dessus et s'ajuste.

Cette convergence n'est pas parfaite. Le cycle OODA de Boyd comporte une phase d'**Orientation** (la synthèse interprétative des données brutes) que notre cadre subsume dans VOIR. Ce choix est délibéré : en ingénierie de production, l'Orientation est en grande partie le fruit du packaging du contexte (quelles informations, dans quel ordre, sous quelle forme). Mais il faut en connaître le prix : un cadre sans Orientation explicite tend à traiter l'agent comme un système purement réactif. Pour les tâches nécessitant un raisonnement architectural profond, cette dimension mériterait d'être isolée.

De même, ce cadre diffère délibérément des architectures cognitives proposées par la recherche académique. L'architecture de référence de Lilian Weng (« LLM Powered Autonomous Agents », 2023) divise l'agent en Planification, Mémoire et Utilisation d'Outils [^weng]. Le paradigme ReAct de Shunyu Yao dissocie la trace de raisonnement de l'action externe. Notre cadre sacrifie cette granularité au profit de l'applicabilité : en production, la planification et la mémoire sont des problèmes de packaging du contexte (combien d'historique injecter, sous quelle forme, à quel moment). Mais si votre agent requiert de l'auto-réflexion (la capacité à réviser ses propres décisions indépendamment des signaux de correction externes), cette dimension mériterait un traitement séparé.

[^weng]: Lilian Weng, « LLM Powered Autonomous Agents », Lil'Log, 2023. https://lilianweng.github.io/posts/2023-06-23-agent/

L'assimilation de CORRIGER à la récompense en RL n'est pas métaphorique. Les travaux récents sur l'entraînement par renforcement des modèles de code utilisent explicitement des fonctions de récompense composites combinant correction fonctionnelle (tests unitaires), correction syntaxique (linters) et structure sémantique (graphes de flux de données) [^syncode]. Vos linters et vos tests ne sont pas de simples filtres a posteriori ; ils fonctionnent comme des fonctions de récompense multidimensionnelles qui façonnent la politique de génération de l'agent.

Tout le reste (fichiers d'instructions, serveurs MCP, hooks, skills, prompt files) n'est que l'implémentation de ces trois leviers. Les leviers sont l'invariant ; les mécanismes sont contingents.

[^syncode]: « Domain-Adaptable Reinforcement Learning for Code Generation with Dense Rewards », arXiv:2605.21180, 2025. https://arxiv.org/abs/2605.21180

### Le principe du maillon faible

**Un agent n'est jamais meilleur que le plus faible de ses trois leviers.**

Imaginez un diagramme en barres. Chaque levier a une hauteur. La performance globale de l'agent est plafonnée par la barre la plus courte.

Ce principe est massivement documenté. Une étude interne d'Anthropic portant sur 132 ingénieurs utilisant Claude Code a mesuré une augmentation de 67 % du nombre de Pull Requests individuelles fusionnées par jour [^anthropic-study]. Pourtant, les métriques globales de livraison de l'organisation (throughput, DORA metrics) n'ont pas évolué. C'est le « paradoxe de la productivité » : l'augmentation de la vitesse de génération (AGIR) a entraîné une augmentation proportionnelle de la charge de revue (le temps de code review a augmenté de 91 %) parce que la boucle de validation (CORRIGER) n'a pas suivi. Le maillon faible s'est déplacé de l'ingénieur écrivant le code vers l'ingénieur devant auditer un volume massif de code généré de manière asynchrone. Des plateformes d'observabilité comme Axify et Faros AI constatent le même phénomène : sans refonte de l'infrastructure de test et de revue, le taux de retravail augmente [^axify] [^faros].

Customiser un seul axe ne sert à rien. Donner des dizaines d'outils à un agent qui ne comprend pas votre contexte, c'est lui confier les clés d'une usine sans lui donner le plan. Lui donner toute la visibilité sans boucle de correction, c'est le laisser produire sans contrôle qualité. L'équilibre entre les trois leviers est ce qui détermine si l'agent est un outil de production ou un générateur de pull requests à retravailler.

[^anthropic-study]: « How AI Is Transforming Work at Anthropic », Anthropic Research, 2025. Étude portant sur 53 entretiens et 200 000 transcripts de sessions Claude Code (février–août 2025). https://www.anthropic.com/research/how-ai-is-transforming-work-at-anthropic
[^axify]: « AI coding tools' impact: Metrics, ROI, and Review Signals in 2026 », Axify. https://axify.io/blog/ai-coding-tools-impact
[^faros]: « Measuring Claude Code ROI », Faros AI. https://www.faros.ai/blog/how-to-measure-claude-code-roi-developer-productivity-insights-with-faros-ai

---

## 2. Le contexte est une ressource physique : packagez pour scaler

La fenêtre de contexte d'un agent est une ressource fixe et coûteuse. La plupart des moteurs d'agents opèrent aujourd'hui autour de 200 000 tokens effectifs ; certains modèles annoncent des fenêtres bien plus larges (Gemini 2.5 revendique jusqu'à 1 million), mais la fenêtre *exploitable* par un agent en session de travail reste contrainte par le coût et la latence. Tout ce qu'on y injecte (outils, instructions, conversation, fichiers lus) **entre en compétition avec le travail réel**.

Comme dans tout système contraint, chaque octet compte (je ne fais pas tourner un `npm install` sur un Arduino). Un coût caché s'y ajoute : l'évolution des tokenizers peut influer sur la facture sans changement de comportement. Le budget token n'est pas seulement une question de ce que vous injectez, mais aussi de comment le modèle le consomme.

Ce n'est pas un artefact d'un outil particulier. C'est une conséquence de l'architecture transformer (le mécanisme d'attention et le KV cache) qui sous-tend **tous** les moteurs actuels. Deux règles en découlent.

### Règle 1 : Ne payez pas pour ce que vous n'utilisez pas

Chaque token injecté dans le contexte a un coût, direct (facturation) et indirect (il occupe de l'espace que la tâche aurait pu utiliser). L'anti-pattern classique : une équipe injecte 50 000 tokens de documentation d'architecture à chaque tâche, y compris pour un simple renommage de variable. Le contexte est saturé avant même que l'agent ne commence à travailler.

### Règle 2 : Stable en haut, volatile en bas

Pour générer chaque token, le modèle recalcule les relations d'attention entre tous les tokens précédents. Le **KV cache** (Key-Value cache) stocke le résultat de ces calculs pour ne pas les refaire à chaque tour. C'est ce qui permet à une conversation longue de rester fluide et abordable. Mais ce cache fonctionne séquentiellement : modifier quelque chose tôt dans le prompt invalide tout ce qui suit. Le coût de cette invalidation est documenté publiquement :

| Fournisseur | Réduction du coût (cache hit vs miss) | Détail |
| :---- | :---- | :---- |
| Anthropic (Claude) | **90 %** | Cache read à 0.1× le prix de base. Surcharge de 25 % à la première écriture. TTL 5 min [^anthropic-cache] |
| DeepSeek | **~90 %** | Cache hit à ~0.003 $/M tokens vs ~0.14 $/M tokens en miss [^deepseek-price] |
| Google (Gemini) | **~75 %** | TTL configurable, optimisé pour contextes massifs [^llm-pricing] |

Dans les flux de travail agentiques où un agent effectue des dizaines d'allers-retours pour une seule tâche, le prefix caching réduit la consommation financière de plus de 80 % par session [^mindstudio-cache]. L'enjeu n'est pas marginal. D'où la règle de placement :

[^anthropic-cache]: « Prompt caching », documentation Claude API. https://platform.claude.com/docs/en/build-with-claude/prompt-caching
[^deepseek-price]: « DeepSeek Pricing Explained », Flowith Blog. https://flowith.io/blog/deepseek-pricing-explained-most-tokens-per-dollar/
[^llm-pricing]: « LLM API Pricing 2026 », PE Collective. https://pecollective.com/blog/llm-api-pricing-comparison/
[^mindstudio-cache]: « What Is Prompt Caching in Claude Code? », MindStudio. https://www.mindstudio.ai/blog/prompt-caching-claude-code-token-savings

![Règle de placement : stable en haut, volatile en bas](/assets/voir-agir-corriger/context-placement.svg)

Pour rendre cela concret, voici la différence entre un contexte bien structuré et un contexte mal structuré :

![Comparaison d'un contexte bien structuré vs mal structuré](/assets/voir-agir-corriger/context-structure.svg)

Le premier schéma recalcule la quasi-totalité du contexte à chaque tour. Le second ne recalcule que la partie basse, celle qui change de toute façon.

**Avertissement : ce contrôle n'existe pas partout.** Dans les agents CLI (Aider, Claude Code), l'ingénieur maîtrise l'ordonnancement via des commandes explicites et des fichiers de configuration à la racine (CLAUDE.md). Dans les IDE agentiques (Cursor, Copilot Workspace), l'utilisateur définit des intentions via des fichiers de règles (.cursorrules, copilot-instructions.md), mais c'est le moteur RAG interne qui décide quels fragments sont injectés et dans quel ordre. Dans les agents autonomes cloud (Devin), le packaging est une boîte noire inaccessible [^tools-comparison]. Avant d'investir dans l'optimisation du placement, vérifiez que votre outil vous en donne le contrôle.

[^tools-comparison]: « Every AI Coding Tool Compared », Developers Digest. https://www.developersdigest.tech/blog/ai-coding-tools-comparison-matrix-2026

### Packager pour scaler

Passer à l'échelle, c'est traiter le contexte comme du code embarqué. Structurer. Hiérarchiser. Charger à la demande.

Pour chaque ressource que vous injectez dans le contexte, posez-vous la question : **« est-ce que cette approche tient si le projet grossit d'un facteur 10 ? »** Si aujourd'hui vous chargez en bloc vos 20 conventions, vos 5 ADRs et les descriptions de vos 8 outils MCP, ça passe. Mais quand le projet aura 200 conventions et 50 outils, cette même approche saturera le contexte avant que l'agent ne commence à travailler. La bonne abstraction charge ce dont l'agent a besoin, quand il en a besoin, pas tout, tout le temps.

Les mécanismes varient selon les outils : fichiers d'instructions scopés par répertoire, skills chargées paresseusement (lazy loading), recherche d'outils à la demande plutôt que catalogue exhaustif. Mais le principe est le même partout : **packagez le contexte comme du code embarqué**, avec la même discipline d'optimisation.

Un dernier point sur la connaissance : le fine-tuning est rarement la bonne réponse pour injecter du contexte métier. Gekhman et al. (EMNLP 2024) montrent que les LLMs acquièrent l'essentiel de leurs connaissances factuelles durant le pré-entraînement. Le fine-tuning supervisé apprend principalement la forme et le format de l'interaction. Lorsque de nouvelles connaissances factuelles finissent par être assimilées via le fine-tuning, elles augmentent de façon linéaire la propension du modèle à halluciner [^gekhman].

Le fine-tuning conserve une supériorité pour les compétences comportementales et stylistiques : conventions syntaxiques spécifiques, maîtrise d'un DSL non représenté dans les données d'entraînement, adaptation du ton de revue de code. Mais pour de la connaissance factuelle volatile (conventions d'équipe, état de l'architecture, décisions récentes), il est trop coûteux, trop lent, et les modèles évoluent trop vite pour amortir l'investissement. La voie privilégiée reste l'in-context learning, c'est-à-dire des fichiers texte bien conçus, injectés au bon moment dans le contexte.

[^gekhman]: Gekhman et al., « Does Fine-Tuning LLMs on New Knowledge Encourage Hallucinations? », EMNLP 2024. https://aclanthology.org/2024.emnlp-main.444/

---

## 3. La question qui survit à tous les outils

Les agents de code proposent une multitude de mécanismes d'extension : fichiers d'instructions, prompt files, outils MCP, hooks, scripts, skills. Il est tentant de les évaluer un par un. Mais il existe une grille de lecture plus durable.

Tous ces mécanismes se placent sur un seul spectre :

![Le spectre déclaratif / programmatique](/assets/voir-agir-corriger/audit-spectrum.svg)


À gauche : vous décrivez ce que vous voulez. C'est portable, simple à mettre en place, mais fondamentalement fragile : vous *espérez* que le modèle suive vos consignes.

À droite : vous exécutez du code et vous vérifiez le résultat. C'est plus puissant, mais cela demande du travail d'ingénierie et un moteur qui le supporte.

### Compensatoire vs amplificateur

Ce spectre déclaratif/programmatique se croise avec une seconde dimension, indépendante : la **durabilité** de l'investissement face à l'évolution des modèles.

- **Les outils compensatoires** pallient les limites actuelles du modèle. Ils deviennent **moins** utiles quand le modèle progresse.
- **Les outils amplificateurs** exploitent les capacités du modèle. Ils deviennent **plus** utiles quand le modèle progresse.

Ces deux axes sont indépendants. Un outil peut être déclaratif *et* amplificateur, ou programmatique *et* compensatoire :

|  | **Compensatoire** (perd de la valeur quand le modèle progresse) | **Amplificateur** (gagne de la valeur quand le modèle progresse) |
| :---- | :---- | :---- |
| **Déclaratif** | Instruction de 200 lignes listant les fichiers interdits | Description des invariants architecturaux du système ; ADRs injectés en contexte |
| **Programmatique** | Hook qui force un formatage spécifique que le modèle futur comprendra nativement | Hook pre-commit qui lance les tests et réinjecte les échecs ; scanner de sécurité qui remonte ses findings |

Cette distinction est formalisée dans la recherche en interaction humain-machine sous les termes d'effet « égalisateur » vs « amplificateur cognitif ». Des études empiriques montrent que l'IA générative agit comme outil compensatoire (égalisateur) sur les tâches routinières (nivelant les performances entre novices et experts) mais comme pur amplificateur cognitif sur les tâches complexes, où la qualité de la sortie dépend de la capacité de l'humain à définir le problème et guider le raffinement itératif [^equalizer]. Le concept de « scaffolding » (échafaudage), omniprésent dans les frameworks comme LangChain, est intrinsèquement compensatoire : c'est une structure procédurale temporaire conçue pour pallier l'incapacité momentanée du modèle. L'histoire récente le confirme : les frameworks de chaînage de requêtes construits en 2023 pour contourner les fenêtres de 4K tokens de GPT-3.5 sont devenus des béquilles inutiles face aux fenêtres de 200K+ tokens.

La question à se poser pour chaque extension que vous construisez :

**« Sera-t-elle plus utile ou inutile quand le modèle sera deux fois plus intelligent ? »**

Si la réponse est « moins utile », vous investissez dans de la dette transitoire. Si c'est « plus utile », vous construisez un avantage durable. Cette question vaut indépendamment de la position sur l'axe déclaratif/programmatique.

L'implication pour vos investissements est directe : **priorisez les outils amplificateurs**, qu'ils soient déclaratifs ou programmatiques. Et gardez en tête le risque symétrique : un amplificateur non gouverné est un amplificateur d'erreurs systémiques. Si l'agent génère du code intégrant des architectures obsolètes ou des bibliothèques vulnérables, et que l'ingénieur l'utilise sans jugement critique, l'outil cimente et amplifie les pires pratiques à l'échelle de l'organisation. C'est pourquoi la recommandation d'investir dans les amplificateurs est indissociable de l'impératif de contrôle (section 4).

[^equalizer]: « AI as Cognitive Amplifier: Rethinking Human Judgment in the Age of Generative AI », arXiv:2512.10961. https://arxiv.org/abs/2512.10961

---

## 4. Les garde-fous : permettre l'autonomie dans le respect des règles

Les fournisseurs d'agents investissent massivement dans les leviers VOIR (fenêtres de contexte plus larges, RAG, intégrations) et AGIR (plus d'outils, plus d'autonomie). Le levier CORRIGER reste en retrait. Les raisons sont structurelles : la vérification robuste est spécifique à chaque organisation (vos règles, vos conventions, vos contraintes réglementaires), et elle entre en tension avec le narratif commercial de l'autonomie. Regardez les pages produit des principaux agents : le nombre d'intégrations et la taille de la fenêtre de contexte sont mis en avant ; les mécanismes de vérification et de contrôle sont rarement mentionnés.

C'est votre angle mort. Et c'est votre responsabilité.

### L'autonomie n'est pas un objectif, c'est un moyen

L'objectif n'est pas que l'agent fasse ce qu'il veut. L'objectif final reste de produire des solutions logicielles fiables. L'autonomie de l'agent n'est que le moyen d'accélérer cette production, d'améliorer efficacité et/ou efficience en le faisant travailler **dans le cadre que vous avez défini**, sans supervision constante. La nuance est fondamentale.

Observez le contraste : le code de votre extension VS Code passe par une CI complète (lint, typecheck, tests, intégration, sécurité), une couverture de tests minimale et des analytics hebdomadaires. Le code généré par votre agent ? Souvent aucune vérification post-génération, aucune feedback loop programmatique, aucune observabilité (tokens, coût, qualité), aucun audit trail.

Poser des garde-fous n'est pas brider l'agent. C'est **ce qui permet l'autonomie**. Un agent sans rails est un risque qu'on surveille en permanence. Un agent avec des hooks qui bloquent les actions interdites, des linters qui corrigent en temps réel, des tests qui valident avant chaque commit et des permissions granulaires. Cet agent-là, on peut le laisser travailler seul.

Les garde-fous ne ralentissent pas l'agent. Ils lui permettent de courir plus vite sur une piste balisée.

### Six dimensions à forcer

Aucun fournisseur ne traite correctement ces questions, parce qu'elles sont commercialement inconfortables :

| Dimension | Question à se poser |
| :---- | :---- |
| **Autonomie** | L'agent propose ou exécute ? Avec quelle réversibilité ? |
| **Responsabilité** | Qui répond quand l'agent a tort ? |
| **Vérification** | Où se déplace le goulot quand la production casse ? |
| **Coût / valeur** | Le compute consommé justifie-t-il le gain réel ? |
| **Sécurité** | Quelle surface d'attaque ouvre chaque accès donné ? |
| **Observabilité** | Combien de tokens par tâche ? Quel taux de succès des boucles de correction ? Quelle dérive progressive dans la qualité du code ? |

Le principe directeur :

**« Chaque accès supplémentaire donné à l'agent exige un investissement proportionnel en vérification. »**

Un accès sans vérification n'est pas un gain de productivité. C'est une dette de sécurité. Exemple : une équipe active le mode autonome et donne à l'agent un accès en lecture à Slack, sans audit post-action. L'agent, exposé à un message contenant une injection de prompt, modifie un fichier de déploiement. Personne ne s'en aperçoit. Contre-exemple : la même équipe, avec des hooks de validation sur les fichiers critiques, des tests avant chaque commit, des permissions scopées par répertoire. L'agent travaille seul, vite, et dans les règles.

Ce risque est documenté : des études montrent des taux de réussite de 84 % pour les injections indirectes de prompt sur des éditeurs de code agentiques [^aishell], et l'OWASP (organisation de référence en sécurité applicative) classe ces vecteurs dans son Top 10 LLM 2025 [^owasp].

[^aishell]: « Your AI, My Shell: Demystifying Prompt Injection Attacks on Agentic AI Coding Editors », arXiv:2509.22040. https://arxiv.org/abs/2509.22040
[^owasp]: OWASP Top 10 for LLM Applications v2025. https://owasp.org/www-project-top-10-for-large-language-model-applications/

### Les garde-fous ne sont pas infaillibles : la défense en profondeur

Soyons honnêtes : un hook de validation ne détectera pas une injection de prompt sophistiquée qui produit une modification syntaxiquement correcte mais sémantiquement malveillante. Les garde-fous programmatiques sont une couche de défense, pas un pare-feu absolu.

C'est pourquoi la posture recommandée est la **défense en profondeur**, c'est-à-dire plusieurs couches indépendantes, chacune attrapant ce que les autres laissent passer :

- **Moindre privilège** : l'agent n'accède qu'aux ressources strictement nécessaires à sa tâche. Pas d'accès réseau large, pas de credentials globaux.
- **Isolation** : l'agent travaille dans un sandbox (container, worktree dédié) dont les effets de bord sont limités par construction.
- **Audit log immuable** : chaque action de l'agent est journalisée de façon non-altérable. Même si une action malveillante passe, elle est traçable et réversible.
- **Revue humaine sur les chemins critiques** : les modifications de fichiers de déploiement, d'infrastructure ou de sécurité déclenchent une revue explicite, quel que soit le niveau d'autonomie par ailleurs.

Aucune couche seule ne suffit. C'est leur combinaison qui rend l'autonomie de l'agent raisonnablement sûre.

---

## 5. Quand les agents se multiplient : le modèle à plusieurs

Tout ce qui précède raisonne sur un agent unique, en session unique. Mais l'état de l'art pousse vers les architectures **multi-agents** : un orchestrateur qui répartit le travail entre sous-agents spécialisés, chacun dans son propre contexte (un worktree, un container, une branche).

Le modèle Voir / Agir / Corriger reste valide, mais il se **distribue**, et c'est là que de nouvelles questions apparaissent.

### Les trois leviers se répartissent entre agents

Un orchestrateur et ses sous-agents ne partagent pas les mêmes leviers au même niveau :

- **VOIR** se segmente : chaque sous-agent ne voit que son périmètre (un module, un fichier, une tâche). L'orchestrateur, lui, doit maintenir une vision d'ensemble, mais son contexte est tout aussi fini. Le packaging devient critique : que transmettre à chaque sous-agent, et sous quelle forme ?
- **AGIR** se parallélise : plusieurs agents modifient le code simultanément. L'isolation (worktrees, branches séparées) évite les conflits, mais la **réconciliation** (merge, résolution de conflits sémantiques) devient un problème à part entière.
- **CORRIGER** se connecte en chaîne : le résultat d'un sous-agent est le signal d'entrée de l'orchestrateur. Le levier CORRIGER du sous-agent (ses tests, ses linters) alimente le levier VOIR de l'orchestrateur. Si cette boucle est mal conçue (par exemple si le sous-agent renvoie un résumé textuel au lieu d'un résultat structuré et vérifiable), l'orchestrateur perd sa capacité de correction.

### Le piège de la complexité

La dégradation du contexte (*Context Degradation*) est le risque principal des architectures multi-agents. Une recherche de JetBrains présentée au workshop NeurIPS 2025 sur le Deep Learning pour le Code (« The Complexity Trap ») démontre que l'accumulation itérative de traces d'exécution, d'erreurs de terminal et de sorties d'outils sature la fenêtre de contexte, provoquant une chute d'efficacité et une explosion des coûts [^complexity-trap].

L'approche instinctive (utiliser un sous-agent LLM pour résumer l'historique sémantique) est souvent un piège de complexité en soi. L'étude montre qu'une stratégie plus rustique, le « masquage d'observation » (*observation masking*, troncature et omission des anciennes sorties brutes d'outils), réduit les coûts de plus de 50 % tout en maintenant voire améliorant le taux de réussite. Comme le résume un analyste : « Votre fichier d'instructions n'est pas une base de connaissances, c'est un budget d'attention » [^agents-md].

### Le maillon faible se déplace

Dans une architecture multi-agents, le goulot n'est plus nécessairement dans un levier d'un agent unique. Il se déplace vers les **interfaces entre agents** : la qualité du prompt que l'orchestrateur envoie (VOIR du sous-agent), la structure du résultat que le sous-agent renvoie (VOIR de l'orchestrateur), et surtout la question de **qui vérifie le travail final** quand aucun agent individuel n'a vu l'ensemble.

Le principe du maillon faible reste le même, mais l'unité d'analyse change : ce n'est plus un levier d'un agent, c'est le **maillon le plus faible de la chaîne d'agents**.

Sur le plan de la sécurité, la dynamique connue sous le nom de *weakest-link dynamics* s'applique : un seul serveur MCP compromis fournissant des définitions d'outils empoisonnées peut contaminer l'ensemble de la chaîne [^multi-agent-sec]. Les frontières entre agents doivent être conçues comme des disjoncteurs (*circuit breakers*), pas comme de simples divisions du travail.

[^complexity-trap]: « The Complexity Trap », JetBrains Research, NeurIPS 2025 Workshop on Deep Learning for Code. https://arxiv.org/abs/2508.21433
[^agents-md]: « Your AGENTS.md is a Liability », paddo.dev. https://paddo.dev/blog/your-agents-md-is-a-liability/
[^multi-agent-sec]: « Open Challenges in Multi-Agent Security », arXiv:2505.02077. https://arxiv.org/abs/2505.02077

---

## 6. Quatre étapes pour évaluer n'importe quel agent

Les principes précédents se condensent en une méthode applicable à n'importe quel moteur : Claude Code, Copilot, Cursor, Cody, Aider, ou ce qui viendra demain.

**Étape 1 – Mapper.** Identifiez les primitives de l'agent et placez-les sur les trois leviers. Quels mécanismes alimentent VOIR ? Lesquels étendent AGIR ? Lesquels font tourner CORRIGER ? Si un levier n'a aucun mécanisme, vous avez trouvé votre goulot d'étranglement.

**Étape 2 – Identifier qui contrôle le packaging.** L'ingénieur contrôle-t-il ce qui entre dans le contexte, ou est-ce le fournisseur qui décide ? Cette réponse détermine si l'économie du contexte est un levier actionnable pour vous, ou une boîte noire que vous subissez. Le niveau de contrôle varie radicalement :

| Niveau de contrôle | Outils | Caractéristiques |
| :---- | :---- | :---- |
| **Transparence totale** (CLI) | Aider, Claude Code | Contrôle chirurgical : commandes /add, /drop, fichier CLAUDE.md. L'ingénieur décide de chaque fragment injecté |
| **Contrôle déclaratif intermédié** (IDE) | Cursor, Copilot Workspace | Fichiers de règles (.cursorrules, copilot-instructions.md), mais le moteur RAG interne orchestre le packaging final |
| **Boîte noire** (Cloud) | Devin, Factory | Prompt système, modèle sous-jacent et règles d'inclusion du contexte propriétaires et inaccessibles |

**Étape 3 – Classer sur le spectre.** Pour chaque extension que vous envisagez de construire, posez la question compensatoire / amplificateur. Priorisez les investissements amplificateurs – ce sont les seuls dont la valeur croît avec le temps.

**Étape 4 – Forcer la couche Contrôle.** Concevez activement les limites d'autonomie, le suivi des coûts, les frontières de sécurité et les pipelines de vérification. Aucun fournisseur ne le fera à votre place.

### Le trade-off stratégique : Contrôle vs Intégration

Derrière ces étapes se cache un arbitrage central. Certains agents vous donnent le **contrôle** total : vous décidez du packaging du contexte, vous construisez vos boucles de feedback, vous choisissez vos outils. D'autres vous donnent l'**intégration** : l'écosystème est préconfiguré, les connexions existent, mais vous acceptez les choix du fournisseur.

Le bon choix dépend de deux variables : la spécificité de votre contexte et la force de vos contraintes d'intégration.

|  | **Contraintes d'intégration faibles** | **Contraintes d'intégration fortes** |
| :---- | :---- | :---- |
| **Contexte propriétaire fort** (conventions spécifiques, compliance, outillage interne) | **Contrôle** : vous avez besoin de packager un contexte que personne d'autre ne comprend. Exemple : une équipe fintech avec des règles de compliance spécifiques à son régulateur. | **Hybride** : vous avez besoin des deux. L'intégration native pour les flux standard, le contrôle pour le contexte métier. Coût le plus élevé, mais parfois inévitable. |
| **Contexte propriétaire faible** (stack standard, conventions classiques) | **Léger** : un agent vanilla avec quelques fichiers d'instructions suffit. Ne sur-investissez pas. Exemple : une startup early-stage sur une stack React/Node standard. | **Intégration** : votre contexte n'est pas spécial, mais vos outils doivent communiquer. Exemple : une équipe full-GitHub qui veut que l'agent gère PRs, issues et CI sans configuration. |

L'erreur la plus fréquente : choisir le contrôle total par réflexe d'ingénieur, alors que le contexte propriétaire ne le justifie pas. Le contrôle a un coût d'opportunité : chaque boucle de feedback que vous construisez est une boucle que vous maintenez.

La méthode en quatre étapes, elle, est indépendante du choix. Elle fonctionne dans les deux cas.

---

## Pour conclure

Trois principes durables résument ce cadre :

1. **Un agent n'est jamais meilleur que son levier le plus faible** : ce qu'il voit, ce qu'il peut faire, ce qui le corrige.
2. **Investissez dans les outils amplificateurs**, pas dans les compensatoires. Ce sont les seuls qui gagnent en valeur quand les modèles progressent.
3. **La couche de contrôle est votre responsabilité.** Les garde-fous ne brident pas l'agent ; ils sont la condition de son autonomie.

Les modèles continueront de progresser. Une partie de l'ingénierie de contexte actuelle deviendra inutile, dette transitoire que le modèle absorbera. Mais le cadre en trois leviers, le test compensatoire/amplificateur et l'impératif de contrôle resteront valides. Ce ne sont pas des détails d'implémentation. Ce sont des heuristiques structurelles suffisamment robustes pour avoir survécu à plusieurs générations d'outils. Elles ont des conditions d'invalidation : des agents avec mémoire persistante et apprentissage en ligne, des agents multi-modaux opérant sur des interfaces visuelles, ou des agents capables d'auto-réflexion sans signal de correction externe pourraient exiger un modèle étendu. Pour l'état de l'art actuel, les trois leviers tiennent.

**Customiser un agent de code, ce n'est pas tout lui donner. C'est décider, sous contrainte de coût, ce qu'il voit, ce qu'il peut faire, et ce qui le corrige.**

---

## Annexe : les anti-patterns les plus fréquents

| Anti-pattern | Levier défaillant | Symptôme | Correction |
| :---- | :---- | :---- | :---- |
| 50K tokens de docs injectés à chaque tâche | VOIR (saturation) | Contexte épuisé avant le début du travail | Lazy loading, scoping par répertoire |
| 15 outils MCP, aucun linter | CORRIGER (absent) | Code plausible, casse en CI | Hook pre-commit + tests automatiques |
| Mode autonome + accès Slack non audité | CORRIGER (absent) | Exposition à l'injection indirecte | Permissions scopées, audit log, revue sur chemins critiques |
| Contrôle total par réflexe d'ingénieur | Sur-investissement | Chaque boucle construite = boucle maintenue | Évaluer le contexte propriétaire avant d'investir (matrice section 6) |
| Résumé textuel entre sous-agents | VOIR de l'orchestrateur (dégradé) | Perte de vérifiabilité, dérive | Résultat structuré et vérifiable (exit code, JSON, diff) |
| Aucune mesure du coût par tâche | Observabilité (absente) | Risque non quantifié déguisé en productivité | Tracking tokens/tâche, taux de retravail, dérive architecturale |
