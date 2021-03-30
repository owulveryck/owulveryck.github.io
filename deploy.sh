#!/bin/zsh

hugo --buildDrafts -d ./s3-drafts
hugo -d ./s3
#for f in $(ls s3/*.html s3/**/*.html)
#do
#    mv $f "${f%%.*}"
#done

s3cmd sync --acl-public --guess-mime-type s3/ s3://blog.owulveryck.info
