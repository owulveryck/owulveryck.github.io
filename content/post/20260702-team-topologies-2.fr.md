---
title: "Team Topologies à l'ère agentique : la boucle, la plateforme et l'enabling team"
slug: "sdlc-team-topologies-fr"
date: 2026-07-02T10:00:00+02:00
images: [/assets/sdlc-team-topologies/sdlc-team-topologies.svg]
draft: false
summary: "Une exploration visuelle du SDLC agentique : comment la boucle agentique (capter, planifier, agir, observer) s'inscrit dans une organisation Team Topologies — la stream-aligned team qui conçoit la solution, la plateforme qui propulse et outille le système agentique, et l'enabling team qui pose le cadre et capitalise."
tags: ["architecture", "A2A", "AP2", "agents"]
categories: ["dev"]
author: "Olivier Wulveryck"
toc: false
comment: false
mathjax: false
---

Note préalable: quand je te dis d'implémenter des points, ce sont des indicateurs qui passent le long d'un chemin entre deux éléments. Pour une phase particulière, le flow complet doit tourner en boucle (par exemple, on passe de a -> b -> c -> d et on recommence a -> b -> c -> d)
Entre chaque phase, les éléments affichés doivent rester sauf information contraire par contre les points de la phase doivent être complètement supprimés

## Introduction

{{< scrollytelling svg="/assets/sdlc-team-topologies/sdlc-team-topologies.fr.svg" >}}

{{< scrollytelling-step phase="1" id="phase-1" >}}

## Rappel sur la boucle agentique

Dans cette partie, je ne veux voir que la boucle "système agentique (automome)" (tu vas cacher la partie "Implémente la solution"). Je veux voir toute la boucle avec les cases capter, planifier, agir et observer avec l'auto-correctif et le résultat insatisfaisant et la fleche verte qui sort (et qui pointe sur un résultat et tu affiche la case résultat (mais pas la case stream aligned)).

Je veux un point qui parte de capter et qui passe à planifier et agir. A chaque fois qu'il passe sur une boite, la boite passe en surbrillance pour montrer le passage.
{{< /scrollytelling-step >}}
{{< scrollytelling-step phase="2" id="phase-2" >}}

### L'auto auto-correction

Je veux que tu laisses les choses affichées de la sorte, mais desormais les points ne s'affichent plus. Il part de agir vers observer qui passe en surbrillance.
Ensuite il passe vers planifier via la ligne "auto-correctif", puis planifier en surbrillance, puis point vers agir, et vers observer.

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="3" id="phase-3" >}}

### Complément d'informations

Le point part de observer et remonte vers capter, puis vers planifier, puis vers agir

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="4" id="phase-4" >}}

### Condition de sortie

Le point part de observer et va vers resultat

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="5" id="phase-5" >}}

## la plateforme comme support d'execution

on fait apparaitre la plateforme sans le cylindre "standards" (ni la courroie, ni le cylindre qui fait tourner la roue à gauche). on affiche pas non plus la ligne api, mcp, a2a, le x-as-a-service et les capacités cloud. on affiche aussi que platform et pas "garantit cohérence fiabilité et confiance."

je voudrais une animation qui montre que llm tourne et que c'est lui qui active la courroie et qui fait touner la roue "système agentique". j'imagine des pointillés qui vont de droite à gauche au dessus et de gauche à droite en dessous et pareil le cercle système agentique en pointillé qui tourne en sens antihoraire

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="6" id="phase-6" >}}

## la plateforme comme fournisseur d'outils

On fait apparaitre les capacités cloud, les connaissances du SI, les capacités du SI ainsi que le X-as-a-service, le API, MCP et A2A et on fait une animation qui part de Agir vers x-as-a-service et vers les capacités

EDIT: on doit continuer de voir l'animation qui fait tourner
{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="7" id="phase-7" >}}

## Conception de solution

Je veux que tous les éléments soient cachés. ..
Et qu'on affiche que la problématique métier toute la boite jaune avec son contenu sauf "résultat" et "Livraison". Tu vas créer un point qui va de problématique vers conception pui de conception à specs et de specs à expression.

EDIT: enleve la fleche "enrichit la solution" et le fond jaune ainsi que la stream aligned team

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="7.2" id="phase-7.2" >}}

Je veux que tu ajoutes la boite capter avec le lien, mais au lieu d'écrire "capter", je veux que tu écrives "développement de la solution"
 
{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="8" id="phase-8" >}}

### Le SDLC Agentique

Tu vas alors afficher toute la partie "système agentique" (avec "implémente la solution") pour que l'on voit le système complet stream-aligned + système agentique avec un point qui va faire le parcours suivant:

Problématique -> conception -> specs -> expression -> capter -> planifier -> agir -> observer -> resultat 

EDIT: tu vas enlever le cadre jaune et le stream-aligned team

{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="8.1" id="phase-8.1" >}}
## Organisation

### L'équipe stream-aligned

Tu vas alors afficher toute la partie "système agentique" (avec "implémente la solution") pour que l'on voit le système complet stream-aligned + système agentique avec un point qui va faire le parcours suivant:

Problématique -> conception -> specs -> expression -> capter -> planifier -> agir -> observer -> resultat 


{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="8.2" id="phase-8.2" >}}

#### L'équipe stream-aligned constuit aussi des agents

Je veux que que la livraison clignote avec un affichage d'un icone robot dessus

{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="9" id="phase-9" >}}

## La besoin d'enabling

Je ne veux voir que tout le bloc jaune (stream-aligner) et le sytème agentique mais pas la plateforme 

EDIT: faire apparaitre l'équipe enabling avec le titre sans les liens
EDIT2: je veux aussi voir "pose le cadre conceptuel..."

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="10" id="phase-10" >}}

### les apports fondamentaux

Je veux voir apparaitre en plus le cadre mauve enabling team avec collaboration (mais sans "capitalisation ...").
On a aussi la fleche parametre et pas la fleche contexte technique

On va avoir un point qui va de enabling team vers le système agentique pour le paramétrage

EDIT: dans l'intersection entre enabling et stream-aligned en dessous de collaboration, tu peux écrire aide à la conception de la solution et à la ligne: mise en place des standards pour compléter le contexte

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="11" id="phase-11" >}}

### les apports fondamentaux
On fait apparaitre la fleche contexte technique avec une animation de enabling vers le systeme agentique

{{< /scrollytelling-step >}}

{{< scrollytelling-step phase="12" id="phase-12" >}}

On fait apparaitre capitalisation (skills agents etc...)

{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="13" id="phase-13" >}}

On fait réapparaitre la platefore (tous les éléments précédents). En plus on ajoute le cylindre standards mais sans la courroie.

EDIT: on ajoute l'animation qui fait tourner les roues comme dans la phase 5

{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="14" id="phase-14" >}}

### La standardisation

On ajoute la courroie et on affiche une animation qui montre que la standardisation et le LLM font tourner la roue du système agentique

EDIT: je veux que la ligne "mise en place des standards" descende jusque capitalisation et disapraisse.

{{< /scrollytelling-step >}}


{{< scrollytelling-step phase="15" id="phase-15" >}}

### La réussite de l'enabling team

On supprime tout ce qui concerne l'enabling team

{{< /scrollytelling-step >}}


{{< /scrollytelling >}}
