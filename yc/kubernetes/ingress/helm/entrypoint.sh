#!/usr/bin/env sh

set -ex

# From list https://gitea.com/gitea/helm-chart/releases.
giteaHelmReleaseUrl=$1
if [ -z "${giteaHelmReleaseUrl}" ]
then
  echo "ERROR: Release archive URL is required!"
  exit 1
fi

fileName=$(basename ${giteaHelmReleaseUrl})
giteaHelm="/tmp/gitea/helm"
rm -rf "${giteaHelm}"
mkdir -p "${giteaHelm}"
giteaHelmArchive="${giteaHelm}/${fileName}"

wget "${giteaHelmReleaseUrl}" \
  --output-document="${giteaHelmArchive}" \
  --quiet

giteaHelmSources="${giteaHelm}/src"
mkdir -p "${giteaHelmSources}"  
tar -xvf "${giteaHelmArchive}" -C "${giteaHelmSources}"

cd "${giteaHelmSources}/ingress-nginx/"
helm repo add bitnami https://raw.githubusercontent.com/bitnami/charts/pre-2022/bitnami
helm dependency build
echo "Rendering"
helm template .
echo "Done"
