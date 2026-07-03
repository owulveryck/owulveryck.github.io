# Infrastructure CloudFront

Distribution ID : `EJVG52EKHMUNU`

## Configuration actuelle

- **Origin** : S3 REST (`blog.owulveryck.info.s3.amazonaws.com`), ACL `public-read`
- **Compression** : Brotli + Gzip activée (`Compress: true`)
- **HTTP** : HTTP/3 (QUIC)
- **Cache** : AWS managed Cache Policy `CachingOptimized` (`658327ea-f89d-4fab-a63d-7e88639e58f6`), TTL respecte les headers `Cache-Control` S3
- **Security headers** : AWS managed Response Headers Policy `SecurityHeadersPolicy` (`67f7725c-6f97-4210-82d7-5512b31e9d03`) — HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, X-XSS-Protection
- **Logging** : activé, bucket `logs.blog.owulveryck.info`, préfixe `cf-logs/`, rétention 90 jours
- **URL rewrite** : CloudFront Function `blog-url-rewrite` sur `viewer-request`
- **Price Class** : `PriceClass_All`

## URL Rewrite (`url-rewrite.js`)

CloudFront Function qui résout le problème des URLs sans `index.html` sur S3.
S3 en mode REST API ne sert pas automatiquement `index.html`
pour les chemins se terminant par `/`.

La fonction réécrit les URIs sur l'événement `viewer-request` :
- `/fr/` → `/fr/index.html`
- `/fr` → `/fr/index.html`
- `/post.html` → inchangé

### Déploiement (one-time)

```bash
aws cloudfront create-function \
  --name blog-url-rewrite \
  --function-config '{"Comment":"Rewrite URIs for S3 index document resolution","Runtime":"cloudfront-js-2.0"}' \
  --function-code fileb://cloudfront-functions/url-rewrite.js

aws cloudfront publish-function \
  --name blog-url-rewrite \
  --if-match <ETAG>
```

### Mise à jour

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

## Logging (access logs)

Bucket S3 dédié en `us-east-1`, lifecycle policy de 90 jours.

### Setup (one-time)

```bash
# 1. Créer le bucket (us-east-1, région par défaut)
aws s3api create-bucket --bucket logs.blog.owulveryck.info

# 2. Activer les ACL (CloudFront écrit les logs via le compte awslogsdelivery)
aws s3api put-bucket-ownership-controls \
  --bucket logs.blog.owulveryck.info \
  --ownership-controls '{"Rules":[{"ObjectOwnership":"BucketOwnerPreferred"}]}'

# 3. Activer le logging sur la distribution
# Récupérer l'ETag et extraire uniquement DistributionConfig (sans l'enveloppe)
ETAG=$(aws cloudfront get-distribution-config --id EJVG52EKHMUNU \
  | tee /tmp/dist-full.json | jq -r '.ETag')
jq '.DistributionConfig' /tmp/dist-full.json > /tmp/dist-config.json

# Modifier le bloc "Logging" dans /tmp/dist-config.json :
#   "Logging": {
#     "Enabled": true,
#     "IncludeCookies": false,
#     "Bucket": "logs.blog.owulveryck.info.s3.amazonaws.com",
#     "Prefix": "cf-logs/"
#   }
# Ou directement avec jq :
jq '.Logging = {
  "Enabled": true,
  "IncludeCookies": false,
  "Bucket": "logs.blog.owulveryck.info.s3.amazonaws.com",
  "Prefix": "cf-logs/"
}' /tmp/dist-config.json > /tmp/dist-config-updated.json

aws cloudfront update-distribution \
  --id EJVG52EKHMUNU \
  --if-match "$ETAG" \
  --distribution-config file:///tmp/dist-config-updated.json

# 4. Lifecycle policy : suppression automatique après 90 jours
aws s3api put-bucket-lifecycle-configuration \
  --bucket logs.blog.owulveryck.info \
  --lifecycle-configuration '{
    "Rules": [{
      "ID": "expire-old-logs",
      "Status": "Enabled",
      "Filter": {"Prefix": "cf-logs/"},
      "Expiration": {"Days": 90}
    }]
  }'
```

### Vérification

```bash
# Bucket existe
aws s3 ls s3://logs.blog.owulveryck.info/

# Logging activé
aws cloudfront get-distribution --id EJVG52EKHMUNU \
  | jq '.Distribution.DistributionConfig.Logging'

# Logs présents (attendre quelques minutes après activation)
aws s3 ls s3://logs.blog.owulveryck.info/cf-logs/
```

## Stratégie de cache (définie dans `.github/workflows/publish.yml`)

| Contenu | TTL | Header |
|---|---|---|
| CSS/JS/fonts fingerprinted | 1 an | `public, max-age=31536000, immutable` |
| Images (`assets/*`, `content-images/*`) | 30 jours | `public, max-age=2592000` |
| HTML, XML, tout le reste | 1 jour | `public, max-age=86400` |

L'invalidation CloudFront (`/*`) est lancée à chaque déploiement pour garantir la fraîcheur du contenu HTML.

## Prochaine étape : Origin Access Control (OAC)

Migrer de `--acl public-read` vers OAC pour bloquer l'accès S3 direct :

```bash
# 1. Créer l'OAC
aws cloudfront create-origin-access-control \
  --origin-access-control-config '{
    "Name": "blog-owulveryck-oac",
    "Description": "OAC for blog.owulveryck.info S3 bucket",
    "SigningProtocol": "sigv4",
    "SigningBehavior": "always",
    "OriginAccessControlOriginType": "s3"
  }'

# 2. Attacher à la distribution (dans DistributionConfig)
#    Origins.Items[0].OriginAccessControlId = "<OAC_ID>"

# 3. Bucket policy autorisant CloudFront
aws s3api put-bucket-policy --bucket blog.owulveryck.info --policy '{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {"Service": "cloudfront.amazonaws.com"},
    "Action": "s3:GetObject",
    "Resource": "arn:aws:s3:::blog.owulveryck.info/*",
    "Condition": {
      "StringEquals": {
        "AWS:SourceArn": "arn:aws:cloudfront::566930824077:distribution/EJVG52EKHMUNU"
      }
    }
  }]
}'

# 4. Vérifier que le site fonctionne via CloudFront

# 5. Bloquer l'accès public S3
aws s3api put-public-access-block --bucket blog.owulveryck.info \
  --public-access-block-configuration \
  'BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true'

# 6. Retirer --acl public-read du workflow GitHub Actions
```
