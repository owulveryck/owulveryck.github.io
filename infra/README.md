# Infrastructure - CloudFront Function

## URL Rewrite (`url-rewrite.js`)

CloudFront Function qui résout le problème des URLs sans `index.html` sur S3.
S3 en mode REST API (origin CloudFront) ne sert pas automatiquement `index.html`
pour les chemins se terminant par `/`.

La fonction réécrit les URIs sur l'événement `viewer-request` :
- `/fr/` → `/fr/index.html`
- `/fr` → `/fr/index.html`
- `/post.html` → inchangé

## Déploiement (one-time)

```bash
# 1. Créer la fonction
aws cloudfront create-function \
  --name blog-url-rewrite \
  --function-config '{"Comment":"Rewrite URIs for S3 index document resolution","Runtime":"cloudfront-js-2.0"}' \
  --function-code fileb://cloudfront-functions/url-rewrite.js

# 2. Publier (passer de DEVELOPMENT à LIVE)
# Utiliser l'ETag retourné par create-function
aws cloudfront publish-function \
  --name blog-url-rewrite \
  --if-match <ETAG>

# 3. Associer à la distribution CloudFront
# Récupérer la config actuelle :
aws cloudfront get-distribution-config --id <DISTRIBUTION_ID> > dist-config.json

# Ajouter FunctionAssociations dans DefaultCacheBehavior du JSON :
#   "FunctionAssociations": {
#     "Quantity": 1,
#     "Items": [{
#       "FunctionARN": "<ARN de l'étape 1>",
#       "EventType": "viewer-request"
#     }]
#   }

# Appliquer :
aws cloudfront update-distribution \
  --id <DISTRIBUTION_ID> \
  --distribution-config file://dist-config-modified.json \
  --if-match <ETAG>
```

## Mise à jour

Si le code de la fonction change :

```bash
aws cloudfront update-function \
  --name blog-url-rewrite \
  --function-config '{"Comment":"Rewrite URIs for S3 index document resolution","Runtime":"cloudfront-js-2.0"}' \
  --function-code fileb://cloudfront-functions/url-rewrite.js \
  --if-match <CURRENT_ETAG>

aws cloudfront publish-function \
  --name blog-url-rewrite \
  --if-match <ETAG>
```
