---
title: "Comment activer l'effet volant de valeur avec vos données"
date: 2024-10-15T12:15:33+01:00
lastmod: 2024-10-15T12:15:33+01:00
images: [/assets/value-flywheel/value-flywheel-data.webp]
draft: false
keywords: []
summary:
tags: []
categories: []
author: "Olivier Wulveryck"
comment: false
toc: true
autoCollapseToc: false
# You can also define another contentCopyright.
contentCopyright: false
reward: false
mathjax: false
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/20241015-value-flywheel-data/).

Dans le monde hyper-compétitif d'aujourd'hui, les entreprises ne s'appuient plus uniquement sur des décisions instinctives ou intuitives ; elles dépendent d'**insights basés sur les données** pour rester agiles et prendre des décisions rapides et intelligentes.
Cependant, les données seules ne sont pas la réponse ; elles sont le facilitateur pour créer un élan sur un **volant d'inertie commercial et technologique** :
un modèle où les données guident les décisions, les décisions guident les actions, et ces actions génèrent de la valeur, propulsant l'entreprise vers l'avant dans un cycle auto-renforçant.

Dans un [article précédent](https://blog.owulveryck.info/2024/04/09/data-as-a-product-and-data-contract-an-evolutionary-approach-to-data-maturity.html), j'ai utilisé un modèle pour expliquer comment les données pourraient traverser les frontières des applications et des domaines pour apporter une valeur croissante au niveau organisationnel.

Voici une synthèse du développement de cet article : l'axe X concerne l'augmentation de la certitude de la décision corrélée à son offre, et l'axe Y la diffusion des données corrélée à sa demande.

![Une courbe en S représentant l'évolution des données, le X est la certitude et Y est l'ubiquité. Il y a une division : le bas de la S contient des données brutes, le milieu est organisé, et le haut fait autorité. Le point d'inflexion est indiqué comme contrat de données.](/assets/data_certitude.svg)

Dans cet article, j'utilise YAML (_Yet Another ModeL_ - sic) pour déterminer comment une culture, une technologie et une organisation des données peuvent activer un effet volant de valeur d'un point de vue technologique et commercial conjoint.

À la fin, vous obtiendrez des insights actionnables sur l'analyse et la mise en œuvre efficace de ces principes au sein de votre propre organisation.

### Description du problème

Dans le monde de l'**analytique**, les données recèlent un immense potentiel pour générer de la **valeur commerciale** en produisant des insights actionnables. Pourtant, elles restent souvent limitées à être simplement un **facilitateur de solutions** au sein d'un **domaine commercial** spécifique (l'espace problème). Parce que les données sont généralement traitées comme un outil pour résoudre des défis commerciaux isolés, elles franchissent rarement les frontières de leur domaine. Cela limite leur capacité à contribuer à des **solutions technologiques** plus larges ou à générer de la valeur dans l'ensemble de l'organisation.

Ce désalignement provient souvent de **silos organisationnels**. Pire encore, les équipes de données et les équipes commerciales fonctionnent fréquemment de manière isolée, conduisant à des efforts disjoints où les insights, aussi robustes soient-ils, sont soit mal compris, soit sous-utilisés.

Pour remédier à cela, les entreprises doivent se concentrer sur la création d'un **système** qui favorise la collaboration entre ces équipes. C'est là que l'analogie d'un **système de poulie à courroie** devient pertinente. Tout comme un système de poulie transforme le mouvement et la puissance, une approche alignée des données transforme les informations brutes en valeur commerciale.

Imaginez un **volant d'inertie** comme la **roue entraînée** du système—cela représente la **valeur technologique et commerciale** combinée. Les **données** servent de **courroie**, facilitant le transfert de connaissances, de contexte et de valeur à travers l'organisation. La clé, cependant, est d'identifier la **roue motrice** : un volant d'inertie plus petit qui se concentre spécifiquement sur la **maturité des données**, accélérant les résultats commerciaux.

{{< figure src="/assets/value-flywheel/value-flywheel-data.webp" alt="Un système de poulie à courroie : un volant de valeur décrit comme les données sont le moteur, un volant commercial et technologique conjoint est entraîné par une courroie qui alimente le moteur avec des initiatives commerciales. Le volant entraîné est alimenté avec des données-en-tant-que-produit" >}}

Voyons maintenant comment les entreprises peuvent définir et construire cette **roue motrice centrée sur les données** en quatre étapes, qui agira comme le moteur alimentant le volant de valeur plus large de l'organisation.

### **L'effet volant de valeur comme modèle**

#### Le concept original

J'ai rencontré le concept du volant de valeur dans le livre [The Value Flywheel Effect](https://itrevolution.com/product/the-value-flywheel-effect/) de David Anderson, avec Mark McCann et Michael O'Reilly.

J'aime l'introduction avec l'idée d'élan et je vais simplement copier ici le tout début du livre que vous pouvez trouver sur le [site web d'IT Revolution](https://itrevolution.com/articles/what-is-the-value-flywheel/) gratuitement :

> L'élan est une chose étrange.
> Il est difficile d'imaginer ce que cela fera et demande beaucoup d'efforts à réaliser.
> Quand nous apprenons à faire du vélo, par exemple, cela semble d'abord maladroit et gênant.
> Il est difficile de faire tourner les roues au début, et notre frustration est souvent évidente.
> Mais notre professeur nous assure que cela passera.
> Quand nous commençons enfin à prendre de l'élan, l'exaltation nous coupe le souffle.
> Chaque poussée de pédale devient plus facile et demande moins d'effort.
> Soudain, nous pouvons nous concentrer sur l'expérience plus large de glisser à travers une belle forêt ou une rue bordée d'arbres.
> La valeur de notre dur labeur est évidente, et nous pouvons maintenant continuer à en récolter les bénéfices avec de moins en moins de peine.

Pour comprendre pleinement **l'effet volant de valeur**, considérons-le d'abord d'un point de vue commercial plus large.
L'application la plus connue de l'effet volant d'inertie vient d'Amazon, où de petites victoires en satisfaction client mènent à plus de trafic, ce qui attire des vendeurs tiers, ce qui à son tour entraîne des coûts et des prix plus bas, créant un élan supplémentaire.
Chaque élément renforce le suivant, créant un centre d'accélération perpétuelle de croissance, d'efficacité et de valeur différenciée.

Le modèle du volant d'inertie décrit un **système auto-renforçant** d'élan, devenant plus fort et plus facile à maintenir à mesure que différentes actions s'accumulent au fil du temps.
Au lieu de s'appuyer sur des changements ou des modifications gigantesques ponctuels, le volant d'inertie fonctionne en multipliant les effets d'**améliorations petites et cohérentes**.
Une fois qu'il est en mouvement, il est difficile à arrêter.

#### Les principes clés et les personas

Je ferai référence au modèle du livre *The Value Flywheel Effect*, qui est la troisième itération de ce concept (suivant l'adaptation d'Amazon et l'idée originale de [Jim Collins](https://www.jimcollins.com/concepts/the-flywheel.html)). Ce modèle est expliqué en détail dans le livre, et je recommande vivement de le lire. En résumé, le volant de valeur est décomposé en quatre étapes :

**Clarté de l'objectif** → **Défi et paysage** → **Prochaines meilleures actions** → **Valeur à long terme**

Chacune des quatre étapes entraîne la suivante. **Étape par étape**, **itération par itération**, nous surmontons l'inertie (d'une manière vraiment agile) et créons de l'élan.

{{< figure src="/assets/value-flywheel/the-value-flywheel-effect.webp" alt="Le volant de valeur en quatre phases : I clarté de l'objectif, II défi et paysage, III prochaine meilleure action, IV valeur à long terme. La roue relie les quatre phases et boucle de la phase IV à la phase I." title="Illustration de l'effet volant de valeur (Adapté du modèle de David Anderson)" >}}

Pour résumer, voici les principes clés et les personas de l'effet volant de valeur :

{{< figure src="https://itrevolution.com/wp-content/uploads/2022/10/Screen-Shot-2022-08-03-at-3.06.48-PM-1024x819.png.webp)" title="12 principes clés de l'effet volant de valeur - David Anderson, Michael O'Reilly, Mark McCann - (c) ITRevolution" >}}
[Source](https://itrevolution.com/articles/12-key-tenets-of-the-value-flywheel-effect/)

Voyons maintenant si nous pourrions adapter ces principes clés à la roue motrice des données. Le premier élément est de considérer l'approche existante en place dans l'organisation (Adapter > Adopter)

### Le volant de valeur des données : Organiser l'écosystème des données

Pour appliquer avec succès **l'effet volant de valeur** au sein d'un écosystème de données, les entreprises ont besoin d'une approche structurée pour gérer les données à travers l'organisation. Cela est généralement mis en œuvre à travers le concept d'**usine de données**. Une usine de données fonctionne comme le moteur central qui alimente la prise de décision basée sur les données, garantissant que les données circulent librement entre les équipes commerciales et techniques, les transformant en insights actionnables.

À sa base, l'usine de données est responsable du **cycle de vie complet des données** : de la **collecte de données** et de la **curation** au **stockage**, à la **gouvernance** et à l'**exposition** pour les insights. Elle agit comme un pont entre les domaines, permettant le transfert transparent des données depuis des espaces problématiques commerciaux isolés vers des applications technologiques et stratégiques plus larges à travers l'organisation.

Cependant, pour que l'usine de données soit efficace, l'**équipe de données** doit jouer un rôle central au sein de chaque domaine. Selon la gouvernance de l'organisation, il existe divers modèles de gestion des données qui peuvent être adoptés—**centralisé**, **distribué**, **fédéré** ou **maillage de données**. Bien que ces modèles diffèrent dans la façon dont ils gèrent et organisent les données à travers l'entreprise, une constante demeure : il y a toujours une **équipe centrale de gestion des données** responsable d'assurer la cohérence des données, la gouvernance et l'alignement avec les objectifs commerciaux plus larges.

Indépendamment du modèle utilisé, cette équipe centrale a des responsabilités cruciales. Ils servent de gardiens de l'**écosystème de données**, fournissant l'infrastructure, les politiques et les processus nécessaires pour un flux de données et une collaboration fluides.

Avec **l'effet volant de valeur**, nous pouvons découvrir les **principes communs** qui guident l'équipe centrale de gestion des données, lui permettant de générer efficacement de l'élan et d'accélérer la création de valeur dans toute l'organisation.

Maintenant, explorons comment les **quatre phases du volant de valeur** peuvent être appliquées à (potentiellement) toute entreprise aspirant à être axée sur les données.

Pour la clarté et la cohérence, je vais adhérer à la structure de l'effet volant d'inertie original, décrivant chaque phase en termes de ses principes clés, des personnes responsables et de l'objectif principal.

### Phase 1 - Clarté de l'objectif

La première phase est dirigée par le **Chief Data Officer (CDO)**, dont la responsabilité clé est de siéger au niveau exécutif et de définir ce que signifie pour l'entreprise d'être **axée sur les données**. Cette phase consiste à créer une **métrique North Star** qui encapsule la vision de l'entreprise d'être axée sur les données, particulièrement d'un point de vue technologique. Le CDO collabore avec d'autres dirigeants pour comprendre comment les données peuvent soutenir les objectifs stratégiques de l'entreprise, garantissant que chaque initiative s'aligne sur les objectifs commerciaux plus larges.

- **Principes clés** : Alignement sur les **objectifs commerciaux**, établissement d'une **North Star axée sur les données** qui définit clairement ce que signifie pour l'entreprise de tirer parti des données comme avantage concurrentiel et répond à la question : _que signifie pour nous d'être axé sur les données ?_.
- **Personnes en charge** : Chief Data Officer (CDO), Leadership, Experts de domaine.
- **Objectif** : Créer un alignement entre les équipes de données et commerciales, garantissant que les initiatives de données soutiennent les objectifs stratégiques et préviennent le désalignement entre les attentes et les résultats.

Cette clarté d'objectif pose les fondements de toutes les activités ultérieures dans l'écosystème de données, garantissant que les efforts de l'**usine de données** sont toujours axés sur la livraison de valeur commerciale.

### Phase 2 - Défi et paysage

Une fois que l'**objectif** est clair, la phase suivante est une évaluation **dirigée par l'ingénierie** du **paysage technique**. Ici, l'**équipe d'ingénierie** doit évaluer les systèmes et l'infrastructure actuels pour s'assurer qu'ils peuvent répondre aux besoins d'une organisation axée sur les données. Cela inclut la compréhension des capacités techniques de l'entreprise telles que l'utilisation de **lacs de données**, d'**architectures d'échange de données**, et si l'entreprise s'appuie sur une **architecture pilotée par les événements**. L'équipe doit identifier ce qui peut être construit en interne pour soutenir la stratégie de données définie dans la Phase 1 et quels composants sont mieux approvisionnés en tant que **commodités**.

- **Principes clés** : Évaluation de la **dette** technologique, des goulots d'étranglement, de l'infrastructure de données et de la capacité de l'entreprise à stocker, gouverner et partager des données à travers les domaines.
- **Personnes en charge** : Responsables technologiques, Ingénieurs de données, Architectes.
- **Objectif** : Comprendre quelles lacunes existent dans la pile technologique actuelle, garantir que les bonnes fondations sont en place (par exemple, architecture pilotée par les événements, lacs de données), et décider du développement interne par rapport aux services de commodité.

L'**usine de données** dans cette phase doit construire les capacités fondamentales requises pour le succès à long terme, qu'il s'agisse de mettre à l'échelle la pile technologique existante ou d'aborder les déficiences techniques.

### Phase 3 - Prochaine meilleure action

La troisième phase est dirigée par les **équipes produit** qui doivent maintenant se concentrer sur des **victoires rapides** qui peuvent démontrer une valeur commerciale immédiate renforcée par les données. L'objectif ici est d'identifier les **produits de données** ou les **insights de données** qui peuvent **catalyser les initiatives commerciales** à court terme. Ces victoires devraient être étroitement alignées avec les objectifs commerciaux, tels que l'amélioration des prévisions de ventes pour donner un avantage concurrentiel en logistique ou le développement d'insights basés sur les données pour soutenir des initiatives stratégiques clés.

- **Principes clés** : **Progrès itératif**, concentration sur la livraison de **victoires rapides**, responsabilité décentralisée.
- **Personnes en charge** : Chefs de produit, Leaders de domaine, Ingénieurs de données.
- **Objectif** : Développer et déployer des produits de données qui délivrent une **valeur commerciale immédiate**, en se concentrant sur les initiatives qui soutiennent les objectifs commerciaux à court terme tout en construisant de l'élan pour le succès à long terme.

Ici, l'**usine de données** doit fournir un support évolutif et agile pour permettre aux équipes de domaine de développer leurs propres produits de données. L'accent est mis sur la prise de mesures progressives pour débloquer de la valeur, garantissant que les produits de données sont alignés avec les initiatives commerciales.

### Phase 4 - Valeur à long terme

La phase finale se concentre sur la sécurisation de la **valeur à long terme**, dirigée par le **CTO**. À ce stade, l'objectif est de garantir que les systèmes et processus de données en place peuvent soutenir l'évolution continue de l'entreprise. Le CTO travaille en étroite collaboration avec le **CDO** et le conseil d'administration pour garantir que l'**usine de données** continue à se développer comme un facilitateur clé des futurs produits de données et innovations. L'accent est mis sur l'évolutivité, la gouvernance et le maintien de l'avantage concurrentiel acquis en étant une entreprise axée sur les données.

- **Principes clés** : Évolutivité, gouvernance et création d'un **écosystème de données durable**.
- **Personnes en charge** : Chief Technology Officer (CTO), Architectes de données, Équipes de données de domaine.
- **Objectif** : Architecturer un système qui garantit une **valeur à long terme** en permettant aux équipes d'innover continuellement et en faisant de l'usine de données un **fournisseur de plateforme** pour les futurs besoins commerciaux axés sur les données.

Dans cette phase, l'usine de données passe à un rôle plus **de soutien**, permettant une innovation continue tout en garantissant que l'infrastructure est robuste, évolutive et adaptable. Cette phase alimente également les insights au **CDO** et à l'équipe exécutive, garantissant que la stratégie de données reste un différenciateur concurrentiel.

### Conclusion : Activer le volant d'inertie

Prenons un peu de recul pour examiner ce mécanisme :

{{< figure src="/assets/value-flywheel/value-flywheel-data.webp" alt="Un système de poulie à courroie : un volant de valeur décrit comme les données sont le moteur, un volant commercial et technologique conjoint est entraîné par une courroie qui alimente le moteur avec des initiatives commerciales. Le volant entraîné est alimenté avec des données-en-tant-que-produit" >}}

Le **volant de valeur** est activé phase par phase, avec l'**usine de données** centrale jouant un rôle crucial tout au long du voyage vers un système de gestion de données fédéré. Cela garantit une **cohérence inter-domaines** et la livraison de solutions évolutives. En suivant cette approche structurée, le paysage des données se transforme en un cadre cohérent où chaque partie de l'organisation—pas seulement l'équipe centrale de données—contribue au **succès à long terme** de l'entreprise.

Ce modèle fournit un schéma utile pour expliquer le potentiel d'un système de données au sein d'une organisation. Cependant, la mise en œuvre de chaque étape dépasse la portée de cet article. Il existe de nombreux outils disponibles—tels que le cadre North Star et la cartographie Wardley—qui peuvent soutenir l'application de cette structure à chaque étape.