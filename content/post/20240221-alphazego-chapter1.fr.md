---
title: "AlphaZero à partir de zéro Chapitre 1 : MCTS de base"
date: 2024-02-18T17:56:25+01:00
lastmod: 2024-02-18T17:56:25+01:00
draft: true
keywords: []
description: ""
tags: []
categories: []
author: ""

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S. comment can only be closed
comment: false
toc: false
autoCollapseToc: false
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

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/20240221-alphazego-chapter1/).

J'ai toujours été fasciné par l'Intelligence artificielle appliquée aux jeux.

Récemment, j'ai regardé le film ["Wargames"](https://en.wikipedia.org/wiki/WarGames) de 1983. Cela m'a donné envie de jouer à nouveau avec ces jeux de plateau.
J'aime l'idée d'un ordinateur jouant à un jeu sans être explicitement programmé pour le faire.

J'ai également (re)regardé le film ["AlphaGo"](https://www.youtube.com/watch?v=WXuK6gekU1Y&themeRefresh=1). J'ai été fasciné par AlphaGo (et plus tard par AlphaZero).
L'idée principale que j'ai appréciée est celle de l'apprentissage par renforcement : la machine apprend par elle-même en tirant des leçons de ses erreurs.