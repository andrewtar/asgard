#!/usr/bin/env bash

set -e

account_name="drone"
registry_name="registry"

accounts=($(yc iam service-account list --format json | jq -r '.[].name'))
if [[ " ${accounts[*]} " =~ " ${account_name} " ]]; then
  echo "The account ${account_name} already exists. Skipped."
else
  yc iam service-account create \
    --name ${account_name} \
    --description "Service account for Drone CI"
  echo "The account ${account_name} created. OK."
fi

roles=(
  container-registry.images.puller
  container-registry.images.pusher
)
account_id=$(yc iam service-account get --name ${account_name} --format json | jq .id -r)
for role in "${roles[@]}"
do
  yc container registry add-access-binding \
    --name "${registry_name}" \
    --service-account-name "${account_name}" \
    --role "${role}"
  echo "Role binding"${role}" for account ${account_name}. OK."
done
