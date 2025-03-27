---
title: "Diffusion en continu de la reMarkable 2"
date: 2021-03-30T17:44:10+02:00
draft: false
summary: "Cet article est une traduction automatique. L'article original a été écrit en anglais. Cet article décrit le câblage de l'outil que j'ai créé pour diffuser le contenu de la remarkable 2 sur un ordinateur. Du système de fichiers proc à l'implémentation gRPC sur HTTP/2 via la génération de certificats."
tags: ["remarkable", "grpc", "protobuf", "linux", "go"]
---

> **Note**: Cet article est une traduction automatique. [L'article original a été écrit en anglais](/post/20210329-remarkable/).

Je suis l'heureux propriétaire d'une tablette [_reMarkable 2_](https://remarkable.com/). L'appareil est facile à utiliser dès sa sortie de la boîte. La seule chose qui me manque est un moyen approprié de diffuser le contenu sur mon ordinateur portable pour le diffuser lors d'une visioconférence.

Différentes solutions existent pour cela, mais je voulais quelque chose de facile à déployer avec très peu de dépendances et de configurations. De plus, je suis toujours à la recherche de projets à coder et de nouvelles choses à apprendre. Coder un outil pour répondre à mon besoin est un moyen parfait d'atteindre ces deux objectifs.

Cet article explique comment fonctionne l'outil [goMarkableStream](https://github.com/owulveryck/goMarkableStream).

Dans ce post, vous trouverez :

- des informations sur le système de fichiers `/proc/` sous Linux ;
- la génération de client et serveur gRPC à partir d'une définition protobuf ;
- Une paire de certificats intégrés pour l'authentification mutuelle.

## Obtenir une image de la tablette

La première chose à déterminer est comment obtenir une image de la reMarkable.

La remarkable est un appareil basé sur armv7 fonctionnant sous un système d'exploitation Linux. L'accès SSH est fourni, il est donc assez facile de se connecter en tant que `root` sur la tablette.

La méthode habituelle pour obtenir une image est d'interroger le [framebuffer](https://en.wikipedia.org/wiki/Framebuffer). Le noyau Linux expose un [périphérique de framebuffer](https://www.kernel.org/doc/Documentation/fb/framebuffer.txt) adressable via un nœud de périphérique (typiquement `/dev/fb0`). Ce périphérique vise à fournir une abstraction, _afin que le logiciel n'ait pas besoin de connaître quoi que ce soit sur les registres matériels de bas niveau._

Ma première tentative a échoué : l'interrogation du périphérique `/dev/fb0` ne fonctionne pas sur la reMarkable 2. Des personnes brillantes ont fait de l'ingénierie inverse et ont fourni une bonne explication sur ce [site web](https://remarkablewiki.com/tech/rm2_framebuffer). En résumé :

> La rm2 n'utilise pas le epdc intégré (Electronic Paper Display Controller) de l'imx7. Au lieu de cela, l'écran e-Ink est directement connecté au contrôleur LCD. Cela signifie que toutes les fonctions que le epdc ferait normalement sont maintenant effectuées par logiciel...

Cela signifie que le framebuffer n'est pas exposé dans `/dev/fb0` par le noyau mais par logiciel.

Pour obtenir une image, nous devons obtenir l'adresse de la portion de RAM contenant le bitmap de l'image de la tablette, et nous savons qu'elle n'est pas référencée par le noyau.

### L'adresse du framebuffer

Pour obtenir l'adresse globale du framebuffer dans la RAM, nous allons interroger un processus qui la connaît déjà. L'application principale de la remarkable gérant l'interface graphique s'appelle [`xochitl`](https://remarkablewiki.com/tech/xochitl). C'est un logiciel à code fermé ; il n'est donc pas possible de trouver ce que nous cherchons en modifiant le code.

_Note :_ Ce n'est pas tout à fait exact. Il est possible de hacker le processus, mais cela dépasse largement mes compétences. Voir le [remarkable2-framebuffer](https://github.com/ddvk/remarkable2-framebuffer) pour plus d'informations.

Le noyau Linux trace le mappage mémoire par processus et l'expose dans le pseudo-fichier `proc/[pid]/maps` (voir [man 5 procfs](https://man7.org/linux/man-pages/man5/procfs.5.html)).

En analysant les maps, il apparaît que le processus `xochitl` mappe virtuellement l'adresse du framebuffer vers le pseudo-périphérique.

```shell
grep -C1 '/dev/fb0' /proc/$(pidof xochitl)/maps
72086000-72886000 rw-p 00000000 00:00 0
72886000-74044000 rw-s a8100000 00:06 248        /dev/fb0
74044000-747d2000 rw-p 00000000 00:00 0
```

Le framebuffer global est donc situé à `0x74044000` dans la RAM. La RAM du processus `xochitl` est accessible via un appel à `/proc/[pid]/mem` (encore une fois, voir [man 5 procfs](https://man7.org/linux/man-pages/man5/procfs.5.html)).

Maintenant, combien d'octets devons-nous extraire ?

La résolution de la reMarkable 2 est de 1404x1872. Par conséquent, récupérons 2628288 octets :

```shell
reMarkable: ~/ echo $((0x74044000))
1946435584
reMarkable: ~/ dd if=/proc/$(pidof xochitl)/mem of=image.raw count=2628288 bs=1 skip=1946435584
2628288+0 records in
2628288+0 records out
reMarkable: ~/ ls -lrth image.raw
-rw-r--r--    1 root     root        2.5M Mar 31 07:43 image.raw
```

### Notre première capture d'écran

Récupérons le fichier `image.raw` et convertissons-le dans un format lisible avec imagemagick :

```shell
 convert -depth 8 -size 1872x1404+0 gray:image.raw image.png
 ```

 Nous pouvons alors afficher l'image qui pourrait ressembler à ceci :
 
 ![hello reMarkable](/assets/remarkable_hello.png)

 ## Construction d'une application

Maintenant que nous sommes capables de saisir une image, construisons une application pour saisir un flux en temps réel.

 ### Architecture globale et principe 

L'application fonctionne en mode client/serveur. Le serveur obtient les images brutes dans une boucle infinie et les sert sur le réseau.
C'est ensuite la responsabilité du client de récupérer les images brutes du réseau et de les encoder dans un flux vidéo.

Une implémentation triviale consisterait à ouvrir une connexion réseau au niveau 4 et à utiliser le protocole TCP comme support pour le flux d'octets.
Néanmoins, cela impliquerait un travail pour mettre en place des délimiteurs entre chaque trame et gérer les mauvais messages.

Par conséquent, c'est une bonne idée d'encapsuler chaque image dans un message et de s'appuyer sur les capacités d'un framework pour faire un codage/décodage approprié.

Jusqu'à présent, l'option la plus large est d'utiliser les protocol buffers car ils utiliseront un mécanisme de typage décent tout en restant compacts et faciles à utiliser.

Le message représente une image et est défini comme suit :

```proto
message image {
    int64 width = 1;
    int64 height = 2;
    bytes image_data = 4;
}
```

Traiter le flux des messages pour gérer une image une par une fait partie d'un protocole de niveau 7. Au lieu d'écrire le nôtre, continuons à travailler avec protobuf en utilisant [gRPC](https://grpc.io/).
gRPC est un framework RPC universel haute performance et open-source qui fonctionne sur HTTP/2. La surcharge réseau est donc faible, et la communication entre le client et le serveur reste efficace.

Notre service de streaming exposera une fonction `GetImage` qui saisira l'image de la mémoire et l'enverra sur le réseau :

```proto
message Input {}

service Stream {
  rpc GetImage(Input) returns (image) {}
}
```

### Implémentation

L'implémentation du client et du serveur est faite en Go.

L'outil `protoc` génère le squelette du service de streaming :

```shell
protoc --gofast_out=plugins=grpc:.  defs.proto3
```

Parmi quelques utilitaires pour gérer la sérialisation et la désérialisation du message protobuf (voir la doc [Image](https://pkg.go.dev/github.com/owulveryck/goMarkableStream@v0.2.1/stream#Image) pour plus d'informations), le framework gRPC expose certains

Le [`StreamServer`](https://pkg.go.dev/github.com/owulveryck/goMarkableStream@v0.2.1/stream#StreamServer) est une interface. C'est maintenant notre responsabilité de créer une structure qui remplit l'interface et qui implémente réellement le mécanisme `GetImage` (obtenir l'image de la mémoire comme exposé précédemment)

```go
type StreamServer interface {
	GetImage(context.Context, *Input) (*Image, error)
}
```

Notre serveur est une structure de base gérant quelques éléments :

```go
// Server implementation
type Server struct {
	imagePool   sync.Pool
	r           io.ReaderAt
	pointerAddr int64
	runnable    chan struct{}
}
```

Le champ `r` est un pointeur vers le fichier `/proc/[pid]/mem` à partir duquel nous lirons les données. `pointerAddr` est l'emplacement du framebuffer dans ce fichier (0x74044000) dans notre exemple, et `runnable` est un canal utilisé pour gérer les requêtes et éviter de surcharger le CPU de la reMarkable (TL;DR : deux appels consécutifs à `GetImage` devront attendre pour pouvoir consommer `runnable` et une goroutine met un événement toutes les x millisecondes dans la file d'attente runnable).

Fondamentalement, l'implémentation de `GetImage` est triviale :

```go
// GetImage input is nil
func (s *Server) GetImage(ctx context.Context, in *Input) (*Image, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-s.runnable:
		img := s.imagePool.Get().(*Image)
		_, err := s.r.ReadAt(img.ImageData, s.pointerAddr)
		if err != nil {
			s.imagePool.Put(img)
			return nil, err
		}
		return img, nil
	}
}
```

La magie consiste simplement à lire les octets, à les mettre dans une image et à les renvoyer à l'appelant.
Exposer le service consiste simplement à instancier les objets et à utiliser les outils construits par le framework gRPC :

```go
ln, _ := net.Listen("tcp", ":2000") // open a listener on TCP on ":2000"
s := stream.NewServer(file, addr) // create the stram object
s.Start() // start the gorouting that feeds the `runnable` channel every x ms
grpcServer := grpc.NewServer(grpc.Creds(grpcCreds)) // create the gRPC server
stream.RegisterStreamServer(grpcServer, s) // register our stream object so it is used by our server
grpcServer.Serve(ln); err != nil { // make the server listen on a TCP connection
```

Le client compose simplement le serveur et appelle la procédure distante `GetImage` dans une boucle sans fin :

```go
conn, err := grpc.Dial("localhost:2000") // Dial the server
client := stream.NewStreamClient(conn)

var img image.Gray
for err == nil {
    response, err := client.GetImage(context.Background(), &stream.Input{})
```

Ensuite, il encode la `response` dans un fichier JPEG et l'ajoute à un flux MJPEG.

```go
var img image.Gray
var b bytes.Buffer
img.Pix = response.ImageData
img.Stride = int(response.Width)
img.Rect = image.Rect(0, 0, int(response.Width), int(response.Height))
jpeg.Encode(&b, &img, nil)
mjpegStream.Update(b.Bytes())
```

La création et l'exposition du flux MJPEG ne sont pas détaillées dans ce post car elles sont légèrement hors contexte. Veuillez consulter le code si vous voulez plus d'informations.

### Sécurité

Même si HTTP/2 ne nécessite pas de chiffrement (voir [ici](https://http2.github.io/faq/#does-http2-require-encryption)), de nombreuses implémentations ne prennent en charge le protocole que s'il est utilisé sur une connexion chiffrée.
L'implémentation Go de gRPC nécessite par défaut un canal de chiffrement (qui peut être contourné avec l'utilisation d'une méthode `Insecure`, mais nous savons tous que ce n'est pas une bonne façon de procéder ;)).

C'est donc une bonne pratique d'implémenter ce mécanisme de sécurité qui évitera l'espionnage des images du wifi si vous utilisez l'outil sur un réseau non fiable.

Comme je ne veux rien de difficile à maintenir, je génère un certificat auto-signé que j'intègre à la fois sur le client et le serveur avec la nouvelle commande `embed` du langage Go.

J'implémente également un mécanisme d'authentification mutuelle. Par conséquent, seul un client connu peut se connecter au serveur.
Le certificat est généré par build (via un ensemble de commandes `go:generate`). Par conséquent, si vous souhaitez améliorer la sécurité, c'est votre responsabilité de générer de nouveaux binaires et de les stocker dans un endroit sûr, quelque part sur votre ordinateur (car ils contiennent le certificat).
Je conviens que ce n'est pas l'option la plus sécurisée, mais elle est suffisamment bonne pour la plupart des cas d'utilisation.

### Génération du certificat

Le certificat est généré en code Go pur :

- Un package interne est en charge de la sorcellerie des certificats (voir la [documentation du certificat](https://pkg.go.dev/github.com/owulveryck/goMarkableStream@v0.2.1/internal/certificate)).
- Une CLI simple génère le fichier (voir [le code](https://github.com/owulveryck/goMarkableStream/blob/v0.2.1/certs/cmd/main.go)).
- Un package `cert` (voir [la documentation ici](https://pkg.go.dev/github.com/owulveryck/goMarkableStream@v0.2.1/certs) expose une seule fonction `GetCertificateWrapper()` renvoyant une configuration prête à l'emploi basée sur le certificat intégré (`*certificate.CertConfigCarrier`).

Intégrer TLS dans le serveur gRPC est simple :

1. Pour le serveur :

```go
cert, err := certs.GetCertificateWrapper() // Get the certificate configuration with the embeded certificate
grpcCreds := &callInfoAuthenticator{credentials.NewTLS(cert.ServerTLSConf)} // callInfoAuthenticator is fulfiling the interface https://pkg.go.dev/google.golang.org/grpc@v1.36.1/credentials#TransportCredentials and do the validation of the cerficiate of the client
grpcServer := grpc.NewServer(grpc.Creds(grpcCreds)) // creates the server with the validation mechanism
```

2. Pour le client :

```go
cert, err := certs.GetCertificateWrapper()
grpcCreds := credentials.NewTLS(cert.ClientTLSConf)
// Create a connection with the TLS credentials
conn, err := grpc.Dial(c.ServerAddr, grpc.WithTransportCredentials(grpcCreds), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
//...
```

## C'est tout, les amis !

L'outil semble fonctionner comme prévu pour la plupart des utilisateurs. Au moins, il est assez bon pour moi. Je ne prévois pas d'ajouter de fonctionnalités fantaisistes. N'hésitez pas à l'essayer si vous possédez une tablette :

[https://github.com/owulveryck/goMarkableStream](https://github.com/owulveryck/goMarkableStream)

Le dépôt contient également un fichier `goreleaser` si vous souhaitez créer votre propre version avec vos propres certificats.

Voici une vidéo du produit final :
{{< youtube c4-hJ6xRzg4 >}}