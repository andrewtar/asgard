#!/usr/bin/env bash

# Example:
# tools/save_image_to_registry.sh \
#   --source "k8s.gcr.io/ingress-nginx/controller:v1.2.0@sha256:d8196e3bc1e72547c5dec66d6556c0ff92a23f6d0919b206be170bc90d5f9185" \
#   --name ingress-nginx/controller:v1.2

set -e

echo "Running $0 $*..."

PARSED=$(getopt -o "" --longoptions source:,name: --name "$0" -- "$@")

eval set -- "$PARSED"
while true; do
  case "$1" in
    --source)
      source="$2"
      shift 2
      ;;
    --name)
      name="$2"
      shift 2
      ;;
    *) break ;;
  esac
done

base_repository="cr.yandex/crp1l4j9no209t82ra7l"

if [ ! "${source}" ]; then
  echo "Error: Source image shoud be specified."
  exit 1
fi

if [ ! "${name}" ]; then
  echo "Error: Destination image name shoud be specified."
  exit 1
fi

new_iamge_name="${base_repository}/${name}"

docker pull "${source}"
docker tag "${source}" "${new_iamge_name}"
docker push "${new_iamge_name}"

echo "Done. Image ${source} was saved as ${new_iamge_name}"
