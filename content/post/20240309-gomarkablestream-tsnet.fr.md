---
title: "Après le BYOD, le BYOC (Bring Your Own Cloud): un voyage de la Maison au Monde"
date: 2024-03-09T12:15:33+01:00
lastmod: 2024-03-09T12:15:33+01:00
draft: false
keywords: []
summary: Cet article est une traduction automatique. L'article original a été écrit en anglais. Découvrez comment j'ai transformé ma tablette reMarkable en tableau blanc portable 📒✨, accessible de n'importe où via un VPN WireGuard sécurisé (tailscale) et une configuration de proxy inverse basée sur le cloud.
  
  Du confort du télétravail au monde dynamique de la mobilité, découvrez la technologie derrière la solution.
tags: []
categories: []
author: ""

# Uncomment to pin article to front page
# weight: 1
# You can also close(false) or open(true) something for this content.
# P.S.
comment: false
toc: true
autoCollapseToc: false
# You can also define another contentCopyright.
contentCopyright: false
reward: false
mathjax: false

# Uncomment to add to the homepage's dropdown menu; weight = order of article
# menu:
#   main:
#     parent: "docs"
#     weight: 1
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/20240309-gomarkablestream-tsnet/).

## Contexte

À l'ère du télétravail, j'ai développé goMarkableStream, un outil conçu pour diffuser de manière transparente le contenu de ma tablette reMarkable pendant les appels vidéo.
L'objectif était de remplacer le griffonnage sur tableau blanc physique lors des réunions à distance.

L'outil est passé d'une preuve de concept à une partie de ma boîte à outils quotidienne. L'idée derrière l'outil est :

Un service fonctionne sur l'appareil reMarkable et capture l'image.
Il expose un service qui sert l'image via HTTP(s) avec une implémentation personnalisée.
Ensuite, un moteur de rendu est encodé dans le navigateur en WebGL/JS pour afficher le contenu de l'écran.
Pendant un appel vidéo, je peux partager un onglet du navigateur et donc partager ce que j'écris avec l'audience.

En effet, la solution a apporté de la valeur pour la collaboration à distance et le partage d'idées en temps réel.

J'ai partagé les détails de **ce voyage**, soulignant comment peu à peu, il **a comblé le fossé entre les réunions physiques et virtuelles**.

Cependant, alors que la phase de télétravail à temps plein a diminué et que nous sommes passés à un mode de vie plus mobile, la solution qui expose un service sur un réseau local atteint ses limites.

## Problème

Avec le travail hybride, et même avec une situation de travail depuis n'importe où (maison, bureau, sites clients), j'ai rencontré des situations où je partageais ma connexion mobile avec ma tablette et j'étais simplement incapable de diffuser le contenu en raison de limitations.

J'ai donc besoin d'un **changement de paradigme** : le service doit passer **d'un outil interne** à **un service** auquel je dois **accéder de n'importe où**.
La diffusion sur Internet est la voie à suivre.

Comme je ne peux pas simplement exposer le service de streaming hébergé sur ma tablette à Internet pour des raisons évidentes,
j'ai besoin de m'appuyer sur un tiers pour gérer la connexion de l'extérieur vers la tablette.

Cet article décrit le parcours pour y parvenir, d'un simple proxy inverse via NGrok à une solution VPN basée sur WireGuard.

Je vais d'abord exposer la solution basée sur un proxy inverse propulsé par [NGrok](https://ngrok.com/).
Puis j'expliquerai les limitations qui m'ont conduit à la solution d'accéder au service via un VPN propulsé par Tailscale.
Cette partie donnera des indices sur le mécanisme wireguard et exposera les éléments de base de l'infrastructure en place pour exposer le service de streaming.

Avant la pandémie, nous utilisions des VPN pour nous connecter au bureau depuis la maison... Maintenant, j'ai changé de paradigme pour me connecter à la maison depuis le bureau.
Je suppose que c'est la suite de l'évolution du "bring your own device" (BYOD).

## Première solution : NGrok

Comme [je l'ai écrit sur mon blog il y a quelques mois](https://blog.owulveryck.info/2023/10/10/rethinking-presentations-beyond-static-slides.html), j'utilise ma tablette comme support pour les présentations.
Cela fonctionne parfaitement sur mon propre réseau, mais j'ai rencontré des problèmes lorsque je me suis déplacé sur un site avec des limitations.
J'ai pensé que je pourrais toujours apporter mon propre ordinateur portable avec moi, mais ce n'est pas toujours le cas. J'avais donc besoin d'un moyen d'exposer le service de streaming sur Internet et de donner l'adresse aux personnes chargées de présenter le contenu.

La première étape facile que j'ai trouvée comme solution était d'intégrer le service NGrok dans mon outil.
En fait, la promesse de NGrok est :

> Se connecter à des réseaux externes de manière cohérente, sécurisée et répétable sans nécessiter de changements dans les configurations réseau.
> - Connectivité Bring Your Own Cloud (BYOC)
> - Connectivité IoT

L'implémentation était assez facile à intégrer dans l'outil.

_Note :_ J'intègre cela dans l'outil car je veux que l'application soit autonome, soit moins intrusive dans le système natif de la tablette, et donc facile à installer et à exécuter.

En fait, comme il existe un [SDK Go pour NGrok](https://ngrok.com/docs/using-ngrok-with/go/) et que mon outil est écrit en Go, je n'ai qu'à importer et initier le service.

Fondamentalement, NGrok implémente un [`Listener`](https://pkg.go.dev/net#Listener), et tout ce que j'ai besoin de faire est de remplacer le listener de base du service HTTP pour utiliser ce listener à la place. La magie se produit sous le capot (connexion au service NGrok, etc.).

Voici une fonction auxiliaire pour initialiser le listener basé sur une structure de configuration :

```go 
func setupListener(ctx context.Context, c *configuration) (net.Listener, error) {
        switch c.BindAddr {
        case "ngrok":
                l, err := ngrok.Listen(ctx,
                        config.HTTPEndpoint(),
                        ngrok.WithAuthtokenFromEnv(),
                )
                c.BindAddr = l.Addr().String()
                c.TLS = false
                return l, err
        default:
                return net.Listen("tcp", c.BindAddr)
        }
}
```

Et voici son utilisation dans la boucle principale (`handler` a été configuré auparavant) :

```go 
l, err := setupListener(context.Background(), &c)
// ...
log.Fatal(http.Serve(l, handler))
```

Lorsque je lance l'outil, avec les variables d'environnement correctes, il se connecte au service NGrok et affiche l'URL externe pour s'y connecter.
Et voilà : ça marche !

Cependant, il y a des contraintes et des limitations :

- Tout d'abord, avec la version gratuite de NGrok, le réseau est limité. Je ne pourrai pas utiliser mon outil tout le mois, mais je pourrais vivre avec.
- Le deuxième problème est que je ne peux pas configurer le DNS du point de terminaison sur la version gratuite. Et chaque fois qu'il redémarre, l'URL du point de terminaison change. C'est ennuyeux.

Tous ces problèmes auraient été résolus en payant pour le service NGrok, mais c'est beaucoup trop cher pour mes besoins et en effet, n'aurait pas résolu le dernier problème :

Mais le plus gros problème est que **la solution ne gère pas bien l'itinérance** (changement de réseaux) et **les longues pauses** (lorsque la tablette dort pendant longtemps). Cela a rendu la solution peu fiable.

J'ai donc cherché une autre solution.

## Solution suivante : un VPN ?

Une solution potentielle consiste à rendre le service accessible sur Internet en utilisant un nom cohérent. Cependant, plusieurs défis se posent :

- Les appareils se connectent souvent à un réseau privé et accèdent à Internet via une passerelle.
- Exposer directement le service à Internet pose des risques de sécurité.

Une solution à mon problème implique une passerelle qui dirige le trafic externe vers le service spécifié sur mon appareil au sein du réseau privé.
Mais, pour s'adapter à l'itinérance, la passerelle doit soit :

- Être "intelligente" et suivre l'adresse de l'appareil, soit
- S'assurer que l'adresse de l'appareil au sein du réseau reste statique.

Une **passerelle intelligente** crée une **forte dépendance** vis-à-vis du service et **nécessite une couche de persistance** pour surveiller l'emplacement de l'appareil, une approche que je préfère éviter.

Alternativement, exploiter l'infrastructure pour attribuer une adresse statique à l'appareil est facilement réalisable en établissant un VPN.
Ce VPN étendra le réseau privé sur Internet, maintenant constante l'adresse IP de l'appareil, quelle que soit la topologie de connexion.

**Dans les protocoles VPN conventionnels** comme IPsec ou OpenVPN, la **connexion** du VPN **dépend généralement de l'adresse IP de l'appareil qui se connecte**.
Si l'adresse IP de l'appareil change (par exemple, lors du passage d'un réseau à un autre), une connexion VPN typique serait interrompue, nécessitant le rétablissement de la connexion sous la nouvelle adresse IP.
Cette procédure peut causer des retards et des perturbations dans la connectivité.

Heureusement, une alternative moderne aux VPN traditionnels existe : Wireguard !

#### L'approche de WireGuard

WireGuard adopte une approche différente des VPN traditionnels qui prend intrinsèquement en charge l'itinérance transparente :

- **Identification de la connexion :** WireGuard identifie les connexions non pas par les adresses IP source ou destination mais par l'identité cryptographique des pairs (c'est-à-dire leurs clés publiques).
Cela signifie que tant que l'identité cryptographique reste la même, WireGuard ne se soucie pas si l'adresse IP réelle d'un appareil change.
- **Persistance de session :** Lorsqu'un client WireGuard se déplace vers un réseau différent et obtient une nouvelle adresse IP, il envoie simplement des paquets authentifiés depuis sa nouvelle IP au serveur WireGuard (ou pair).
Le serveur reconnaît le client par sa clé publique et poursuit la session sans interruption.
Le serveur met ensuite automatiquement à jour sa table de routage interne avec la nouvelle adresse IP du client, maintenant le tunnel chiffré sans avoir besoin de rétablir la connexion.
- **Réponse rapide :** Ce mécanisme permet un basculement presque instantané entre les réseaux.
Les utilisateurs ne remarquent généralement aucune perturbation dans leur connexion VPN alors qu'ils se déplacent à travers différents environnements réseau, ce qui rend WireGuard particulièrement adapté aux appareils mobiles qui changent fréquemment d'environnements réseau.

WireGuard est entièrement implémenté dans [Tailscale](https://tailscale.com/).

Tailscale implémente un _réseau défini par logiciel (SDN)_.
À sa base, il établit un périphérique réseau virtuel au niveau du noyau du système d'exploitation, fournissant ainsi un service réseau accessible à toutes les applications.

### Défis et solutions dans l'intégration

Tailscale est développé en Go, tirant parti du support du langage pour les applications autonomes.
Cette approche signifie qu'un seul binaire peut englober toutes les fonctionnalités de Tailscale.
La complétude de Turing de Go facilite la facilité de compilation croisée et de portage du code à travers différentes architectures.

Vous exécutez simplement `./tailscale` qui gère le processus et crée ou rejoint un réseau IP appelé "_tailnet_"

Par conséquent, il existe une version de Tailscale compatible avec l'appareil reMarkable, qui est un système basé sur Linux fonctionnant sur un processeur ARM v7.

Malheureusement, le noyau Linux de reMarkable ne prend pas en charge le pilote de périphérique [tun/tap](https://docs.kernel.org/networking/tuntap.html), et il est donc impossible d'exécuter tailscale tel quel.

_Note_ : il a été signalé sur Reddit qu'exécuter Tailscale sur le reMarkable est en fait possible, comme expliqué [ici](https://remarkable.guide/tech/tailscale.html).

Cependant, comme Tailscale fonctionne comme un SDN, il existe une méthode alternative pour se connecter au service sans dépendre du support du noyau, purement en espace utilisateur : [_tsnet_](https://tailscale.com/kb/1244/tsnet).

## Introduction à la bibliothèque tsnet

> tsnet est une bibliothèque qui vous permet d'intégrer Tailscale dans un programme Go.
Cela utilise une pile réseau TCP/IP en espace utilisateur et établit des connexions directes à vos nœuds sur votre tailnet comme le ferait n'importe quelle autre machine sur votre tailnet.
Combiné à d'autres fonctionnalités de Tailscale, cela vous permet de créer des façons nouvelles et intéressantes d'utiliser des ordinateurs auxquelles vous n'auriez jamais pensé auparavant.

### Implémentation de la solution

Comme NGrok, tsnet implémente un listener, nous permettant de modifier la fonction que nous avons précédemment définie pour accommoder le scénario "tailscale".

Il y a une astuce intéressante.
Lors de la première connexion, pour enregistrer le service sur le tailnet, le framework affiche une URL pour l'authentification via Single Sign-On (SSO).
Si nous désactivons la journalisation, cette information cruciale n'apparaît plus.
Bien qu'il existe plusieurs façons de gérer cette situation, la solution la plus simple est d'initier le service en "mode développement" pour la première utilisation (en activant un drapeau spécifique),
et ensuite de supprimer la journalisation lorsque ce drapeau est désactivé (par exemple, lors du démarrage en tant que service).

Voici l'implémentation proposée :

```go
func setupListener(ctx context.Context, c *configuration) (net.Listener, error) {
        switch c.BindAddr {
        case "tailscale":
                srv := new(tsnet.Server)
                srv.Hostname = "gomarkablestream"
                // Disable logs when not in devmode
                if !c.DevMode {
                        srv.Logf = func(string, ...any) {}
                }
                return srv.Listen("tcp", ":2001")
        case "ngrok":
                l, err := ngrok.Listen(ctx,
                        config.HTTPEndpoint(),
                        ngrok.WithAuthtokenFromEnv(),
                )
                c.BindAddr = l.Addr().String()
                c.TLS = false
                return l, err
        default:
                return net.Listen("tcp", s.BindAddr)
        }
}
```

Lorsque le service démarre, il expose le service et apparaît sur la console tailscale :

![Panneau d'administration de tailscale avec une liste de machines connectées, et une mise en évidence sur le service gomarkablestream](/assets/tsnet-gomarkablestream.png)

Le service est alors accessible via un appel http à `100.81.233.46` (dans l'exemple).

## Le reste de l'infrastructure

Maintenant que le service est exposé dans le VPN, j'ai besoin de mettre en place une passerelle pour y accéder depuis un autre réseau et éventuellement depuis Internet.

J'utiliserai [`Caddy`](https://caddyserver.com/) comme proxy inverse sur un nœud de mon tailnet. Ce nœud aura à la fois une connexion au tailnet et une connexion au réseau cible (celui où j'ai besoin d'obtenir le flux).

### Caddy comme proxy inverse

Le service Caddy fonctionnera sur une instance EC2 sur Internet, avec Tailscale installé pour s'assurer que la machine rejoigne mon tailnet.
J'attribuerai ensuite un nom DNS à l'instance EC2 (pour cet exemple, utilisons myremarkable.chezmoi.com).

Cet exemple de configuration Caddy (Caddyfile) démarrera le service, obtiendra automatiquement un certificat de Let's Encrypt et configurera l'authentification de base.
Une fois authentifié, il acheminera le trafic vers l'appareil remarkable.

```Caddyfile
{
        admin off
}
// This is the external name of the node
myremarkable.chezmoi.com {
        reverse_proxy gomarkablestream:2001

        # Basic authentication
        basicauth /* {
                user $#ENCRYPTEDPASSWORD
        }
}
```
Cette configuration garantit que l'accès à https://myremarkable.chezmoi.com depuis n'importe où sur Internet affichera de manière sécurisée le contenu de ma tablette,
à condition que la tablette soit connectée à Internet.
Le service s'accommode de l'itinérance ; ainsi, où que je sois, je peux connecter ma tablette à Internet (par exemple, via mon téléphone) et simplement accéder à l'URL pour me connecter de manière transparente.

### Conclusion

Cette solution est fonctionnelle mais peut être améliorée dans plusieurs dimensions.

Premièrement, faire fonctionner une machine en continu dans le cloud n'est ni rentable ni écologique. Une alternative plus efficace serait d'adopter une approche sans serveur, en encapsulant la passerelle dans un conteneur Docker qui est lancé "à la demande".

Concernant la sécurité, l'approche actuelle implique de faire confiance aux appareils au sein du tailnet. Une stratégie préférable impliquerait de passer à une [architecture zero trust](https://cyber.gouv.fr/publications/le-modele-zero-trust), qui ne fait pas intrinsèquement confiance à aucune entité à l'intérieur ou à l'extérieur du réseau.

Je prévois de consulter mon collègue [François](https://www.linkedin.com/in/%F0%9F%90%B3-fran%C3%A7ois-hamy-00433411a/) pour obtenir de l'aide dans cette direction.

Côté fonctionnalités, il y a un intérêt à améliorer l'"auth basique" de la passerelle vers un mécanisme plus robuste et à introduire la capacité de générer un accès temporaire. Cela permettrait d'accorder à des tiers la permission de visualiser temporairement le flux.

#### Mot final

Comme toujours, je trouve de la joie à affronter les nouvelles contraintes présentées par notre monde en évolution. Ces défis ne posent pas seulement des problèmes, mais me poussent également à embrasser mon rôle d'ingénieur : identifier et exécuter des solutions. Ce processus représente le summum de l'apprentissage.