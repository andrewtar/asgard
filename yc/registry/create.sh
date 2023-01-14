#!/usr/bin/env bash

name='registry'

registries=($(yc container registry list --format json | jq -r '.[].name'))
if [[ " ${registries[*]} " =~ " ${name} " ]]; then
  echo "The docker registry ${name} already exists. Skipped."
else
  yc container registry create \
    --name ${name}
    --description "Storage for docker images."
  echo "The docker registry ${name} created. OK."
fi
