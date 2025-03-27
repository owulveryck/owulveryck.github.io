---
title: "Repenser les présentations : Au-delà des diapositives statiques"
date: 2023-10-10T07:55:21+02:00
lastmod: 2023-10-10T07:55:21+02:00
draft: false
images: [/assets/crowdasleep_small.png]
videos: [/assets/present.webm]
keywords: []
summary: À l'ère numérique, les présentations PowerPoint traditionnelles ne parviennent souvent pas à capter l'attention des audiences en raison de leur nature statique. 
  

  La recherche suggère que l'attention du public s'estompe après seulement 10 minutes sans engagement. 
  
  La méthode proposée dans cet article vise à revitaliser les présentations en
  
  * Incorporant le dessin en direct avec l'aide d'outils comme la tablette reMarkable pour une interaction en temps réel.
  
  * Utilisant un script pour créer des PDF à partir d'images, combinant la structure familière des diapositives avec le dessin spontané sur le moment.
   

  Le résultat est une expérience de présentation plus authentique, engageante et percutante, bien qu'elle nécessite une préparation plus approfondie et une compréhension du sujet. 
  
  Les outils et méthodes mis en évidence visent à déplacer l'accent de l'esthétique vers un véritable engagement avec le contenu.
tags: ["keynote", "reMarkable", "powerpoint", "presentation", "talk"]
categories: ["tools"]
author: "Olivier Wulveryck"

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: false
toc: true
autoCollapseToc: true
# You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
contentCopyright: false
reward: false
mathjax: false

# Uncomment to add to the homepage's dropdown menu; weight = order of article
# menu:
#   main:
#     parent: "docs"
#     weight: 1
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/20231010-nextgen-presentation/).

## Introduction
Dans le monde rapide d'aujourd'hui, les présentations traditionnelles ne parviennent souvent pas à engager les audiences.

Leur nature statique et préparée à l'avance manque de l'interactivité et du dynamisme nécessaires pour capter l'attention.

Nous n'exposons plus des "_power points_" ou des "_key notes_" comme nous le faisions à l'époque des _transparents_ et des _rétroprojecteurs_ 
(voir l'impressionnante présentation [_growing a language_](https://www.youtube.com/watch?v=_ahvzDzKdB0) de _Guy L. Steele Jr._ comme illustration, même si le contenu n'est pas lié à cet article, son utilisation des transparents est intéressante).

Les **diapositives** ont évolué **d'un rôle de soutien** pour devenir le centre d'attention, où **elles sont la présentation**.

J'ai assisté à d'innombrables présentations ennuyeuses où les participants sont souvent trouvés endormis devant des diapositives statiques sans fin.

![Illustration d'une foule somnolente en écoutant une présentation monotone](/assets/crowdasleep_small.png)

En revanche, lors des conférences technologiques où j'ai le plaisir d'assister à des démonstrations en direct, le public est visiblement plus enthousiaste. Ce n'est pas une surprise que le framework _scrum_ intègre des démos.

Cependant, toutes les présentations ne peuvent pas être centrées sur la démonstration. Donc, j'explore une nouvelle approche pour améliorer l'interactivité et l'engagement tout en maintenant un format structuré.

Cet article explore mon approche et les outils que j'utilise. Bien que je ne prétende pas avoir **la méthode définitive** (car il en existe de nombreuses autres), je partage mes techniques personnelles sur ce blog technique dans l'espoir qu'elles puissent être utiles à quelqu'un.

## La situation 

### La science derrière l'engagement

Selon "Brain Rules" du Dr. John Medina, l'attention commence à diminuer après environ 10 minutes sans engagement[^1].

Ainsi, un exposé devrait soit être plus court que 10 minutes, soit il devrait y avoir des efforts pour réengager le public.

Une étude menée par l'Université Carnegie Mellon a démontré que les méthodes d'apprentissage actif améliorent significativement les performances académiques, mettant l'accent sur les cadres interactifs plutôt que sur les cours traditionnels[^2].

Par conséquent, dans un exposé durant plus de 10 minutes, engager le public à travers un processus d'apprentissage actif peut aider à maintenir l'attention.

### Appliquer la science

Je possède une tablette reMarkable et, comme expliqué précédemment sur mon blog, je l'utilise pour les appels vidéo.

Maintenant que la période de télétravail complet est terminée et que les gens reviennent aux conférences, je peux étendre cette méthode de présentation pour l'utiliser en direct, dans la vie réelle.

Je l'ai déjà utilisée dans une série de BBLs comme complément à certaines diapositives, et les retours ont été positifs.

Je peux aller plus loin en remplaçant les diapositives traditionnelles par des éléments stockés sur la reMarkable, me permettant d'écrire dessus, tout en laissant des pages blanches entre elles pour dessiner des informations complémentaires.

Je peux inciter le public à réagir et adapter le dessin, favorisant l'apprentissage actif.

![](/assets/present_looped.webp)

### Les compromis

Tandis que les diapositives traditionnelles fournissent un filet de sécurité pour les présentateurs, les aidant à rester sur la bonne voie, le dessin en direct nécessite une préparation approfondie. 
On doit avoir une compréhension profonde de son sujet pour maintenir la fluidité et la confiance tout au long de la présentation.

Et soyons honnêtes, bien que je sois relativement fluide lorsque j'ai besoin d'expliquer quelque chose dans lequel je suis compétent à un petit groupe, présenter devant une foule pose un défi différent.
Je peux difficilement compter sur les retours (qu'ils soient intentionnels ou non).

Ainsi, j'ai besoin de maintenir un certain cadre fourni par l'outillage.

## Outillage et ingénierie : Donner vie à la présentation

Explorons maintenant ma boîte à outils. Je vais l'essayer demain à [Cloud Nord](https://cloudnord.fr/) dans ma ville natale. Je sais que les gens sont gentils et amicaux ici, ils pardonneront toute erreur.

### Streaming avec reMarkable

L'un des outils clés permettant ce style de présentation dynamique est la tablette reMarkable.

Ce n'est pas simplement un appareil pour la prise de notes ou la lecture ; c'est un puissant outil de streaming qui met en avant l'essence du dessin en direct.

Avec [goMarkableStream](https://github.com/owulveryck/goMarkableStream), je peux facilement diffuser le contenu de la tablette.

Récemment, j'ai ajouté une fonctionnalité qui permet de diffuser sur Internet en intégrant la capacité de tunneling [ngrok](https://ngrok.com/).
Cela permet de diffuser sur un ordinateur qui n'est pas nécessairement sur le même réseau que la tablette, bien que cela nécessite une connectivité Internet décente.

J'ai essayé le partage de connexion sur un réseau mobile médiocre, et le résultat n'était pas assez fiable pour une présentation. Cependant, si la salle dispose d'un wifi adéquat, cela peut suffire.

Néanmoins, chaque fois que possible, j'adhère à une connexion filaire à mon ordinateur portable et je diffuse depuis mon ordinateur portable. C'est un choix plus sûr.

### Encadrer la présentation dans un PDF

Pour fusionner le traditionnel avec le nouveau, j'emploie un script qui convertit une collection d'images en PDF.

Cela simule la structure familière des diapositives mais avec une touche – ces "diapositives" peuvent être écrites, assurant que bien qu'il y ait une base pour la présentation, l'interaction en temps réel reste intacte.

[Ce gist sur GitHub](https://gist.github.com/owulveryck/1317f9b22433aa18778b673000159141) présente un `Makefile` qui prend un ensemble d'images et les convertit dans un format adapté à la reMarkable.
Il convertit les images en niveaux de gris et à une résolution de 1872x1404.

Ensuite, il ajoute une bordure et une annotation au bas de chaque page (juste parce que j'aime le style).

Après cela, il assemble les images en un PDF.

Pour obtenir l'ordre correct, le script lit les images à partir d'un fichier nommé `slides.txt`.

Je peux également ajouter une page blanche (et la page blanche est également générée par le script).

Je tenterai de partager le résultat ici une fois que j'aurai créé une présentation complète avec cela.

## Conclusion

Le vieil adage, "La forme suit la fonction", est valable même pour les présentations. 
Alors que la méthode traditionnelle privilégie l'esthétique, l'approche dynamique met l'accent sur le contenu et l'engagement authentique.

Elle recentre la présentation sur la communication de l'essence d'un sujet au lieu de simplement lire des diapositives.

Cette approche peut demander plus d'efforts mais promet une expérience de présentation plus authentique et percutante.

Une autre possibilité est d'enregistrer une présentation pour la diffusion... C'est en fait ce que j'ai fait pour [cette illustration](/assets/present.webm).

---

### Références

[^1]: [Brain Rules, règle des 10 minutes](https://brainrules.blogspot.com/2009/03/10-minute-rule.html)

[^2]: [Une nouvelle recherche montre que l'apprentissage est plus efficace lorsqu'il est actif](https://www.cmu.edu/news/stories/archives/2021/october/active-learning.htm)