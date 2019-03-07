#!/usr/bin/env bash

set -e
set -u

echo ${BUILDPACKS}

mkdir /var/lib/buildpacks
echo "${BUILDPACKS}" | jq '.[] | {"name": .name, "key": .key, "url": .url}' | jq -s .  > /var/lib/buildpacks/config.json

echo "${BUILDPACKS}" | jq -c '.[]' | while read row; do
  name=$(echo "${row}" | jq -r '.name')
  url=$(echo "${row}" | jq -r '.url')
  curl --cacert /etc/config/certs/internal-ca-cert -fsSLo /tmp/buildpack.zip "${url}"
  unzip -qq /tmp/buildpack.zip -d "/var/lib/buildpacks/$(echo -n "${name}" | md5sum | awk '{ print $1 }')"
  rm /tmp/buildpack.zip
done

/packs/recipe
