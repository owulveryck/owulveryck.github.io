# Infrastructure CloudFront

Distribution ID : `EJVG52EKHMUNU`

## Configuration actuelle

- **Origin** : S3 REST (`blog.owulveryck.info.s3.amazonaws.com`), ACL `public-read`
- **Compression** : Brotli + Gzip activée (`Compress: true`)
- **HTTP** : HTTP/3 (QUIC)
- **Cache** : legacy settings, TTL respecte les headers `Cache-Control` S3
- **Security headers** : Response Headers Policy `blog-security-headers` (`8f8054a8-14ac-452e-942a-f4e3c438807b`)
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
